package main

import (
	"fmt"
	"time"
)

// Sampler is an abstract kind of struct to unify different kind of Samplers, e.g., for intervals.
type Sampler struct {
	inputChan  chan interface{}
	outputChan chan interface{}

	stopSyn chan struct{}
	stopAck chan struct{}
}

// Chan is the output channel.
func (sampler *Sampler) Chan() chan interface{} {
	return sampler.outputChan
}

// Close this Sampler and notify the worker process to stop reading from the input channel.
func (sampler *Sampler) Close() error {
	close(sampler.stopSyn)
	<-sampler.stopAck

	return nil
}

// NewIntervalSampler samples the Signaler's channel to intervals (time.Duration).
func NewIntervalSampler(inputChan chan interface{}, samplingRate time.Duration) (sampler *Sampler) {
	sampler = &Sampler{
		inputChan:  inputChan,
		outputChan: make(chan interface{}),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go func(sampler *Sampler, samplingRate time.Duration) {
		var firstInput, lastInput time.Time

		for {
			select {
			case <-sampler.stopSyn:
				close(sampler.stopAck)
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

// NewMorseSampler samples time.Durations, as received from the IntervalSampler, to morse code words.
func NewMorseSampler(inputChan chan interface{}, maxDotDuration, minIdleDuration time.Duration) (sampler *Sampler) {
	sampler = &Sampler{
		inputChan:  inputChan,
		outputChan: make(chan interface{}),

		stopSyn: make(chan struct{}),
		stopAck: make(chan struct{}),
	}

	go func(sampler *Sampler, maxDotDuration, minIdleDuration time.Duration) {
		var tmpWord string
		var lastInput time.Time

		for {
			select {
			case <-sampler.stopSyn:
				close(sampler.stopAck)
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
