package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/oxzi/pedal/ipc"
)

const serverSocket = "/tmp/pedal.sock"

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Expecting two arguments or more.")
	}

	kind := os.Args[1]
	args := os.Args[2:]

	message := ipc.Message{}

	switch kind {
	case "input":
		message.Kind = ipc.Input

	case "mode":
		message.Kind = ipc.Mode

	default:
		log.WithField("Argument", kind).Fatal("Unknown kind")
	}

	message.Payload = strings.Join(args, " ")

	if err := ipc.SendMessage(serverSocket, message); err != nil {
		log.WithError(err).Fatal("Sending Message failed.")
	}
}
