package ocpp16

import (
	"reflect"
)

// -------------------- Remote Stop Transaction (CS -> CP) --------------------

type RemoteStopTransactionRequest struct {
	TransactionId int `json:"transactionId" validate:"gte=0"`
}

type RemoteStopTransactionConfirmation struct {
	Status RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus"`
}

type RemoteStopTransactionFeature struct{}

func (f RemoteStopTransactionFeature) GetFeatureName() string {
	return RemoteStopTransactionFeatureName
}

func (f RemoteStopTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RemoteStopTransactionRequest{})
}

func (f RemoteStopTransactionFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(RemoteStopTransactionConfirmation{})
}

func (r RemoteStopTransactionRequest) GetFeatureName() string {
	return RemoteStopTransactionFeatureName
}

func (c RemoteStopTransactionConfirmation) GetFeatureName() string {
	return RemoteStopTransactionFeatureName
}

func NewRemoteStopTransactionRequest(transactionId int) *RemoteStopTransactionRequest {
	return &RemoteStopTransactionRequest{TransactionId: transactionId}
}

func NewRemoteStopTransactionConfirmation(status RemoteStartStopStatus) *RemoteStopTransactionConfirmation {
	return &RemoteStopTransactionConfirmation{Status: status}
}
