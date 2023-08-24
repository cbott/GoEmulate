package main

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/oto"
)

// Actual sample rate for sound played out of your speakers
const AudioSampleRate = 44100
const CounterMaximum = 2048
const Freq = 1048576 // Hz, clock rate of Game Boy audio channels
// Game Boy CPU cycle rate
const FPS = 60
const CyclesPerSample = CpuSpeed / float64(AudioSampleRate)
const SweepTickLength = 1 / 128.0  // Seconds
const LengthTickLength = 1 / 256.0 // Seconds
const VolumeTickLength = 1 / 64.0  // Seconds

const bufferSeconds = 120

var (
	mu          sync.Mutex
	player      *oto.Player
	audioBuffer chan [2]uint8
)

var DutyCycle = map[uint8]float64{
	0b00: 7 / 8.0,
	0b01: 6 / 8.0,
	0b10: 4 / 8.0,
	0b11: 2 / 8.0,
}

type Channel1 struct {
	running bool
	duty    uint8
	// 11-bit period value
	periodValue  uint16
	cycleCounter float64
	time         float64
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

func getSampleRate(periodValue uint16) float64 {
	return Freq / (CounterMaximum - float64(periodValue))
}

func (c *Channel1) processCPUCycles(cycles int) {
	c.cycleCounter += float64(cycles)
	if c.cycleCounter < CyclesPerSample {
		return
	}
	c.cycleCounter -= CyclesPerSample

	// Time to take a sample
	left, right := c.GetSample()
	audioBuffer <- [2]uint8{left, right}
}

func (c *Channel1) generator(t float64) byte {
	// Where t is the wave we are on (2.5 = half way through the 3rd wave)
	_, frac := math.Modf(t)
	if DutyCycle[c.duty] < frac {
		return 0xFF
	}
	return 0
}

func (c *Channel1) Trigger() {
	c.running = true
	c.lengthCounter = float64(c.initialSoundLength) * LengthTickLength
	c.currentVolume = c.initialVolumeEnvelope
	c.volumeEnvelopeCounter = 0
	c.sweepTimeCounter = 0
}

// Get sample for left and right, at the global audio sample rate
func (c *Channel1) GetSample() (uint8, uint8) {
	if !c.running {
		return 0, 0
	}

	// Get sample
	var output uint8
	chanFrequency := 131072 / (CounterMaximum - float64(c.periodValue))
	// Step, radians, how far through 1 wave we got with this sample
	// If frequency is lower, step will be lower, so we step through a wave more slowly
	step := chanFrequency / float64(AudioSampleRate)
	// time isn't really time, it's cumulative number of waves we've created
	c.time += step
	// Take the sample value from the generator at the current point along our wave
	output = uint8(float64(c.generator(c.time)) * float64(c.currentVolume) / 0xF)

	// Update the frequency sweep
	if c.sweepTime > 0 {
		c.sweepTimeCounter += 1.0 / AudioSampleRate
		numTicks := c.sweepTimeCounter / SweepTickLength
		if numTicks >= float64(c.sweepTime) {
			// Period should get changed
			// F = F + F / 2^n
			periodChange := c.periodValue >> uint16(c.sweepSlope)
			if c.sweepDirection == 0 {
				// increasing sweep
				c.periodValue += periodChange
				if c.periodValue > 0x7FF {
					c.periodValue = 0x7FF
					c.running = false
				}
			} else {
				// decreasing sweep
				c.periodValue -= periodChange
			}
			// Reset the sweep time counter
			c.sweepTimeCounter -= float64(c.sweepTime) * SweepTickLength
		}
	}

	// Update the volume envelope
	if c.volumeSweepPace > 0 {
		c.volumeEnvelopeCounter += 1.0 / AudioSampleRate
		numTicks := c.volumeEnvelopeCounter / VolumeTickLength
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
			c.volumeEnvelopeCounter -= float64(c.volumeSweepPace) * VolumeTickLength
		}
	}

	// Update the sound play length
	c.lengthCounter += 1.0 / AudioSampleRate
	if c.soundLengthEnable && (c.lengthCounter/LengthTickLength) >= 64 {
		// If sound length is enabled and
		c.running = false
	}

	return output, output
}

func (n Channel1) Err() error {
	return nil
}

const (
	// 2 channels: game boy supports stereo sound
	ChannelsPerSample = 2
	// 1-byte sound depth: tyring this for now, don't think game boy has higher fidelity
	BytesPerChannel = 1
	BytesPerSample  = ChannelsPerSample * BytesPerChannel
)

func main() {
	// Arbitrarily going to choose to buffer 200 samples which is about 1/4 frame
	var bufferSamples int = 200
	context, err := oto.NewContext(AudioSampleRate, ChannelsPerSample, BytesPerChannel, bufferSamples*BytesPerSample)
	if err != nil {
		log.Fatalf("Audio initialization error: %v", err)
	}
	player = context.NewPlayer()

	// // At most we will allow up to 5000 samples to get backed up, hopefully this won't happen
	audioBuffer = make(chan [2]uint8, 5000)

	frameTime := time.Second / time.Duration(bufferSeconds)
	ticker := time.NewTicker(frameTime)
	targetSamples := AudioSampleRate / bufferSeconds
	go func() {
		var reading [2]byte
		var buffer []byte
		for range ticker.C { // 1/120s
			fbLen := len(audioBuffer)
			if fbLen >= targetSamples/2 { // 184
				newBuffer := make([]byte, fbLen*2)
				for i := 0; i < fbLen*2; i += 2 {
					reading = <-audioBuffer
					newBuffer[i], newBuffer[i+1] = reading[0], reading[1]
				}
				buffer = newBuffer
			}

			_, err := player.Write(buffer)
			if err != nil {
				log.Printf("error sampling: %v", err)
			}
		}
	}()

	myChan := Channel1{periodValue: 0x600, cycleCounter: 0}

	myChan.sweepTime = 0b000
	myChan.sweepSlope = 0b111
	myChan.sweepDirection = 0
	myChan.duty = 0b10

	myChan.initialSoundLength = 0b000000
	myChan.soundLengthEnable = false

	myChan.initialVolumeEnvelope = 0b0000
	myChan.volumeEnvelopeDirection = 1
	myChan.volumeSweepPace = 0b001

	myChan.Trigger()

	fpsTicker := time.NewTicker(time.Second / FPS)
	for range fpsTicker.C {
		for j := 0; j < 69905; j++ {
			// Number of cycles per frame
			myChan.processCPUCycles(1)
		}
	}
}
