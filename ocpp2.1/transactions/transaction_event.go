package transactions

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Transaction Event (CS -> CSMS) --------------------

const TransactionEventFeatureName = "TransactionEvent"

// TransactionEventRequest is the payload for a TransactionEvent request from CS to CSMS.
type TransactionEventRequest struct {
	EventType             types.TransactionEvent      `json:"eventType" validate:"required,transactionEvent21"`
	Timestamp             *types.DateTime             `json:"timestamp" validate:"required"`
	TriggerReason         types.TriggerReason         `json:"triggerReason" validate:"required,triggerReason21"`
	SeqNo                 int                         `json:"seqNo" validate:"required,gte=0"`
	TransactionInfo       types.Transaction           `json:"transactionInfo" validate:"required"`
	Offline               bool                        `json:"offline,omitempty"`
	NumberOfPhasesUsed    *int                        `json:"numberOfPhasesUsed,omitempty" validate:"omitempty,gte=0,lte=3"`
	CableMaxCurrent       *int                        `json:"cableMaxCurrent,omitempty"`
	ReservationID         *int                        `json:"reservationId,omitempty" validate:"omitempty,gte=0"`
	EVSE                  *types.EVSE                 `json:"evse,omitempty" validate:"omitempty"`
	IdToken               *types.IdToken              `json:"idToken,omitempty" validate:"omitempty"`
	MeterValue            []types.MeterValue          `json:"meterValue,omitempty" validate:"omitempty,dive"`
	PreconditioningStatus types.PreconditioningStatus `json:"preconditioningStatus,omitempty" validate:"omitempty,preconditioningStatus21"`
	EvseSleep             bool                        `json:"evseSleep,omitempty"`
	CostDetails           *types.CostDetails          `json:"costDetails,omitempty" validate:"omitempty"`
	CustomData            *types.CustomDataType       `json:"customData,omitempty" validate:"omitempty"`
}

// TransactionEventResponse is the payload for a TransactionEvent response from CSMS to CS.
type TransactionEventResponse struct {
	TotalCost                   *float64                `json:"totalCost,omitempty"`
	ChargingPriority            *int                    `json:"chargingPriority,omitempty" validate:"omitempty,gte=-9,lte=9"`
	IdTokenInfo                 *types.IdTokenInfo      `json:"idTokenInfo,omitempty" validate:"omitempty"`
	TransactionLimit            *types.TransactionLimit `json:"transactionLimit,omitempty" validate:"omitempty"`
	UpdatedPersonalMessage      *types.MessageContent   `json:"updatedPersonalMessage,omitempty" validate:"omitempty"`
	UpdatedPersonalMessageExtra []types.MessageContent  `json:"updatedPersonalMessageExtra,omitempty" validate:"omitempty,min=1,max=4,dive"`
	CustomData                  *types.CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// TransactionEventFeature represents the TransactionEvent feature.
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

// NewTransactionEventRequest creates a new TransactionEventRequest.
func NewTransactionEventRequest(eventType types.TransactionEvent, timestamp *types.DateTime, triggerReason types.TriggerReason, seqNo int, transactionInfo types.Transaction) *TransactionEventRequest {
	return &TransactionEventRequest{
		EventType:       eventType,
		Timestamp:       timestamp,
		TriggerReason:   triggerReason,
		SeqNo:           seqNo,
		TransactionInfo: transactionInfo,
	}
}

// NewTransactionEventResponse creates a new TransactionEventResponse.
func NewTransactionEventResponse() *TransactionEventResponse {
	return &TransactionEventResponse{}
}
