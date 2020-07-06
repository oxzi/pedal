package ipc

import (
	"encoding/gob"
	"net"
)

// SendMessage to a Server on socketname.
func SendMessage(socketname string, message Message) error {
	conn, connErr := net.Dial("unix", socketname)
	if connErr != nil {
		return connErr
	}

	if err := gob.NewEncoder(conn).Encode(message); err != nil {
		return err
	}

	return conn.Close()
}
