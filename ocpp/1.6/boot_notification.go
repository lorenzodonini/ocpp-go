package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"reflect"
	"time"
)

// -------------------- Boot Notification --------------------
type RegistrationStatus string

const (
	RegistrationStatusAccepted RegistrationStatus = "Accepted"
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusRejected RegistrationStatus = "Rejected"
)

type BootNotificationRequest struct {
	ocpp.Request            `json:"-"`
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

//TODO: add custom validator for registration status & interval
type BootNotificationConfirmation struct {
	ocpp.Confirmation `json:"-"`
	CurrentTime       time.Time          `json:"currentTime" validate:"required"`
	Interval          int                `json:"interval" validate:"required,gte=0"`
	Status            RegistrationStatus `json:"status" validate:"required"`
}

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

func NewBootNotificationRequest(chargePointModel string, chargePointVendor string) *BootNotificationRequest {
	return &BootNotificationRequest{ChargePointModel: chargePointModel, ChargePointVendor: chargePointVendor}
}

func NewBootNotificationConfirmation(currentTime time.Time, interval int, status RegistrationStatus) *BootNotificationConfirmation {
	return &BootNotificationConfirmation{CurrentTime: currentTime, Interval: interval, Status: status}
}
