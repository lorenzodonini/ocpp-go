package remotecontrol

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// -------------------- Trigger Message (CSMS -> CS) --------------------

const UnlockConnectorFeatureName = "UnlockConnector"

// Status in UnlockConnectorResponse.
type UnlockStatus string

const (
	UnlockStatusUnlocked                     UnlockStatus = "Unlocked"                     // Connector has successfully been unlocked.
	UnlockStatusUnlockFailed                 UnlockStatus = "UnlockFailed"                 // Failed to unlock the connector.
	UnlockStatusOngoingAuthorizedTransaction UnlockStatus = "OngoingAuthorizedTransaction" // The connector is not unlocked, because there is still an authorized transaction ongoing.
	UnlockStatusUnknownConnector             UnlockStatus = "UnknownConnector"             // The specified connector is not known by the Charging Station.
)

func isValidUnlockStatus(fl validator.FieldLevel) bool {
	status := UnlockStatus(fl.Field().String())
	switch status {
	case UnlockStatusUnlocked,
		UnlockStatusUnlockFailed,
		UnlockStatusOngoingAuthorizedTransaction,
		UnlockStatusUnknownConnector:
		return true
	default:
		return false
	}
}

// The field definition of the UnlockConnector request payload sent by the CSMS to the Charging Station.
type UnlockConnectorRequest struct {
	EvseID      int `json:"evseId" validate:"gte=0"`
	ConnectorID int `json:"connectorId" validate:"gte=0"`
}

// This field definition of the UnlockConnector response payload, sent by the Charging Station to the CSMS in response to a UnlockConnectorRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type UnlockConnectorResponse struct {
	Status     UnlockStatus      `json:"status" validate:"required,unlockStatus201"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty"`
}

// It sometimes happens that a connector of a Charging Station socket does not unlock correctly.
// This happens most of the time when there is tension on the charging cable.
// This means the driver cannot unplug his charging cable from the Charging Station.
// To help a driver, the CSO can send a UnlockConnectorRequest to the Charging Station.
// The Charging Station will then try to unlock the connector again and respond with an UnlockConnectorResponse.
type UnlockConnectorFeature struct{}

func (f UnlockConnectorFeature) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (f UnlockConnectorFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorRequest{})
}

func (f UnlockConnectorFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(UnlockConnectorResponse{})
}

func (r UnlockConnectorRequest) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

func (c UnlockConnectorResponse) GetFeatureName() string {
	return UnlockConnectorFeatureName
}

// Creates a new UnlockConnectorRequest, containing all required fields. There are no optional fields for this message.
func NewUnlockConnectorRequest(evseID int, connectorID int) *UnlockConnectorRequest {
	return &UnlockConnectorRequest{EvseID: evseID, ConnectorID: connectorID}
}

// Creates a new UnlockConnectorResponse, containing all required fields. Optional fields may be set afterwards.
func NewUnlockConnectorResponse(status UnlockStatus) *UnlockConnectorResponse {
	return &UnlockConnectorResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("unlockStatus201", isValidUnlockStatus)
}
