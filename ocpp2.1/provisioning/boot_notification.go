package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Boot Notification (CS -> CSMS) --------------------

const BootNotificationFeatureName = "BootNotification"

// RegistrationStatus represents the result of a registration request.
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

// BootReason represents the reason for sending a BootNotification.
type BootReason string

const (
	BootReasonApplicationReset BootReason = "ApplicationReset"
	BootReasonFirmwareUpdate   BootReason = "FirmwareUpdate"
	BootReasonLocalReset       BootReason = "LocalReset"
	BootReasonPowerUp          BootReason = "PowerUp"
	BootReasonRemoteReset      BootReason = "RemoteReset"
	BootReasonScheduledReset   BootReason = "ScheduledReset"
	BootReasonTriggered        BootReason = "Triggered"
	BootReasonUnknown          BootReason = "Unknown"
	BootReasonWatchdog         BootReason = "Watchdog"
)

func isValidBootReason(fl validator.FieldLevel) bool {
	reason := BootReason(fl.Field().String())
	switch reason {
	case BootReasonApplicationReset, BootReasonFirmwareUpdate, BootReasonLocalReset,
		BootReasonPowerUp, BootReasonRemoteReset, BootReasonScheduledReset,
		BootReasonTriggered, BootReasonUnknown, BootReasonWatchdog:
		return true
	default:
		return false
	}
}

// ModemType defines parameters for wireless communication.
type ModemType struct {
	Iccid      string                `json:"iccid,omitempty" validate:"omitempty,max=20"`
	Imsi       string                `json:"imsi,omitempty" validate:"omitempty,max=20"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// ChargingStationType represents the physical system where an EV can be charged.
type ChargingStationType struct {
	SerialNumber    string                `json:"serialNumber,omitempty" validate:"omitempty,max=25"`
	Model           string                `json:"model" validate:"required,max=20"`
	VendorName      string                `json:"vendorName" validate:"required,max=50"`
	FirmwareVersion string                `json:"firmwareVersion,omitempty" validate:"omitempty,max=50"`
	Modem           *ModemType            `json:"modem,omitempty" validate:"omitempty"`
	CustomData      *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// BootNotificationRequest is the payload for a BootNotification request from CS to CSMS.
type BootNotificationRequest struct {
	Reason          BootReason            `json:"reason" validate:"required,bootReason21"`
	ChargingStation ChargingStationType   `json:"chargingStation" validate:"required,dive"`
	CustomData      *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// BootNotificationResponse is the payload for a BootNotification response from CSMS to CS.
type BootNotificationResponse struct {
	CurrentTime *types.DateTime       `json:"currentTime" validate:"required"`
	Interval    int                   `json:"interval" validate:"gte=0"`
	Status      RegistrationStatus    `json:"status" validate:"required,registrationStatus21"`
	StatusInfo  *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData  *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// BootNotificationFeature represents the BootNotification feature.
type BootNotificationFeature struct{}

func (f BootNotificationFeature) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (f BootNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(BootNotificationRequest{})
}

func (f BootNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(BootNotificationResponse{})
}

func (r BootNotificationRequest) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (c BootNotificationResponse) GetFeatureName() string {
	return BootNotificationFeatureName
}

// NewBootNotificationRequest creates a new BootNotificationRequest.
func NewBootNotificationRequest(reason BootReason, model string, vendorName string) *BootNotificationRequest {
	return &BootNotificationRequest{
		Reason: reason,
		ChargingStation: ChargingStationType{
			Model:      model,
			VendorName: vendorName,
		},
	}
}

// NewBootNotificationResponse creates a new BootNotificationResponse.
func NewBootNotificationResponse(currentTime *types.DateTime, interval int, status RegistrationStatus) *BootNotificationResponse {
	return &BootNotificationResponse{
		CurrentTime: currentTime,
		Interval:    interval,
		Status:      status,
	}
}

func init() {
	_ = types.Validate.RegisterValidation("registrationStatus21", isValidRegistrationStatus)
	_ = types.Validate.RegisterValidation("bootReason21", isValidBootReason)
}
