package main

import (
	"fmt"
	"time"
)

func main() {
	const samplingRate = 50 * time.Millisecond
	const morseMaxDotDuration = 300 * time.Millisecond
	const morseMinIdleDuration = 1000 * time.Millisecond

	signaler, err := NewSignaler("/dev/ttyUSB0", samplingRate)
	if err != nil {
		panic(err)
	}
	intervalSampler := NewIntervalSampler(signaler.Chan(), samplingRate)
	morseSampler := NewMorseSampler(intervalSampler.Chan(), morseMaxDotDuration, morseMinIdleDuration)

	go func() {
		for msg := range morseSampler.Chan() {
			switch msg := msg.(type) {
			case error:
				panic(msg)

			default:
				fmt.Println(msg)
			}
		}
	}()

	time.Sleep(30 * time.Second)

	morseSampler.Close()
	intervalSampler.Close()

	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
