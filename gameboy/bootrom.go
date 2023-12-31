package gameboy

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
