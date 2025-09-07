package ws

import "fmt"

// WsError provides additional context for websocket errors.
type WsError struct {
	ClientID  string // ID of the client (if applicable)
	Operation string // Operation being performed (ping, write, close, etc.)
	Err       error  // Underlying error
}

func (e *WsError) Error() string {
	msg := "websocket error"
	if e.ClientID != "" {
		msg += fmt.Sprintf(" [client: %s]", e.ClientID)
	}
	if e.Operation != "" {
		msg += fmt.Sprintf(" [op: %s]", e.Operation)
	}
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
}

// Helper to create a WsError
func NewWsError(clientID, operation string, err error) *WsError {
	return &WsError{
		ClientID:  clientID,
		Operation: operation,
		Err:       err,
	}
}
