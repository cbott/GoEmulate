package main

/*
0000              8000        A000           C000       E000         FE00      FEA0     FF00    FF80
| ROM (cartridge) | Video RAM | External RAM | RAM      | Unused     | OAM RAM | Unused | I/O   | HRAM
*/

const (
	CartridgeEndAddress     = 0x8000
	ExternalRAMStartAddress = 0xA000
	ExternalRAMEndAddress   = 0xC000
	// Memory address of flag indicating whether to load from Boot ROM or cartrige address space
	BOOT = 0xFF50
)

type Memory struct {
	memory    [0x10000]uint8
	cartridge Cartridge
	apu       *APU
	bootrom   [0x100]uint8

	// Internal counter to keep track of when the DIV register should increment
	divAccumulator int
	// Stores the state of each joypad button (down/up/left/right/start/select/B/A)
	ButtonStates uint8
}

// Write a value to memory
func (m *Memory) set(address uint16, value uint8) {
	if address == DIV {
		// Writing any value to the DIV register sets it to 0
		m.divAccumulator = 0
		m.memory[address] = 0
	} else if address == DMA {
		// Initiate a DMA transfer, register should never be read so we will leave it at 0
		m.performDMATransfer(value)
	} else if address == JOYPAD {
		// Only bits 4 and 5 of the P1 register are writeable
		m.memory[address] = (m.memory[address] & 0xF) | (value & 0b00110000)
	} else if address < CartridgeEndAddress ||
		address >= ExternalRAMStartAddress && address < ExternalRAMEndAddress {
		m.cartridge.WriteTo(address, value)
	} else if address >= NR10 && address < WaveRAMStart {
		// Sound controls
		m.apu.WriteTo(address, value)
	} else {
		m.memory[address] = value
	}
}

// Read a value from memory
func (m *Memory) get(address uint16) uint8 {
	// Address space 0-FF is mapped to Boot ROM untill fully booted
	if (address < 0x100) && (m.memory[BOOT] == 0) {
		return m.bootrom[address]
	}

	// Address space 0-8000 is mapped to the cartridge ROM
	// Address space A000-C000 is mapped to cartridge RAM
	if address < CartridgeEndAddress {
		if m.cartridge == nil {
			panic("Attempted to access cartridge before one is loaded")
		}
		return m.cartridge.ReadFrom(address)
	}

	if address >= ExternalRAMStartAddress && address < ExternalRAMEndAddress {
		return m.cartridge.ReadFrom(address)
	}

	// Top 3 bits of IF register are unused and always read high
	if address == IF {
		return m.memory[address] | 0b11100000
	}

	// STAT bit 7 is unused
	if address == STAT {
		return m.memory[address] | 0b10000000
	}

	// Joypad return value depends on value in the register
	if address == JOYPAD {
		return m.GetP1Value()
	}

	// Sound register access
	if address >= NR10 && address < WaveRAMStart {
		return m.apu.ReadFrom(address)
	}

	// In most cases we just read a raw value
	return m.memory[address]

	// TODO: (future) Implement E000 as echo RAM, prevent access to illegal addresses
	// or access restricted memory during DMA or specific PPU modes
}

func (m *Memory) Init() {
	m.bootrom = BootRom
	// TODO: Implementing the APU like this feels hacky
	m.apu = NewAPU()
	// TODO: this is kind of a funky way to do things?
	m.apu.channel3.waveRAM = m.memory[WaveRAMStart:]
}

// Perform a DMA transfer into OAM RAM from the specified source address (divided by 0x100)
func (m *Memory) performDMATransfer(source uint8) {
	// We are performing this all at once for convenience
	source_address := uint16(source) << 8

	var index uint16 = 0
	for index = 0; index <= 0x9F; index++ {
		m.set(OAMRamAddressStart+index, m.get(source_address+index))
	}
}

// Set memory to the state it would be in after boot ROM runs
// if skipping normal bootrom execution we can run this instead
func (m *Memory) BypassBootROM() {
	m.apu.BypassBootROM()
	m.memory[BOOT] = 1
	m.memory[DIV] = 0x1E // Does not currently match my timer implementation
	m.memory[IF] = 0xE1
	m.memory[TIMA] = 0x00
	m.memory[TMA] = 0x00
	m.memory[TAC] = 0x00
	m.memory[LCDC] = 0x91
	m.memory[SCY] = 0x00
	m.memory[SCX] = 0x00
	m.memory[LYC] = 0x00
	m.memory[BGP] = 0xFC
	m.memory[WY] = 0x00
	m.memory[WX] = 0x00
	m.memory[IE] = 0x00
}
