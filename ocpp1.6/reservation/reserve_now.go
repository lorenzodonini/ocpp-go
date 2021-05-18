package reservation

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Reserve Now (CS -> CP) --------------------

const ReserveNowFeatureName = "ReserveNow"

// Status reported in ReserveNowConfirmation.
type ReservationStatus string

const (
	ReservationStatusAccepted    ReservationStatus = "Accepted"
	ReservationStatusFaulted     ReservationStatus = "Faulted"
	ReservationStatusOccupied    ReservationStatus = "Occupied"
	ReservationStatusRejected    ReservationStatus = "Rejected"
	ReservationStatusUnavailable ReservationStatus = "Unavailable"
)

func isValidReservationStatus(fl validator.FieldLevel) bool {
	status := ReservationStatus(fl.Field().String())
	switch status {
	case ReservationStatusAccepted, ReservationStatusFaulted, ReservationStatusOccupied, ReservationStatusRejected, ReservationStatusUnavailable:
		return true
	default:
		return false
	}
}

// The field definition of the ReserveNow request payload sent by the Central System to the Charge Point.
type ReserveNowRequest struct {
	ConnectorId   int             `json:"connectorId" validate:"gte=0"`
	ExpiryDate    *types.DateTime `json:"expiryDate" validate:"required"`
	IdTag         string          `json:"idTag" validate:"required,max=20"`
	ParentIdTag   string          `json:"parentIdTag,omitempty" validate:"max=20"`
	ReservationId int             `json:"reservationId"`
}

// This field definition of the ReserveNow confirmation payload, sent by the Charge Point to the Central System in response to a ReserveNowRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ReserveNowConfirmation struct {
	Status ReservationStatus `json:"status" validate:"required,reservationStatus"`
}

// A Central System can issue a ReserveNowRequest to a Charge Point to reserve a connector for use by a specific idTag.
// The Central System MAY specify a connector to be reserved.
// If the reservationId in the request matches a reservation in the Charge Point, then the Charge Point SHALL replace that reservation with the new reservation in the request.
// If the reservationId does not match any reservation in the Charge Point, then the Charge Point SHALL return the status value ‘Accepted’ if it succeeds in reserving a connector.
// The Charge Point SHALL return ‘Occupied’ if the Charge Point or the specified connector are occupied.
// The Charge Point SHALL also return ‘Occupied’ when the Charge Point or connector has been reserved for the same or another idTag.
// The Charge Point SHALL return ‘Faulted’ if the Charge Point or the connector are in the Faulted state.
// The Charge Point SHALL return ‘Unavailable’ if the Charge Point or connector are in the Unavailable state.
// The Charge Point SHALL return ‘Rejected’ if it is configured not to accept reservations.
// If the Charge Point accepts the reservation request, then it SHALL refuse charging for all incoming idTags on the reserved connector, except when the incoming idTag or the parent idTag match the idTag or parent idTag of the reservation.
// A reservation SHALL be terminated on the Charge Point when either (1) a transaction is started for the reserved idTag or parent idTag and on the reserved connector or any connector when the reserved connectorId is 0,
// or (2) when the time specified in expiryDate is reached, or (3) when the Charge Point or connector are set to Faulted or Unavailable.
// If a transaction for the reserved idTag is started, then Charge Point SHALL send the reservationId in the StartTransactionRequest payload (see Start Transaction) to notify the Central System that the reservation is terminated.
// When a reservation expires, the Charge Point SHALL terminate the reservation and make the connector available.
type ReserveNowFeature struct{}

func (f ReserveNowFeature) GetFeatureName() string {
	return ReserveNowFeatureName
}

func (f ReserveNowFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ReserveNowRequest{})
}

func (f ReserveNowFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ReserveNowConfirmation{})
}

func (r ReserveNowRequest) GetFeatureName() string {
	return ReserveNowFeatureName
}

func (c ReserveNowConfirmation) GetFeatureName() string {
	return ReserveNowFeatureName
}

// Creates a new ReserveNowRequest, containing all required fields. Optional fields may be set afterwards.
func NewReserveNowRequest(connectorId int, expiryDate *types.DateTime, idTag string, reservationId int) *ReserveNowRequest {
	return &ReserveNowRequest{ConnectorId: connectorId, ExpiryDate: expiryDate, IdTag: idTag, ReservationId: reservationId}
}

// Creates a new ReserveNowConfirmation, containing all required fields. There are no optional fields for this message.
func NewReserveNowConfirmation(status ReservationStatus) *ReserveNowConfirmation {
	return &ReserveNowConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("reservationStatus", isValidReservationStatus)
}
