package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Report (CSMS -> CS) --------------------

const GetReportFeatureName = "GetReport"

// GetReportRequest is the payload for a GetReport request from CSMS to CS.
type GetReportRequest struct {
	RequestID         int                       `json:"requestId" validate:"gte=0"`
	ComponentCriteria []ComponentCriterionType  `json:"componentCriteria,omitempty" validate:"omitempty,max=4,dive,componentCriterion21"`
	ComponentVariable []types.ComponentVariable `json:"componentVariable,omitempty" validate:"omitempty,dive"`
	CustomData        *types.CustomDataType     `json:"customData,omitempty" validate:"omitempty"`
}

// GetReportResponse is the payload for a GetReport response from CS to CSMS.
type GetReportResponse struct {
	Status     types.GenericDeviceModelStatus `json:"status" validate:"required,genericDeviceModelStatus21"`
	StatusInfo *types.StatusInfo              `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType          `json:"customData,omitempty" validate:"omitempty"`
}

// GetReportFeature represents the GetReport feature.
type GetReportFeature struct{}

func (f GetReportFeature) GetFeatureName() string {
	return GetReportFeatureName
}

func (f GetReportFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetReportRequest{})
}

func (f GetReportFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetReportResponse{})
}

func (r GetReportRequest) GetFeatureName() string {
	return GetReportFeatureName
}

func (c GetReportResponse) GetFeatureName() string {
	return GetReportFeatureName
}

// NewGetReportRequest creates a new GetReportRequest.
func NewGetReportRequest(requestID int) *GetReportRequest {
	return &GetReportRequest{RequestID: requestID}
}

// NewGetReportResponse creates a new GetReportResponse.
func NewGetReportResponse(status types.GenericDeviceModelStatus) *GetReportResponse {
	return &GetReportResponse{Status: status}
}
