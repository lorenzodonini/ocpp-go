package securefirmware

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// -------------------- Firmware Status Notification (CS -> CSMS) --------------------

const SignedFirmwareStatusNotificationFeatureName = "SignedFirmwareStatusNotification"

// Status reported in SignedFirmwareStatusNotificationRequest.
type FirmwareStatus string

const (
	FirmwareStatusDownloaded                FirmwareStatus = "Downloaded"
	FirmwareStatusDownloadFailed            FirmwareStatus = "DownloadFailed"
	FirmwareStatusDownloading               FirmwareStatus = "Downloading"
	FirmwareStatusDownloadScheduled         FirmwareStatus = "DownloadScheduled"
	FirmwareStatusDownloadPaused            FirmwareStatus = "DownloadPaused"
	FirmwareStatusIdle                      FirmwareStatus = "Idle"
	FirmwareStatusInstallationFailed        FirmwareStatus = "InstallationFailed"
	FirmwareStatusInstalling                FirmwareStatus = "Installing"
	FirmwareStatusInstalled                 FirmwareStatus = "Installed"
	FirmwareStatusInstallRebooting          FirmwareStatus = "InstallRebooting"
	FirmwareStatusInstallScheduled          FirmwareStatus = "InstallScheduled"
	FirmwareStatusInstallVerificationFailed FirmwareStatus = "InstallVerificationFailed"
	FirmwareStatusInvalidSignature          FirmwareStatus = "InvalidSignature"
	FirmwareStatusSignatureVerified         FirmwareStatus = "SignatureVerified"
	FirmwareStatusCertificateVerified       FirmwareStatus = "CertificateVerified"
	FirmwareStatusInvalidCertificate        FirmwareStatus = "InvalidCertificate"
	FirmwareStatusRevokedCertificate        FirmwareStatus = "RevokedCertificate"
)

func isValidFirmwareStatus(fl validator.FieldLevel) bool {
	status := FirmwareStatus(fl.Field().String())
	switch status {
	case FirmwareStatusDownloaded,
		FirmwareStatusDownloadFailed,
		FirmwareStatusDownloading,
		FirmwareStatusDownloadScheduled,
		FirmwareStatusDownloadPaused,
		FirmwareStatusIdle,
		FirmwareStatusInstallationFailed,
		FirmwareStatusInstalling,
		FirmwareStatusInstalled,
		FirmwareStatusInstallRebooting,
		FirmwareStatusInstallScheduled,
		FirmwareStatusInstallVerificationFailed,
		FirmwareStatusInvalidSignature,
		FirmwareStatusSignatureVerified:
		return true
	default:
		return false
	}
}

// The field definition of the FirmwareStatusNotification request payload sent by the Charging Station to the CSMS.
type SignedFirmwareStatusNotificationRequest struct {
	Status    FirmwareStatus `json:"status" validate:"required,signedFirmwareStatus"`
	RequestID *int           `json:"requestId,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the FirmwareStatusNotification response payload, sent by the CSMS to the Charging Station in response to a SignedFirmwareStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SignedFirmwareStatusNotificationResponse struct {
}

// The Charging Station sends a notification to inform the CSMS about the progress of the downloading and installation of a firmware update.
// The Charging Station SHALL only send the status Idle after receipt of a TriggerMessage for a Firmware Status Notification, when it is not busy downloading/installing firmware.
// The FirmwareStatusNotification requests SHALL be sent to keep the CSMS updated with the status of the update process.
type SignedFirmwareStatusNotificationFeature struct{}

func (f SignedFirmwareStatusNotificationFeature) GetFeatureName() string {
	return SignedFirmwareStatusNotificationFeatureName
}

func (f SignedFirmwareStatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SignedFirmwareStatusNotificationRequest{})
}

func (f SignedFirmwareStatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SignedFirmwareStatusNotificationResponse{})
}

func (r SignedFirmwareStatusNotificationRequest) GetFeatureName() string {
	return SignedFirmwareStatusNotificationFeatureName
}

func (c SignedFirmwareStatusNotificationResponse) GetFeatureName() string {
	return SignedFirmwareStatusNotificationFeatureName
}

// Creates a new SignedFirmwareStatusNotificationRequest, containing all required fields. Optional fields may be set afterwards.
func NewFirmwareStatusNotificationRequest(status FirmwareStatus) *SignedFirmwareStatusNotificationRequest {
	return &SignedFirmwareStatusNotificationRequest{Status: status}
}

// Creates a new SignedFirmwareStatusNotificationResponse, which doesn't contain any required or optional fields.
func NewFirmwareStatusNotificationResponse() *SignedFirmwareStatusNotificationResponse {
	return &SignedFirmwareStatusNotificationResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("signedFirmwareStatus", isValidFirmwareStatus)
}
