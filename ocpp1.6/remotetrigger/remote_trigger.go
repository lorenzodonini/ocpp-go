// Contains support for remote triggering of Charge Point initiated messages.
package remotetrigger

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	TriggerMessageFeatureName = "TriggerMessage"
)

type CentralSystemRemoteTriggerHandler interface {
}

type ChargePointRemoteTriggerHandler interface {
	OnTriggerMessage(request *TriggerMessageRequest) (confirmation *TriggerMessageConfirmation, err error)
}

const ProfileName = "RemoteTrigger"

var Profile = ocpp.NewProfile(
	ProfileName,
	TriggerMessageFeature{})
