package main

import (
	"log"
	"time"

	"github.com/tarm/serial"
)

type Signaler struct {
	serialPort   *serial.Port
	samplingRate time.Duration

	closeReaderSyn chan struct{}
	closeReaderAck chan struct{}
	closeWriterSyn chan struct{}
	closeWriterAck chan struct{}
}

func NewSignaler(serialDevice string, samplingRate time.Duration) (signaler *Signaler, err error) {
	signaler = &Signaler{
		samplingRate: samplingRate,

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
				log.Fatalf("Reading errored: %v", err)
			} else if now := time.Now(); lastSend.Add(signaler.samplingRate).Before(now) {
				lastSend = now
				log.Print("ACK")
			}
		}
	}
}

func (signaler *Signaler) backgroundWriter() {
	defer close(signaler.closeWriterAck)

	for {
		select {
		case <-signaler.closeWriterSyn:
			return

		default:
			if _, err := signaler.serialPort.Write([]byte{0xFF}); err != nil {
				log.Fatalf("Writing errored: %v", err)
			}
			if err := signaler.serialPort.Flush(); err != nil {
				log.Fatalf("Flusing errored: %v", err)
			}
		}
	}
}

func (signaler *Signaler) Close() (err error) {
	close(signaler.closeReaderSyn)
	close(signaler.closeWriterSyn)

	for _, closeAck := range []chan struct{}{signaler.closeReaderAck, signaler.closeWriterAck} {
		select {
		case <-closeAck:
		case <-time.After(time.Second):
		}
	}

	return signaler.serialPort.Close()
}
