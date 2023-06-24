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
F: Flag register
SP: Stack Pointer
PC: Program Counter
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

// TODO: this is not the best way to deal with flags
// probably delete it as a standard register and replace it with a struct of bools or similar
const (
	// 7 (Z): Zero Flag
	FlagZ = 7
	// 6 (N): Subtract Flag (BDC)
	FlagN = 6
	// 5 (H): Half Carry Flag (BCD)
	FlagH = 5
	// 4 (C): Carry Flag
	FlagC = 4
	// 3-0: Unused
	//Values marked BCD are used for Binary-Coded Decimal operations only
)

// Write a value to the flag register
func (r *CpuRegisters) set_flag(flag uint8, value bool) {
	if value {
		// set
		r.F |= 1 << flag

	} else {
		// clear
		r.F &= (0xFF ^ (1 << flag))
	}
}

// Read a value from the flag register
func (r *CpuRegisters) get_flag(flag uint8) bool {
	return (r.F & (1 << flag)) != 0
}

// TODO: we could just directly access these?
// TODO: should we be using camelCase instead of underscores in fn names?
// TODO: "set" is kind of an overloaded word, usually indicates writing a 1, "write" is better?

// Set an 8-bit register
func (r *CpuRegisters) set_register(letter string, value uint8) error {
	switch letter {
	case "A":
		r.A = value
	case "F":
		// Lower 4 bits of F cannot be set
		r.F = value & 0xF0
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

// Set a 16-bit register
func (r *CpuRegisters) set_register16(letter string, value uint16) error {
	switch letter {
	case "AF":
		r.A = uint8(value >> 8)
		// Lower 4 bits of F cannot be set
		r.F = uint8(value & 0xF0)
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

// Read an 8-bit register
func (r *CpuRegisters) get_register(letter string) uint8 {
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

// Read a 16-bit register
func (r *CpuRegisters) get_register16(letter string) uint16 {
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

// Add n to A and set the appropriate flags for the result,
// if carry is set to true the value of the carry flag will be added as well
func (r *CpuRegisters) addToRegisterA(n uint8, carry bool) {
	a := r.get_register("A")
	var carrybit uint8 = 0
	if carry && r.get_flag(FlagC) {
		carrybit = 1
	}
	result := a + n + carrybit
	r.set_flag(FlagZ, result == 0)
	r.set_flag(FlagN, false)
	r.set_flag(FlagH, ((a&0xF)+(n&0xF)+carrybit) > 0xF)
	r.set_flag(FlagC, a > result)
	r.set_register("A", result)
}

// Subtract n from A and set the appropriate flags for the result
func (r *CpuRegisters) subtractFromRegisterA(n uint8, carry bool) {
	a := r.get_register("A")
	var carrybit uint8 = 0
	if carry && r.get_flag(FlagC) {
		carrybit = 1
	}
	var signedResult int16 = int16(a) - int16(n) - int16(carrybit)
	var signedHalfResult int16 = int16(a&0xF) - int16(n&0xF) - int16(carrybit)
	result := uint8(signedResult)

	r.set_flag(FlagZ, result == 0)
	r.set_flag(FlagN, true)
	r.set_flag(FlagH, signedHalfResult < 0)
	r.set_flag(FlagC, signedResult < 0)
	r.set_register("A", result)
}

// Perform bitwise AND with register A and store the result in A
func (r *CpuRegisters) andA(n uint8) {
	a := r.get_register("A")
	result := a & n
	r.set_flag(FlagZ, result == 0)
	r.set_flag(FlagN, false)
	r.set_flag(FlagH, true)
	r.set_flag(FlagC, false)
	r.set_register("A", result)
}

// Perform bitwise OR with register A and store the result in A
func (r *CpuRegisters) orA(n uint8) {
	a := r.get_register("A")
	result := a | n
	r.set_flag(FlagZ, result == 0)
	r.set_flag(FlagN, false)
	r.set_flag(FlagH, false)
	r.set_flag(FlagC, false)
	r.set_register("A", result)
}

// Perform bitwise XOR with register A and store the result in A
func (r *CpuRegisters) xorA(n uint8) {
	a := r.get_register("A")
	result := a ^ n
	r.set_flag(FlagZ, result == 0)
	r.set_flag(FlagN, false)
	r.set_flag(FlagH, false)
	r.set_flag(FlagC, false)
	r.set_register("A", result)
}

// Compare n with register A and set the approprate flags based on the result
func (r *CpuRegisters) compareA(n uint8) {
	a := r.get_register("A")
	r.set_flag(FlagZ, a == n)
	r.set_flag(FlagN, true)
	r.set_flag(FlagH, (a&0xF) < (n&0xF))
	r.set_flag(FlagC, a < n)
}

// Increment an 8 bit register and set the appropriate flags for the result,
func (r *CpuRegisters) incrementRegister(letter string) {
	value := r.get_register(letter)
	r.set_flag(FlagZ, value == 0xFF)
	r.set_flag(FlagN, false)
	r.set_flag(FlagH, (value&0xF) == 0xF)
	r.set_register(letter, value+1)
}

// Decrement an 8 bit register and set the appropriate flags for the result,
func (r *CpuRegisters) decrementRegister(letter string) {
	value := r.get_register(letter)
	r.set_flag(FlagZ, value == 0x01)
	r.set_flag(FlagN, true)
	r.set_flag(FlagH, (value&0xF) == 0x0)
	r.set_register(letter, value-1)
}

// Add n to HL and set the appropriate flags for the result
func (r *CpuRegisters) addToRegisterHL(n uint16) {
	hl := r.get_register16("HL")
	result := hl + n
	r.set_flag(FlagN, false)
	r.set_flag(FlagH, (hl&0x0FFF)+(n&0x0FFF) > 0x0FFF)
	r.set_flag(FlagC, hl > result)
	r.set_register16("HL", result)
}

// Set registers to default Game Boy startup values
func (c *CpuRegisters) Init() {
	c.set_register16("AF", 0x01B0)
	c.set_register16("BC", 0x0013)
	c.set_register16("DE", 0x00D8)
	c.set_register16("HL", 0x014D)
	c.set_register16("SP", 0xFFFE)
	c.set_register16("PC", 0x0100)
}
