package ocpp16

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	GetLocalListVersionFeatureName = "GetLocalListVersion"
	SendLocalListFeatureName       = "SendLocalList"
)

type CentralSystemLocalAuthListListener interface {
}

type ChargePointLocalAuthListListener interface {
	OnGetLocalListVersion(request *GetLocalListVersionRequest) (confirmation *GetLocalListVersionConfirmation, err error)
	OnSendLocalList(request *SendLocalListRequest) (confirmation *SendLocalListConfirmation, err error)
}

const LocalAuthListProfileName = "localAuthList"

var LocalAuthListProfile = ocpp.NewProfile(
	LocalAuthListProfileName,
	GetLocalListVersionFeature{},
	SendLocalListFeature{})
