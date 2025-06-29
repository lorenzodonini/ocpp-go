// The Smart charging functional block contains OCPP 2.0 features that enable the CSO (or a third party) to influence the charging current/power transferred during a transaction, or set limits to the amount of current/power a Charging Station can draw from the grid.
package smartcharging

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Smart charging profile.
type CSMSHandler interface {
	// OnClearedChargingLimit is called on the CSMS whenever a ClearedChargingLimitRequest is received from a charging station.
	OnClearedChargingLimit(chargingStationID string, request *ClearedChargingLimitRequest) (response *ClearedChargingLimitResponse, err error)
	// OnNotifyChargingLimit is called on the CSMS whenever a NotifyChargingLimitRequest is received from a charging station.
	OnNotifyChargingLimit(chargingStationID string, request *NotifyChargingLimitRequest) (response *NotifyChargingLimitResponse, err error)
	// OnNotifyEVChargingNeeds is called on the CSMS whenever a NotifyEVChargingNeedsRequest is received from a charging station.
	OnNotifyEVChargingNeeds(chargingStationID string, request *NotifyEVChargingNeedsRequest) (response *NotifyEVChargingNeedsResponse, err error)
	// OnNotifyEVChargingSchedule is called on the CSMS whenever a NotifyEVChargingScheduleRequest is received from a charging station.
	OnNotifyEVChargingSchedule(chargingStationID string, request *NotifyEVChargingScheduleRequest) (response *NotifyEVChargingScheduleResponse, err error)
	// OnReportChargingProfiles is called on the CSMS whenever a ReportChargingProfilesRequest is received from a charging station.
	OnReportChargingProfiles(chargingStationID string, request *ReportChargingProfilesRequest) (reponse *ReportChargingProfilesResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Smart charging profile.
type ChargingStationHandler interface {
	// OnClearChargingProfile is called on a charging station whenever a ClearChargingProfileRequest is received from the CSMS.
	OnClearChargingProfile(request *ClearChargingProfileRequest) (response *ClearChargingProfileResponse, err error)
	// OnGetChargingProfiles is called on a charging station whenever a GetChargingProfilesRequest is received from the CSMS.
	OnGetChargingProfiles(request *GetChargingProfilesRequest) (response *GetChargingProfilesResponse, err error)
	// OnGetCompositeSchedule is called on a charging station whenever a GetCompositeScheduleRequest is received from the CSMS.
	OnGetCompositeSchedule(request *GetCompositeScheduleRequest) (response *GetCompositeScheduleResponse, err error)
	// OnSetChargingProfile is called on a charging station whenever a SetChargingProfileRequest is received from the CSMS.
	OnSetChargingProfile(request *SetChargingProfileRequest) (response *SetChargingProfileResponse, err error)
}

const ProfileName = "SmartCharging"

var Profile = ocpp.NewProfile(
	ProfileName,
	ClearChargingProfileFeature{},
	ClearedChargingLimitFeature{},
	GetChargingProfilesFeature{},
	GetCompositeScheduleFeature{},
	NotifyChargingLimitFeature{},
	NotifyEVChargingNeedsFeature{},
	NotifyEVChargingScheduleFeature{},
	ReportChargingProfilesFeature{},
	SetChargingProfileFeature{},
)
