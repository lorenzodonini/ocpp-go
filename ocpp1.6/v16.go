// The package contains an implementation of the OCPP 1.6 communication protocol between a Charge Point and a Central System in an EV charging infrastructure.
package ocpp16

import (
	"crypto/tls"
	"net"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
)

type ChargePointConnection interface {
	ID() string
	RemoteAddr() net.Addr
	TLSConnectionState() *tls.ConnectionState
}

type ChargePointConnectionHandler func(chargePoint ChargePointConnection)

// -------------------- v1.6 Charge Point --------------------

// A Charge Point represents the physical system where an EV can be charged.
// You can instantiate a default Charge Point struct by calling NewClient.
//
// The logic for incoming messages needs to be implemented, and the message handlers need to be registered with the charge point:
//
//	handler := &ChargePointHandler{}
//	client.SetCoreHandler(handler)
//
// Refer to the ChargePointHandler interfaces in the respective core, firmware, localauth, remotetrigger, reservation and smartcharging profiles for the implementation requirements.
//
// A charge point can be started and stopped using the Start and Stop functions.
// While running, messages can be sent to the Central system by calling the Charge point's functions, e.g.
//
//	bootConf, err := client.BootNotification("model1", "vendor1")
//
// All messages are synchronous blocking, and return either the response from the Central system or an error.
// To send asynchronous messages and avoid blocking the calling thread, refer to SendRequestAsync.
type ChargePoint interface {
	// Sends a BootNotificationRequest to the central system, along with information about the charge point.
	BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *core.BootNotificationRequest)) (*core.BootNotificationConfirmation, error)
	// Requests explicit authorization to the central system, provided a valid IdTag (typically the client's). The central system may either authorize or reject the client.
	Authorize(idTag string, props ...func(request *core.AuthorizeRequest)) (*core.AuthorizeConfirmation, error)
	// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
	DataTransfer(vendorId string, props ...func(request *core.DataTransferRequest)) (*core.DataTransferConfirmation, error)
	// Notifies the central system that the charge point is still online. The central system's response is used for time synchronization purposes. It is recommended to perform this operation once every 24 hours.
	Heartbeat(props ...func(request *core.HeartbeatRequest)) (*core.HeartbeatConfirmation, error)
	// Sends a batch of collected meter values to the central system, for billing and analysis. May be done periodically during ongoing transactions.
	MeterValues(connectorId int, meterValues []types.MeterValue, props ...func(request *core.MeterValuesRequest)) (*core.MeterValuesConfirmation, error)
	// Requests to start a transaction for a specific connector. The central system will verify the client's IdTag and either accept or reject the transaction.
	StartTransaction(connectorId int, idTag string, meterStart int, timestamp *types.DateTime, props ...func(request *core.StartTransactionRequest)) (*core.StartTransactionConfirmation, error)
	// Stops an ongoing transaction. Typically a batch of meter values is passed along with this message.
	StopTransaction(meterStop int, timestamp *types.DateTime, transactionId int, props ...func(request *core.StopTransactionRequest)) (*core.StopTransactionConfirmation, error)
	// Notifies the central system of a status update. This may apply to the entire charge point or to a single connector.
	StatusNotification(connectorId int, errorCode core.ChargePointErrorCode, status core.ChargePointStatus, props ...func(request *core.StatusNotificationRequest)) (*core.StatusNotificationConfirmation, error)
	// Notifies the central system of a status change in the upload of diagnostics data.
	DiagnosticsStatusNotification(status firmware.DiagnosticsStatus, props ...func(request *firmware.DiagnosticsStatusNotificationRequest)) (*firmware.DiagnosticsStatusNotificationConfirmation, error)
	// Notifies the central system of a status change during the download of a new firmware version.
	FirmwareStatusNotification(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationConfirmation, error)

	// Registers a handler for incoming core profile messages
	SetCoreHandler(listener core.ChargePointHandler)
	// Registers a handler for incoming local authorization profile messages
	SetLocalAuthListHandler(listener localauth.ChargePointHandler)
	// Registers a handler for incoming firmware management profile messages
	SetFirmwareManagementHandler(listener firmware.ChargePointHandler)
	// Registers a handler for incoming reservation profile messages
	SetReservationHandler(listener reservation.ChargePointHandler)
	// Registers a handler for incoming remote trigger profile messages
	SetRemoteTriggerHandler(listener remotetrigger.ChargePointHandler)
	// Registers a handler for incoming smart charging profile messages
	SetSmartChargingHandler(listener smartcharging.ChargePointHandler)
	// Sends a request to the central system.
	// The central system will respond with a confirmation, or with an error if the request was invalid or could not be processed.
	// In case of network issues (i.e. the remote host couldn't be reached), the function also returns an error.
	//
	// The request is synchronous blocking.
	SendRequest(request ocpp.Request) (ocpp.Response, error)
	// Sends an asynchronous request to the central system.
	// The central system will respond with a confirmation messages, or with an error if the request was invalid or could not be processed.
	// This result is propagated via a callback, called asynchronously.
	// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never called.
	SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Response, protoError error)) error
	// Connects to the central system and starts the charge point routine.
	// The function doesn't block and returns right away, after having attempted to open a connection to the central system.
	// If the connection couldn't be opened, an error is returned.
	//
	// Optional client options must be set before calling this function. Refer to NewClient.
	//
	// No auto-reconnect logic is implemented as of now, but is planned for the future.
	Start(centralSystemUrl string) error
	// Stops the charge point routine, disconnecting it from the central system.
	// Any pending requests are discarded.
	Stop()
	// Returns true if the charge point is currently connected to the central system, false otherwise.
	// While automatically reconnecting to the central system, the method returns false.
	IsConnected() bool
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	// The channel is closed by the charge point when stopped.
	Errors() <-chan error
}

