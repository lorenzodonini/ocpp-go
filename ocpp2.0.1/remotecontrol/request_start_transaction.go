package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Request Start Transaction (CSMS -> CS) --------------------

const RequestStartTransactionFeatureName = "RequestStartTransaction"

// Status reported in RequestStartTransactionResponse.
type RequestStartStopStatus string

const (
	RequestStartStopStatusAccepted RequestStartStopStatus = "Accepted"
	RequestStartStopStatusRejected RequestStartStopStatus = "Rejected"
)

func isValidRequestStartStopStatus(fl validator.FieldLevel) bool {
	status := RequestStartStopStatus(fl.Field().String())
	switch status {
	case RequestStartStopStatusAccepted, RequestStartStopStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the RequestStartTransaction request payload sent by the CSMS to the Charging Station.
type RequestStartTransactionRequest struct {
	EvseID          *int                   `json:"evseId,omitempty" validate:"omitempty,gt=0"`
	RemoteStartID   int                    `json:"remoteStartId" validate:"gte=0"`
	IDToken         types.IdToken          `json:"idToken"`
	ChargingProfile *types.ChargingProfile `json:"chargingProfile,omitempty"`
	GroupIdToken    *types.IdToken         `json:"groupIdToken,omitempty" validate:"omitempty,dive"`
}

// This field definition of the RequestStartTransaction response payload, sent by the Charging Station to the CSMS in response to a RequestStartTransactionRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type RequestStartTransactionResponse struct {
	Status        RequestStartStopStatus `json:"status" validate:"required,requestStartStopStatus"`
	TransactionID string                 `json:"transactionId,omitempty" validate:"max=36"`
	StatusInfo    *types.StatusInfo      `json:"statusInfo,omitempty"`
}

// The CSMS may remotely start a transaction for a user.
// This functionality may be triggered by:
//   - a CSO, to help out a user, that is having trouble starting a transaction
//   - a third-party event (e.g. mobile app)
//   - a previously set ChargingProfile
//
// The CSMS sends a RequestStartTransactionRequest to the Charging Station.
// The Charging Stations will reply with a RequestStartTransactionResponse.
type RequestStartTransactionFeature struct{}

func (f RequestStartTransactionFeature) GetFeatureName() string {
	return RequestStartTransactionFeatureName
}

func (f RequestStartTransactionFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(RequestStartTransactionRequest{})
}

func (f RequestStartTransactionFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(RequestStartTransactionResponse{})
}

func (r RequestStartTransactionRequest) GetFeatureName() string {
	return RequestStartTransactionFeatureName
}

func (c RequestStartTransactionResponse) GetFeatureName() string {
	return RequestStartTransactionFeatureName
}

// Creates a new RequestStartTransactionRequest, containing all required fields. Optional fields may be set afterwards.
func NewRequestStartTransactionRequest(remoteStartID int, IdToken types.IdToken) *RequestStartTransactionRequest {
	return &RequestStartTransactionRequest{RemoteStartID: remoteStartID, IDToken: IdToken}
}

// Creates a new RequestStartTransactionResponse, containing all required fields. Optional fields may be set afterwards.
func NewRequestStartTransactionResponse(status RequestStartStopStatus) *RequestStartTransactionResponse {
	return &RequestStartTransactionResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("requestStartStopStatus", isValidRequestStartStopStatus)
}
