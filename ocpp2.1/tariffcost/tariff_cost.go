// The authorization functional block contains OCPP 2.1 features that show tariff and costs to an EV driver, when supported by the charging station.
package tariffcost

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 Tariff and cost profile.
type CSMSHandler interface {
	// OnNotifySettlement is called on the CSMS whenever a NotifySettlementRequest is received from a charging station.
	OnNotifySettlement(chargingStationID string, request *NotifySettlementRequest) (response *NotifySettlementResponse, err error)
	// OnNotifyWebPaymentStarted is called on the CSMS whenever a NotifyWebPaymentStartedRequest is received from a charging station.
	OnNotifyWebPaymentStarted(chargingStationID string, request *NotifyWebPaymentStartedRequest) (response *NotifyWebPaymentStartedResponse, err error)
	// OnVatNumberValidation is called on the CSMS whenever a VatNumberValidationRequest is received from a charging station.
	OnVatNumberValidation(chargingStationID string, request *VatNumberValidationRequest) (response *VatNumberValidationResponse, err error)
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
	ChangeTransactionTariffFeature{},
	ClearTariffsFeature{},
	CostUpdatedFeature{},
	GetTariffsFeature{},
	NotifySettlementFeature{},
	NotifyWebPaymentStartedFeature{},
	SetDefaultTariffFeature{},
	VatNumberValidationFeature{},
)
