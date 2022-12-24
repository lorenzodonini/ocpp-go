package reservation

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Reserve Now (CSMS -> CS) --------------------

const ReserveNowFeatureName = "ReserveNow"

// Status reported in ReserveNowResponse.
type ReserveNowStatus string

const (
	ReserveNowStatusAccepted    ReserveNowStatus = "Accepted"
	ReserveNowStatusFaulted     ReserveNowStatus = "Faulted"
	ReserveNowStatusOccupied    ReserveNowStatus = "Occupied"
	ReserveNowStatusRejected    ReserveNowStatus = "Rejected"
	ReserveNowStatusUnavailable ReserveNowStatus = "Unavailable"
)

func isValidReserveNowStatus(fl validator.FieldLevel) bool {
	status := ReserveNowStatus(fl.Field().String())
	switch status {
	case ReserveNowStatusAccepted, ReserveNowStatusFaulted, ReserveNowStatusOccupied, ReserveNowStatusRejected, ReserveNowStatusUnavailable:
		return true
	default:
		return false
	}
}

// Allowed ConnectorType, as supported by most charging station vendors.
// The OCPP protocol directly supports the most widely known connector types. For not mentioned types,
// refer to the Other1PhMax16A, Other1PhOver16A and Other3Ph fallbacks.
type ConnectorType string

const (
	ConnectorTypeCCS1              ConnectorType = "cCCS1"           // Combined Charging System 1 (captive cabled) a.k.a. Combo 1
	ConnectorTypeCCS2              ConnectorType = "cCCS2"           // Combined Charging System 2 (captive cabled) a.k.a. Combo 2
	ConnectorTypeG105              ConnectorType = "cG105"           // JARI G105-1993 (captive cabled) a.k.a. CHAdeMO
	ConnectorTypeTesla             ConnectorType = "cTesla"          // Tesla Connector
	ConnectorTypeCType1            ConnectorType = "cType1"          // IEC62196-2 Type 1 connector (captive cabled) a.k.a. J1772
	ConnectorTypeCType2            ConnectorType = "cType2"          // IEC62196-2 Type 2 connector (captive cabled) a.k.a. Mennekes connector
	ConnectorType3091P16A          ConnectorType = "s309-1P-16A"     // 16A 1 phase IEC60309 socket
	ConnectorType3091P32A          ConnectorType = "s309-1P-32A"     // 32A 1 phase IEC60309 socket
	ConnectorType3093P16A          ConnectorType = "s309-3P-16A"     // 16A 3 phase IEC60309 socket
	ConnectorType3093P32A          ConnectorType = "s309-3P-32A"     // 32A 3 phase IEC60309 socket
	ConnectorTypeBS1361            ConnectorType = "sBS1361"         // UK domestic socket a.k.a. 13Amp
	ConnectorTypeCEE77             ConnectorType = "sCEE-7-7"        // CEE 7/7 16A socket. May represent 7/4 & 7/5 a.k.a Schuko
	ConnectorTypeSType2            ConnectorType = "sType2"          // EC62196-2 Type 2 socket a.k.a. Mennekes connector
	ConnectorTypeSType3            ConnectorType = "sType3"          // IEC62196-2 Type 2 socket a.k.a. Scame
	ConnectorTypeOther1PhMax16A    ConnectorType = "Other1PhMax16A"  // Other single phase (domestic) sockets not mentioned above, rated at no more than 16A. CEE7/17, AS3112, NEMA 5-15, NEMA 5-20, JISC8303, TIS166, SI 32, CPCS-CCC, SEV1011, etc.
	ConnectorTypeOther1PhOver16A   ConnectorType = "Other1PhOver16A" // Other single phase sockets not mentioned above (over 16A)
	ConnectorTypeOther3Ph          ConnectorType = "Other3Ph"        // Other 3 phase sockets not mentioned above. NEMA14-30, NEMA14-50.
	ConnectorTypePan               ConnectorType = "Pan"             // Pantograph connector
	ConnectorTypeWirelessInductive ConnectorType = "wInductive"      // Wireless inductively coupled connection
	ConnectorTypeWirelessResonant  ConnectorType = "wResonant"       // Wireless resonant coupled connection
	ConnectorTypeUndetermined      ConnectorType = "Undetermined"    // Yet to be determined (e.g. before plugged in)
	ConnectorTypeUnknown           ConnectorType = "Unknown"         // Unknown; not determinable
)

