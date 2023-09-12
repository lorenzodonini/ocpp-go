package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Boot Notification (CP -> CS) --------------------

const BootNotificationFeatureName = "BootNotification"

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
	CurrentTime *types.DateTime    `json:"currentTime" validate:"required"`
	Interval    int                `json:"interval" validate:"gte=0"`
	Status      RegistrationStatus `json:"status" validate:"required,registrationStatus16"`
}

// After each (re)boot, a Charge Point SHALL send a request to the Central System with information about its configuration (e.g. version, vendor, etc.).
// The Central System SHALL respond to indicate whether it will accept the Charge Point.
// Between the physical power-on/reboot and the successful completion of a BootNotification, where Central System returns Accepted or Pending, the Charge Point SHALL NOT send any other request to the Central System.
// When the Central System responds with a BootNotification.conf with a status Accepted, the Charge Point will adjust the heartbeat
// interval in accordance with the interval from the response PDU and it is RECOMMENDED to synchronize its internal clock with the supplied Central System’s current time.
// If that interval value is zero, the Charge Point chooses a waiting interval on its own, in a way that avoids flooding the Central System with requests.
// If the Central System returns the Pending status, the communication channel SHOULD NOT be closed by either the Charge Point or the Central System.
// The Central System MAY send request messages to retrieve information from the Charge Point or change its configuration.
type BootNotificationFeature struct{}

func (f BootNotificationFeature) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (f BootNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(BootNotificationRequest{})
}

func (f BootNotificationFeature) GetResponseType() reflect.Type {
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

// Creates a new BootNotificationConfirmation. There are no optional fields for this message.
func NewBootNotificationConfirmation(currentTime *types.DateTime, interval int, status RegistrationStatus) *BootNotificationConfirmation {
	return &BootNotificationConfirmation{CurrentTime: currentTime, Interval: interval, Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("registrationStatus16", isValidRegistrationStatus)
}
