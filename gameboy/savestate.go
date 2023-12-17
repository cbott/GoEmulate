package gameboy

import (
	"github.com/cbott/GoEmulate/cartridges"
)

// Number of distinct save states we will allow storing, arbitrary limit, more would just make the gameboy object larger
const NumSaveStates = 3

// SaveState holds a snapshot of the current Game Boy state, allowing it to be returned to later
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

// StoreState saves the current memory and CPU state to the internal storage array at index i
// if index falls outside the range 0 <= i < NumSaveStates the operation will be ignored
// Returns whether the state was successfully stored
func (gb *Gameboy) StoreState(i int) bool {
	if i < 0 || i >= NumSaveStates {
		return false
	}

	save := SaveState{}

	save.memory = gb.memory.memory
	save.divAccumulator = gb.memory.divAccumulator
	save.ButtonStates = gb.memory.buttonStates

	save.cpu = *gb.cpu
	save.currentScanCycles = gb.currentScanCycles
	save.timerAccumulator = gb.timerAccumulator
	save.halted = gb.halted
	save.interruptMasterEnable = gb.interruptMasterEnable
	save.pendingInterruptEnable = gb.pendingInterruptEnable
	save.screenCleared = gb.screenCleared
	save.displayEnabled = gb.displayEnabled

	save.ram, save.ramBank, save.ramEnabled, save.romBank = gb.memory.cartridge.GetState()

	gb.savestates[i] = &save
	return true
}

// RecallState overwrites the current memory and CPU state with values previously stored at index i with StoreState
// if index has no previously saved state, the operation will be ignored
// Returns whether the state was successfully restored
func (gb *Gameboy) RecallState(i int) bool {
	if i < 0 || i >= NumSaveStates {
		return false
	}
	if gb.savestates[i] == nil {
		return false
	}

	var state *SaveState = gb.savestates[i]

	gb.memory.memory = state.memory
	gb.memory.divAccumulator = state.divAccumulator
	gb.memory.buttonStates = state.ButtonStates

	gb.cpu = &state.cpu
	gb.currentScanCycles = state.currentScanCycles
	gb.timerAccumulator = state.timerAccumulator
	gb.halted = state.halted
	gb.interruptMasterEnable = state.interruptMasterEnable
	gb.pendingInterruptEnable = state.pendingInterruptEnable
	gb.screenCleared = state.screenCleared
	gb.displayEnabled = state.displayEnabled

	gb.memory.cartridge.SetState(state.ram, state.ramBank, state.ramEnabled, state.romBank)

	return true
}
