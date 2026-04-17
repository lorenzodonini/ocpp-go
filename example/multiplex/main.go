// Example demonstrating protocol multiplexing for OCPP 1.6 and 2.0.1 on a single port.
//
// This server can handle both OCPP 1.6 charge points and OCPP 2.0.1 charging stations
// simultaneously using WebSocket subprotocol negotiation.
//
// When a client connects with "Sec-WebSocket-Protocol: ocpp2.0.1, ocpp1.6",
// the server negotiates the first mutually-supported protocol.
package main

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/multiplex"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/ws"
)

// OCPP 1.6 handler
type ocpp16Handler struct{}

func (h *ocpp16Handler) OnAuthorize(chargePointId string, request *core.AuthorizeRequest) (*core.AuthorizeConfirmation, error) {
	fmt.Printf("[OCPP 1.6] Authorize from %s: %s\n", chargePointId, request.IdTag)
	return core.NewAuthorizationConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted)), nil
}

func (h *ocpp16Handler) OnBootNotification(chargePointId string, request *core.BootNotificationRequest) (*core.BootNotificationConfirmation, error) {
	fmt.Printf("[OCPP 1.6] BootNotification from %s: %s %s\n", chargePointId, request.ChargePointVendor, request.ChargePointModel)
	return core.NewBootNotificationConfirmation(types.NewDateTime(time.Now()), 60, core.RegistrationStatusAccepted), nil
}

func (h *ocpp16Handler) OnDataTransfer(chargePointId string, request *core.DataTransferRequest) (*core.DataTransferConfirmation, error) {
	fmt.Printf("[OCPP 1.6] DataTransfer from %s\n", chargePointId)
	return core.NewDataTransferConfirmation(core.DataTransferStatusAccepted), nil
}

func (h *ocpp16Handler) OnHeartbeat(chargePointId string, request *core.HeartbeatRequest) (*core.HeartbeatConfirmation, error) {
	fmt.Printf("[OCPP 1.6] Heartbeat from %s\n", chargePointId)
	return core.NewHeartbeatConfirmation(types.NewDateTime(time.Now())), nil
}

func (h *ocpp16Handler) OnMeterValues(chargePointId string, request *core.MeterValuesRequest) (*core.MeterValuesConfirmation, error) {
	fmt.Printf("[OCPP 1.6] MeterValues from %s\n", chargePointId)
	return core.NewMeterValuesConfirmation(), nil
}

func (h *ocpp16Handler) OnStartTransaction(chargePointId string, request *core.StartTransactionRequest) (*core.StartTransactionConfirmation, error) {
	fmt.Printf("[OCPP 1.6] StartTransaction from %s\n", chargePointId)
	return core.NewStartTransactionConfirmation(types.NewIdTagInfo(types.AuthorizationStatusAccepted), 1), nil
}

func (h *ocpp16Handler) OnStatusNotification(chargePointId string, request *core.StatusNotificationRequest) (*core.StatusNotificationConfirmation, error) {
	fmt.Printf("[OCPP 1.6] StatusNotification from %s: connector %d is %s\n", chargePointId, request.ConnectorId, request.Status)
	return core.NewStatusNotificationConfirmation(), nil
}

func (h *ocpp16Handler) OnStopTransaction(chargePointId string, request *core.StopTransactionRequest) (*core.StopTransactionConfirmation, error) {
	fmt.Printf("[OCPP 1.6] StopTransaction from %s\n", chargePointId)
	return core.NewStopTransactionConfirmation(), nil
}

// OCPP 2.0.1 handler
type ocpp201Handler struct{}

func (h *ocpp201Handler) OnBootNotification(chargingStationId string, request *provisioning.BootNotificationRequest) (*provisioning.BootNotificationResponse, error) {
	fmt.Printf("[OCPP 2.0.1] BootNotification from %s: %s %s\n", chargingStationId, request.ChargingStation.VendorName, request.ChargingStation.Model)
	return provisioning.NewBootNotificationResponse(types2.NewDateTime(time.Now()), 60, provisioning.RegistrationStatusAccepted), nil
}

func (h *ocpp201Handler) OnNotifyReport(chargingStationId string, request *provisioning.NotifyReportRequest) (*provisioning.NotifyReportResponse, error) {
	fmt.Printf("[OCPP 2.0.1] NotifyReport from %s\n", chargingStationId)
	return provisioning.NewNotifyReportResponse(), nil
}

func main() {
	// Create multi-protocol server
	server := multiplex.NewMultiProtocolServer()

	// Register OCPP 1.6 handlers
	server.OCPP16Server().SetCoreHandler(&ocpp16Handler{})

	// Register OCPP 2.0.1 handlers
	server.OCPP201Server().SetProvisioningHandler(&ocpp201Handler{})

	// Track connections and their protocol versions
	server.SetNewClientHandler(func(channel ws.Channel) {
		fmt.Printf("New client connected: %s (protocol: %s)\n", channel.ID(), channel.Subprotocol())
	})

	server.SetDisconnectedClientHandler(func(channel ws.Channel) {
		fmt.Printf("Client disconnected: %s\n", channel.ID())
	})

	// Start listening on port 8080
	fmt.Println("Starting multi-protocol OCPP server on :8080/ocpp/{id}")
	fmt.Println("Supported protocols: ocpp1.6, ocpp2.0.1")
	fmt.Println()
	fmt.Println("Clients can connect with either protocol. The server will negotiate")
	fmt.Println("based on the Sec-WebSocket-Protocol header.")

	server.Start(8080, "/ocpp/{id}")
}
