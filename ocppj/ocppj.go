// Contains an implementation of OCPP message dispatcher via JSON over WebSocket.
package ocppj

import (
	"encoding/json"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	errors2 "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"math/rand"
	"reflect"
)

// The validator, used for validating incoming/outgoing OCPP messages.
var Validate = validator.New()

func init() {
	_ = Validate.RegisterValidation("errorCode", IsErrorCodeValid)
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
	return json.Marshal(fields)
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
	return json.Marshal(fields)
}

// -------------------- Call Error --------------------

// An OCPP-J CallError message, containing an OCPP Error.
type CallError struct {
	Message
	MessageTypeId    MessageType    `json:"messageTypeId" validate:"required,eq=4"`
	UniqueId         string         `json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ocpp.ErrorCode `json:"errorCode" validate:"errorCode"`
	ErrorDescription string         `json:"errorDescription" validate:"required"`
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
	fields[4] = callError.ErrorDetails
	return ocppMessageToJson(fields)
}

const (
	NotImplemented                ocpp.ErrorCode = "NotImplemented"                // Requested Action is not known by receiver.
	NotSupported                  ocpp.ErrorCode = "NotSupported"                  // Requested Action is recognized but not supported by the receiver.
	InternalError                 ocpp.ErrorCode = "InternalError"                 // An internal error occurred and the receiver was not able to process the requested Action successfully.
	MessageTypeNotSupported       ocpp.ErrorCode = "MessageTypeNotSupported"       // A message with an Message Type Number received that is not supported by this implementation.
	ProtocolError                 ocpp.ErrorCode = "ProtocolError"                 // Payload for Action is incomplete.
	SecurityError                 ocpp.ErrorCode = "SecurityError"                 // During the processing of Action a security issue occurred preventing receiver from completing the Action successfully.
	FormationViolation            ocpp.ErrorCode = "FormationViolation"            // Payload for Action is syntactically incorrect or not conform the PDU structure for Action.
	PropertyConstraintViolation   ocpp.ErrorCode = "PropertyConstraintViolation"   // Payload is syntactically correct but at least one field contains an invalid value.
	OccurrenceConstraintViolation ocpp.ErrorCode = "OccurrenceConstraintViolation" // Payload for Action is syntactically correct but at least one of the fields violates occurrence constraints.
	TypeConstraintViolation       ocpp.ErrorCode = "TypeConstraintViolation"       // Payload for Action is syntactically correct but at least one of the fields violates data type constraints (e.g. “somestring”: 12).
	GenericError                  ocpp.ErrorCode = "GenericError"                  // Any other error not covered by the previous ones.
)

func IsErrorCodeValid(fl validator.FieldLevel) bool {
	code := ocpp.ErrorCode(fl.Field().String())
	switch code {
	case NotImplemented, NotSupported, InternalError, MessageTypeNotSupported, ProtocolError, SecurityError, FormationViolation, PropertyConstraintViolation, OccurrenceConstraintViolation, TypeConstraintViolation, GenericError:
		return true
	}
	return false
}

// -------------------- Logic --------------------

// Unmarshals an OCPP-J json object from a byte array.
// Returns the array of elements contained in the message.
func ParseRawJsonMessage(dataJson []byte) []interface{} {
	var arr []interface{}
	err := json.Unmarshal(dataJson, &arr)
	if err != nil {
		// TODO: return error
		log.Fatal(err)
	}
	return arr
}

// Unmarshals an OCPP-J json object from a JSON string.
// Returns the array of elements contained in the message.
func ParseJsonMessage(dataJson string) []interface{} {
	rawJson := []byte(dataJson)
	return ParseRawJsonMessage(rawJson)
}

func ocppMessageToJson(message interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	jsonData[0] = '['
	jsonData[len(jsonData)-1] = ']'
	return jsonData, nil
}

func getValueLength(value interface{}) int {
	switch value.(type) {
	case int:
		return value.(int)
	case string:
		return len(value.(string))
	default:
		return 0
	}
}

func errorFromValidation(validationErrors validator.ValidationErrors, messageId string) *ocpp.Error {
	for _, el := range validationErrors {
		switch el.ActualTag() {
		case "required":
			return ocpp.NewError(OccurrenceConstraintViolation, fmt.Sprintf("Field %v required but not found", el.Namespace()), messageId)
		case "max":
			return ocpp.NewError(PropertyConstraintViolation, fmt.Sprintf("Field %v must be maximum %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value())), messageId)
		case "min":
			return ocpp.NewError(PropertyConstraintViolation, fmt.Sprintf("Field %v must be minimum %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value())), messageId)
		case "gte":
			return ocpp.NewError(PropertyConstraintViolation, fmt.Sprintf("Field %v must be >= %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value())), messageId)
		case "gt":
			return ocpp.NewError(PropertyConstraintViolation, fmt.Sprintf("Field %v must be > %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value())), messageId)
		case "lte":
			return ocpp.NewError(PropertyConstraintViolation, fmt.Sprintf("Field %v must be <= %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value())), messageId)
		case "lt":
			return ocpp.NewError(PropertyConstraintViolation, fmt.Sprintf("Field %v must be < %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value())), messageId)
		}
	}
	return ocpp.NewError(GenericError, fmt.Sprintf("%v", validationErrors.Error()), messageId)
}

// -------------------- Endpoint --------------------

// An OCPP-J endpoint is one of the two entities taking part in the communication.
// The endpoint keeps state for supported OCPP profiles and current pending requests.
type Endpoint struct {
	Profiles            []*ocpp.Profile
	PendingRequestState PendingRequestState
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
func (endpoint *Endpoint) ParseMessage(arr []interface{}) (Message, *ocpp.Error) {
	// Checking message fields
	if len(arr) < 3 {
		return nil, ocpp.NewError(FormationViolation, "Invalid message. Expected array length >= 3", "")
	}
	rawTypeId, ok := arr[0].(float64)
	if !ok {
		return nil, ocpp.NewError(FormationViolation, fmt.Sprintf("Invalid element %v at 0, expected message type (int)", arr[0]), "")
	}
	typeId := MessageType(rawTypeId)
	uniqueId, ok := arr[1].(string)
	if !ok {
		return nil, ocpp.NewError(FormationViolation, fmt.Sprintf("Invalid element %v at 1, expected unique ID (string)", arr[1]), uniqueId)
	}
	// Parse message
	if typeId == CALL {
		if len(arr) != 4 {
			return nil, ocpp.NewError(FormationViolation, "Invalid Call message. Expected array length 4", uniqueId)
		}
		action := arr[2].(string)
		profile, ok := endpoint.GetProfileForFeature(action)
		if !ok {
			return nil, ocpp.NewError(NotSupported, fmt.Sprintf("Unsupported feature %v", action), uniqueId)
		}
		request, err := profile.ParseRequest(action, arr[3], parseRawJsonRequest)
		if err != nil {
			return nil, ocpp.NewError(FormationViolation, err.Error(), uniqueId)
		}
		call := Call{
			MessageTypeId: CALL,
			UniqueId:      uniqueId,
			Action:        action,
			Payload:       request,
		}
		err = Validate.Struct(call)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId)
		}
		return &call, nil
	} else if typeId == CALL_RESULT {
		request, ok := endpoint.PendingRequestState.GetPendingRequest(uniqueId)
		if !ok {
			log.Printf("No previous request %v sent. Discarding response message", uniqueId)
			return nil, nil
		}
		profile, _ := endpoint.GetProfileForFeature(request.GetFeatureName())
		confirmation, err := profile.ParseResponse(request.GetFeatureName(), arr[2], parseRawJsonConfirmation)
		if err != nil {
			return nil, ocpp.NewError(FormationViolation, err.Error(), uniqueId)
		}
		callResult := CallResult{
			MessageTypeId: CALL_RESULT,
			UniqueId:      uniqueId,
			Payload:       confirmation,
		}
		err = Validate.Struct(callResult)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId)
		}
		return &callResult, nil
	} else if typeId == CALL_ERROR {
		_, ok := endpoint.PendingRequestState.GetPendingRequest(uniqueId)
		if !ok {
			log.Printf("No previous request %v sent. Discarding error message", uniqueId)
			return nil, nil
		}
		if len(arr) < 4 {
			return nil, ocpp.NewError(FormationViolation, "Invalid Call Error message. Expected array length >= 4", uniqueId)
		}
		var details interface{}
		if len(arr) > 4 {
			details = arr[4]
		}
		rawErrorCode := arr[2].(string)
		errorCode := ocpp.ErrorCode(rawErrorCode)
		callError := CallError{
			MessageTypeId:    CALL_ERROR,
			UniqueId:         uniqueId,
			ErrorCode:        errorCode,
			ErrorDescription: arr[3].(string),
			ErrorDetails:     details,
		}
		err := Validate.Struct(callError)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId)
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
		return nil, errors2.Errorf("Couldn't create Call for unsupported action %v", action)
	}
	// TODO: handle collisions?
	uniqueId := messageIdGenerator()
	call := Call{
		MessageTypeId: CALL,
		UniqueId:      uniqueId,
		Action:        action,
		Payload:       request,
	}
	err := Validate.Struct(call)
	if err != nil {
		return nil, err
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
		return nil, errors2.Errorf("Couldn't create Call Result for unsupported action %v", action)
	}
	callResult := CallResult{
		MessageTypeId: CALL_RESULT,
		UniqueId:      uniqueId,
		Payload:       confirmation,
	}
	err := Validate.Struct(callResult)
	if err != nil {
		return nil, err
	}
	return &callResult, nil
}

// Creates a CallError message, given the message's unique ID and the error.
func (endpoint *Endpoint) CreateCallError(uniqueId string, code ocpp.ErrorCode, description string, details interface{}) *CallError {
	callError := CallError{
		MessageTypeId:    CALL_ERROR,
		UniqueId:         uniqueId,
		ErrorCode:        code,
		ErrorDescription: description,
		ErrorDetails:     details,
	}
	return &callError
}
