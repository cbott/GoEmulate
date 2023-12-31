package cartridges

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	ROMEndAddress           = 0x8000
	ExternalRAMStartAddress = 0xA000
	ExternalRAMEndAddress   = 0xC000

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

// Cartridge interface represents a game cartridge containing ROM and (optionally) RAM banks
type Cartridge interface {
	ReadFrom(address uint16) uint8
	WriteTo(address uint16, value uint8)
	LoadRAM()
	SaveRAM()
	GetState() ([][RAMBankSize]uint8, uint8, bool, uint16)
	SetState([][RAMBankSize]uint8, uint8, bool, uint16)
}

// Common base for all cartridge types defining ROM and RAM banks
type CartridgeCore struct {
	filename string
	rom      []uint8
	ram      [][RAMBankSize]uint8

	// Number of available 16MiB ROM banks we can switch between
	numRomBanks uint16
	// Number of available 8MiB RAM banks we can switch between
	numRamBanks uint8

	// Currently selected ROM bank for 4000-7FFF
	romBank uint16
	// Currently selected RAM bank A000-BFFF
	ramBank uint8

	// Whether RAM reading and writing are enabled
	ramEnabled bool
}

// Return the current state of the cartridge
func (c CartridgeCore) GetState() ([][RAMBankSize]uint8, uint8, bool, uint16) {
	return c.ram, c.ramBank, c.ramEnabled, c.romBank
}

// Set the state of the cartridge
func (c *CartridgeCore) SetState(ram [][RAMBankSize]uint8, ramBank uint8, ramEnabled bool, romBank uint16) {
	c.ram = ram
	c.ramBank = ramBank
	c.ramEnabled = ramEnabled
	c.romBank = romBank
}

// Read a cartridge binary file and return the correct cartridge type containing the file contents
func Make(filename string) Cartridge {
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
	// Title length can vary by cartridge type so we will just stop at the first null character
	var title string = strings.Split(string(data[TitleAddress:TitleAddress+TitleLength]), "\x00")[0]

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
	case 0x01, 0x02, 0x03:
		return NewMBC1Cartridge(filename, data)
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		return NewMBC3Cartridge(filename, data)
	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		return NewMBC5Cartridge(filename, data)
	default:
		panic(fmt.Sprintf("Cartridge type %d not implemented", cartridgeType))
	}
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
