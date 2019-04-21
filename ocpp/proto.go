package ocpp

import (
	"encoding/json"
	"log"
	"reflect"
)

type Validatable interface {
	validate() error
}

type Feature interface {
	GetFeatureName() string
	GetRequestType() reflect.Type
	GetConfirmationType() reflect.Type
}

type Request interface {
	Validatable
	GetFeature() Feature
}

type Confirmation interface {
	Validatable
	GetFeature() Feature
}

/*
 Profile
 */
type Profile struct {
	Features map[string]Feature
}

func (p* Profile) SupportsFeature(name string) bool {
	_, ok := p.Features[name]
	return ok
}

func (p* Profile) ParseRequest(featureName string, rawRequest interface{}) *Request {
	feature, ok := p.Features[featureName]
	if !ok {
		log.Printf("Feature %s not found", featureName)
		return nil
	}
	requestType := feature.GetRequestType()
	bytes, ok := rawRequest.([]byte)
	if !ok {
		log.Printf("Couldn't cast raw request to bytes")
		return nil
	}
	request := reflect.New(requestType)
	err := json.Unmarshal(bytes, &request)
	if err != nil {
		log.Printf("Error while parsing json %v", err)
	}
	return request.Interface().(*Request)
}

func (p* Profile) ParseConfirmation(featureName string, rawConfirmation interface{}) *Confirmation {
	feature, ok := p.Features[featureName]
	if !ok {
		log.Printf("Feature %s not found", featureName)
		return nil
	}
	requestType := feature.GetConfirmationType()
	bytes, ok := rawConfirmation.([]byte)
	if !ok {
		log.Printf("Couldn't cast raw confirmation to bytes")
		return nil
	}
	confirmation := reflect.New(requestType)
	err := json.Unmarshal(bytes, &confirmation)
	if err != nil {
		log.Printf("Error while parsing json %v", err)
	}
	return confirmation.Interface().(*Confirmation)
}

var Profiles []Profile

/*
 Message
 */
type MessageType int

const (
	CALL 		= 2
	CALL_RESULT = 3
	CALL_ERROR 	= 4
)

type Message struct {
	MessageTypeId MessageType	`json:"messageTypeId"`
	UniqueId string 			`json:"uniqueId"`	//Max 36 chars
	Validatable
}

func (m* Message) validate() error {
	return nil
}

type Call struct {
	Message
	Action string				`json:"action"`
	Payload* Request			`json:"payload"`
}

type CallResult struct {
	Message
	Payload* Confirmation		`json:"payload"`
}

type CallError struct {
	Message
	ErrorCode ErrorCode			`json:"errorCode"`
	ErrorDescription string		`json:"errorDescription"`
	ErrorDetails interface{}	`json:"errorDetails"`
}

func ParseJsonMessage(dataJson string) []interface{} {
	var arr []interface{}
	err := json.Unmarshal([]byte(dataJson), &arr)
	if err != nil {
		log.Fatal(err)
	}
	return arr
}

func ParseMessage(arr []interface{}) (error, interface{}) {
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
	message := Message{MessageTypeId: MessageType(typeId), UniqueId: uniqueId}
	if typeId == CALL {
		call := Call{
			Message: message,
			Action:  arr[2].(string),
			Payload: nil,
		}
		return nil, call
	} else if typeId == CALL_RESULT {
		callResult := CallResult{
			Message: message,
			Payload: nil,
		}
		return nil, callResult
	} else if typeId == CALL_ERROR {
		callError := CallError{
			Message: message,
			ErrorCode: arr[2].(ErrorCode),
			ErrorDescription: arr[3].(string),
			ErrorDetails: arr[4],
		}
		return nil, callError
	} else {
		//TODO: return custom error
		return nil, nil
	}
}