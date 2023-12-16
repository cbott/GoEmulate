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
// TODO: probably a better way to do this
type ButtonState struct {
	BtnA, BtnB, BtnSelect, BtnStart   bool
	BtnRight, BtnLeft, BtnUp, BtnDown bool
}

// SetButtonState allows an external caller to input the current joypad values into the emulator
// should we pass these in as 8 bools? that's the most clear
// could do uint8 but that seems impossible to understand
// maybe a struct with 8 fields? or dictionary?
// Update the Joypad button states with what is currently pressed on the keyboard
// and requests an interrupt if a button just became pressed
func (gb *Gameboy) SetButtonStates(state *ButtonState) {
	// Gameboy reads values from register as 1=unpressed, 0=pressed
	// we invert at the end so bit operations are easier
	var reg uint8 = 0
	// TODO: ideally we move the bit shift constants somewhere else, like a map or const
	if state.BtnDown {
		reg |= 1 << 7
	}
	if state.BtnUp {
		reg |= 1 << 6
	}
	if state.BtnLeft {
		reg |= 1 << 5
	}
	if state.BtnRight {
		reg |= 1 << 4
	}
	if state.BtnStart {
		reg |= 1 << 3
	}
	if state.BtnSelect {
		reg |= 1 << 2
	}
	if state.BtnB {
		reg |= 1 << 1
	}
	if state.BtnA {
		reg |= 1 << 0
	}
	reg = ^reg

	// Perform an interrupt if any button went from unpressed to pressed
	var doInterrupt bool = false
	for i := 0; i < 8; i++ {
		if (gb.memory.ButtonStates&(1<<i) != 0) && (reg&(1<<i) == 0) {
			doInterrupt = true
		}
	}

	// TODO: rethink this variable, at least make it private
	gb.memory.ButtonStates = reg

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
		state |= ^((m.ButtonStates >> 4) & 0xF)
	}
	// If both action and direction buttons are selected we will set the pair as pressed if either is pressed
	if p1&JOYPAD_action_buttons == 0 {
		state |= ^(m.ButtonStates & 0xF)
	}

	p1 |= (^state) & 0xF

	return p1
}
