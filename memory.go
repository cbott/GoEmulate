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

// Game Boy boot ROM program
var BootRom = [0x100]uint8{
	// Initialize RAM
	// Writes 0 to all of video RAM
	0x31, 0xFE, 0xFF, // LD SP $FFFE
	0xAF,             // XOR A
	0x21, 0xFF, 0x9F, // LD HL $9FFF
	// 7
	0x32,       // LD (HL-),A
	0xCB, 0x7C, // BIT 7,H
	0x20, 0xFB, // JR NZ,$FB -- jump to 7

	// Initialize sound
	0x21, 0x26, 0xFF, // LD HL $FF26
	0x0E, 0x11, // LD C $11
	0x3E, 0x80, // LD A $80
	0x32,       // LD (HL-),A
	0xE2,       // LD ($FF00+C),A
	0x0C,       // INC C
	0x3E, 0xF3, // LD A $F3
	0xE2,       // LD ($FF00+C),A
	0x32,       // LD (HL-),A
	0x3E, 0x77, // LD A $77
	0x77, // LD (HL),A

	// Set up logo
	0x3E, 0xFC, // LD A $FC
	0xE0, 0x47, // LD ($FF47),A	-- Initialize BGP to 0xFC
	0x11, 0x04, 0x01, // LD DE $0104
	0x21, 0x10, 0x80, // LD HL $8010
	// 27
	0x1A,             // LD A,(DE)
	0xCD, 0x95, 0x00, // CALL $0095
	0xCD, 0x96, 0x00, // CALL $0096
	0x13,       // INC DE
	0x7B,       // LD A,E
	0xFE, 0x34, // CP A,$34
	0x20, 0xF3, // JR NZ,$F3 -- jump to 27
	0x11, 0xD8, 0x00, // LD DE $00D8
	0x06, 0x08, // LD B $08
	// 39
	0x1A,       // LD A,(DE)
	0x13,       // INC DE
	0x22,       // LD (HL+),A
	0x23,       // INC HL
	0x05,       // DEC B
	0x20, 0xF9, // JR NZ,$F9 -- jump to 39
	0x3E, 0x19, // LD A $19
	0xEA, 0x10, 0x99, // LD ($9910),A
	0x21, 0x2F, 0x99, // LD HL $992F
	// 48
	0x0E, 0x0C, // LD C $0C
	// 4A
	0x3D,       // DEC A
	0x28, 0x08, // JR Z,$08 -- jump to 55
	0x32,       // LD (HL-),A
	0x0D,       // DEC C
	0x20, 0xF9, // JR NZ,$F9 -- jump to 4A
	0x2E, 0x0F, // LD L $0F
	0x18, 0xF3, // JR $F3 -- jump to 48

	// Scroll Logo
	// 55
	0x67,       // LD H,A
	0x3E, 0x64, // LD A $64
	0x57,       // LD D,A
	0xE0, 0x42, // LD ($FF42),A -- Set SCY to 0x64
	0x3E, 0x91, // LD A $91
	0xE0, 0x40, // LD ($FF40),A -- Set LCDC to 0x91
	0x04, // INC B
	// 60
	0x1E, 0x02, // LD E $02
	// 62
	0x0E, 0x0C, // LD C $0C
	// 64
	0xF0, 0x44, // LD A,($FF44)
	0xFE, 0x90, // CP A,$90
	0x20, 0xFA, // JR NZ,$FA -- jump to 64
	0x0D,       // DEC C
	0x20, 0xF7, // JR NZ,$F7 -- jump to 64
	0x1D,       // DEC E
	0x20, 0xF2, // JR NZ,$F2 -- jump to 62
	0x0E, 0x13, // LD C $13
	0x24,       // INC H
	0x7C,       // LD A,H
	0x1E, 0x83, // LD E $83
	0xFE, 0x62, // CP A,$62
	0x28, 0x06, // JR Z,$06 -- jump to 80
	0x1E, 0xC1, // LD E $C1
	0xFE, 0x64, // CP A,$64
	0x20, 0x06, // JR NZ,$06 -- jump to 86

	// Play sound
	// 80
	0x7B,       // LD A,E
	0xE2,       // LD ($FF00+C),A
	0x0C,       // INC C
	0x3E, 0x87, // LD A $87
	0xE2, // LD ($FF00+C),A

	// Scroll logo
	// 86
	0xF0, 0x42, // LD A,($FF42) -- Store SCY to A
	0x90,       // SUB B -- Subtract B from A
	0xE0, 0x42, // LD ($FF42),A -- Write the value of A to SCY
	0x15,       // DEC D
	0x20, 0xD2, // JR NZ,$D2 -- jump to 60
	0x05,       // DEC B
	0x20, 0x4F, // JR NZ,$4F -- jump to E0
	0x16, 0x20, // LD D $20
	0x18, 0xCB, // JR $CB -- jump to 60

	// Decode logo
	0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
	0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9,

	// Logo data
	// A8
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
	0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
	0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	// D8
	0x3C, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,

	// Compare logo
	// E0
	0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x00, 0x00, 0x23, 0x7D, 0xFE, 0x34, 0x20,
	0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x00, 0x00, 0x3E, 0x01, 0xE0, 0x50,
}

