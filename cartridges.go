package main

import (
	"fmt"
	"io/ioutil"
)

const (
	CartridgeTypeAddress = 0x0147
	ROMSizeAddress       = 0x0148
	RAMSizeAddress       = 0x0149
	TitleAddress         = 0x0134
	TitleLength          = 16
	ROMBankSize          = 0x4000 // 16 KiB
	RAMBankSize          = 0x2000 // 8 KiB
)

var cartridgeTypeMap = map[uint8]string{
	0x00: "ROM Only",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

// RAM Size in KiB
var ramSizeMap = map[uint8]uint32{
	0: 0,
	2: 8,
	3: 32,
	4: 128,
}

// define Cartridge interface
type Cartridge interface {
	ReadFrom(address uint16) uint8
	WriteTo(address uint16, value uint8)
}

// Read a cartridge binary file and return the correct cartridge type containing the file contents
func parseCartridgeFile(filename string) Cartridge {
	// TODO: better error handling

	// Load cartridge binary data
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Unable to load ROM file %s", filename))
	}

	// Parse out cartridge header attributes
	cartridgeType := data[CartridgeTypeAddress]
	cartridgeTypeString, ok := cartridgeTypeMap[cartridgeType]
	if !ok {
		panic(fmt.Sprintf("Unknown cartridge type %d", cartridgeType))
	}

	ramSizeKey := data[RAMSizeAddress]
	ramSize, ok := ramSizeMap[ramSizeKey]
	if !ok {
		panic(fmt.Sprintf("Unknown RAM Size code %d", ramSizeKey))
	}

	var romSize int = 32 * (1 << data[ROMSizeAddress])
	var title string = string(data[TitleAddress : TitleAddress+TitleLength])

	fmt.Printf("Cartridge file: %s\n", filename)
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Type: %s\n", cartridgeTypeString)
	fmt.Printf("ROM Size: %d KiB\n", romSize)
	fmt.Printf("RAM Size: %d KiB\n", ramSize)

	// Return correct cartridge type for this file
	switch cartridgeType {
	case 0x00:
		return &ROMOnlyCartridge{rom: data}
	case 0x01:
		return MakeMBC1Cartridge(data)
	default:
		panic(fmt.Sprintf("Cartridge type %d not implemented", cartridgeType))
	}
}

// Load an initialized Cartdridge struct into Game Boy memory
func (gb *Gameboy) LoadCartridge(c Cartridge) {
	gb.memory.cartridge = c
}

// Basic cartridge containing 32KiB ROM from 0000-7FFF
type ROMOnlyCartridge struct {
	rom []uint8
}

func (c *ROMOnlyCartridge) ReadFrom(address uint16) uint8 {
	return c.rom[address]
}

func (c *ROMOnlyCartridge) WriteTo(address uint16, value uint8) {
	// TODO: do we need to support RAM here too?
	panic("Cannot write to ROM Only Cartridge")
}

// Memory Bank Controller 1 Cartridge
// 2MiB ROM / 32KiB RAM
type MemoryBankController1Cartridge struct {
	rom []uint8

	// Each RAM bank is 8KiB
	ram [][RAMBankSize]uint8

	// currently selected ROM bank for 4000-7FFF
	romBank uint8
	// currently selected RAM bank for A000-BFFF
	// used as most significant 2 bits of ROM bank if ramMode is false
	ramBank uint8

	ramEnabled bool

	// ramMode
	// false -> ROM Banking Mode, up to 8KiB RAM and 2MiB ROM
	// true  -> RAM Banking Mode, up to 32KiB RAM and 512KiB ROM
	ramMode bool
}

func MakeMBC1Cartridge(data []uint8) *MemoryBankController1Cartridge {
	c := MemoryBankController1Cartridge{rom: data}
	// TODO: reduce duplication with cartridge detection
	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	// 8KiB per bank
	numBanks := ramSize / 8
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, numBanks)

	return &c
}

func (c *MemoryBankController1Cartridge) ReadFrom(address uint16) uint8 {
	// Bank 0 is fixed
	if address < ROMBankSize {
		return c.rom[address]
	}

	// Bank 1 is switched
	if address < CartridgeEndAddress {
		var bank uint8 = c.romBank

		// ROM bank 0 cannot be selected, hardware will use bank 1 instead
		// Note: we intentionally do this before adding bits 4/5 below to match hardware behavior
		if bank == 0 {
			bank = 1
		}

		if !c.ramMode {
			// We are in ROM Banking Mode, use ramBank as bits 4 and 5 of bank number
			bank &= c.ramBank << 5
		}

		offset := uint32(bank-1) * ROMBankSize

		return c.rom[uint32(address)+offset]
	}

	// RAM
	// TODO: reduce duplication with Write
	if address >= ExternalRAMStartAddress && address < ExternalRAMEndAddress {
		// Reading from RAM when not enabled is undefined
		if !c.ramEnabled {
			return 0xFF
		}

		// In ROM banking mode we only have access to one RAM bank
		var bank uint8 = 0
		if c.ramMode {
			bank = c.ramBank
		}

		// Read the value from the appropriate RAM bank
		return c.ram[bank][address-ExternalRAMStartAddress]
	}

	panic(fmt.Sprintf("Attempted to read from undefined Cartridge address 0x%X", address))
}

func (c *MemoryBankController1Cartridge) WriteTo(address uint16, value uint8) {
	// RAM Enable Select
	if address <= 0x1FFF {
		c.ramEnabled = (value & 0xF) == 0xA
	}

	// ROM Bank Select
	if address >= 0x2000 && address <= 0x3FFF {
		// TOOD: mask correctly based on actual ROM size
		c.romBank = value & 0b11111
	}

	// RAM Bank Select
	if address >= 4000 && address <= 0x5FFF {
		c.ramBank = value & 0b11
	}

	// ROM/RAM Mode Select
	if address >= 0x6000 && address <= 0x7FFF {
		c.ramMode = (value & 1) == 1
	}

	// RAM
	if address >= ExternalRAMStartAddress && address < ExternalRAMEndAddress {
		// RAM must be enabled before writing to it
		if !c.ramEnabled {
			return
		}

		// TODO: Addressing other RAM banks should only happen if they exist, do nothing otherwise

		// In ROM banking mode we only have access to one RAM bank
		var bank uint8 = 0
		if c.ramMode {
			bank = c.ramBank
		}

		// Set the value in the appropriate RAM bank
		c.ram[bank][address-ExternalRAMStartAddress] = value
	}

	// TODO: handle writes to invalid address?
}
