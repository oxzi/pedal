package main

import (
	"fmt"
	"time"
)

func main() {
	const samplingRate = 50 * time.Millisecond
	const dashMinDuration = 250 * time.Millisecond
	const morseMaxDuration = time.Second

	signaler, err := NewSignaler("/dev/ttyUSB0", samplingRate)
	if err != nil {
		panic(err)
	}

	intervals := intervalSampler(signaler.Chan(), samplingRate)
	morse := morseSampler(intervals, dashMinDuration)
	morseWord := morseWordSampler(morse, morseMaxDuration)

	go func() {
		for msg := range morseWord {
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
