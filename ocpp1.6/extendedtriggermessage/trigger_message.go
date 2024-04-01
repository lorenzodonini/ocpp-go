// The diagnostics functional block contains OCPP 2.0 features than enable remote diagnostics of problems with a charging station.
package extendedtriggermessage

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by Charging stations for handling messages part of the OCPP 1.6j security extension.
type ChargePointHandler interface {
	// OnExtendedTriggerMessage is called on a charging station whenever a ExtendedTriggerMessageRequest is received from the CSMS.
	OnExtendedTriggerMessage(request *ExtendedTriggerMessageRequest) (response *ExtendedTriggerMessageResponse, err error)
}

const ProfileName = "ExtendedTriggerMessage"

var Profile = ocpp.NewProfile(
	ProfileName,
	ExtendedTriggerMessageFeature{},
)
