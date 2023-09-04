package sound

import (
	"math"
)

const (
	// Period of one tick for a sound channel length timer
	LengthTickTime = 1 / 256.0
	// Period of one tick for a sound channel frequency sweep
	SweepTickTime = 1 / 128.0
	// Period of one tick for a sound channel volume envelope
	VolumeTickTime = 1 / 64.0
)

// WaveDuty defines the fraction of the wave which is high vs low
var WaveDuty = map[uint8]float64{
	0b00: 7 / 8.0, // _-------_-------
	0b01: 6 / 8.0, // __------__------
	0b10: 4 / 8.0, // ____----____----
	0b11: 2 / 8.0, // ______--______--
}

// Any one of the 4 Game Boy sound channels
type SoundChannel struct {
	// TODO: could use some cleanup
	// Sets whether this channel is pulse (channels 1 or 2), wave (channel 3), or noise (channel 4)
	channelNumber int

	// Whether or not the sound channel is enabled
	on bool

	// Channel 1 only
	// Currently selected wave duty cycle (0-3)
	duty uint8
	// Channel 1 sweep control register
	nr10RegisterValue uint8
	// Time between frequency changes in 128Hz ticks
	// 0-7, 0=no update 1=fast change 7=slow change
	sweepTime uint8
	// Time since starting the frequency sweep (s)
	sweepTimeCounter float64
	// 0-7, 0=fast change 7=slow change
	sweepSlope uint8
	// 0: increasing, 1: decreasing
	sweepDirection uint8

	// 11-bit value which corresponds to frequency
	frequencyValue uint16

	// Fractional number of waves we have generated up to this point
	waveCounter float64

	// Volume and Envelope control register (Channels 1, 2, 4 only)
	nrX2RegisterValue uint8

	// If sound length is enabled, initial sound length defines start point for increasing sound counter
	// Ch 1, 2, 4: 0-63, 0 = long sound, 63 = short sound
	// Ch 3: 0-255, 0 = long sound, 255 = short sound
	initialSoundLength uint8
	soundLengthEnable  bool
	// Time since starting this sound (s)
	lengthCounter float64

	// 4 bit 0-F, 0=sound off F=full volume
	initialVolumeEnvelope uint8
	// 1 bit, 0=decreasing 1=increasing
	volumeEnvelopeDirection uint8
	// 3 bit 0-7, 0=no volume change, 1=fast change 7=slow change
	volumeSweepPace       uint8
	currentVolume         uint8
	volumeEnvelopeCounter float64

	// Selects volume scaling (channel 3 only)
	// 2 bit: 0 = mute, 1 = 100%, 2 = 50%, 3 = 25%
	outputLevel uint8
	// Used for channel 3 only
	waveRAM []uint8

	// Channel 4 only
	// LFSR for noise generation
	shiftRegister uint16
	// 0 = 15 bits, 1 = 7 bits
	shiftRegisterWidth uint8
	// 4 bit clock shift (s)
	shiftRegisterClockShift uint8
	// 3 bit clock divider (r)
	shiftRegisterClockRatio uint8
	// Noise frequency is f = 262144 / (max(r, 0.5) * 2**s)
}

// Get the next sample for a pulse channel
func (c *SoundChannel) pulseGenerator(t float64) uint8 {
	// Where t is the wave we are on (2.5 = half way through the 3rd wave)
	_, frac := math.Modf(t)
	if WaveDuty[c.duty] < frac {
		return 0xFF
	}
	return 0
}

// Get the next sample for a wave output channel
func (c *SoundChannel) waveGenerator(t float64) uint8 {
	// Where t is the wave we are on (2.5 = half way through the 3rd wave)
	_, frac := math.Modf(t)
	sampleIndex := int(frac * WaveRAMSize)
	waveByte := c.waveRAM[sampleIndex]

	// Each sample is 4 bits each, we read through Wave RAM high nybble first
	var value uint8
	if int(frac*WaveRAMSize*2)%2 == 0 {
		// In the first half of this byte
		value = waveByte & 0xF0
	} else {
		// In the second half of this byte
		value = waveByte << 4
	}

	return value
}

// Get the next sample for the noise channel
func (c *SoundChannel) noiseGenerator() uint8 {
	if (c.shiftRegister & 1) == 0 {
		return 0
	} else {
		return 0xFF
	}
}

// Trigger the sound channel
func (c *SoundChannel) Trigger() {
	c.applyNR10()
	c.applyNRx2()

	c.on = true
	c.lengthCounter = float64(c.initialSoundLength) * LengthTickTime
	c.currentVolume = c.initialVolumeEnvelope
	c.volumeEnvelopeCounter = 0
	c.shiftRegister = 0
	c.sweepTimeCounter = 0
}

// Channel 4 only, apply 1 shift and feeback step to the noise LFSR
func (c *SoundChannel) updateShiftRegister() {
	// xnor bit 0 with bit 1
	value := ^(c.shiftRegister ^ (c.shiftRegister >> 1)) & 1

	// Set bit 15
	c.shiftRegister = (c.shiftRegister & 0x7FFF) | (value << 15)
	if c.shiftRegisterWidth == 1 {
		// Set bit 7
		c.shiftRegister = (c.shiftRegister & 0xFF7F) | (value << 7)
	}

	// Shift the register
	c.shiftRegister >>= 1
}

