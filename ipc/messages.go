package ipc

import (
	"encoding/gob"
)

// Target defines a Message's kind.
type Target int

const (
	// Input defines the Signaler's input, e.g., /dev/ttyUSB0.
	Input Target = iota

	// Mode defines a Mode, e.g., modes.Trigger.
	Mode Target = iota
)

func (target Target) String() string {
	switch target {
	case Input:
		return "Input"
	case Mode:
		return "Mode"
	default:
		return "unknown"
	}
}

// Message to be sent from a Client (SendMessage) to a Server.
type Message struct {
	Kind    Target
	Payload string
}

func init() {
	gob.Register(Message{})
}
