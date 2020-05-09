package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Unlock Connector (CS -> CP) --------------------

const UnlockConnectorFeatureName = "UnlockConnector"

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
// When a EV driver calls the CPO help-desk, an operator could manually trigger the sending of an UnlockConnectorRequest to the Charge Point, forcing a new attempt to unlock the connector.
// Hopefully this time the connector unlocks and the EV driver can unplug the cable and drive away.
// The UnlockConnectorRequest SHOULD NOT be used to remotely stop a running transaction, use the Remote Stop Transaction instead.
// Upon receipt of an UnlockConnectorRequest, the Charge Point SHALL respond with a UnlockConnectorConfirmation.
// The response payload SHALL indicate whether the Charge Point was able to unlock its connector.
// If there was a transaction in progress on the specific connector, then Charge Point SHALL finish the transaction first as described in Stop Transaction.
type UnlockConnectorFeature struct{}

func (f UnlockConnectorFeature) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (f UnlockConnectorFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorRequest{})
}

func (f UnlockConnectorFeature) GetResponseType() reflect.Type {
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
	_ = types.Validate.RegisterValidation("unlockStatus", isValidUnlockStatus)
}
