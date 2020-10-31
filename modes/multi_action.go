package modes

import (
	"fmt"
	"time"

	"github.com/oxzi/pedal/pedal"

	"github.com/micmonay/keybd_event"
)

// MultiAction performs different action based on a Morse code input.
type MultiAction struct {
	actions map[string]Action

	interval *pedal.Sampler
	morse    *pedal.Sampler

	errors chan error

	stopSyn chan struct{}
	stopAck chan struct{}
}

// NewMultiAction based on a map of Morse codes to Actions, the Signaler's sampling rate and durations for Morse code parsing.
func NewMultiAction(signalChan chan interface{}, actions map[string]Action, samplingRate, maxUnit time.Duration) (multiAction *MultiAction) {
	multiAction = &MultiAction{
		actions: actions,

		errors: make(chan error),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	multiAction.interval = pedal.NewIntervalSampler(signalChan, samplingRate)
	multiAction.morse = pedal.NewMorseSampler(multiAction.interval.Chan(), maxUnit)

	go multiAction.worker()

	return
}

func NewMorseKeyboard(signalChan chan interface{}, samplingRate, maxUnit time.Duration) (multiAction *MultiAction, err error) {
	morseMap := map[string]int{
		"._":    keybd_event.VK_A,
		"_...":  keybd_event.VK_B,
		"_._.":  keybd_event.VK_C,
		"_..":   keybd_event.VK_D,
		".":     keybd_event.VK_E,
		".._.":  keybd_event.VK_F,
		"__.":   keybd_event.VK_G,
		"....":  keybd_event.VK_H,
		"..":    keybd_event.VK_I,
		".___":  keybd_event.VK_J,
		"_._":   keybd_event.VK_K,
		"._..":  keybd_event.VK_L,
		"__":    keybd_event.VK_M,
		"_.":    keybd_event.VK_N,
		"___":   keybd_event.VK_O,
		".__.":  keybd_event.VK_P,
		"__._":  keybd_event.VK_Q,
		"._.":   keybd_event.VK_R,
		"...":   keybd_event.VK_S,
		"_":     keybd_event.VK_T,
		".._":   keybd_event.VK_U,
		"..._":  keybd_event.VK_V,
		".__":   keybd_event.VK_W,
		"_.._":  keybd_event.VK_X,
		"_.__":  keybd_event.VK_Y,
		"__..":  keybd_event.VK_Z,
		".____": keybd_event.VK_1,
		"..___": keybd_event.VK_2,
		"...__": keybd_event.VK_3,
		"...._": keybd_event.VK_4,
		".....": keybd_event.VK_5,
		"_....": keybd_event.VK_6,
		"__...": keybd_event.VK_7,
		"___..": keybd_event.VK_8,
		"____.": keybd_event.VK_9,
		"_____": keybd_event.VK_0,
	}

	actions := map[string]Action{}
	for k, v := range morseMap {
		if kbdAction, kbdErr := NewKeyboardPressAction([]int{v}); kbdErr != nil {
			err = kbdErr
			return
		} else {
			actions[k] = kbdAction
		}
	}

	multiAction = NewMultiAction(signalChan, actions, samplingRate, maxUnit)
	return
}

// Errors is channel to pass raising errors.
func (multiAction *MultiAction) Errors() chan error {
	return multiAction.errors
}

// Close this Mode and all its internal workers.
func (multiAction *MultiAction) Close() (err error) {
	close(multiAction.stopSyn)
	<-multiAction.stopAck

	close(multiAction.errors)

	if err = multiAction.morse.Close(); err != nil {
		return
	}
	if err = multiAction.interval.Close(); err != nil {
		return
	}

	return
}

// worker to execute Actions based on the parsed Morse code.
func (multiAction *MultiAction) worker() {
	defer close(multiAction.stopAck)

	for {
		select {
		case <-multiAction.stopSyn:
			return

		case input := <-multiAction.morse.Chan():
			switch input := input.(type) {
			case string:
				if action, ok := multiAction.actions[input]; ok {
					if actionErr := action.Execute(); actionErr != nil {
						multiAction.errors <- actionErr
						return
					}
				}

			case error:
				multiAction.errors <- input
				return

			default:
				multiAction.errors <- fmt.Errorf("MultiAction: unsupported type %T", input)
				return
			}
		}
	}
}
