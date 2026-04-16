package availability

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Status Notification (CS -> CSMS) --------------------

const StatusNotificationFeatureName = "StatusNotification"

// ConnectorStatus represents the status of a connector.
type ConnectorStatus string

const (
	ConnectorStatusAvailable   ConnectorStatus = "Available"
	ConnectorStatusOccupied    ConnectorStatus = "Occupied"
	ConnectorStatusReserved    ConnectorStatus = "Reserved"
	ConnectorStatusUnavailable ConnectorStatus = "Unavailable"
	ConnectorStatusFaulted     ConnectorStatus = "Faulted"
)

func isValidConnectorStatus(fl validator.FieldLevel) bool {
	s := ConnectorStatus(fl.Field().String())
	switch s {
	case ConnectorStatusAvailable, ConnectorStatusOccupied, ConnectorStatusReserved,
		ConnectorStatusUnavailable, ConnectorStatusFaulted:
		return true
	default:
		return false
	}
}

// StatusNotificationRequest is the payload for a StatusNotification request from CS to CSMS.
type StatusNotificationRequest struct {
	Timestamp       types.DateTime        `json:"timestamp" validate:"required"`
	ConnectorStatus ConnectorStatus       `json:"connectorStatus" validate:"required,connectorStatus21"`
	EvseID          int                   `json:"evseId" validate:"gte=0"`
	ConnectorID     int                   `json:"connectorId" validate:"gte=0"`
	CustomData      *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// StatusNotificationResponse is the payload for a StatusNotification response from CSMS to CS.
type StatusNotificationResponse struct {
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// StatusNotificationFeature represents the StatusNotification feature.
type StatusNotificationFeature struct{}

func (f StatusNotificationFeature) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func (f StatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(StatusNotificationRequest{})
}

func (f StatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(StatusNotificationResponse{})
}

func (r StatusNotificationRequest) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func (c StatusNotificationResponse) GetFeatureName() string {
	return StatusNotificationFeatureName
}

// NewStatusNotificationRequest creates a new StatusNotificationRequest.
func NewStatusNotificationRequest(timestamp types.DateTime, status ConnectorStatus, evseID int, connectorID int) *StatusNotificationRequest {
	return &StatusNotificationRequest{
		Timestamp:       timestamp,
		ConnectorStatus: status,
		EvseID:          evseID,
		ConnectorID:     connectorID,
	}
}

// NewStatusNotificationResponse creates a new StatusNotificationResponse.
func NewStatusNotificationResponse() *StatusNotificationResponse {
	return &StatusNotificationResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("connectorStatus21", isValidConnectorStatus)
}
