package main

import (
	"fmt"
	"time"
)

// intervalSampler reads the Signaler's channel and samples its values to intervals (time.Duration) or passes errors.
func intervalSampler(inputChan chan error, samplingRate time.Duration) (outputChan chan interface{}) {
	outputChan = make(chan interface{})

	go func() {
		var firstInput, lastInput time.Time

		for {
			select {
			case <-time.After(samplingRate):
				if lastInput != (time.Time{}) && lastInput.Add(2*samplingRate).Before(time.Now()) {
					outputChan <- lastInput.Sub(firstInput)

					firstInput = time.Time{}
					lastInput = time.Time{}
				}

			case input := <-inputChan:
				if input != nil {
					outputChan <- input
					return
				}

				if firstInput == (time.Time{}) {
					firstInput = time.Now()
				}
				lastInput = time.Now()
			}
		}
	}()

	return
}

func morseSampler(inputChan chan interface{}, dashMinDuration time.Duration) (outputChan chan interface{}) {
	outputChan = make(chan interface{})

	go func() {
		for input := range inputChan {
			switch input := input.(type) {
			case time.Duration:
				if input < dashMinDuration {
					outputChan <- '.'
				} else {
					outputChan <- '_'
				}

			case error:
				outputChan <- input
				return

			default:
				outputChan <- fmt.Errorf("morseSampler: unsupported type %T", input)
				return
			}
		}
	}()

	return
}

func morseWordSampler(inputChan chan interface{}, maxDuration time.Duration) (outputChan chan interface{}) {
	outputChan = make(chan interface{})

	go func() {
		var tmpWord string
		var lastInput time.Time

		for {
			select {
			case <-time.After(maxDuration):
				if tmpWord != "" && lastInput.Add(maxDuration).Before(time.Now()) {
					outputChan <- tmpWord

					tmpWord = ""
					lastInput = time.Time{}
				}

			case input := <-inputChan:
				switch input := input.(type) {
				case rune:
					tmpWord += string(input)
					lastInput = time.Now()

				case error:
					outputChan <- input
					return

				default:
					outputChan <- fmt.Errorf("morseWordSampler: unsupported type %T", input)
					return
				}
			}
		}
	}()

	return
}
