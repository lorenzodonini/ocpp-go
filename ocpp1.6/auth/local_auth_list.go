// Contains features to manage the local authorization list in Charge Points.
package auth

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	GetLocalListVersionFeatureName = "GetLocalListVersion"
	SendLocalListFeatureName       = "SendLocalList"
)

type CentralSystemLocalAuthListHandler interface {
}

type ChargePointLocalAuthListHandler interface {
	OnGetLocalListVersion(request *GetLocalListVersionRequest) (confirmation *GetLocalListVersionConfirmation, err error)
	OnSendLocalList(request *SendLocalListRequest) (confirmation *SendLocalListConfirmation, err error)
}

const ProfileName = "localAuthList"

var Profile = ocpp.NewProfile(
	ProfileName,
	GetLocalListVersionFeature{},
	SendLocalListFeature{})
