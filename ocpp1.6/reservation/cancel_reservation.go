package reservation

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Cancel Reservation (CS -> CP) --------------------

const CancelReservationFeatureName = "CancelReservation"

// Status reported in CancelReservationConfirmation.
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

// The field definition of the CancelReservation request payload sent by the Central System to the Charge Point.
type CancelReservationRequest struct {
	ReservationId int `json:"reservationId"`
}

// This field definition of the CancelReservation confirmation payload, sent by the Charge Point to the Central System in response to a CancelReservationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type CancelReservationConfirmation struct {
	Status CancelReservationStatus `json:"status" validate:"required,cancelReservationStatus16"`
}

// To cancel a reservation the Central System SHALL send an CancelReservationRequest to the Charge Point.
// If the Charge Point has a reservation matching the reservationId in the request payload, it SHALL return status ‘Accepted’.
// Otherwise it SHALL return ‘Rejected’.
type CancelReservationFeature struct{}

func (f CancelReservationFeature) GetFeatureName() string {
	return CancelReservationFeatureName
}

func (f CancelReservationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(CancelReservationRequest{})
}

func (f CancelReservationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(CancelReservationConfirmation{})
}

func (r CancelReservationRequest) GetFeatureName() string {
	return CancelReservationFeatureName
}

func (c CancelReservationConfirmation) GetFeatureName() string {
	return CancelReservationFeatureName
}

// Creates a new CancelReservationRequest, containing all required fields. There are no optional fields for this message.
func NewCancelReservationRequest(reservationId int) *CancelReservationRequest {
	return &CancelReservationRequest{ReservationId: reservationId}
}

// Creates a new CancelReservationConfirmation, containing all required fields. There are no optional fields for this message.
func NewCancelReservationConfirmation(status CancelReservationStatus) *CancelReservationConfirmation {
	return &CancelReservationConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("cancelReservationStatus16", isValidCancelReservationStatus)
}
