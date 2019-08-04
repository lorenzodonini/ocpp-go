package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Change Availability --------------------

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

type ChangeAvailabilityRequest struct {
	ConnectorId int              `json:"connectorId" validate:"gte=0"`
	Type        AvailabilityType `json:"type" validate:"required,availabilityType"`
}

type ChangeAvailabilityConfirmation struct {
	Status AvailabilityStatus `json:"status" validate:"required,availabilityStatus"`
}

type ChangeAvailabilityFeature struct{}

func (f ChangeAvailabilityFeature) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (f ChangeAvailabilityFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityRequest{})
}

func (f ChangeAvailabilityFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityConfirmation{})
}

func (r ChangeAvailabilityRequest) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (c ChangeAvailabilityConfirmation) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func NewChangeAvailabilityRequest(connectorId int, availabilityType AvailabilityType) *ChangeAvailabilityRequest {
	return &ChangeAvailabilityRequest{ConnectorId: connectorId, Type: availabilityType}
}

func NewChangeAvailabilityConfirmation(status AvailabilityStatus) *ChangeAvailabilityConfirmation {
	return &ChangeAvailabilityConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("availabilityType", isValidAvailabilityType)
	_ = Validate.RegisterValidation("availabilityStatus", isValidAvailabilityStatus)
}
