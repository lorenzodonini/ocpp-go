// The transactions functional block contains OCPP 2.0 features related to OCPP transactions.
package transactions

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Transactions profile.
type CSMSHandler interface {
	// OnTransactionEvent is called on the CSMS whenever a TransactionEventRequest is received from a charging station.
	OnTransactionEvent(chargingStationID string, request *TransactionEventRequest) (response *TransactionEventResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Transactions profile.
type ChargingStationHandler interface {
	// OnGetTransactionStatusResponse is called on a charging station whenever a OnGetTransactionStatusRequest is received from the CSMS.
	OnGetTransactionStatus(request *GetTransactionStatusRequest) (response *GetTransactionStatusResponse, err error)
}

const ProfileName = "transactions"

var Profile = ocpp.NewProfile(
	ProfileName,
	GetTransactionStatusFeature{},
	TransactionEventFeature{},
)
