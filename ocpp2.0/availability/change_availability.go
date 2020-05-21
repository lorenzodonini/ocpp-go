package availability

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Change Availability (CSMS -> CS) --------------------

const ChangeAvailabilityFeatureName = "ChangeAvailability"

// Requested availability change in ChangeAvailabilityRequest.
type OperationalStatus string

const (
	OperationalStatusInoperative OperationalStatus = "Inoperative"
	OperationalStatusOperative   OperationalStatus = "Operative"
)

func isValidOperationalStatus(fl validator.FieldLevel) bool {
	status := OperationalStatus(fl.Field().String())
	switch status {
	case OperationalStatusInoperative, OperationalStatusOperative:
		return true
	default:
		return false
	}
}

// Status returned in response to ChangeAvailabilityRequest
type ChangeAvailabilityStatus string

const (
	ChangeAvailabilityStatusAccepted  ChangeAvailabilityStatus = "Accepted"
	ChangeAvailabilityStatusRejected  ChangeAvailabilityStatus = "Rejected"
	ChangeAvailabilityStatusScheduled ChangeAvailabilityStatus = "Scheduled"
)

func isValidChangeAvailabilityStatus(fl validator.FieldLevel) bool {
	status := ChangeAvailabilityStatus(fl.Field().String())
	switch status {
	case ChangeAvailabilityStatusAccepted, ChangeAvailabilityStatusRejected, ChangeAvailabilityStatusScheduled:
		return true
	default:
		return false
	}
}

// The field definition of the ChangeAvailability request payload sent by the CSMS to the Charging Station.
type ChangeAvailabilityRequest struct {
	EvseID            int               `json:"evseId" validate:"gte=0"`
	OperationalStatus OperationalStatus `json:"operationalStatus" validate:"required,operationalStatus"`
}

// This field definition of the ChangeAvailability response payload, sent by the Charging Station to the CSMS in response to a ChangeAvailabilityRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ChangeAvailabilityResponse struct {
	Status ChangeAvailabilityStatus `json:"status" validate:"required,changeAvailabilityStatus"`
}

// CSMS can request a Charging Station to change its availability.
// A Charging Station is considered available (“operative”) when it is charging or ready for charging.
// A Charging Station is considered unavailable when it does not allow any charging.
// The CSMS SHALL send a ChangeAvailabilityRequest for requesting a Charging Station to change its availability.
// The CSMS can change the availability to available or unavailable.
type ChangeAvailabilityFeature struct{}

func (f ChangeAvailabilityFeature) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (f ChangeAvailabilityFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityRequest{})
}

func (f ChangeAvailabilityFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ChangeAvailabilityResponse{})
}

func (r ChangeAvailabilityRequest) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

func (c ChangeAvailabilityResponse) GetFeatureName() string {
	return ChangeAvailabilityFeatureName
}

// Creates a new ChangeAvailabilityRequest, containing all required fields. There are no optional fields for this message.
func NewChangeAvailabilityRequest(evseID int, operationalStatus OperationalStatus) *ChangeAvailabilityRequest {
	return &ChangeAvailabilityRequest{EvseID: evseID, OperationalStatus: operationalStatus}
}

// Creates a new ChangeAvailabilityResponse, containing all required fields. There are no optional fields for this message.
func NewChangeAvailabilityResponse(status ChangeAvailabilityStatus) *ChangeAvailabilityResponse {
	return &ChangeAvailabilityResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("operationalStatus", isValidOperationalStatus)
	_ = types.Validate.RegisterValidation("changeAvailabilityStatus", isValidChangeAvailabilityStatus)
}
