package display

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Clear Display Message (CSMS -> CS) --------------------

const ClearDisplayMessageFeatureName = "ClearDisplayMessage"

// Status returned in response to ClearDisplayRequest.
type ClearMessageStatus string

const (
	ClearMessageStatusAccepted ClearMessageStatus = "Accepted"
	ClearMessageStatusUnknown  ClearMessageStatus = "Unknown"
)

func isValidClearMessageStatus(fl validator.FieldLevel) bool {
	status := ClearMessageStatus(fl.Field().String())
	switch status {
	case ClearMessageStatusAccepted, ClearMessageStatusUnknown:
		return true
	default:
		return false
	}
}

// The field definition of the ClearDisplay request payload sent by the CSMS to the Charging Station.
type ClearDisplayRequest struct {
	ID int `json:"id"` // Id of the message that SHALL be removed from the Charging Station.
}

// This field definition of the ClearDisplay response payload, sent by the Charging Station to the CSMS in response to a ClearDisplayRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearDisplayResponse struct {
	Status     ClearMessageStatus `json:"status" validate:"required,clearMessageStatus21"`
	StatusInfo *types.StatusInfo  `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSMS asks the Charging Station to clear a display message that has been configured in the Charging Station to be cleared/removed.
// The Charging station checks for a message with the requested ID and removes it.
// The Charging station then responds with a ClearDisplayResponse. The response payload indicates whether the Charging Station was able to remove the message from display or not.
type ClearDisplayFeature struct{}

func (f ClearDisplayFeature) GetFeatureName() string {
	return ClearDisplayMessageFeatureName
}

func (f ClearDisplayFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearDisplayRequest{})
}

func (f ClearDisplayFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearDisplayResponse{})
}

func (r ClearDisplayRequest) GetFeatureName() string {
	return ClearDisplayMessageFeatureName
}

func (c ClearDisplayResponse) GetFeatureName() string {
	return ClearDisplayMessageFeatureName
}

// Creates a new ClearDisplayRequest, containing all required fields. There are no optional fields for this message.
func NewClearDisplayRequest(id int) *ClearDisplayRequest {
	return &ClearDisplayRequest{ID: id}
}

// Creates a new ClearDisplayResponse, containing all required fields. Optional fields may be set afterwards.
func NewClearDisplayResponse(status ClearMessageStatus) *ClearDisplayResponse {
	return &ClearDisplayResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("clearMessageStatus21", isValidClearMessageStatus)
}
