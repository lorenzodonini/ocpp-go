package transactions

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Transaction Event (CS -> CSMS) --------------------

const TransactionEventFeatureName = "TransactionEvent"

// The type of a transaction event.
type TransactionEvent string

// Reason that triggered a TransactionEventRequest.
type TriggerReason string

// The state of the charging process.
type ChargingState string

// Reason for stopping a transaction.
type Reason string

const (
	TransactionEventStarted TransactionEvent = "Started" // First event of a transaction.
	TransactionEventUpdated TransactionEvent = "Updated" // Transaction event in between 'Started' and 'Ended'.
	TransactionEventEnded   TransactionEvent = "Ended"   // Last event of a transaction

	TriggerReasonAuthorized           TriggerReason = "Authorized"           // Charging is authorized, by any means.
	TriggerReasonCablePluggedIn       TriggerReason = "CablePluggedIn"       // Cable is plugged in and EVDetected.
	TriggerReasonChargingRateChanged  TriggerReason = "ChargingRateChanged"  // Rate of charging changed by more than LimitChangeSignificance.
	TriggerReasonChargingStateChanged TriggerReason = "ChargingStateChanged" // Charging state changed.
	TriggerReasonDeAuthorized         TriggerReason = "Deauthorized"         // The transaction was stopped because of the authorization status in the response to a transactionEventRequest.
	TriggerReasonEnergyLimitReached   TriggerReason = "EnergyLimitReached"   // Maximum energy of charging reached. For example: in a pre-paid charging solution
	TriggerReasonEVCommunicationLost  TriggerReason = "EVCommunicationLost"  // Communication with EV lost, for example: cable disconnected.
	TriggerReasonEVConnectTimeout     TriggerReason = "EVConnectTimeout"     // EV not connected before the connection is timed out.
	TriggerReasonMeterValueClock      TriggerReason = "MeterValueClock"      // Needed to send a clock aligned meter value.
	TriggerReasonMeterValuePeriodic   TriggerReason = "MeterValuePeriodic"   // Needed to send a periodic meter value.
	TriggerReasonTimeLimitReached     TriggerReason = "TimeLimitReached"     // Maximum time of charging reached. For example: in a pre-paid charging solution
	TriggerReasonTrigger              TriggerReason = "Trigger"              // Requested by the CSMS via a TriggerMessageRequest.
	TriggerReasonUnlockCommand        TriggerReason = "UnlockCommand"        // CSMS sent an Unlock Connector command.
	TriggerReasonStopAuthorized       TriggerReason = "StopAuthorized"       // An EV Driver has been authorized to stop charging.
	TriggerReasonEVDeparted           TriggerReason = "EVDeparted"           // EV departed. For example: When a departing EV triggers a parking bay detector.
	TriggerReasonEVDetected           TriggerReason = "EVDetected"           // EV detected. For example: When an arriving EV triggers a parking bay detector.
	TriggerReasonRemoteStop           TriggerReason = "RemoteStop"           // A RequestStopTransactionRequest has been sent.
	TriggerReasonRemoteStart          TriggerReason = "RemoteStart"          // A RequestStartTransactionRequest has been sent.
	TriggerReasonAbnormalCondition    TriggerReason = "AbnormalCondition"    // An Abnormal Error or Fault Condition has occurred.
	TriggerReasonSignedDataReceived   TriggerReason = "SignedDataReceived"   // Signed data is received from the energy meter.
	TriggerReasonResetCommand         TriggerReason = "ResetCommand"         // CSMS sent a Reset Charging Station command.

	ChargingStateCharging      ChargingState = "Charging"      // The contactor of the Connector is closed and energy is flowing to between EVSE and EV.
	ChargingStateEVConnected   ChargingState = "EVConnected"   // There is a connection between EV and EVSE (wired or wireless).
	ChargingStateSuspendedEV   ChargingState = "SuspendedEV"   // When the EV is connected to the EVSE and the EVSE is offering energy but the EV is not taking any energy.
	ChargingStateSuspendedEVSE ChargingState = "SuspendedEVSE" // When the EV is connected to the EVSE but the EVSE is not offering energy to the EV (e.g. due to smart charging, power constraints, authorization status).
	ChargingStateIdle          ChargingState = "Idle"          // There is no connection between EV and EVSE.

	ReasonDeAuthorized       Reason = "DeAuthorized"       // The transaction was stopped because of the authorization status in the response to a transactionEventRequest.
	ReasonEmergencyStop      Reason = "EmergencyStop"      // Emergency stop button was used.
	ReasonEnergyLimitReached Reason = "EnergyLimitReached" // EV charging session reached a locally enforced maximum energy transfer limit.
	ReasonEVDisconnected     Reason = "EVDisconnected"     // Disconnecting of cable, vehicle moved away from inductive charge unit.
	ReasonGroundFault        Reason = "GroundFault"        // A GroundFault has occurred.
	ReasonImmediateReset     Reason = "ImmediateReset"     // A Reset(Immediate) command was received.
	ReasonLocal              Reason = "Local"              // Stopped locally on request of the EV Driver at the Charging Station. This is a regular termination of a transaction.
	ReasonLocalOutOfCredit   Reason = "LocalOutOfCredit"   // A local credit limit enforced through the Charging Station has been exceeded.
	ReasonMasterPass         Reason = "MasterPass"         // The transaction was stopped using a token with a MasterPassGroupId.
	ReasonOther              Reason = "Other"              // Any other reason.
	ReasonOvercurrentFault   Reason = "OvercurrentFault"   // A larger than intended electric current has occurred.
	ReasonPowerLoss          Reason = "PowerLoss"          // Complete loss of power.
	ReasonPowerQuality       Reason = "PowerQuality"       // Quality of power too low, e.g. voltage too low/high, phase imbalance, etc.
	ReasonReboot             Reason = "Reboot"             // A locally initiated reset/reboot occurred.
	ReasonRemote             Reason = "Remote"             // Stopped remotely on request of the CSMS. This is a regular termination of a transaction.
	ReasonSOCLimitReached    Reason = "SOCLimitReached"    // Electric vehicle has reported reaching a locally enforced maximum battery State of Charge (SOC).
	ReasonStoppedByEV        Reason = "StoppedByEV"        // The transaction was stopped by the EV.
	ReasonTimeLimitReached   Reason = "TimeLimitReached"   // EV charging session reached a locally enforced time limit.
	ReasonTimeout            Reason = "Timeout"            // EV not connected within timeout.
)

