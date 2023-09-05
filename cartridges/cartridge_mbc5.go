package cartridges

import (
	"fmt"
	"log"
)

// Memory Bank Controller 5 Cartridge
// Up to 8MiB ROM (512 banks) / 128KiB RAM (16 banks)
// MBC5 cartridges can support rumble but we will be ignoring that for now
type MemoryBankController5Cartridge struct {
	filename string
	rom      []uint8
	ram      [][RAMBankSize]uint8

	// Number of available 16MiB ROM banks we can switch between (2-512)
	numRomBanks uint16
	// Number of available 8MiB RAM banks we can switch between (0-4)
	numRamBanks uint8

	// Currently selected ROM bank for 4000-7FFF (0-512)
	romBank uint16
	// Currently selected RAM bank A000-BFFF (0-15)
	ramBank uint8

	// Whether RAM reading and writing are enabled
	ramEnabled bool
}

func NewMBC5Cartridge(filename string, data []uint8) *MemoryBankController5Cartridge {
	c := MemoryBankController5Cartridge{rom: data, filename: filename}
	c.numRomBanks = 1 << (data[ROMSizeAddress] + 1)

	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	c.numRamBanks = uint8(ramSize / 8) // 8KiB per bank
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, c.numRamBanks)

	// Load RAM state
	c.LoadRAM()

	return &c
}

// Read a value from MBC3 ROM or RAM
func (c *MemoryBankController5Cartridge) ReadFrom(address uint16) uint8 {
	// Read from ROM Bank 0 (fixed)
	if address < ROMBankSize {
		return c.rom[address]
	}

	// Read from ROM Bank 1 (switched)
	if address < ROMEndAddress {
		var bank uint16 = c.romBank
		// If bank 0 is selected we will have a negative offset allowing to read bank 0 from bank 1 address space
		offset := (int32(bank) - 1) * ROMBankSize
		return c.rom[int32(address)+offset]
	}

	// Read from RAM
	if address >= ExternalRAMStartAddress && address < ExternalRAMEndAddress {
		// Reading from RAM when not enabled is undefined
		if !c.ramEnabled {
			return 0xFF
		}

		// We have selected a RAM bank to be active
		if c.ramBank < c.numRamBanks {
			// Read from selcted RAM bank
			return c.ram[c.ramBank][address-ExternalRAMStartAddress]
		}

		panic(fmt.Sprintf("Attempted to read from invalid RAM bank 0x%X", c.ramBank))
	}

	// If reading from invalid address we will panic for debugging but hardware may behave differently
	panic(fmt.Sprintf("Attempted to read from undefined Cartridge address 0x%X", address))
}

// Write a value to MBC3 control registers or RAM
func (c *MemoryBankController5Cartridge) WriteTo(address uint16, value uint8) {
	switch address >> 12 {
	case 0, 1:
		// RAM Enable Select (0000-1FFF)
		c.ramEnabled = (value & 0xF) == 0xA
	case 2:
		// ROM Bank lower 8 bits (2000-3FFF)
		c.romBank = (c.romBank & 0x100) | uint16(value)
	case 3:
		// ROM Bank bit 8 (2000-3FFF)
		c.romBank = (c.romBank & 0x0FF) | (uint16(value&1) << 8)
	case 4, 5:
		// RAM Bank Select (4000-5FFF)
		// Select RAM bank 0-15
		// Mask to lower 4 bits only, though this should never matter, unclear what proper handling is
		// TODO: if cartidge supports rumble bit 3 sets rumble state instead
		c.ramBank = value & 0xF
	case 0xA, 0xB:
		// Write to RAM (A000-BFFF)
		// Writing to RAM when not enabled does nothing
		if !c.ramEnabled {
			return
		}

		if c.ramBank < c.numRamBanks {
			// We have selected a RAM bank to be active
			// Set the value in the appropriate RAM bank
			c.ram[c.ramBank][address-ExternalRAMStartAddress] = value
		}
	default:
		// Our cartridge will ignore writes to invalid addresses
		return
	}
}

// Save cartridge RAM contents to a file
func (c *MemoryBankController5Cartridge) SaveRAM() {
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
func (c *MemoryBankController5Cartridge) LoadRAM() {
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
