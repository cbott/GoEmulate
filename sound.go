package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/oto"
)

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
	// Channel 4: Noise
	NR41 = 0xFF20 // Channel 4 length timer
	NR42 = 0xFF21 // Channel 4 volume and envelope
	NR43 = 0xFF22 // Channel 4 frequency and randomness
	NR44 = 0xFF23 // Channel 4 control
	// Global sound controls
	NR50 = 0xFF24 // Master volume and VIN panning
	// Channel 3 Wave RAM is FF30-FF3F (16 bytes)
	WaveRAMStart = 0xFF30
	WaveRAMSize  = 16
)

const (
	NR51               = 0xFF25 // Sound panning
	NR51_mix_ch1_right = 1 << 0
	NR51_mix_ch2_right = 1 << 1
	NR51_mix_ch3_right = 1 << 2
	NR51_mix_ch4_right = 1 << 3
	NR51_mix_ch1_left  = 1 << 5
	NR51_mix_ch2_left  = 1 << 5
	NR51_mix_ch3_left  = 1 << 6
	NR51_mix_ch4_left  = 1 << 7
)

// TODO: do we want to break out all the constants like this? Or just deal with it in the read/write section
const (
	NR52            = 0xFF26 // Sound on/off
	NR52_apu_enable = 1 << 7 // 0 = all sound off
	NR52_ch4_on     = 1 << 3 // read only
	NR52_ch3_on     = 1 << 2 // read only
	NR52_ch2_on     = 1 << 1 // read only
	NR52_ch1_on     = 1 << 0 // read only
)

// Audio player configuration constants
const (
	// Actual sample rate for sound played out of your speakers
	AudioSampleRate = 44100
	bufferSeconds   = 120
	CyclesPerSample = CpuSpeed / float64(AudioSampleRate)
)

// Audio Processing Unit
type APU struct {
	// Whether or not APU audio output is enabled - false = All sound off
	on bool

	// Number of CPU cycles processed since the last sample, resets back to 0 when a sample is taken
	cycleCounter float64

	channel1 *SoundChannel
	channel2 *SoundChannel
	channel3 *SoundChannel
	channel4 *SoundChannel

	player      *oto.Player
	audioBuffer chan [2]uint8

	// NR51 controls whether or not to mix each of the sound channels into the left and right audio outputs
	nr51RegisterValue uint8

	// Master volume is a value 0-7 where 0=very quiet, 7=full volume (no scaling)
	leftVolume  uint8
	rightVolume uint8
}

func (apu *APU) Init() {
	// TODO: change to being MakeAPU()?
	apu.on = true

	// Context Settings
	// 44100 Hz Sample rate: Standard audio frequency
	// 2 channels: This is stereo audio
	// 1 Byte bit depth: Game Boy audio channels have 8-bit output
	// 400 Byte buffer: stores 200 audio samples, comes out to ~1/4 frame length
	context, err := oto.NewContext(AudioSampleRate, 2, 1, 400)
	if err != nil {
		log.Fatalf("Audio initialization error: %v", err)
	}
	apu.player = context.NewPlayer()

	// At most we will allow up to 5000 samples to get backed up, hopefully this won't happen
	apu.audioBuffer = make(chan [2]uint8, 5000)

	// Initialize our sound channels
	apu.channel1 = &SoundChannel{channelNumber: 1}
	apu.channel2 = &SoundChannel{channelNumber: 2}
	apu.channel3 = &SoundChannel{channelNumber: 3}
	apu.channel4 = &SoundChannel{channelNumber: 4}

	// Start the go function which will continually pull samples from the audio buffer and play them
	frameTime := time.Second / time.Duration(bufferSeconds)
	ticker := time.NewTicker(frameTime)
	targetSamples := AudioSampleRate / bufferSeconds

	go func() {
		var reading [2]byte
		var buffer []byte
		for range ticker.C { // 1/120s
			fbLen := len(apu.audioBuffer)
			if fbLen >= targetSamples/2 { // 184
				newBuffer := make([]byte, fbLen*2)
				for i := 0; i < fbLen*2; i += 2 {
					reading = <-apu.audioBuffer
					newBuffer[i], newBuffer[i+1] = reading[0], reading[1]
				}
				buffer = newBuffer
			}

			_, err := apu.player.Write(buffer)
			if err != nil {
				log.Printf("error sampling: %v", err)
			}
		}
	}()
}

