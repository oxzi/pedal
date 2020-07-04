package main

import (
	"time"

	"github.com/tarm/serial"
)

// Signaler monitors a serial device and returns a press of the pedal via a channel.
type Signaler struct {
	serialPort   *serial.Port
	samplingRate time.Duration

	signalChan chan interface{}

	closeReaderSyn chan struct{}
	closeReaderAck chan struct{}
	closeWriterSyn chan struct{}
	closeWriterAck chan struct{}
}

// NewSignaler creates a new Signaler for a serial device and a sample rate.
func NewSignaler(serialDevice string, samplingRate time.Duration) (signaler *Signaler, err error) {
	signaler = &Signaler{
		samplingRate: samplingRate,

		signalChan: make(chan interface{}),

		closeReaderSyn: make(chan struct{}),
		closeReaderAck: make(chan struct{}),
		closeWriterSyn: make(chan struct{}),
		closeWriterAck: make(chan struct{}),
	}

	serialConf := &serial.Config{Name: serialDevice, Baud: 110}
	if signaler.serialPort, err = serial.OpenPort(serialConf); err != nil {
		signaler = nil
		return
	}

	go signaler.backgroundReader()
	go signaler.backgroundWriter()

	return
}

// backgroundReader is a goroutine to read from the serial device.
func (signaler *Signaler) backgroundReader() {
	var buf = make([]byte, 128)
	var lastSend time.Time

	defer close(signaler.closeReaderAck)

	for {
		select {
		case <-signaler.closeReaderSyn:
			return

		default:
			if _, err := signaler.serialPort.Read(buf); err != nil {
				signaler.signalChan <- err
				return
			} else if now := time.Now(); lastSend.Add(signaler.samplingRate).Before(now) {
				lastSend = now
				signaler.signalChan <- nil
			}
		}
	}
}

// backgroundWriter is a goroutine to write to the serial device.
func (signaler *Signaler) backgroundWriter() {
	defer close(signaler.closeWriterAck)

	for {
		select {
		case <-signaler.closeWriterSyn:
			return

		default:
			if _, err := signaler.serialPort.Write([]byte{0xFF}); err != nil {
				signaler.signalChan <- err
				return
			} else if err := signaler.serialPort.Flush(); err != nil {
				signaler.signalChan <- err
				return
			}
		}
	}
}

// Chan is the feedback channel.
//
// Each press within the sample rate emits a nil. However, in case of an error, this error is sent.
func (signaler *Signaler) Chan() chan interface{} {
	return signaler.signalChan
}

// Close this Signaler with its underlying serial connection.
func (signaler *Signaler) Close() (err error) {
	close(signaler.closeReaderSyn)
	close(signaler.closeWriterSyn)

	for _, closeAck := range []chan struct{}{signaler.closeReaderAck, signaler.closeWriterAck} {
		select {
		case <-closeAck:
		case <-time.After(time.Second):
		}
	}

	close(signaler.signalChan)

	return signaler.serialPort.Close()
}
