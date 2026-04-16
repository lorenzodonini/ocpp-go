package availability

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Change Availability (CSMS -> CS) --------------------

const ChangeAvailabilityFeatureName = "ChangeAvailability"

// OperationalStatus represents the requested operational status.
type OperationalStatus string

const (
	OperationalStatusInoperative OperationalStatus = "Inoperative"
	OperationalStatusOperative   OperationalStatus = "Operative"
)

func isValidOperationalStatus(fl validator.FieldLevel) bool {
	s := OperationalStatus(fl.Field().String())
	switch s {
	case OperationalStatusInoperative, OperationalStatusOperative:
		return true
	default:
		return false
	}
}

// ChangeAvailabilityStatus represents the response status.
type ChangeAvailabilityStatus string

const (
	ChangeAvailabilityStatusAccepted  ChangeAvailabilityStatus = "Accepted"
	ChangeAvailabilityStatusRejected  ChangeAvailabilityStatus = "Rejected"
	ChangeAvailabilityStatusScheduled ChangeAvailabilityStatus = "Scheduled"
)

func isValidChangeAvailabilityStatus(fl validator.FieldLevel) bool {
	s := ChangeAvailabilityStatus(fl.Field().String())
	switch s {
	case ChangeAvailabilityStatusAccepted, ChangeAvailabilityStatusRejected, ChangeAvailabilityStatusScheduled:
		return true
	default:
		return false
	}
}

// ChangeAvailabilityRequest is the payload for a ChangeAvailability request from CSMS to CS.
type ChangeAvailabilityRequest struct {
	OperationalStatus OperationalStatus     `json:"operationalStatus" validate:"required,operationalStatus21"`
	EVSE              *types.EVSE           `json:"evse,omitempty" validate:"omitempty"`
	CustomData        *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ChangeAvailabilityResponse is the payload for a ChangeAvailability response from CS to CSMS.
type ChangeAvailabilityResponse struct {
	Status     ChangeAvailabilityStatus `json:"status" validate:"required,changeAvailabilityStatus21"`
	StatusInfo *types.StatusInfo        `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType    `json:"customData,omitempty" validate:"omitempty"`
}

// ChangeAvailabilityFeature represents the ChangeAvailability feature.
type ChangeAvailabilityFeature struct{}

func (f ChangeAvailabilityFeature) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (f ChangeAvailabilityFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityRequest{})
}

func (f ChangeAvailabilityFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityResponse{})
}

func (r ChangeAvailabilityRequest) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (c ChangeAvailabilityResponse) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

// NewChangeAvailabilityRequest creates a new ChangeAvailabilityRequest.
func NewChangeAvailabilityRequest(operationalStatus OperationalStatus) *ChangeAvailabilityRequest {
	return &ChangeAvailabilityRequest{OperationalStatus: operationalStatus}
}

// NewChangeAvailabilityResponse creates a new ChangeAvailabilityResponse.
func NewChangeAvailabilityResponse(status ChangeAvailabilityStatus) *ChangeAvailabilityResponse {
	return &ChangeAvailabilityResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("operationalStatus21", isValidOperationalStatus)
	_ = types.Validate.RegisterValidation("changeAvailabilityStatus21", isValidChangeAvailabilityStatus)
}
