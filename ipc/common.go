package ipc

import "os"

// DefaultServerSocket is the default Unix socket; overwritten by ServerSocketEnv.
const DefaultServerSocket = "/tmp/pedal.sock"

// ServerSocketEnv is the environment variable's key to override the DefaultServerSocket.
const ServerSocketEnv = "SOCKET"

// ServerSocketPath is the file path of the Unix socket used for IPC.
//
// By default, DefaultServerSocket is used, which might be overwritten by the ServerSocketPath environment variable.
func ServerSocketPath() (socket string) {
	if socket = os.Getenv(ServerSocketEnv); socket == "" {
		socket = DefaultServerSocket
	}
	return
}