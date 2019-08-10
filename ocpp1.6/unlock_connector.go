package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Unlock Connector (CS -> CP) --------------------

type UnlockStatus string

const (
	UnlockStatusUnlocked     UnlockStatus = "Unlocked"
	UnlockStatusUnlockFailed UnlockStatus = "UnlockFailed"
	UnlockStatusNotSupported UnlockStatus = "NotSupported"
)

func isValidUnlockStatus(fl validator.FieldLevel) bool {
	status := UnlockStatus(fl.Field().String())
	switch status {
	case UnlockStatusUnlocked, UnlockStatusUnlockFailed, UnlockStatusNotSupported:
		return true
	default:
		return false
	}
}

type UnlockConnectorRequest struct {
	ConnectorId int `json:"connectorId" validate:"gt=0"`
}

type UnlockConnectorConfirmation struct {
	Status UnlockStatus `json:"status" validate:"required,unlockStatus"`
}

type UnlockConnectorFeature struct{}

func (f UnlockConnectorFeature) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (f UnlockConnectorFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorRequest{})
}

func (f UnlockConnectorFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorConfirmation{})
}

func (r UnlockConnectorRequest) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (c UnlockConnectorConfirmation) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func NewUnlockConnectorRequest(connectorId int) *UnlockConnectorRequest {
	return &UnlockConnectorRequest{ConnectorId: connectorId}
}

func NewUnlockConnectorConfirmation(status UnlockStatus) *UnlockConnectorConfirmation {
	return &UnlockConnectorConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("unlockStatus", isValidUnlockStatus)
}
