package main

import "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/transactions"

func (c *CSMSHandler) OnTransactionEvent(chargingStationID string, request *transactions.TransactionEventRequest) (response *transactions.TransactionEventResponse, err error) {
	switch request.EventType {
	case transactions.TransactionEventStarted:
		logDefault(chargingStationID, request.GetFeatureName()).Infof("transaction %v started, reason: %v, state: %v", request.TransactionInfo.TransactionID, request.TriggerReason, request.TransactionInfo.ChargingState)
		break
	case transactions.TransactionEventUpdated:
		logDefault(chargingStationID, request.GetFeatureName()).Infof("transaction %v updated, reason: %v, state: %v\n", request.TransactionInfo.TransactionID, request.TriggerReason, request.TransactionInfo.ChargingState)
		for _, mv := range request.MeterValue {
			logDefault(chargingStationID, request.GetFeatureName()).Printf("%v", mv)
		}
		break
	case transactions.TransactionEventEnded:
		logDefault(chargingStationID, request.GetFeatureName()).Infof("transaction %v stopped, reason: %v, state: %v\n", request.TransactionInfo.TransactionID, request.TriggerReason, request.TransactionInfo.ChargingState)
		break
	}
	response = transactions.NewTransactionEventResponse()
	return
}
