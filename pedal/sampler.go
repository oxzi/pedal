package pedal

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
