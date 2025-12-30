// The Remote control functional block contains OCPP 2.0 features for remote-control management from the CSMS.
package remotecontrol

import "github.com/xBlaz3kx/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Remote control profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Remote control profile.
type ChargingStationHandler interface {
	// OnRequestStartTransaction is called on a charging station whenever a RequestStartTransactionRequest is received from the CSMS.
	OnRequestStartTransaction(request *RequestStartTransactionRequest) (response *RequestStartTransactionResponse, err error)
	// OnRequestStopTransaction is called on a charging station whenever a RequestStopTransactionRequest is received from the CSMS.
	OnRequestStopTransaction(request *RequestStopTransactionRequest) (response *RequestStopTransactionResponse, err error)
	// OnTriggerMessage is called on a charging station whenever a TriggerMessageRequest is received from the CSMS.
	OnTriggerMessage(request *TriggerMessageRequest) (response *TriggerMessageResponse, err error)
	// OnUnlockConnector is called on a charging station whenever a UnlockConnectorRequest is received from the CSMS.
	OnUnlockConnector(request *UnlockConnectorRequest) (response *UnlockConnectorResponse, err error)
}

const ProfileName = "RemoteControl"

var Profile = ocpp.NewProfile(
	ProfileName,
	RequestStartTransactionFeature{},
	RequestStopTransactionFeature{},
	TriggerMessageFeature{},
	UnlockConnectorFeature{},
)
