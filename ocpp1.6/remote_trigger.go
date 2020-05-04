package ocpp16

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	TriggerMessageFeatureName = "TriggerMessage"
)

type CentralSystemRemoteTriggerHandler interface {
}

type ChargePointRemoteTriggerHandler interface {
	OnTriggerMessage(request *TriggerMessageRequest) (confirmation *TriggerMessageConfirmation, err error)
}

const RemoteTriggerProfileName = "RemoteTrigger"

var RemoteTriggerProfile = ocpp.NewProfile(
	RemoteTriggerProfileName,
	TriggerMessageFeature{})