// Advance audio processing unit by the specified number of machine cycles (4MHz)
func (apu *APU) RunAudioProcess(cycles int) {
	if !apu.on {
		return
	}

	// Accumulate time until we have reached the length of 1 audio sample
	apu.cycleCounter += float64(cycles)
	if apu.cycleCounter < CyclesPerSample {
		return
	}
	apu.cycleCounter -= CyclesPerSample

	// Sample each channel
	sample1 := float64(apu.channel1.GetSample())
	sample2 := float64(apu.channel2.GetSample())
	sample3 := float64(apu.channel3.GetSample())
	sample4 := float64(apu.channel4.GetSample())

	// Mix into left and right audio outputs
	var leftUnscaled, rightUnscaled float64
	if (apu.nr51RegisterValue & NR51_mix_ch1_left) != 0 {
		leftUnscaled += sample1
	}
	if (apu.nr51RegisterValue & NR51_mix_ch1_right) != 0 {
		rightUnscaled += sample1
	}
	if (apu.nr51RegisterValue & NR51_mix_ch2_left) != 0 {
		leftUnscaled += sample2
	}
	if (apu.nr51RegisterValue & NR51_mix_ch2_right) != 0 {
		rightUnscaled += sample2
	}
	if (apu.nr51RegisterValue & NR51_mix_ch3_left) != 0 {
		leftUnscaled += sample3
	}
	if (apu.nr51RegisterValue & NR51_mix_ch3_right) != 0 {
		rightUnscaled += sample3
	}
	if (apu.nr51RegisterValue & NR51_mix_ch4_left) != 0 {
		leftUnscaled += sample4
	}
	if (apu.nr51RegisterValue & NR51_mix_ch4_right) != 0 {
		if (apu.nr51RegisterValue & NR51_mix_ch3_left) != 0 {
			leftUnscaled += sample3
		}
		if (apu.nr51RegisterValue & NR51_mix_ch3_right) != 0 {
			rightUnscaled += sample3
		}
		rightUnscaled += sample4
	}

	// Normalize
	leftUnscaled /= 4
	rightUnscaled /= 4

	// Apply the master volume scaling
	var left uint8 = uint8(leftUnscaled * float64(apu.leftVolume+1) / 8.0)
	var right uint8 = uint8(rightUnscaled * float64(apu.rightVolume+1) / 8.0)

	apu.audioBuffer <- [2]uint8{left, right}
}

