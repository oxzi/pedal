package main

import (
	"fmt"
	"time"
)

func main() {
	const samplingRate = 50 * time.Millisecond
	const morseMaxDotDuration = 250 * time.Millisecond
	const morseMinIdleDuration = 1000 * time.Millisecond

	signaler, err := NewSignaler("/dev/ttyUSB0", samplingRate)
	if err != nil {
		panic(err)
	}

	intervals := intervalSampler(signaler.Chan(), samplingRate)
	morse := morseSampler(intervals, morseMaxDotDuration, morseMinIdleDuration)

	go func() {
		for msg := range morse {
			switch msg := msg.(type) {
			case error:
				panic(msg)

			default:
				fmt.Println(msg)
			}
		}
	}()

	time.Sleep(10 * time.Second)

	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
