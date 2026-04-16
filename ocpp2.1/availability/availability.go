// Package availability contains the availability functional block for OCPP 2.1.
// It contains features for notifying the CSMS of availability and status changes.
package availability

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the availability functional block.
const ProfileName = "availability"

// ChargingStationHandler contains callbacks for handling incoming CSMS requests.
type ChargingStationHandler interface {
	// OnChangeAvailability handles a ChangeAvailabilityRequest from the CSMS.
	OnChangeAvailability(request *ChangeAvailabilityRequest) (response *ChangeAvailabilityResponse, err error)
}

// CSMSHandler contains callbacks for handling incoming Charging Station requests.
type CSMSHandler interface {
	// OnHeartbeat handles a HeartbeatRequest from a Charging Station.
	OnHeartbeat(chargingStationID string, request *HeartbeatRequest) (response *HeartbeatResponse, err error)
	// OnStatusNotification handles a StatusNotificationRequest from a Charging Station.
	OnStatusNotification(chargingStationID string, request *StatusNotificationRequest) (response *StatusNotificationResponse, err error)
}

// Profile returns the availability profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		ChangeAvailabilityFeature{},
		HeartbeatFeature{},
		StatusNotificationFeature{},
	)
}
