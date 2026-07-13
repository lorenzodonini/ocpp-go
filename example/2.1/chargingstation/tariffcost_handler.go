package main

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/tariffcost"
)

func (handler *ChargingStationHandler) OnCostUpdated(request *tariffcost.CostUpdatedRequest) (response *tariffcost.CostUpdatedResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("accepted request to display cost for transaction %v: %v", request.TransactionID, request.TotalCost)
	// TODO: update internal display to show updated cost for transaction
	return tariffcost.NewCostUpdatedResponse(), nil
}
