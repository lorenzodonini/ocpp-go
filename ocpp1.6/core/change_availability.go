package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Change Availability (CS -> CP) --------------------

const ChangeAvailabilityFeatureName = "ChangeAvailability"

// Requested availability change in ChangeAvailabilityRequest.
type AvailabilityType string

const (
	AvailabilityTypeOperative   AvailabilityType = "Operative"
	AvailabilityTypeInoperative AvailabilityType = "Inoperative"
)

func isValidAvailabilityType(fl validator.FieldLevel) bool {
	status := AvailabilityType(fl.Field().String())
	switch status {
	case AvailabilityTypeOperative, AvailabilityTypeInoperative:
		return true
	default:
		return false
	}
}

// Status returned in response to ChangeAvailabilityRequest
type AvailabilityStatus string

const (
	AvailabilityStatusAccepted  AvailabilityStatus = "Accepted"
	AvailabilityStatusRejected  AvailabilityStatus = "Rejected"
	AvailabilityStatusScheduled AvailabilityStatus = "Scheduled"
)

func isValidAvailabilityStatus(fl validator.FieldLevel) bool {
	status := AvailabilityStatus(fl.Field().String())
	switch status {
	case AvailabilityStatusAccepted, AvailabilityStatusRejected, AvailabilityStatusScheduled:
		return true
	default:
		return false
	}
}

// The field definition of the ChangeAvailability request payload sent by the Central System to the Charge Point.
type ChangeAvailabilityRequest struct {
	ConnectorId int              `json:"connectorId" validate:"gte=0"`
	Type        AvailabilityType `json:"type" validate:"required,availabilityType"`
}

// This field definition of the ChangeAvailability confirmation payload, sent by the Charge Point to the Central System in response to a ChangeAvailabilityRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ChangeAvailabilityConfirmation struct {
	Status AvailabilityStatus `json:"status" validate:"required,availabilityStatus"`
}

// Central System can request a Charge Point to change its availability.
// A Charge Point is considered available (“operative”) when it is charging or ready for charging.
// A Charge Point is considered unavailable when it does not allow any charging.
// The Central System SHALL send a ChangeAvailabilityRequest for requesting a Charge Point to change its availability.
// The Central System can change the availability to available or unavailable.
type ChangeAvailabilityFeature struct{}

func (f ChangeAvailabilityFeature) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (f ChangeAvailabilityFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityRequest{})
}

func (f ChangeAvailabilityFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityConfirmation{})
}

func (r ChangeAvailabilityRequest) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (c ChangeAvailabilityConfirmation) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

// Creates a new ChangeAvailabilityRequest, containing all required fields. There are no optional fields for this message.
func NewChangeAvailabilityRequest(connectorId int, availabilityType AvailabilityType) *ChangeAvailabilityRequest {
	return &ChangeAvailabilityRequest{ConnectorId: connectorId, Type: availabilityType}
}

// Creates a new ChangeAvailabilityConfirmation, containing all required fields. There are no optional fields for this message.
func NewChangeAvailabilityConfirmation(status AvailabilityStatus) *ChangeAvailabilityConfirmation {
	return &ChangeAvailabilityConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("availabilityType", isValidAvailabilityType)
	_ = types.Validate.RegisterValidation("availabilityStatus", isValidAvailabilityStatus)
}
