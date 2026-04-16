package meter

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Meter Values (CS -> CSMS) --------------------

const MeterValuesFeatureName = "MeterValues"

// MeterValuesRequest is the payload for a MeterValues request from CS to CSMS.
type MeterValuesRequest struct {
	EvseID     int                   `json:"evseId" validate:"gte=0"`
	MeterValue []types.MeterValue    `json:"meterValue" validate:"required,min=1,dive"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// MeterValuesResponse is the payload for a MeterValues response from CSMS to CS.
type MeterValuesResponse struct {
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// MeterValuesFeature represents the MeterValues feature.
type MeterValuesFeature struct{}

func (f MeterValuesFeature) GetFeatureName() string {
	return MeterValuesFeatureName
}

func (f MeterValuesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(MeterValuesRequest{})
}

func (f MeterValuesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(MeterValuesResponse{})
}

func (r MeterValuesRequest) GetFeatureName() string {
	return MeterValuesFeatureName
}

func (c MeterValuesResponse) GetFeatureName() string {
	return MeterValuesFeatureName
}

// NewMeterValuesRequest creates a new MeterValuesRequest.
func NewMeterValuesRequest(evseID int, meterValue []types.MeterValue) *MeterValuesRequest {
	return &MeterValuesRequest{EvseID: evseID, MeterValue: meterValue}
}

// NewMeterValuesResponse creates a new MeterValuesResponse.
func NewMeterValuesResponse() *MeterValuesResponse {
	return &MeterValuesResponse{}
}
