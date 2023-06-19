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

func (gb *Gameboy) popPC() uint8 {
	pc := gb.cpu.get_register16("PC")
	gb.cpu.set_register16("PC", pc+1)
	return gb.memory.get(pc)
}

func (gb *Gameboy) RunNextFrame() {
	// Run Gameboy processes up to the next complete frame to be displayed
	for totalCycles := 0; totalCycles < CyclesPerFrame; {
		operationCycles := gb.RunNextOpcode()
		gb.RunGraphicsProcess(operationCycles)
		totalCycles += operationCycles
	}
}

func (gb *Gameboy) RunNextOpcode() int {
	// Returns the number of clock cycles to complete (4MHz cycles)
	opcode := gb.popPC()
	return gb.Opcode(opcode) * 4
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game Boy Emulator",
		Bounds: pixel.R(0, 0, 1024, 768),
		// VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	// win.SetSmooth(true)

	var gb = Gameboy{}

	// TODO: I think cpu init shouldn't run until after bootrom
	gb.cpu.Init()
	gb.memory.Init()

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
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
