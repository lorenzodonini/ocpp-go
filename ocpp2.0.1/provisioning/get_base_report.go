package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// -------------------- Get Base Report (CSMS -> CS) --------------------

const GetBaseReportFeatureName = "GetBaseReport"

// Requested availability change in GetBaseReportRequest.
type ReportBaseType string

const (
	ReportTypeConfigurationInventory ReportBaseType = "ConfigurationInventory"
	ReportTypeFullInventory          ReportBaseType = "FullInventory"
	ReportTypeSummaryInventory       ReportBaseType = "SummaryInventory"
)

func isValidReportBaseType(fl validator.FieldLevel) bool {
	status := ReportBaseType(fl.Field().String())
	switch status {
	case ReportTypeConfigurationInventory, ReportTypeFullInventory, ReportTypeSummaryInventory:
		return true
	default:
		return false
	}
}

// The field definition of the GetBaseReport request payload sent by the CSMS to the Charging Station.
type GetBaseReportRequest struct {
	RequestID  int            `json:"requestId"`
	ReportBase ReportBaseType `json:"reportBase" validate:"required,reportBaseType"`
}

// This field definition of the GetBaseReport response payload, sent by the Charging Station to the CSMS in response to a GetBaseReportRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetBaseReportResponse struct {
	Status     types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus"`
	StatusInfo *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSO may trigger the CSMS to request a report from a Charging Station.
// The CSMS shall then request a Charging Station to send a predefined report as defined in ReportBase.
// The Charging Station responds with GetBaseReportResponse.
// The result will be returned asynchronously in one or more NotifyReportRequest messages (one for each report part).
type GetBaseReportFeature struct{}

func (f GetBaseReportFeature) GetFeatureName() string {
	return GetBaseReportFeatureName
}

func (f GetBaseReportFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetBaseReportRequest{})
}

func (f GetBaseReportFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetBaseReportResponse{})
}

func (r GetBaseReportRequest) GetFeatureName() string {
	return GetBaseReportFeatureName
}

func (c GetBaseReportResponse) GetFeatureName() string {
	return GetBaseReportFeatureName
}

// Creates a new GetBaseReportRequest, containing all required fields. There are no optional fields for this message.
func NewGetBaseReportRequest(requestID int, reportBase ReportBaseType) *GetBaseReportRequest {
	return &GetBaseReportRequest{RequestID: requestID, ReportBase: reportBase}
}

// Creates a new GetBaseReportResponse, containing all required fields. There are no optional fields for this message.
func NewGetBaseReportResponse(status types.GenericDeviceModelStatus) *GetBaseReportResponse {
	return &GetBaseReportResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("reportBaseType", isValidReportBaseType)
}
