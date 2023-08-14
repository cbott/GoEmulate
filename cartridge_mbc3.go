package main

import "fmt"

// Real Time Clock registers
// Address	Value				Range
// 08		Seconds   			0-59
// 09		Minutes   			0-59
// 0A		Hours     			0-23
// 0B		Day Counter Low		0x00-0xFF
// 0C		Day Counter High	Bit 0 = Day Counter MSB
//                              Bit 6 = Halt (0=Active, 1=Stop Timer)
//                              Bit 7 = Day Counter Carry Bit (1=Counter Overflow)

const (
	RTCBankStart = 0x08
	NumRTCBanks  = 5
)

// Memory Bank Controller 3 Cartridge
// 2MiB ROM / 32KiB RAM, Timer
type MemoryBankController3Cartridge struct {
	rom        []uint8
	ram        [][RAMBankSize]uint8
	rtc        [NumRTCBanks]uint8
	latchedrtc [NumRTCBanks]uint8

	// Number of available 16MiB ROM banks we can switch between (2-512)
	numRomBanks uint16
	// Number of available 8MiB RAM banks we can switch between (0-4)
	numRamBanks uint8

	// Currently selected ROM bank for 4000-7FFF
	romBank uint8
	// Currently selected RAM bank or RTC register for A000-BFFF
	ramBank uint8

	// Whether RAM/RTC reading and writing are enabled
	ramEnabled bool
	// Whether this cartridge supports a real time clock
	hasRTC bool
}

func NewMBC3Cartridge(data []uint8) *MemoryBankController3Cartridge {
	c := MemoryBankController3Cartridge{rom: data}
	c.numRomBanks = 1 << (data[ROMSizeAddress] + 1)

	ramSizeKey := data[RAMSizeAddress]
	ramSize := ramSizeMap[ramSizeKey]
	c.numRamBanks = uint8(ramSize / 8) // 8KiB per bank
	// Initialize RAM banks
	c.ram = make([][RAMBankSize]uint8, c.numRamBanks)

	// Cartridge types 0x0F and 0x10 have RTC hardware
	cartridgeType := data[CartridgeTypeAddress]
	if cartridgeType == 0x0F || cartridgeType == 0x10 {
		c.hasRTC = true
	}

	return &c
}

// Read a value from MBC3 ROM or RAM
func (c *MemoryBankController3Cartridge) ReadFrom(address uint16) uint8 {
	// Read from ROM Bank 0 (fixed)
	if address < ROMBankSize {
		return c.rom[address]
	}

	// Read from ROM Bank 1 (switched)
	if address < CartridgeEndAddress {
		var bank uint8 = c.romBank

		// ROM bank 0 cannot be selected, hardware will use bank 1 instead
		if bank == 0 {
			bank = 1
		}

		offset := uint32(bank-1) * ROMBankSize
		return c.rom[uint32(address)+offset]
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

		// We have selected a RTC register to be active
		if c.hasRTC && c.ramBank >= RTCBankStart && c.ramBank < RTCBankStart+NumRTCBanks {
			return c.rtc[c.ramBank-RTCBankStart]
		}

		panic(fmt.Sprintf("Attempted to read from invalid RAM bank 0x%X", c.ramBank))
	}

	// If reading from invalid address we will panic for debugging but hardware may behave differently
	panic(fmt.Sprintf("Attempted to read from undefined Cartridge address 0x%X", address))
}

// Write a value to MBC3 control registers or RAM
func (c *MemoryBankController3Cartridge) WriteTo(address uint16, value uint8) {
	switch address >> 12 {
	case 0, 1:
		// RAM Enable Select (0000-1FFF)
		c.ramEnabled = (value & 0xF) == 0xA
	case 2, 3:
		// ROM Bank Select (2000-3FFF)
		c.romBank = value & 0b1111111
	case 4, 5:
		// RAM Bank Select (4000-5FFF)
		// Can set 0-3 to select RAM bank or 8-C to select a RTC register
		// Mask to lower 4 bits only, though this should never matter, unclear what proper handling is
		c.ramBank = value & 0xF
	case 6, 7:
		// Latch clock data (6000-7FFF)
		// The proper method for latching the clock is actually to write 0x00 followed by 0x01
		// but we will be more permissive here, some sources suggest writing 0x01 is all you need
		if value == 0x01 {
			for i := 0; i < NumRTCBanks; i++ {
				c.latchedrtc[i] = c.rtc[i]
			}
		}
	case 0xA, 0xB:
		// Write to RAM or RTC Register (A000-BFFF)
		// Writing to RAM/RTC when not enabled does nothing
		if !c.ramEnabled {
			return
		}

		if c.ramBank < c.numRamBanks {
			// We have selected a RAM bank to be active
			// Set the value in the appropriate RAM bank
			c.ram[c.ramBank][address-ExternalRAMStartAddress] = value
		} else if c.hasRTC && c.ramBank >= RTCBankStart && c.ramBank < RTCBankStart+NumRTCBanks {
			// We have selected a RTC register to be active
			// Write the value in the RTC register
			c.rtc[c.ramBank-RTCBankStart] = value
		}
	default:
		// Our cartridge will ignore writes to invalid addresses
		return
	}
}
