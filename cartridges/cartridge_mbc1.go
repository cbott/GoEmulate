package cartridges

import (
	"fmt"
	"log"
)

// Memory Bank Controller 1 Cartridge
// Up to 2MiB ROM (128 banks) / 8KiB RAM (1 bank)
// OR
// Up to 512KiB ROM (32 banks) / 32KiB RAM (4 banks)
type MemoryBankController1Cartridge struct {
	filename string
	rom      []uint8
	ram      [][RAMBankSize]uint8

	// Number of available 16MiB ROM banks we can switch between (2-128)
	numRomBanks uint8
	// Number of available 8MiB RAM banks we can switch between (0-4)
	numRamBanks uint8

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

func NewMBC1Cartridge(filename string, data []uint8) *MemoryBankController1Cartridge {
	c := MemoryBankController1Cartridge{rom: data, filename: filename}
	c.numRomBanks = 1 << (data[ROMSizeAddress] + 1)

	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	c.numRamBanks = uint8(ramSize / 8) // 8KiB per bank
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, c.numRamBanks)

	c.LoadRAM()

	return &c
}

func (c *MemoryBankController1Cartridge) ReadFrom(address uint16) uint8 {
	// Read from ROM Bank 0 (fixed)
	// TODO: support multi-cartridges which can switch this bank in some cases
	if address < ROMBankSize {
		return c.rom[address]
	}

	// Bank 1 is switched
	if address < ROMEndAddress {
		var bank uint8 = c.romBank

		// ROM bank 0 cannot be selected, hardware will use bank 1 instead
		// Note: we intentionally do this before adding bits 4/5 below to match hardware behavior
		if bank == 0 {
			bank = 1
		}

		if !c.ramMode {
			// We are in ROM Banking Mode, use ramBank as bits 4 and 5 of bank number
			bank |= c.ramBank << 5
		}

		// Mask bank to the number of banks available
		bank &= c.numRomBanks - 1

		offset := uint32(bank-1) * ROMBankSize
		return c.rom[uint32(address)+offset]
	}

	// RAM
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

		// Limit access to RAM banks that actually exist
		// TODO: verify correct behavior
		bank %= c.numRamBanks

		// Read the value from the appropriate RAM bank
		return c.ram[bank][address-ExternalRAMStartAddress]
	}

	panic(fmt.Sprintf("Attempted to read from undefined Cartridge address 0x%X", address))
}

func (c *MemoryBankController1Cartridge) WriteTo(address uint16, value uint8) {
	switch address >> 12 {
	case 0, 1:
		// RAM Enable Select (0000-1FFF)
		c.ramEnabled = (value & 0xF) == 0xA
	case 2, 3:
		// ROM Bank Select (2000-3FFF)
		c.romBank = value & 0b11111
	case 4, 5:
		// RAM Bank Select (4000-5FFF)
		// 0-3, select RAM bank or upper 2 bits of ROM bank
		c.ramBank = value & 0b11
	case 6, 7:
		// ROM/RAM Mode Select
		c.ramMode = (value & 1) == 1
	case 0xA, 0xB:
		// Write to RAM (A000-BFFF)
		// Writing to RAM when not enabled does nothing
		if !c.ramEnabled {
			return
		}

		// In ROM banking mode we only have access to one RAM bank
		var bank uint8 = 0
		if c.ramMode {
			bank = c.ramBank
		}

		if c.ramBank < c.numRamBanks {
			// Set the value in the appropriate RAM bank
			c.ram[bank][address-ExternalRAMStartAddress] = value
		}
	default:
		// Our cartridge will ignore writes to invalid addresses
		return
	}
}

// Save cartridge RAM contents to a file
func (c *MemoryBankController1Cartridge) SaveRAM() {
	if c.numRamBanks == 0 {
		log.Printf("Cartridge does not have any RAM banks to save\n")
		return
	}
	// Note: saving is enabled here even if the physical cartridge wouldn't have had the battery to support it
	err := WriteRAMToFile(c.filename, c.ram)
	if err != nil {
		log.Printf("Unable save RAM: %v\n", err)
	}
}

// Load cartridge RAM from a file
// TODO: can we use generics here to not duplicate between mbc1 and mbc3?
func (c *MemoryBankController1Cartridge) LoadRAM() {
	// If cartridge does not have RAM we will skip any sort of loading
	if c.numRamBanks == 0 {
		return
	}

	err := ReadRAMFromFile(c.filename, c.ram)
	if err != nil {
		// We will be permissive here continue running after logging the issue
		log.Printf("Unable to load RAM from file: %v\n", err)
	}
}
