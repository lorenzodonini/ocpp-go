package smartcharging

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Clear Charging Profile (CS -> CP) --------------------

const ClearChargingProfileFeatureName = "ClearChargingProfile"

// Status reported in ClearChargingProfileConfirmation.
type ClearChargingProfileStatus string

const (
	ClearChargingProfileStatusAccepted ClearChargingProfileStatus = "Accepted"
	ClearChargingProfileStatusUnknown  ClearChargingProfileStatus = "Unknown"
)

func isValidClearChargingProfileStatus(fl validator.FieldLevel) bool {
	status := ClearChargingProfileStatus(fl.Field().String())
	switch status {
	case ClearChargingProfileStatusAccepted, ClearChargingProfileStatusUnknown:
		return true
	default:
		return false
	}
}

// The field definition of the ClearChargingProfile request payload sent by the Central System to the Charge Point.
type ClearChargingProfileRequest struct {
	Id                     *int                             `json:"id,omitempty" validate:"omitempty"`
	ConnectorId            *int                             `json:"connectorId,omitempty" validate:"omitempty,gte=0"`
	ChargingProfilePurpose types.ChargingProfilePurposeType `json:"chargingProfilePurpose,omitempty" validate:"omitempty,chargingProfilePurpose"`
	StackLevel             *int                             `json:"stackLevel,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the ClearChargingProfile confirmation payload, sent by the Charge Point to the Central System in response to a ClearChargingProfileRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearChargingProfileConfirmation struct {
	Status ClearChargingProfileStatus `json:"status" validate:"required,chargingProfileStatus"`
}

// If the Central System wishes to clear some or all of the charging profiles that were previously sent the Charge Point,
// it SHALL send a ClearChargingProfileRequest.
// The Central System can use this message to clear (remove) either a specific charging profile (denoted by id) or a selection of
// charging profiles that match with the values of the optional connectorId, stackLevel and chargingProfilePurpose fields.
// The Charge Point SHALL respond with a ClearChargingProfileConfirmation payload specifying whether it was able to process the request.
type ClearChargingProfileFeature struct{}

func (f ClearChargingProfileFeature) GetFeatureName() string {
	return ClearChargingProfileFeatureName
}

func (f ClearChargingProfileFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearChargingProfileRequest{})
}

func (f ClearChargingProfileFeature) GetResponseType() reflect.Type {
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
	_ = types.Validate.RegisterValidation("clearChargingProfileStatus", isValidClearChargingProfileStatus)
}
