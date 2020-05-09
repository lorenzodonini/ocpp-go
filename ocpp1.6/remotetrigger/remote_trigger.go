// Contains support for remote triggering of Charge Point initiated messages.
package remotetrigger

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by Central systems for handling messages part of the OCPP 1.6 RemoteTrigger profile.
type CentralSystemRemoteTriggerHandler interface {
}

// Needs to be implemented by Charge points for handling messages part of the OCPP 1.6 RemoteTrigger profile.
type ChargePointRemoteTriggerHandler interface {
	OnTriggerMessage(request *TriggerMessageRequest) (confirmation *TriggerMessageConfirmation, err error)
}

// The profile name
const ProfileName = "RemoteTrigger"

// Provides support for remote triggering of Charge Point initiated messages.
var Profile = ocpp.NewProfile(
	ProfileName,
	TriggerMessageFeature{})
