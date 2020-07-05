package pedal

import "time"

// NewCooldownSampler discards inputs for a cooldown time, excluding errors.
func NewCooldownSampler(inputChan chan interface{}, cooldownDuration time.Duration) (sampler *Sampler) {
	sampler = &Sampler{
		inputChan:  inputChan,
		outputChan: make(chan interface{}),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go func(sampler *Sampler, cooldownDuration time.Duration) {
		defer close(sampler.stopAck)

		var lastInput time.Time

		for {
			select {
			case <-sampler.stopSyn:
				return

			case input := <-sampler.inputChan:
				switch input := input.(type) {
				case error:
					sampler.outputChan <- input
					return

				default:
					if now := time.Now(); lastInput.Add(cooldownDuration).Before(now) {
						sampler.outputChan <- input
						lastInput = now
					}
				}
			}
		}
	}(sampler, cooldownDuration)

	return
}
