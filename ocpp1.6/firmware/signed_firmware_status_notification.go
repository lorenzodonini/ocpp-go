package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Signed Firmware Status Notification (CP -> CS) --------------------

const SignedFirmwareStatusNotificationFeatureName = "SignedFirmwareStatusNotification"

// Status reported in SignedFirmwareStatusNotificationRequest.
type SignedFirmwareStatus string

const (
	SignedFirmwareStatusDownloaded                SignedFirmwareStatus = "Downloaded"
	SignedFirmwareStatusDownloadFailed            SignedFirmwareStatus = "DownloadFailed"
	SignedFirmwareStatusDownloading               SignedFirmwareStatus = "Downloading"
	SignedFirmwareStatusDownloadScheduled         SignedFirmwareStatus = "DownloadScheduled"
	SignedFirmwareStatusDownloadPaused            SignedFirmwareStatus = "DownloadPaused"
	SignedFirmwareStatusIdle                      SignedFirmwareStatus = "Idle"
	SignedFirmwareStatusInstallationFailed        SignedFirmwareStatus = "InstallationFailed"
	SignedFirmwareStatusInstalling                SignedFirmwareStatus = "Installing"
	SignedFirmwareStatusInstalled                 SignedFirmwareStatus = "Installed"
	SignedFirmwareStatusInstallRebooting          SignedFirmwareStatus = "InstallRebooting"
	SignedFirmwareStatusInstallScheduled          SignedFirmwareStatus = "InstallScheduled"
	SignedFirmwareStatusInstallVerificationFailed SignedFirmwareStatus = "InstallVerificationFailed"
	SignedFirmwareStatusInvalidSignature          SignedFirmwareStatus = "InvalidSignature"
	SignedFirmwareStatusSignatureVerified         SignedFirmwareStatus = "SignatureVerified"
)

func isValidSignedFirmwareStatus(fl validator.FieldLevel) bool {
	status := SignedFirmwareStatus(fl.Field().String())
	switch status {
	case SignedFirmwareStatusDownloaded, SignedFirmwareStatusDownloadFailed, SignedFirmwareStatusDownloading, SignedFirmwareStatusDownloadScheduled,
		SignedFirmwareStatusDownloadPaused, SignedFirmwareStatusIdle, SignedFirmwareStatusInstallationFailed, SignedFirmwareStatusInstalling,
		SignedFirmwareStatusInstalled, SignedFirmwareStatusInstallRebooting, SignedFirmwareStatusInstallScheduled, SignedFirmwareStatusInstallVerificationFailed,
		SignedFirmwareStatusInvalidSignature, SignedFirmwareStatusSignatureVerified:
		return true
	default:
		return false
	}
}

type SignedFirmwareStatusNotificationRequest struct {
	Status    SignedFirmwareStatus `json:"status"`
	RequestId int                  `json:"requestId"`
}

// This field definition of the SignedFirmwareStatusNotification confirmation payload, sent by the Central System to the Charge Point in response to a SignedFirmwareStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SignedFirmwareStatusNotificationConfirmation struct {
}

type SignedFirmwareStatusNotificationFeature struct{}

func (f SignedFirmwareStatusNotificationFeature) GetFeatureName() string {
	return SignedFirmwareStatusNotificationFeatureName
}

func (f SignedFirmwareStatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SignedFirmwareStatusNotificationRequest{})
}

func (f SignedFirmwareStatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SignedFirmwareStatusNotificationRequest{})
}

func (r SignedFirmwareStatusNotificationRequest) GetFeatureName() string {
	return SignedFirmwareStatusNotificationFeatureName
}

func (c SignedFirmwareStatusNotificationConfirmation) GetFeatureName() string {
	return SignedFirmwareStatusNotificationFeatureName
}

// Creates a new SignedFirmwareStatusNotificationRequest, containing all required fields.
func NewSignedFirmwareStatusNotificationRequest(status SignedFirmwareStatus, requestId int) *SignedFirmwareStatusNotificationRequest {
	return &SignedFirmwareStatusNotificationRequest{
		Status:    status,
		RequestId: requestId,
	}
}

// Creates a new FirmwareStatusNotificationConfirmation, which doesn't contain any required or optional fields.
func NewSignedFirmwareStatusNotificationConfirmation() *SignedFirmwareStatusNotificationConfirmation {
	return &SignedFirmwareStatusNotificationConfirmation{}
}

func init() {
	_ = types.Validate.RegisterValidation("signedFirmwareStatus", isValidSignedFirmwareStatus)
}
