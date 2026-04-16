// Package remotecontrol contains the remote control functional block for OCPP 2.1.
package remotecontrol

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the remote control functional block.
const ProfileName = "remotecontrol"

// ChargingStationHandler contains callbacks for handling incoming CSMS requests.
type ChargingStationHandler interface {
	// OnRequestStartTransaction handles a RequestStartTransactionRequest from the CSMS.
	OnRequestStartTransaction(request *RequestStartTransactionRequest) (response *RequestStartTransactionResponse, err error)
	// OnRequestStopTransaction handles a RequestStopTransactionRequest from the CSMS.
	OnRequestStopTransaction(request *RequestStopTransactionRequest) (response *RequestStopTransactionResponse, err error)
	// OnTriggerMessage handles a TriggerMessageRequest from the CSMS.
	OnTriggerMessage(request *TriggerMessageRequest) (response *TriggerMessageResponse, err error)
	// OnUnlockConnector handles an UnlockConnectorRequest from the CSMS.
	OnUnlockConnector(request *UnlockConnectorRequest) (response *UnlockConnectorResponse, err error)
}

// Profile returns the remote control profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		RequestStartTransactionFeature{},
		RequestStopTransactionFeature{},
		TriggerMessageFeature{},
		UnlockConnectorFeature{},
	)
}
