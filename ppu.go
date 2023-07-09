package main

// Pixel Processing Unit

/*
160 (W) x 144 (H) pixels

20 x 18 tiles of 8x8 pixels each

40 Sprites, 10 per line


Layers
------
Background
Window
Sprites


OAM RAM (Object Attribute Map = Sprite details)
     X Y # Flags
FE00 - - - -
FE04
...
FEC9

Flags: Priority, Flip X, Flip Y, Palette

PPU Registers
-------------
Adr    Name
FF40   [LCDC] LCD Control
   7   display enable
   6   window tile map address select (0=9800-9BFF, 1=9C00-9FFF)
   5   window enable
   4   BG and window tile data select (0=8800-97FF, 1=8000-8FFF)
   3   BG tile map address select (0=9800-9BFF, 1=9C00-9FFF)
   2   OBJ (sprite) size (0=8x8, 1=8x16)
   1   OBJ enable
   0   BG enable
FF41   [STAT] LCDC Status
   6   LYC=LY interrupt
   5   Mode 2 OAM interrupt
   4   Mode 1 V-Blank interrupt
   3   Mode 0 H-Blank interrupt
   2   LYC=LY flag
 1-0   Mode
FF42   [SCY] Scroll Y
FF43   [SCX] Scroll X
FF44   [LY] LCDC Y-Coordinate
FF45   [LYC] LY Compare
FF46   [DMA] DMA Transfer and Start
FF47   [BGP] Background Palette
 6-7   - defines which color 0b11 maps to (00=white, 01=light gray, 10=dark gray, 11=black)
 4-5   - color for 0b10
 2-3   - color for 0b01
 0-1   - color for 0b00
FF48   [OBP0] Object Palette 0
FF49   [OBP1] Object Palette 1
FF4A   [WY] Window Y Position
FF4B   [WX] Window X Position minus 7 (WX=7 corresponds to starting all the way on the left at column 0)
*/

const (
	ScreenWidth  = 160
	ScreenHeight = 144
)

// Vertical blanking lines below the screen
const VBlankLines = 10

const (
	SCY  = 0xFF42
	SCX  = 0xFF43
	LY   = 0xFF44
	LYC  = 0xFF45
	DMA  = 0xFF46
	BGP  = 0xFF47
	OBP0 = 0xFF48
	OBP1 = 0xFF49
	WY   = 0xFF4A
	WX   = 0xFF4B
)

const (
	LCDC                   = 0xFF40
	LCDC_display_enable    = 1 << 7
	LCDC_window_map_select = 1 << 6
	LCDC_window_enable     = 1 << 5
	LCDC_tile_data_select  = 1 << 4
	LCDC_bg_map_select     = 1 << 3
	LCDC_obj_size          = 1 << 2
	LCDC_obj_enable        = 1 << 1
	LCDC_bg_enable         = 1 << 0
)

const (
	STAT                     = 0xFF41
	STAT_lyc_eq_ly_flag      = 1 << 2
	STAT_hblank_interrupt    = 1 << 3
	STAT_vblank_interrupt    = 1 << 4
	STAT_oam_interrupt       = 1 << 5
	STAT_lyc_eq_ly_interrupt = 1 << 6
)

// PPU timing, in 4MHz cycles
const (
	OAMSearchCycles     = 20 * 4
	PixelTransferCycles = 43 * 4
	CyclesPerLine       = 114 * 4 // The rest of the cycles in a non-vblank line are h-blank
)

// Tile data is stored in VRAM starting at 0x8000 for Sprite tiles
// and either 0x8000 or 0x8800 for Background tiles (LCDC_tile_data_select)
const (
	TileDataAddressLow  = 0x8000
	TileDataAddressHigh = 0x8800
)

const (
	DisplayModeHBlank        = 0b00
	DisplayModeVBlank        = 0b01
	DisplayModeOAMSearch     = 0b10
	DisplayModePixelTransfer = 0b11
)

// Read display mode from the LCD Status register
func (gb *Gameboy) GetDisplayMode() uint8 {
	return gb.memory.get(STAT) & 0b11
}

// Set display mode in the LCD Status register
func (gb *Gameboy) SetDisplayMode(mode uint8) {
	gb.memory.set(STAT, (gb.memory.get(STAT)&0b11111100)|(mode&0b11))
}