func isValidConnectorType(fl validator.FieldLevel) bool {
	status := ConnectorType(fl.Field().String())
	switch status {
	case ConnectorTypeCCS1, ConnectorTypeCCS2, ConnectorTypeG105, ConnectorTypeTesla, ConnectorTypeCType1,
		ConnectorTypeCType2, ConnectorType3091P16A, ConnectorType3091P32A, ConnectorType3093P16A, ConnectorType3093P32A,
		ConnectorTypeBS1361, ConnectorTypeCEE77, ConnectorTypeSType2, ConnectorTypeSType3, ConnectorTypeOther1PhMax16A,
		ConnectorTypeOther1PhOver16A, ConnectorTypeOther3Ph, ConnectorTypePan, ConnectorTypeWirelessInductive,
		ConnectorTypeWirelessResonant, ConnectorTypeUndetermined, ConnectorTypeUnknown:
		return true
	default:
		return false
	}
}

// The field definition of the ReserveNow request payload sent by the CSMS to the Charging Station.
type ReserveNowRequest struct {
	ID             int             `json:"id" validate:"gte=0"` // ID of reservation
	ExpiryDateTime *types.DateTime `json:"expiryDateTime" validate:"required"`
	ConnectorType  ConnectorType   `json:"connectorType,omitempty" validate:"omitempty,connectorType"`
	EvseID         *int            `json:"evseId,omitempty" validate:"omitempty,gte=0"`
	IdToken        types.IdToken   `json:"idToken" validate:"required,dive"`
	GroupIdToken   *types.IdToken  `json:"groupIdToken,omitempty" validate:"omitempty,dive"`
}

// This field definition of the ReserveNow response payload, sent by the Charging Station to the CSMS in response to a ReserveNowRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ReserveNowResponse struct {
	Status     ReserveNowStatus  `json:"status" validate:"required,reserveNowStatus"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty"`
}

// To ensure an EV drive can charge their EV at a charging station, the EV driver may make a reservation until
// a certain expiry time. A user may reserve a specific EVSE.
//
// The EV driver asks the CSMS to reserve an unspecified EVSE at a charging station.
// The CSMS sends a ReserveNowRequest to a charging station.
// The charging station responds with ReserveNowResponse, with an according status.
//
// After confirming a reservation, the charging station shall asynchronously send a
// StatusNotificationRequest to the CSMS.
type ReserveNowFeature struct{}

func (f ReserveNowFeature) GetFeatureName() string {
	return ReserveNowFeatureName
}

func (f ReserveNowFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ReserveNowRequest{})
}

func (f ReserveNowFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ReserveNowResponse{})
}

func (r ReserveNowRequest) GetFeatureName() string {
	return ReserveNowFeatureName
}

func (c ReserveNowResponse) GetFeatureName() string {
	return ReserveNowFeatureName
}

// Creates a new ReserveNowRequest, containing all required fields. Optional fields may be set afterwards.
func NewReserveNowRequest(id int, expiryDateTime *types.DateTime, idToken types.IdToken) *ReserveNowRequest {
	return &ReserveNowRequest{ID: id, ExpiryDateTime: expiryDateTime, IdToken: idToken}
}

// Creates a new ReserveNowResponse, containing all required fields. Optional fields may be set afterwards.
func NewReserveNowResponse(status ReserveNowStatus) *ReserveNowResponse {
	return &ReserveNowResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("reserveNowStatus", isValidReserveNowStatus)
	_ = types.Validate.RegisterValidation("connectorType", isValidConnectorType)
}
