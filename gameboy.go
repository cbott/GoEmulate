package main

import (
	"fmt"
	"os"
)

const (
	CpuSpeed        = 4194304 // Hz
	FramesPerSecond = 60.0
	CyclesPerFrame  = CpuSpeed / int(FramesPerSecond)
)

type Gameboy struct {
	cpu    *CpuRegisters
	memory *Memory

	// Array of RGB triplets for each pixel on the Game Boy screen
	// This is filled in throughout the PPU processes and then displayed
	screenData [ScreenWidth][ScreenHeight][3]uint8

	// scan is the horizontal "position" in clock cycles where the PPU is currently drawing
	currentScanX int
	// Internal counter to keep track of when the DIV register should increment
	timerAccumulator int
	// When halted CPU will not execute instructions except for interrupts
	halted bool
	// Interrupt Master Enable sets whether interrupts are enabled globally
	interruptMasterEnable bool
	// Some instructions set interrupt state with a 1 operation delay, these bools track that state
	pendingInterruptEnable bool

	screenCleared bool

	// TODO: (debug) delete debug counter
	debugCounter int32
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

func (gb *Gameboy) RunNextFrame() {
	// Run Gameboy processes up to the next complete frame to be displayed
	var totalCycles int
	var lastTotalCycles int

	for totalCycles = 0; totalCycles < CyclesPerFrame; {
		var operationCycles int
		if gb.halted {
			// TODO: verify this is the correct behavior while halted
			operationCycles = 4
		} else {
			operationCycles = gb.RunNextOpcode()
		}

		totalCycles += operationCycles

		// We want to run these processes for the time it took to do the current opcode
		// plus any time we spent on interrupts the last loop
		cyclesSinceLast := totalCycles - lastTotalCycles
		lastTotalCycles = totalCycles

		gb.RunGraphicsProcess(cyclesSinceLast)
		gb.RunTimers(cyclesSinceLast)

		// Evaulate interrupt state after this round of graphics and timer updates
		totalCycles += gb.RunInterrupts()
	}
}

// Returns the number of clock cycles to complete (4MHz cycles)
func (gb *Gameboy) RunNextOpcode() int {
	gb.debugCounter++

	if gb.debugCounter == -1 {
		fmt.Printf("debug\n")
	}

	if gb.debugCounter < -4000000 {
		f, err := os.OpenFile("gb_results.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("unable to open log")
		}
		defer f.Close()
		fmt.Fprintf(f, "X:%v A:%02X F:%02X B:%02X C:%02X D:%02X E:%02X H:%02X L:%02X SP:%04X PC:%04X MEM:%02X,%02X,%02X,%02X\n",
			gb.currentScanX, gb.cpu.A, gb.cpu.F, gb.cpu.B, gb.cpu.C, gb.cpu.D, gb.cpu.E, gb.cpu.H, gb.cpu.L,
			gb.cpu.SP, gb.cpu.PC, gb.memory.get(gb.cpu.PC), gb.memory.get(gb.cpu.PC+1),
			gb.memory.get(gb.cpu.PC+2), gb.memory.get(gb.cpu.PC+3))
	}

	opcode := gb.popPC()
	return gb.Opcode(opcode) * 4
}
