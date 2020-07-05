package main

import (
	"fmt"
	"time"

	"github.com/geistesk/pedal/modes"
	"github.com/geistesk/pedal/pedal"
)

func main() {
	const samplingRate = 50 * time.Millisecond
	// const morseMaxDotDuration = 300 * time.Millisecond
	// const morseMinIdleDuration = 1000 * time.Millisecond

	signaler, err := pedal.NewSignaler("/dev/ttyUSB0", samplingRate)
	if err != nil {
		panic(err)
	}

	micMuteAction := modes.NewCommandAction("amixer -c 0 set Capture toggle")
	trigger := modes.NewTrigger(signaler.Chan(), micMuteAction, 500*time.Millisecond)

	go func() {
		for err := range trigger.Errors() {
			panic(err)
		}
	}()

	time.Sleep(30 * time.Second)
	fmt.Println("Closing down..")

	if err := trigger.Close(); err != nil {
		panic(err)
	}
	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
