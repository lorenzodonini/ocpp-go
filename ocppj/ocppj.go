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

var Validate = validator.New()

func init() {
	_ = Validate.RegisterValidation("errorCode", IsErrorCodeValid)
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

var messageIdGenerator = func() string {
	return fmt.Sprintf("%v", rand.Uint32())
}

func SetMessageIdGenerator(generator func() string) {
	if generator != nil {
		messageIdGenerator = generator
	}
}

// -------------------- Call --------------------
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
type CallResult struct {
	Message
	MessageTypeId MessageType       `json:"messageTypeId" validate:"required,eq=3"`
	UniqueId      string            `json:"uniqueId" validate:"required,max=36"`
	Payload       ocpp.Confirmation `json:"payload" validate:"required"`
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
	NotImplemented                ocpp.ErrorCode = "NotImplemented"
	NotSupported                  ocpp.ErrorCode = "NotSupported"
	InternalError                 ocpp.ErrorCode = "InternalError"
	ProtocolError                 ocpp.ErrorCode = "ProtocolError"
	SecurityError                 ocpp.ErrorCode = "SecurityError"
	FormationViolation            ocpp.ErrorCode = "FormationViolation"
	PropertyConstraintViolation   ocpp.ErrorCode = "PropertyConstraintViolation"
	OccurrenceConstraintViolation ocpp.ErrorCode = "OccurrenceConstraintViolation"
	TypeConstraintViolation       ocpp.ErrorCode = "TypeConstraintViolation"
	GenericError                  ocpp.ErrorCode = "GenericError"
)

func IsErrorCodeValid(fl validator.FieldLevel) bool {
	code := ocpp.ErrorCode(fl.Field().String())
	switch code {
	case NotImplemented, NotSupported, InternalError, ProtocolError, SecurityError, FormationViolation, PropertyConstraintViolation, OccurrenceConstraintViolation, TypeConstraintViolation, GenericError:
		return true
	}
	return false
}

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
type Endpoint struct {
	Profiles        []*ocpp.Profile
	pendingRequests map[string]ocpp.Request
}

func (endpoint *Endpoint) AddProfile(profile *ocpp.Profile) {
	endpoint.Profiles = append(endpoint.Profiles, profile)
}

func (endpoint *Endpoint) GetProfile(name string) (*ocpp.Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

func (endpoint *Endpoint) GetProfileForFeature(featureName string) (*ocpp.Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.SupportsFeature(featureName) {
			return p, true
		}
	}
	return nil, false
}

func (endpoint *Endpoint) AddPendingRequest(id string, request ocpp.Request) {
	endpoint.pendingRequests[id] = request
}

func (endpoint *Endpoint) GetPendingRequest(id string) (ocpp.Request, bool) {
	request, ok := endpoint.pendingRequests[id]
	return request, ok
}

func (endpoint *Endpoint) DeletePendingRequest(id string) {
	delete(endpoint.pendingRequests, id)
}

func (endpoint *Endpoint) clearPendingRequests() {
	endpoint.pendingRequests = map[string]ocpp.Request{}
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

func parseRawJsonConfirmation(raw interface{}, confirmationType reflect.Type) (ocpp.Confirmation, error) {
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	confirmation := reflect.New(confirmationType).Interface()
	err = json.Unmarshal(bytes, &confirmation)
	if err != nil {
		return nil, err
	}
	result := confirmation.(ocpp.Confirmation)
	return result, nil
}

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
		request, ok := endpoint.pendingRequests[uniqueId]
		if !ok {
			log.Printf("No previous request %v sent. Discarding response message", uniqueId)
			return nil, nil
		}
		profile, _ := endpoint.GetProfileForFeature(request.GetFeatureName())
		confirmation, err := profile.ParseConfirmation(request.GetFeatureName(), arr[2], parseRawJsonConfirmation)
		if err != nil {
			return nil, ocpp.NewError(FormationViolation, err.Error(), uniqueId)
		}
		callResult := CallResult{
			MessageTypeId: CALL_RESULT,
			UniqueId:      uniqueId,
			Payload:       confirmation,
		}
		endpoint.DeletePendingRequest(callResult.GetUniqueId())
		err = Validate.Struct(callResult)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId)
		}
		return &callResult, nil
	} else if typeId == CALL_ERROR {
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
		endpoint.DeletePendingRequest(callError.GetUniqueId())
		err := Validate.Struct(callError)
		if err != nil {
			return nil, errorFromValidation(err.(validator.ValidationErrors), uniqueId)
		}
		return &callError, nil
	} else {
		return nil, ocpp.NewError(FormationViolation, fmt.Sprintf("Invalid message type ID %v", typeId), uniqueId)
	}
}

func (endpoint *Endpoint) CreateCall(request ocpp.Request) (*Call, error) {
	action := request.GetFeatureName()
	profile, _ := endpoint.GetProfileForFeature(action)
	if profile == nil {
		return nil, errors2.Errorf("Couldn't create Call for unsupported action %v", action)
	}
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
	endpoint.AddPendingRequest(uniqueId, request)
	return &call, nil
}

func (endpoint *Endpoint) CreateCallResult(confirmation ocpp.Confirmation, uniqueId string) (*CallResult, error) {
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
