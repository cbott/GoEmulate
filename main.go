package main

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func RenderScreen(window *pixelgl.Window, picture *pixel.PictureData, data *[ScreenWidth][ScreenHeight][3]uint8) {
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			col := data[x][y]
			rgb := color.RGBA{R: col[0], G: col[1], B: col[2], A: 0xFF}
			picture.Pix[(ScreenHeight-1-y)*ScreenWidth+x] = rgb
		}
	}

	// TODO: figure out why we're doing this, for now just clearing with red so we see if it is a problem
	// Seems like sprites should just go over it all
	// bg := color.RGBA{R: 0x08, G: 0x18, B: 0x20, A: 0xFF}
	bg := color.RGBA{R: 0xB0, G: 0x00, B: 0x00, A: 0xFF}
	window.Clear(bg)

	sprite := pixel.NewSprite(pixel.Picture(picture), pixel.R(0, 0, ScreenWidth, ScreenHeight))
	sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()))

	window.Update()
}

func CreateGameBoy() *Gameboy {
	var gb = Gameboy{}
	gb.cpu = &CpuRegisters{}
	gb.memory = &Memory{}
	gb.memory.Init()
	// TODO: Does interrupt master enable default to True?
	return &gb
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game Boy Emulator",
		Bounds: pixel.R(0, 0, ScreenWidth, ScreenHeight),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetSmooth(true)

	gb := CreateGameBoy()

	// TODO: Allow for skipping boot ROM if we want
	skipBoot := true

	if skipBoot {
		gb.memory.set(BOOT, 1)
		gb.memory.BypassBootROM()
		gb.cpu.BypassBootROM()
	}
	gb.memory.LoadROMFile("roms/tetris.gb")

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
			gb.keyboardToJoypad(win)
			gb.RunNextFrame()
			RenderScreen(win, picture, &gb.PreparedData)
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
