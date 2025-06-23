package transactions

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Transaction Event (CS -> CSMS) --------------------

const TransactionEventFeatureName = "TransactionEvent"

// Contains transaction specific information.
type Transaction struct {
	TransactionID     string              `json:"transactionId" validate:"required,max=36"`
	ChargingState     ChargingState       `json:"chargingState,omitempty" validate:"omitempty,chargingState21"`
	TimeSpentCharging *int                `json:"timeSpentCharging,omitempty" validate:"omitempty"` // Contains the total time that energy flowed from EVSE to EV during the transaction (in seconds).
	StoppedReason     Reason              `json:"stoppedReason,omitempty" validate:"omitempty,stoppedReason21"`
	RemoteStartID     *int                `json:"remoteStartId,omitempty" validate:"omitempty"`
	OperationMode     types.OperationMode `json:"operationMode,omitempty" validate:"omitempty,operationMode"`
	TariffId          *string             `json:"tariffId,omitempty" validate:"omitempty,max=60"`
	TransactionLimit  *TransactionLimit   `json:"transactionLimit,omitempty" validate:"omitempty,dive"`
}

// The field definition of the TransactionEvent request payload sent by the Charging Station to the CSMS.
type TransactionEventRequest struct {
	EventType             TransactionEvent       `json:"eventType" validate:"required,transactionEven21"`
	Timestamp             *types.DateTime        `json:"timestamp" validate:"required"`
	TriggerReason         TriggerReason          `json:"triggerReason" validate:"required,triggerReason21"`
	SequenceNo            int                    `json:"seqNo" validate:"gte=0"`
	Offline               bool                   `json:"offline,omitempty"`
	NumberOfPhasesUsed    *int                   `json:"numberOfPhasesUsed,omitempty" validate:"omitempty,gte=0"`
	CableMaxCurrent       *int                   `json:"cableMaxCurrent,omitempty"`                                                  // The maximum current of the connected cable in Ampere (A).
	ReservationID         *int                   `json:"reservationId,omitempty"`                                                    // The ID of the reservation that terminates as a result of this transaction.
	PreconditioningStatus *PreconditioningStatus `json:"preconditioningStatus,omitempty" validate:"omitempty,preconditioningStatus"` // The preconditioning status of the EV.
	EvseSleep             *bool                  `json:"evseSleep,omitempty"`
	TransactionInfo       Transaction            `json:"transactionInfo" validate:"required"` // Contains transaction specific information.
	IDToken               *types.IdToken         `json:"idToken,omitempty" validate:"omitempty,dive"`
	Evse                  *types.EVSE            `json:"evse,omitempty" validate:"omitempty"`             // Identifies which evse (and connector) of the Charging Station is used.
	MeterValue            []types.MeterValue     `json:"meterValue,omitempty" validate:"omitempty,dive"`  // Contains the relevant meter values.
	CostDetails           *types.CostDetails     `json:"costDetails,omitempty" validate:"omitempty,dive"` // Contains the cost details for this transaction. This can be used to inform the CSMS about the cost of this transaction.
}

// This field definition of the TransactionEventResponse payload, sent by the CSMS to the Charging Station in response to a TransactionEventRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type TransactionEventResponse struct {
	TotalCost                   *float64               `json:"totalCost,omitempty" validate:"omitempty,gte=0"`                        // SHALL only be sent when charging has ended. Final total cost of this transaction, including taxes. To indicate a free transaction, the CSMS SHALL send 0.00.
	ChargingPriority            *int                   `json:"chargingPriority,omitempty" validate:"omitempty,min=-9,max=9"`          // Priority from a business point of view. Default priority is 0, The range is from -9 to 9.
	IDTokenInfo                 *types.IdTokenInfo     `json:"idTokenInfo,omitempty" validate:"omitempty"`                            // Is required when the transactionEventRequest contained an idToken.
	UpdatedPersonalMessage      *types.MessageContent  `json:"updatedPersonalMessage,omitempty" validate:"omitempty,dive"`            // This can contain updated personal message that can be shown to the EV Driver. This can be used to provide updated tariff information.
	UpdatedPersonalMessageExtra []types.MessageContent `json:"updatedPersonalMessageExtra,omitempty" validate:"omitempty,max=4,dive"` // This can contain updated personal message that can be shown to the EV Driver. This can be used to provide updated tariff information.
	TransactionLimit            *TransactionLimit      `json:"transactionLimit,omitempty" validate:"omitempty,dive"`                  // Contains the transaction limit for this transaction. This can be used to inform the Charging Station about the maximum cost, time, energy or SoC for this transaction.
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
	_ = types.Validate.RegisterValidation("transactionEvent21", isValidTransactionEvent)
	_ = types.Validate.RegisterValidation("triggerReason21", isValidTriggerReason)
	_ = types.Validate.RegisterValidation("chargingState21", isValidChargingState)
	_ = types.Validate.RegisterValidation("stoppedReason21", isValidReason)
	_ = types.Validate.RegisterValidation("preconditioningStatus", isValidPreconditioningStatus)
}
