// Package multiplex provides protocol multiplexing support for OCPP servers.
// It allows a single WebSocket server to handle multiple OCPP protocol versions
// on the same port, using WebSocket subprotocol negotiation to route connections
// to the appropriate handler.
package multiplex

import "github.com/lorenzodonini/ocpp-go/ws"

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

// Server provides a shared WebSocket server for multiple OCPP protocol versions.
// Create OCPP servers (CentralSystem, CSMS) using the shared WebSocketServer(),
// then call Start() on each one. The ws.Server.Start() method is idempotent,
// so only the first call actually starts the server.
//
// Example usage:
//
//	// Create multiplex server
//	mux := multiplex.NewServer()
//
//	// Create only the OCPP servers you need
//	cs := ocpp16.NewCentralSystem(nil, mux.WebSocketServer())
//	csms := ocpp2.NewCSMS(nil, mux.WebSocketServer())
//
//	// Register handlers
//	cs.SetCoreHandler(&myOCPP16Handler{})
//	csms.SetProvisioningHandler(&myOCPP201Handler{})
//
//	// Start servers (ws.Server.Start is idempotent)
//	go cs.Start(8080, "/ocpp/{id}")
//	csms.Start(8080, "/ocpp/{id}")
type Server interface {
	// WebSocketServer returns the underlying ws.Server.
	// Use this when creating OCPP servers to share the same WebSocket connection.
	WebSocketServer() ws.Server

	// SetSubprotocolSelector sets a custom callback for choosing which subprotocol
	// to use when a client requests multiple subprotocols.
	SetSubprotocolSelector(selector SubprotocolSelector)

	// SetBasicAuthHandler enables HTTP Basic Authentication for all connections.
	SetBasicAuthHandler(handler func(username, password string) bool)

	// SetCheckClientHandler sets a validation handler for incoming connections.
	SetCheckClientHandler(handler ws.CheckClientHandler)
}
