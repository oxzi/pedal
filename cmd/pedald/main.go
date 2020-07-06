package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/geistesk/pedal/ipc"
	"github.com/geistesk/pedal/modes"
	"github.com/geistesk/pedal/pedal"
)

const serverSocket = "/tmp/pedal.sock"

const samplingRate = 50 * time.Millisecond
const morseMaxUnit = 250 * time.Millisecond

var (
	server *ipc.Server
	mutex  sync.Mutex

	signaler *pedal.Signaler
	mode     modes.Mode
)

// signalerClose checks if signaler is set, and closes it.
func signalerClose() {
	if signaler != nil {
		log.Info("Closing old Signaler, this might take some seconds")
		if err := signaler.Close(); err != nil {
			log.WithError(err).Error("Closing Signaler errored")
		}

		signaler = nil
	}
}

// modeClose checks if mode is set, and closes it.
func modeClose() {
	if mode != nil {
		log.Info("Closing old Mode")
		if err := mode.Close(); err != nil {
			log.WithError(err).Error("Closing Mode errored")
		}

		mode = nil
	}
}

// signalerCallback is called for IPC signalerCallbacks.
func signalerCallback(tty string) {
	mutex.Lock()
	defer mutex.Unlock()

	modeClose()
	signalerClose()

	if s, err := pedal.NewSignaler(tty, samplingRate); err != nil {
		log.WithError(err).Error("Updating Signaler errored")
	} else {
		signaler = s
		log.Info("Updated Signaler; please now configure a Mode")
	}
}

// modeCallback is called fo rIPC modeCallbacks.
func modeCallback(payload string) {
	mutex.Lock()
	defer mutex.Unlock()

	if signaler == nil {
		log.Warn("A Signaler must be set before configuring a Mode")
		return
	}

	modeClose()

	// TODO: cases
	micMuteAction := modes.NewCommandAction("amixer -c 0 set Capture toggle")
	mode = modes.NewTrigger(signaler.Chan(), micMuteAction, 500*time.Millisecond)
	log.Info("Updated Mode")
}

// waitInterrupt waits for a SIGINT.
func waitInterrupt() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func main() {
	log.Info("Starting up..")

	if s, err := ipc.NewServer(serverSocket, signalerCallback, modeCallback); err != nil {
		log.WithError(err).Fatal("Starting server failed")
	} else {
		server = s
	}

	waitInterrupt()
	log.Info("Closing down..")

	if err := server.Close(); err != nil {
		log.WithError(err).Error("Closing server failed")
	}

	mutex.Lock()
	modeClose()
	signalerClose()
	mutex.Unlock()
}
