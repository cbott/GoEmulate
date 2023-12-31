package cartridges

// Basic cartridge containing 32KiB ROM from 0000-7FFF
type ROMOnlyCartridge struct {
	CartridgeCore
}

func NewROMOnlyCartridge(data []uint8) *ROMOnlyCartridge {
	c := ROMOnlyCartridge{}
	c.rom = data
	return &c
}

func (c *ROMOnlyCartridge) ReadFrom(address uint16) uint8 {
	return c.rom[address]
}

func (c *ROMOnlyCartridge) WriteTo(address uint16, value uint8) {
	// Writes to ROM cartridge are no-ops
}

// Save cartridge RAM contents to a file
func (c *ROMOnlyCartridge) SaveRAM() {
	// ROM only cartridge does not have RAM
}

// Load cartridge RAM from a file
func (c *ROMOnlyCartridge) LoadRAM() {
	// ROM only cartridge does not have RAM
}
