package core

import (
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"reflect"
)

// -------------------- Change Availability --------------------

type AvailabilityType string

const (
	AvailabilityTypeOperative = "Operative"
	AvailabilityTypeInoperative = "Inoperative"
)

type AvailabilityStatus string

const (
	AvailabilityStatusAccepted = "Accepted"
	AvailabilityStatusRejected = "Rejected"
	AvailabilityStatusScheduled = "Scheduled"
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
	return v16.ChangeAvailabilityFeatureName
}

func (f ChangeAvailabilityFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityRequest{})
}

func (f ChangeAvailabilityFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityConfirmation{})
}

func (r ChangeAvailabilityRequest) GetFeatureName() string {
	return v16.ChangeAvailabilityFeatureName
}

func (c ChangeAvailabilityConfirmation) GetFeatureName() string {
	return v16.ChangeAvailabilityFeatureName
}
