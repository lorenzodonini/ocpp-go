package ocpp

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