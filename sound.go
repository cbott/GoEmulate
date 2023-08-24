package main

// Control for the Audio Processing Unit (APU)

// Noise Register naming: NRxy
// x = channel number (1-4,  5 = global)
// y = register ID
const (
	// Channel 1: Pulse A (with period sweep)
	NR10 = 0xFF10 // Channel 1 sweep
	NR11 = 0xFF11 // Channel 1 length timer and duty cycle
	NR12 = 0xFF12 // Channel 1 volume and envelope
	NR13 = 0xFF13 // Channel 1 period low
	NR14 = 0xFF14 // Channel 1 period high and control
	// Channel 2: Pulse B
	NR21 = 0xFF16 // Channel 2 length timer and duty cycle
	NR22 = 0xFF17 // Channel 2 volume and envelope
	NR23 = 0xFF18 // Channel 2 period low
	NR24 = 0xFF19 // Channel 2 period high and control
	// Channel 3: Wave output
	NR30 = 0xFF1A // Channel 3 DAC enable
	NR31 = 0xFF1B // Channel 3 length timer
	NR32 = 0xFF1C // Channel 3 output level
	NR33 = 0xFF1D // Channel 3 period low
	NR34 = 0xFF1E // Channel 3 period high and control
	// Channel 3 Wave RAM is FF30-FF3F (16 bytes)
	WaveRAMStart = 0xFF30
	WaveRAMSize  = 16
	// Channel 4: Noise
	NR41 = 0xFF20 // Channel 4 length timer
	NR42 = 0xFF21 // Channel 4 volume and envelope
	NR43 = 0xFF22 // Channel 4 frequency and randomness
	NR44 = 0xFF23 // Channel 4 control
	// Global sound controls
	NR50 = 0xFF24 // Master volume and VIN panning
	NR51 = 0xFF25 // Sound panning
)

const (
	NR52            = 0xFF26 // Sound on/off
	NR52_apu_enable = 1 << 7 // 0 = all sound off
	NR52_ch4_on     = 1 << 3 // read only
	NR52_ch3_on     = 1 << 2 // read only
	NR52_ch2_on     = 1 << 1 // read only
	NR52_ch1_on     = 1 << 0 // read only
)

// Audio Processing Unit
type APU struct {
	// 0 = All sound off
	on bool
}
