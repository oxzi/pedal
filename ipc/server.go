package ipc

import (
	"encoding/gob"
	"io"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Server listens on an Unix socket for Message. Those are delegated to a callback function.
type Server struct {
	listener net.Listener

	signalerCallback func(string)
	modeCallback     func(string)
}

// NewServer to be started on the given socketname with the two callback functions.
func NewServer(socketname string, signalerCallback func(string), modeCallback func(string)) (server *Server, err error) {
	defer func() {
		if err != nil {
			server = nil
		}
	}()

	server = &Server{
		signalerCallback: signalerCallback,
		modeCallback:     modeCallback,
	}

	if _, fileInfoErr := os.Stat(socketname); !os.IsNotExist(fileInfoErr) {
		if err = os.Remove(socketname); err != nil {
			return
		}
	}

	if server.listener, err = net.Listen("unix", socketname); err != nil {
		return
	}

	go server.serve()

	return
}

// Close this Server and its active connections.
func (server *Server) Close() error {
	return server.listener.Close()
}

// serve incoming connections and dispatch those to serveConn.
func (server *Server) serve() {
	for {
		if conn, err := server.listener.Accept(); err != nil {
			const closeErr = "use of closed network connection"
			if strings.Contains(err.Error(), closeErr) {
				return
			}

			log.WithError(err).Error("Accepting socket connection errored")
		} else {
			go server.serveConn(conn)
		}
	}
}

// serveConn handles incoming sessions.
func (server *Server) serveConn(conn net.Conn) {
	defer func() {
		log.WithField("Conn", conn).Debug("Closing socket connection")
		_ = conn.Close()
	}()

	log.WithField("Conn", conn).Debug("Starting new socket connection")

	var message Message
	var decoder = gob.NewDecoder(conn)

	for {
		if err := decoder.Decode(&message); err != nil {
			if err != io.EOF {
				log.WithError(err).Error("Parsing message from socket errored")
			}

			return
		}

		logger := log.WithField("Message", message)
		switch message.Kind {
		case Input:
			logger.Info("Processing input change")
			server.signalerCallback(message.Payload)

		case Mode:
			logger.Info("Processing modechange")
			server.modeCallback(message.Payload)

		default:
			logger.Warn("Unknown kind in received message.")
			return
		}
	}
}
