package diagnostics

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Set Monitoring Base (CSMS -> CS) --------------------

const SetMonitoringBaseFeatureName = "SetMonitoringBase"

// Monitoring base to be set within the Charging Station.
type MonitoringBase string

const (
	MonitoringBaseAll            MonitoringBase = "All"
	MonitoringBaseFactoryDefault MonitoringBase = "FactoryDefault"
	MonitoringBaseHardWiredOnly  MonitoringBase = "HardWiredOnly"
)

func isValidMonitoringBase(fl validator.FieldLevel) bool {
	status := MonitoringBase(fl.Field().String())
	switch status {
	case MonitoringBaseAll, MonitoringBaseFactoryDefault, MonitoringBaseHardWiredOnly:
		return true
	default:
		return false
	}
}

// The field definition of the SetMonitoringBase request payload sent by the CSMS to the Charging Station.
type SetMonitoringBaseRequest struct {
	MonitoringBase MonitoringBase `json:"monitoringBase" validate:"required,monitoringBase21"` // Specifies which monitoring base will be set.
}

// This field definition of the SetMonitoringBase response payload, sent by the Charging Station to the CSMS in response to a SetMonitoringBaseRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetMonitoringBaseResponse struct {
	Status     types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus21"` // Indicates whether the Charging Station was able to accept the request.
	StatusInfo *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`             // Detailed status information.
}

// A CSMS has the ability to request the Charging Station to activate a set of preconfigured
// monitoring settings, as denoted by the value of MonitoringBase. This is achieved by sending a
// SetMonitoringBaseRequest to the charging station. The charging station will respond with a
// SetMonitoringBaseResponse message.
//
// It is up to the manufacturer of the Charging Station to define which monitoring settings are activated
// by All, FactoryDefault and HardWiredOnly.
type SetMonitoringBaseFeature struct{}

func (f SetMonitoringBaseFeature) GetFeatureName() string {
	return SetMonitoringBaseFeatureName
}

func (f SetMonitoringBaseFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetMonitoringBaseRequest{})
}

func (f SetMonitoringBaseFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetMonitoringBaseResponse{})
}

func (r SetMonitoringBaseRequest) GetFeatureName() string {
	return SetMonitoringBaseFeatureName
}

func (c SetMonitoringBaseResponse) GetFeatureName() string {
	return SetMonitoringBaseFeatureName
}

// Creates a new SetMonitoringBaseRequest, containing all required fields. There are no optional fields for this message.
func NewSetMonitoringBaseRequest(monitoringBase MonitoringBase) *SetMonitoringBaseRequest {
	return &SetMonitoringBaseRequest{MonitoringBase: monitoringBase}
}

// Creates a new SetMonitoringBaseResponse, containing all required fields. Optional fields may be set afterwards.
func NewSetMonitoringBaseResponse(status types.GenericDeviceModelStatus) *SetMonitoringBaseResponse {
	return &SetMonitoringBaseResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("monitoringBase21", isValidMonitoringBase)
}
