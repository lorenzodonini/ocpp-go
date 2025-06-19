package v2x

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- AFRRSignal (CSMS -> CS) --------------------

const AFRRSignal = "AFRRSignal"

// The field definition of the AFRRSignalRequest request payload sent by the CSMS to the Charging Station.
type AFRRSignalRequest struct {
	Timestamp *types.DateTime `json:"timestamp" validate:"required"`
	Signal    int             `json:"signal" validate:"required"`
}

// This field definition of the AFRRSignalResponse
type AFRRSignalResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty"`
}

type AFRRSignalFeature struct{}

func (f AFRRSignalFeature) GetFeatureName() string {
	return AFRRSignal
}

func (f AFRRSignalFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(AFRRSignalRequest{})
}

func (f AFRRSignalFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(AFRRSignalResponse{})
}

func (r AFRRSignalRequest) GetFeatureName() string {
	return AFRRSignal
}

func (c AFRRSignalResponse) GetFeatureName() string {
	return AFRRSignal
}

// Creates a new AFRRSignalRequest, containing all required fields. Optional fields may be set afterwards.
func NewAFRRSignalRequest(timestamp *types.DateTime, signal int) *AFRRSignalRequest {
	return &AFRRSignalRequest{
		Timestamp: timestamp,
		Signal:    signal,
	}
}

// Creates a new NewAFFRSignalResponse, containing all required fields. Optional fields may be set afterwards.
func NewAFRRSignalResponse(status types.GenericStatus) *AFRRSignalResponse {
	return &AFRRSignalResponse{Status: status}
}
