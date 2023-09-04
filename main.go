package main

import (
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Emulator settings
// DefaultScale: initial size of the Game Boy window
// SkipBoot: skip running the bootrom and begin execution from cartridge immediately
// UseDebugColors: color sprites red, window green, background blue
const (
	DefaultScale   = 4
	SkipBoot       = false
	UseDebugColors = false
)

func RenderScreen(window *pixelgl.Window, picture *pixel.PictureData, data *[ScreenWidth][ScreenHeight][3]uint8) {
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			col := data[x][y]
			rgb := color.RGBA{R: col[0], G: col[1], B: col[2], A: 0xFF}
			picture.Pix[(ScreenHeight-1-y)*ScreenWidth+x] = rgb
		}
	}

	// Clear the screen, also sets color for areas of window not filled by Game Boy screen
	bg := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
	window.Clear(bg)

	// Scale the Game Boy screen to maximize the size within the window
	// scale = min(windowX/gameboyX, windowY/gameboyY)
	windowSize := window.Bounds().Size()
	divisor := pixel.V(1.0/ScreenWidth, 1.0/ScreenHeight)
	scale := math.Min(windowSize.ScaledXY(divisor).XY())

	// Draw the Game Boy screen to the window
	sprite := pixel.NewSprite(pixel.Picture(picture), pixel.R(0, 0, ScreenWidth, ScreenHeight))
	sprite.Draw(window, pixel.IM.Scaled(pixel.ZV, scale).Moved(window.Bounds().Center()))

	window.Update()
}

func NewGameBoy() *Gameboy {
	var gb = Gameboy{}
	gb.cpu = &CpuRegisters{}
	gb.memory = &Memory{}
	gb.memory.Init()

	if SkipBoot {
		gb.memory.BypassBootROM()
		gb.cpu.BypassBootROM()
	}

	return &gb
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Game Boy Emulator",
		Bounds:    pixel.R(0, 0, ScreenWidth*DefaultScale, ScreenHeight*DefaultScale),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	gb := NewGameBoy()
	gb.LoadCartridge(parseCartridgeFile("roms/pokemon_red.gb"))

	picture := &pixel.PictureData{
		Pix:    make([]color.RGBA, ScreenWidth*ScreenHeight),
		Stride: ScreenWidth,
		Rect:   pixel.R(0, 0, ScreenWidth, ScreenHeight),
	}

	// Ticker will execute once per Gameboy frame
	ticker := time.NewTicker(time.Second / FramesPerSecond)

	for !win.Closed() {
		select {
		case <-ticker.C:
			gb.ReadKeyboard(win)
			gb.RunNextFrame()
			RenderScreen(win, picture, &gb.screenData)
			// Pressing 'P' will save RAM contents to a file
			if win.JustPressed(pixelgl.KeyP) {
				gb.memory.cartridge.SaveRAM()
			}
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
