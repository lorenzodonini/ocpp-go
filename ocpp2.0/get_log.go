package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Log (CSMS -> CS) --------------------

// LogType represents the type of log file that the Charging Station should send. It is used in GetLogRequest.
type LogType string

// LogStatus represents the status returned by a Charging Station in a GetLogConfirmation.
type LogStatus string

const (
	LogTypeDiagnostics        LogType   = "DiagnosticsLog"   // This contains the field definition of a diagnostics log file
	LogTypeSecurity           LogType   = "SecurityLog"      // Sent by the CSMS to the Charging Station to request that the Charging Station uploads the security log
	LogStatusAccepted         LogStatus = "Accepted"         // Accepted this log upload. This does not mean the log file is uploaded is successfully, the Charging Station will now start the log file upload.
	LogStatusRejected         LogStatus = "Rejected"         // Log update request rejected.
	LogStatusAcceptedCanceled LogStatus = "AcceptedCanceled" // Accepted this log upload, but in doing this has canceled an ongoing log file upload.
)

func isValidLogType(fl validator.FieldLevel) bool {
	status := LogType(fl.Field().String())
	switch status {
	case LogTypeDiagnostics, LogTypeSecurity:
		return true
	default:
		return false
	}
}

func isValidLogStatus(fl validator.FieldLevel) bool {
	status := LogStatus(fl.Field().String())
	switch status {
	case LogStatusAccepted, LogStatusRejected, LogStatusAcceptedCanceled:
		return true
	default:
		return false
	}
}

// LogParameters specifies the requested log and the location to which the log should be sent. It is used in GetLogRequest.
type LogParameters struct {
	RemoteLocation  string    `json:"remoteLocation" validate:"required,max=512,url"`
	OldestTimestamp *DateTime `json:"oldestTimestamp,omitempty" validate:"omitempty"`
	LatestTimestamp *DateTime `json:"latestTimestamp,omitempty" validate:"omitempty"`
}

// The field definition of the GetLog request payload sent by the CSMS to the Charging Station.
type GetLogRequest struct {
	LogType       LogType       `json:"logType" validate:"required,logType"`
	RequestID     int           `json:"requestId" validate:"gte=0"`
	Retries       *int          `json:"retries,omitempty" validate:"omitempty,gte=0"`
	RetryInterval *int          `json:"retryInterval,omitempty" validate:"omitempty,gte=0"`
	Log           LogParameters `json:"log" validate:"required"`
}

// This field definition of the GetLog confirmation payload, sent by the Charging Station to the CSMS in response to a GetLogRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetLogConfirmation struct {
	Status   LogStatus `json:"status" validate:"required,logStatus"`            // This field indicates whether the Charging Station was able to accept the request.
	Filename string    `json:"filename,omitempty" validate:"omitempty,max=256"` // This contains the name of the log file that will be uploaded. This field is not present when no logging information is available.
}

// The CSO may trigger the CSMS to request a report from a Charging Station.
// The CSMS shall then request a Charging Station to send a predefined report as defined in ReportBase.
// The Charging Station responds with GetLogConfirmation.
// The result will be returned asynchronously in one or more NotifyReportRequest messages (one for each report part).
type GetLogFeature struct{}

func (f GetLogFeature) GetFeatureName() string {
	return GetLogFeatureName
}

func (f GetLogFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetLogRequest{})
}

func (f GetLogFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetLogConfirmation{})
}

func (r GetLogRequest) GetFeatureName() string {
	return GetLogFeatureName
}

func (c GetLogConfirmation) GetFeatureName() string {
	return GetLogFeatureName
}

// Creates a new GetLogRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetLogRequest(logType LogType, requestID int, logParameters LogParameters) *GetLogRequest {
	return &GetLogRequest{LogType: logType, RequestID: requestID, Log: logParameters}
}

// Creates a new GetLogConfirmation, containing all required fields. Optional fields may be set afterwards.
func NewGetLogConfirmation(status LogStatus) *GetLogConfirmation {
	return &GetLogConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("logType", isValidLogType)
	_ = Validate.RegisterValidation("logStatus", isValidLogStatus)
}