// Creates a new OCPP 1.6 charge point client.
// The id parameter is required to uniquely identify the charge point.
//
// The endpoint and client parameters may be omitted, in order to use a default configuration:
//
//	client := NewClient("someUniqueId", nil, nil)
//
// Additional networking parameters (e.g. TLS or proxy configuration) may be passed, by creating a custom client.
// Here is an example for a client using TLS configuration with a self-signed certificate:
//
//	certPool := x509.NewCertPool()
//	data, err := os.ReadFile("serverSelfSignedCertFilename")
//	if err != nil {
//		log.Fatal(err)
//	}
//	ok = certPool.AppendCertsFromPEM(data)
//	if !ok {
//		log.Fatal("couldn't parse PEM certificate")
//	}
//	cp := NewClient("someUniqueId", nil, ws.NewTLSClient(&tls.Config{
//		RootCAs: certPool,
//	})
//
// For more advanced options, or if a customer networking/occpj layer is required,
// please refer to ocppj.Client and ws.WsClient.
func NewChargePoint(id string, endpoint *ocppj.Client, client ws.WsClient) ChargePoint {
	if client == nil {
		client = ws.NewClient()
	}
	client.SetRequestedSubProtocol(types.V16Subprotocol)
	cp := chargePoint{confirmationHandler: make(chan ocpp.Response, 1), errorHandler: make(chan error, 1), callbacks: callbackqueue.New()}

	if endpoint == nil {
		dispatcher := ocppj.NewDefaultClientDispatcher(ocppj.NewFIFOClientQueue(0))
		endpoint = ocppj.NewClient(id, client, dispatcher, nil, core.Profile, localauth.Profile, firmware.Profile, reservation.Profile, remotetrigger.Profile, smartcharging.Profile)
	}
	// Callback invoked by dispatcher, whenever a queued request is canceled, due to timeout.
	endpoint.SetOnRequestCanceled(cp.onRequestTimeout)
	cp.client = endpoint

	cp.client.SetResponseHandler(func(confirmation ocpp.Response, requestId string) {
		cp.confirmationHandler <- confirmation
	})
	cp.client.SetErrorHandler(func(err *ocpp.Error, details interface{}) {
		cp.errorHandler <- err
	})
	cp.client.SetRequestHandler(cp.handleIncomingRequest)
	return &cp
}

// -------------------- v1.6 Central System --------------------

