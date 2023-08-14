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

// Update the Joypad button states with what is currently pressed on the keyboard
// and requests an interrupt if a button just became pressed
func (gb *Gameboy) ReadKeyboard(pixelWindow *pixelgl.Window) {
	// 1=unpressed, 0=pressed
	// invert at the end so bit operations are easier
	var state uint8 = 0
	// Read direction pad
	if pixelWindow.Pressed(BUTTON_DOWN) {
		state |= 1 << 7
	}
	if pixelWindow.Pressed(BUTTON_UP) {
		state |= 1 << 6
	}
	if pixelWindow.Pressed(BUTTON_LEFT) {
		state |= 1 << 5
	}
	if pixelWindow.Pressed(BUTTON_RIGHT) {
		state |= 1 << 4
	}
	if pixelWindow.Pressed(BUTTON_START) {
		state |= 1 << 3
	}
	if pixelWindow.Pressed(BUTTON_SELECT) {
		state |= 1 << 2
	}
	if pixelWindow.Pressed(BUTTON_B) {
		state |= 1 << 1
	}
	if pixelWindow.Pressed(BUTTON_A) {
		state |= 1 << 0
	}

	state = ^state

	// Perform an interrupt if any button went from unpressed to pressed
	var doInterrupt bool = false
	for i := 0; i < 8; i++ {
		if (gb.memory.ButtonStates&(1<<i) != 0) && (state&(1<<i) == 0) {
			doInterrupt = true
		}
	}

	gb.memory.ButtonStates = state

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
