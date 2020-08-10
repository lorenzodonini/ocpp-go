package core

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// -------------------- MeterValues (CP -> CS) --------------------

const MeterValuesFeatureName = "MeterValues"

// The field definition of the MeterValues request payload sent by the Charge Point to the Central System.
type MeterValuesRequest struct {
	ConnectorId   int                `json:"connectorId" validate:"gte=0"`
	TransactionId *int               `json:"transactionId,omitempty"`
	MeterValue    []types.MeterValue `json:"meterValue" validate:"required,min=1,dive"`
}

// This field definition of the Authorize confirmation payload, sent by the Charge Point to the Central System in response to an AuthorizeRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type MeterValuesConfirmation struct {
}

// A Charge Point MAY sample the electrical meter or other sensor/transducer hardware to provide extra information about its meter values.
// It is up to the Charge Point to decide when it will send meter values.
// This can be configured using the ChangeConfiguration message to specify data acquisition intervals and specify data to be acquired & reported.
// The Charge Point SHALL send a MeterValuesRequest for offloading meter values. The request PDU SHALL contain for each sample:
// 1. The id of the Connector from which samples were taken. If the connectorId is 0, it is associated with the entire Charge Point.
// If the connectorId is 0 and the Measurand is energy related, the sample SHOULD be taken from the main energy meter.
// 2. The transactionId of the transaction to which these values are related, if applicable.
// If there is no transaction in progress or if the values are taken from the main meter, then transaction id may be omitted.
// 3. One or more meterValue elements, of type MeterValue, each representing a set of one or more data values taken at a particular point in time.
// Each MeterValue element contains a timestamp and a set of one or more individual sampledValue elements, all captured at the same point in time.
// Each sampledValue element contains a single value datum. The nature of each sampledValue is determined by the optional measurand, context, location, unit, phase, and format fields.
// The optional measurand field specifies the type of value being measured/reported. The optional context field specifies the reason/event triggering the reading.
// The optional location field specifies where the measurement is taken (e.g. Inlet, Outlet).
// The optional phase field specifies to which phase or phases of the electric installation the value applies.
// The Charging Point SHALL report all phase number dependent values from the electrical meter (or grid connection when absent) point of view.
// For individual connector phase rotation information, the Central System MAY query the ConnectorPhaseRotation configuration key on the Charging Point via GetConfiguration.
// The Charge Point SHALL report the phase rotation in respect to the grid connection. Possible values per connector are:
// NotApplicable, Unknown, RST, RTS, SRT, STR, TRS and TSR. see section Standard Configuration Key Names & Values for more information.
// The EXPERIMENTAL optional format field specifies whether the data is represented in the normal (default) form as a simple numeric value ("Raw"), or as “SignedData”, an opaque digitally signed binary data block, represented as hex data. This experimental field may be deprecated and subsequently removed in later versions, when a more mature solution alternative is provided.
// To retain backward compatibility, the default values of all of the optional fields on a sampledValue element are such that a value without any additional fields will be interpreted, as a register reading of active import energy in Wh (Watt-hour) units.
// Upon receipt of a MeterValuesRequest, the Central System SHALL respond with a MeterValuesConfirmation.
// It is likely that The Central System applies sanity checks to the data contained in a MeterValuesRequest it received. The outcome of such sanity checks SHOULD NOT ever cause the Central System to not respond with a MeterValuesConfirmation. Failing to respond with a MeterValues.conf will only cause the Charge Point to try the same message again as specified in Error responses to transaction-related messages.
type MeterValuesFeature struct{}

func (f MeterValuesFeature) GetFeatureName() string {
	return MeterValuesFeatureName
}

func (f MeterValuesFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(MeterValuesRequest{})
}

func (f MeterValuesFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(MeterValuesConfirmation{})
}

func (r MeterValuesRequest) GetFeatureName() string {
	return MeterValuesFeatureName
}

func (c MeterValuesConfirmation) GetFeatureName() string {
	return MeterValuesFeatureName
}

// Creates a new MeterValuesRequest, containing all required fields. Optional fields may be set afterwards.
func NewMeterValuesRequest(connectorId int, meterValues []types.MeterValue) *MeterValuesRequest {
	return &MeterValuesRequest{ConnectorId: connectorId, MeterValue: meterValues}
}

// Creates a new MeterValuesConfirmation, which doesn't contain any required or optional fields.
func NewMeterValuesConfirmation() *MeterValuesConfirmation {
	return &MeterValuesConfirmation{}
}