// A Central System manages Charge Points and has the information for authorizing users for using its Charge Points.
// You can instantiate a default Central System struct by calling the NewServer function.
//
// The logic for handling incoming messages needs to be implemented, and the message handlers need to be registered with the central system:
//
//	handler := &CentralSystemHandler{}
//	server.SetCoreHandler(handler)
//
// Refer to the CentralSystemHandler interfaces in the respective core, firmware, localauth, remotetrigger, reservation and smartcharging profiles for the implementation requirements.
//
// A Central system can be started by using the Start function.
// To be notified of incoming (dis)connections from charge points refer to the SetNewClientHandler and SetChargePointDisconnectedHandler functions.
//
// While running, messages can be sent to a charge point by calling the Central system's functions, e.g.:
//
//	callback := func(conf *ChangeAvailabilityConfirmation, err error) {
//		// handle the response...
//	}
//	changeAvailabilityConf, err := server.ChangeAvailability("cs0001", callback, 1, AvailabilityTypeOperative)
//
// All messages are sent asynchronously and do not block the caller.
type CentralSystem interface {
	// Instructs a charge point to change its availability. The target availability can be set for a single connector of for the whole charge point.
	ChangeAvailability(clientId string, callback func(*core.ChangeAvailabilityConfirmation, error), connectorId int, availabilityType core.AvailabilityType, props ...func(*core.ChangeAvailabilityRequest)) error
	// Changes the configuration of a charge point, by setting a specific key-value pair.
	// The configuration key must be supported by the target charge point, in order for the configuration to be accepted.
	ChangeConfiguration(clientId string, callback func(*core.ChangeConfigurationConfirmation, error), key string, value string, props ...func(*core.ChangeConfigurationRequest)) error
	// Instructs the charge point to clear its current authorization cache. All authorization saved locally will be invalidated.
	ClearCache(clientId string, callback func(*core.ClearCacheConfirmation, error), props ...func(*core.ClearCacheRequest)) error
	// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
	DataTransfer(clientId string, callback func(*core.DataTransferConfirmation, error), vendorId string, props ...func(*core.DataTransferRequest)) error
	// Retrieves the configuration values for the provided configuration keys.
	GetConfiguration(clientId string, callback func(*core.GetConfigurationConfirmation, error), keys []string, props ...func(*core.GetConfigurationRequest)) error
	// Instructs a charge point to start a transaction for a specified client on a provided connector.
	// Depending on the configuration, an explicit authorization message may still be required, before the transaction can start.
	RemoteStartTransaction(clientId string, callback func(*core.RemoteStartTransactionConfirmation, error), idTag string, props ...func(*core.RemoteStartTransactionRequest)) error
	// Instructs a charge point to stop an ongoing transaction, given the transaction's ID.
	RemoteStopTransaction(clientId string, callback func(*core.RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *core.RemoteStopTransactionRequest)) error
	// Forces a charge point to perform an internal hard or soft reset. In both cases, all ongoing transactions are stopped.
	Reset(clientId string, callback func(*core.ResetConfirmation, error), resetType core.ResetType, props ...func(*core.ResetRequest)) error
	// Attempts to unlock a specific connector on a charge point. Used for remote support purposes.
	UnlockConnector(clientId string, callback func(*core.UnlockConnectorConfirmation, error), connectorId int, props ...func(*core.UnlockConnectorRequest)) error
	// Queries the current version of the local authorization list from a charge point.
	GetLocalListVersion(clientId string, callback func(*localauth.GetLocalListVersionConfirmation, error), props ...func(request *localauth.GetLocalListVersionRequest)) error
	// Sends or updates a local authorization list on a charge point. Versioning rules must be followed.
	SendLocalList(clientId string, callback func(*localauth.SendLocalListConfirmation, error), version int, updateType localauth.UpdateType, props ...func(request *localauth.SendLocalListRequest)) error
	// Requests diagnostics data from a charge point. The data will be uploaded out-of-band to the provided URL location.
	GetDiagnostics(clientId string, callback func(*firmware.GetDiagnosticsConfirmation, error), location string, props ...func(request *firmware.GetDiagnosticsRequest)) error
	// Instructs the charge point to download and install a new firmware version. The firmware binary will be downloaded out-of-band from the provided URL location.
	UpdateFirmware(clientId string, callback func(*firmware.UpdateFirmwareConfirmation, error), location string, retrieveDate *types.DateTime, props ...func(request *firmware.UpdateFirmwareRequest)) error
	// Instructs the charge point to reserve a connector for a specific IdTag (client). The connector, or the entire charge point, will be reserved until the provided expiration time.
	ReserveNow(clientId string, callback func(*reservation.ReserveNowConfirmation, error), connectorId int, expiryDate *types.DateTime, idTag string, reservationId int, props ...func(request *reservation.ReserveNowRequest)) error
	// Cancels a previously reserved charge point or connector, given the reservation ID.
	CancelReservation(clientId string, callback func(*reservation.CancelReservationConfirmation, error), reservationId int, props ...func(request *reservation.CancelReservationRequest)) error
	// Instructs a charge point to send a specific message to the central system. This is used for forcefully triggering status updates, when the last known state is either too old or not clear to the central system.
	TriggerMessage(clientId string, callback func(*remotetrigger.TriggerMessageConfirmation, error), requestedMessage remotetrigger.MessageTrigger, props ...func(request *remotetrigger.TriggerMessageRequest)) error
	// Sends a smart charging profile to a charge point. Refer to the smart charging documentation for more information.
	SetChargingProfile(clientId string, callback func(*smartcharging.SetChargingProfileConfirmation, error), connectorId int, chargingProfile *types.ChargingProfile, props ...func(request *smartcharging.SetChargingProfileRequest)) error
	// Removes one or more charging profiles from a charge point.
	ClearChargingProfile(clientId string, callback func(*smartcharging.ClearChargingProfileConfirmation, error), props ...func(request *smartcharging.ClearChargingProfileRequest)) error
	// Queries a charge point to the composite smart charging schedules and rules for a specified time interval.
	GetCompositeSchedule(clientId string, callback func(*smartcharging.GetCompositeScheduleConfirmation, error), connectorId int, duration int, props ...func(request *smartcharging.GetCompositeScheduleRequest)) error

	// Registers a handler for incoming core profile messages.
	SetCoreHandler(handler core.CentralSystemHandler)
	// Registers a handler for incoming local authorization profile messages.
	SetLocalAuthListHandler(handler localauth.CentralSystemHandler)
	// Registers a handler for incoming firmware management profile messages.
	SetFirmwareManagementHandler(handler firmware.CentralSystemHandler)
	// Registers a handler for incoming reservation profile messages.
	SetReservationHandler(handler reservation.CentralSystemHandler)
	// Registers a handler for incoming remote trigger profile messages.
	SetRemoteTriggerHandler(handler remotetrigger.CentralSystemHandler)
	// Registers a handler for incoming smart charging profile messages.
	SetSmartChargingHandler(handler smartcharging.CentralSystemHandler)
	// Registers a handler for new incoming charge point connections.
	SetNewChargePointHandler(handler ChargePointConnectionHandler)
	// Registers a handler for charge point disconnections.
	SetChargePointDisconnectedHandler(handler ChargePointConnectionHandler)
	// Sends an asynchronous request to the charge point.
	// The charge point will respond with a confirmation message, or with an error if the request was invalid or could not be processed.
	// This result is propagated via a callback, called asynchronously.
	// In case of network issues (i.e. the remote host couldn't be reached), the function returns an error directly. In this case, the callback is never called.
	SendRequestAsync(clientId string, request ocpp.Request, callback func(ocpp.Response, error)) error
	// Starts running the central system on the specified port and URL.
	// The central system runs as a daemon and handles incoming charge point connections and messages.
	//
	// The function blocks forever, so it is suggested to wrap it in a goroutine, in case other functionality needs to be executed on the main program thread.
	Start(listenPort int, listenPath string)
	// Errors returns a channel for error messages. If it doesn't exist it es created.
	Errors() <-chan error
}