func (gb *Gameboy) RunGraphicsProcess(cycles int) {
	status := gb.memory.get(STAT)
	currentLine := gb.memory.get(LY)
	mode := gb.GetDisplayMode()

	if (gb.memory.get(LCDC) & LCDC_display_enable) == 0 {
		// LCD is not enabled
		// TODO: verify nothing needs to run in this case
		// fmt.Println("LCD not enabled")
		return
	}

	var newMode uint8
	var interrupt bool = false

	if currentLine >= ScreenHeight {
		// Current line is in V-Blank section
		newMode = DisplayModeVBlank
		interrupt = (status & STAT_vblank_interrupt) != 0
	} else if gb.currentScanX < OAMSearchCycles {
		// Current line is a displayed row, and current scan is in OAM Search section
		newMode = DisplayModeOAMSearch
		interrupt = (status & STAT_oam_interrupt) != 0
	} else if gb.currentScanX < PixelTransferCycles {
		// Current line is a displayed row, and current scan is in Pixel Transfer section
		newMode = DisplayModePixelTransfer
		// There are no interrupts triggered on Pixel Transfer mode
	} else {
		// Current line is a displayed row, and current scan is in H-Blank section
		newMode = DisplayModeHBlank
		interrupt = (status & STAT_hblank_interrupt) != 0
	}

	if newMode != mode {
		// Our emulator will not process commands throughout a mode,
		// instead we just do all the work up front for simplicity
		if newMode == DisplayModePixelTransfer {
			// Write a line, or sections of the line
			gb.renderLine(currentLine)
		}
		if interrupt {
			gb.SetInterruptRequestFlag(Interrupt_lcd_stat)
		}
		gb.SetDisplayMode(newMode)
	}

	if currentLine == gb.memory.get(LYC) {
		// Set the LYC=LY flag
		status |= STAT_lyc_eq_ly_flag
		// Trigger an interrupt on this, if enabled
		if (status & STAT_lyc_eq_ly_interrupt) != 0 {
			// TODO: I think this is only supposed to trigger at the start of a line
			// but as-written we set it multiple times per line?
			gb.SetInterruptRequestFlag(Interrupt_lcd_stat)
		}
	} else {
		// Clear the LYC=LY flag
		status &= (0xFF ^ STAT_lyc_eq_ly_flag)
	}
	gb.memory.set(STAT, status)

	gb.currentScanX += cycles
	if gb.currentScanX >= CyclesPerLine {
		// If we get to the end of a line, move Y coordinate down to the next row and start back at the left
		currentLine++

		if currentLine > ScreenHeight+VBlankLines {
			// We are past the bottom of the screen, so we've now drawn the full frame
			gb.PreparedData = gb.screenData
			// TODO: is this really the best way to clear it?
			// and do we even need to clear it? We'll write over it later right?
			gb.screenData = [ScreenWidth][ScreenHeight][3]uint8{}

			currentLine = 0
		}

		// TODO: I think the idea here is that we don't set to 0 in case we didn't run this
		// function right at the end of a row. We don't want to accumulate timing error.
		gb.currentScanX -= CyclesPerLine

		gb.memory.set(LY, currentLine)

		if currentLine == ScreenHeight {
			// The CPU triggers an interrupt when it enters the vblank section
			gb.SetInterruptRequestFlag(Interrupt_vblank)
		}
	}
}

func (gb *Gameboy) renderLine(lineNumber uint8) {
	// Fill in a single line of the screen buffer
	control := gb.memory.get(LCDC)
	if (control & LCDC_bg_enable) != 0 {
		gb.renderLineTiles(lineNumber)
	}
	if (control & LCDC_obj_enable) != 0 {
		gb.renderLineSprites(lineNumber)
	}
}

