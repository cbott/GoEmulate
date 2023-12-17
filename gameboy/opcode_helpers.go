package gameboy

// Returns the number of clock cycles to complete (4MHz cycles)
func (gb *Gameboy) RunNextOpcode() int {
	opcode := gb.popPC()
	return gb.Opcode(opcode) * 4
}

// Return the 8 bit value in memory at address (PC) and then increment PC
func (gb *Gameboy) popPC() uint8 {
	pc := gb.cpu.getRegister16(regPC)
	gb.cpu.setRegister16(regPC, pc+1)
	return gb.memory.get(pc)
}

// Read the 16 bit value in memory at address (PC, PC+1) and increment PC twice
func (gb *Gameboy) popPC16() uint16 {
	lsb := uint16(gb.popPC())
	msb := uint16(gb.popPC())
	return msb<<8 | lsb
}

// Push a 16 bit value onto the stack as two separate parts and update the stack pointer
func (gb *Gameboy) pushToStack(high uint8, low uint8) {
	sp := gb.cpu.getRegister16(regSP)
	gb.memory.set(sp-1, high)
	gb.memory.set(sp-2, low)
	// Decrement stack pointer twice
	gb.cpu.setRegister16(regSP, sp-2)
}

// Push a single 16 bit value onto the stack and update the stack pointer
func (gb *Gameboy) pushToStack16(value uint16) {
	gb.pushToStack(uint8(value>>8), uint8(value&0xFF))
}

// Pop a 16 bit value off of the stack and update the stack pointer
func (gb *Gameboy) popFromStack() uint16 {
	sp := gb.cpu.getRegister16(regSP)
	low := uint16(gb.memory.get(sp))
	high := uint16(gb.memory.get(sp + 1))
	// Increment stack pointer twice
	gb.cpu.setRegister16(regSP, sp+2)
	return (high << 8) | low
}

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
