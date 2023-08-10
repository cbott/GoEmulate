package main

import "fmt"

// TODO: merge functionality with MBC1?

// Real Time Clock registers
// Address	Value				Range
// 08		Seconds   			0-59
// 09		Minutes   			0-59
// 0A		Hours     			0-23
// 0B		Day Counter Low		0x00-0xFF
// 0C		Day Counter High	Bit 0 = Day Counter MSB
//                              Bit 6 = Halt (0=Active, 1=Stop Timer)
//                              Bit 7 = Day Counter Carry Bit (1=Counter Overflow)

// Memory Bank Controller 3 Cartridge
// 2MiB ROM / 32KiB RAM, Timer
type MemoryBankController3Cartridge struct {
	rom []uint8
	ram [][RAMBankSize]uint8

	// Number of available 16MiB ROM banks we can switch between (2-512)
	numRomBanks uint16

	// Number of available 8MiB RAM banks we can switch between (0-4)
	numRamBanks uint8

	// currently selected ROM bank for 4000-7FFF
	romBank uint8
	// currently selected RAM bank for A000-BFFF
	// used as most significant 2 bits of ROM bank if ramMode is false
	ramBank uint8

	ramEnabled bool
}

func MakeMBC3Cartridge(data []uint8) *MemoryBankController3Cartridge {
	c := MemoryBankController3Cartridge{rom: data}
	// TODO: reduce duplication with cartridge detection
	// TODO: validate cartridge values match actual file size/headers
	c.numRomBanks = 1 << (data[ROMSizeAddress] + 1)

	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	// 8KiB per bank
	c.numRamBanks = uint8(ramSize / 8)
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, c.numRamBanks)

	return &c
}

func (c *MemoryBankController3Cartridge) ReadFrom(address uint16) uint8 {
	// Bank 0 is fixed
	if address < ROMBankSize {
		return c.rom[address]
	}

	// Bank 1 is switched
	if address < CartridgeEndAddress {
		var bank uint8 = c.romBank

		// ROM bank 0 cannot be selected, hardware will use bank 1 instead
		if bank == 0 {
			bank = 1
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

		// We have selected a RAM bank to be active
		if c.ramBank < 8 {
			if c.ramBank >= c.numRamBanks {
				// TODO: define correct behavior
				panic("Attempted to read from RAM bank outside of range")
			}
			// Set the value in the appropriate RAM bank
			return c.ram[c.ramBank][address-ExternalRAMStartAddress]
		}

		// We have selected a RTC register to be active
		if c.ramBank < 0xD {
			// TODO: implement RTC
			return 0x00
		}

		// TODO: define correct behavior
		panic("Attempted to read from invalid RAM bank")
	}

	panic(fmt.Sprintf("Attempted to read from undefined Cartridge address 0x%X", address))
}

func (c *MemoryBankController3Cartridge) WriteTo(address uint16, value uint8) {
	// RAM Enable Select
	if address <= 0x1FFF {
		c.ramEnabled = (value & 0xF) == 0xA
	}

	// ROM Bank Select
	if address >= 0x2000 && address <= 0x3FFF {
		// TOOD: mask correctly based on actual ROM size (compare to numRomBanks)
		c.romBank = value & 0b1111111
	}

	// RAM Bank Select
	if address >= 4000 && address <= 0x5FFF {
		// Can set 0-3 to select RAM bank or 8-C to instead read the RTC registers
		c.ramBank = value & 0xF
	}

	// Latch Clock Data
	// if address >= 0x6000 && address <= 0x7FFF {
	// }

	// RAM
	if address >= ExternalRAMStartAddress && address < ExternalRAMEndAddress {
		// RAM must be enabled before writing to it
		if !c.ramEnabled {
			return
		}

		// We have selected a RAM bank to be active
		if c.ramBank < 8 {
			if c.ramBank >= c.numRamBanks {
				// TODO: check whether this is correct behavior
				// for now we will ignore writes to invalid RAM banks
				return
			}
			// Set the value in the appropriate RAM bank
			c.ram[c.ramBank][address-ExternalRAMStartAddress] = value
		}

		// We have selected a RTC register to be active
		// if c.ramBank < 0xD {
		// TODO: implement RTC
		// }

		// If a RAM bank is selected which is not valid we will ignore it
	}

	// TODO: handle writes to invalid address? - some other emulators just do nothing
}
