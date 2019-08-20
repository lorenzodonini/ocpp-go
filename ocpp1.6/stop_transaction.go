package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Stop Transaction (CP -> CS) --------------------
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

type StopTransactionRequest struct {
	IdTag           string       `json:"idTag,omitempty" validate:"max=20"`
	MeterStop       int          `json:"meterStop" validate:"gte=0"`
	Timestamp       *DateTime    `json:"timestamp" validate:"required"`
	TransactionId   int          `json:"transactionId" validate:"gte=0"`
	Reason          Reason       `json:"reason,omitempty" validate:"omitempty,reason"`
	TransactionData []MeterValue `json:"transactionData,omitempty" validate:"omitempty,dive"`
}

type StopTransactionConfirmation struct {
	IdTagInfo     *IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"`
}

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

func NewStopTransactionRequest(meterStop int, timestamp *DateTime, transactionId int) *StopTransactionRequest {
	return &StopTransactionRequest{MeterStop: meterStop, Timestamp: timestamp, TransactionId: transactionId}
}

func NewStopTransactionConfirmation() *StopTransactionConfirmation {
	return &StopTransactionConfirmation{}
}

func init() {
	_ = Validate.RegisterValidation("reason", isValidReason)
}
