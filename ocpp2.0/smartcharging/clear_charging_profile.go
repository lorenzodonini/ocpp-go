package smartcharging

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
)

// -------------------- Clear Charging Profile (CSMS -> CS) --------------------

const ClearChargingProfileFeatureName = "ClearChargingProfile"

// Status reported in ClearChargingProfileResponse.
type ClearChargingProfileStatus string

const (
	ClearChargingProfileStatusAccepted ClearChargingProfileStatus = "Accepted"
	ClearChargingProfileStatusUnknown  ClearChargingProfileStatus = "Unknown"
)

type ClearChargingProfileType struct {
	EvseID                 *int                             `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfilePurpose types.ChargingProfilePurposeType `json:"chargingProfilePurpose,omitempty" validate:"omitempty,chargingProfilePurpose"`
	StackLevel             *int                             `json:"stackLevel,omitempty" validate:"omitempty,gt=0"`
}

func isValidClearChargingProfileStatus(fl validator.FieldLevel) bool {
	status := ClearChargingProfileStatus(fl.Field().String())
	switch status {
	case ClearChargingProfileStatusAccepted, ClearChargingProfileStatusUnknown:
		return true
	default:
		return false
	}
}

// The field definition of the ClearChargingProfile request payload sent by the CSMS to the Charging Station.
type ClearChargingProfileRequest struct {
	ChargingProfileID       *int                      `json:"chargingProfileId,omitempty" validate:"omitempty"`
	ChargingProfileCriteria *ClearChargingProfileType `json:"chargingProfileCriteria,omitempty" validate:"omitempty"`
}

// This field definition of the ClearChargingProfile response payload, sent by the Charging Station to the CSMS in response to a ClearChargingProfileRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearChargingProfileResponse struct {
	Status     ClearChargingProfileStatus `json:"status" validate:"required,clearChargingProfileStatus"`
	StatusInfo *types.StatusInfo          `json:"statusInfo,omitempty" validate:"omitempty"`
}

// If the CSMS wishes to clear some or all of the charging profiles that were previously sent the Charging Station,
// it SHALL send a ClearChargingProfileRequest.
// The CSMS can use this message to clear (remove) either a specific charging profile (denoted by id) or a selection of
// charging profiles that match with the values of the optional connectorId, stackLevel and chargingProfilePurpose fields.
// The Charging Station SHALL respond with a ClearChargingProfileResponse payload specifying whether it was able to process the request.
type ClearChargingProfileFeature struct{}

func (f ClearChargingProfileFeature) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

func (f ClearChargingProfileFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearChargingProfileRequest{})
}

func (f ClearChargingProfileFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearChargingProfileResponse{})
}

func (r ClearChargingProfileRequest) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

func (c ClearChargingProfileResponse) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

// Creates a new ClearChargingProfileRequest. All fields are optional and may be set afterwards.
func NewClearChargingProfileRequest() *ClearChargingProfileRequest {
	return &ClearChargingProfileRequest{}
}

// Creates a new ClearChargingProfileResponse, containing all required fields. There are no optional fields for this message.
func NewClearChargingProfileResponse(status ClearChargingProfileStatus) *ClearChargingProfileResponse {
	return &ClearChargingProfileResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("clearChargingProfileStatus", isValidClearChargingProfileStatus)
}