func isValidTransactionEvent(fl validator.FieldLevel) bool {
	status := TransactionEvent(fl.Field().String())
	switch status {
	case TransactionEventStarted, TransactionEventUpdated, TransactionEventEnded:
		return true
	default:
		return false
	}
}

func isValidTriggerReason(fl validator.FieldLevel) bool {
	status := TriggerReason(fl.Field().String())
	switch status {
	case TriggerReasonAuthorized, TriggerReasonCablePluggedIn, TriggerReasonChargingRateChanged,
		TriggerReasonChargingStateChanged, TriggerReasonDeAuthorized, TriggerReasonEnergyLimitReached,
		TriggerReasonEVCommunicationLost, TriggerReasonEVConnectTimeout, TriggerReasonMeterValueClock,
		TriggerReasonMeterValuePeriodic, TriggerReasonTimeLimitReached, TriggerReasonTrigger,
		TriggerReasonUnlockCommand, TriggerReasonStopAuthorized, TriggerReasonEVDeparted,
		TriggerReasonEVDetected, TriggerReasonRemoteStop, TriggerReasonRemoteStart,
		TriggerReasonAbnormalCondition, TriggerReasonSignedDataReceived, TriggerReasonResetCommand:
		return true
	default:
		return false
	}
}

func isValidChargingState(fl validator.FieldLevel) bool {
	status := ChargingState(fl.Field().String())
	switch status {
	case ChargingStateCharging, ChargingStateEVConnected, ChargingStateSuspendedEV, ChargingStateSuspendedEVSE, ChargingStateIdle:
		return true
	default:
		return false
	}
}

func isValidReason(fl validator.FieldLevel) bool {
	status := Reason(fl.Field().String())
	switch status {
	case ReasonDeAuthorized, ReasonEmergencyStop, ReasonEnergyLimitReached, ReasonEVDisconnected,
		ReasonGroundFault, ReasonImmediateReset, ReasonLocal, ReasonLocalOutOfCredit, ReasonMasterPass,
		ReasonOther, ReasonOvercurrentFault, ReasonPowerLoss, ReasonPowerQuality, ReasonReboot, ReasonRemote,
		ReasonSOCLimitReached, ReasonStoppedByEV, ReasonTimeLimitReached, ReasonTimeout:
		return true
	default:
		return false
	}
}

