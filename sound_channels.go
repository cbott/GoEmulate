package main

import (
	"math"
)

const (
	LengthTickTime = 1 / 256.0 // Seconds
	SweepTickTime  = 1 / 128.0 // Seconds
	VolumeTickTime = 1 / 64.0  // Seconds
)

// WaveDuty defines the fraction of the wave which is high vs low
var WaveDuty = map[uint8]float64{
	0b00: 7 / 8.0, // _-------_-------
	0b01: 6 / 8.0, // __------__------
	0b10: 4 / 8.0, // ____----____----
	0b11: 2 / 8.0, // ______--______--
}

type SoundChannel struct {
	// Whether or not the sound channel is enabled
	on bool

	// Currently selected wave duty cycle (0-3)
	duty uint8

	// 11-bit value which corresponds to frequency as 131072/(2048-x) Hz
	frequencyValue uint16

	// Fractional number of waves we have generated up to this point
	waveCounter float64

	// Time between frequency changes in 128Hz ticks
	// 0-7, 0=no update 1=fast change 7=slow change
	sweepTime uint8
	// Time since starting the frequency sweep (s)
	sweepTimeCounter float64
	// 0-7, 0=fast change 7=slow change
	sweepSlope uint8
	// 0: increasing, 1: decreasing
	sweepDirection uint8

	// If sound length is enabled, initial sound length defines start point for increasing sound counter
	// 0-63, 0 = long sound, 63 = short sound
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
}

func (c *SoundChannel) generator(t float64) byte {
	// Where t is the wave we are on (2.5 = half way through the 3rd wave)
	_, frac := math.Modf(t)
	if WaveDuty[c.duty] < frac {
		return 0xFF
	}
	return 0
}

func (c *SoundChannel) Trigger() {
	c.on = true
	c.lengthCounter = float64(c.initialSoundLength) * LengthTickTime
	c.currentVolume = c.initialVolumeEnvelope
	c.volumeEnvelopeCounter = 0
	c.sweepTimeCounter = 0
}

// Sample the audio channel at the global audio sample rate
func (c *SoundChannel) GetSample() uint8 {
	if !c.on {
		return 0
	}

	// Get sample
	var output uint8
	// Each wave is 8 "samples" long and channels are clocked at 1048576 Hz -> 1048576 / 8 = 131072
	// frequencyValue rolls over each time it reaches 2048, triggering 1 "sample"
	chanFrequency := 131072.0 / (2048.0 - float64(c.frequencyValue))
	// Step is how far through 1 wave we got with this sample
	// If frequency is lower, step will be lower, so we step through a wave more slowly
	step := chanFrequency / float64(AudioSampleRate)
	c.waveCounter += step
	// Take the sample value from the generator at the current point along our wave
	output = uint8(float64(c.generator(c.waveCounter)) * float64(c.currentVolume) / 0xF)

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
	if c.soundLengthEnable && (c.lengthCounter/LengthTickTime) >= 64 {
		// If sound length is enabled and
		c.on = false
	}

	return output
}
