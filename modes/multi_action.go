package modes

import (
	"fmt"
	"time"

	"github.com/geistesk/pedal/pedal"
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
func NewMultiAction(signalChan chan interface{}, actions map[string]Action, samplingRate, maxDotDuration, minIdleDuration time.Duration) (multiAction *MultiAction) {
	multiAction = &MultiAction{
		actions: actions,

		errors: make(chan error),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	multiAction.interval = pedal.NewIntervalSampler(signalChan, samplingRate)
	multiAction.morse = pedal.NewMorseSampler(multiAction.interval.Chan(), maxDotDuration, minIdleDuration)

	go multiAction.worker()

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
