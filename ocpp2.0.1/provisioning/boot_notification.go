package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// -------------------- Boot Notification (CS -> CSMS) --------------------

const BootNotificationFeatureName = "BootNotification"

// Result of registration in response to a BootNotification request.
type RegistrationStatus string

// The reason for sending a BootNotification event to the CSMS.
type BootReason string

const (
	RegistrationStatusAccepted RegistrationStatus = "Accepted"
	RegistrationStatusPending  RegistrationStatus = "Pending"
	RegistrationStatusRejected RegistrationStatus = "Rejected"
	BootReasonApplicationReset BootReason         = "ApplicationReset"
	BootReasonFirmwareUpdate   BootReason         = "FirmwareUpdate"
	BootReasonLocalReset       BootReason         = "LocalReset"
	BootReasonPowerUp          BootReason         = "PowerUp"
	BootReasonRemoteReset      BootReason         = "RemoteReset"
	BootReasonScheduledReset   BootReason         = "ScheduledReset"
	BootReasonTriggered        BootReason         = "Triggered"
	BootReasonUnknown          BootReason         = "Unknown"
	BootReasonWatchdog         BootReason         = "Watchdog"
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

func isValidBootReason(fl validator.FieldLevel) bool {
	reason := BootReason(fl.Field().String())
	switch reason {
	case BootReasonApplicationReset, BootReasonFirmwareUpdate, BootReasonLocalReset, BootReasonPowerUp, BootReasonRemoteReset, BootReasonScheduledReset, BootReasonTriggered, BootReasonUnknown, BootReasonWatchdog:
		return true
	default:
		return false
	}
}

// Defines parameters required for initiating and maintaining wireless communication with other devices.
type ModemType struct {
	Iccid string `json:"iccid,omitempty" validate:"max=20"`
	Imsi  string `json:"imsi,omitempty" validate:"max=20"`
}

// The physical system where an Electrical Vehicle (EV) can be charged.
type ChargingStationType struct {
	SerialNumber    string     `json:"serialNumber,omitempty" validate:"max=25"`
	Model           string     `json:"model" validate:"required,max=20"`
	VendorName      string     `json:"vendorName" validate:"required,max=50"`
	FirmwareVersion string     `json:"firmwareVersion,omitempty" validate:"max=50"`
	Modem           *ModemType `json:"modem,omitempty"`
}

// The field definition of the BootNotification request payload sent by the Charging Station to the CSMS.
type BootNotificationRequest struct {
	Reason          BootReason          `json:"reason" validate:"required,bootReason"`
	ChargingStation ChargingStationType `json:"chargingStation" validate:"required,dive"`
}

// The field definition of the BootNotification response payload, sent by the CSMS to the Charging Station in response to a BootNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type BootNotificationResponse struct {
	CurrentTime *types.DateTime    `json:"currentTime" validate:"required"`
	Interval    int                `json:"interval" validate:"gte=0"`
	Status      RegistrationStatus `json:"status" validate:"required,registrationStatus"`
	StatusInfo  *types.StatusInfo  `json:"statusInfo,omitempty" validate:"omitempty"`
}

// After each (re)boot, a Charging Station SHALL send a request to the CSMS with information about its configuration (e.g. version, vendor, etc.).
// The CSMS SHALL respond to indicate whether it will accept the Charging Station.
// Between the physical power-on/reboot and the successful completion of a BootNotification, where CSMS returns Accepted or Pending, the Charging Station SHALL NOT send any other request to the CSMS.
//
// When the CSMS responds with a BootNotificationResponse with a status Accepted, the Charging Station will adjust the heartbeat
// interval in accordance with the interval from the response PDU and it is RECOMMENDED to synchronize its internal clock with the supplied CSMS’s current time.
//
// If that interval value is zero, the Charging Station chooses a waiting interval on its own, in a way that avoids flooding the CSMS with requests.
// If the CSMS returns the Pending status, the communication channel SHOULD NOT be closed by either the Charging Station or the CSMS.
//
// The CSMS MAY send request messages to retrieve information from the Charging Station or change its configuration.
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

// Creates a new BootNotificationRequest, containing all required fields. Optional fields may be set afterwards.
func NewBootNotificationRequest(reason BootReason, model string, vendorName string) *BootNotificationRequest {
	return &BootNotificationRequest{Reason: reason, ChargingStation: ChargingStationType{Model: model, VendorName: vendorName}}
}

// Creates a new BootNotificationResponse. There are no optional fields for this message.
func NewBootNotificationResponse(currentTime *types.DateTime, interval int, status RegistrationStatus) *BootNotificationResponse {
	return &BootNotificationResponse{CurrentTime: currentTime, Interval: interval, Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("registrationStatus", isValidRegistrationStatus)
	_ = types.Validate.RegisterValidation("bootReason", isValidBootReason)
}
