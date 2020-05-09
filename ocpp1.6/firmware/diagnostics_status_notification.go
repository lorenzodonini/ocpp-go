package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Diagnostics Status Notification (CP -> CS) --------------------

const DiagnosticsStatusNotificationFeatureName = "DiagnosticsStatusNotification"

// Status reported in DiagnosticsStatusNotificationRequest.
type DiagnosticsStatus string

const (
	DiagnosticsStatusIdle         DiagnosticsStatus = "Idle"
	DiagnosticsStatusUploaded     DiagnosticsStatus = "Uploaded"
	DiagnosticsStatusUploadFailed DiagnosticsStatus = "UploadFailed"
	DiagnosticsStatusUploading    DiagnosticsStatus = "Uploading"
)

func isValidDiagnosticsStatus(fl validator.FieldLevel) bool {
	status := DiagnosticsStatus(fl.Field().String())
	switch status {
	case DiagnosticsStatusIdle, DiagnosticsStatusUploaded, DiagnosticsStatusUploadFailed, DiagnosticsStatusUploading:
		return true
	default:
		return false
	}
}

// The field definition of the DiagnosticsStatusNotification request payload sent by the Charge Point to the Central System.
type DiagnosticsStatusNotificationRequest struct {
	Status DiagnosticsStatus `json:"status" validate:"required,diagnosticsStatus"`
}

// This field definition of the DiagnosticsStatusNotification confirmation payload, sent by the Central System to the Charge Point in response to a DiagnosticsStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type DiagnosticsStatusNotificationConfirmation struct {
}

// The Charge Point sends a notification to inform the Central System that the upload of diagnostics is busy or has finished successfully or failed.
// The Charge Point SHALL only send the status Idle after receipt of a TriggerMessage for a Diagnostics Status Notification, when it is not busy uploading diagnostics.
type DiagnosticsStatusNotificationFeature struct{}

func (f DiagnosticsStatusNotificationFeature) GetFeatureName() string {
	return DiagnosticsStatusNotificationFeatureName
}

func (f DiagnosticsStatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(DiagnosticsStatusNotificationRequest{})
}

func (f DiagnosticsStatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(DiagnosticsStatusNotificationConfirmation{})
}

func (r DiagnosticsStatusNotificationRequest) GetFeatureName() string {
	return DiagnosticsStatusNotificationFeatureName
}

func (c DiagnosticsStatusNotificationConfirmation) GetFeatureName() string {
	return DiagnosticsStatusNotificationFeatureName
}

// Creates a new DiagnosticsStatusNotificationRequest, containing all required fields.
func NewDiagnosticsStatusNotificationRequest(status DiagnosticsStatus) *DiagnosticsStatusNotificationRequest {
	return &DiagnosticsStatusNotificationRequest{Status: status}
}

// Creates a new DiagnosticsStatusNotificationConfirmation, which doesn't contain any required or optional fields.
func NewDiagnosticsStatusNotificationConfirmation() *DiagnosticsStatusNotificationConfirmation {
	return &DiagnosticsStatusNotificationConfirmation{}
}

func init() {
	_ = types.Validate.RegisterValidation("diagnosticsStatus", isValidDiagnosticsStatus)
}
