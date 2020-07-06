package pedal

import (
	"fmt"
	"time"
)

// NewMorseSampler samples time.Durations, as received from the IntervalSampler, to morse code words.
//
// The maxUnit describes the maximum time for one unit, which equals a dot. A dash might take up to
// four units. Longer delays are interpreted as the character's end.
func NewMorseSampler(inputChan chan interface{}, maxUnit time.Duration) (sampler *Sampler) {
	sampler = &Sampler{
		inputChan:  inputChan,
		outputChan: make(chan interface{}),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go func(sampler *Sampler, maxUnit time.Duration) {
		defer close(sampler.stopAck)

		var tmpWord string
		var lastInput time.Time

		for {
			select {
			case <-sampler.stopSyn:
				return

			case <-time.After(maxUnit):
				if tmpWord != "" && lastInput.Add(4*maxUnit).Before(time.Now()) {
					sampler.outputChan <- tmpWord

					tmpWord = ""
					lastInput = time.Time{}
				}

			case input := <-sampler.inputChan:
				switch input := input.(type) {
				case time.Duration:
					if input <= maxUnit {
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
	}(sampler, maxUnit)

	return
}
