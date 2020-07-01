package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Status Notification (CP -> CS) --------------------

const StatusNotificationFeatureName = "StatusNotification"

// Charge Point status reported in StatusNotificationRequest.
type ChargePointErrorCode string

// Status reported in StatusNotificationRequest.
// A status can be reported for the Charge Point main controller (connectorId = 0) or for a specific connector.
// Status for the Charge Point main controller is a subset of the enumeration: Available, Unavailable or Faulted.
type ChargePointStatus string

const (
	ConnectorLockFailure           ChargePointErrorCode = "ConnectorLockFailure"
	EVCommunicationError           ChargePointErrorCode = "EVCommunicationError"
	GroundFailure                  ChargePointErrorCode = "GroundFailure"
	HighTemperature                ChargePointErrorCode = "HighTemperature"
	InternalError                  ChargePointErrorCode = "InternalError"
	LocalListConflict              ChargePointErrorCode = "LocalListConflict"
	NoError                        ChargePointErrorCode = "NoError"
	OtherError                     ChargePointErrorCode = "OtherError"
	OverCurrentFailure             ChargePointErrorCode = "OverCurrentFailure"
	OverVoltage                    ChargePointErrorCode = "OverVoltage"
	PowerMeterFailure              ChargePointErrorCode = "PowerMeterFailure"
	PowerSwitchFailure             ChargePointErrorCode = "PowerSwitchFailure"
	ReaderFailure                  ChargePointErrorCode = "ReaderFailure"
	ResetFailure                   ChargePointErrorCode = "ResetFailure"
	UnderVoltage                   ChargePointErrorCode = "UnderVoltage"
	WeakSignal                     ChargePointErrorCode = "WeakSignal"
	ChargePointStatusAvailable     ChargePointStatus    = "Available"
	ChargePointStatusPreparing     ChargePointStatus    = "Preparing"
	ChargePointStatusCharging      ChargePointStatus    = "Charging"
	ChargePointStatusSuspendedEVSE ChargePointStatus    = "SuspendedEVSE"
	ChargePointStatusSuspendedEV   ChargePointStatus    = "SuspendedEV"
	ChargePointStatusFinishing     ChargePointStatus    = "Finishing"
	ChargePointStatusReserved      ChargePointStatus    = "Reserved"
	ChargePointStatusUnavailable   ChargePointStatus    = "Unavailable"
	ChargePointStatusFaulted       ChargePointStatus    = "Faulted"
)

func isValidChargePointStatus(fl validator.FieldLevel) bool {
	status := ChargePointStatus(fl.Field().String())
	switch status {
	case ChargePointStatusAvailable, ChargePointStatusPreparing, ChargePointStatusCharging, ChargePointStatusFaulted, ChargePointStatusFinishing, ChargePointStatusReserved, ChargePointStatusSuspendedEV, ChargePointStatusSuspendedEVSE, ChargePointStatusUnavailable:
		return true
	default:
		return false
	}
}

func isValidChargePointErrorCode(fl validator.FieldLevel) bool {
	status := ChargePointErrorCode(fl.Field().String())
	switch status {
	case ConnectorLockFailure, EVCommunicationError, GroundFailure, HighTemperature, InternalError, LocalListConflict, NoError, OtherError, OverVoltage, OverCurrentFailure, PowerMeterFailure, PowerSwitchFailure, ReaderFailure, ResetFailure, UnderVoltage, WeakSignal:
		return true
	default:
		return false
	}
}

// The field definition of the StatusNotification request payload sent by the Charge Point to the Central System.
type StatusNotificationRequest struct {
	ConnectorId     int                  `json:"connectorId" validate:"gte=0"`
	ErrorCode       ChargePointErrorCode `json:"errorCode" validate:"required,chargePointErrorCode"`
	Info            string               `json:"info,omitempty" validate:"max=50"`
	Status          ChargePointStatus    `json:"status" validate:"required,chargePointStatus"`
	Timestamp       *types.DateTime      `json:"timestamp,omitempty" validate:"omitempty"`
	VendorId        string               `json:"vendorId,omitempty" validate:"max=255"`
	VendorErrorCode string               `json:"vendorErrorCode,omitempty" validate:"max=50"`
}

// This field definition of the StatusNotification confirmation payload, sent by the Central System to the Charge Point in response to a StatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type StatusNotificationConfirmation struct {
}

// A Charge Point sends a notification to the Central System to inform the Central System about a status change or an error within the Charge Point.
// Refer to the official OCPP specification for a full list of status changes that may trigger such a message.
// To limit the number of transitions, the Charge Point MAY omit sending a StatusNotificationRequest if it was active for less time than defined in the optional configuration key MinimumStatusDuration.
// This way, a Charge Point MAY choose not to send certain StatusNotificationRequest payloads.
// The Charge Point MAY send a StatusNotificationRequest to inform the Central System of fault conditions.
// When the 'status' field is not Faulted, the condition should be considered a warning since charging operations are still possible.
// When a Charge Point is configured with StopTransactionOnEVSideDisconnect set to false, a transaction is running and the EV becomes disconnected on EV side, then a StatusNotificationRequest with the state: SuspendedEV SHOULD be send to the Central System, with the 'errorCode' field set to: 'NoError'.
// The Charge Point SHOULD add additional information in the 'info' field, Notifying the Central System with the reason of suspension: 'EV side disconnected'. The current transaction is not stopped.
// When a Charge Point is configured with StopTransactionOnEVSideDisconnect set to true, a transaction is running and the EV becomes disconnected on EV side, then a StatusNotificationRequest with the state: 'Finishing' SHOULD be send to the Central System, with the 'errorCode' field set to: 'NoError'.
// The Charge Point SHOULD add additional information in the 'info' field, Notifying the Central System with the reason of stopping: 'EV side disconnected'. The current transaction is stopped.
// When a Charge Point connects to a Central System after having been offline, it updates the Central System about its status according to the following rules:
// 1. The Charge Point SHOULD send a StatusNotificationRequest with its current status if the status changed while the Charge Point was offline.
// 2. The Charge Point MAY send a StatusNotificationRequest to report an error that occurred while the Charge Point was offline.
// 3. The Charge Point SHOULD NOT send StatusNotificationRequest payloads for historical status change events that happened while the Charge Point was offline and that do not inform the Central System of Charge Point errors or the Charge Point’s current status.
// 4. The StatusNotificationRequest messages MUST be sent in the order in which the events that they describe occurred.
// Upon receipt of a StatusNotificationRequest, the Central System SHALL respond with a StatusNotificationConfirmation.
type StatusNotificationFeature struct{}

func (f StatusNotificationFeature) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func (f StatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(StatusNotificationRequest{})
}

func (f StatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(StatusNotificationConfirmation{})
}

func (r StatusNotificationRequest) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func (c StatusNotificationConfirmation) GetFeatureName() string {
	return StatusNotificationFeatureName
}

// Creates a new StatusNotificationRequest, containing all required fields.
// Optional fields may be set directly on the created request.
func NewStatusNotificationRequest(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus) *StatusNotificationRequest {
	return &StatusNotificationRequest{ConnectorId: connectorId, ErrorCode: errorCode, Status: status}
}

// Creates a new StatusNotificationConfirmation, which doesn't contain any required or optional fields.
func NewStatusNotificationConfirmation() *StatusNotificationConfirmation {
	return &StatusNotificationConfirmation{}
}

func init() {
	_ = types.Validate.RegisterValidation("chargePointErrorCode", isValidChargePointErrorCode)
	_ = types.Validate.RegisterValidation("chargePointStatus", isValidChargePointStatus)
}
