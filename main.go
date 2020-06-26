package main

import (
	"time"
)

func main() {
	signaler, err := NewSignaler("/dev/ttyUSB0", 50*time.Millisecond)
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
