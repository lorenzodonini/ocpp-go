// The provisioning functional block contains features that help a CSO to provision their Charging Stations, allowing them on their network and retrieving configuration information from these Charging Stations.
// Additionally, it contains features for retrieving information about the configuration of Charging Stations, make changes to the configuration, resetting it etc.
package provisioning

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Provisioning profile.
type CSMSHandler interface {
	// OnBootNotification is called on the CSMS whenever a BootNotificationRequest is received from a charging station.
	OnBootNotification(chargingStationID string, request *BootNotificationRequest) (response *BootNotificationResponse, err error)
	// OnNotifyReport is called on the CSMS whenever a NotifyReportRequest is received from a charging station.
	OnNotifyReport(chargingStationID string, request *NotifyReportRequest) (response *NotifyReportResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Provisioning profile.
type ChargingStationHandler interface {
	// OnGetBaseReport is called on a charging station whenever a GetBaseReportRequest is received from the CSMS.
	OnGetBaseReport(request *GetBaseReportRequest) (response *GetBaseReportResponse, err error)
	// OnGetReport is called on a charging station whenever a GetReportRequest is received from the CSMS.
	OnGetReport(request *GetReportRequest) (response *GetReportResponse, err error)
	// OnGetVariables is called on a charging station whenever a GetVariablesRequest is received from the CSMS.
	OnGetVariables(request *GetVariablesRequest) (response *GetVariablesResponse, err error)
	// OnReset is called on a charging station whenever a ResetRequest is received from the CSMS.
	OnReset(request *ResetRequest) (response *ResetResponse, err error)
	// OnSetNetworkProfile is called on a charging station whenever a SetNetworkProfileRequest is received from the CSMS.
	OnSetNetworkProfile(request *SetNetworkProfileRequest) (response *SetNetworkProfileResponse, err error)
	// OnSetVariables is called on a charging station whenever a SetVariablesRequest is received from the CSMS.
	OnSetVariables(request *SetVariablesRequest) (response *SetVariablesResponse, err error)
}

const ProfileName = "provisioning"

var Profile = ocpp.NewProfile(
	ProfileName,
	BootNotificationFeature{},
	GetBaseReportFeature{},
	GetReportFeature{},
	GetVariablesFeature{},
	NotifyReportFeature{},
	ResetFeature{},
	SetNetworkProfileFeature{},
	SetVariablesFeature{},
)
