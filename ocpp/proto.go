package ocpp

import (
	"encoding/json"
	"fmt"
	errors2 "github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"math/rand"
	"reflect"
)

type Feature interface {
	GetFeatureName() string
	GetRequestType() reflect.Type
	GetConfirmationType() reflect.Type
}

type Request interface {
	GetFeatureName() string
}

type Confirmation interface {
	GetFeatureName() string
}

type ProtoError struct {
	Error     error
	ErrorCode ErrorCode
	MessageId string
}

var validate = validator.New()

// -------------------- Profile --------------------
type Profile struct {
	Name     string
	Features map[string]Feature
}

func NewProfile(name string, features ...Feature) *Profile {
	profile := Profile{Name: name, Features: make(map[string]Feature)}
	for _, feature := range features {
		profile.AddFeature(feature)
	}
	return &profile
}

func (p *Profile) AddFeature(feature Feature) {
	p.Features[feature.GetFeatureName()] = feature
}

func (p *Profile) SupportsFeature(name string) bool {
	_, ok := p.Features[name]
	return ok
}

func (p *Profile) GetFeature(name string) Feature {
	return p.Features[name]
}

func (p *Profile) ParseRequest(featureName string, rawRequest interface{}) Request {
	feature, ok := p.Features[featureName]
	if !ok {
		log.Printf("Feature %s not found", featureName)
		return nil
	}
	requestType := feature.GetRequestType()
	bytes, _ := json.Marshal(rawRequest)
	//bytes := []byte(rawRequest)
	if !ok {
		log.Printf("Couldn't cast raw request to bytes")
		return nil
	}
	request := reflect.New(requestType).Interface()
	err := json.Unmarshal(bytes, &request)
	if err != nil {
		log.Printf("Error while parsing json %v", err)
	}
	log.Print(reflect.TypeOf(request))
	result := request.(Request)
	log.Print(reflect.TypeOf(result))
	return result
}

func (p *Profile) ParseConfirmation(featureName string, rawConfirmation interface{}) Confirmation {
	feature, ok := p.Features[featureName]
	if !ok {
		log.Printf("Feature %s not found", featureName)
		return nil
	}
	requestType := feature.GetConfirmationType()
	bytes, _ := json.Marshal(rawConfirmation)
	if !ok {
		log.Printf("Couldn't cast raw confirmation to bytes")
		return nil
	}
	confirmation := reflect.New(requestType).Interface()
	err := json.Unmarshal(bytes, &confirmation)
	if err != nil {
		log.Printf("Error while parsing json %v", err)
	}
	result := confirmation.(Confirmation)
	return result
}

// -------------------- Message --------------------
type MessageType int

const (
	CALL        MessageType = 2
	CALL_RESULT MessageType = 3
	CALL_ERROR  MessageType = 4
)

type Message interface {
	GetMessageTypeId() MessageType
	GetUniqueId() string
	json.Marshaler
}

