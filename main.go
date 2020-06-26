package main

import (
	"fmt"
	"time"
)

func main() {
	signaler, err := NewSignaler("/dev/ttyUSB0", 50*time.Millisecond)
	if err != nil {
		panic(err)
	}

	go func() {
		for signal := range signaler.Chan() {
			if signal == nil {
				fmt.Println("ACK")
			} else {
				panic(signal)
			}
		}
	}()

	time.Sleep(10 * time.Second)

	if err := signaler.Close(); err != nil {
		panic(err)
	}
}
