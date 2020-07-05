package pedal

import (
	"fmt"
	"time"
)

// NewMorseSampler samples time.Durations, as received from the IntervalSampler, to morse code words.
func NewMorseSampler(inputChan chan interface{}, maxDotDuration, minIdleDuration time.Duration) (sampler *Sampler) {
	sampler = &Sampler{
		inputChan:  inputChan,
		outputChan: make(chan interface{}),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go func(sampler *Sampler, maxDotDuration, minIdleDuration time.Duration) {
		defer close(sampler.stopAck)

		var tmpWord string
		var lastInput time.Time

		for {
			select {
			case <-sampler.stopSyn:
				return

			case <-time.After(minIdleDuration / 4):
				if tmpWord != "" && lastInput.Add(minIdleDuration).Before(time.Now()) {
					sampler.outputChan <- tmpWord

					tmpWord = ""
					lastInput = time.Time{}
				}

			case input := <-sampler.inputChan:
				switch input := input.(type) {
				case time.Duration:
					if input <= maxDotDuration {
						tmpWord += "."
					} else {
						tmpWord += "_"
					}
					lastInput = time.Now()

				case error:
					sampler.outputChan <- input
					return

				default:
					sampler.outputChan <- fmt.Errorf("morseSampler: unsupported type %T", input)
					return
				}
			}
		}
	}(sampler, maxDotDuration, minIdleDuration)

	return
}
