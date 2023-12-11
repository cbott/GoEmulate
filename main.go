package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/cbott/GoEmulate/cartridges"
	"github.com/cbott/GoEmulate/gameboy"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

// initial size of the Game Boy window, overridden by command line flag
const DefaultScale = 4

// Maximum allowable speed multiplier for emulation
const MaximumSpeed = 10

func RenderScreen(window *opengl.Window, picture *pixel.PictureData, data *[gameboy.ScreenWidth][gameboy.ScreenHeight][3]uint8) {
	for y := 0; y < gameboy.ScreenHeight; y++ {
		for x := 0; x < gameboy.ScreenWidth; x++ {
			col := data[x][y]
			rgb := color.RGBA{R: col[0], G: col[1], B: col[2], A: 0xFF}
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
	sprite := pixel.NewSprite(pixel.Picture(picture), pixel.R(0, 0, gameboy.ScreenWidth, gameboy.ScreenHeight))
	sprite.Draw(window, pixel.IM.Scaled(pixel.ZV, scale).Moved(window.Bounds().Center()))

	window.Update()
}

func run() {
	// Parse cmd line args
	runBootROM := flag.Bool("bootrom", false, "run boot ROM prior to cartridge")
	useDebugColors := flag.Bool("debug", false, "use debug colors (color sprites red, window green, background blue)")
	scaleflag := flag.Int("scale", DefaultScale, "window scale factor")
	flag.Parse()

	romFile := flag.Arg(0)
	if romFile == "" {
		fmt.Println("ROM file must be specified")
		os.Exit(1)
	}

	var scale float64 = float64(*scaleflag)
	cfg := opengl.WindowConfig{
		Title:     "Game Boy Emulator",
		Bounds:    pixel.R(0, 0, gameboy.ScreenWidth*scale, gameboy.ScreenHeight*scale),
		VSync:     true,
		Resizable: true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	gb := gameboy.NewGameBoy(!*runBootROM, *useDebugColors)
	gb.LoadCartridge(cartridges.Make(romFile))

	var savestate *gameboy.SaveState

	picture := &pixel.PictureData{
		Pix:    make([]color.RGBA, gameboy.ScreenWidth*gameboy.ScreenHeight),
		Stride: gameboy.ScreenWidth,
		Rect:   pixel.R(0, 0, gameboy.ScreenWidth, gameboy.ScreenHeight),
	}

	// Ticker will execute once per Gameboy frame
	var factor float64 = gameboy.FramesPerSecond
	ticker := time.NewTicker(time.Nanosecond * time.Duration(int64(1e9/factor)))

	var frameSpeedUp int = 1

	for !win.Closed() {
		select {
		case <-ticker.C:
			gb.ReadKeyboard(win)
			for i := 0; i < frameSpeedUp; i++ {
				gb.RunNextFrame()
			}
			RenderScreen(win, picture, &gb.ScreenData)

			// Pressing 'P' will save RAM contents to a file
			if win.JustPressed(pixel.KeyP) {
				gb.SaveCartridgeRAM()
			}
			// Pressing '+/=' increases speed-up
			if win.JustPressed(pixel.KeyEqual) && frameSpeedUp < MaximumSpeed {
				// Limit to 10x speed, on my machine we start to drop framerate around 13x
				frameSpeedUp++
				fmt.Printf("Increased speed to %v\n", frameSpeedUp)
			}
			// Pressing '-/_' decreases speed-up
			if win.JustPressed(pixel.KeyMinus) && frameSpeedUp > 1 {
				frameSpeedUp--
				fmt.Printf("Decreased speed to %v\n", frameSpeedUp)
			}
			// Pressing 1 saves the CPU state
			// Pressing Shift+1 loads the CPU state
			if win.JustPressed(pixel.Key1) {
				if win.Pressed(pixel.KeyLeftShift) || win.Pressed(pixel.KeyRightShift) {
					err := gameboy.RestoreState(gb, savestate)
					if err != nil {
						fmt.Println("Unable to restore state")
					} else {
						fmt.Println("Restored State 1")
					}
				} else {
					savestate = gameboy.NewSaveState(gb)
					fmt.Println("Saved State 1")
				}
			}

		default:
		}
	}
}

func main() {
	opengl.Run(run)
}
