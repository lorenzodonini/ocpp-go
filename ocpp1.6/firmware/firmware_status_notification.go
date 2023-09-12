package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Firmware Status Notification (CP -> CS) --------------------

const FirmwareStatusNotificationFeatureName = "FirmwareStatusNotification"

// Status reported in FirmwareStatusNotificationRequest.
type FirmwareStatus string

const (
	FirmwareStatusDownloaded         FirmwareStatus = "Downloaded"
	FirmwareStatusDownloadFailed     FirmwareStatus = "DownloadFailed"
	FirmwareStatusDownloading        FirmwareStatus = "Downloading"
	FirmwareStatusIdle               FirmwareStatus = "Idle"
	FirmwareStatusInstallationFailed FirmwareStatus = "InstallationFailed"
	FirmwareStatusInstalling         FirmwareStatus = "Installing"
	FirmwareStatusInstalled          FirmwareStatus = "Installed"
)

func isValidFirmwareStatus(fl validator.FieldLevel) bool {
	status := FirmwareStatus(fl.Field().String())
	switch status {
	case FirmwareStatusDownloaded, FirmwareStatusDownloadFailed, FirmwareStatusDownloading, FirmwareStatusIdle, FirmwareStatusInstallationFailed, FirmwareStatusInstalling, FirmwareStatusInstalled:
		return true
	default:
		return false
	}
}

// The field definition of the FirmwareStatusNotification request payload sent by the Charge Point to the Central System.
type FirmwareStatusNotificationRequest struct {
	Status FirmwareStatus `json:"status" validate:"required,firmwareStatus16"`
}

// This field definition of the FirmwareStatusNotification confirmation payload, sent by the Central System to the Charge Point in response to a FirmwareStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type FirmwareStatusNotificationConfirmation struct {
}

// The Charge Point sends a notification to inform the Central System about the progress of the downloading and installation of a firmware update.
// The Charge Point SHALL only send the status Idle after receipt of a TriggerMessage for a Firmware Status Notification, when it is not busy downloading/installing firmware.
// The FirmwareStatusNotification requests SHALL be sent to keep the Central System updated with the status of the update process.
type FirmwareStatusNotificationFeature struct{}

func (f FirmwareStatusNotificationFeature) GetFeatureName() string {
	return FirmwareStatusNotificationFeatureName
}

func (f FirmwareStatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(FirmwareStatusNotificationRequest{})
}

func (f FirmwareStatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(FirmwareStatusNotificationConfirmation{})
}

func (r FirmwareStatusNotificationRequest) GetFeatureName() string {
	return FirmwareStatusNotificationFeatureName
}

func (c FirmwareStatusNotificationConfirmation) GetFeatureName() string {
	return FirmwareStatusNotificationFeatureName
}

// Creates a new FirmwareStatusNotificationRequest, containing all required fields.
func NewFirmwareStatusNotificationRequest(status FirmwareStatus) *FirmwareStatusNotificationRequest {
	return &FirmwareStatusNotificationRequest{Status: status}
}

// Creates a new FirmwareStatusNotificationConfirmation, which doesn't contain any required or optional fields.
func NewFirmwareStatusNotificationConfirmation() *FirmwareStatusNotificationConfirmation {
	return &FirmwareStatusNotificationConfirmation{}
}

func init() {
	_ = types.Validate.RegisterValidation("firmwareStatus16", isValidFirmwareStatus)
}
