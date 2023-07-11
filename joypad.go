package main

import "github.com/faiface/pixel/pixelgl"

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

// Button-to-keyboard mapping
const (
	BUTTON_DOWN   = pixelgl.KeyDown
	BUTTON_UP     = pixelgl.KeyUp
	BUTTON_LEFT   = pixelgl.KeyLeft
	BUTTON_RIGHT  = pixelgl.KeyRight
	BUTTON_START  = pixelgl.KeyEnter
	BUTTON_SELECT = pixelgl.KeyS
	BUTTON_B      = pixelgl.KeyB
	BUTTON_A      = pixelgl.KeyA
)

// Interpret keyboard presses as Game Boy joypad inputs
func (gb *Gameboy) keyboardToJoypad(pixelWindow *pixelgl.Window) {
	// TODO: this whole function is not very clean right now
	joypad_register := gb.memory.get(JOYPAD) & 0xF0
	var state uint8 = 0b1111

	if joypad_register&JOYPAD_direction_buttons == 0 {
		// Read direction pad
		if pixelWindow.Pressed(BUTTON_RIGHT) {
			state &= 0b1110
		}
		if pixelWindow.Pressed(BUTTON_LEFT) {
			state &= 0b1101
		}
		if pixelWindow.Pressed(BUTTON_UP) {
			state &= 0b1011
		}
		if pixelWindow.Pressed(BUTTON_DOWN) {
			state &= 0b0111
		}
	}
	if joypad_register&JOYPAD_action_buttons == 0 {
		// Read A/B/Select/Start
		if pixelWindow.Pressed(BUTTON_A) {
			state &= 0b1110
		}
		if pixelWindow.Pressed(BUTTON_B) {
			state &= 0b1101
		}
		if pixelWindow.Pressed(BUTTON_SELECT) {
			state &= 0b1011
		}
		if pixelWindow.Pressed(BUTTON_START) {
			state &= 0b0111
		}
	}

	// Perform an interrupt if any button went from unpressed to pressed
	var doInterrupt bool = false
	for i := 0; i < 4; i++ {
		if (gb.joypadState&(1<<i) != 0) && (state&(1<<i) == 0) {
			doInterrupt = true
		}
	}

	if doInterrupt {
		gb.SetInterruptRequestFlag(Interrupt_joypad)
	}

	// Update joypad register with button states
	gb.memory.set(JOYPAD, joypad_register|state)

	// Store state for next iteration
	// We could also just compare to the prior register state if we assume/enforce it is read only
	gb.joypadState = state
}