type Memory struct {
	memory    [0x10000]uint8
	cartridge Cartridge
	bootrom   [0x100]uint8

	// TODO: move div accumulator to GB object if we keep the reference to the gb object?
	// Internal counter to keep track of when the DIV register should increment
	divAccumulator int

	// TODO: see if we can avoid this circular reference
	gameboy *Gameboy
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
	} else if (address >= 0xFF10 && address <= 0xFF26) || (address >= 0xFF30 && address <= 0xFF3F) {
		// TODO: temporary - skip sound things
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

	if (address >= 0xFF10 && address <= 0xFF26) || (address >= 0xFF30 && address <= 0xFF3F) {
		// TODO: temporary, disable sound stuff
		return 0x00
	}

	// Top 3 bits of IF register are unused and always read high
	// TODO: confirm if that's true or if they just default high on startup
	if address == IF {
		return m.memory[address] | 0b11100000
	}

	// STAT bit 7 is unused
	if address == STAT {
		return m.memory[address] | 0b10000000
	}

	// Joypad return value depends on value in the register
	if address == JOYPAD {
		return m.gameboy.joypad.GetP1Value(m.memory[address])
	}

	// In most cases we just read a raw value
	return m.memory[address]

	// TODO: Implement E000 as echo RAM?
}

func (m *Memory) Init() {
	m.bootrom = BootRom
}

// Perform a DMA transfer into OAM RAM from the specified source address (divided by 0x100)
func (m *Memory) performDMATransfer(source uint8) {
	// TODO: We are performing this all at once for convenience
	// should implement checks to ensure we don't access restricted memory during this time
	// maybe make a flag that gets checked in get/set and clears after a certain number of cycles
	source_address := uint16(source) << 8

	var index uint16 = 0
	for index = 0; index <= 0x9F; index++ {
		m.set(OAMRamAddressStart+index, m.get(source_address+index))
	}
}

// Set memory to the state it would be in after boot ROM runs
// if skipping normal bootrom execution we can run this instead
func (m *Memory) BypassBootROM() {
	// TODO: figure out what these constants mean
	// switch to setting by register name like m.set(SCY, 0)
	// TODO: compare these to what we get when we actually run the boot ROM
	m.memory[DIV] = 0x1E // Needs a double check
	m.set(IF, 0xE1)      // Conflicting sources on this one
	m.set(0xFF05, 0x00)  // TIMA
	m.set(0xFF06, 0x00)  // TMA
	m.set(0xFF07, 0x00)  // TAC
	m.set(0xFF10, 0x80)  // NR10
	m.set(0xFF11, 0xBF)  // NR11
	m.set(0xFF12, 0xF3)  // NR12
	m.set(0xFF14, 0xBF)  // NR14
	m.set(0xFF16, 0x3F)  // NR21
	m.set(0xFF17, 0x00)  // NR22
	m.set(0xFF19, 0xBF)  // NR24
	m.set(0xFF1A, 0x7F)  // NR30
	m.set(0xFF1B, 0xFF)  // NR31
	m.set(0xFF1C, 0x9F)  // NR32
	m.set(0xFF1E, 0xBF)  // NR33
	m.set(0xFF20, 0xFF)  // NR41
	m.set(0xFF21, 0x00)  // NR42
	m.set(0xFF22, 0x00)  // NR43
	m.set(0xFF23, 0xBF)  // NR30
	m.set(0xFF24, 0x77)  // NR50
	m.set(0xFF25, 0xF3)  // NR51
	m.set(0xFF26, 0xF1)  // NR52
	m.set(0xFF40, 0x91)  // LCDC
	m.set(0xFF42, 0x00)  // SCY
	m.set(0xFF43, 0x00)  // SCX
	m.set(0xFF45, 0x00)  // LYC
	m.set(0xFF47, 0xFC)  // BGP
	m.set(0xFF4A, 0x00)  // WY
	m.set(0xFF4B, 0x00)  // WX
	m.set(0xFFFF, 0x00)  // Interrupt Enable register
}
