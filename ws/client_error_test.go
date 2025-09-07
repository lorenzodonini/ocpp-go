package ws

import (
	"errors"
	"testing"
)

func TestClientErrorChannelCompatibility(t *testing.T) {
	c := &client{
		errC: make(chan error, 1),
	}
	baseErr := errors.New("something went wrong")
	wsErr := NewWsError("test-client", "write", baseErr)

	// Send WsError to errC
	c.errC <- wsErr

	recv := <-c.errC
	// Should be accessible as error
	if recv.Error() == "" {
		t.Errorf("Error message should not be empty")
	}
	// Should be accessible as *WsError
	typed, ok := recv.(*WsError)
	if !ok {
		t.Errorf("Type assertion to *WsError failed")
	}
	if typed.ClientID != "test-client" {
		t.Errorf("Expected ClientID 'test-client', got '%s'", typed.ClientID)
	}
	if typed.Operation != "write" {
		t.Errorf("Expected Operation 'write', got '%s'", typed.Operation)
	}
	if typed.Err.Error() != "something went wrong" {
		t.Errorf("Expected base error, got '%s'", typed.Err.Error())
	}
}
