// Contains an implementation of OCPP message dispatcher via JSON over WebSocket.
package ocppj

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sync/atomic"

	"gopkg.in/go-playground/validator.v9"

	"github.com/xBlaz3kx/ocpp-go/ocpp"
)

var EscapeHTML atomic.Bool

func init() {
	EscapeHTML.Store(true)
}

// Allows an instance of ocppj to configure if the message is Marshaled by escaping special caracters like "<", ">", "&" etc
// For more info https://pkg.go.dev/encoding/json#HTMLEscape
func SetHTMLEscape(flag bool) {
	EscapeHTML.Store(flag)
}

// MessageType identifies the type of message exchanged between two OCPP endpoints.
type MessageType int

const (
	CALL              MessageType = 2
	CALL_RESULT       MessageType = 3
	CALL_ERROR        MessageType = 4
	CALL_RESULT_ERROR MessageType = 5
	SEND              MessageType = 6
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

// -------------------- Call Result Error --------------------

// An OCPP-J CallResultError message, containing an OCPP Result Error.
type CallResultError struct {
	Message
	MessageTypeId    MessageType    `json:"messageTypeId" validate:"required,eq=5"`
	UniqueId         string         `json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ocpp.ErrorCode `json:"errorCode" validate:"errorCode"`
	ErrorDescription string         `json:"errorDescription" validate:"omitempty,max=255"`
	ErrorDetails     interface{}    `json:"errorDetails" validate:"omitempty"`
}

func (callError *CallResultError) GetMessageTypeId() MessageType {
	return callError.MessageTypeId
}

func (callError *CallResultError) GetUniqueId() string {
	return callError.UniqueId
}

func (callError *CallResultError) MarshalJSON() ([]byte, error) {
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

// -------------------- Send --------------------

// An OCPP-J SEND message, containing an OCPP Request.
type Send struct {
	Message       `validate:"-"`
	MessageTypeId MessageType  `json:"messageTypeId" validate:"required,eq=6"`
	UniqueId      string       `json:"uniqueId" validate:"required,max=36"`
	Action        string       `json:"action" validate:"required,max=36"`
	Payload       ocpp.Request `json:"payload" validate:"required"`
}

func (call *Send) GetMessageTypeId() MessageType {
	return call.MessageTypeId
}

func (call *Send) GetUniqueId() string {
	return call.UniqueId
}

func (call *Send) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 4)
	fields[0] = int(call.MessageTypeId)
	fields[1] = call.UniqueId
	fields[2] = call.Action
	fields[3] = call.Payload
	return jsonMarshal(fields)
}

const (
	NotImplemented                   ocpp.ErrorCode = "NotImplemented"                // Requested Action is not known by receiver.
	NotSupported                     ocpp.ErrorCode = "NotSupported"                  // Requested Action is recognized but not supported by the receiver.
	InternalError                    ocpp.ErrorCode = "InternalError"                 // An internal error occurred and the receiver was not able to process the requested Action successfully.
	MessageTypeNotSupported          ocpp.ErrorCode = "MessageTypeNotSupported"       // A message with a Message Type Number received that is not supported by this implementation.
	ProtocolError                    ocpp.ErrorCode = "ProtocolError"                 // Payload for Action is incomplete.
	RpcFrameworkError                ocpp.ErrorCode = "RpcFrameworkError"             // Content of the call is not a valid RPC Request, for example: MessageId could not be read.
	SecurityError                    ocpp.ErrorCode = "SecurityError"                 // During the processing of Action a security issue occurred preventing receiver from completing the Action successfully.
	PropertyConstraintViolation      ocpp.ErrorCode = "PropertyConstraintViolation"   // Payload is syntactically correct but at least one field contains an invalid value.
	OccurrenceConstraintViolationV2  ocpp.ErrorCode = "OccurrenceConstraintViolation" // Payload for Action is syntactically correct but at least one of the fields violates occurrence constraints.
	OccurrenceConstraintViolationV16 ocpp.ErrorCode = "OccurenceConstraintViolation"  // Payload for Action is syntactically correct but at least one of the fields violates occurrence constraints. Contains a typo in OCPP 1.6
	TypeConstraintViolation          ocpp.ErrorCode = "TypeConstraintViolation"       // Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12).
	GenericError                     ocpp.ErrorCode = "GenericError"                  // Any other error not covered by the previous ones.
	FormatViolationV2                ocpp.ErrorCode = "FormatViolation"               // Payload for Action is syntactically incorrect. This is only valid for OCPP 2.0.1
	FormatViolationV16               ocpp.ErrorCode = "FormationViolation"            // Payload for Action is syntactically incorrect or not conform the PDU structure for Action. This is only valid for OCPP 1.6
)

type dialector interface {
	Dialect() ocpp.Dialect
}

func FormatErrorType(d dialector) ocpp.ErrorCode {
	switch d.Dialect() {
	case ocpp.V16:
		return FormatViolationV16
	case ocpp.V2, ocpp.V21:
		return FormatViolationV2
	default:
		panic(fmt.Sprintf("invalid dialect: %v", d))
	}
}

func OccurrenceConstraintErrorType(d dialector) ocpp.ErrorCode {
	switch d.Dialect() {
	case ocpp.V16:
		return OccurrenceConstraintViolationV16
	case ocpp.V2, ocpp.V21:
		return OccurrenceConstraintViolationV2
	default:
		panic(fmt.Sprintf("invalid dialect: %v", d))
	}
}

func IsErrorCodeValid(fl validator.FieldLevel) bool {
	code := ocpp.ErrorCode(fl.Field().String())
	switch code {
	case NotImplemented, NotSupported, InternalError, MessageTypeNotSupported,
		ProtocolError, SecurityError, FormatViolationV16, FormatViolationV2,
		PropertyConstraintViolation, OccurrenceConstraintViolationV16, OccurrenceConstraintViolationV2,
		TypeConstraintViolation, GenericError, RpcFrameworkError:
		return true
	}
	return false
}

// -------------------- Logic --------------------

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

func occurrenceViolation(d dialector, fieldError validator.FieldError, messageId, feature string) *ocpp.Error {
	description := fmt.Sprintf("Field %s required but not found", fieldError.Namespace())
	if feature != "" {
		description = fmt.Sprintf("%s for feature %s", description, feature)
	}
	return ocpp.NewError(OccurrenceConstraintErrorType(d), description, messageId)
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

func errorFromValidation(d dialector, validationErrors validator.ValidationErrors, messageId, feature string) *ocpp.Error {
	for _, el := range validationErrors {
		switch el.ActualTag() {
		case "required":
			return occurrenceViolation(d, el, messageId, feature)
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

	// Parse message base on type
	switch typeId {
	case CALL:
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
			return nil, errorFromValidation(endpoint, err.(validator.ValidationErrors), uniqueId, action)
		}
		return &call, nil
	case CALL_RESULT:
		request, ok := pendingRequestState.GetPendingRequest(uniqueId)
		if !ok {
			// No logger available in Endpoint.ParseMessage - this is expected to be called from Server/Client which have loggers
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
			return nil, errorFromValidation(endpoint, err.(validator.ValidationErrors), uniqueId, request.GetFeatureName())
		}
		return &callResult, nil
	case CALL_ERROR:
		_, ok := pendingRequestState.GetPendingRequest(uniqueId)
		if !ok {
			// No logger available in Endpoint.ParseMessage - this is expected to be called from Server/Client which have loggers
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
			return nil, errorFromValidation(endpoint, err.(validator.ValidationErrors), uniqueId, "")
		}
		return &callError, nil
	case CALL_RESULT_ERROR:
		// Onl
		if endpoint.Dialect() != ocpp.V21 {
			return nil, ocpp.NewError(MessageTypeNotSupported, "CALL_RESULT_ERROR message is not supported in this OCPP version", uniqueId)
		}

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

		callResultError := CallResultError{
			MessageTypeId:    CALL_ERROR,
			UniqueId:         uniqueId,
			ErrorCode:        errorCode,
			ErrorDescription: errorDescription,
			ErrorDetails:     details,
		}
		err := Validate.Struct(callResultError)
		if err != nil {
			return nil, errorFromValidation(endpoint, err.(validator.ValidationErrors), uniqueId, "")
		}
		return &callResultError, nil
	case SEND:
		// SEND can be only sent in OCPP 2.1
		if endpoint.Dialect() != ocpp.V21 {
			return nil, ocpp.NewError(MessageTypeNotSupported, "SEND message is not supported in this OCPP version", uniqueId)
		}

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
			return nil, err
		}

		return &Send{
			MessageTypeId: SEND,
			UniqueId:      uniqueId,
			Action:        action,
			Payload:       request,
		}, nil
	default:
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
	if validationEnabled.Load() {
		err := Validate.Struct(call)
		if err != nil {
			return nil, err
		}
	}
	return &call, nil
}

func (endpoint *Endpoint) CreateSend(request ocpp.Request) (*Send, error) {
	action := request.GetFeatureName()
	profile, _ := endpoint.GetProfileForFeature(action)
	if profile == nil {
		return nil, fmt.Errorf("Couldn't create Send for unsupported action %v", action)
	}
	// TODO: handle collisions?
	uniqueId := messageIdGenerator()
	send := Send{
		MessageTypeId: SEND,
		UniqueId:      uniqueId,
		Action:        action,
		Payload:       request,
	}
	if validationEnabled {
		err := Validate.Struct(send)
		if err != nil {
			return nil, err
		}
	}
	return &send, nil
}

// Creates a CallResultError message, given the message's unique ID and the error.
func (endpoint *Endpoint) CreateCallResultError(uniqueId string, code ocpp.ErrorCode, description string, details interface{}) (*CallResultError, error) {
	callError := CallResultError{
		MessageTypeId:    CALL_RESULT_ERROR,
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
	if validationEnabled.Load() {
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
	if validationEnabled.Load() {
		err := Validate.Struct(callError)
		if err != nil {
			return nil, err
		}
	}
	return &callError, nil
}
