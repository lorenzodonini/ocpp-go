package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Boot Notification (CP -> CS) --------------------

// Result of registration in response to a BootNotification request.
type RegistrationStatus string

const (
	RegistrationStatusAccepted RegistrationStatus = "Accepted"
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusRejected RegistrationStatus = "Rejected"
)

func isValidRegistrationStatus(fl validator.FieldLevel) bool {
	status := RegistrationStatus(fl.Field().String())
	switch status {
	case RegistrationStatusAccepted, RegistrationStatusPending, RegistrationStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the BootNotification request payload sent by the Charge Point to the Central System.
type BootNotificationRequest struct {
	ChargeBoxSerialNumber   string `json:"chargeBoxSerialNumber,omitempty" validate:"max=25"`
	ChargePointModel        string `json:"chargePointModel" validate:"required,max=20"`
	ChargePointSerialNumber string `json:"chargePointSerialNumber,omitempty" validate:"max=25"`
	ChargePointVendor       string `json:"chargePointVendor" validate:"required,max=20"`
	FirmwareVersion         string `json:"firmwareVersion,omitempty" validate:"max=50"`
	Iccid                   string `json:"iccid,omitempty" validate:"max=20"`
	Imsi                    string `json:"imsi,omitempty" validate:"max=20"`
	MeterSerialNumber       string `json:"meterSerialNumber,omitempty" validate:"max=25"`
	MeterType               string `json:"meterType,omitempty" validate:"max=25"`
}

// This field definition of the BootNotification confirmation payload, sent by the Central System to the Charge Point in response to a BootNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type BootNotificationConfirmation struct {
	CurrentTime *DateTime          `json:"currentTime" validate:"required"`
	Interval    int                `json:"interval" validate:"gte=0"`
	Status      RegistrationStatus `json:"status" validate:"required,registrationStatus"`
}

// After each (re)boot, a Charge Point SHALL send a request to the Central System with information about its configuration (e.g. version, vendor, etc.).
// The Central System SHALL respond to indicate whether it will accept the Charge Point.
type BootNotificationFeature struct{}

func (f BootNotificationFeature) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (f BootNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(BootNotificationRequest{})
}

func (f BootNotificationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(BootNotificationConfirmation{})
}

func (r BootNotificationRequest) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (c BootNotificationConfirmation) GetFeatureName() string {
	return BootNotificationFeatureName
}

// Creates a new BootNotificationRequest, containing all required fields. Optional fields may be set afterwards.
func NewBootNotificationRequest(chargePointModel string, chargePointVendor string) *BootNotificationRequest {
	return &BootNotificationRequest{ChargePointModel: chargePointModel, ChargePointVendor: chargePointVendor}
}

// Creates a new BootNotificationConfirmation. Optional fields may be set afterwards.
func NewBootNotificationConfirmation(currentTime *DateTime, interval int, status RegistrationStatus) *BootNotificationConfirmation {
	return &BootNotificationConfirmation{CurrentTime: currentTime, Interval: interval, Status: status}
}

func init() {
	_ = Validate.RegisterValidation("registrationStatus", isValidRegistrationStatus)
}
