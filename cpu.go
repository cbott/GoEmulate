package main

import (
	"fmt"
)

/*
CPU Registers

 A | F
------
 B | C
------
 D | E
------
 H | L
------
  SP
------
  PC
------

A: Arithmetic register for logic and math
F: Flags register (Zero, Carry)
PC: Program counter
*/

const CpuSpeed = 4194304 // Hz

// Define registers
type CpuRegisters struct {
	A, F uint8
	B, C uint8
	D, E uint8
	H, L uint8
	SP   uint16
	PC   uint16 // TODO: Maybe don't have PC treated as a register
}

// TODO: we could just directly access these?
func (r *CpuRegisters) set_register(letter string, value uint8) error {
	// Set an 8-bit register
	switch letter {
	case "A":
		r.A = value
	case "F":
		r.F = value
	case "B":
		r.B = value
	case "C":
		r.C = value
	case "D":
		r.D = value
	case "E":
		r.E = value
	case "H":
		r.H = value
	case "L":
		r.L = value
	default:
		return fmt.Errorf("tried to set nonexistent 8 bit register %s", letter)
	}
	return nil
}

func (r *CpuRegisters) set_register16(letter string, value uint16) error {
	// Set a 16-bit register
	switch letter {
	case "AF":
		r.A = uint8(value >> 8)
		r.F = uint8(value & 0xFF)
	case "BC":
		r.B = uint8(value >> 8)
		r.C = uint8(value & 0xFF)
	case "DE":
		r.D = uint8(value >> 8)
		r.E = uint8(value & 0xFF)
	case "HL":
		r.H = uint8(value >> 8)
		r.L = uint8(value & 0xFF)
	case "SP":
		r.SP = value
	case "PC":
		r.PC = value
	default:
		return fmt.Errorf("tried to set nonexistent 16 bit register %s", letter)
	}
	return nil
}

func (r *CpuRegisters) get_register(letter string) uint8 {
	// Read an 8-bit register
	switch letter {
	case "A":
		return r.A
	case "F":
		return r.F
	case "B":
		return r.B
	case "C":
		return r.C
	case "D":
		return r.D
	case "E":
		return r.E
	case "H":
		return r.H
	case "L":
		return r.L
	default:
		panic(fmt.Sprintf("tried to get nonexistent 8-bit register %s", letter))
	}
}

func (r *CpuRegisters) get_register16(letter string) uint16 {
	// Read a 16-bit register
	switch letter {
	case "AF":
		return uint16(r.A)<<8 | uint16(r.F)
	case "BC":
		return uint16(r.B)<<8 | uint16(r.C)
	case "DE":
		return uint16(r.D)<<8 | uint16(r.E)
	case "HL":
		return uint16(r.H)<<8 | uint16(r.L)
	case "SP":
		return r.SP
	case "PC":
		return r.PC
	default:
		panic(fmt.Sprintf("tried to get nonexistent 16-bit register %s", letter))
	}
}

func (c *CpuRegisters) Init() {
	// Set registers to default Game Boy startup values
	c.set_register16("AF", 0x01B0)
	c.set_register16("BC", 0x0013)
	c.set_register16("DE", 0x00D8)
	c.set_register16("HL", 0x014D)
	c.set_register16("SP", 0xFFFE)
	c.set_register16("PC", 0x0100)
}
