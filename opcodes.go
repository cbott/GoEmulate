package main

import "fmt"

// Opcodes data from https://pastraiser.com/cpu/gameboy/gameboy_opcodes.html
// and http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf

func (gb *Gameboy) Opcode(opcode uint8) int {
	// Execute a single opcode and return the number of CPU cycles it took (1MHz CPU cycles)
	// TODO: standardize whicy type of cycle we're talking about
	switch opcode {
	//////////////// 8-bit loads ////////////////
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
	case 0x58:
		// LD E,B
		gb.cpu.set_register("E", gb.cpu.get_register("B"))
		return 1
	case 0x59:
		// LD E,C
		gb.cpu.set_register("E", gb.cpu.get_register("C"))
		return 1
	case 0x5A:
		// LD E,E
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
	//////////////// 16-bit loads ////////////////
	case 0x01:
		// LD BC,nn
		gb.cpu.set_register("C", gb.popPC()) // low byte
		gb.cpu.set_register("B", gb.popPC()) // high byte
		return 3
	case 0x00:
		// NOOP
		return 1
	default:
		panic(fmt.Sprintf("opcode %X not implemented", opcode))
	}
}
