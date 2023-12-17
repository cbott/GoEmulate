package gameboy

/*
Joypad inputs start at memory address FF00

P14 | P15 (0 = select)
------------
 v  | Start  | P13 (0 = pressed)
 ^  | Select | P12
 <  | B      | P11
 >  | A      | P10

*/

const (
	JOYPAD                   = 0xFF00
	JOYPAD_action_buttons    = 1 << 5
	JOYPAD_direction_buttons = 1 << 4
)

// Convenience for passing in current joypad button status
type ButtonState struct {
	BtnA, BtnB, BtnSelect, BtnStart   bool
	BtnRight, BtnLeft, BtnUp, BtnDown bool
}

// SetButtonState sets the state of the Game Boy joypad buttons for the emulator
// and requests an interrupt if a button just became pressed
func (gb *Gameboy) SetButtonStates(state *ButtonState) {
	// Gameboy reads values from register as 1=unpressed, 0=pressed
	// we invert at the end so bit operations are easier
	var reg uint8 = 0
	for index, element := range []bool{state.BtnA, state.BtnB, state.BtnSelect, state.BtnStart,
		state.BtnRight, state.BtnLeft, state.BtnUp, state.BtnDown} {
		if element {
			reg |= 1 << index
		}
	}
	reg = ^reg

	// Perform an interrupt if any button went from unpressed to pressed
	var doInterrupt bool = false
	for i := 0; i < 8; i++ {
		if (gb.memory.buttonStates&(1<<i) != 0) && (reg&(1<<i) == 0) {
			doInterrupt = true
		}
	}

	gb.memory.buttonStates = reg

	if doInterrupt {
		// For maximum parity with hardware we should only trigger this if the particular input
		// row is enabled (P14/P15) but since we only run this function once prior to each frame
		// and select bits can be changed during the frame we will assume games only check this
		// interrupt if they are also enabling the inputs
		gb.SetInterruptRequestFlag(Interrupt_joypad)
	}
}

// Given register P1 with select bits set, return value of P1 with button state bits set as well
func (m *Memory) GetP1Value() uint8 {
	// Read the current P1 register value
	p1 := m.memory[JOYPAD]

	// Set top 2 bits (unused) and clear bottom 4 bits (button states)
	p1 = (p1 & 0xF0) | 0xC0

	// Fill in the bottom 4 bits with button states
	var state uint8
	if p1&JOYPAD_direction_buttons == 0 {
		state |= ^((m.buttonStates >> 4) & 0xF)
	}
	// If both action and direction buttons are selected we will set the pair as pressed if either is pressed
	if p1&JOYPAD_action_buttons == 0 {
		state |= ^(m.buttonStates & 0xF)
	}

	p1 |= (^state) & 0xF

	return p1
}
