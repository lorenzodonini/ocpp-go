// The authorization functional block contains OCPP 2.1 features that show tariff and costs to an EV driver, when supported by the charging station.
package tariffcost

import "github.com/xBlaz3kx/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 Tariff and cost profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.1 Tariff and cost profile.
type ChargingStationHandler interface {
	// OnCostUpdated is called on a charging station whenever a CostUpdatedRequest is received from the CSMS.
	OnCostUpdated(request *CostUpdatedRequest) (confirmation *CostUpdatedResponse, err error)
	// OnSetDefaultTariff is called on a charging station whenever a SetDefaultTariffRequest is received from the CSMS.
	OnSetDefaultTariff(request *SetDefaultTariffRequest) (confirmation *SetDefaultTariffResponse, err error)
	// OnGetTariffs is called on a charging station whenever a GetTariffsRequest is received from the CSMS.
	OnGetTariffs(request *GetTariffsRequest) (confirmation *GetTariffsResponse, err error)
	// OnClearTariffs is called on a charging station whenever a ClearTariffsResponse is received from the CSMS.
	OnClearTariffs(request *ClearTariffsRequest) (confirmation *ClearTariffsResponse, err error)
	// OnChangeTransactionTariff is called on a charging station whenever a ChangeTransactionTariffRequest is received from the CSMS.
	OnChangeTransactionTariff(request *ChangeTransactionTariffRequest) (confirmation *ChangeTransactionTariffResponse, err error)
}

const ProfileName = "TariffCost"

var Profile = ocpp.NewProfile(
	ProfileName,
	CostUpdatedFeature{},
	SetDefaultTariffFeature{},
	GetTariffsFeature{},
	ClearTariffsFeature{},
	ChangeTransactionTariffFeature{},
)
