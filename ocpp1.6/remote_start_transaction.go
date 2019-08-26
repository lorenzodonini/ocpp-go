package ocpp16

import (
	"reflect"
)

// -------------------- Remote Start Transaction (CS -> CP) --------------------

// The field definition of the RemoteStartTransaction request payload sent by the Central System to the Charge Point.
type RemoteStartTransactionRequest struct {
	ConnectorId     int              `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	IdTag           string           `json:"idTag" validate:"required,max=20"`
	ChargingProfile *ChargingProfile `json:"chargingProfile,omitempty"`
}

// This field definition of the RemoteStartTransaction confirmation payload, sent by the Charge Point to the Central System in response to a RemoteStartTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type RemoteStartTransactionConfirmation struct {
	Status RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus"`
}

// Central System can request a Charge Point to start a transaction by sending a RemoteStartTransactionRequest.
// Upon receipt, the Charge Point SHALL reply with RemoteStartTransactionConfirmation and a status indicating whether it has accepted the request and will attempt to start a transaction.
type RemoteStartTransactionFeature struct{}

func (f RemoteStartTransactionFeature) GetFeatureName() string {
	return RemoteStartTransactionFeatureName
}

func (f RemoteStartTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RemoteStartTransactionRequest{})
}

func (f RemoteStartTransactionFeature) GetConfirmationType() reflect.Type {
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
func NewRemoteStartTransactionConfirmation(status RemoteStartStopStatus) *RemoteStartTransactionConfirmation {
	return &RemoteStartTransactionConfirmation{Status: status}
}
