package ocpp16

import "github.com/lorenzodonini/ocpp-go/ocpp"

const (
	SetChargingProfileFeatureName   = "SetChargingProfile"
	ClearChargingProfileFeatureName = "ClearChargingProfile"
	GetCompositeScheduleFeatureName = "GetCompositeSchedule"
)

type CentralSystemSmartChargingHandler interface {
}

type ChargePointSmartChargingHandler interface {
	OnSetChargingProfile(request *SetChargingProfileRequest) (confirmation *SetChargingProfileConfirmation, err error)
	OnClearChargingProfile(request *ClearChargingProfileRequest) (confirmation *ClearChargingProfileConfirmation, err error)
	OnGetCompositeSchedule(request *GetCompositeScheduleRequest) (confirmation *GetCompositeScheduleConfirmation, err error)
}

const SmartChargingProfileName = "SmartCharging"

var SmartChargingProfile = ocpp.NewProfile(
	SmartChargingProfileName,
	SetChargingProfileFeature{},
	ClearChargingProfileFeature{},
	GetCompositeScheduleFeature{})