// Contains transaction specific information.
type Transaction struct {
	TransactionID     string        `json:"transactionId" validate:"required,max=36"`
	ChargingState     ChargingState `json:"chargingState,omitempty" validate:"omitempty,chargingState"`
	TimeSpentCharging *int          `json:"timeSpentCharging,omitempty" validate:"omitempty"`
	StoppedReason     Reason        `json:"stoppedReason,omitempty" validate:"omitempty,stoppedReason"`
	RemoteStartID     *int          `json:"remoteStartId,omitempty" validate:"omitempty"`
}

// The field definition of the TransactionEvent request payload sent by the Charging Station to the CSMS.
type TransactionEventRequest struct {
	EventType          TransactionEvent   `json:"eventType" validate:"required,transactionEvent"`
	Timestamp          *types.DateTime    `json:"timestamp" validate:"required"`
	TriggerReason      TriggerReason      `json:"triggerReason" validate:"required,triggerReason"`
	SequenceNo         int                `json:"seqNo" validate:"gte=0"`
	Offline            bool               `json:"offline,omitempty"`
	NumberOfPhasesUsed *int               `json:"numberOfPhasesUsed,omitempty" validate:"omitempty,gte=0"`
	CableMaxCurrent    *int               `json:"cableMaxCurrent,omitempty"`           // The maximum current of the connected cable in Ampere (A).
	ReservationID      *int               `json:"reservationId,omitempty"`             // The Id of the reservation that terminates as a result of this transaction.
	TransactionInfo    Transaction        `json:"transactionInfo" validate:"required"` // Contains transaction specific information.
	IDToken            *types.IdToken     `json:"idToken,omitempty" validate:"omitempty"`
	Evse               *types.EVSE        `json:"evse,omitempty" validate:"omitempty"`            // Identifies which evse (and connector) of the Charging Station is used.
	MeterValue         []types.MeterValue `json:"meterValue,omitempty" validate:"omitempty,dive"` // Contains the relevant meter values.
}

// This field definition of the TransactionEventResponse payload, sent by the CSMS to the Charging Station in response to a TransactionEventRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type TransactionEventResponse struct {
	TotalCost              *float64              `json:"totalCost,omitempty" validate:"omitempty,gte=0"`               // SHALL only be sent when charging has ended. Final total cost of this transaction, including taxes. To indicate a free transaction, the CSMS SHALL send 0.00.
	ChargingPriority       *int                  `json:"chargingPriority,omitempty" validate:"omitempty,min=-9,max=9"` // Priority from a business point of view. Default priority is 0, The range is from -9 to 9.
	IDTokenInfo            *types.IdTokenInfo    `json:"idTokenInfo,omitempty" validate:"omitempty"`                   // Is required when the transactionEventRequest contained an idToken.
	UpdatedPersonalMessage *types.MessageContent `json:"updatedPersonalMessage,omitempty" validate:"omitempty"`        // This can contain updated personal message that can be shown to the EV Driver. This can be used to provide updated tariff information.
}

// Gives the CSMS information that will later be used to bill a transaction.
// For this purpose, status changes and additional transaction-related information is sent, such as
// retrying and sequence number messages.
//
// A Charging Station notifies the CSMS using a TransactionEventRequest. The CSMS then responds with a
// TransactionEventResponse.
type TransactionEventFeature struct{}

func (f TransactionEventFeature) GetFeatureName() string {
	return TransactionEventFeatureName
}

func (f TransactionEventFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(TransactionEventRequest{})
}

func (f TransactionEventFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(TransactionEventResponse{})
}

func (r TransactionEventRequest) GetFeatureName() string {
	return TransactionEventFeatureName
}

func (c TransactionEventResponse) GetFeatureName() string {
	return TransactionEventFeatureName
}

// Creates a new TransactionEventRequest, containing all required fields. Optional fields may be set afterwards.
func NewTransactionEventRequest(t TransactionEvent, timestamp *types.DateTime, reason TriggerReason, seqNo int, info Transaction) *TransactionEventRequest {
	return &TransactionEventRequest{EventType: t, Timestamp: timestamp, TriggerReason: reason, SequenceNo: seqNo, TransactionInfo: info}
}

// Creates a new TransactionEventResponse, containing all required fields. Optional fields may be set afterwards.
func NewTransactionEventResponse() *TransactionEventResponse {
	return &TransactionEventResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("transactionEvent", isValidTransactionEvent)
	_ = types.Validate.RegisterValidation("triggerReason", isValidTriggerReason)
	_ = types.Validate.RegisterValidation("chargingState", isValidChargingState)
	_ = types.Validate.RegisterValidation("stoppedReason", isValidReason)
}
