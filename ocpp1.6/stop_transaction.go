package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Stop Transaction (CP -> CS) --------------------

// Reason for stopping a transaction in StopTransactionRequest.
type Reason string

const (
	ReasonDeAuthorized   Reason = "DeAuthorized"
	ReasonEmergencyStop  Reason = "EmergencyStop"
	ReasonEVDisconnected Reason = "EVDisconnected"
	ReasonHardReset      Reason = "HardReset"
	ReasonLocal          Reason = "Local"
	ReasonOther          Reason = "Other"
	ReasonPowerLoss      Reason = "PowerLoss"
	ReasonReboot         Reason = "Reboot"
	ReasonRemote         Reason = "Remote"
	ReasonSoftReset      Reason = "SoftReset"
	ReasonUnlockCommand  Reason = "UnlockCommand"
)

func isValidReason(fl validator.FieldLevel) bool {
	reason := Reason(fl.Field().String())
	switch reason {
	case ReasonDeAuthorized, ReasonEmergencyStop, ReasonEVDisconnected, ReasonHardReset, ReasonLocal, ReasonOther, ReasonPowerLoss, ReasonReboot, ReasonRemote, ReasonSoftReset, ReasonUnlockCommand:
		return true
	default:
		return false
	}
}

// The field definition of the StopTransaction request payload sent by the Charge Point to the Central System.
type StopTransactionRequest struct {
	IdTag           string       `json:"idTag,omitempty" validate:"max=20"`
	MeterStop       int          `json:"meterStop" validate:"gte=0"`
	Timestamp       *DateTime    `json:"timestamp" validate:"required"`
	TransactionId   int          `json:"transactionId" validate:"gte=0"`
	Reason          Reason       `json:"reason,omitempty" validate:"omitempty,reason"`
	TransactionData []MeterValue `json:"transactionData,omitempty" validate:"omitempty,dive"`
}

// This field definition of the StopTransaction confirmation payload, sent by the Central System to the Charge Point in response to a StopTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type StopTransactionConfirmation struct {
	IdTagInfo *IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"`
}

// When a transaction is stopped, the Charge Point SHALL send a StopTransactionRequest, notifying to the Central System that the transaction has stopped.
type StopTransactionFeature struct{}

func (f StopTransactionFeature) GetFeatureName() string {
	return StopTransactionFeatureName
}

func (f StopTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(StopTransactionRequest{})
}

func (f StopTransactionFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(StopTransactionConfirmation{})
}

func (r StopTransactionRequest) GetFeatureName() string {
	return StopTransactionFeatureName
}

func (c StopTransactionConfirmation) GetFeatureName() string {
	return StopTransactionFeatureName
}

// Creates a new StopTransactionRequest, containing all required fields.
// Optional fields may be set directly on the created request.
func NewStopTransactionRequest(meterStop int, timestamp *DateTime, transactionId int) *StopTransactionRequest {
	return &StopTransactionRequest{MeterStop: meterStop, Timestamp: timestamp, TransactionId: transactionId}
}

// Creates a new StopTransactionConfirmation. Optional fields may be set afterwards.
func NewStopTransactionConfirmation() *StopTransactionConfirmation {
	return &StopTransactionConfirmation{}
}

//TODO: advanced validation
func init() {
	_ = Validate.RegisterValidation("reason", isValidReason)
}
