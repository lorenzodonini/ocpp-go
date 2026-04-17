// Package multiplex provides protocol multiplexing support for OCPP servers.
// It allows a single WebSocket server to handle both OCPP 1.6 and OCPP 2.0.1 clients
// on the same port, using WebSocket subprotocol negotiation to route connections
// to the appropriate handler.
package multiplex

import (
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	ocpp2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1"
	"github.com/lorenzodonini/ocpp-go/ws"
)

// ProtocolVersion represents an OCPP protocol version.
type ProtocolVersion string

const (
	// V16 represents OCPP 1.6
	V16 ProtocolVersion = "ocpp1.6"
	// V201 represents OCPP 2.0.1
	V201 ProtocolVersion = "ocpp2.0.1"
)

// Subprotocol constants for WebSocket negotiation.
const (
	V16Subprotocol  = "ocpp1.6"
	V201Subprotocol = "ocpp2.0.1"
)

// SubprotocolSelector is a callback function that allows the application to choose
// which subprotocol to use when a client requests multiple subprotocols.
//
// Parameters:
//   - clientID: The identifier of the connecting client (extracted from the URL path)
//   - requestedSubprotocols: The list of subprotocols requested by the client
//     (from the Sec-WebSocket-Protocol header, e.g., ["ocpp2.0.1", "ocpp1.6"])
//
// Returns the subprotocol to use for this connection (e.g., "ocpp1.6" or "ocpp2.0.1").
// The returned value must be one of the supported subprotocols.
// If an empty string is returned, the default behavior is used (first mutually-supported protocol).
//
// Example:
//
//	server.SetSubprotocolSelector(func(clientID string, requested []string) string {
//	    // Always prefer OCPP 2.0.1 if the client supports it
//	    for _, p := range requested {
//	        if p == multiplex.V201Subprotocol {
//	            return p
//	        }
//	    }
//	    // Fall back to first requested
//	    if len(requested) > 0 {
//	        return requested[0]
//	    }
//	    return ""
//	})
type SubprotocolSelector func(clientID string, requestedSubprotocols []string) string

// MultiProtocolServer is a server that can handle both OCPP 1.6 and OCPP 2.0.1 clients
// on a single port. It uses WebSocket subprotocol negotiation to determine which
// protocol version each client uses.
type MultiProtocolServer interface {
	// Start begins listening for connections on the specified port and path.
	// The function blocks until Stop is called.
	//
	// Example:
	//   go server.Start(8080, "/ocpp/{id}")
	Start(port int, listenPath string)

	// Stop gracefully shuts down the server.
	Stop()

	// OCPP16Server returns the underlying OCPP 1.6 Central System.
	// Use this to register handlers for OCPP 1.6 messages.
	OCPP16Server() ocpp16.CentralSystem

	// OCPP201Server returns the underlying OCPP 2.0.1 CSMS.
	// Use this to register handlers for OCPP 2.0.1 messages.
	OCPP201Server() ocpp2.CSMS

	// SetSubprotocolSelector sets a custom callback for choosing which subprotocol
	// to use when a client requests multiple subprotocols. This allows the application
	// to implement custom protocol selection logic (e.g., prefer OCPP 2.0.1 over 1.6).
	// If not set, the default behavior is used (first mutually-supported protocol).
	SetSubprotocolSelector(selector SubprotocolSelector)

	// SetNewClientHandler sets a callback for all new client connections,
	// regardless of protocol version. The callback receives the WebSocket channel
	// which can be used to determine the protocol via channel.Subprotocol().
	SetNewClientHandler(handler func(channel ws.Channel))

	// SetDisconnectedClientHandler sets a callback for all client disconnections,
	// regardless of protocol version.
	SetDisconnectedClientHandler(handler func(channel ws.Channel))

	// SetBasicAuthHandler enables HTTP Basic Authentication for all connections.
	SetBasicAuthHandler(handler func(username, password string) bool)

	// SetCheckClientHandler sets a validation handler for incoming connections.
	// Return false to reject the connection.
	SetCheckClientHandler(handler ws.CheckClientHandler)
}
