// The authorization functional block contains OCPP 2.0 features that show tariff and costs to an EV driver, when supported by the charging station.
package tariffcost

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Tariff and cost profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Tariff and cost profile.
type ChargingStationHandler interface {
	// OnCostUpdated is called on a charging station whenever a CostUpdatedRequest is received from the CSMS.
	OnCostUpdated(request *CostUpdatedRequest) (confirmation *CostUpdatedResponse, err error)
}

const ProfileName = "TariffCost"

var Profile = ocpp.NewProfile(
	ProfileName,
	CostUpdatedFeature{},
)
