package gameboy

import "github.com/cbott/GoEmulate/cartridges"

type SaveState struct {
	cpu                    CpuRegisters
	currentScanCycles      int
	timerAccumulator       int
	halted                 bool
	interruptMasterEnable  bool
	pendingInterruptEnable bool
	screenCleared          bool
	displayEnabled         bool

	memory         [0x10000]uint8
	divAccumulator int
	ButtonStates   uint8

	// Cartridge state
	ram        [][cartridges.RAMBankSize]uint8
	ramBank    uint8
	ramEnabled bool
	romBank    uint16
}

func NewSaveState(gb *Gameboy) *SaveState {
	save := SaveState{}

	save.memory = gb.memory.memory
	save.divAccumulator = gb.memory.divAccumulator
	save.ButtonStates = gb.memory.ButtonStates

	save.cpu = *gb.cpu
	save.currentScanCycles = gb.currentScanCycles
	save.timerAccumulator = gb.timerAccumulator
	save.halted = gb.halted
	save.interruptMasterEnable = gb.interruptMasterEnable
	save.pendingInterruptEnable = gb.pendingInterruptEnable
	save.screenCleared = gb.screenCleared
	save.displayEnabled = gb.displayEnabled

	save.ram, save.ramBank, save.ramEnabled, save.romBank = gb.memory.cartridge.GetState()

	return &save
}

func RestoreState(gb *Gameboy, state *SaveState) {
	// Making a copy so we don't mess with the real save state
	var referenceState SaveState = *state

	gb.memory.memory = referenceState.memory
	gb.memory.divAccumulator = state.divAccumulator
	gb.memory.ButtonStates = state.ButtonStates

	gb.cpu = &referenceState.cpu
	gb.currentScanCycles = referenceState.currentScanCycles
	gb.timerAccumulator = referenceState.timerAccumulator
	gb.halted = referenceState.halted
	gb.interruptMasterEnable = referenceState.interruptMasterEnable
	gb.pendingInterruptEnable = referenceState.pendingInterruptEnable
	gb.screenCleared = referenceState.screenCleared
	gb.displayEnabled = referenceState.displayEnabled

	gb.memory.cartridge.SetState(state.ram, state.ramBank, state.ramEnabled, state.romBank)
}
