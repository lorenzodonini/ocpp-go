// The provisioning functional block contains features that help a CSO to provision their Charging Stations, allowing them on their network and retrieving configuration information from these Charging Stations.
// Additionally, it contains features for retrieving information about the configuration of Charging Stations, make changes to the configuration, resetting it etc.
package provisioning

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Provisioning profile.
type CSMSHandler interface {
	// OnBootNotification is called on the CSMS whenever a BootNotificationRequest is received from a charging station.
	OnBootNotification(chargingStationID string, request *BootNotificationRequest) (confirmation *BootNotificationResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Provisioning profile.
type ChargingStationHandler interface {
	// OnGetBaseReport is called on a charging station whenever a GetBaseReportRequest is received from the CSMS.
	OnGetBaseReport(request *GetBaseReportRequest) (response *GetBaseReportResponse, err error)
	// OnGetReport is called on a charging station whenever a GetReportRequest is received from the CSMS.
	OnGetReport(request *GetReportRequest) (response *GetReportResponse, err error)
}

const ProfileName = "provisioning"

var Profile = ocpp.NewProfile(
	ProfileName,
	BootNotificationFeature{},
	GetBaseReportFeature{},
	GetReportFeature{},
	// SetVariables
)
