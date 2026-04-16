package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Base Report (CSMS -> CS) --------------------

const GetBaseReportFeatureName = "GetBaseReport"

// ReportBaseType indicates the type of report requested.
type ReportBaseType string

const (
	ReportBaseConfigurationInventory ReportBaseType = "ConfigurationInventory"
	ReportBaseFullInventory          ReportBaseType = "FullInventory"
	ReportBaseSummaryInventory       ReportBaseType = "SummaryInventory"
)

func isValidReportBaseType(fl validator.FieldLevel) bool {
	t := ReportBaseType(fl.Field().String())
	switch t {
	case ReportBaseConfigurationInventory, ReportBaseFullInventory, ReportBaseSummaryInventory:
		return true
	default:
		return false
	}
}

// GetBaseReportRequest is the payload for a GetBaseReport request from CSMS to CS.
type GetBaseReportRequest struct {
	RequestID  int                   `json:"requestId" validate:"gte=0"`
	ReportBase ReportBaseType        `json:"reportBase" validate:"required,reportBaseType21"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// GetBaseReportResponse is the payload for a GetBaseReport response from CS to CSMS.
type GetBaseReportResponse struct {
	Status     types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus21"`
	StatusInfo *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType          `json:"customData,omitempty" validate:"omitempty"`
}

// GetBaseReportFeature represents the GetBaseReport feature.
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

// NewGetBaseReportRequest creates a new GetBaseReportRequest.
func NewGetBaseReportRequest(requestID int, reportBase ReportBaseType) *GetBaseReportRequest {
	return &GetBaseReportRequest{RequestID: requestID, ReportBase: reportBase}
}

// NewGetBaseReportResponse creates a new GetBaseReportResponse.
func NewGetBaseReportResponse(status types.GenericDeviceModelStatus) *GetBaseReportResponse {
	return &GetBaseReportResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("reportBaseType21", isValidReportBaseType)
}