func (gb *Gameboy) renderLineTiles(lineNumber uint8) [ScreenWidth]bool {
	// Render the background and window tiles in a single line
	// Returns an array with an element for each pixel indicating if Sprites can draw over it
	scrollX := gb.memory.get(SCX)
	scrollY := gb.memory.get(SCY)
	windowX := gb.memory.get(WX) - 7 // WX has a 7 pixel offset
	windowY := gb.memory.get(WY)
	control := gb.memory.get(LCDC)

	// This row contains some of the window if window drawing is enabled
	// and we are on or below the starting row of the window
	drawWindow := (control&LCDC_window_enable) != 0 && (lineNumber >= windowY)

	var tileDataStartAddress uint16
	if (control & LCDC_tile_data_select) == 0 {
		// The first half of the BG/Window tiles overlap with the last half of the Sprite tiles
		tileDataStartAddress = TileDataAddressHigh
	} else {
		// BG/Window tiles and Sprite tiles fully share the same address space
		tileDataStartAddress = TileDataAddressLow
	}

	var tileMapStartAddress uint16
	if drawWindow {
		// If window is drawn on this line, use window map
		// TODO: Does this work if window just covers part of the width?
		if (control & LCDC_window_map_select) == 0 {
			tileMapStartAddress = 0x9800
		} else {
			tileMapStartAddress = 0x9C00
		}
	} else {
		// If just background, use background map
		if (control & LCDC_bg_map_select) == 0 {
			tileMapStartAddress = 0x9800
		} else {
			tileMapStartAddress = 0x9C00
		}
	}

	var relativeY uint8
	if drawWindow {
		// If we're drawing the window, y position is referenced relative to first line in the window
		relativeY = lineNumber - windowY
	} else {
		// If drawing background, y position is relative to where we are scrolled to in the 32x32 background map
		relativeY = lineNumber + scrollY
	}

	// Determine which row of the 32x32 grid this tile is in
	tileRow := relativeY / 8

	palette := gb.memory.get(BGP)

	// Array with each value representing whether or not the corresponding pixel
	// is drawn with pallete entry 0 and will therefore be drawn over by sprites with priority 1
	lineBGPixelLowPriority := [ScreenWidth]bool{}
	var absoluteX uint8

	// Set pixel colors for this line
	for absoluteX = 0; absoluteX < ScreenWidth; absoluteX++ {
		var relativeX uint8
		if drawWindow && absoluteX >= windowX {
			relativeX = absoluteX - windowX
		} else {
			relativeX = absoluteX + scrollX
		}

		// Determine which column of the 32x32 grid this tile is in
		tileCol := relativeX / 8

		// Find the BG or Window map entry for this tile to see where in tile data to look
		var tileNumber uint8 = gb.memory.get(tileMapStartAddress + uint16(tileRow)*32 + uint16(tileCol))

		var tileAddress uint16
		if tileDataStartAddress == TileDataAddressLow {
			// If the data table is 0x8000-0x8FFF then tile number is 0-255 offset from 0x8000
			// each tile occupies 16 bytes, 2 bytes per line
			tileAddress = tileDataStartAddress + uint16(tileNumber)*16
		} else {
			// If the data table is 0x8800-0x97FF then tile number is -128-127 offset from 0x9000
			tileAddress = tileDataStartAddress + uint16((int16(int8(tileNumber))+128)*16)
		}

		// TODO: are we supposed to implement X flip or Y flip for BG tiles? I think that's sprites only

		// Each line in the tile is defined by 2 bytes, first byte holds the least significant bit of each pixel,
		// second byte hold the most significant bit, bit 7 being leftmost, bit 0 rightmost
		rowInTile := relativeY % 8
		lineLSB := gb.memory.get(tileAddress + uint16(rowInTile)*2)
		lineMSB := gb.memory.get(tileAddress + uint16(rowInTile)*2 + 1)

		columnInTile := relativeX % 8
		// pixelColor is the 2-bit value that we use to index into palette to get the displayed color
		var pixelColor uint8 = 0b00
		if lineLSB&(0b10000000>>columnInTile) != 0 {
			pixelColor |= 0b01
		}
		if lineMSB&(0b10000000>>columnInTile) != 0 {
			pixelColor |= 0b10
		}

		// Keep track of which pixels in this row used palette color 0, as these will be drawn
		// over by sprites with priority 1
		lineBGPixelLowPriority[absoluteX] = (pixelColor == 0b00)

		// Set the appropriate pixel of the screen buffer
		red, green, blue := getColorFromPalette(pixelColor, palette)
		gb.screenData[absoluteX][lineNumber][0] = red
		gb.screenData[absoluteX][lineNumber][1] = green
		gb.screenData[absoluteX][lineNumber][2] = blue
	}
	return lineBGPixelLowPriority
}

func getColorFromPalette(colorIndex uint8, palette uint8) (uint8, uint8, uint8) {
	// colorIndex: 2-bit index in the palette
	// palette: byte representing the 4 colors we have to choose from
	// Returns RGB values for the selected color
	color := 0b11 & (palette >> (colorIndex * 2))
	if color == 0b00 {
		// white
		return 255, 255, 255
	} else if color == 0b01 {
		// light gray
		return 170, 170, 170
	} else if color == 0b10 {
		// dark gray
		return 85, 85, 85
	} else {
		// black
		return 0, 0, 0
	}
}

func (gb *Gameboy) renderLineSprites(lineNumber uint8) {
}
