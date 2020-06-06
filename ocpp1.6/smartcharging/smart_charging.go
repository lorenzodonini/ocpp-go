// Contains support for basic Smart Charging, for instance using control pilot.
package smartcharging

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by Central systems for handling messages part of the OCPP 1.6 SmartCharging profile.
type CentralSystemHandler interface {
}

// Needs to be implemented by Charge points for handling messages part of the OCPP 1.6 SmartCharging profile.
type ChargePointHandler interface {
	OnSetChargingProfile(request *SetChargingProfileRequest) (confirmation *SetChargingProfileConfirmation, err error)
	OnClearChargingProfile(request *ClearChargingProfileRequest) (confirmation *ClearChargingProfileConfirmation, err error)
	OnGetCompositeSchedule(request *GetCompositeScheduleRequest) (confirmation *GetCompositeScheduleConfirmation, err error)
}

// The profile name
const ProfileName = "SmartCharging"

// Provides support for basic Smart Charging, for instance using control pilot.
var Profile = ocpp.NewProfile(
	ProfileName,
	SetChargingProfileFeature{},
	ClearChargingProfileFeature{},
	GetCompositeScheduleFeature{})
