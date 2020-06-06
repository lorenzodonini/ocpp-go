// Contains features to manage the local authorization list in Charge Points.
package localauth

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by Central systems for handling messages part of the OCPP 1.6 LocalAuthList profile.
type CentralSystemHandler interface {
}

// Needs to be implemented by Charge points for handling messages part of the OCPP 1.6 LocalAuthList profile.
type ChargePointHandler interface {
	OnGetLocalListVersion(request *GetLocalListVersionRequest) (confirmation *GetLocalListVersionConfirmation, err error)
	OnSendLocalList(request *SendLocalListRequest) (confirmation *SendLocalListConfirmation, err error)
}

// The profile name
const ProfileName = "localAuthList"

// Provides support for managing the local authorization list in Charge Points.
var Profile = ocpp.NewProfile(
	ProfileName,
	GetLocalListVersionFeature{},
	SendLocalListFeature{})
