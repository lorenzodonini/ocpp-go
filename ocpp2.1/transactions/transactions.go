// Package transactions contains the transactions functional block for OCPP 2.1.
package transactions

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the transactions functional block.
const ProfileName = "transactions"

// ChargingStationHandler contains callbacks for handling incoming CSMS requests.
type ChargingStationHandler interface {
	// OnGetTransactionStatus handles a GetTransactionStatusRequest from the CSMS.
	OnGetTransactionStatus(request *GetTransactionStatusRequest) (response *GetTransactionStatusResponse, err error)
}

// CSMSHandler contains callbacks for handling incoming Charging Station requests.
type CSMSHandler interface {
	// OnTransactionEvent handles a TransactionEventRequest from the Charging Station.
	OnTransactionEvent(request *TransactionEventRequest) (response *TransactionEventResponse, err error)
}

// Profile returns the transactions profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		TransactionEventFeature{},
		GetTransactionStatusFeature{},
	)
}
