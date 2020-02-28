package ocpp2

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Clear Charging Profile (CS -> CP) --------------------

// Status reported in ClearChargingProfileConfirmation.
type ClearChargingProfileStatus string

const (
	ClearChargingProfileStatusAccepted ClearChargingProfileStatus = "Accepted"
	ClearChargingProfileStatusUnknown  ClearChargingProfileStatus = "Unknown"
)

type ClearChargingProfileType struct {
	ID                     int                        `json:"id,omitempty" validate:"gte=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose,omitempty" validate:"omitempty,chargingProfilePurpose"`
	StackLevel             int                        `json:"stackLevel,omitempty" validate:"omitempty,gt=0"`
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
	EvseID          *int                      `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfile *ClearChargingProfileType `json:"chargingProfile,omitempty" validate:"omitempty,dive"`
}

// This field definition of the ClearChargingProfile confirmation payload, sent by the Charging Station to the CSMS in response to a ClearChargingProfileRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearChargingProfileConfirmation struct {
	Status ClearChargingProfileStatus `json:"status" validate:"required,clearChargingProfileStatus"`
}

// If the CSMS wishes to clear some or all of the charging profiles that were previously sent the Charging Station,
// it SHALL send a ClearChargingProfileRequest.
// The CSMS can use this message to clear (remove) either a specific charging profile (denoted by id) or a selection of
// charging profiles that match with the values of the optional connectorId, stackLevel and chargingProfilePurpose fields.
// The Charging Station SHALL respond with a ClearChargingProfileConfirmation payload specifying whether it was able to process the request.
type ClearChargingProfileFeature struct{}

func (f ClearChargingProfileFeature) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

func (f ClearChargingProfileFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearChargingProfileRequest{})
}

func (f ClearChargingProfileFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(ClearChargingProfileConfirmation{})
}

func (r ClearChargingProfileRequest) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

func (c ClearChargingProfileConfirmation) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

// Creates a new ClearChargingProfileRequest. All fields are optional and may be set afterwards.
func NewClearChargingProfileRequest() *ClearChargingProfileRequest {
	return &ClearChargingProfileRequest{}
}

// Creates a new ClearChargingProfileConfirmation, containing all required fields. There are no optional fields for this message.
func NewClearChargingProfileConfirmation(status ClearChargingProfileStatus) *ClearChargingProfileConfirmation {
	return &ClearChargingProfileConfirmation{Status: status}
}

func init() {
	_ = Validate.RegisterValidation("clearChargingProfileStatus", isValidClearChargingProfileStatus)
}
