package gameboy

// Define register constants for each available CPU register

type register8 interface {
	register8Method()
}

type register8Name string

func (r register8Name) register8Method() {}

const (
	regA register8Name = "A"
	regF register8Name = "F"
	regB register8Name = "B"
	regC register8Name = "C"
	regD register8Name = "D"
	regE register8Name = "E"
	regH register8Name = "H"
	regL register8Name = "L"
)

type register16 interface {
	register16Method()
}

type register16Name string

func (r register16Name) register16Method() {}

const (
	regAF register16Name = "AF"
	regBC register16Name = "BC"
	regDE register16Name = "DE"
	regHL register16Name = "HL"
	regSP register16Name = "SP"
	regPC register16Name = "PC"
)
