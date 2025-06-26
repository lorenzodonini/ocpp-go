package der

import (
	"reflect"
)

// -------------------- ReportDERControl (CS -> CSMS) --------------------

const ReportDERControl = "ReportDERControl"

// The field definition of the ReportDERControlRequest request payload sent by the CSMS to the Charging Station.
type ReportDERControlRequest struct {
	RequestId         int                    `json:"requestId" validate:"required"`
	Tbc               *bool                  `json:"tbc,omitempty" validate:"omitempty"`
	FixedPFAbsorb     []FixedPFGet           `json:"fixedPFAbsorb,omitempty" validate:"omitempty,max=24,dive"`
	FixedPFInject     []FixedPFGet           `json:"fixedPFInject,omitempty" validate:"omitempty,max=24,dive"`
	FixedVar          []FixedVarGet          `json:"fixedVar,omitempty" validate:"omitempty,max=24,dive"`
	LimitMaxDischarge []LimitMaxDischargeGet `json:"limitMaxDischarge,omitempty" validate:"omitempty,max=24,dive"`
	FreqDroop         []FreqDroopGet         `json:"freqDroop,omitempty" validate:"omitempty,max=24,dive"`
	EnterService      []EnterServiceGet      `json:"enterService,omitempty" validate:"omitempty,max=24,dive"`
	Gradient          []GradientGet          `json:"gradient,omitempty" validate:"omitempty,max=24,dive"`
	Curve             []DERCurveGet          `json:"curve,omitempty" validate:"omitempty,max=24,dive"`
}

// This field definition of the ReportDERControlResponse
type ReportDERControlResponse struct {
}

type ReportDERControlFeature struct{}

func (f ReportDERControlFeature) GetFeatureName() string {
	return ReportDERControl
}

func (f ReportDERControlFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ReportDERControlRequest{})
}

func (f ReportDERControlFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ReportDERControlResponse{})
}

func (r ReportDERControlRequest) GetFeatureName() string {
	return ReportDERControl
}

func (c ReportDERControlResponse) GetFeatureName() string {
	return ReportDERControl
}

// Creates a new ReportDERControlRequest, containing all required fields. Optional fields may be set afterwards.
func NewReportDERControlRequest(requestId int) *ReportDERControlRequest {
	return &ReportDERControlRequest{
		RequestId: requestId,
	}
}

// Creates a new ReportDERControlResponse, containing all required fields. Optional fields may be set afterwards.
func NewReportDERControlResponse() *ReportDERControlResponse {
	return &ReportDERControlResponse{}
}
