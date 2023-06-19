package main

// define Cartridge interface
type Cartridge interface {
	ReadFrom(address uint16) uint8
}

type ROMOnlyCartridge struct {
	// Basic cartridge containing 32kB ROM from 0000-7FFF
	rom []uint8
}

func (c *ROMOnlyCartridge) ReadFrom(address uint16) uint8 {
	return c.rom[address]
}

type MemoryBankController1Cartridge struct {
	// Memory Bank Controller 1 Cartridge
	// 0000-3FFF Non switchable ROM bank
	// 4000-7FFF Switchable ROM banks
	rom []uint8
}
