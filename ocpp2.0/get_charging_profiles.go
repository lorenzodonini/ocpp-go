package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Charging Profile (CSMS -> Charging Station) --------------------

// Status reported in GetChargingProfilesConfirmation.
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
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose,omitempty" validate:"omitempty,chargingProfilePurpose"`
	StackLevel             *int                       `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
	ChargingProfileID      []int                      `json:"chargingProfileId,omitempty" validate:"omitempty,dive,gte=0"` // This field SHALL NOT contain more ids than set in ChargingProfileEntries.maxLimit
	ChargingLimitSource    []ChargingLimitSourceType  `json:"chargingLimitSource,omitempty" validate:"omitempty,max=4,dive,chargingLimitSource"`
}

// The field definition of the GetChargingProfiles request payload sent by the CSMS to the Charging Station.
type GetChargingProfilesRequest struct {
	RequestID       *int                     `json:"requestId,omitempty" validate:"omitempty,gte=0"`
	EvseID          *int                     `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfile ChargingProfileCriterion `json:"chargingProfile" validate:"required"`
}

// This field definition of the GetChargingProfiles confirmation payload, sent by the Charging Station to the CSMS in response to a GetChargingProfilesRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetChargingProfilesConfirmation struct {
	Status GetChargingProfileStatus `json:"status" validate:"required,getChargingProfileStatus"`
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

func (f GetChargingProfilesFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetChargingProfilesConfirmation{})
}

func (r GetChargingProfilesRequest) GetFeatureName() string {
	return GetChargingProfilesFeatureName
}

func (c GetChargingProfilesConfirmation) GetFeatureName() string {
	return GetChargingProfilesFeatureName
}

// Creates a new GetChargingProfilesRequest, containing all required fields. Optional fields may be set afterwards.
func NewGetChargingProfilesRequest(chargingProfile ChargingProfileCriterion) *GetChargingProfilesRequest {
	return &GetChargingProfilesRequest{ChargingProfile: chargingProfile}
}

// Creates a new GetChargingProfilesConfirmation, containing all required fields. There are no optional fields for this message.
func NewGetChargingProfilesConfirmation(status GetChargingProfileStatus) *GetChargingProfilesConfirmation {
	return &GetChargingProfilesConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("getChargingProfileStatus", isValidGetChargingProfileStatus)
}
