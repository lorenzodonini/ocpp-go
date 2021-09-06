package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Stop Transaction (CP -> CS) --------------------

const StopTransactionFeatureName = "StopTransaction"

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
	IdTag           string             `json:"idTag,omitempty" validate:"max=20"`
	MeterStop       int                `json:"meterStop"`
	Timestamp       *types.DateTime    `json:"timestamp" validate:"required"`
	TransactionId   int                `json:"transactionId"`
	Reason          Reason             `json:"reason,omitempty" validate:"omitempty,reason"`
	TransactionData []types.MeterValue `json:"transactionData,omitempty" validate:"omitempty,dive"`
}

// This field definition of the StopTransaction confirmation payload, sent by the Central System to the Charge Point in response to a StopTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type StopTransactionConfirmation struct {
	IdTagInfo *types.IdTagInfo `json:"idTagInfo,omitempty" validate:"omitempty"`
}

// When a transaction is stopped, the Charge Point SHALL send a StopTransactionRequest, notifying to the Central System that the transaction has stopped.
// A StopTransactionRequest MAY contain an optional TransactionData element to provide more details about transaction usage.
// The optional TransactionData element is a container for any number of MeterValues, using the same data structure as the meterValue elements of the MeterValuesRequest payload.
// Upon receipt of a StopTransactionRequest, the Central System SHALL respond with a StopTransactionConfirmation.
// The idTag in the request payload MAY be omitted when the Charge Point itself needs to stop the transaction. For instance, when the Charge Point is requested to reset.
// If a transaction is ended in a normal way (e.g. EV-driver presented his identification to stop the transaction), the Reason element MAY be omitted and the Reason SHOULD be assumed 'Local'.
// If the transaction is not ended normally, the Reason SHOULD be set to a correct value. As part of the normal transaction termination, the Charge Point SHALL unlock the cable (if not permanently attached).
// The Charge Point MAY unlock the cable (if not permanently attached) when the cable is disconnected at the EV.
// If supported, this functionality is reported and controlled by the configuration key UnlockConnectorOnEVSideDisconnect.
// The Charge Point MAY stop a running transaction when the cable is disconnected at the EV. If supported, this functionality is reported and controlled by the configuration key StopTransactionOnEVSideDisconnect.
// If StopTransactionOnEVSideDisconnect is set to false, the transaction SHALL not be stopped when the cable is disconnected from the EV.
// If the EV is reconnected, energy transfer is allowed again. In this case there is no mechanism to prevent other EVs from charging and disconnecting during that same ongoing transaction.
// With UnlockConnectorOnEVSideDisconnect set to false, the Connector SHALL remain locked at the Charge Point until the user presents the identifier.
// By setting StopTransactionOnEVSideDisconnect to true, the transaction SHALL be stopped when the cable is disconnected from the EV.
// If the EV is reconnected, energy transfer is not allowed until the transaction is stopped and a new transaction is started.
// If UnlockConnectorOnEVSideDisconnect is set to true, also the Connector on the Charge Point will be unlocked.
// It is likely that The Central System applies sanity checks to the data contained in a StopTransactionRequest it received.
// The outcome of such sanity checks SHOULD NOT ever cause the Central System to not respond with a StopTransactionConfirmation.
// Failing to respond with a StopTransactionConfirmation will only cause the Charge Point to try the same message again as specified in Error responses to transaction-related messages.
// If Charge Point has implemented an Authorization Cache, then upon receipt of a StopTransactionConfirmation the Charge Point SHALL update the cache entry, if the idTag is not in the Local Authorization List, with the IdTagInfo value from the response as described under Authorization Cache.
type StopTransactionFeature struct{}

func (f StopTransactionFeature) GetFeatureName() string {
	return StopTransactionFeatureName
}

func (f StopTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(StopTransactionRequest{})
}

func (f StopTransactionFeature) GetResponseType() reflect.Type {
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
func NewStopTransactionRequest(meterStop int, timestamp *types.DateTime, transactionId int) *StopTransactionRequest {
	return &StopTransactionRequest{MeterStop: meterStop, Timestamp: timestamp, TransactionId: transactionId}
}

// Creates a new StopTransactionConfirmation. Optional fields may be set afterwards.
func NewStopTransactionConfirmation() *StopTransactionConfirmation {
	return &StopTransactionConfirmation{}
}

//TODO: advanced validation
func init() {
	_ = types.Validate.RegisterValidation("reason", isValidReason)
}
