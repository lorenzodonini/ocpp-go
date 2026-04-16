// Package provisioning contains the provisioning functional block for OCPP 2.1.
// The provisioning functional block contains features that help a CSO provision
// their Charging Stations, allowing them on their network and retrieving basic
// information about the capabilities and configuration of a Charging Station.
package provisioning

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the provisioning functional block.
const ProfileName = "provisioning"

// ChargingStationHandler contains callbacks for handling incoming CSMS requests
// related to the provisioning functional block.
type ChargingStationHandler interface {
	// OnGetBaseReport handles a GetBaseReportRequest from the CSMS.
	OnGetBaseReport(request *GetBaseReportRequest) (response *GetBaseReportResponse, err error)
	// OnGetReport handles a GetReportRequest from the CSMS.
	OnGetReport(request *GetReportRequest) (response *GetReportResponse, err error)
	// OnGetVariables handles a GetVariablesRequest from the CSMS.
	OnGetVariables(request *GetVariablesRequest) (response *GetVariablesResponse, err error)
	// OnSetVariables handles a SetVariablesRequest from the CSMS.
	OnSetVariables(request *SetVariablesRequest) (response *SetVariablesResponse, err error)
	// OnSetNetworkProfile handles a SetNetworkProfileRequest from the CSMS.
	OnSetNetworkProfile(request *SetNetworkProfileRequest) (response *SetNetworkProfileResponse, err error)
	// OnReset handles a ResetRequest from the CSMS.
	OnReset(request *ResetRequest) (response *ResetResponse, err error)
}

// CSMSHandler contains callbacks for handling incoming Charging Station requests
// related to the provisioning functional block.
type CSMSHandler interface {
	// OnBootNotification handles a BootNotificationRequest from a Charging Station.
	OnBootNotification(chargingStationID string, request *BootNotificationRequest) (response *BootNotificationResponse, err error)
	// OnNotifyReport handles a NotifyReportRequest from a Charging Station.
	OnNotifyReport(chargingStationID string, request *NotifyReportRequest) (response *NotifyReportResponse, err error)
}

// Profile returns the provisioning profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		BootNotificationFeature{},
		ResetFeature{},
		GetBaseReportFeature{},
		GetReportFeature{},
		NotifyReportFeature{},
		GetVariablesFeature{},
		SetVariablesFeature{},
		SetNetworkProfileFeature{},
	)
}