// Parse the value in NR10 to assign channel sweep control values
func (c *SoundChannel) applyNR10() {
	// Sweep Time: bits 4-6
	c.sweepTime = (c.nr10RegisterValue >> 4) & 0b111
	// Sweep Direction: bit 3
	c.sweepDirection = (c.nr10RegisterValue >> 3) & 0b1
	// Sweep Slope: bits 0-2
	c.sweepSlope = c.nr10RegisterValue & 0b111
}

// Parse the value in NR12/22/42 to assign channel volume envelope values
func (c *SoundChannel) applyNRx2() {
	// Initial volume of envelope: bits 4-7
	c.initialVolumeEnvelope = (c.nrX2RegisterValue >> 4) & 0b1111
	// Envelope Direction: bit 3
	c.volumeEnvelopeDirection = (c.nrX2RegisterValue >> 3) & 1
	// Volume Sweep Pace: bits 0-2
	c.volumeSweepPace = c.nrX2RegisterValue & 0b111
}

// Sample the audio channel at the global audio sample rate
func (c *SoundChannel) GetSample() uint8 {
	if !c.on {
		return 0
	}

	// Get sample
	var output uint8

	// Channel 1 and 2: Each wave is 8 "samples" long and channels are clocked at 1048576 Hz -> 1048576 / 8 = 131072
	// Channel 3: Each wave is 32 "samples" long and channels are clocked at 2097152 Hz -> 2097152 / 32 = 65536
	// TODO: these are just powers of two off, same with noise frequency, should make use of that to simplify
	var chanFrequency float64
	if c.channelNumber == 1 || c.channelNumber == 2 {
		// frequencyValue rolls over each time it reaches 2048, triggering 1 "sample"
		chanFrequency = 131072.0 / (2048.0 - float64(c.frequencyValue))
	} else if c.channelNumber == 3 {
		chanFrequency = 65536.0 / (2048.0 - float64(c.frequencyValue))
	} else if c.channelNumber == 4 {
		var s int = 1 << c.shiftRegisterClockShift
		chanFrequency = 262144.0 / (math.Max(float64(c.shiftRegisterClockRatio), 0.5) * float64(s))
	}
	// Step is how far through 1 wave we got with this sample
	// If frequency is lower, step will be lower, so we step through a wave more slowly
	step := chanFrequency / float64(AudioSampleRate)
	c.waveCounter += step

	// Take the sample value from the generator at the current point along our wave
	if c.channelNumber == 1 || c.channelNumber == 2 {
		output = uint8(float64(c.pulseGenerator(c.waveCounter)) * float64(c.currentVolume) / 0xF)
	} else if c.channelNumber == 3 {
		if c.outputLevel != 0 {
			// Scale output volume based on selected output level
			output = c.waveGenerator(c.waveCounter) >> (c.outputLevel - 1)
		}
	} else if c.channelNumber == 4 {
		ipart := int(c.waveCounter)

		for i := 0; i < ipart; i++ {
			c.updateShiftRegister()
		}
		// Take the sample value from the generator at the current point along our wave
		output = uint8(float64(c.noiseGenerator()) * float64(c.currentVolume) / 0xF)

		// Remove the whole number of waves completed, since we've already run the shift register for them
		c.waveCounter -= float64(ipart)
	}

	// Update the frequency sweep
	if c.sweepTime > 0 {
		c.sweepTimeCounter += 1.0 / AudioSampleRate
		numTicks := c.sweepTimeCounter / SweepTickTime
		if numTicks >= float64(c.sweepTime) {
			// Period should get changed
			// F = F + F / 2^n
			periodChange := c.frequencyValue >> uint16(c.sweepSlope)
			if c.sweepDirection == 0 {
				// increasing sweep
				c.frequencyValue += periodChange
				if c.frequencyValue > 0x7FF {
					c.frequencyValue = 0x7FF
					c.on = false
				}
			} else {
				// decreasing sweep
				c.frequencyValue -= periodChange
			}
			// Reset the sweep time counter
			c.sweepTimeCounter -= float64(c.sweepTime) * SweepTickTime
		}
	}

	// Update the volume envelope
	if c.volumeSweepPace > 0 {
		c.volumeEnvelopeCounter += 1.0 / AudioSampleRate
		numTicks := c.volumeEnvelopeCounter / VolumeTickTime
		if numTicks >= float64(c.volumeSweepPace) {
			// Volume should get changed
			if c.volumeEnvelopeDirection == 0 {
				// Volume should decrease
				if c.currentVolume > 0 {
					c.currentVolume -= 1
				}
			} else if c.currentVolume < 0xF {
				// Volume should increase
				c.currentVolume += 1
			}
			c.volumeEnvelopeCounter -= float64(c.volumeSweepPace) * VolumeTickTime
		}
	}

	// Update the sound play length
	c.lengthCounter += 1.0 / AudioSampleRate
	var lengthMax float64
	if c.channelNumber == 1 || c.channelNumber == 2 {
		lengthMax = 64
	} else if c.channelNumber == 3 {
		lengthMax = 256
	}
	if c.soundLengthEnable && (c.lengthCounter/LengthTickTime) >= lengthMax {
		// If sound length is enabled and
		c.on = false
	}

	return output
}
