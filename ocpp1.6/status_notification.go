package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"time"
)

// -------------------- Status Notification (CP -> CS) --------------------
type ChargePointErrorCode string
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

type StatusNotificationRequest struct {
	ConnectorId     int                  `json:"connectorId" validate:"gte=0"`
	ErrorCode       ChargePointErrorCode `json:"errorCode" validate:"required,chargePointErrorCode"`
	Info            string               `json:"info,omitempty" validate:"max=50"`
	Status          ChargePointStatus    `json:"status" validate:"required,chargePointStatus"`
	Timestamp       DateTime             `json:"timestamp,omitempty"`
	VendorId        string               `json:"vendorId,omitempty" validate:"max=255"`
	VendorErrorCode string               `json:"vendorErrorCode,omitempty" validate:"max=50"`
}

type StatusNotificationConfirmation struct {
}

type StatusNotificationFeature struct{}

func (f StatusNotificationFeature) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func (f StatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(StatusNotificationRequest{})
}

func (f StatusNotificationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(StatusNotificationConfirmation{})
}

func (r StatusNotificationRequest) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func (c StatusNotificationConfirmation) GetFeatureName() string {
	return StatusNotificationFeatureName
}

func NewStatusNotificationRequest(connectorId int, errorCode ChargePointErrorCode, status ChargePointStatus) *StatusNotificationRequest {
	return &StatusNotificationRequest{ConnectorId: connectorId, ErrorCode: errorCode, Status: status}
}

func NewStatusNotificationConfirmation() *StatusNotificationConfirmation {
	return &StatusNotificationConfirmation{}
}

func validateStatusNotificationRequest(sl validator.StructLevel) {
	request := sl.Current().Interface().(StatusNotificationRequest)
	if !dateTimeIsNull(request.Timestamp) && !validateDateTimeLt(request.Timestamp, time.Now()) {
		sl.ReportError(request.Timestamp, "Timestamp", "timestamp", "lt", "")
	}
}

func init() {
	//Validate.RegisterStructValidation(validateStatusNotificationRequest, StatusNotificationRequest{})
	_ = Validate.RegisterValidation("chargePointErrorCode", isValidChargePointErrorCode)
	_ = Validate.RegisterValidation("chargePointStatus", isValidChargePointStatus)
}
