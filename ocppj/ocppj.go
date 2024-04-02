// Contains an implementation of OCPP message dispatcher via JSON over WebSocket.
package ocppj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/lorenzodonini/ocpp-go/logging"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// The validator, used for validating incoming/outgoing OCPP messages.
var Validate = validator.New()

// The internal validation settings. Enabled by default.
var validationEnabled bool

// The internal verbose logger
var log logging.Logger

var EscapeHTML = true

func init() {
	_ = Validate.RegisterValidation("errorCode", IsErrorCodeValid)
	log = &logging.VoidLogger{}
	validationEnabled = true
}

// Sets a custom Logger implementation, allowing the ocpp-j package to log events.
// By default, a VoidLogger is used, so no logs will be sent to any output.
//
// The function panics, if a nil logger is passed.
func SetLogger(logger logging.Logger) {
	if logger == nil {
		panic("cannot set a nil logger")
	}
	log = logger
}

// Allows an instance of ocppj to configure if the message is Marshaled by escaping special caracters like "<", ">", "&" etc
// For more info https://pkg.go.dev/encoding/json#HTMLEscape
func SetHTMLEscape(flag bool) {
	EscapeHTML = flag
}

// Allows to enable/disable automatic validation for OCPP messages
// (this includes the field constraints defined for every request/response).
// The feature may be useful when working with OCPP implementations that don't fully comply to the specs.
//
// Validation is enabled by default.
//
// ⚠️ Use at your own risk! When disabled, outgoing and incoming OCPP messages will not be validated anymore,
// potentially leading to errors.
func SetMessageValidation(enabled bool) {
	validationEnabled = enabled
}

// MessageType identifies the type of message exchanged between two OCPP endpoints.
type MessageType int

const (
	CALL        MessageType = 2
	CALL_RESULT MessageType = 3
	CALL_ERROR  MessageType = 4
)

// An OCPP-J message.
type Message interface {
	// Returns the message type identifier of the message.
	GetMessageTypeId() MessageType
	// Returns the unique identifier of the message.
	GetUniqueId() string
	json.Marshaler
}

var messageIdGenerator = func() string {
	return fmt.Sprintf("%v", rand.Uint32())
}

// SetMessageIdGenerator sets a lambda function for generating unique IDs for new messages.
// The function is invoked automatically when creating a new Call.
//
// Settings this overrides the default behavior, which is:
//
//	fmt.Sprintf("%v", rand.Uint32())
func SetMessageIdGenerator(generator func() string) {
	if generator != nil {
		messageIdGenerator = generator
	}
}

// -------------------- Call --------------------

// An OCPP-J Call message, containing an OCPP Request.
type Call struct {
	Message       `validate:"-"`
	MessageTypeId MessageType  `json:"messageTypeId" validate:"required,eq=2"`
	UniqueId      string       `json:"uniqueId" validate:"required,max=36"`
	Action        string       `json:"action" validate:"required,max=36"`
	Payload       ocpp.Request `json:"payload" validate:"required"`
}

func (call *Call) GetMessageTypeId() MessageType {
	return call.MessageTypeId
}

func (call *Call) GetUniqueId() string {
	return call.UniqueId
}

func (call *Call) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 4)
	fields[0] = int(call.MessageTypeId)
	fields[1] = call.UniqueId
	fields[2] = call.Action
	fields[3] = call.Payload
	return jsonMarshal(fields)
}

// -------------------- Call Result --------------------

// An OCPP-J CallResult message, containing an OCPP Response.
type CallResult struct {
	Message
	MessageTypeId MessageType   `json:"messageTypeId" validate:"required,eq=3"`
	UniqueId      string        `json:"uniqueId" validate:"required,max=36"`
	Payload       ocpp.Response `json:"payload" validate:"required"`
}

func (callResult *CallResult) GetMessageTypeId() MessageType {
	return callResult.MessageTypeId
}

