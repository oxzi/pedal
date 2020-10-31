package modes

import (
	"time"

	"github.com/oxzi/pedal/pedal"
)

// Trigger executes an Action after pressing the pedal with a configurable cooldown.
type Trigger struct {
	action   Action
	cooldown *pedal.Sampler

	errors chan error

	stopSyn chan struct{}
	stopAck chan struct{}
}

// NewTrigger creates a Trigger from a pedal's signal, some Action and a cooldown.
func NewTrigger(signalChan chan interface{}, action Action, cooldownDuration time.Duration) (trigger *Trigger) {
	trigger = &Trigger{
		action:   action,
		cooldown: pedal.NewCooldownSampler(signalChan, cooldownDuration),

		errors: make(chan error),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go trigger.worker()

	return
}

// Errors is channel to pass raising errors.
func (trigger *Trigger) Errors() chan error {
	return trigger.errors
}

// Close this Mode and all its internal workers.
func (trigger *Trigger) Close() (err error) {
	close(trigger.stopSyn)
	<-trigger.stopAck

	close(trigger.errors)

	if err = trigger.cooldown.Close(); err != nil {
		return
	}

	return
}

// worker is the internal worker routine.
func (trigger *Trigger) worker() {
	defer close(trigger.stopAck)

	for {
		select {
		case <-trigger.stopSyn:
			return

		case input := <-trigger.cooldown.Chan():
			switch input := input.(type) {
			case error:
				trigger.errors <- input
				return

			default:
				if actionErr := trigger.action.Execute(); actionErr != nil {
					trigger.errors <- actionErr
					return
				}
			}
		}
	}
}
