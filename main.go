package main

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const FramesPerSecond = 60.0
const CyclesPerFrame = CpuSpeed / int(FramesPerSecond)

type Gameboy struct {
	cpu    *CpuRegisters
	memory *Memory

	// TODO: re-evaluate this matrix
	// Matrix of pixel data which is used while the screen is rendering. When a
	// frame has been completed, this data is copied into the PreparedData matrix.
	screenData [ScreenWidth][ScreenHeight][3]uint8

	// PreparedData is a matrix of screen pixel data for a single frame which has
	// been fully rendered.
	PreparedData [ScreenWidth][ScreenHeight][3]uint8

	// scan is the horizontal "position" in clock cycles where the PPU is currently drawing
	// TODO: should this be stored in the PPU object if we make one?
	currentScanX int
}

// TODO: Move all *Gameboy functions to a separate file

// Return the value in memory pointed to by PC and then increment PC
func (gb *Gameboy) popPC() uint8 {
	pc := gb.cpu.get_register16("PC")
	gb.cpu.set_register16("PC", pc+1)
	return gb.memory.get(pc)
}

// Push a 16 bit value onto the stack and update the stack pointer
func (gb *Gameboy) pushToStack(high uint8, low uint8) {
	sp := gb.cpu.get_register16("SP")
	gb.memory.set(sp-1, high)
	gb.memory.set(sp-2, low)
	// Decrement stack pointer twice
	gb.cpu.set_register16("SP", sp-2)
}

// Pop a 16 bit value off of the stack and update the stack pointer
func (gb *Gameboy) popFromStack() uint16 {
	sp := gb.cpu.get_register16("SP")
	low := uint16(gb.memory.get(sp))
	high := uint16(gb.memory.get(sp + 1))
	// Increment stack pointer twice
	gb.cpu.set_register16("SP", sp+2)
	return (high << 8) | low
}

func (gb *Gameboy) RunNextFrame() {
	// Run Gameboy processes up to the next complete frame to be displayed
	for totalCycles := 0; totalCycles < CyclesPerFrame; {
		operationCycles := gb.RunNextOpcode()
		gb.RunGraphicsProcess(operationCycles)
		// TODO: handle hardware timers
		totalCycles += operationCycles
		// TODO: totalCycles += run interrupt service routines
	}
}

func (gb *Gameboy) RunNextOpcode() int {
	// Returns the number of clock cycles to complete (4MHz cycles)
	opcode := gb.popPC()
	return gb.Opcode(opcode) * 4
}

func RenderScreen(window *pixelgl.Window, picture *pixel.PictureData, data *[ScreenWidth][ScreenHeight][3]uint8) {
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			col := data[x][y]
			rgb := color.RGBA{R: col[0], G: col[1], B: col[2], A: 0xFF}
			picture.Pix[(ScreenHeight-1-y)*ScreenWidth+x] = rgb
		}
	}

	// TODO: figure out why we're doing this
	// Seems like sprites should just go over it all
	bg := color.RGBA{R: 0x08, G: 0x18, B: 0x20, A: 0xFF}
	window.Clear(bg)

	sprite := pixel.NewSprite(pixel.Picture(picture), pixel.R(0, 0, ScreenWidth, ScreenHeight))
	sprite.Draw(window, pixel.IM)

	window.Update()
}

func CreateGameBoy() *Gameboy {
	var gb = Gameboy{}
	gb.cpu = &CpuRegisters{}
	gb.memory = &Memory{}
	gb.memory.LoadBootROM()
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

	// TODO: I think cpu init shouldn't run until after bootrom
	// gb.cpu.Init()
	// gb.memory.Init()

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
			gb.RunNextFrame()
			RenderScreen(win, picture, &gb.PreparedData)
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
