package cartridges

import "fmt"

// Memory Bank Controller 1 Cartridge
// 2MiB ROM / 32KiB RAM
type MemoryBankController1Cartridge struct {
	rom []uint8
	ram [][RAMBankSize]uint8

	// Number of available 16MiB ROM banks we can switch between (2-512)
	numRomBanks uint16

	// Number of available 8MiB RAM banks we can switch between (0-4)
	numRamBanks uint16

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

func NewMBC1Cartridge(data []uint8) *MemoryBankController1Cartridge {
	c := MemoryBankController1Cartridge{rom: data}
	// TODO: reduce duplication with cartridge detection
	// TODO: validate cartridge values match actual file size/headers
	c.numRomBanks = 1 << (data[ROMSizeAddress] + 1)

	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	// 8KiB per bank
	c.numRamBanks = ramSize / 8
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, c.numRamBanks)

	return &c
}

func (c *MemoryBankController1Cartridge) ReadFrom(address uint16) uint8 {
	// Bank 0 is fixed
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
		// TODO: compare to numRamBanks

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
		// TOOD: mask correctly based on actual ROM size (compare to numRomBanks)
		c.romBank = value & 0b11111
	}

	// RAM Bank Select
	if address >= 4000 && address <= 0x5FFF {
		// TODO: I think we actually need to check ramEnabled here and only set RAM bank if enabled
		// otherwise set ROM bank
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

		// TODO: compare to numRamBanks

		// In ROM banking mode we only have access to one RAM bank
		var bank uint8 = 0
		if c.ramMode {
			bank = c.ramBank
		}

		// Set the value in the appropriate RAM bank
		c.ram[bank][address-ExternalRAMStartAddress] = value
	}

	// TODO: handle writes to invalid address? - some other emulators just do nothing
}

// Save cartridge RAM contents to a file
func (c *MemoryBankController1Cartridge) SaveRAM() {

}

// Load cartridge RAM from a file
func (c *MemoryBankController1Cartridge) LoadRAM() {

}
