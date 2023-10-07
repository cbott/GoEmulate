package cartridges

import (
	"fmt"
	"log"
)

// Memory Bank Controller 5 Cartridge
// Up to 8MiB ROM (512 banks) / 128KiB RAM (16 banks), optional rumble
type MemoryBankController5Cartridge struct {
	CartridgeCore
	// Whether this cartridge supports rumble
	hasRumble bool
}

func NewMBC5Cartridge(filename string, data []uint8) *MemoryBankController5Cartridge {
	c := MemoryBankController5Cartridge{}
	c.rom = data
	c.filename = filename
	c.numRomBanks = 1 << (data[ROMSizeAddress] + 1)

	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	c.numRamBanks = uint8(ramSize / 8) // 8KiB per bank
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, c.numRamBanks)

	// Cartridge types 0x1C/0x1D/0x1E have rumble motor
	cartridgeType := data[CartridgeTypeAddress]
	if cartridgeType == 0x1C || cartridgeType == 0x1D || cartridgeType == 0x1E {
		c.hasRumble = true
	}

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
		if c.hasRumble {
			// Lower 3 bits set RAM bank, bit 3 controls rumble motor (ignored by emulator)
			c.ramBank = value & 0x7
		} else {
			// Lower 4 bits set RAM bank
			c.ramBank = value & 0xF
		}
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
