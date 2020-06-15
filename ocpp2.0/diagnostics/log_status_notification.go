package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Log Status Notification (CS -> CSMS) --------------------

const LogStatusNotificationFeatureName = "LogStatusNotification"

// UploadLogStatus represents the current status of the log-upload procedure, reported by a Charging Station in a LogStatusNotificationRequest.
type UploadLogStatus string

const (
	UploadLogStatusBadMessage       UploadLogStatus = "BadMessage"            // A badly formatted packet or other protocol incompatibility was detected.
	UploadLogStatusIdle             UploadLogStatus = "Idle"                  // The Charging Station is not uploading a log file. Idle SHALL only be used when the message was triggered by a TriggerMessageRequest.
	UploadLogStatusNotSupportedOp   UploadLogStatus = "NotSupportedOperation" // The server does not support the operation.
	UploadLogStatusPermissionDenied UploadLogStatus = "PermissionDenied"      // Insufficient permissions to perform the operation.
	UploadLogStatusUploaded         UploadLogStatus = "Uploaded"              // File has been uploaded successfully.
	UploadLogStatusUploadFailure    UploadLogStatus = "UploadFailure"         // Failed to upload the requested file.
	UploadLogStatusUploading        UploadLogStatus = "Uploading"             // File is being uploaded.
)

func isValidUploadLogStatus(fl validator.FieldLevel) bool {
	status := UploadLogStatus(fl.Field().String())
	switch status {
	case UploadLogStatusBadMessage, UploadLogStatusIdle, UploadLogStatusNotSupportedOp, UploadLogStatusPermissionDenied, UploadLogStatusUploaded, UploadLogStatusUploadFailure, UploadLogStatusUploading:
		return true
	default:
		return false
	}
}

// The field definition of the LogStatusNotification request payload sent by a Charging Station to the CSMS.
type LogStatusNotificationRequest struct {
	Status    UploadLogStatus `json:"status" validate:"required,uploadLogStatus"`
	RequestID int             `json:"requestId" validate:"gte=0"`
}

// This field definition of the LogStatusNotification response payload, sent by the CSMS to the Charging Station in response to a LogStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type LogStatusNotificationResponse struct {
}

// A Charging Station shall send LogStatusNotification requests to update the CSMS with the current status of a log-upload procedure.
// The CSMS shall respond with a LogStatusNotificationResponse acknowledging the status update request.
//
// After a successful log upload, the The Charging Station returns to Idle status.
type LogStatusNotificationFeature struct{}

func (f LogStatusNotificationFeature) GetFeatureName() string {
	return LogStatusNotificationFeatureName
}

func (f LogStatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(LogStatusNotificationRequest{})
}

func (f LogStatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(LogStatusNotificationResponse{})
}

func (r LogStatusNotificationRequest) GetFeatureName() string {
	return LogStatusNotificationFeatureName
}

func (c LogStatusNotificationResponse) GetFeatureName() string {
	return LogStatusNotificationFeatureName
}

// Creates a new LogStatusNotificationRequest, containing all required fields. There are no optional fields for this message.
func NewLogStatusNotificationRequest(status UploadLogStatus, requestID int) *LogStatusNotificationRequest {
	return &LogStatusNotificationRequest{Status: status, RequestID: requestID}
}

// Creates a new LogStatusNotificationResponse, which doesn't contain any required or optional fields.
func NewLogStatusNotificationResponse() *LogStatusNotificationResponse {
	return &LogStatusNotificationResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("uploadLogStatus", isValidUploadLogStatus)
}
