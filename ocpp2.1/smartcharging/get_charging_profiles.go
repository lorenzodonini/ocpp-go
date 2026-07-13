package smartcharging

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Charging Profiles (CSMS -> Charging Station) --------------------

const GetChargingProfilesFeatureName = "GetChargingProfiles"

// Status reported in GetChargingProfilesResponse.
type GetChargingProfileStatus string

const (
	GetChargingProfileStatusAccepted   GetChargingProfileStatus = "Accepted"
	GetChargingProfileStatusNoProfiles GetChargingProfileStatus = "NoProfiles"
)

func isValidGetChargingProfileStatus(fl validator.FieldLevel) bool {
	status := GetChargingProfileStatus(fl.Field().String())
	switch status {
	case GetChargingProfileStatusAccepted, GetChargingProfileStatusNoProfiles:
		return true
	default:
		return false
	}
}

// ChargingProfileCriterion specifies the charging profile within a GetChargingProfilesRequest.
// A ChargingProfile consists of ChargingSchedule, describing the amount of power or current that can be delivered per time interval.
type ChargingProfileCriterion struct {
	ChargingProfilePurpose types.ChargingProfilePurposeType `json:"chargingProfilePurpose,omitempty" validate:"omitempty,chargingProfilePurpose21"`
	StackLevel             *int                             `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
	ChargingProfileID      []int                            `json:"chargingProfileId,omitempty" validate:"omitempty"` // This field SHALL NOT contain more ids than set in ChargingProfileEntries.maxLimit
	ChargingLimitSource    []types.ChargingLimitSourceType  `json:"chargingLimitSource,omitempty" validate:"omitempty,max=4,dive,chargingLimitSource21"`
}

// The field definition of the GetChargingProfiles request payload sent by the CSMS to the Charging Station.
type GetChargingProfilesRequest struct {
	RequestID       int                      `json:"requestId"`
	EvseID          *int                     `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfile ChargingProfileCriterion `json:"chargingProfile" validate:"required"`
}

// This field definition of the GetChargingProfiles response payload, sent by the Charging Station to the CSMS in response to a GetChargingProfilesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetChargingProfilesResponse struct {
	Status     GetChargingProfileStatus `json:"status" validate:"required,getChargingProfileStatus21"`
	StatusInfo *types.StatusInfo        `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSMS MAY ask a Charging Station to report all, or a subset of all the install Charging Profiles from the different possible sources, by sending a GetChargingProfilesRequest.
// This can be used for some automatic smart charging control system, or for debug purposes by a CSO.
// The Charging Station SHALL respond, indicating if it can report Charging Schedules by sending a GetChargingProfilesResponse message.
type GetChargingProfilesFeature struct{}

func (f GetChargingProfilesFeature) GetFeatureName() string {
	return GetChargingProfilesFeatureName
}

func (f GetChargingProfilesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetChargingProfilesRequest{})
}

func (f GetChargingProfilesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetChargingProfilesResponse{})
}

func (r GetChargingProfilesRequest) GetFeatureName() string {
	return GetChargingProfilesFeatureName
}

func (c GetChargingProfilesResponse) GetFeatureName() string {
	return GetChargingProfilesFeatureName
}

// Creates a new GetChargingProfilesRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetChargingProfilesRequest(chargingProfile ChargingProfileCriterion) *GetChargingProfilesRequest {
	return &GetChargingProfilesRequest{ChargingProfile: chargingProfile}
}

// Creates a new GetChargingProfilesResponse, containing all required fields. Optional fields may be set afterwards.
func NewGetChargingProfilesResponse(status GetChargingProfileStatus) *GetChargingProfilesResponse {
	return &GetChargingProfilesResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("getChargingProfileStatus21", isValidGetChargingProfileStatus)
}
