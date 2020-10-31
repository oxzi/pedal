package main

import (
	"time"

	"github.com/oxzi/pedal/modes"
)

const serverSocket = "/tmp/pedal.sock"

// samplingRate to be used for both Samplers and Modes.
const samplingRate = 50 * time.Millisecond

// morseMaxUnit is the maximum length of a Morse unit, equals a dot.
const morseMaxUnit = 250 * time.Millisecond

// modeMessages are the supported Modes by their IPC name.
var modeMessages = map[string](func(chan interface{}) (modes.Mode, error)){
	"mic-toggle": func(signalerChan chan interface{}) (modes.Mode, error) {
		micMuteAction := modes.NewCommandAction("amixer -c 0 set Capture toggle")
		return modes.NewTrigger(signalerChan, micMuteAction, 500*time.Millisecond), nil
	},

	"morse-keyboard": func(signalerChan chan interface{}) (modes.Mode, error) {
		return modes.NewMorseKeyboard(signalerChan, samplingRate, morseMaxUnit)
	},

	"scroll-down": func(signalerChan chan interface{}) (modes.Mode, error) {
		if kbdDownAction, err := modes.NewKeyboardPressAction([]int{108}); err != nil {
			return nil, err
		} else {
			return modes.NewTrigger(signalerChan, kbdDownAction, 200*time.Millisecond), nil
		}
	},

	"arrow-right": func(signalerChan chan interface{}) (modes.Mode, error) {
		if kbdDownAction, err := modes.NewKeyboardPressAction([]int{106}); err != nil {
			return nil, err
		} else {
			return modes.NewTrigger(signalerChan, kbdDownAction, 500*time.Millisecond), nil
		}
	},
}
