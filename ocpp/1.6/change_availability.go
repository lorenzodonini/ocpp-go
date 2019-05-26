package v16

import (
	"reflect"
)

// -------------------- Change Availability --------------------

type AvailabilityType string

const (
	AvailabilityTypeOperative AvailabilityType = "Operative"
	AvailabilityTypeInoperative AvailabilityType = "Inoperative"
)

type AvailabilityStatus string

const (
	AvailabilityStatusAccepted AvailabilityStatus = "Accepted"
	AvailabilityStatusRejected AvailabilityStatus = "Rejected"
	AvailabilityStatusScheduled AvailabilityStatus = "Scheduled"
)

type ChangeAvailabilityRequest struct {
	ConnectorId int				`json:"connectorId" validate:"required,gt=0"`
	Type AvailabilityType		`json:"type" validate:"required"`
}

type ChangeAvailabilityConfirmation struct {
	Status AvailabilityStatus	`json:"status" validate:"required"`
}

type ChangeAvailabilityFeature struct {}

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
