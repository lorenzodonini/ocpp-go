package ws

import (
	"errors"
	"testing"
)

func TestWsErrorImplementsError(t *testing.T) {
	err := errors.New("base error")
	wsErr := NewWsError("client123", "ping", err)

	var e error = wsErr
	if e.Error() == "" {
		t.Errorf("WsError.Error() should return a non-empty string")
	}
}

func TestWsErrorTypeAssertion(t *testing.T) {
	err := errors.New("base error")
	wsErr := NewWsError("client456", "write", err)

	ch := make(chan error, 1)
	ch <- wsErr

	recv := <-ch
	wsErr2, ok := recv.(*WsError)
	if !ok {
		t.Errorf("Type assertion to *WsError failed")
	}
	if wsErr2.ClientID != "client456" {
		t.Errorf("Expected ClientID 'client456', got '%s'", wsErr2.ClientID)
	}
	if wsErr2.Operation != "write" {
		t.Errorf("Expected Operation 'write', got '%s'", wsErr2.Operation)
	}
	if wsErr2.Err.Error() != "base error" {
		t.Errorf("Expected base error, got '%s'", wsErr2.Err.Error())
	}
}
