package main

import (
	"fmt"
	"time"

	"github.com/geistesk/pedal/modes"
	"github.com/geistesk/pedal/pedal"
)

func main() {
	const samplingRate = 50 * time.Millisecond
	const morseMaxUnit = 250 * time.Millisecond

	signaler, err := pedal.NewSignaler("/dev/ttyUSB0", samplingRate)
	if err != nil {
		panic(err)
	}

	/*
		micMuteAction := modes.NewCommandAction("amixer -c 0 set Capture toggle")
		kbdSpaceAction, kbdErr := modes.NewKeyboardPressAction([]int{57})
		if kbdErr != nil {
			panic(kbdErr)
		}

		actionMap := map[string]modes.Action{
			".": micMuteAction,
			"_": kbdSpaceAction,
		}
	*/

	// multiAction := modes.NewMultiAction(signaler.Chan(), actionMap, samplingRate, morseMaxDotDuration, morseMinIdleDuration)
	// trigger := modes.NewTrigger(signaler.Chan(), kbdSpaceAction, 500*time.Millisecond)

	morseAction, morseErr := modes.NewMorseKeyboard(signaler.Chan(), samplingRate, morseMaxUnit)
	if morseErr != nil {
		panic(morseErr)
	}

	go func() {
		for err := range morseAction.Errors() {
			panic(err)
		}
	}()

	time.Sleep(30 * time.Second)
	fmt.Println("Closing down..")

	if err := morseAction.Close(); err != nil {
		panic(err)
	}
	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
