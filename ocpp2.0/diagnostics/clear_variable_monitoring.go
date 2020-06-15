package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Clear Variable Monitoring (CSMS -> CS) --------------------

const ClearVariableMonitoringFeatureName = "ClearVariableMonitoring"

// Status contained inside a ClearMonitoringResult struct.
type ClearMonitoringStatus string

const (
	ClearMonitoringStatusAccepted ClearMonitoringStatus = "Accepted"
	ClearMonitoringStatusRejected ClearMonitoringStatus = "Rejected"
	ClearMonitoringStatusNotFound ClearMonitoringStatus = "NotFound"
)

func isValidClearMonitoringStatus(fl validator.FieldLevel) bool {
	status := ClearMonitoringStatus(fl.Field().String())
	switch status {
	case ClearMonitoringStatusAccepted, ClearMonitoringStatusRejected, ClearMonitoringStatusNotFound:
		return true
	default:
		return false
	}
}

type ClearMonitoringResult struct {
	ID     int                   `json:"id" validate:"required,gte=0"`
	Status ClearMonitoringStatus `json:"status" validate:"required,clearMonitoringStatus"`
}

// The field definition of the ClearVariableMonitoring request payload sent by the CSMS to the Charging Station.
type ClearVariableMonitoringRequest struct {
	ID []int `json:"id" validate:"required,min=1,dive,gte=0"` // List of the monitors to be cleared, identified by their Id.
}

// This field definition of the ClearVariableMonitoring response payload, sent by the Charging Station to the CSMS in response to a ClearVariableMonitoringRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearVariableMonitoringResponse struct {
	ClearMonitoringResult []ClearMonitoringResult `json:"clearMonitoringResult" validate:"required,min=1,dive"` // List of result statuses per monitor.
}

// The CSMS asks the Charging Station to clear a display message that has been configured in the Charging Station to be cleared/removed.
// The Charging station checks for a message with the requested ID and removes it.
// The Charging station then responds with a ClearVariableMonitoringResponse. The response payload indicates whether the Charging Station was able to remove the message from display or not.
type ClearVariableMonitoringFeature struct{}

func (f ClearVariableMonitoringFeature) GetFeatureName() string {
	return ClearVariableMonitoringFeatureName
}

func (f ClearVariableMonitoringFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearVariableMonitoringRequest{})
}

func (f ClearVariableMonitoringFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearVariableMonitoringResponse{})
}

func (r ClearVariableMonitoringRequest) GetFeatureName() string {
	return ClearVariableMonitoringFeatureName
}

func (c ClearVariableMonitoringResponse) GetFeatureName() string {
	return ClearVariableMonitoringFeatureName
}

// Creates a new ClearVariableMonitoringRequest, containing all required fields. There are no optional fields for this message.
func NewClearVariableMonitoringRequest(id []int) *ClearVariableMonitoringRequest {
	return &ClearVariableMonitoringRequest{ID: id}
}

// Creates a new ClearVariableMonitoringResponse, containing all required fields. There are no optional fields for this message.
func NewClearVariableMonitoringResponse(result []ClearMonitoringResult) *ClearVariableMonitoringResponse {
	return &ClearVariableMonitoringResponse{ClearMonitoringResult: result}
}

func init() {
	_ = types.Validate.RegisterValidation("clearMonitoringStatus", isValidClearMonitoringStatus)
}
