// User Interface
// this file handles graphics and keyboard input for the GoEmulate emulator
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/cbott/GoEmulate/gameboy"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

// Button-to-keyboard mapping
const (
	// Game Boy joypad
	KEY_DOWN   = pixel.KeyDown
	KEY_UP     = pixel.KeyUp
	KEY_LEFT   = pixel.KeyLeft
	KEY_RIGHT  = pixel.KeyRight
	KEY_START  = pixel.KeyEnter
	KEY_SELECT = pixel.KeyS
	KEY_B      = pixel.KeyZ
	KEY_A      = pixel.KeyX
	// Emulator controls
	KEY_WRITE_RAM  = pixel.KeyP
	KEY_SPEED_UP   = pixel.KeyEqual
	KEY_SPEED_DOWN = pixel.KeyMinus
	KEY_SAVESTATE1 = pixel.Key1 // SHIFT + SAVESTATE restores a previously saved state
	KEY_SAVESTATE2 = pixel.Key2
	KEY_SAVESTATE3 = pixel.Key3
)

var saveStateKeys = [...]pixel.Button{KEY_SAVESTATE1, KEY_SAVESTATE2, KEY_SAVESTATE3}

// Emulator defines a convenient structure for passing logic core, display, and state information around together
type Emulator struct {
	console *gameboy.Gameboy
	window  *opengl.Window
	speed   int
}

// update runs 1 or more frames worth of CPU cycles on the emulator core (depending on specified speed),
// processes inputs from the keyboard, and updates the display to match the new state of the emulator
func update(emulator *Emulator) {
	// Run the console for 1 frame, or multiple frames for "fast-forwarding"/speed-up
	for i := 0; i < emulator.speed; i++ {
		emulator.console.RunNextFrame()
	}
	render(emulator.window, &emulator.console.ScreenData)

	joypadstate := gameboy.ButtonState{
		BtnA:      emulator.window.Pressed(KEY_A),
		BtnB:      emulator.window.Pressed(KEY_B),
		BtnSelect: emulator.window.Pressed(KEY_SELECT),
		BtnStart:  emulator.window.Pressed(KEY_START),
		BtnRight:  emulator.window.Pressed(KEY_RIGHT),
		BtnLeft:   emulator.window.Pressed(KEY_LEFT),
		BtnUp:     emulator.window.Pressed(KEY_UP),
		BtnDown:   emulator.window.Pressed(KEY_DOWN),
	}
	emulator.console.SetButtonStates(&joypadstate)

	// Save to cartridge
	if emulator.window.JustPressed(KEY_WRITE_RAM) {
		emulator.console.SaveCartridgeRAM()
	}

	// Emulation speed
	if emulator.window.JustPressed(KEY_SPEED_UP) && emulator.speed < MaximumSpeed {
		// No specific reason to limit speed but it seemed reasonable for usability
		// On my machine we start to drop framerate around 13x
		emulator.speed++
		fmt.Printf("Increased speed to %v\n", emulator.speed)
	}
	if emulator.window.JustPressed(KEY_SPEED_DOWN) && emulator.speed > 1 {
		emulator.speed--
		fmt.Printf("Decreased speed to %v\n", emulator.speed)
	}

	// CPU States
	for i := 0; i < 3; i++ {
		if emulator.window.JustPressed(saveStateKeys[i]) {
			var success bool
			var action string

			if emulator.window.Pressed(pixel.KeyLeftShift) || emulator.window.Pressed(pixel.KeyRightShift) {
				// Recall
				success = emulator.console.RecallState(i)
				action = "recall"
			} else {
				// Store
				success = emulator.console.StoreState(i)
				action = "store"
			}

			if success {
				fmt.Printf("action %v:%d succeeded\n", action, i+1)
			} else {
				fmt.Printf("action %v:%d failed\n", action, i+1)
			}
		}
	}
}

// render displays a 2D array of RGB triplets, data, to the window with appropriate scaling
func render(window *opengl.Window, data *[gameboy.ScreenWidth][gameboy.ScreenHeight][3]uint8) {
	// Convert RGB array to PictureData that can be consumed by pixel
	picture := pixel.PictureData{
		Pix:    make([]color.RGBA, gameboy.ScreenWidth*gameboy.ScreenHeight),
		Stride: gameboy.ScreenWidth,
		Rect:   pixel.R(0, 0, gameboy.ScreenWidth, gameboy.ScreenHeight),
	}

	for x := 0; x < gameboy.ScreenWidth; x++ {
		column := data[x]
		for y := 0; y < gameboy.ScreenHeight; y++ {
			rgb := color.RGBA{R: column[y][0], G: column[y][1], B: column[y][2], A: 0xFF}
			picture.Pix[(gameboy.ScreenHeight-1-y)*gameboy.ScreenWidth+x] = rgb
		}
	}

	// Clear the screen, also sets color for areas of window not filled by Game Boy screen
	bg := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	window.Clear(bg)

	// Scale the Game Boy screen to maximize the size within the window
	// scale = min(windowX/gameboyX, windowY/gameboyY)
	windowSize := window.Bounds().Size()
	divisor := pixel.V(1.0/gameboy.ScreenWidth, 1.0/gameboy.ScreenHeight)
	scale := math.Min(windowSize.ScaledXY(divisor).XY())

	// Draw the Game Boy screen to the window
	sprite := pixel.NewSprite(&picture, pixel.R(0, 0, gameboy.ScreenWidth, gameboy.ScreenHeight))
	sprite.Draw(window, pixel.IM.Scaled(pixel.ZV, scale).Moved(window.Bounds().Center()))

	window.Update()
}
