package main

// Basic cartridge containing 32KiB ROM from 0000-7FFF
type ROMOnlyCartridge struct {
	rom []uint8
}

func MakeROMOnlyCartridge(data []uint8) *ROMOnlyCartridge {
	return &ROMOnlyCartridge{rom: data}
}

func (c *ROMOnlyCartridge) ReadFrom(address uint16) uint8 {
	return c.rom[address]
}

func (c *ROMOnlyCartridge) WriteTo(address uint16, value uint8) {
	// TODO: do we need to support RAM here too?
	panic("Cannot write to ROM Only Cartridge")
}
