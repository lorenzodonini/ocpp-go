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
	NOT_IMPLEMENTED = "NotImplemented"
	NOT_SUPPORTED = "NotSupported"
	INTERNAL_ERROR = "InternalError"
	PROTOCOL_ERROR = "ProtocolError"
	SECURITY_ERROR = "SecurityError"
	FORMATION_VIOLATION = "FormationViolation"
	PROPERTY_CONSTRAINT_VIOLATION = "PropertyConstraintViolation"
	OCCURRENCE_CONSTRAINT_VIOLATION = "OccurrenceConstraintViolation"
	TYPE_CONSTRAINT_VIOLATION = "TypeConstraintViolation"
	GENERIC_ERROR = "GenericError"
)