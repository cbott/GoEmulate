package gameboy

import "github.com/cbott/GoEmulate/cartridges"

const (
	CpuSpeed        = 4194304                                      // Hz
	CyclesPerFrame  = (ScreenHeight + VBlankLines) * CyclesPerLine // 154 lines
	FramesPerSecond = float64(CpuSpeed) / CyclesPerFrame
)

type Gameboy struct {
	cpu    *CpuRegisters
	memory *Memory

	// Array of RGB triplets for each pixel on the Game Boy screen
	// This is filled in throughout the PPU processes and then displayed
	ScreenData [ScreenWidth][ScreenHeight][3]uint8

	// Scan Cycles tracks the number of cycles elapsed drawing the current frame so far
	currentScanCycles int
	// Internal counter to keep track of when the DIV register should increment
	timerAccumulator int
	// When halted CPU will not execute instructions except for interrupts
	halted bool
	// Interrupt Master Enable sets whether interrupts are enabled globally
	interruptMasterEnable bool
	// Some instructions set interrupt state with a 1 operation delay, these bools track that state
	pendingInterruptEnable bool
	// Storage for CPU save states
	savestates [NumSaveStates]*SaveState

	screenCleared  bool
	displayEnabled bool
	debugColors    bool
}

// Create and initialize a Game Boy struct
func NewGameBoy(skipboot bool, debugColors bool) *Gameboy {
	var gb = Gameboy{}
	gb.cpu = &CpuRegisters{}
	gb.memory = &Memory{}
	gb.memory.Init()

	if skipboot {
		gb.memory.BypassBootROM()
		gb.cpu.BypassBootROM()
	}

	gb.debugColors = debugColors

	return &gb
}

// Load an initialized Cartridge struct into Game Boy memory
func (gb *Gameboy) LoadCartridge(c cartridges.Cartridge) {
	gb.memory.cartridge = c
}

// Write cartridge RAM contents to the save file
func (gb *Gameboy) SaveCartridgeRAM() {
	gb.memory.cartridge.SaveRAM()
}

// RunNextFrame executes Game Boy processes up to the next complete frame to be displayed
func (gb *Gameboy) RunNextFrame() {
	var totalCycles int
	var lastTotalCycles int

	// Clear screen and restart rendering at dot zero
	// A bit of a hack to force display timing to match up perfectly with PPU process
	gb.clearScreen()
	// Set flag to resume rendering if we recently enabled the LCD
	gb.displayEnabled = true

	for totalCycles = 0; totalCycles < CyclesPerFrame; {
		var operationCycles int
		if gb.halted {
			// Behavior while halted is approximate, not really tested
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
		gb.memory.apu.RunAudioProcess(cyclesSinceLast)

		// Evaulate interrupt state after this round of graphics and timer updates
		totalCycles += gb.RunInterrupts()
	}
}
