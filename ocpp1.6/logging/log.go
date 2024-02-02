// The diagnostics functional block contains OCPP 2.0 features than enable remote diagnostics of problems with a charging station.
package logging

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 1.6j security extension.
type CentralSystemHandler interface {
	// OnLogStatusNotification is called on the CSMS whenever a LogStatusNotificationRequest is received from a Charging Station.
	OnLogStatusNotification(chargingStationID string, request *LogStatusNotificationRequest) (response *LogStatusNotificationResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 1.6j security extension.
type ChargePointHandler interface {
	// OnGetLog is called on a charging station whenever a GetLogRequest is received from the CSMS.
	OnGetLog(request *GetLogRequest) (response *GetLogResponse, err error)
}

const ProfileName = "Log"

var Profile = ocpp.NewProfile(
	ProfileName,
	GetLogFeature{},
	LogStatusNotificationFeature{},
)
