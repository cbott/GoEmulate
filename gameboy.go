package main

const FramesPerSecond = 60.0
const CyclesPerFrame = CpuSpeed / int(FramesPerSecond)

type Gameboy struct {
	cpu    *CpuRegisters
	memory *Memory

	// TODO: re-evaluate this matrix
	// Matrix of pixel data which is used while the screen is rendering. When a
	// frame has been completed, this data is copied into the PreparedData matrix.
	screenData [ScreenWidth][ScreenHeight][3]uint8

	// PreparedData is a matrix of screen pixel data for a single frame which has
	// been fully rendered.
	PreparedData [ScreenWidth][ScreenHeight][3]uint8

	// scan is the horizontal "position" in clock cycles where the PPU is currently drawing
	// TODO: should this be stored in the PPU object if we make one?
	currentScanX int

	// Interrupt Master Enable sets whether interrupts are enabled globally
	interruptMasterEnable bool
}

// TODO: Move all *Gameboy functions to a separate file

// Return the value in memory pointed to by PC and then increment PC
func (gb *Gameboy) popPC() uint8 {
	pc := gb.cpu.get_register16("PC")
	gb.cpu.set_register16("PC", pc+1)
	return gb.memory.get(pc)
}

// Push a 16 bit value onto the stack as two separate parts and update the stack pointer
func (gb *Gameboy) pushToStack(high uint8, low uint8) {
	sp := gb.cpu.get_register16("SP")
	gb.memory.set(sp-1, high)
	gb.memory.set(sp-2, low)
	// Decrement stack pointer twice
	gb.cpu.set_register16("SP", sp-2)
}

// Push a single 16 bit value onto the stack and update the stack pointer
func (gb *Gameboy) pushToStack16(value uint16) {
	gb.pushToStack(uint8(value>>8), uint8(value&0xFF))
}

// Pop a 16 bit value off of the stack and update the stack pointer
func (gb *Gameboy) popFromStack() uint16 {
	sp := gb.cpu.get_register16("SP")
	low := uint16(gb.memory.get(sp))
	high := uint16(gb.memory.get(sp + 1))
	// Increment stack pointer twice
	gb.cpu.set_register16("SP", sp+2)
	return (high << 8) | low
}

func (gb *Gameboy) RunNextFrame() {
	// Run Gameboy processes up to the next complete frame to be displayed
	for totalCycles := 0; totalCycles < CyclesPerFrame; {
		operationCycles := gb.RunNextOpcode()
		gb.RunGraphicsProcess(operationCycles)
		// TODO: handle hardware timers
		totalCycles += operationCycles
		// TODO: totalCycles += run interrupt service routines
	}
}

func (gb *Gameboy) RunNextOpcode() int {
	// Returns the number of clock cycles to complete (4MHz cycles)
	opcode := gb.popPC()
	return gb.Opcode(opcode) * 4
}