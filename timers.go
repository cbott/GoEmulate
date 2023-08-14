// timers.go
// Handle Game Boy timer and divider registers

package main

const (
	DIV              = 0xFF04 // Divider register
	TIMA             = 0xFF05 // Timer counter
	TMA              = 0xFF06 // Timer modulo
	TAC              = 0xFF07 // Timer control
	TAC_timer_enable = 1 << 2
	// Bit  2   - Timer Enable
	// Bits 1-0 - Input Clock Select
	//            00: CPU Clock / 1024 ->   4096 Hz
	//            01: CPU Clock / 16   -> 262144 Hz
	//            10: CPU Clock / 64   ->  65536 Hz
	//            11: CPU Clock / 256  ->  16384 Hz
)

// Get timer divider ratio specified by the timer control register
func getClockSelect(tac_value uint8) int {
	switch tac_value & 0b11 {
	case 0b00:
		return 1024
	case 0b01:
		return 16
	case 0b10:
		return 64
	default:
		return 256
	}
}

// Advance hardware timers by the specified number of machine cycles (4MHz)
func (gb *Gameboy) RunTimers(cycles int) {
	// DIV is incremented at a rate of 16384Hz (1/256 of the clock rate)
	gb.memory.divAccumulator += cycles
	// Access memory directly as calls to set() will reset DIV
	// We're overcomplicating things a bit here with the % stuff, probably safe to assume
	// we call this function often enough that we'll never increment by more than 1 at a time
	gb.memory.memory[DIV] = gb.memory.memory[DIV] + uint8(gb.memory.divAccumulator/256)
	gb.memory.divAccumulator %= 256

	tac_value := gb.memory.get(TAC)
	if tac_value&TAC_timer_enable != 0 {
		// Timer is enabled
		gb.timerAccumulator += cycles
		clock := getClockSelect(tac_value)
		// Note: This may accumulate a small error when switching clock rates if timerAccumulator wasn't 0
		// likely irrelevant for most use cases
		var incrementedTIMA int = int(gb.memory.get(TIMA)) + (gb.timerAccumulator / clock)
		gb.timerAccumulator %= clock

		var newTIMA uint8
		if incrementedTIMA > 0xFF {
			// TIMA has overflowed, reset to the value in TMA and request an interrupt
			newTIMA = gb.memory.get(TMA)
			gb.SetInterruptRequestFlag(Interrupt_timer)
		} else {
			// Otherwise just increment TIMA
			newTIMA = uint8(incrementedTIMA)
		}

		gb.memory.set(TIMA, newTIMA)
	}
}