func (callResult *CallResult) GetUniqueId() string {
	return callResult.UniqueId
}

func (callResult *CallResult) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 3)
	fields[0] = int(callResult.MessageTypeId)
	fields[1] = callResult.UniqueId
	fields[2] = callResult.Payload
	return jsonMarshal(fields)
}

// -------------------- Call Error --------------------

// An OCPP-J CallError message, containing an OCPP Error.
type CallError struct {
	Message
	MessageTypeId    MessageType    `json:"messageTypeId" validate:"required,eq=4"`
	UniqueId         string         `json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ocpp.ErrorCode `json:"errorCode" validate:"errorCode"`
	ErrorDescription string         `json:"errorDescription" validate:"omitempty"`
	ErrorDetails     interface{}    `json:"errorDetails" validate:"omitempty"`
}

func (callError *CallError) GetMessageTypeId() MessageType {
	return callError.MessageTypeId
}

func (callError *CallError) GetUniqueId() string {
	return callError.UniqueId
}

func (callError *CallError) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 5)
	fields[0] = int(callError.MessageTypeId)
	fields[1] = callError.UniqueId
	fields[2] = callError.ErrorCode
	fields[3] = callError.ErrorDescription
	if callError.ErrorDetails == nil {
		fields[4] = struct{}{}
	} else {
		fields[4] = callError.ErrorDetails
	}
	return ocppMessageToJson(fields)
}

const (
	NotImplemented                ocpp.ErrorCode = "NotImplemented"                // Requested Action is not known by receiver.
	NotSupported                  ocpp.ErrorCode = "NotSupported"                  // Requested Action is recognized but not supported by the receiver.
	InternalError                 ocpp.ErrorCode = "InternalError"                 // An internal error occurred and the receiver was not able to process the requested Action successfully.
	MessageTypeNotSupported       ocpp.ErrorCode = "MessageTypeNotSupported"       // A message with an Message Type Number received that is not supported by this implementation.
	ProtocolError                 ocpp.ErrorCode = "ProtocolError"                 // Payload for Action is incomplete.
	SecurityError                 ocpp.ErrorCode = "SecurityError"                 // During the processing of Action a security issue occurred preventing receiver from completing the Action successfully.
	PropertyConstraintViolation   ocpp.ErrorCode = "PropertyConstraintViolation"   // Payload is syntactically correct but at least one field contains an invalid value.
	OccurrenceConstraintViolation ocpp.ErrorCode = "OccurrenceConstraintViolation" // Payload for Action is syntactically correct but at least one of the fields violates occurrence constraints.
	TypeConstraintViolation       ocpp.ErrorCode = "TypeConstraintViolation"       // Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12).
	GenericError                  ocpp.ErrorCode = "GenericError"                  // Any other error not covered by the previous ones.
	FormatViolationV2             ocpp.ErrorCode = "FormatViolation"               // Payload for Action is syntactically incorrect. This is only valid for OCPP 2.0.1
	FormatViolationV16            ocpp.ErrorCode = "FormationViolation"            // Payload for Action is syntactically incorrect or not conform the PDU structure for Action. This is only valid for OCPP 1.6
)

type dialector interface {
	Dialect() ocpp.Dialect
}

func FormatErrorType(d dialector) ocpp.ErrorCode {
	switch d.Dialect() {
	case ocpp.V16:
		return FormatViolationV16
	case ocpp.V2:
		return FormatViolationV2
	default:
		panic(fmt.Sprintf("invalid dialect: %v", d))
	}
}

func IsErrorCodeValid(fl validator.FieldLevel) bool {
	code := ocpp.ErrorCode(fl.Field().String())
	switch code {
	case NotImplemented, NotSupported, InternalError, MessageTypeNotSupported, ProtocolError, SecurityError, FormatViolationV16, FormatViolationV2, PropertyConstraintViolation, OccurrenceConstraintViolation, TypeConstraintViolation, GenericError:
		return true
	}
	return false
}

// -------------------- Logic --------------------

// Unmarshals an OCPP-J json object from a byte array.
// Returns the array of elements contained in the message.
func ParseRawJsonMessage(dataJson []byte) ([]interface{}, error) {
	var arr []interface{}
	err := json.Unmarshal(dataJson, &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// Unmarshals an OCPP-J json object from a JSON string.
// Returns the array of elements contained in the message.
func ParseJsonMessage(dataJson string) ([]interface{}, error) {
	rawJson := []byte(dataJson)
	return ParseRawJsonMessage(rawJson)
}

func ocppMessageToJson(message interface{}) ([]byte, error) {
	jsonData, err := jsonMarshal(message)
	if err != nil {
		return nil, err
	}
	jsonData[0] = '['
	jsonData[len(jsonData)-1] = ']'
	return jsonData, nil
}

func getValueLength(value interface{}) int {
	switch value := value.(type) {
	case int:
		return value
	case string:
		return len(value)
	default:
		return 0
	}
}

func occurenceViolation(fieldError validator.FieldError, messageId, feature string) *ocpp.Error {
	description := fmt.Sprintf("Field %s required but not found", fieldError.Namespace())
	if feature != "" {
		description = fmt.Sprintf("%s for feature %s", description, feature)
	}
	return ocpp.NewError(OccurrenceConstraintViolation, description, messageId)
}

func propertyConstraintViolation(fieldError validator.FieldError, details, messageId, feature string) *ocpp.Error {
	description := fmt.Sprintf("Field %s must be %s %s, but was %d", fieldError.Namespace(), details, fieldError.Param(), getValueLength(fieldError.Value()))
	if feature != "" {
		description = fmt.Sprintf("%s for feature %s", description, feature)
	}
	return ocpp.NewError(
		PropertyConstraintViolation,
		description,
		messageId,
	)
}

func errorFromValidation(validationErrors validator.ValidationErrors, messageId string, feature string) *ocpp.Error {
	for _, el := range validationErrors {
		switch el.ActualTag() {
		case "required":
			return occurenceViolation(el, messageId, feature)
		case "max":
			return propertyConstraintViolation(el, "maximum", messageId, feature)
		case "min":
			return propertyConstraintViolation(el, "minimum", messageId, feature)
		case "gte":
			return propertyConstraintViolation(el, ">=", messageId, feature)
		case "gt":
			return propertyConstraintViolation(el, ">", messageId, feature)
		case "lte":
			return propertyConstraintViolation(el, "<=", messageId, feature)
		case "lt":
			return propertyConstraintViolation(el, "<", messageId, feature)
		}
	}
	return ocpp.NewError(GenericError, fmt.Sprintf("%v", validationErrors.Error()), messageId)
}

// Marshals data by manipulating EscapeHTML property of encoder
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(EscapeHTML)
	err := encoder.Encode(t)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

// -------------------- Endpoint --------------------

// An OCPP-J endpoint is one of the two entities taking part in the communication.
// The endpoint keeps state for supported OCPP profiles and current pending requests.
type Endpoint struct {
	dialect  ocpp.Dialect
	Profiles []*ocpp.Profile
}

// Sets endpoint dialect.
func (endpoint *Endpoint) SetDialect(d ocpp.Dialect) {
	endpoint.dialect = d
}

// Gets endpoint dialect.
func (endpoint *Endpoint) Dialect() ocpp.Dialect {
	return endpoint.dialect
}

// Adds support for a new profile on the endpoint.
func (endpoint *Endpoint) AddProfile(profile *ocpp.Profile) {
	endpoint.Profiles = append(endpoint.Profiles, profile)
}

// Retrieves a profile for the given profile name.
// Returns a false flag in case no profile matching the specified name was found.
func (endpoint *Endpoint) GetProfile(name string) (*ocpp.Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

// Retrieves a profile for a given feature.
// Returns a false flag in case no profile supporting the specified feature was found.
func (endpoint *Endpoint) GetProfileForFeature(featureName string) (*ocpp.Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.SupportsFeature(featureName) {
			return p, true
		}
	}
	return nil, false
}

func parseRawJsonRequest(raw interface{}, requestType reflect.Type) (ocpp.Request, error) {
	if raw == nil {
		raw = &struct{}{}
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	request := reflect.New(requestType).Interface()
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		return nil, err
	}
	result := request.(ocpp.Request)
	return result, nil
}

func parseRawJsonConfirmation(raw interface{}, confirmationType reflect.Type) (ocpp.Response, error) {
	if raw == nil {
		raw = &struct{}{}
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	confirmation := reflect.New(confirmationType).Interface()
	err = json.Unmarshal(bytes, &confirmation)
	if err != nil {
		return nil, err
	}
	result := confirmation.(ocpp.Response)
	return result, nil
}

// Parses an OCPP-J message. The function expects an array of elements, as contained in the JSON message.
//
// Pending requests are automatically cleared, in case the received message is a CallResponse or CallError.
func (endpoint *Endpoint) ParseMessage(arr []interface{}, pendingRequestState ClientState) (Message, error) {
	// Checking message fields
	if len(arr) < 3 {
		return nil, ocpp.NewError(FormatErrorType(endpoint), "Invalid message. Expected array length >= 3", "")
	}
	rawTypeId, ok := arr[0].(float64)
	if !ok {
		return nil, ocpp.NewError(FormatErrorType(endpoint), fmt.Sprintf("Invalid element %v at 0, expected message type (int)", arr[0]), "")
	}
	typeId := MessageType(rawTypeId)
	uniqueId, ok := arr[1].(string)
	if !ok {
		return nil, ocpp.NewError(FormatErrorType(endpoint), fmt.Sprintf("Invalid element %v at 1, expected unique ID (string)", arr[1]), uniqueId)
	}
	if uniqueId == "" {
		return nil, ocpp.NewError(FormatErrorType(endpoint), "Invalid unique ID, cannot be empty", uniqueId)
	}
	// Parse message
	if typeId == CALL {
		if len(arr) != 4 {
			return nil, ocpp.NewError(FormatErrorType(endpoint), "Invalid Call message. Expected array length 4", uniqueId)
		}
		action, ok := arr[2].(string)
		if !ok {
			return nil, ocpp.NewError(FormatErrorType(endpoint), fmt.Sprintf("Invalid element %v at 2, expected action (string)", arr[2]), uniqueId)
		}

		profile, ok := endpoint.GetProfileForFeature(action)
		if !ok {
			return nil, ocpp.NewError(NotSupported, fmt.Sprintf("Unsupported feature %v", action), uniqueId)
		}
		request, err := profile.ParseRequest(action, arr[3], parseRawJsonRequest)
		if err != nil {
			return nil, ocpp.NewError(FormatErrorType(endpoint), err.Error(), uniqueId)
		}
		call := Call{
			MessageTypeId: CALL,
			UniqueId:      uniqueId,
			Action:        action,
			Payload:       request,
		}
		err = Validate.Struct(call)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId, action)
		}
		return &call, nil
	} else if typeId == CALL_RESULT {
		request, ok := pendingRequestState.GetPendingRequest(uniqueId)
		if !ok {
			log.Infof("No previous request %v sent. Discarding response message", uniqueId)
			return nil, nil
		}
		profile, _ := endpoint.GetProfileForFeature(request.GetFeatureName())
		confirmation, err := profile.ParseResponse(request.GetFeatureName(), arr[2], parseRawJsonConfirmation)
		if err != nil {
			return nil, ocpp.NewError(FormatErrorType(endpoint), err.Error(), uniqueId)
		}
		callResult := CallResult{
			MessageTypeId: CALL_RESULT,
			UniqueId:      uniqueId,
			Payload:       confirmation,
		}
		err = Validate.Struct(callResult)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId, request.GetFeatureName())
		}
		return &callResult, nil
	} else if typeId == CALL_ERROR {
		_, ok := pendingRequestState.GetPendingRequest(uniqueId)
		if !ok {
			log.Infof("No previous request %v sent. Discarding error message", uniqueId)
			return nil, nil
		}
		if len(arr) < 4 {
			return nil, ocpp.NewError(FormatErrorType(endpoint), "Invalid Call Error message. Expected array length >= 4", uniqueId)
		}
		var details interface{}
		if len(arr) > 4 {
			details = arr[4]
		}
		rawErrorCode, ok := arr[2].(string)
		if !ok {
			return nil, ocpp.NewError(FormatErrorType(endpoint), fmt.Sprintf("Invalid element %v at 2, expected rawErrorCode (string)", arr[2]), rawErrorCode)
		}
		errorCode := ocpp.ErrorCode(rawErrorCode)
		errorDescription := ""
		if v, ok := arr[3].(string); ok {
			errorDescription = v
		}
		callError := CallError{
			MessageTypeId:    CALL_ERROR,
			UniqueId:         uniqueId,
			ErrorCode:        errorCode,
			ErrorDescription: errorDescription,
			ErrorDetails:     details,
		}
		err := Validate.Struct(callError)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId, "")
		}
		return &callError, nil
	} else {
		return nil, ocpp.NewError(MessageTypeNotSupported, fmt.Sprintf("Invalid message type ID %v", typeId), uniqueId)
	}
}

// Creates a Call message, given an OCPP request. A unique ID for the message is automatically generated.
// Returns an error in case the request's feature is not supported on this endpoint.
//
// The created call is not automatically scheduled for transmission and is not added to the list of pending requests.
func (endpoint *Endpoint) CreateCall(request ocpp.Request) (*Call, error) {
	action := request.GetFeatureName()
	profile, _ := endpoint.GetProfileForFeature(action)
	if profile == nil {
		return nil, fmt.Errorf("Couldn't create Call for unsupported action %v", action)
	}
	// TODO: handle collisions?
	uniqueId := messageIdGenerator()
	call := Call{
		MessageTypeId: CALL,
		UniqueId:      uniqueId,
		Action:        action,
		Payload:       request,
	}
	if validationEnabled {
		err := Validate.Struct(call)
		if err != nil {
			return nil, err
		}
	}
	return &call, nil
}

// Creates a CallResult message, given an OCPP response and the message's unique ID.
//
// Returns an error in case the response's feature is not supported on this endpoint.
func (endpoint *Endpoint) CreateCallResult(confirmation ocpp.Response, uniqueId string) (*CallResult, error) {
	action := confirmation.GetFeatureName()
	profile, _ := endpoint.GetProfileForFeature(action)
	if profile == nil {
		return nil, ocpp.NewError(NotSupported, fmt.Sprintf("couldn't create Call Result for unsupported action %v", action), uniqueId)
	}
	callResult := CallResult{
		MessageTypeId: CALL_RESULT,
		UniqueId:      uniqueId,
		Payload:       confirmation,
	}
	if validationEnabled {
		err := Validate.Struct(callResult)
		if err != nil {
			return nil, err
		}
	}
	return &callResult, nil
}

// Creates a CallError message, given the message's unique ID and the error.
func (endpoint *Endpoint) CreateCallError(uniqueId string, code ocpp.ErrorCode, description string, details interface{}) (*CallError, error) {
	callError := CallError{
		MessageTypeId:    CALL_ERROR,
		UniqueId:         uniqueId,
		ErrorCode:        code,
		ErrorDescription: description,
		ErrorDetails:     details,
	}
	if validationEnabled {
		err := Validate.Struct(callError)
		if err != nil {
			return nil, err
		}
	}
	return &callError, nil
}
