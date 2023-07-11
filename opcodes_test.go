package main

import (
	"testing"
)

type teststructure struct {
	initialValue uint8
	initialN     bool
	initialH     bool
	initialC     bool
	finalValue   uint8
	finalZ       bool
	finalH       bool
	finalC       bool
}

// Test add with carry instruction
func TestADC(t *testing.T) {
	// We will just test with register B since they all use the same function internally
	type testADCStructure struct {
		initialA uint8
		initialB uint8
		initialC bool
		finalA   uint8
		finalZ   bool
		finalN   bool
		finalH   bool
		finalC   bool
	}
	testcases := []testADCStructure{
		{0, 0, false, 0, true, false, false, false},
		{1, 127, false, 128, false, false, true, false},
		{1, 128, false, 129, false, false, false, false},
		{0, 0, true, 1, false, false, false, false},
		{0, 255, true, 0, true, false, true, true},
	}

	gb := CreateGameBoy()

	for i, testcase := range testcases {
		gb.cpu.set_register("A", testcase.initialA)
		gb.cpu.set_register("B", testcase.initialB)
		gb.cpu.set_flag(FlagC, testcase.initialC)

		gb.Opcode(0x88) // ADC B

		if gb.cpu.get_register("A") != testcase.finalA ||
			gb.cpu.get_flag(FlagZ) != testcase.finalZ ||
			gb.cpu.get_flag(FlagN) != testcase.finalN ||
			gb.cpu.get_flag(FlagH) != testcase.finalH ||
			gb.cpu.get_flag(FlagC) != testcase.finalC {
			t.Fatalf("Failed ADC test %v\nExpected {A=0x%x, Z=%t, N=%t, H=%t, C=%t}, Got {A=0x%x, Z=%t, N=%t, H=%t, C=%t}",
				i, testcase.finalA, testcase.finalZ, testcase.finalN, testcase.finalH, testcase.finalC,
				gb.cpu.get_register("A"), gb.cpu.get_flag(FlagZ), gb.cpu.get_flag(FlagN), gb.cpu.get_flag(FlagH), gb.cpu.get_flag(FlagC))
		}
	}

}

func TestDAA(t *testing.T) {
	var testcases = []teststructure{
		{0x99, false, false, false, 0x99, false, false, false},
		{0xFA, false, false, false, 0x60, false, false, true},
		{0x00, false, false, false, 0x00, true, false, false},
		{0x00, false, true, true, 0x66, false, false, true},
		{0x9A, false, false, false, 0x00, true, false, true},
		{0x9A, true, false, false, 0x9A, false, false, false},
		{0x33, false, true, false, 0x39, false, false, false},
		{0xB4, false, false, false, 0x14, false, false, true},
		{0xB4, true, true, true, 0x4e, false, false, true},
	}
	gb := CreateGameBoy()

	for _, testcase := range testcases {
		gb.cpu.set_register("A", testcase.initialValue)
		gb.cpu.set_flag(FlagN, testcase.initialN)
		gb.cpu.set_flag(FlagH, testcase.initialH)
		gb.cpu.set_flag(FlagC, testcase.initialC)

		gb.Opcode(0x27)

		finalValue := gb.cpu.get_register("A")
		finalZ := gb.cpu.get_flag(FlagZ)
		finalH := gb.cpu.get_flag(FlagH)
		finalC := gb.cpu.get_flag(FlagC)

		if finalValue != testcase.finalValue || finalZ != testcase.finalZ || finalH != testcase.finalH || finalC != testcase.finalC {
			t.Fatalf("DAA operation on {A=0x%x, N=%t, H=%t, C=%t} resulted in {A=0x%x, Z=%t, H=%t, C=%t}, expected {A=0x%x, Z=%t, H=%t, C=%t}",
				testcase.initialValue, testcase.initialN, testcase.initialH, testcase.initialC,
				finalValue, finalZ, finalH, finalC,
				testcase.finalValue, testcase.finalZ, testcase.finalH, testcase.finalC,
			)
		}
	}
}