// Write to an Audio control register
func (apu *APU) WriteTo(address uint16, value uint8) {
	switch address {
	case NR10:
		// Channel 1 sweep
		// TODO: Values should not take effect until triggered
		// Sweep Time: bits 4-6
		apu.channel1.sweepTime = (value & 0b01110000) >> 4
		// Sweep Direction: bit 3
		apu.channel1.sweepDirection = (value & 0b00001000) >> 3
		// Sweep Slope: bits 0-2
		apu.channel1.sweepSlope = value & 0b00000111
	case NR11:
		// Channel 1 length timer and duty cycle
		// Duty Cycle: bits 6-7
		apu.channel1.duty = (value & 0b11000000) >> 6
		// Initial sound length: bits 0-5
		// TODO: does this need to change the wave length timer immediately?
		apu.channel1.initialSoundLength = value & 0b00111111
	case NR12:
		// Channel 1 volume and envelope
		// TODO: these changes should only take effect when re-triggering the sound
		// Initial volume of envelope: bits 4-7
		apu.channel1.initialVolumeEnvelope = (value & 0b11110000) >> 4
		// Envelope Direction: bit 3
		apu.channel1.volumeEnvelopeDirection = (value & 0b00001000) >> 3
		// Volume Sweep Pace: bits 0-2
		apu.channel1.volumeSweepPace = value & 0b00000111
	case NR13:
		// Channel 1 period low
		// TODO: unsure when to actually update this
		// Low 8 bits of frequency value
		apu.channel1.frequencyValue = (apu.channel1.frequencyValue & 0x0700) | uint16(value)
	case NR14:
		// Channel 1 period high and control
		// High 3 bits of frequency value: bits 0-2
		apu.channel1.frequencyValue = (apu.channel1.frequencyValue & 0x00FF) | (uint16(value&0x7) << 8)
		// Sound Length Enable: bit 6
		// Takes effect immediately
		apu.channel1.soundLengthEnable = (value & (1 << 6)) != 0
		// Sound Trigger: bit 7
		if value&(1<<7) != 0 {
			apu.channel1.Trigger()
		}
	// TODO: reduce duplication across channels
	// likely need to make these methods of SoundChannel unless they access base APU values?
	case NR21:
		// Channel 2 length timer and duty cycle
		// Duty Cycle: bits 6-7
		apu.channel2.duty = (value & 0b11000000) >> 6
		// Initial sound length: bits 0-5
		// TODO: does this need to change the wave length timer immediately?
		apu.channel2.initialSoundLength = value & 0b00111111
	case NR22:
		// Channel 2 volume and envelope
		// TODO: these changes should only take effect when re-triggering the sound
		// Initial volume of envelope: bits 4-7
		apu.channel2.initialVolumeEnvelope = (value & 0b11110000) >> 4
		// Envelope Direction: bit 3
		apu.channel2.volumeEnvelopeDirection = (value & 0b00001000) >> 3
		// Volume Sweep Pace: bits 0-2
		apu.channel2.volumeSweepPace = value & 0b00000111
	case NR23:
		// Channel 2 period low
		// TODO: unsure when to actually update this
		// Low 8 bits of frequency value
		apu.channel2.frequencyValue = (apu.channel2.frequencyValue & 0x0700) | uint16(value)
	case NR24:
		// Channel 2 period high and control
		// High 3 bits of frequency value: bits 0-2
		apu.channel2.frequencyValue = (apu.channel2.frequencyValue & 0x00FF) | (uint16(value&0x7) << 8)
		// Sound Length Enable: bit 6
		// Takes effect immediately
		apu.channel2.soundLengthEnable = (value & (1 << 6)) != 0
		// Sound Trigger: bit 7
		if value&(1<<7) != 0 {
			apu.channel2.Trigger()
		}
	case NR30:
		apu.channel3.on = value&(1<<7) != 0
	case NR31:
		apu.channel3.initialSoundLength = value
	case NR32:
		// Bits 5-6
		apu.channel3.outputLevel = (value >> 5) & 0b11
	case NR33:
		// Channel 3 period low
		// TODO: unsure when to actually update this
		// Low 8 bits of frequency value
		apu.channel3.frequencyValue = (apu.channel3.frequencyValue & 0x0700) | uint16(value)
	case NR34:
		// Channel 3 period high and control
		// High 3 bits of frequency value: bits 0-2
		apu.channel3.frequencyValue = (apu.channel3.frequencyValue & 0x00FF) | (uint16(value&0x7) << 8)
		// Sound Length Enable: bit 6
		// Takes effect immediately
		apu.channel3.soundLengthEnable = (value & (1 << 6)) != 0
		// Sound Trigger: bit 7
		if value&(1<<7) != 0 {
			apu.channel3.Trigger()
		}
	case NR41:
		// Channel 4 length timer
		// Initial sound length: bits 0-5
		apu.channel2.initialSoundLength = value & 0b00111111
	case NR42:
		// Channel 4 volume and envelope
		// Initial volume of envelope: bits 4-7
		apu.channel4.initialVolumeEnvelope = (value & 0b11110000) >> 4
		// Envelope Direction: bit 3
		apu.channel4.volumeEnvelopeDirection = (value & 0b00001000) >> 3
		// Volume Sweep Pace: bits 0-2
		apu.channel4.volumeSweepPace = value & 0b00000111
	case NR43:
		// Channel 4 frequency and randomness
		// Clock Shift: bits 4-7
		apu.channel4.shiftRegisterClockShift = (value >> 4) & 0xF
		// LFSR Width: bit 3
		apu.channel4.shiftRegisterWidth = (value >> 3) & 1
		// Clock Divider: bits 0-2
		apu.channel4.shiftRegisterClockRatio = value & 0b111
	case NR44:
		// Channel 4 control
		// Sound Length Enable: bit 6
		// Takes effect immediately
		apu.channel4.soundLengthEnable = (value & (1 << 6)) != 0
		// Sound Trigger: bit 7
		if value&(1<<7) != 0 {
			apu.channel4.Trigger()
		}
	case NR50:
		// Master volume and VIN panning
		// Note: ignoring bits 3 and 7 which control VIN (audio provided by cartridge)
		// Right output volume: bits 0-2
		apu.rightVolume = value & 0x7
		// Left output volume: bits 4-6
		apu.leftVolume = (value >> 4) & 0x7
	case NR51:
		// Sound panning
		apu.nr51RegisterValue = value
	case NR52:
		// Sound on/off
		apu.on = (value & NR52_apu_enable) != 0
		apu.channel1.on = (value & NR52_ch1_on) != 0
		apu.channel2.on = (value & NR52_ch2_on) != 0
		apu.channel3.on = (value & NR52_ch3_on) != 0
		// apu.channel4.on = (value & NR52_ch4_on) != 0
	default:
		// Writing to invalid address
		return
	}
}