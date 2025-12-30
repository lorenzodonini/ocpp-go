package reservation

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
)

// -------------------- Cancel Reservation (CSMS -> CS) --------------------

const CancelReservationFeatureName = "CancelReservation"

// Status reported in CancelReservationResponse.
type CancelReservationStatus string

const (
	CancelReservationStatusAccepted CancelReservationStatus = "Accepted"
	CancelReservationStatusRejected CancelReservationStatus = "Rejected"
)

func isValidCancelReservationStatus(fl validator.FieldLevel) bool {
	status := CancelReservationStatus(fl.Field().String())
	switch status {
	case CancelReservationStatusAccepted, CancelReservationStatusRejected:
		return true
	default:
		return false
	}
}

// The field definition of the CancelReservation request payload sent by the CSMS to the Charging Station.
type CancelReservationRequest struct {
	ReservationID int `json:"reservationId" validate:"gte=0"`
}

// This field definition of the CancelReservation response payload, sent by the Charging Station to the CSMS in response to a CancelReservationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type CancelReservationResponse struct {
	Status     CancelReservationStatus `json:"status" validate:"required,cancelReservationStatus21"`
	StatusInfo *types.StatusInfo       `json:"statusInfo,omitempty" validate:"omitempty"`
}

// To cancel a reservation the CSMS SHALL send an CancelReservationRequest to the Charging Station.
// If the Charging Station has a reservation matching the reservationId in the request payload, it SHALL return status ‘Accepted’.
// Otherwise it SHALL return ‘Rejected’.
type CancelReservationFeature struct{}

func (f CancelReservationFeature) GetFeatureName() string {
	return CancelReservationFeatureName
}

func (f CancelReservationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(CancelReservationRequest{})
}

func (f CancelReservationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(CancelReservationResponse{})
}

func (r CancelReservationRequest) GetFeatureName() string {
	return CancelReservationFeatureName
}

func (c CancelReservationResponse) GetFeatureName() string {
	return CancelReservationFeatureName
}

// Creates a new CancelReservationRequest, containing all required fields. There are no optional fields for this message.
func NewCancelReservationRequest(reservationId int) *CancelReservationRequest {
	return &CancelReservationRequest{ReservationID: reservationId}
}

// Creates a new CancelReservationResponse, containing all required fields. There are no optional fields for this message.
func NewCancelReservationResponse(status CancelReservationStatus) *CancelReservationResponse {
	return &CancelReservationResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("cancelReservationStatus21", isValidCancelReservationStatus)
}
