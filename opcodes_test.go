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

func TestDAA(t *testing.T) {
	var testcases = []teststructure{
		{0x99, false, false, false, 0x99, false, false, false},
		{0x00, false, false, false, 0x00, false, false, false},
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

		if finalValue != testcase.finalValue || finalH != testcase.finalH || finalC != testcase.finalC {
			t.Fatalf("DAA operation on {A=0x%x, N=%t, H=%t, C=%t} resulted in {A=0x%x, Z=%t, H=%t, C=%t}, expected {A=0x%x, Z=%t, H=%t, C=%t}",
				testcase.initialValue, testcase.initialN, testcase.initialH, testcase.initialC,
				finalValue, finalZ, finalH, finalC,
				testcase.finalValue, testcase.finalZ, testcase.finalH, testcase.finalC,
			)
		}
	}
}
