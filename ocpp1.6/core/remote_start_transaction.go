package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"reflect"
)

// -------------------- Remote Start Transaction (CS -> CP) --------------------

const RemoteStartTransactionFeatureName = "RemoteStartTransaction"

// The field definition of the RemoteStartTransaction request payload sent by the Central System to the Charge Point.
type RemoteStartTransactionRequest struct {
	ConnectorId     *int                   `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	IdTag           string                 `json:"idTag" validate:"required,max=20"`
	ChargingProfile *types.ChargingProfile `json:"chargingProfile,omitempty"`
}

// This field definition of the RemoteStartTransaction confirmation payload, sent by the Charge Point to the Central System in response to a RemoteStartTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type RemoteStartTransactionConfirmation struct {
	Status types.RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus"`
}

// Central System can request a Charge Point to start a transaction by sending a RemoteStartTransactionRequest.
// Upon receipt, the Charge Point SHALL reply with RemoteStartTransactionConfirmation and a status indicating whether it has accepted the request and will attempt to start a transaction.
// The effect of the RemoteStartTransactionRequest message depends on the value of the AuthorizeRemoteTxRequests configuration key in the Charge Point.
//
// • If the value of AuthorizeRemoteTxRequests is true, the Charge Point SHALL behave as if in response to a local action at the Charge Point to start a transaction with the idTag given in the RemoteStartTransactionRequest message. This means that the Charge Point will first try to authorize the idTag, using the Local Authorization List, Authorization Cache and/or an AuthorizeRequest request. A transaction will only be started after authorization was obtained.
//
// • If the value of AuthorizeRemoteTxRequests is false, the Charge Point SHALL immediately try to start a transaction for the idTag given in the RemoteStartTransactionRequest message. Note that after the transaction has been started, the Charge Point will send a StartTransaction request to the Central System, and the Central System will check the authorization status of the idTag when processing this StartTransaction request.
//
// The following typical use cases are the reason for Remote Start Transaction:
//
// • Enable a CPO operator to help an EV driver that has problems starting a transaction.
//
// • Enable mobile apps to control charging transactions via the Central System.
//
// • Enable the use of SMS to control charging transactions via the Central System.
//
// The RemoteStartTransactionRequest SHALL contain an identifier (idTag), which Charge Point SHALL use, if it is able to start a transaction, to send a StartTransactionRequest to Central System.
// The transaction is started in the same way as described in StartTransaction. The RemoteStartTransactionRequest MAY contain a connector id if the transaction is to be started on a specific connector. When no connector id is provided, the Charge Point is in control of the connector selection.
// A Charge Point MAY reject a RemoteStartTransactionRequest without a connector id.
// The Central System MAY include a ChargingProfile in the RemoteStartTransaction request. The purpose of this ChargingProfile SHALL be set to TxProfile. If accepted, the Charge Point SHALL use this ChargingProfile for the transaction.
type RemoteStartTransactionFeature struct{}

func (f RemoteStartTransactionFeature) GetFeatureName() string {
	return RemoteStartTransactionFeatureName
}

func (f RemoteStartTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RemoteStartTransactionRequest{})
}

func (f RemoteStartTransactionFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(RemoteStartTransactionConfirmation{})
}

func (r RemoteStartTransactionRequest) GetFeatureName() string {
	return RemoteStartTransactionFeatureName
}

func (c RemoteStartTransactionConfirmation) GetFeatureName() string {
	return RemoteStartTransactionFeatureName
}

// Creates a new RemoteStartTransactionRequest, containing all required fields. Optional fields may be set afterwards.
func NewRemoteStartTransactionRequest(idTag string) *RemoteStartTransactionRequest {
	return &RemoteStartTransactionRequest{IdTag: idTag}
}

// Creates a new RemoteStartTransactionConfirmation, containing all required fields. There are no optional fields for this message.
func NewRemoteStartTransactionConfirmation(status types.RemoteStartStopStatus) *RemoteStartTransactionConfirmation {
	return &RemoteStartTransactionConfirmation{Status: status}
}
