package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Unlock Connector (CS -> CP) --------------------

// Status in response to UnlockConnectorRequest.
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

// The field definition of the UnlockConnector request payload sent by the Central System to the Charge Point.
type UnlockConnectorRequest struct {
	ConnectorId int `json:"connectorId" validate:"gt=0"`
}

// This field definition of the UnlockConnector confirmation payload, sent by the Charge Point to the Central System in response to an UnlockConnectorRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type UnlockConnectorConfirmation struct {
	Status UnlockStatus `json:"status" validate:"required,unlockStatus"`
}

// Central System can request a Charge Point to unlock a connector. To do so, the Central System SHALL send an UnlockConnectorRequest.
// The purpose of this message: Help EV drivers that have problems unplugging their cable from the Charge Point in case of malfunction of the Connector cable retention.
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

// Creates a new UnlockConnectorRequest, containing all required fields. There are no optional fields for this message.
func NewUnlockConnectorRequest(connectorId int) *UnlockConnectorRequest {
	return &UnlockConnectorRequest{ConnectorId: connectorId}
}

// Creates a new UnlockConnectorConfirmation, containing all required fields. There are no optional fields for this message.
func NewUnlockConnectorConfirmation(status UnlockStatus) *UnlockConnectorConfirmation {
	return &UnlockConnectorConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("unlockStatus", isValidUnlockStatus)
}