// -------------------- Call --------------------
type Call struct {
	Message       `validate:"-"`
	MessageTypeId MessageType `json:"messageTypeId" validate:"required,eq=2"`
	UniqueId      string      `json:"uniqueId" validate:"required,max=36"`
	Action        string      `json:"action" validate:"required,max=36"`
	Payload       Request     `json:"payload" validate:"required"`
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
type CallResult struct {
	Message
	MessageTypeId MessageType  `json:"messageTypeId" validate:"required,eq=3"`
	UniqueId      string       `json:"uniqueId" validate:"required,max=36"`
	Payload       Confirmation `json:"payload" validate:"required"`
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
type ErrorCode string

type CallError struct {
	Message
	MessageTypeId    MessageType `json:"messageTypeId" validate:"required,eq=4"`
	UniqueId         string      `json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ErrorCode   `json:"errorCode" validate:"-"` //TODO: check if error is supported
	ErrorDescription string      `json:"errorDescription" validate:"required"`
	ErrorDetails     interface{} `json:"errorDetails" validate:"omitempty"`
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
	NotImplemented                ErrorCode = "NotImplemented"
	NotSupported                  ErrorCode = "NotSupported"
	InternalError                 ErrorCode = "InternalError"
	ProtocolError                 ErrorCode = "ProtocolError"
	SecurityError                 ErrorCode = "SecurityError"
	FormationViolation            ErrorCode = "FormationViolation"
	PropertyConstraintViolation   ErrorCode = "PropertyConstraintViolation"
	OccurrenceConstraintViolation ErrorCode = "OccurrenceConstraintViolation"
	TypeConstraintViolation       ErrorCode = "TypeConstraintViolation"
	GenericError                  ErrorCode = "GenericError"
)

// -------------------- Logic --------------------
func ParseRawJsonMessage(dataJson []byte) []interface{} {
	var arr []interface{}
	err := json.Unmarshal(dataJson, &arr)
	if err != nil {
		// TODO: return error
		log.Fatal(err)
	}
	return arr
}

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

func newProtoError(validationErrors validator.ValidationErrors, messageId string) *ProtoError {
	for _, el := range validationErrors {
		switch el.ActualTag() {
		case "required":
			return &ProtoError{MessageId: messageId, ErrorCode: OccurrenceConstraintViolation, Error: errors2.Errorf("Field %v required but not found", el.Namespace())}
		case "max":
			return &ProtoError{MessageId: messageId, ErrorCode: PropertyConstraintViolation, Error: errors2.Errorf("Field %v must be maximum %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value()))}
		case "min":
			return &ProtoError{MessageId: messageId, ErrorCode: PropertyConstraintViolation, Error: errors2.Errorf("Field %v must be minimum %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value()))}
		case "gte":
			return &ProtoError{MessageId: messageId, ErrorCode: PropertyConstraintViolation, Error: errors2.Errorf("Field %v must be >= %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value()))}
		case "gt":
			return &ProtoError{MessageId: messageId, ErrorCode: PropertyConstraintViolation, Error: errors2.Errorf("Field %v must be > %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value()))}
		case "lte":
			return &ProtoError{MessageId: messageId, ErrorCode: PropertyConstraintViolation, Error: errors2.Errorf("Field %v must be <= %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value()))}
		case "lt":
			return &ProtoError{MessageId: messageId, ErrorCode: PropertyConstraintViolation, Error: errors2.Errorf("Field %v must be < %v, but was %v", el.Namespace(), el.Param(), getValueLength(el.Value()))}
		}
	}
	return &ProtoError{MessageId: messageId, ErrorCode: GenericError, Error: errors2.Errorf("%v", validationErrors.Error())}
}

// -------------------- Endpoint --------------------
type Endpoint struct {
	Profiles        []*Profile
	pendingRequests map[string]Request
}

func (endpoint *Endpoint) AddProfile(profile *Profile) {
	endpoint.Profiles = append(endpoint.Profiles, profile)
}

func (endpoint *Endpoint) GetProfile(name string) (*Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

func (endpoint *Endpoint) GetProfileForFeature(featureName string) (*Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.SupportsFeature(featureName) {
			return p, true
		}
	}
	return nil, false
}

func (endpoint *Endpoint) AddPendingRequest(id string, request Request) {
	endpoint.pendingRequests[id] = request
}

func (endpoint *Endpoint) GetPendingRequest(id string) (Request, bool) {
	request, ok := endpoint.pendingRequests[id]
	return request, ok
}

func (endpoint *Endpoint) DeletePendingRequest(id string) {
	delete(endpoint.pendingRequests, id)
}

func (endpoint *Endpoint) clearPendingRequests() {
	endpoint.pendingRequests = map[string]Request{}
}

func (endpoint *Endpoint) ParseMessage(arr []interface{}) (Message, *ProtoError) {
	// Checking message fields
	if len(arr) < 3 {
		return nil, &ProtoError{ErrorCode: FormationViolation, Error: errors2.Errorf("Invalid message. Expected array length >= 3")}
	}
	rawTypeId, ok := arr[0].(float64)
	if !ok {
		return nil, &ProtoError{ErrorCode: FormationViolation, Error: errors2.Errorf("Invalid element %v at 0, expected message type (int)", arr[0])}
	}
	typeId := MessageType(rawTypeId)
	uniqueId, ok := arr[1].(string)
	if !ok {
		return nil, &ProtoError{ErrorCode: FormationViolation, Error: errors2.Errorf("Invalid element %v at 1, expected unique ID (string)", arr[1])}
	}
	// Parse message
	if typeId == CALL {
		if len(arr) != 4 {
			return nil, &ProtoError{MessageId: uniqueId, ErrorCode: FormationViolation, Error: errors2.Errorf("Invalid Call message. Expected array length 4")}
		}
		action := arr[2].(string)
		profile, ok := endpoint.GetProfileForFeature(action)
		if !ok {
			return nil, &ProtoError{MessageId: uniqueId, ErrorCode: NotSupported, Error: errors2.Errorf("Unsupported feature %v", action)}
		}
		request := profile.ParseRequest(action, arr[3])
		call := Call{
			MessageTypeId: CALL,
			UniqueId:      uniqueId,
			Action:        action,
			Payload:       request,
		}
		err := validate.Struct(call)
		if err != nil {
			protoError := newProtoError(err.(validator.ValidationErrors), uniqueId)
			return nil, protoError
		}
		return &call, nil
	} else if typeId == CALL_RESULT {
		request, ok := endpoint.pendingRequests[uniqueId]
		if !ok {
			log.Printf("No previous request %v sent. Discarding response message", uniqueId)
			return nil, nil
		}
		profile, _ := endpoint.GetProfileForFeature(request.GetFeatureName())
		confirmation := profile.ParseConfirmation(request.GetFeatureName(), arr[2])
		callResult := CallResult{
			MessageTypeId: CALL_RESULT,
			UniqueId:      uniqueId,
			Payload:       confirmation,
		}
		endpoint.DeletePendingRequest(callResult.GetUniqueId())
		err := validate.Struct(callResult)
		if err != nil {
			protoError := newProtoError(err.(validator.ValidationErrors), uniqueId)
			return nil, protoError
		}
		return &callResult, nil
	} else if typeId == CALL_ERROR {
		if len(arr) < 4 {
			return nil, &ProtoError{MessageId: uniqueId, ErrorCode: FormationViolation, Error: errors2.Errorf("Invalid Call Error message. Expected array length >= 4")}
		}
		var details interface{}
		if len(arr) > 4 {
			details = arr[4]
		}
		rawErrorCode := arr[2].(string)
		errorCode := ErrorCode(rawErrorCode)
		callError := CallError{
			MessageTypeId:    CALL_ERROR,
			UniqueId:         uniqueId,
			ErrorCode:        errorCode,
			ErrorDescription: arr[3].(string),
			ErrorDetails:     details,
		}
		endpoint.DeletePendingRequest(callError.GetUniqueId())
		err := validate.Struct(callError)
		if err != nil {
			protoError := newProtoError(err.(validator.ValidationErrors), uniqueId)
			return nil, protoError
		}
		return &callError, nil
	} else {
		return nil, &ProtoError{MessageId: uniqueId, ErrorCode: FormationViolation, Error: errors2.Errorf("Invalid message type ID %v", typeId)}
	}
}

func (endpoint *Endpoint) CreateCall(request Request) (*Call, error) {
	action := request.GetFeatureName()
	profile, _ := endpoint.GetProfileForFeature(action)
	if profile == nil {
		return nil, errors2.Errorf("Couldn't create Call for unsupported action %v", action)
	}
	uniqueId := fmt.Sprintf("%v", rand.Uint32())
	call := Call{
		MessageTypeId: CALL,
		UniqueId:      uniqueId,
		Action:        action,
		Payload:       request,
	}
	err := validate.Struct(call)
	if err != nil {
		return nil, err
	}
	endpoint.AddPendingRequest(uniqueId, request)
	return &call, nil
}

func (endpoint *Endpoint) CreateCallResult(confirmation Confirmation, uniqueId string) (*CallResult, error) {
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
	err := validate.Struct(callResult)
	if err != nil {
		return nil, err
	}
	return &callResult, nil
}

func (endpoint *Endpoint) CreateCallError(uniqueId string, code ErrorCode, description string, details interface{}) *CallError {
	callError := CallError{
		MessageTypeId:    CALL_ERROR,
		UniqueId:         uniqueId,
		ErrorCode:        code,
		ErrorDescription: description,
		ErrorDetails:     details,
	}
	return &callError
}
