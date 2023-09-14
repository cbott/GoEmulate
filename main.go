package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/cbott/GoEmulate/cartridges"
	"github.com/cbott/GoEmulate/gameboy"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Emulator settings
// DefaultScale: initial size of the Game Boy window
// SkipBoot: skip running the bootrom and begin execution from cartridge immediately
// UseDebugColors: color sprites red, window green, background blue
const (
	// TODO: change to command line flags or at least allow overriding that way?
	DefaultScale   = 4
	SkipBoot       = true
	UseDebugColors = false
)

func RenderScreen(window *pixelgl.Window, picture *pixel.PictureData, data *[gameboy.ScreenWidth][gameboy.ScreenHeight][3]uint8) {
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
	cfg := pixelgl.WindowConfig{
		Title:     "Game Boy Emulator",
		Bounds:    pixel.R(0, 0, gameboy.ScreenWidth*DefaultScale, gameboy.ScreenHeight*DefaultScale),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	gb := gameboy.NewGameBoy(SkipBoot, UseDebugColors)
	gb.LoadCartridge(cartridges.Make("roms/pokemon_yellow.gb"))

	picture := &pixel.PictureData{
		Pix:    make([]color.RGBA, gameboy.ScreenWidth*gameboy.ScreenHeight),
		Stride: gameboy.ScreenWidth,
		Rect:   pixel.R(0, 0, gameboy.ScreenWidth, gameboy.ScreenHeight),
	}

	// Ticker will execute once per Gameboy frame
	ticker := time.NewTicker(time.Second / gameboy.FramesPerSecond)

	var frameSpeedUp int = 1
	last := time.Now()
	avgwindow := 60
	avgcounter := 0

	for !win.Closed() {
		select {
		case <-ticker.C:
			gb.ReadKeyboard(win)

			for i := 0; i < frameSpeedUp; i++ {
				gb.RunNextFrame()
			}

			RenderScreen(win, picture, &gb.ScreenData)
			// Pressing 'P' will save RAM contents to a file
			if win.JustPressed(pixelgl.KeyP) {
				gb.SaveCartridgeRAM()
			}
			// Pressing "+/=" increases speed-up
			if win.JustPressed(pixelgl.KeyEqual) && frameSpeedUp < 10 {
				// Limit to 10x speed, on my machine we start to drop framerate around 13x
				frameSpeedUp++
				fmt.Printf("Increased speed to %v\n", frameSpeedUp)
			}
			// Pressing "-/_" decreases speed-up
			if win.JustPressed(pixelgl.KeyMinus) && frameSpeedUp > 1 {
				frameSpeedUp--
				fmt.Printf("Decreased speed to %v\n", frameSpeedUp)
			}

			avgcounter++
			if avgcounter >= avgwindow {
				avgcounter = 0
				dt := time.Since(last).Seconds()
				last = time.Now()
				win.SetTitle(fmt.Sprintf("FPS: %.f", float64(avgwindow)/dt))

			}
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
