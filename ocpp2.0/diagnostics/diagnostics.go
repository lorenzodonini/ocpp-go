// The diagnostics functional block contains OCPP 2.0 features than enable remote diagnostics of problems with a charging station.
package diagnostics

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Diagnostics profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Diagnostics profile.
type ChargingStationHandler interface {
	// OnClearVariableMonitoring is called on a charging station whenever a ClearVariableMonitoringRequest is received from the CSMS.
	OnClearVariableMonitoring(request *ClearVariableMonitoringRequest) (confirmation *ClearVariableMonitoringConfirmation, err error)
	// OnCustomerInformation is called on a charging station whenever a CustomerInformationRequest is received from the CSMS.
	OnCustomerInformation(request *CustomerInformationRequest) (confirmation *CustomerInformationConfirmation, err error)
}

const ProfileName = "diagnostics"

var Profile = ocpp.NewProfile(
	ProfileName,
	ClearVariableMonitoringFeature{},
	CustomerInformationFeature{},
	)