// Creates a new OCPP 1.6 central system.
//
// The endpoint and server parameters may be omitted, in order to use a default configuration:
//
//	client := NewServer(nil, nil)
//
// It is recommended to use the default configuration, unless a custom networking / ocppj layer is required.
// The default ocppj endpoint supports all OCPP 1.6 profiles out-of-the-box.
//
// If you need a TLS server, you may use the following:
//
//	cs := NewServer(nil, ws.NewTLSServer("certificatePath", "privateKeyPath"))
func NewCentralSystem(endpoint *ocppj.Server, server ws.WsServer) CentralSystem {
	if server == nil {
		server = ws.NewServer()
	}
	server.AddSupportedSubprotocol(types.V16Subprotocol)
	if endpoint == nil {
		endpoint = ocppj.NewServer(server, nil, nil, core.Profile, localauth.Profile, firmware.Profile, reservation.Profile, remotetrigger.Profile, smartcharging.Profile)
	}
	cs := newCentralSystem(endpoint)
	cs.server.SetRequestHandler(func(client ws.Channel, request ocpp.Request, requestId string, action string) {
		cs.handleIncomingRequest(client, request, requestId, action)
	})
	cs.server.SetResponseHandler(func(client ws.Channel, response ocpp.Response, requestId string) {
		cs.handleIncomingConfirmation(client, response, requestId)
	})
	cs.server.SetErrorHandler(func(client ws.Channel, err *ocpp.Error, details interface{}) {
		cs.handleIncomingError(client, err, details)
	})
	cs.server.SetCanceledRequestHandler(func(clientID string, requestID string, request ocpp.Request, err *ocpp.Error) {
		cs.handleCanceledRequest(clientID, request, err)
	})
	return &cs
}
