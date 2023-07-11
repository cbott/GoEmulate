package main

import (
	"fmt"
	"os"
)

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

	// Internal counter to keep track of when the DIV register should increment
	timerAccumulator int

	// When halted CPU will not execute instructions except for interrupts
	halted bool
	// Interrupt Master Enable sets whether interrupts are enabled globally
	interruptMasterEnable bool
	// Some instructions set interrupt state with a 1 operation delay, these bools track that state
	pendingInterruptEnable bool

	// Track button presses so we can raise an interrupt on button press
	joypadState uint8

	screenCleared bool
	debugCounter  int32
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
		var operationCycles int
		if gb.halted {
			// TODO: verify this is the correct behavior while halted
			operationCycles = 4
		} else {
			operationCycles = gb.RunNextOpcode()
		}

		gb.RunGraphicsProcess(operationCycles)
		gb.RunTimers(operationCycles)

		// TODO: see if we have to count these extra cycles for other stuff?
		totalCycles += operationCycles
		totalCycles += gb.RunInterrupts()
	}
}

// Returns the number of clock cycles to complete (4MHz cycles)
func (gb *Gameboy) RunNextOpcode() int {
	gb.debugCounter++

	// if gb.debugCounter == 16508 {
	// 	fmt.Printf("debug")
	// }

	if gb.debugCounter < 0 {
		f, err := os.OpenFile("gb_results.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("unable to open log")
		}
		defer f.Close()
		fmt.Fprintf(f, "A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X PCMEM:%02X,%02X,%02X,%02X\n",
			gb.cpu.A, gb.cpu.F, gb.cpu.B, gb.cpu.C, gb.cpu.D, gb.cpu.E, gb.cpu.H, gb.cpu.L,
			gb.cpu.SP, gb.cpu.PC, gb.memory.get(gb.cpu.PC), gb.memory.get(gb.cpu.PC+1),
			gb.memory.get(gb.cpu.PC+2), gb.memory.get(gb.cpu.PC+3))
	}

	opcode := gb.popPC()
	return gb.Opcode(opcode) * 4
}
