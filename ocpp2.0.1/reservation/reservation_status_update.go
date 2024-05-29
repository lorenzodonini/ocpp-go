package reservation

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Reservation Status Update (CS -> CSMS) --------------------

const ReservationStatusUpdateFeatureName = "ReservationStatusUpdate"

// Status reported in ReservationStatusUpdateRequest.
type ReservationUpdateStatus string

const (
	ReservationUpdateStatusExpired ReservationUpdateStatus = "Expired"
	ReservationUpdateStatusRemoved ReservationUpdateStatus = "Removed"
)

func isValidReservationUpdateStatus(fl validator.FieldLevel) bool {
	status := ReservationUpdateStatus(fl.Field().String())
	switch status {
	case ReservationUpdateStatusExpired, ReservationUpdateStatusRemoved:
		return true
	default:
		return false
	}
}

// The field definition of the ReservationStatusUpdate request payload sent by the Charging Station to the CSMS.
type ReservationStatusUpdateRequest struct {
	ReservationID int                     `json:"reservationId" validate:"gte=0"`
	Status        ReservationUpdateStatus `json:"reservationUpdateStatus" validate:"required,reservationUpdateStatus"`
}

// This field definition of the ReservationStatusUpdate response payload, sent by the CSMS to the Charging Station in response to a ReservationStatusUpdateRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ReservationStatusUpdateResponse struct {
}

// A Charging Station shall cancel an existing reservation when:
//   - the status of a targeted EVSE changes to either Faulted or Unavailable
//   - the reservation has expired, before the EV driver started using the Charging Station
//
// This message is not triggered, if a reservation is explicitly canceled by the user or the CSMS.
//
// The Charging Station sends a ReservationStatusUpdateRequest to the CSMS, with the according status set.
// The CSMS responds with a ReservationStatusUpdateResponse.
type ReservationStatusUpdateFeature struct{}

func (f ReservationStatusUpdateFeature) GetFeatureName() string {
	return ReservationStatusUpdateFeatureName
}

func (f ReservationStatusUpdateFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ReservationStatusUpdateRequest{})
}

func (f ReservationStatusUpdateFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ReservationStatusUpdateResponse{})
}

func (r ReservationStatusUpdateRequest) GetFeatureName() string {
	return ReservationStatusUpdateFeatureName
}

func (c ReservationStatusUpdateResponse) GetFeatureName() string {
	return ReservationStatusUpdateFeatureName
}

// Creates a new ReservationStatusUpdateRequest, containing all required fields. There are no optional fields for this message.
func NewReservationStatusUpdateRequest(reservationID int, status ReservationUpdateStatus) *ReservationStatusUpdateRequest {
	return &ReservationStatusUpdateRequest{ReservationID: reservationID, Status: status}
}

// Creates a new ReservationStatusUpdateResponse, which doesn't contain any required or optional fields.
func NewReservationStatusUpdateResponse() *ReservationStatusUpdateResponse {
	return &ReservationStatusUpdateResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("reservationUpdateStatus", isValidReservationUpdateStatus)
}
