package main

import "fmt"

// Opcodes data from https://pastraiser.com/cpu/gameboy/gameboy_opcodes.html
// and http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf with a few error corrections

// Set a register value equal to current value of PC plus an 8 bit signed immediate value
func (gb *Gameboy) setSPplusN(reg register16) {
	var pc int32 = int32(int8(gb.popPC()))
	var sp int32 = int32(gb.cpu.getRegister16(regSP))
	gb.cpu.setRegister16(reg, uint16(sp+pc))
	// Did we overflow the lower nybble?
	gb.cpu.set_flag(FlagH, (sp&0x000F)+(pc&0x000F) > 0x000F)
	// Did we overflow the upper nybble?
	gb.cpu.set_flag(FlagC, (sp&0x00FF)+(pc&0x00FF) > 0x00FF)
	gb.cpu.set_flag(FlagZ, false)
	gb.cpu.set_flag(FlagN, false)
}

// Execute a single opcode and return the number of CPU cycles it took (1MHz CPU cycles)
func (gb *Gameboy) Opcode(opcode uint8) int {
	switch opcode {
	//////////////// 8-bit loads ////////////////
	case 0x3E:
		// LD A,n
		gb.cpu.setRegister(regA, gb.popPC())
		return 2
	case 0x06:
		// LD B,n
		gb.cpu.setRegister(regB, gb.popPC())
		return 2
	case 0x0E:
		// LD C,n
		gb.cpu.setRegister(regC, gb.popPC())
		return 2
	case 0x16:
		// LD D,n
		gb.cpu.setRegister(regD, gb.popPC())
		return 2
	case 0x1E:
		// LD E,n
		gb.cpu.setRegister(regE, gb.popPC())
		return 2
	case 0x26:
		// LD H,n
		gb.cpu.setRegister(regH, gb.popPC())
		return 2
	case 0x2E:
		// LD L,n
		gb.cpu.setRegister(regL, gb.popPC())
		return 2
	case 0x7F, 0x40, 0x49, 0x52, 0x5B, 0x64, 0x6D:
		// LD X,X (For registers A, B, C, D, E, H, L)
		// Equivalent to NOP
		return 1
	case 0x78:
		// LD A,B
		gb.cpu.setRegister(regA, gb.cpu.getRegister(regB))
		return 1
	case 0x79:
		// LD A,C
		gb.cpu.setRegister(regA, gb.cpu.getRegister(regC))
		return 1
	case 0x7A:
		// LD A,D
		gb.cpu.setRegister(regA, gb.cpu.getRegister(regD))
		return 1
	case 0x7B:
		// LD A,E
		gb.cpu.setRegister(regA, gb.cpu.getRegister(regE))
		return 1
	case 0x7C:
		// LD A,H
		gb.cpu.setRegister(regA, gb.cpu.getRegister(regH))
		return 1
	case 0x7D:
		// LD A,L
		gb.cpu.setRegister(regA, gb.cpu.getRegister(regL))
		return 1
	case 0x7E:
		// LD A,(HL)
		gb.cpu.setRegister(regA, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x41:
		// LD B,C
		gb.cpu.setRegister(regB, gb.cpu.getRegister(regC))
		return 1
	case 0x42:
		// LD B,D
		gb.cpu.setRegister(regB, gb.cpu.getRegister(regD))
		return 1
	case 0x43:
		// LD B,E
		gb.cpu.setRegister(regB, gb.cpu.getRegister(regE))
		return 1
	case 0x44:
		// LD B,H
		gb.cpu.setRegister(regB, gb.cpu.getRegister(regH))
		return 1
	case 0x45:
		// LD B,L
		gb.cpu.setRegister(regB, gb.cpu.getRegister(regL))
		return 1
	case 0x46:
		// LD B,(HL)
		gb.cpu.setRegister(regB, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x47:
		// LD B,A
		gb.cpu.setRegister(regB, gb.cpu.getRegister(regA))
		return 1
	case 0x48:
		// LD C,B
		gb.cpu.setRegister(regC, gb.cpu.getRegister(regB))
		return 1
	case 0x4A:
		// LD C,D
		gb.cpu.setRegister(regC, gb.cpu.getRegister(regD))
		return 1
	case 0x4B:
		// LD C,E
		gb.cpu.setRegister(regC, gb.cpu.getRegister(regE))
		return 1
	case 0x4C:
		// LD C,H
		gb.cpu.setRegister(regC, gb.cpu.getRegister(regH))
		return 1
	case 0x4D:
		// LD C,L
		gb.cpu.setRegister(regC, gb.cpu.getRegister(regL))
		return 1
	case 0x4E:
		// LD C,(HL)
		gb.cpu.setRegister(regC, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x4F:
		// LD C,A
		gb.cpu.setRegister(regC, gb.cpu.getRegister(regA))
		return 1
	case 0x50:
		// LD D,B
		gb.cpu.setRegister(regD, gb.cpu.getRegister(regB))
		return 1
	case 0x51:
		// LD D,C
		gb.cpu.setRegister(regD, gb.cpu.getRegister(regC))
		return 1
	case 0x53:
		// LD D,E
		gb.cpu.setRegister(regD, gb.cpu.getRegister(regE))
		return 1
	case 0x54:
		// LD D,H
		gb.cpu.setRegister(regD, gb.cpu.getRegister(regH))
		return 1
	case 0x55:
		// LD D,L
		gb.cpu.setRegister(regD, gb.cpu.getRegister(regL))
		return 1
	case 0x56:
		// LD D,(HL)
		gb.cpu.setRegister(regD, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x57:
		// LD D,A
		gb.cpu.setRegister(regD, gb.cpu.getRegister(regA))
		return 1
	case 0x58:
		// LD E,B
		gb.cpu.setRegister(regE, gb.cpu.getRegister(regB))
		return 1
	case 0x59:
		// LD E,C
		gb.cpu.setRegister(regE, gb.cpu.getRegister(regC))
		return 1
	case 0x5A:
		// LD E,D
		gb.cpu.setRegister(regE, gb.cpu.getRegister(regD))
		return 1
	case 0x5C:
		// LD E,H
		gb.cpu.setRegister(regE, gb.cpu.getRegister(regH))
		return 1
	case 0x5D:
		// LD E,L
		gb.cpu.setRegister(regE, gb.cpu.getRegister(regL))
		return 1
	case 0x5E:
		// LD E,(HL)
		gb.cpu.setRegister(regE, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x5F:
		// LD E,A
		gb.cpu.setRegister(regE, gb.cpu.getRegister(regA))
		return 1
	case 0x60:
		// LD H,B
		gb.cpu.setRegister(regH, gb.cpu.getRegister(regB))
		return 1
	case 0x61:
		// LD H,C
		gb.cpu.setRegister(regH, gb.cpu.getRegister(regC))
		return 1
	case 0x62:
		// LD H,D
		gb.cpu.setRegister(regH, gb.cpu.getRegister(regD))
		return 1
	case 0x63:
		// LD H,E
		gb.cpu.setRegister(regH, gb.cpu.getRegister(regE))
		return 1
	case 0x65:
		// LD H,L
		gb.cpu.setRegister(regH, gb.cpu.getRegister(regL))
		return 1
	case 0x66:
		// LD H,(HL)
		gb.cpu.setRegister(regH, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x67:
		// LD H,A
		gb.cpu.setRegister(regH, gb.cpu.getRegister(regA))
		return 1
	case 0x68:
		// LD L,B
		gb.cpu.setRegister(regL, gb.cpu.getRegister(regB))
		return 1
	case 0x69:
		// LD L,C
		gb.cpu.setRegister(regL, gb.cpu.getRegister(regC))
		return 1
	case 0x6A:
		// LD L,D
		gb.cpu.setRegister(regL, gb.cpu.getRegister(regD))
		return 1
	case 0x6B:
		// LD L,E
		gb.cpu.setRegister(regL, gb.cpu.getRegister(regE))
		return 1
	case 0x6C:
		// LD L,H
		gb.cpu.setRegister(regL, gb.cpu.getRegister(regH))
		return 1
	case 0x6E:
		// LD L,(HL)
		gb.cpu.setRegister(regL, gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0x6F:
		// LD L,A
		gb.cpu.setRegister(regL, gb.cpu.getRegister(regA))
		return 1
	case 0x70:
		// LD (HL),B
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regB))
		return 2
	case 0x71:
		// LD (HL),C
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regC))
		return 2
	case 0x72:
		// LD (HL),D
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regD))
		return 2
	case 0x73:
		// LD (HL),E
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regE))
		return 2
	case 0x74:
		// LD (HL),H
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regH))
		return 2
	case 0x75:
		// LD (HL),L
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regL))
		return 2
	case 0x36:
		// LD (HL),n
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.popPC())
		return 3
	case 0x0A:
		// LD A,(BC)
		gb.cpu.setRegister(regA, gb.memory.get(gb.cpu.getRegister16(regBC)))
		return 2
	case 0x1A:
		// LD A,(DE)
		gb.cpu.setRegister(regA, gb.memory.get(gb.cpu.getRegister16(regDE)))
		return 2
	case 0xFA:
		// LD A,(nn)
		gb.cpu.setRegister(regA, gb.memory.get(gb.popPC16()))
		return 4
	case 0x02:
		// LD (BC),A
		gb.memory.set(gb.cpu.getRegister16(regBC), gb.cpu.getRegister(regA))
		return 2
	case 0x12:
		// LD (DE),A
		gb.memory.set(gb.cpu.getRegister16(regDE), gb.cpu.getRegister(regA))
		return 2
	case 0x77:
		// LD (HL),A
		gb.memory.set(gb.cpu.getRegister16(regHL), gb.cpu.getRegister(regA))
		return 2
	case 0xEA:
		// LD (nn),A
		gb.memory.set(gb.popPC16(), gb.cpu.getRegister(regA))
		return 4
	case 0xF2:
		// LD A,(0xFF00+C)
		gb.cpu.setRegister(regA, gb.memory.get(0xFF00+uint16(gb.cpu.getRegister(regC))))
		return 2
	case 0xE2:
		// LD (0xFF00+C),A
		gb.memory.set(0xFF00+uint16(gb.cpu.getRegister(regC)), gb.cpu.getRegister(regA))
		return 2
	case 0x3A:
		// LD A,(HL-)
		// Load A with the value at memory address HL, then decrement HL
		currentHL := gb.cpu.getRegister16(regHL)
		gb.cpu.setRegister(regA, gb.memory.get(currentHL))
		gb.cpu.setRegister16(regHL, currentHL-1)
		return 2
	case 0x32:
		// LD (HL-),A
		// Set memory address HL to the value in A, then decrement HL
		currentHL := gb.cpu.getRegister16(regHL)
		gb.memory.set(currentHL, gb.cpu.getRegister(regA))
		gb.cpu.setRegister16(regHL, currentHL-1)
		return 2
	case 0x2A:
		// LD A,(HL+)
		// Load A with the value at memory address HL, then increment HL
		currentHL := gb.cpu.getRegister16(regHL)
		gb.cpu.setRegister(regA, gb.memory.get(currentHL))
		gb.cpu.setRegister16(regHL, currentHL+1)
		return 2
	case 0x22:
		// LD (HL+),A
		// Set memory address HL to the value in A, then increment HL
		currentHL := gb.cpu.getRegister16(regHL)
		gb.memory.set(currentHL, gb.cpu.getRegister(regA))
		gb.cpu.setRegister16(regHL, currentHL+1)
		return 2
	case 0xE0:
		// LD (0xFF00+n),A
		gb.memory.set(0xFF00+uint16(gb.popPC()), gb.cpu.getRegister(regA))
		return 3
	case 0xF0:
		// LD A,(0xFF00+n)
		gb.cpu.setRegister(regA, gb.memory.get(0xFF00+uint16(gb.popPC())))
		return 3
	//////////////// 16-bit loads ////////////////
	case 0x01:
		// LD BC,nn
		gb.cpu.setRegister16(regBC, gb.popPC16())
		return 3
	case 0x11:
		// LD DE,nn
		gb.cpu.setRegister16(regDE, gb.popPC16())
		return 3
	case 0x21:
		// LD HL,nn
		gb.cpu.setRegister16(regHL, gb.popPC16())
		return 3
	case 0x31:
		// LD SP,nn
		gb.cpu.setRegister16(regSP, gb.popPC16())
		return 3
	case 0xF9:
		// LD SP,HL
		gb.cpu.setRegister16(regSP, gb.cpu.getRegister16(regHL))
		return 2
	case 0xF8:
		// LD HL,SP+n
		// n is a signed 8-bit immediate value
		gb.setSPplusN(regHL)
		return 3
	case 0x08:
		// LD (nn),SP
		sp := gb.cpu.getRegister16(regSP)
		lsb := uint8(sp & 0xFF)
		msb := uint8(sp >> 8)
		adr := gb.popPC16()
		gb.memory.set(adr, lsb)
		gb.memory.set(adr+1, msb)
		return 5
	/////////////// Push ////////////////////
	case 0xF5:
		// PUSH AF
		gb.pushToStack(gb.cpu.getRegister(regA), gb.cpu.getRegister(regF))
		return 4
	case 0xC5:
		// PUSH BC
		gb.pushToStack(gb.cpu.getRegister(regB), gb.cpu.getRegister(regC))
		return 4
	case 0xD5:
		// PUSH DE
		gb.pushToStack(gb.cpu.getRegister(regD), gb.cpu.getRegister(regE))
		return 4
	case 0xE5:
		// PUSH HL
		gb.pushToStack(gb.cpu.getRegister(regH), gb.cpu.getRegister(regL))
		return 4
	/////////////// POP ////////////////////
	case 0xF1:
		// POP AF
		gb.cpu.setRegister16(regAF, gb.popFromStack())
		return 3
	case 0xC1:
		// POP BC
		gb.cpu.setRegister16(regBC, gb.popFromStack())
		return 3
	case 0xD1:
		// POP DE
		gb.cpu.setRegister16(regDE, gb.popFromStack())
		return 3
	case 0xE1:
		// POP HL
		gb.cpu.setRegister16(regHL, gb.popFromStack())
		return 3
	/////////////// 8-bit Arithmetic ////////////////////
	case 0x87:
		// ADD A,A
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regA), false)
		return 1
	case 0x80:
		// ADD A,B
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regB), false)
		return 1
	case 0x81:
		// ADD A,C
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regC), false)
		return 1
	case 0x82:
		// ADD A,D
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regD), false)
		return 1
	case 0x83:
		// ADD A,E
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regE), false)
		return 1
	case 0x84:
		// ADD A,H
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regH), false)
		return 1
	case 0x85:
		// ADD A,L
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regL), false)
		return 1
	case 0x86:
		// ADD A,(HL)
		gb.cpu.addToRegisterA(gb.memory.get(gb.cpu.getRegister16(regHL)), false)
		return 2
	case 0xC6:
		// ADD A,n
		gb.cpu.addToRegisterA(gb.popPC(), false)
		return 2
	case 0x8F:
		// ADC A,A
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regA), true)
		return 1
	case 0x88:
		// ADC A,B
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regB), true)
		return 1
	case 0x89:
		// ADC A,C
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regC), true)
		return 1
	case 0x8A:
		// ADC A,D
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regD), true)
		return 1
	case 0x8B:
		// ADC A,E
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regE), true)
		return 1
	case 0x8C:
		// ADC A,H
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regH), true)
		return 1
	case 0x8D:
		// ADC A,L
		gb.cpu.addToRegisterA(gb.cpu.getRegister(regL), true)
		return 1
	case 0x8E:
		// ADC A,(HL)
		gb.cpu.addToRegisterA(gb.memory.get(gb.cpu.getRegister16(regHL)), true)
		return 2
	case 0xCE:
		// ADC A,n
		gb.cpu.addToRegisterA(gb.popPC(), true)
		return 2
	case 0x97:
		// SUB A
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regA), false)
		return 1
	case 0x90:
		// SUB B
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regB), false)
		return 1
	case 0x91:
		// SUB C
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regC), false)
		return 1
	case 0x92:
		// SUB D
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regD), false)
		return 1
	case 0x93:
		// SUB E
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regE), false)
		return 1
	case 0x94:
		// SUB H
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regH), false)
		return 1
	case 0x95:
		// SUB L
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regL), false)
		return 1
	case 0x96:
		// SUB (HL)
		gb.cpu.subtractFromRegisterA(gb.memory.get(gb.cpu.getRegister16(regHL)), false)
		return 2
	case 0xD6:
		// SUB n
		gb.cpu.subtractFromRegisterA(gb.popPC(), false)
		return 2
	case 0x9F:
		// SBC A,A
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regA), true)
		return 1
	case 0x98:
		// SBC A,B
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regB), true)
		return 1
	case 0x99:
		// SBC A,C
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regC), true)
		return 1
	case 0x9A:
		// SBC A,D
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regD), true)
		return 1
	case 0x9B:
		// SBC A,E
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regE), true)
		return 1
	case 0x9C:
		// SBC A,H
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regH), true)
		return 1
	case 0x9D:
		// SBC A,L
		gb.cpu.subtractFromRegisterA(gb.cpu.getRegister(regL), true)
		return 1
	case 0x9E:
		// SBC A,(HL)
		gb.cpu.subtractFromRegisterA(gb.memory.get(gb.cpu.getRegister16(regHL)), true)
		return 2
	case 0xDE:
		// SBC A,n
		gb.cpu.subtractFromRegisterA(gb.popPC(), true)
		return 2
	case 0xA7:
		// AND A
		gb.cpu.andA(gb.cpu.getRegister(regA))
		return 1
	case 0xA0:
		// AND B
		gb.cpu.andA(gb.cpu.getRegister(regB))
		return 1
	case 0xA1:
		// AND C
		gb.cpu.andA(gb.cpu.getRegister(regC))
		return 1
	case 0xA2:
		// AND D
		gb.cpu.andA(gb.cpu.getRegister(regD))
		return 1
	case 0xA3:
		// AND E
		gb.cpu.andA(gb.cpu.getRegister(regE))
		return 1
	case 0xA4:
		// AND H
		gb.cpu.andA(gb.cpu.getRegister(regH))
		return 1
	case 0xA5:
		// AND L
		gb.cpu.andA(gb.cpu.getRegister(regL))
		return 1
	case 0xA6:
		// AND (HL)
		gb.cpu.andA(gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0xE6:
		// AND n
		gb.cpu.andA(gb.popPC())
		return 2
	case 0xB7:
		// OR A
		gb.cpu.orA(gb.cpu.getRegister(regA))
		return 1
	case 0xB0:
		// OR B
		gb.cpu.orA(gb.cpu.getRegister(regB))
		return 1
	case 0xB1:
		// OR C
		gb.cpu.orA(gb.cpu.getRegister(regC))
		return 1
	case 0xB2:
		// OR D
		gb.cpu.orA(gb.cpu.getRegister(regD))
		return 1
	case 0xB3:
		// OR E
		gb.cpu.orA(gb.cpu.getRegister(regE))
		return 1
	case 0xB4:
		// OR H
		gb.cpu.orA(gb.cpu.getRegister(regH))
		return 1
	case 0xB5:
		// OR L
		gb.cpu.orA(gb.cpu.getRegister(regL))
		return 1
	case 0xB6:
		// OR (HL)
		gb.cpu.orA(gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0xF6:
		// OR n
		gb.cpu.orA(gb.popPC())
		return 2
	case 0xAF:
		// XOR A
		gb.cpu.xorA(gb.cpu.getRegister(regA))
		return 1
	case 0xA8:
		// XOR B
		gb.cpu.xorA(gb.cpu.getRegister(regB))
		return 1
	case 0xA9:
		// XOR C
		gb.cpu.xorA(gb.cpu.getRegister(regC))
		return 1
	case 0xAA:
		// XOR D
		gb.cpu.xorA(gb.cpu.getRegister(regD))
		return 1
	case 0xAB:
		// XOR E
		gb.cpu.xorA(gb.cpu.getRegister(regE))
		return 1
	case 0xAC:
		// XOR H
		gb.cpu.xorA(gb.cpu.getRegister(regH))
		return 1
	case 0xAD:
		// XOR L
		gb.cpu.xorA(gb.cpu.getRegister(regL))
		return 1
	case 0xAE:
		// XOR (HL)
		gb.cpu.xorA(gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0xEE:
		// XOR n
		gb.cpu.xorA(gb.popPC())
		return 2
	case 0xBF:
		// CP A
		gb.cpu.compareA(gb.cpu.getRegister(regA))
		return 1
	case 0xB8:
		// CP B
		gb.cpu.compareA(gb.cpu.getRegister(regB))
		return 1
	case 0xB9:
		// CP C
		gb.cpu.compareA(gb.cpu.getRegister(regC))
		return 1
	case 0xBA:
		// CP D
		gb.cpu.compareA(gb.cpu.getRegister(regD))
		return 1
	case 0xBB:
		// CP E
		gb.cpu.compareA(gb.cpu.getRegister(regE))
		return 1
	case 0xBC:
		// CP H
		gb.cpu.compareA(gb.cpu.getRegister(regH))
		return 1
	case 0xBD:
		// CP L
		gb.cpu.compareA(gb.cpu.getRegister(regL))
		return 1
	case 0xBE:
		// CP (HL)
		gb.cpu.compareA(gb.memory.get(gb.cpu.getRegister16(regHL)))
		return 2
	case 0xFE:
		// CP n
		gb.cpu.compareA(gb.popPC())
		return 2
	case 0x3C:
		// INC A
		gb.cpu.incrementRegister(regA)
		return 1
	case 0x04:
		// INC B
		gb.cpu.incrementRegister(regB)
		return 1
	case 0x0C:
		// INC C
		gb.cpu.incrementRegister(regC)
		return 1
	case 0x14:
		// INC D
		gb.cpu.incrementRegister(regD)
		return 1
	case 0x1C:
		// INC E
		gb.cpu.incrementRegister(regE)
		return 1
	case 0x24:
		// INC H
		gb.cpu.incrementRegister(regH)
		return 1
	case 0x2C:
		// INC L
		gb.cpu.incrementRegister(regL)
		return 1
	case 0x34:
		// INC (HL)
		adr := gb.cpu.getRegister16(regHL)
		value := gb.memory.get(adr)
		gb.cpu.set_flag(FlagZ, value == 0xFF)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, (value&0xF) == 0xF)
		gb.memory.set(adr, value+1)
		return 3
	case 0x3D:
		// DEC A
		gb.cpu.decrementRegister(regA)
		return 1
	case 0x05:
		// DEC B
		gb.cpu.decrementRegister(regB)
		return 1
	case 0x0D:
		// DEC C
		gb.cpu.decrementRegister(regC)
		return 1
	case 0x15:
		// DEC D
		gb.cpu.decrementRegister(regD)
		return 1
	case 0x1D:
		// DEC E
		gb.cpu.decrementRegister(regE)
		return 1
	case 0x25:
		// DEC H
		gb.cpu.decrementRegister(regH)
		return 1
	case 0x2D:
		// DEC L
		gb.cpu.decrementRegister(regL)
		return 1
	case 0x35:
		// DEC (HL)
		adr := gb.cpu.getRegister16(regHL)
		value := gb.memory.get(adr)
		gb.cpu.set_flag(FlagZ, value == 0x01)
		gb.cpu.set_flag(FlagN, true)
		gb.cpu.set_flag(FlagH, (value&0xF) == 0x0)
		gb.memory.set(adr, value-1)
		return 3
	/////////////// 16-bit Arithmetic ////////////////////
	case 0x09:
		// ADD HL,BC
		gb.cpu.addToRegisterHL(gb.cpu.getRegister16(regBC))
		return 2
	case 0x19:
		// ADD HL,DE
		gb.cpu.addToRegisterHL(gb.cpu.getRegister16(regDE))
		return 2
	case 0x29:
		// ADD HL,HL
		gb.cpu.addToRegisterHL(gb.cpu.getRegister16(regHL))
		return 2
	case 0x39:
		// ADD HL,SP
		gb.cpu.addToRegisterHL(gb.cpu.getRegister16(regSP))
		return 2
	case 0xE8:
		// ADD SP,n
		// n is a signed 8-bit immediate value
		gb.setSPplusN(regSP)
		return 4
	case 0x03:
		// INC BC
		gb.cpu.setRegister16(regBC, gb.cpu.getRegister16(regBC)+1)
		return 2
	case 0x13:
		// INC DE
		gb.cpu.setRegister16(regDE, gb.cpu.getRegister16(regDE)+1)
		return 2
	case 0x23:
		// INC HL
		gb.cpu.setRegister16(regHL, gb.cpu.getRegister16(regHL)+1)
		return 2
	case 0x33:
		// INC SP
		gb.cpu.setRegister16(regSP, gb.cpu.getRegister16(regSP)+1)
		return 2
	case 0x0B:
		// DEC BC
		gb.cpu.setRegister16(regBC, gb.cpu.getRegister16(regBC)-1)
		return 2
	case 0x1B:
		// DEC DE
		gb.cpu.setRegister16(regDE, gb.cpu.getRegister16(regDE)-1)
		return 2
	case 0x2B:
		// DEC HL
		gb.cpu.setRegister16(regHL, gb.cpu.getRegister16(regHL)-1)
		return 2
	case 0x3B:
		// DEC SP
		gb.cpu.setRegister16(regSP, gb.cpu.getRegister16(regSP)-1)
		return 2
	/////////////// Rotates ////////////////////
	case 0x07:
		// RLCA
		// Rotate Left and set Carry bit, register A
		value := gb.cpu.getRegister(regA)

		carrybit := value >> 7
		result := (value << 1) | carrybit

		gb.cpu.setRegister(regA, result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, carrybit == 1)
		return 1
	case 0x17:
		// RLA
		// Rotate Left through carry flag, register A
		value := gb.cpu.getRegister(regA)

		var oldcarry uint8 = 0
		if gb.cpu.getFlag(FlagC) {
			oldcarry = 1
		}
		newcarry := value >> 7
		result := (value << 1) | oldcarry

		gb.cpu.setRegister(regA, result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, newcarry == 1)
		return 1
	case 0x0F:
		// RRCA
		// Rotate Right and set Carry bit, register A
		value := gb.cpu.getRegister(regA)

		carrybit := value & 1
		result := (value >> 1) | (carrybit << 7)

		gb.cpu.setRegister(regA, result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, carrybit == 1)
		return 1
	case 0x1F:
		// RRA
		// Rotate Right through carry bit, register A
		value := gb.cpu.getRegister(regA)

		var oldcarry uint8 = 0
		if gb.cpu.getFlag(FlagC) {
			oldcarry = 1
		}
		newcarry := value & 1
		result := (value >> 1) | (oldcarry << 7)

		gb.cpu.setRegister(regA, result)
		gb.cpu.set_flag(FlagZ, false)
		gb.cpu.set_flag(FlagN, false)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.set_flag(FlagC, newcarry == 1)
		return 1
	/////////////// Misc ////////////////////
	case 0x27:
		// DAA
		// Decimal adjust register A for binary coded decimal after an add or subtract
		a := uint16(gb.cpu.getRegister(regA))

		if gb.cpu.getFlag(FlagN) {
			// Previous operation was subtraction
			if gb.cpu.getFlag(FlagH) {
				// Underflow lower nybble
				a -= 0x6
			}
			if gb.cpu.getFlag(FlagC) {
				// Underflow upper nybble
				a -= 0x60
			}
		} else {
			// Previous operation was addition
			if (a&0xF) > 0x9 || gb.cpu.getFlag(FlagH) {
				// Overflow lower nybble
				a += 0x6
			}
			if a > 0x9F || gb.cpu.getFlag(FlagC) {
				// Overflow upper nybble
				a += 0x60
				gb.cpu.set_flag(FlagC, true)
			}
		}
		gb.cpu.set_flag(FlagZ, uint8(a) == 0)
		gb.cpu.set_flag(FlagH, false)
		gb.cpu.setRegister(regA, uint8(a))
		return 1
	case 0x2F:
		// CPL
		gb.cpu.setRegister(regA, ^gb.cpu.getRegister(regA))
		gb.cpu.set_flag(FlagN, true)
		gb.cpu.set_flag(FlagH, true)
		return 1
	case 0x3F:
		// CCF
		// Complement Carry Flag
		if gb.cpu.getFlag(FlagC) {
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
	case 0x10:
		// STOP
		// TODO: (future) Handle STOP correctly, treating it as halt for now
		gb.halted = true
		gb.popPC()
		return 1
	/////////////// Jumps ////////////////////
	case 0xC3:
		// JP nn
		gb.cpu.setRegister16(regPC, gb.popPC16())
		return 4
	case 0xC2:
		// JP NZ,nn
		adr := gb.popPC16()
		if !gb.cpu.getFlag(FlagZ) {
			gb.cpu.setRegister16(regPC, adr)
			return 4
		}
		return 3
	case 0xCA:
		// JP Z,nn
		adr := gb.popPC16()
		if gb.cpu.getFlag(FlagZ) {
			gb.cpu.setRegister16(regPC, adr)
			return 4
		}
		return 3
	case 0xD2:
		// JP NC,nn
		adr := gb.popPC16()
		if !gb.cpu.getFlag(FlagC) {
			gb.cpu.setRegister16(regPC, adr)
			return 4
		}
		return 3
	case 0xDA:
		// JP C,nn
		adr := gb.popPC16()
		if gb.cpu.getFlag(FlagC) {
			gb.cpu.setRegister16(regPC, adr)
			return 4
		}
		return 3
	case 0xE9:
		// JP HL
		gb.cpu.setRegister16(regPC, gb.cpu.getRegister16(regHL))
		return 1
	case 0x18:
		// JR n
		// Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		// Using value of PC after incrementing
		var pc int32 = int32(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, uint16(pc+n))
		return 3
	case 0x20:
		// JR NZ,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if !gb.cpu.getFlag(FlagZ) {
			var pc int32 = int32(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, uint16(pc+n))
			return 3
		}
		return 2
	case 0x28:
		// JR Z,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if gb.cpu.getFlag(FlagZ) {
			var pc int32 = int32(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, uint16(pc+n))
			return 3
		}
		return 2
	case 0x30:
		// JR NC,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if !gb.cpu.getFlag(FlagC) {
			var pc int32 = int32(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, uint16(pc+n))
			return 3
		}
		return 2
	case 0x38:
		// JR C,n
		// Conditional Relative Jump, n is a signed 8-bit immediate value
		var n int32 = int32(int8(gb.popPC()))
		if gb.cpu.getFlag(FlagC) {
			var pc int32 = int32(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, uint16(pc+n))
			return 3
		}
		return 2
	case 0xCD:
		// CALL nn
		nn := gb.popPC16()
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, nn)
		return 6
	case 0xC4:
		// CALL NZ,nn
		nn := gb.popPC16()
		if !gb.cpu.getFlag(FlagZ) {
			gb.pushToStack16(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, nn)
			return 6
		}
		return 3
	case 0xCC:
		// CALL Z,nn
		nn := gb.popPC16()
		if gb.cpu.getFlag(FlagZ) {
			gb.pushToStack16(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, nn)
			return 6
		}
		return 3
	case 0xD4:
		// CALL NC,nn
		nn := gb.popPC16()
		if !gb.cpu.getFlag(FlagC) {
			gb.pushToStack16(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, nn)
			return 6
		}
		return 3
	case 0xDC:
		// CALL C,nn
		nn := gb.popPC16()
		if gb.cpu.getFlag(FlagC) {
			gb.pushToStack16(gb.cpu.getRegister16(regPC))
			gb.cpu.setRegister16(regPC, nn)
			return 6
		}
		return 3
	case 0xC7:
		// RST 00H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x00)
		return 4
	case 0xCF:
		// RST 08H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x08)
		return 4
	case 0xD7:
		// RST 10H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x10)
		return 4
	case 0xDF:
		// RST 18H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x18)
		return 4
	case 0xE7:
		// RST 20H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x20)
		return 4
	case 0xEF:
		// RST 28H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x28)
		return 4
	case 0xF7:
		// RST 30H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x30)
		return 4
	case 0xFF:
		// RST 38H
		gb.pushToStack16(gb.cpu.getRegister16(regPC))
		gb.cpu.setRegister16(regPC, 0x38)
		return 4
	case 0xC9:
		// RET
		gb.cpu.setRegister16(regPC, gb.popFromStack())
		return 4
	case 0xC0:
		// RET NZ
		if !gb.cpu.getFlag(FlagZ) {
			gb.cpu.setRegister16(regPC, gb.popFromStack())
			return 5
		}
		return 2
	case 0xC8:
		// RET Z
		if gb.cpu.getFlag(FlagZ) {
			gb.cpu.setRegister16(regPC, gb.popFromStack())
			return 5
		}
		return 2
	case 0xD0:
		// RET NC
		if !gb.cpu.getFlag(FlagC) {
			gb.cpu.setRegister16(regPC, gb.popFromStack())
			return 5
		}
		return 2
	case 0xD8:
		// RET C
		if gb.cpu.getFlag(FlagC) {
			gb.cpu.setRegister16(regPC, gb.popFromStack())
			return 5
		}
		return 2
	case 0xD9:
		// RETI
		// Return and enable interrupts
		gb.cpu.setRegister16(regPC, gb.popFromStack())
		gb.interruptMasterEnable = true
		return 4
	////////////// CB - Extended Instructions /////////////
	case 0xCB:
		// CB
		next_opcode := gb.popPC()
		return gb.CBOpcode(next_opcode)
	default:
		// Intentionally unimplemented opcodes should crash Gameboy
		panic(fmt.Sprintf("opcode 0x%X not implemented", opcode))
	}
}
