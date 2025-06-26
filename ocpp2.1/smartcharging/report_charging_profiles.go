package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Report Charging Profiles (CS -> CSMS) --------------------

const ReportChargingProfilesFeatureName = "ReportChargingProfiles"

// The field definition of the ReportChargingProfiles request payload sent by the Charging Station to the CSMS.
type ReportChargingProfilesRequest struct {
	RequestID           int                           `json:"requestId" validate:"gte=0"`
	ChargingLimitSource types.ChargingLimitSourceType `json:"chargingLimitSource" validate:"required,chargingLimitSource21"`
	Tbc                 bool                          `json:"tbc,omitempty" validate:"omitempty"`
	EvseID              int                           `json:"evseId" validate:"gte=0"`
	ChargingProfile     []types.ChargingProfile       `json:"chargingProfile" validate:"required,min=1,dive"`
}

// This field definition of the ReportChargingProfiles response payload, sent by the CSMS to the Charging Station in
// response to a ReportChargingProfilesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ReportChargingProfilesResponse struct {
}

// The CSMS can ask a Charging Station to report all, or a subset of all the install Charging Profiles
// from the different possible sources. This can be used for some automatic smart charging control system,
// or for debug purposes by a CSO. This is done via the GetChargingProfiles feature.
//
// A Charging Station sends a number of ReportChargingProfilesRequest messages asynchronously to the CSMS,
// after having previously received a GetChargingProfilesRequest. The CSMS acknowledges reception of the
// reports by sending a ReportChargingProfilesResponse to the Charging Station for every received request.
type ReportChargingProfilesFeature struct{}

func (f ReportChargingProfilesFeature) GetFeatureName() string {
	return ReportChargingProfilesFeatureName
}

func (f ReportChargingProfilesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ReportChargingProfilesRequest{})
}

func (f ReportChargingProfilesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ReportChargingProfilesResponse{})
}

func (r ReportChargingProfilesRequest) GetFeatureName() string {
	return ReportChargingProfilesFeatureName
}

func (c ReportChargingProfilesResponse) GetFeatureName() string {
	return ReportChargingProfilesFeatureName
}

// Creates a new ReportChargingProfilesRequest, containing all required fields. Optional fields may be set afterwards.
func NewReportChargingProfilesRequest(requestID int, chargingLimitSource types.ChargingLimitSourceType, evseID int, chargingProfile []types.ChargingProfile) *ReportChargingProfilesRequest {
	return &ReportChargingProfilesRequest{RequestID: requestID, ChargingLimitSource: chargingLimitSource, EvseID: evseID, ChargingProfile: chargingProfile}
}

// Creates a new ReportChargingProfilesResponse, which doesn't contain any required or optional fields.
func NewReportChargingProfilesResponse() *ReportChargingProfilesResponse {
	return &ReportChargingProfilesResponse{}
}
