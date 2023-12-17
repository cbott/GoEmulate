package main

import (
	"flag"
	"fmt"
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

	// Construct Pixel window
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

	// Construct Game Boy emulator
	gb := gameboy.NewGameBoy(!*runBootROM, *useDebugColors)
	gb.LoadCartridge(cartridges.Make(romFile))

	emulator := Emulator{
		console: gb,
		window:  win,
		speed:   1,
	}

	// Ticker will execute once per Game Boy frame
	var factor float64 = gameboy.FramesPerSecond
	ticker := time.NewTicker(time.Nanosecond * time.Duration(int64(1e9/factor)))

	for !win.Closed() {
		select {
		case <-ticker.C:
			update(&emulator)
		default:
		}
	}
}

func main() {
	opengl.Run(run)
}
