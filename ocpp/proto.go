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

var validate = validator.New()

// -------------------- Profile --------------------
type Profile struct {
	Name string
	Features map[string]Feature
}

func NewProfile(name string, features ...Feature) *Profile {
	profile := Profile{Name: name, Features: make(map[string]Feature)}
	for _, feature := range features {
		profile.AddFeature(feature)
	}
	return &profile
}

func (p* Profile) AddFeature(feature Feature) {
	p.Features[feature.GetFeatureName()] = feature
}

func (p* Profile) SupportsFeature(name string) bool {
	_, ok := p.Features[name]
	return ok
}

func (p* Profile) ParseRequest(featureName string, rawRequest interface{}) Request {
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

func (p* Profile) ParseConfirmation(featureName string, rawConfirmation interface{}) Confirmation {
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
	CALL 		= 2
	CALL_RESULT = 3
	CALL_ERROR 	= 4
)

type Message interface {
	GetMessageTypeId() MessageType
	GetUniqueId() string
	json.Marshaler
}

// -------------------- Call --------------------
type Call struct {
	Message					  `validate:"-"`
	MessageTypeId MessageType `json:"messageTypeId" validate:"required,eq=2"`
	UniqueId      string      `json:"uniqueId" validate:"required,max=36"`
	Action        string      `json:"action" validate:"required,max=36"`
	Payload       Request     `json:"payload" validate:"required"`
}

func (call* Call)GetMessageTypeId() MessageType {
	return call.MessageTypeId
}

func (call* Call)GetUniqueId() string {
	return call.UniqueId
}

func (call* Call) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 4)
	fields[0] = call.MessageTypeId
	fields[1] = call.UniqueId
	fields[2] = call.Action
	fields[3] = call.Payload
	return json.Marshal(fields)
}

// -------------------- Call Result --------------------
type CallResult struct {
	Message
	MessageTypeId MessageType 	`json:"messageTypeId" validate:"required,eq=3"`
	UniqueId      string      	`json:"uniqueId" validate:"required,max=36"`
	Payload       Confirmation 	`json:"payload" validate:"required"`
}

func (callResult* CallResult)GetMessageTypeId() MessageType {
	return callResult.MessageTypeId
}

func (callResult* CallResult)GetUniqueId() string {
	return callResult.UniqueId
}

func (callResult *CallResult) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 3)
	fields[0] = callResult.MessageTypeId
	fields[1] = callResult.UniqueId
	fields[2] = callResult.Payload
	return json.Marshal(fields)
}

// -------------------- Call Error --------------------
type CallError struct {
	Message
	MessageTypeId 	 MessageType 	`json:"messageTypeId" validate:"required,eq=4"`
	UniqueId      	 string      	`json:"uniqueId" validate:"required,max=36"`
	ErrorCode        ErrorCode   	`json:"errorCode" validate:"-"` //TODO: check if error is supported
	ErrorDescription string      	`json:"errorDescription" validate:"required"`
	ErrorDetails     interface{} 	`json:"errorDetails" validate:"omitempty"`
}

func (callError* CallError)GetMessageTypeId() MessageType {
	return callError.MessageTypeId
}

func (callError* CallError)GetUniqueId() string {
	return callError.UniqueId
}

func (callError *CallError) MarshalJSON() ([]byte, error) {
	fields := make([]interface{}, 5)
	fields[0] = callError.MessageTypeId
	fields[1] = callError.UniqueId
	fields[2] = callError.ErrorCode
	fields[3] = callError.ErrorDescription
	fields[4] = callError.ErrorDetails
	return ocppMessageToJson(callError)
}


// -------------------- Logic --------------------
func ParseRawJsonMessage(dataJson []byte) []interface{} {
	var arr []interface{}
	err := json.Unmarshal(dataJson, &arr)
	if err != nil {
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
	jsonData[len(jsonData) -1] = ']'
	return jsonData, nil
}

// -------------------- Endpoint --------------------
type Endpoint struct {
	Profiles []*Profile
	PendingRequests map[string]Request
}

func (endpoint *Endpoint)AddProfile(profile *Profile) {
	endpoint.Profiles = append(endpoint.Profiles, profile)
}

func (endpoint *Endpoint)GetProfile(name string) (*Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.Name == name {
			return p, true
		}
	}
	return nil, false
}

func (endpoint *Endpoint)GetProfileForFeature(featureName string) (*Profile, bool) {
	for _, p := range endpoint.Profiles {
		if p.SupportsFeature(featureName) {
			return p, true
		}
	}
	return nil, false
}

func (endpoint *Endpoint)AddPendingRequest(id string, request Request) {
	endpoint.PendingRequests[id] = request
}

func (endpoint *Endpoint)GetPendingRequest(id string) (Request, bool) {
	request, ok := endpoint.PendingRequests[id]
	return request, ok
}

func (endpoint *Endpoint)DeletePendingRequest(id string) {
	delete(endpoint.PendingRequests, id)
}

func (endpoint *Endpoint)ParseMessage(arr []interface{}) (Message, error) {
	// Checking message fields
	if len(arr) < 3 {
		log.Fatal("Invalid message. Expected array length >= 3")
	}
	typeId, ok := arr[0].(float64)
	if !ok {
		log.Printf("Invalid element %v at 0, expected int", arr[0])
	}
	uniqueId, ok := arr[1].(string)
	if !ok {
		log.Printf("Invalid element %v at 1, expected int", arr[1])
	}
	// Parse message
	if typeId == CALL {
		action := arr[2].(string)
		//TODO: check for ok in GetProfileForFeature
		profile, _ := endpoint.GetProfileForFeature(action)
		request := profile.ParseRequest(action, arr[3])
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
		return &call, nil
	} else if typeId == CALL_RESULT {
		request, ok := endpoint.PendingRequests[uniqueId]
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
			return nil, err
		}
		return &callResult, nil
	} else if typeId == CALL_ERROR {
		//TODO: handle error for pending request
		callError := CallError{
			MessageTypeId:    CALL_ERROR,
			UniqueId:         uniqueId,
			ErrorCode:        arr[2].(ErrorCode),
			ErrorDescription: arr[3].(string),
			ErrorDetails:     arr[4],
		}
		endpoint.DeletePendingRequest(callError.GetUniqueId())
		err := validate.Struct(callError)
		if err != nil {
			return nil, err
		}
		return &callError, nil
	} else {
		return nil, errors2.Errorf("Invalid message type ID %v", typeId)
	}
}

func (endpoint *Endpoint)CreateCall(request Request) (*Call, error) {
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
	endpoint.AddPendingRequest(uniqueId, request)
	return &call, nil
}

func (endpoint *Endpoint)CreateCallResult(confirmation Confirmation, uniqueId string) (*CallResult, error) {
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
	return &callResult, nil
}

func (endpoint *Endpoint)CreateCallError(uniqueId string, code ErrorCode, description string, details interface{}) *CallError {
	callError := CallError{
		MessageTypeId:    CALL_ERROR,
		UniqueId:         uniqueId,
		ErrorCode:        code,
		ErrorDescription: description,
		ErrorDetails:     details,
	}
	return &callError
}