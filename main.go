package main

import (
	"fmt"
	"time"

	"github.com/geistesk/pedal/pedal"
)

func main() {
	const samplingRate = 50 * time.Millisecond
	const morseMaxDotDuration = 300 * time.Millisecond
	const morseMinIdleDuration = 1000 * time.Millisecond

	signaler, err := pedal.NewSignaler("/dev/ttyUSB0", samplingRate)
	if err != nil {
		panic(err)
	}

	cooldownSampler := pedal.NewCooldownSampler(signaler.Chan(), time.Second)

	/*
		intervalSampler := NewIntervalSampler(signaler.Chan(), samplingRate)
		morseSampler := NewMorseSampler(intervalSampler.Chan(), morseMaxDotDuration, morseMinIdleDuration)
	*/

	go func() {
		for msg := range cooldownSampler.Chan() {
			switch msg := msg.(type) {
			case error:
				panic(msg)

			default:
				fmt.Println(msg)
			}
		}
	}()

	time.Sleep(30 * time.Second)
	fmt.Println("Closing down..")

	/*
		morseSampler.Close()
		intervalSampler.Close()
	*/
	cooldownSampler.Close()

	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
