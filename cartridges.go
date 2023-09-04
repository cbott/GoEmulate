package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	CartridgeTypeAddress = 0x0147
	ROMSizeAddress       = 0x0148
	RAMSizeAddress       = 0x0149
	TitleAddress         = 0x0134
	TitleLength          = 16
	ROMBankSize          = 0x4000 // 16 KiB
	RAMBankSize          = 0x2000 // 8 KiB
)

//    Available ROM Sizes
// Key		Size		Banks
//  0		32 KiB		2 (no banking)
//  1		64 KiB		4
//  2		128 KiB		8
//  3		256 KiB		16
//  4		512 KiB		32
//  5		1 MiB		64
//  6		2 MiB		128
//  7		4 MiB		256
//  8		8 MiB		512

var cartridgeTypeMap = map[uint8]string{
	0x00: "ROM Only",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x08: "ROM+RAM",
	0x09: "ROM+RAM+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY",
	0x11: "MBC3",
	0x12: "MBC3+RAM",
	0x13: "MBC3+RAM+BATTERY",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
	0xFC: "POCKET CAMERA",
	0xFD: "BANDAI TAMA5",
	0xFE: "HuC3",
	0xFF: "HuC1+RAM+BATTERY",
}

// RAM Size in KiB
var ramSizeMap = map[uint8]uint16{
	0: 0,
	2: 8,
	3: 32,
	4: 128,
}

// define Cartridge interface
type Cartridge interface {
	ReadFrom(address uint16) uint8
	WriteTo(address uint16, value uint8)
	LoadRAM()
	SaveRAM()
}

// Read a cartridge binary file and return the correct cartridge type containing the file contents
func parseCartridgeFile(filename string) Cartridge {
	// Load cartridge binary data
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Unable to load ROM file %s", filename))
	}

	// Parse out cartridge header attributes
	cartridgeType := data[CartridgeTypeAddress]
	cartridgeTypeString, ok := cartridgeTypeMap[cartridgeType]
	if !ok {
		panic(fmt.Sprintf("Unknown cartridge type %d", cartridgeType))
	}

	ramSizeKey := data[RAMSizeAddress]
	ramSize, ok := ramSizeMap[ramSizeKey]
	if !ok {
		panic(fmt.Sprintf("Unknown RAM Size code %d", ramSizeKey))
	}

	var romSize int = 32 * (1 << data[ROMSizeAddress])
	var title string = string(data[TitleAddress : TitleAddress+TitleLength])

	fmt.Printf("Cartridge file: %s\n", filename)
	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Type: %s\n", cartridgeTypeString)
	fmt.Printf("ROM Size: %d KiB\n", romSize)
	fmt.Printf("RAM Size: %d KiB\n", ramSize)

	// Validate ROM Size listed in the cartridge header
	if romSize*1024 != len(data) {
		panic(fmt.Sprintf("ROM size in cartridge header does not match file size\nHeader:\t%d B\nFile:\t%d B",
			romSize*1024, len(data)))
	}

	// Return correct cartridge type for this file
	switch cartridgeType {
	case 0x00:
		return NewROMOnlyCartridge(data)
	case 0x01:
		return NewMBC1Cartridge(data)
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		return NewMBC3Cartridge(filename, data)
	default:
		panic(fmt.Sprintf("Cartridge type %d not implemented", cartridgeType))
	}
}

// Load an initialized Cartridge struct into Game Boy memory
func (gb *Gameboy) LoadCartridge(c Cartridge) {
	gb.memory.cartridge = c
}

// Generate a name for a cartridge RAM save file based on the original ROM file name (filename.ram)
func getSaveFileName(name string) string {
	return name + ".ram"
}

// Write the contents of all RAM banks to a RAM save file (filename.ram)
func WriteRAMToFile(filename string, ramBanks [][RAMBankSize]uint8) error {
	filename = getSaveFileName(filename)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write each RAM bank to the file
	for i := 0; i < len(ramBanks); i++ {
		_, err = f.Write(ramBanks[i][:])
		if err != nil {
			return err
		}
	}

	log.Printf("Saved RAM to file %v\n", filename)
	return nil
}

// Read from a RAM save file to fill RAM banks
func ReadRAMFromFile(filename string, ramBanks [][RAMBankSize]uint8) error {
	filename = getSaveFileName(filename)
	// Load RAM binary file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Check RAM size
	banks := len(ramBanks)
	expectedBytes := banks * RAMBankSize
	if len(data) != expectedBytes {
		return fmt.Errorf("RAM file %s size (%vB) does not match cartrige expectation (%vB)", filename, len(data), expectedBytes)
	}

	for i := 0; i < expectedBytes; i++ {
		ramBanks[i/RAMBankSize][i%RAMBankSize] = data[i]
	}

	log.Printf("Loaded RAM from file %v\n", filename)
	return nil
}
