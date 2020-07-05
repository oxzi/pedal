package pedal

import "time"

// NewIntervalSampler samples the Signaler's channel to intervals (time.Duration).
func NewIntervalSampler(inputChan chan interface{}, samplingRate time.Duration) (sampler *Sampler) {
	sampler = &Sampler{
		inputChan:  inputChan,
		outputChan: make(chan interface{}),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go func(sampler *Sampler, samplingRate time.Duration) {
		defer close(sampler.stopAck)

		var firstInput, lastInput time.Time

		for {
			select {
			case <-sampler.stopSyn:
				return

			case <-time.After(samplingRate):
				if lastInput != (time.Time{}) && lastInput.Add(2*samplingRate).Before(time.Now()) {
					sampler.outputChan <- lastInput.Sub(firstInput)

					firstInput = time.Time{}
					lastInput = time.Time{}
				}

			case input := <-sampler.inputChan:
				if input != nil {
					sampler.outputChan <- input
					return
				}

				if firstInput == (time.Time{}) {
					firstInput = time.Now()
				}
				lastInput = time.Now()
			}
		}
	}(sampler, samplingRate)

	return
}
