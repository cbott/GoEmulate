package main

import "fmt"

// Opcodes data from https://pastraiser.com/cpu/gameboy/gameboy_opcodes.html
// and http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf

// Read the 16 bit address indicated by (PC, PC+1) and increment PC twice
func (gb *Gameboy) popPC16() uint16 {
	// TODO: Not sure if this is common enough to need a helper
	// and also a bit messy to make it a method of Gameboy? Probably okay.
	// TODO: move to same place as popPC, and potentially rename since it's different action
	lsb := uint16(gb.popPC())
	msb := uint16(gb.popPC())
	return msb<<8 | lsb
}

// Execute a single opcode and return the number of CPU cycles it took (1MHz CPU cycles)
func (gb *Gameboy) Opcode(opcode uint8) int {
	// TODO: standardize which type of cycle we're talking about
	switch opcode {
	//////////////// 8-bit loads ////////////////
	case 0x3E:
		// LD A,n
		gb.cpu.set_register("A", gb.popPC())
		return 2
	case 0x06:
		// LD B,n
		gb.cpu.set_register("B", gb.popPC())
		return 2
	case 0x0E:
		// LD C,n
		gb.cpu.set_register("C", gb.popPC())
		return 2
	case 0x16:
		// LD D,n
		gb.cpu.set_register("D", gb.popPC())
		return 2
	case 0x1E:
		// LD E,n
		gb.cpu.set_register("E", gb.popPC())
		return 2
	case 0x26:
		// LD H,n
		gb.cpu.set_register("H", gb.popPC())
		return 2
	case 0x2E:
		// LD L,n
		gb.cpu.set_register("L", gb.popPC())
		return 2
	case 0x7F, 0x40, 0x49, 0x52, 0x5B, 0x64, 0x6D:
		// LD X,X (For registers A, B, C, D, E, H, L)
		// Equivalent to NOP
		return 1
	case 0x78:
		// LD A,B
		gb.cpu.set_register("A", gb.cpu.get_register("B"))
		return 1
	case 0x79:
		// LD A,C
		gb.cpu.set_register("A", gb.cpu.get_register("C"))
		return 1
	case 0x7A:
		// LD A,D
		gb.cpu.set_register("A", gb.cpu.get_register("D"))
		return 1
	case 0x7B:
		// LD A,E
		gb.cpu.set_register("A", gb.cpu.get_register("E"))
		return 1
	case 0x7C:
		// LD A,H
		gb.cpu.set_register("A", gb.cpu.get_register("H"))
		return 1
	case 0x7D:
		// LD A,L
		gb.cpu.set_register("A", gb.cpu.get_register("L"))
		return 1
	case 0x7E:
		// LD A,(HL)
		gb.cpu.set_register("A", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x41:
		// LD B,C
		gb.cpu.set_register("B", gb.cpu.get_register("C"))
		return 1
	case 0x42:
		// LD B,D
		gb.cpu.set_register("B", gb.cpu.get_register("D"))
		return 1
	case 0x43:
		// LD B,E
		gb.cpu.set_register("B", gb.cpu.get_register("E"))
		return 1
	case 0x44:
		// LD B,H
		gb.cpu.set_register("B", gb.cpu.get_register("H"))
		return 1
	case 0x45:
		// LD B,L
		gb.cpu.set_register("B", gb.cpu.get_register("L"))
		return 1
	case 0x46:
		// LD B,(HL)
		gb.cpu.set_register("B", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x47:
		// LD B,A
		gb.cpu.set_register("B", gb.cpu.get_register("A"))
		return 1
	case 0x48:
		// LD C,B
		gb.cpu.set_register("C", gb.cpu.get_register("B"))
		return 1
	case 0x4A:
		// LD C,D
		gb.cpu.set_register("C", gb.cpu.get_register("D"))
		return 1
	case 0x4B:
		// LD C,E
		gb.cpu.set_register("C", gb.cpu.get_register("E"))
		return 1
	case 0x4C:
		// LD C,H
		gb.cpu.set_register("C", gb.cpu.get_register("H"))
		return 1
	case 0x4D:
		// LD C,L
		gb.cpu.set_register("C", gb.cpu.get_register("L"))
		return 1
	case 0x4E:
		// LD C,(HL)
		gb.cpu.set_register("C", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x4F:
		// LD C,A
		gb.cpu.set_register("C", gb.cpu.get_register("A"))
		return 1
	case 0x50:
		// LD D,B
		gb.cpu.set_register("D", gb.cpu.get_register("B"))
		return 1
	case 0x51:
		// LD D,C
		gb.cpu.set_register("D", gb.cpu.get_register("C"))
		return 1
	case 0x53:
		// LD D,E
		gb.cpu.set_register("D", gb.cpu.get_register("E"))
		return 1
	case 0x54:
		// LD D,H
		gb.cpu.set_register("D", gb.cpu.get_register("H"))
		return 1
	case 0x55:
		// LD D,L
		gb.cpu.set_register("D", gb.cpu.get_register("L"))
		return 1
	case 0x56:
		// LD D,(HL)
		gb.cpu.set_register("D", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x57:
		// LD D,A
		gb.cpu.set_register("D", gb.cpu.get_register("A"))
		return 1
	case 0x58:
		// LD E,B
		gb.cpu.set_register("E", gb.cpu.get_register("B"))
		return 1
	case 0x59:
		// LD E,C
		gb.cpu.set_register("E", gb.cpu.get_register("C"))
		return 1
	case 0x5A:
		// LD E,D
		gb.cpu.set_register("E", gb.cpu.get_register("D"))
		return 1
	case 0x5C:
		// LD E,H
		gb.cpu.set_register("E", gb.cpu.get_register("H"))
		return 1
	case 0x5D:
		// LD E,L
		gb.cpu.set_register("E", gb.cpu.get_register("L"))
		return 1
	case 0x5E:
		// LD E,(HL)
		gb.cpu.set_register("E", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x5F:
		// LD E,A
		gb.cpu.set_register("E", gb.cpu.get_register("A"))
		return 1
	case 0x60:
		// LD H,B
		gb.cpu.set_register("H", gb.cpu.get_register("B"))
		return 1
	case 0x61:
		// LD H,C
		gb.cpu.set_register("H", gb.cpu.get_register("C"))
		return 1
	case 0x62:
		// LD H,D
		gb.cpu.set_register("H", gb.cpu.get_register("D"))
		return 1
	case 0x63:
		// LD H,E
		gb.cpu.set_register("H", gb.cpu.get_register("E"))
		return 1
	case 0x65:
		// LD H,L
		gb.cpu.set_register("H", gb.cpu.get_register("L"))
		return 1
	case 0x66:
		// LD H,(HL)
		gb.cpu.set_register("H", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x67:
		// LD H,A
		gb.cpu.set_register("H", gb.cpu.get_register("A"))
		return 1
	case 0x68:
		// LD L,B
		gb.cpu.set_register("L", gb.cpu.get_register("B"))
		return 1
	case 0x69:
		// LD L,C
		gb.cpu.set_register("L", gb.cpu.get_register("C"))
		return 1
	case 0x6A:
		// LD L,D
		gb.cpu.set_register("L", gb.cpu.get_register("D"))
		return 1
	case 0x6B:
		// LD L,E
		gb.cpu.set_register("L", gb.cpu.get_register("E"))
		return 1
	case 0x6C:
		// LD L,H
		gb.cpu.set_register("L", gb.cpu.get_register("H"))
		return 1
	case 0x6E:
		// LD L,(HL)
		gb.cpu.set_register("L", gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0x6F:
		// LD L,A
		gb.cpu.set_register("L", gb.cpu.get_register("A"))
		return 1
	case 0x70:
		// LD (HL),B
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("B"))
		return 2
	case 0x71:
		// LD (HL),C
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("C"))
		return 2
	case 0x72:
		// LD (HL),D
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("D"))
		return 2
	case 0x73:
		// LD (HL),E
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("E"))
		return 2
	case 0x74:
		// LD (HL),H
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("H"))
		return 2
	case 0x75:
		// LD (HL),L
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("L"))
		return 2
	case 0x36:
		// LD (HL),n
		gb.memory.set(gb.cpu.get_register16("HL"), gb.popPC())
		return 3
	case 0x0A:
		// LD A,(BC)
		gb.cpu.set_register("A", gb.memory.get(gb.cpu.get_register16("BC")))
		return 2
	case 0x1A:
		// LD A,(DE)
		gb.cpu.set_register("A", gb.memory.get(gb.cpu.get_register16("DE")))
		return 2
	case 0xFA:
		// LD A,(nn)
		gb.cpu.set_register("A", gb.memory.get(gb.popPC16()))
		return 4
	case 0x02:
		// LD (BC),A
		gb.memory.set(gb.cpu.get_register16("BC"), gb.cpu.get_register("A"))
		return 2
	case 0x12:
		// LD (DE),A
		gb.memory.set(gb.cpu.get_register16("DE"), gb.cpu.get_register("A"))
		return 2
	case 0x77:
		// LD (HL),A
		gb.memory.set(gb.cpu.get_register16("HL"), gb.cpu.get_register("A"))
		return 2
	case 0xEA:
		// LD (nn),A
		gb.memory.set(gb.popPC16(), gb.cpu.get_register("A"))
		return 4
	case 0xF2:
		// LD A,(0xFF00+C)
		gb.cpu.set_register("A", gb.memory.get(0xFF00+uint16(gb.cpu.get_register("C"))))
		return 2
	case 0xE2:
		// LD (0xFF00+C),A
		gb.memory.set(0xFF00+uint16(gb.cpu.get_register("C")), gb.cpu.get_register("A"))
		return 2
	case 0x3A:
		// LD A,(HL-)
		// Load A with the value at memory address HL, then decrement HL
		currentHL := gb.cpu.get_register16("HL")
		gb.cpu.set_register("A", gb.memory.get(currentHL))
		gb.cpu.set_register16("HL", currentHL-1)
		return 2
	case 0x32:
		// LD (HL-),A
		// Set memory address HL to the value in A, then decrement HL
		currentHL := gb.cpu.get_register16("HL")
		gb.memory.set(currentHL, gb.cpu.get_register("A"))
		gb.cpu.set_register16("HL", currentHL-1)
		return 2
	case 0x2A:
		// LD A,(HL+)
		// Load A with the value at memory address HL, then increment HL
		currentHL := gb.cpu.get_register16("HL")
		gb.cpu.set_register("A", gb.memory.get(currentHL))
		gb.cpu.set_register16("HL", currentHL+1)
		return 2
	case 0x22:
		// LD (HL+),A
		// Set memory address HL to the value in A, then increment HL
		currentHL := gb.cpu.get_register16("HL")
		gb.memory.set(currentHL, gb.cpu.get_register("A"))
		gb.cpu.set_register16("HL", currentHL+1)
		return 2
	case 0xE0:
		// LD (0xFF00+n),A
		gb.memory.set(0xFF00+uint16(gb.popPC()), gb.cpu.get_register("A"))
		return 3
	case 0xF0:
		// LD A,(0xFF00+n)
		gb.cpu.set_register("A", gb.memory.get(0xFF00+uint16(gb.popPC())))
		return 3
	//////////////// 16-bit loads ////////////////
	case 0x01:
		// LD BC,nn
		gb.cpu.set_register16("BC", gb.popPC16())
		return 3
	case 0x11:
		// LD DE,nn
		gb.cpu.set_register16("DE", gb.popPC16())
		return 3
	case 0x21:
		// LD HL,nn
		gb.cpu.set_register16("HL", gb.popPC16())
		return 3
	case 0x31:
		// LD SP,nn
		gb.cpu.set_register16("SP", gb.popPC16())
		return 3
	case 0xF9:
		// LD SP,HL
		gb.cpu.set_register16("SP", gb.cpu.get_register16("HL"))
		return 2
	case 0xF8:
		// LD HL,SP+n
		// n is a signed 8-bit immediate value
		var pc int32 = int32(int8(gb.popPC()))
		var sp int32 = int32(gb.cpu.get_register16("SP"))
		gb.cpu.set_register16("HL", uint16(sp+pc))
		// Did we overflow the lower nybble?
		gb.cpu.set_flag(FlagH, (sp&0x000F)+(pc&0x000F) > 0x000F)
		// Did we overflow the upper nybble?
		gb.cpu.set_flag(FlagC, (sp&0x00FF)+(pc&0x00FF) > 0x00FF)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		return 3
	case 0x08:
		// LD (nn),SP
		sp := gb.cpu.get_register16("SP")
		lsb := uint8(sp & 0xFF)
		msb := uint8(sp >> 8)
		adr := gb.popPC16()
		gb.memory.set(adr, lsb)
		gb.memory.set(adr+1, msb)
		return 5
	/////////////// Push ////////////////////
	case 0xF5:
		// PUSH AF
		gb.pushToStack(gb.cpu.get_register("A"), gb.cpu.get_register("F"))
		return 4
	case 0xC5:
		// PUSH BC
		gb.pushToStack(gb.cpu.get_register("B"), gb.cpu.get_register("C"))
		return 4
	case 0xD5:
		// PUSH DE
		gb.pushToStack(gb.cpu.get_register("D"), gb.cpu.get_register("E"))
		return 4
	case 0xE5:
		// PUSH HL
		gb.pushToStack(gb.cpu.get_register("H"), gb.cpu.get_register("L"))
		return 4
	/////////////// POP ////////////////////
	case 0xF1:
		// POP AF
		gb.cpu.set_register16("AF", gb.popFromStack())
		return 3
	case 0xC1:
		// POP BC
		gb.cpu.set_register16("BC", gb.popFromStack())
		return 3
	case 0xD1:
		// POP DE
		gb.cpu.set_register16("DE", gb.popFromStack())
		return 3
	case 0xE1:
		// POP HL
		gb.cpu.set_register16("HL", gb.popFromStack())
		return 3
	/////////////// 8-bit Arithmetic ////////////////////
	case 0x87:
		// ADD A,A
		gb.cpu.addToRegisterA(gb.cpu.get_register("A"), false)
		return 1
	case 0x80:
		// ADD A,B
		gb.cpu.addToRegisterA(gb.cpu.get_register("B"), false)
		return 1
	case 0x81:
		// ADD A,C
		gb.cpu.addToRegisterA(gb.cpu.get_register("C"), false)
		return 1
	case 0x82:
		// ADD A,D
		gb.cpu.addToRegisterA(gb.cpu.get_register("D"), false)
		return 1
	case 0x83:
		// ADD A,E
		gb.cpu.addToRegisterA(gb.cpu.get_register("E"), false)
		return 1
	case 0x84:
		// ADD A,H
		gb.cpu.addToRegisterA(gb.cpu.get_register("H"), false)
		return 1
	case 0x85:
		// ADD A,L
		gb.cpu.addToRegisterA(gb.cpu.get_register("L"), false)
		return 1
	case 0x86:
		// ADD A,(HL)
		gb.cpu.addToRegisterA(gb.memory.get(gb.cpu.get_register16("HL")), false)
		return 2
	case 0xC6:
		// ADD A,n
		gb.cpu.addToRegisterA(gb.popPC(), false)
		return 2
	case 0x8F:
		// ADC A,A
		gb.cpu.addToRegisterA(gb.cpu.get_register("A"), true)
		return 1
	case 0x88:
		// ADC A,B
		gb.cpu.addToRegisterA(gb.cpu.get_register("B"), true)
		return 1
	case 0x89:
		// ADC A,C
		gb.cpu.addToRegisterA(gb.cpu.get_register("C"), true)
		return 1
	case 0x8A:
		// ADC A,D
		gb.cpu.addToRegisterA(gb.cpu.get_register("D"), true)
		return 1
	case 0x8B:
		// ADC A,E
		gb.cpu.addToRegisterA(gb.cpu.get_register("E"), true)
		return 1
	case 0x8C:
		// ADC A,H
		gb.cpu.addToRegisterA(gb.cpu.get_register("H"), true)
		return 1
	case 0x8D:
		// ADC A,L
		gb.cpu.addToRegisterA(gb.cpu.get_register("L"), true)
		return 1
	case 0x8E:
		// ADC A,(HL)
		gb.cpu.addToRegisterA(gb.memory.get(gb.cpu.get_register16("HL")), true)
		return 2
	case 0xCE:
		// ADC A,n
		gb.cpu.addToRegisterA(gb.popPC(), true)
		return 2
	case 0x97:
		// SUB A
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("A"), false)
		return 1
	case 0x90:
		// SUB B
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("B"), false)
		return 1
	case 0x91:
		// SUB C
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("C"), false)
		return 1
	case 0x92:
		// SUB D
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("D"), false)
		return 1
	case 0x93:
		// SUB E
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("E"), false)
		return 1
	case 0x94:
		// SUB H
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("H"), false)
		return 1
	case 0x95:
		// SUB L
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("L"), false)
		return 1
	case 0x96:
		// SUB (HL)
		gb.cpu.subtractFromRegisterA(gb.memory.get(gb.cpu.get_register16("HL")), false)
		return 2
	case 0xD6:
		// SUB n
		gb.cpu.subtractFromRegisterA(gb.popPC(), false)
		return 2
	case 0x9F:
		// SBC A,A
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("A"), true)
		return 1
	case 0x98:
		// SBC A,B
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("B"), true)
		return 1
	case 0x99:
		// SBC A,C
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("C"), true)
		return 1
	case 0x9A:
		// SBC A,D
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("D"), true)
		return 1
	case 0x9B:
		// SBC A,E
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("E"), true)
		return 1
	case 0x9C:
		// SBC A,H
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("H"), true)
		return 1
	case 0x9D:
		// SBC A,L
		gb.cpu.subtractFromRegisterA(gb.cpu.get_register("L"), true)
		return 1
	case 0x9E:
		// SBC A,(HL)
		gb.cpu.subtractFromRegisterA(gb.memory.get(gb.cpu.get_register16("HL")), true)
		return 2
	case 0xDE:
		// SBC A,n
		gb.cpu.subtractFromRegisterA(gb.popPC(), true)
		return 2
	case 0xA7:
		// AND A
		gb.cpu.andA(gb.cpu.get_register("A"))
		return 1
	case 0xA0:
		// AND B
		gb.cpu.andA(gb.cpu.get_register("B"))
		return 1
	case 0xA1:
		// AND C
		gb.cpu.andA(gb.cpu.get_register("C"))
		return 1
	case 0xA2:
		// AND D
		gb.cpu.andA(gb.cpu.get_register("D"))
		return 1
	case 0xA3:
		// AND E
		gb.cpu.andA(gb.cpu.get_register("E"))
		return 1
	case 0xA4:
		// AND H
		gb.cpu.andA(gb.cpu.get_register("H"))
		return 1
	case 0xA5:
		// AND L
		gb.cpu.andA(gb.cpu.get_register("L"))
		return 1
	case 0xA6:
		// AND (HL)
		gb.cpu.andA(gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0xE6:
		// AND n
		gb.cpu.andA(gb.popPC())
		return 2
	case 0xB7:
		// OR A
		gb.cpu.orA(gb.cpu.get_register("A"))
		return 1
	case 0xB0:
		// OR B
		gb.cpu.orA(gb.cpu.get_register("B"))
		return 1
	case 0xB1:
		// OR C
		gb.cpu.orA(gb.cpu.get_register("C"))
		return 1
	case 0xB2:
		// OR D
		gb.cpu.orA(gb.cpu.get_register("D"))
		return 1
	case 0xB3:
		// OR E
		gb.cpu.orA(gb.cpu.get_register("E"))
		return 1
	case 0xB4:
		// OR H
		gb.cpu.orA(gb.cpu.get_register("H"))
		return 1
	case 0xB5:
		// OR L
		gb.cpu.orA(gb.cpu.get_register("L"))
		return 1
	case 0xB6:
		// OR (HL)
		gb.cpu.orA(gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0xF6:
		// OR n
		gb.cpu.orA(gb.popPC())
		return 2
	case 0xAF:
		// XOR A
		gb.cpu.xorA(gb.cpu.get_register("A"))
		return 1
	case 0xA8:
		// XOR B
		gb.cpu.xorA(gb.cpu.get_register("B"))
		return 1
	case 0xA9:
		// XOR C
		gb.cpu.xorA(gb.cpu.get_register("C"))
		return 1
	case 0xAA:
		// XOR D
		gb.cpu.xorA(gb.cpu.get_register("D"))
		return 1
	case 0xAB:
		// XOR E
		gb.cpu.xorA(gb.cpu.get_register("E"))
		return 1
	case 0xAC:
		// XOR H
		gb.cpu.xorA(gb.cpu.get_register("H"))
		return 1
	case 0xAD:
		// XOR L
		gb.cpu.xorA(gb.cpu.get_register("L"))
		return 1
	case 0xAE:
		// XOR (HL)
		gb.cpu.xorA(gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0xEE:
		// XOR n
		gb.cpu.xorA(gb.popPC())
		return 2
	case 0xBF:
		// CP A
		gb.cpu.compareA(gb.cpu.get_register("A"))
		return 1
	case 0xB8:
		// CP B
		gb.cpu.compareA(gb.cpu.get_register("B"))
		return 1
	case 0xB9:
		// CP C
		gb.cpu.compareA(gb.cpu.get_register("C"))
		return 1
	case 0xBA:
		// CP D
		gb.cpu.compareA(gb.cpu.get_register("D"))
		return 1
	case 0xBB:
		// CP E
		gb.cpu.compareA(gb.cpu.get_register("E"))
		return 1
	case 0xBC:
		// CP H
		gb.cpu.compareA(gb.cpu.get_register("H"))
		return 1
	case 0xBD:
		// CP L
		gb.cpu.compareA(gb.cpu.get_register("L"))
		return 1
	case 0xBE:
		// CP (HL)
		gb.cpu.compareA(gb.memory.get(gb.cpu.get_register16("HL")))
		return 2
	case 0xFE:
		// CP n
		gb.cpu.compareA(gb.popPC())
		return 2
	case 0x3C:
		// INC A
		gb.cpu.incrementRegister("A")
		return 1
	case 0x04:
		// INC B
		gb.cpu.incrementRegister("B")
		return 1
	case 0x0C:
		// INC C
		gb.cpu.incrementRegister("C")
		return 1
	case 0x14:
		// INC D
		gb.cpu.incrementRegister("D")
		return 1
	case 0x1C:
		// INC E
		gb.cpu.incrementRegister("E")
		return 1
	case 0x24:
		// INC H
		gb.cpu.incrementRegister("H")
		return 1
	case 0x2C:
		// INC L
		gb.cpu.incrementRegister("L")
		return 1
	case 0x34:
		// INC (HL)
		adr := gb.cpu.get_register16("HL")
		value := gb.memory.get(adr)
		gb.cpu.set_flag(FlagZ, value == 0xFF)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, (value&0xF) == 0xF)
		gb.memory.set(adr, value+1)
		return 3
	case 0x3D:
		// DEC A
		gb.cpu.decrementRegister("A")
		return 1
	case 0x05:
		// DEC B
		gb.cpu.decrementRegister("B")
		return 1
	case 0x0D:
		// DEC C
		gb.cpu.decrementRegister("C")
		return 1
	case 0x15:
		// DEC D
		gb.cpu.decrementRegister("D")
		return 1
	case 0x1D:
		// DEC E
		gb.cpu.decrementRegister("E")
		return 1
	case 0x25:
		// DEC H
		gb.cpu.decrementRegister("H")
		return 1
	case 0x2D:
		// DEC L
		gb.cpu.decrementRegister("L")
		return 1
	case 0x35:
		// DEC (HL)
		adr := gb.cpu.get_register16("HL")
		value := gb.memory.get(adr)
		gb.cpu.set_flag(FlagZ, value == 0x01)
		gb.cpu.set_flag(FlagN, true)
		gb.cpu.set_flag(FlagH, (value&0xF) == 0x0)
		gb.memory.set(adr, value-1)
		return 3
	/////////////// 16-bit Arithmetic ////////////////////
	case 0x09:
		// ADD HL,BC
		gb.cpu.addToRegisterHL(gb.cpu.get_register16("BC"))
		return 2
	case 0x19:
		// ADD HL,DE
		gb.cpu.addToRegisterHL(gb.cpu.get_register16("DE"))
		return 2
	case 0x29:
		// ADD HL,HL
		gb.cpu.addToRegisterHL(gb.cpu.get_register16("HL"))
		return 2
	case 0x39:
		// ADD HL,SP
		gb.cpu.addToRegisterHL(gb.cpu.get_register16("SP"))
		return 2
	case 0xE8:
		// ADD SP,n
		// n is a signed 8-bit immediate value
		// TODO: break this off into a function since it is shared with opcode 0xF8?
		var pc int32 = int32(int8(gb.popPC()))
		var sp int32 = int32(gb.cpu.get_register16("SP"))
		gb.cpu.set_register16("SP", uint16(sp+pc))
		// Did we overflow the lower nybble?
		gb.cpu.set_flag(FlagH, (sp&0x000F)+(pc&0x000F) > 0x000F)
		// Did we overflow the upper nybble?
		gb.cpu.set_flag(FlagC, (sp&0x00FF)+(pc&0x00FF) > 0x00FF)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		return 4
	case 0x03:
		// INC BC
		gb.cpu.set_register16("BC", gb.cpu.get_register16("BC")+1)
		return 2
	case 0x13:
		// INC DE
		gb.cpu.set_register16("DE", gb.cpu.get_register16("DE")+1)
		return 2
	case 0x23:
		// INC HL
		gb.cpu.set_register16("HL", gb.cpu.get_register16("HL")+1)
		return 2
	case 0x33:
		// INC SP
		gb.cpu.set_register16("SP", gb.cpu.get_register16("SP")+1)
		return 2
	case 0x0B:
		// DEC BC
		gb.cpu.set_register16("BC", gb.cpu.get_register16("BC")-1)
		return 2
	case 0x1B:
		// DEC DE
		gb.cpu.set_register16("DE", gb.cpu.get_register16("DE")-1)
		return 2
	case 0x2B:
		// DEC HL
		gb.cpu.set_register16("HL", gb.cpu.get_register16("HL")-1)
		return 2
	case 0x3B:
		// DEC SP
		gb.cpu.set_register16("SP", gb.cpu.get_register16("SP")-1)
		return 2
	/////////////// Rotates ////////////////////
	case 0x07:
		// RLCA
		// Rotate Left and set Carry bit, register A
		value := gb.cpu.get_register("A")

		carrybit := value >> 7
		result := (value << 1) | carrybit

		gb.cpu.set_register("A", result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, carrybit == 1)
		return 1
	case 0x17:
		// RLA
		// Rotate Left through carry flag, register A
		value := gb.cpu.get_register("A")

		var oldcarry uint8 = 0
		if gb.cpu.get_flag(FlagC) {
			oldcarry = 1
		}
		newcarry := value >> 7
		result := (value << 1) | oldcarry

		gb.cpu.set_register("A", result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, newcarry == 1)
		return 1
	case 0x0F:
		// RRCA
		// Rotate Right and set Carry bit, register A
		value := gb.cpu.get_register("A")

		carrybit := value & 1
		result := (value >> 1) | (carrybit << 7)

		gb.cpu.set_register("A", result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, carrybit == 1)
		return 1
	case 0x1F:
		// RRA
		// Rotate Right through carry bit, register A
		value := gb.cpu.get_register("A")

		var oldcarry uint8 = 0
		if gb.cpu.get_flag(FlagC) {
			oldcarry = 1
		}
		newcarry := value & 1
		result := (value >> 1) | (oldcarry << 7)

		gb.cpu.set_register("A", result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, newcarry == 1)
		return 1
	/////////////// Misc ////////////////////
	case 0x27:
		// DAA
		// Decimal adjust register A for binary coded decimal after an add or subtract
		a := uint16(gb.cpu.get_register("A"))

		if gb.cpu.get_flag(FlagN) {
			// Previous operation was subtraction
			if gb.cpu.get_flag(FlagH) {
				// Underflow lower nybble
				a -= 0x6
			}
			if gb.cpu.get_flag(FlagC) {
				// Underflow upper nybble
				a -= 0x60
			}
		} else {
			// Previous operation was addition
			if (a&0xF) > 0x9 || gb.cpu.get_flag(FlagH) {
				// Overflow lower nybble
				a += 0x6
			}
			if a > 0x9F || gb.cpu.get_flag(FlagC) {
				// Overflow upper nybble
				a += 0x60
				gb.cpu.set_flag(FlagC, true)
			}
		}
		gb.cpu.set_flag(FlagZ, uint8(a) == 0)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_register("A", uint8(a))
		return 1
	case 0x2F:
		// CPL
		gb.cpu.set_register("A", ^gb.cpu.get_register("A"))
		gb.cpu.set_flag(FlagN, true)
		gb.cpu.set_flag(FlagH, true)
		return 1
	case 0x3F:
		// CCF
		// Complement Carry Flag
		if gb.cpu.get_flag(FlagC) {
			gb.cpu.set_flag(FlagC, false)
		} else {
			gb.cpu.set_flag(FlagC, true)
		}
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		return 1
	case 0x37:
		// SCF
		// Set Carry Flag
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, true)
		return 1
	case 0x00:
		// NOP
		return 1
	case 0xF3:
		// DI
		// Disable Interrupts
		// based on available documentation it seems that Game Boy Color and later models
		// introduce a delay here, but DMG disables interrupts immediately
		gb.interruptMasterEnable = false
		return 1
	case 0xFB:
		// EI
		// Enable Interrupts
		gb.pendingInterruptEnable = true
		return 1
	case 0x76:
		// HALT
		gb.halted = true
		return 1
	// TODO: implement stop, others
	/////////////// Jumps ////////////////////
	case 0xC3:
		// JP nn
		gb.cpu.set_register16("PC", gb.popPC16())
		return 4
	case 0xC2:
		// JP NZ,nn
		adr := gb.popPC16()
		if !gb.cpu.get_flag(FlagZ) {
			gb.cpu.set_register16("PC", adr)
			return 4
		}
		return 3
	case 0xCA:
		// JP Z,nn
		adr := gb.popPC16()
		if gb.cpu.get_flag(FlagZ) {
			gb.cpu.set_register16("PC", adr)
			return 4
		}
		return 3
	case 0xD2:
		// JP NC,nn
		adr := gb.popPC16()
		if !gb.cpu.get_flag(FlagC) {
			gb.cpu.set_register16("PC", adr)
			return 4
		}
		return 3
	case 0xDA:
		// JP C,nn
		adr := gb.popPC16()
		if gb.cpu.get_flag(FlagC) {
			gb.cpu.set_register16("PC", adr)
			return 4
		}
		return 3
	case 0xE9:
		// JP HL
		gb.cpu.set_register16("PC", gb.cpu.get_register16("HL"))
		return 1
	case 0x18:
		// JR n
		// Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		// Using value of PC after incrementing
		var pc int32 = int32(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", uint16(pc+n))
		return 3
	case 0x20:
		// JR NZ,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if !gb.cpu.get_flag(FlagZ) {
			var pc int32 = int32(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", uint16(pc+n))
			return 3
		}
		return 2
	case 0x28:
		// JR Z,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if gb.cpu.get_flag(FlagZ) {
			var pc int32 = int32(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", uint16(pc+n))
			return 3
		}
		return 2
	case 0x30:
		// JR NC,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if !gb.cpu.get_flag(FlagC) {
			var pc int32 = int32(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", uint16(pc+n))
			return 3
		}
		return 2
	case 0x38:
		// JR C,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if gb.cpu.get_flag(FlagC) {
			var pc int32 = int32(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", uint16(pc+n))
			return 3
		}
		return 2
	case 0xCD:
		// CALL nn
		nn := gb.popPC16()
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", nn)
		return 6
	case 0xC4:
		// CALL NZ,nn
		nn := gb.popPC16()
		if !gb.cpu.get_flag(FlagZ) {
			gb.pushToStack16(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", nn)
			return 6
		}
		return 3
	case 0xCC:
		// CALL Z,nn
		nn := gb.popPC16()
		if gb.cpu.get_flag(FlagZ) {
			gb.pushToStack16(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", nn)
			return 6
		}
		return 3
	case 0xD4:
		// CALL NC,nn
		nn := gb.popPC16()
		if !gb.cpu.get_flag(FlagC) {
			gb.pushToStack16(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", nn)
			return 6
		}
		return 3
	case 0xDC:
		// CALL C,nn
		nn := gb.popPC16()
		if gb.cpu.get_flag(FlagC) {
			gb.pushToStack16(gb.cpu.get_register16("PC"))
			gb.cpu.set_register16("PC", nn)
			return 6
		}
		return 3
	case 0xC7:
		// RST 00H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x00)
		return 4
	case 0xCF:
		// RST 08H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x08)
		return 4
	case 0xD7:
		// RST 10H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x10)
		return 4
	case 0xDF:
		// RST 18H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x18)
		return 4
	case 0xE7:
		// RST 20H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x20)
		return 4
	case 0xEF:
		// RST 28H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x28)
		return 4
	case 0xF7:
		// RST 30H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x30)
		return 4
	case 0xFF:
		// RST 38H
		gb.pushToStack16(gb.cpu.get_register16("PC"))
		gb.cpu.set_register16("PC", 0x38)
		return 4
	case 0xC9:
		// RET
		gb.cpu.set_register16("PC", gb.popFromStack())
		return 2
	case 0xC0:
		// RET NZ
		if !gb.cpu.get_flag(FlagZ) {
			gb.cpu.set_register16("PC", gb.popFromStack())
			return 5
		}
		return 2
	case 0xC8:
		// RET Z
		if gb.cpu.get_flag(FlagZ) {
			gb.cpu.set_register16("PC", gb.popFromStack())
			return 5
		}
		return 2
	case 0xD0:
		// RET NC
		if !gb.cpu.get_flag(FlagC) {
			gb.cpu.set_register16("PC", gb.popFromStack())
			return 5
		}
		return 2
	case 0xD8:
		// RET C
		if !gb.cpu.get_flag(FlagC) {
			gb.cpu.set_register16("PC", gb.popFromStack())
			return 5
		}
		return 2
	case 0xD9:
		// RETI
		// Return and enable interrupts
		gb.cpu.set_register16("PC", gb.popFromStack())
		gb.interruptMasterEnable = true
		return 2
	////////////// CB - Extended Instructions /////////////
	case 0xCB:
		// CB
		// TODO: We could also handle this by just setting a flag to indicate that the
		// last opcode was CB and then dispatch that when this function is called again
		// to allow other processes to occur between them
		next_opcode := gb.popPC()
		return gb.CBOpcode(next_opcode) + 1
	// TODO: Handle intentionally unimplemented opcodes (should crash Gameboy, but maybe handle gracefully)
	default:
		panic(fmt.Sprintf("opcode 0x%X not implemented", opcode))
	}
}
