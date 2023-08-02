package main

var cbRegisterOrder = []register8{
	regB,
	regC,
	regD,
	regE,
	regH,
	regL,
	nil, // (HL)
	regA,
}

var cbInstructionOrder = []func(*Gameboy, register8){
	cbRLC, cbRRC,
	cbRL, cbRR,
	cbSLA, cbSRA,
	cbSWAP, cbSRL,
	cbPartial(cbBIT, 0), cbPartial(cbBIT, 1),
	cbPartial(cbBIT, 2), cbPartial(cbBIT, 3),
	cbPartial(cbBIT, 4), cbPartial(cbBIT, 5),
	cbPartial(cbBIT, 6), cbPartial(cbBIT, 7),
	cbPartial(cbRES, 0), cbPartial(cbRES, 1),
	cbPartial(cbRES, 2), cbPartial(cbRES, 3),
	cbPartial(cbRES, 4), cbPartial(cbRES, 5),
	cbPartial(cbRES, 6), cbPartial(cbRES, 7),
	cbPartial(cbSET, 0), cbPartial(cbSET, 1),
	cbPartial(cbSET, 2), cbPartial(cbSET, 3),
	cbPartial(cbSET, 4), cbPartial(cbSET, 5),
	cbPartial(cbSET, 6), cbPartial(cbSET, 7),
}

// Wrapper on get_register to account for (HL) access
func cbRegisterGet(gb *Gameboy, reg register8) uint8 {
	if reg == nil {
		// (HL)
		return gb.memory.get(gb.cpu.get_register16(regHL))
	} else {
		return gb.cpu.get_register(reg)
	}
}

// Wrapper on set_register to account for (HL) access
func cbRegisterSet(gb *Gameboy, reg register8, value uint8) {
	if reg == nil {
		// (HL)
		gb.memory.set(gb.cpu.get_register16(regHL), value)
	} else {
		gb.cpu.set_register(reg, value)
	}
}

// Partial() allows us to pre-specify the bit argument for BIT/RES/SET calls
func cbPartial(f func(*Gameboy, register8, uint8), bit uint8) func(*Gameboy, register8) {
	partialfunc := func(gb *Gameboy, reg register8) {
		f(gb, reg, bit)
	}
	return partialfunc
}

// Rotate Left and set Carry bit
func cbRLC(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	carrybit := value >> 7
	result := (value << 1) | carrybit

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, carrybit == 1)
}

// Rotate Right and set Carry bit
func cbRRC(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	carrybit := value & 1
	result := (value >> 1) | (carrybit << 7)

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, carrybit == 1)
}

// Rotate Left through carry bit
func cbRL(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	var oldcarry uint8 = 0
	if gb.cpu.get_flag(FlagC) {
		oldcarry = 1
	}
	newcarry := value >> 7
	result := (value << 1) | oldcarry

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, newcarry == 1)
}

// Rotate Right through carry bit
func cbRR(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	var oldcarry uint8 = 0
	if gb.cpu.get_flag(FlagC) {
		oldcarry = 1
	}
	newcarry := value & 1
	result := (value >> 1) | (oldcarry << 7)

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, newcarry == 1)
}

// Shift Left into carry, set lsb to 0
func cbSLA(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	carrybit := value >> 7
	result := value << 1

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, carrybit == 1)
}

// Shift Right into carry, do not change msb
func cbSRA(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	carrybit := value & 1
	result := (value & 0b10000000) | (value >> 1)

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, carrybit == 1)
}

// Swap low and high nybbles
func cbSWAP(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	result := (value << 4) | (value >> 4)

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, false)
}

// Shift Right into carry, set msb to 0
func cbSRL(gb *Gameboy, reg register8) {
	value := cbRegisterGet(gb, reg)

	carrybit := value & 1
	result := value >> 1

	cbRegisterSet(gb, reg, result)
	gb.cpu.set_flag(FlagZ, result == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, false)
	gb.cpu.set_flag(FlagC, carrybit == 1)
}

// Test bit in register
func cbBIT(gb *Gameboy, reg register8, bit uint8) {
	gb.cpu.set_flag(FlagZ, cbRegisterGet(gb, reg)&(1<<bit) == 0)
	gb.cpu.set_flag(FlagN, false)
	gb.cpu.set_flag(FlagH, true)
}

// Reset bit in register
func cbRES(gb *Gameboy, reg register8, bit uint8) {
	cbRegisterSet(gb, reg, cbRegisterGet(gb, reg) & ^(1<<bit))
}

// Set bit in register
func cbSET(gb *Gameboy, reg register8, bit uint8) {
	cbRegisterSet(gb, reg, cbRegisterGet(gb, reg)|(1<<bit))
}

// Execute a single CB-Prefixed opcode and return the number of CPU cycles it took (1MHz CPU cycles)
func (gb *Gameboy) CBOpcode(opcode uint8) int {
	// CB instructions have a regular pattern so we can avoid hard coding things
	column := (opcode & 0xF) % 8
	reg := cbRegisterOrder[column]
	function := cbInstructionOrder[opcode/8]

	function(gb, reg)

	if opcode == 0x4e || opcode == 0x5e || opcode == 0x6e || opcode == 0x7e {
		// TODO: Several emulators I've looked at call these operations 3 cycles
		// but most documentation says 4, need to pick the right one and implement cleanly
		return 3
	}
	if reg == nil {
		// (HL)
		return 4
	}
	return 2
}
