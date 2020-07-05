package modes

import (
	"os/exec"
	"time"

	"github.com/geistesk/pedal/pedal"
)

// Trigger executes a command after pressing the pedal with a configurable cooldown.
type Trigger struct {
	commandStr string
	cooldown   *pedal.Sampler

	errors chan error

	stopSyn chan struct{}
	stopAck chan struct{}
}

// NewTrigger creates a Trigger from a pedal's signal, some (shell) command and a cooldown.
func NewTrigger(signalChan chan interface{}, commandStr string, cooldownDuration time.Duration) (trigger *Trigger) {
	trigger = &Trigger{
		commandStr: commandStr,
		cooldown:   pedal.NewCooldownSampler(signalChan, cooldownDuration),

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
func (trigger *Trigger) Close() error {
	close(trigger.stopSyn)
	<-trigger.stopAck

	return trigger.cooldown.Close()
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
				cmd := exec.Command("sh", "-c", trigger.commandStr)
				if cmdErr := cmd.Start(); cmdErr != nil {
					trigger.errors <- cmdErr
					return
				}
			}
		}
	}
}
