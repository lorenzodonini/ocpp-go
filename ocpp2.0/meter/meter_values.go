package meter

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Meter Values (CS -> CSMS) --------------------

const MeterValuesFeatureName = "MeterValues"

// The field definition of the MeterValues request payload sent by the Charge Point to the Central System.
type MeterValuesRequest struct {
	EvseID     int                `json:"evseId" validate:"gte=0"` // This contains a number (>0) designating an EVSE of the Charging Station. ‘0’ (zero) is used to designate the main power meter.
	MeterValue []types.MeterValue `json:"meterValue" validate:"required,min=1,dive"`
}

// This field definition of the Authorize confirmation payload, sent by the Charge Point to the Central System in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type MeterValuesResponse struct {
}


// The message is used to sample the electrical meter or other sensor/transducer hardware to provide information about the Charging Stations' Meter Values, outside of a transaction.
// The Charging Station is configured to send Meter values every XX seconds.
//
// The Charging Station samples the electrical meter or other sensor/transducer hardware to provide information about its Meter Values.
// Depending on configuration settings, the Charging Station MAY send a MeterValues request, for offloading Meter Values to the CSMS.
// Upon receipt of a MeterValuesRequest message, the CSMS responds with a MeterValuesResponse message
//
// The MeterValuesRequest and MeterValuesResponse messages are deprecated in OCPP 2.0.
// It is advised to start using Device Management Monitoring instead, see the diagnostics functional block.
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

// Creates a new MeterValuesRequest, containing all required fields. Optional fields may be set afterwards.
func NewMeterValuesRequest(evseID int, meterValues []types.MeterValue) *MeterValuesRequest {
	return &MeterValuesRequest{EvseID: evseID, MeterValue: meterValues}
}

// Creates a new MeterValuesResponse, which doesn't contain any required or optional fields.
func NewMeterValuesResponse() *MeterValuesResponse {
	return &MeterValuesResponse{}
}
