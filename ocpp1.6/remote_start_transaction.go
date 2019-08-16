package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Unlock Connector (CS -> CP) --------------------

type RemoteStartStopStatus string

const (
	RemoteStartStopStatusAccepted RemoteStartStopStatus = "Accepted"
	RemoteStartStopStatusRejected RemoteStartStopStatus = "Rejected"
)

func isValidRemoteStartStopStatus(fl validator.FieldLevel) bool {
	status := RemoteStartStopStatus(fl.Field().String())
	switch status {
	case RemoteStartStopStatusAccepted, RemoteStartStopStatusRejected:
		return true
	default:
		return false
	}
}

type RemoteStartTransactionRequest struct {
	ConnectorId     int              `json:"connectorId,omitempty" validate:"omitempty,gt=0"`
	IdTag           string           `json:"idTag" validate:"required,max=20"`
	ChargingProfile *ChargingProfile `json:"chargingProfile,omitempty"`
}

type RemoteStartTransactionConfirmation struct {
	Status RemoteStartStopStatus `json:"status" validate:"required,remoteStartStopStatus"`
}

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

func NewRemoteStartTransactionRequest(idTag string) *RemoteStartTransactionRequest {
	return &RemoteStartTransactionRequest{IdTag: idTag}
}

func NewRemoteStartTransactionConfirmation(status RemoteStartStopStatus) *RemoteStartTransactionConfirmation {
	return &RemoteStartTransactionConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("remoteStartStopStatus", isValidRemoteStartStopStatus)
}
