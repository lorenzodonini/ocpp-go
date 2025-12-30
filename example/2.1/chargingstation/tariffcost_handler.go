package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/tariffcost"
)

func (handler *ChargingStationHandler) OnCostUpdated(request *tariffcost.CostUpdatedRequest) (response *tariffcost.CostUpdatedResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("accepted request to display cost for transaction %v: %v", request.TransactionID, request.TotalCost)
	// TODO: update internal display to show updated cost for transaction
	return tariffcost.NewCostUpdatedResponse(), nil
}

func (handler *ChargingStationHandler) OnSetDefaultTariff(request *tariffcost.SetDefaultTariffRequest) (confirmation *tariffcost.SetDefaultTariffResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("accepted request to set default tariff for EVSE %d: %s", request.EvseId, request.Tariff.TariffId)
	// TODO: store default tariff for EVSE
	return tariffcost.NewSetDefaultTariffResponse(tariffcost.TariffSetStatusAccepted), nil
}

func (handler *ChargingStationHandler) OnGetTariffs(request *tariffcost.GetTariffsRequest) (confirmation *tariffcost.GetTariffsResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("request to get tariffs for EVSE %d", request.EvseId)
	// TODO: retrieve stored tariffs for EVSE
	// For now, return no tariffs
	return tariffcost.NewGetTariffsResponse(tariffcost.TariffGetStatusNoTariff), nil
}

func (handler *ChargingStationHandler) OnClearTariffs(request *tariffcost.ClearTariffsRequest) (confirmation *tariffcost.ClearTariffsResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("request to clear tariffs (EVSE: %d, Tariff IDs: %v)", request.EvseId, request.TariffIds)
	// TODO: clear stored tariffs
	// Build results for each tariff ID, or single result if clearing all
	var results []tariffcost.ClearTariffsResult
	if len(request.TariffIds) > 0 {
		results = make([]tariffcost.ClearTariffsResult, len(request.TariffIds))
		for i, tariffId := range request.TariffIds {
			results[i] = tariffcost.ClearTariffsResult{
				TariffId: tariffId,
				Status:   tariffcost.TariffClearStatusAccepted,
			}
		}
	} else {
		// Clearing all tariffs
		results = []tariffcost.ClearTariffsResult{
			{
				Status: tariffcost.TariffClearStatusAccepted,
			},
		}
	}
	return tariffcost.NewClearTariffsResponse(results), nil
}

func (handler *ChargingStationHandler) OnChangeTransactionTariff(request *tariffcost.ChangeTransactionTariffRequest) (confirmation *tariffcost.ChangeTransactionTariffResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("accepted request to change tariff for transaction %s: %s", request.TransactionId, request.Tariff.TariffId)
	// TODO: update tariff for active transaction
	return tariffcost.NewChangeTransactionTariffResponse(tariffcost.TariffSetStatusAccepted), nil
}
