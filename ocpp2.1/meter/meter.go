// Package meter contains the meter functional block for OCPP 2.1.
package meter

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the meter functional block.
const ProfileName = "meter"

// CSMSHandler contains callbacks for handling incoming Charging Station requests.
type CSMSHandler interface {
	// OnMeterValues handles a MeterValuesRequest from the Charging Station.
	OnMeterValues(request *MeterValuesRequest) (response *MeterValuesResponse, err error)
}

// Profile returns the meter profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		MeterValuesFeature{},
	)
}
