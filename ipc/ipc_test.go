package ipc

import (
	"io/ioutil"
	"testing"
	"time"
)

func randomFile(t *testing.T) string {
	if dir, err := ioutil.TempDir("", "*"); err != nil {
		t.Fatal(err)
		return ""
	} else {
		return dir + "/sock"
	}
}

func TestIpcConnect(t *testing.T) {
	inputFuncChan := make(chan struct{})
	inputPayload := "hello world"

	inputFunc := func(payload string) {
		if payload == inputPayload {
			close(inputFuncChan)
		} else {
			t.Fatalf("Invalid payload: %s", payload)
		}
	}

	modeFuncChan := make(chan struct{})
	modePayload := "foo bar buz"

	modeFunc := func(payload string) {
		if payload == modePayload {
			close(modeFuncChan)
		} else {
			t.Fatalf("Invalid payload: %s", payload)
		}
	}

	socket := randomFile(t)

	server, serverErr := NewServer(socket, inputFunc, modeFunc)
	if serverErr != nil {
		t.Fatal(serverErr)
	}

	if err := SendMessage(socket, Message{Kind: Input, Payload: inputPayload}); err != nil {
		t.Fatal(err)
	}

	if err := SendMessage(socket, Message{Kind: Mode, Payload: modePayload}); err != nil {
		t.Fatal(err)
	}

	select {
	case <-inputFuncChan:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("No feedback from inputFunc")
	}

	select {
	case <-modeFuncChan:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("No feedback from modeFuncChan")
	}

	if err := server.Close(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)
}
