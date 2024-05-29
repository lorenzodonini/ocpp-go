package availability

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Status Notification (CS -> CSMS) --------------------

const StatusNotificationFeatureName = "StatusNotification"

type ConnectorStatus string

const (
	ConnectorStatusAvailable   ConnectorStatus = "Available"   // When a Connector becomes available for a new User (Operative)
	ConnectorStatusOccupied    ConnectorStatus = "Occupied"    // When a Connector becomes occupied, so it is not available for a new EV driver. (Operative)
	ConnectorStatusReserved    ConnectorStatus = "Reserved"    // When a Connector becomes reserved as a result of ReserveNow command (Operative)
	ConnectorStatusUnavailable ConnectorStatus = "Unavailable" // When a Connector becomes unavailable as the result of a Change Availability command or an event upon which the Charging Station transitions to unavailable at its discretion.
	ConnectorStatusFaulted     ConnectorStatus = "Faulted"     // When a Connector (or the EVSE or the entire Charging Station it belongs to) has reported an error and is not available for energy delivery. (Inoperative).
)

func isValidConnectorStatus(fl validator.FieldLevel) bool {
	status := ConnectorStatus(fl.Field().String())
	switch status {
	case ConnectorStatusAvailable, ConnectorStatusOccupied, ConnectorStatusReserved, ConnectorStatusUnavailable, ConnectorStatusFaulted:
		return true
	default:
		return false
	}
}

// The field definition of the StatusNotification request payload sent by the Charging Station to the CSMS.
type StatusNotificationRequest struct {
	Timestamp       *types.DateTime `json:"timestamp" validate:"required"`
	ConnectorStatus ConnectorStatus `json:"connectorStatus" validate:"required,connectorStatus"`
	EvseID          int             `json:"evseId" validate:"gte=0"`
	ConnectorID     int             `json:"connectorId" validate:"gte=0"`
}

// This field definition of the StatusNotification response payload, sent by the CSMS to the Charging Station in response to a StatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type StatusNotificationResponse struct {
}

// The Charging Station notifies the CSMS about a connector status change.
// This may typically be after on of the following events:
//   - (re)boot
//   - reset
//   - any transaction event (start/stop/authorization)
//   - reservation events
//   - change availability operations
//   - remote triggers
//
// The charging station sends a StatusNotificationRequest to the CSMS with information about the new status.
// The CSMS responds with a StatusNotificationResponse.
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

// Creates a new StatusNotificationRequest, containing all required fields. There are no optional fields for this message.
func NewStatusNotificationRequest(timestamp *types.DateTime, status ConnectorStatus, evseID int, connectorID int) *StatusNotificationRequest {
	return &StatusNotificationRequest{Timestamp: timestamp, ConnectorStatus: status, EvseID: evseID, ConnectorID: connectorID}
}

// Creates a new StatusNotificationResponse, which doesn't contain any required or optional fields.
func NewStatusNotificationResponse() *StatusNotificationResponse {
	return &StatusNotificationResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("connectorStatus", isValidConnectorStatus)
}
