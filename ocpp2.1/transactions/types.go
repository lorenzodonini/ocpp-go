package transactions

import "gopkg.in/go-playground/validator.v9"

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

	ReasonDeAuthorized         Reason = "DeAuthorized"         // The transaction was stopped because of the authorization status in the response to a transactionEventRequest.
	ReasonEmergencyStop        Reason = "EmergencyStop"        // Emergency stop button was used.
	ReasonEnergyLimitReached   Reason = "EnergyLimitReached"   // EV charging session reached a locally enforced maximum energy transfer limit.
	ReasonEVDisconnected       Reason = "EVDisconnected"       // Disconnecting of cable, vehicle moved away from inductive charge unit.
	ReasonGroundFault          Reason = "GroundFault"          // A GroundFault has occurred.
	ReasonImmediateReset       Reason = "ImmediateReset"       // A Reset(Immediate) command was received.
	ReasonLocal                Reason = "Local"                // Stopped locally on request of the EV Driver at the Charging Station. This is a regular termination of a transaction.
	ReasonLocalOutOfCredit     Reason = "LocalOutOfCredit"     // A local credit limit enforced through the Charging Station has been exceeded.
	ReasonMasterPass           Reason = "MasterPass"           // The transaction was stopped using a token with a MasterPassGroupId.
	ReasonOther                Reason = "Other"                // Any other reason.
	ReasonOvercurrentFault     Reason = "OvercurrentFault"     // A larger than intended electric current has occurred.
	ReasonPowerLoss            Reason = "PowerLoss"            // Complete loss of power.
	ReasonPowerQuality         Reason = "PowerQuality"         // Quality of power too low, e.g. voltage too low/high, phase imbalance, etc.
	ReasonReboot               Reason = "Reboot"               // A locally initiated reset/reboot occurred.
	ReasonRemote               Reason = "Remote"               // Stopped remotely on request of the CSMS. This is a regular termination of a transaction.
	ReasonSOCLimitReached      Reason = "SOCLimitReached"      // Electric vehicle has reported reaching a locally enforced maximum battery State of Charge (SOC).
	ReasonStoppedByEV          Reason = "StoppedByEV"          // The transaction was stopped by the EV.
	ReasonTimeLimitReached     Reason = "TimeLimitReached"     // EV charging session reached a locally enforced time limit.
	ReasonTimeout              Reason = "Timeout"              // EV not connected within timeout.
	ReasonCostLimitReached     Reason = "CostLimitReached"     // Maximum cost has been reached, as defined by transactionLimit.maxCost.
	ReasonLimitSet             Reason = "LimitSet"             // Limit of cost/time/energy/SoC for transaction has set or changed
	ReasonOperationModeChanged Reason = "OperationModeChanged" // V2X operation mode has changed (at start of a new charging schedule period).
	ReasonRunningCost          Reason = "RunningCost"          // Trigger used when TransactionEvent is sent (only) to report a running cost update.
	ReasonSoCLimitReached      Reason = "SoCLimitReached"      // State of charge limit has been reached, as defined by transactionLimit.maxSoC
	ReasonTariffChanged        Reason = "TariffChanged"        // Tariff for transaction has changed.
	ReasonTariffNotAccepted    Reason = "TariffNotAccepted"    // Trigger to notify that EV Driver has not accepted the tariff for transaction. idToken becomes deauthorized.
	ReasonTxResumed            Reason = "TxResumed"            // Transaction has resumed after reset or power outage.
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
		ReasonSOCLimitReached, ReasonStoppedByEV, ReasonTimeLimitReached, ReasonTimeout,
		ReasonCostLimitReached, ReasonLimitSet, ReasonOperationModeChanged,
		ReasonRunningCost, ReasonSoCLimitReached, ReasonTariffChanged, ReasonTariffNotAccepted, ReasonTxResumed:
		return true
	default:
		return false
	}
}

type PreconditioningStatus string

const (
	PreconditioningStatusNotReady        PreconditioningStatus = "NotReady"
	PreconditioningStatusReady           PreconditioningStatus = "Ready"
	PreconditioningStatusUnknown         PreconditioningStatus = "Unknown"
	PreconditioningStatusPreconditioning PreconditioningStatus = "Preconditioning"
)

func isValidPreconditioningStatus(fl validator.FieldLevel) bool {
	status := PreconditioningStatus(fl.Field().String())
	switch status {
	case PreconditioningStatusNotReady, PreconditioningStatusReady, PreconditioningStatusUnknown, PreconditioningStatusPreconditioning:
		return true
	default:
		return false
	}
}

type TransactionLimit struct {
	MaxCost   *float64 `json:"maxCost,omitempty" validate:"omitempty"`
	MaxEnergy *float64 `json:"maxEnergy,omitempty" validate:"omitempty"`
	MaxTime   *int     `json:"maxTime,omitempty" validate:"omitempty"`
	MaxSoC    *int     `json:"maxSoC,omitempty" validate:"omitempty,gte=0,lte=100"` // Percentage of battery charge
}
