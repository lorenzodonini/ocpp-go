package ocpp

import "time"

const (
	ISO8601 = "2006-01-02T15:04:05Z"
)

type PropertyViolation struct {
	error
	Property string

}

func (e* PropertyViolation) Error() string {
	return ""
}

type ErrorCode string

const (
	NotImplemented = "NotImplemented"
	NotSupported = "NotSupported"
	InternalError = "InternalError"
	ProtocolError = "ProtocolError"
	SecurityError = "SecurityError"
	FormationViolation = "FormationViolation"
	PropertyConstraintViolation = "PropertyConstraintViolation"
	OccurrenceConstraintViolation = "OccurrenceConstraintViolation"
	TypeConstraintViolation = "TypeConstraintViolation"
	GenericError = "GenericError"
)

type RegistrationStatus string

const (
	RegistrationStatusAccepted = "Accepted"
	RegistrationStatusPending = "Pending"
	RegistrationStatusRejected = "Rejected"
)

type AuthorizationStatus string

const (
	AuthorizationStatusAccepted = "Accepted"
	AuthorizationStatusBlocked = "Blocked"
	AuthorizationStatusExpired = "Expired"
	AuthorizationStatusInvalid = "Invalid"
	AuthorizationStatusConcurrentTx = "ConcurrentTx"
)

type IdTagInfo struct {
	ExpiryDate time.Time		`json:"expiryDate" validate:"omitempty,gt"`
	ParentIdTag string			`json:"parentIdTag" validate:"omitempty,max=20"`
	Status AuthorizationStatus	`json:"status" validate:"required"`
}