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
	GetFeatureName() string
}

type Confirmation interface {
	Validatable
	GetFeatureName() string
}

// -------------------- Profile --------------------
type Profile struct {
	Features map[string]Feature
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
	Payload Request				`json:"payload"`
}

type CallResult struct {
	Message
	Payload Confirmation		`json:"payload"`
}

type CallError struct {
	Message
	ErrorCode ErrorCode			`json:"errorCode"`
	ErrorDescription string		`json:"errorDescription"`
	ErrorDetails interface{}	`json:"errorDetails"`
}

// -------------------- Global Variables --------------------
var Profiles []*Profile
var PendingRequests = map[string]Request{}


// -------------------- Logic --------------------
func AddProfile(profile *Profile) {
	Profiles = append(Profiles, profile)
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
		action := arr[2].(string)
		profile := GetProfileForFeature(action)
		request := profile.ParseRequest(action, arr[3])
		call := Call{
			Message: message,
			Action:  action,
			Payload: request,
		}
		return nil, call
	} else if typeId == CALL_RESULT {
		request, ok := PendingRequests[message.UniqueId]
		if !ok {
			log.Printf("No previous request %v sent. Discarding response message", message.UniqueId)
			return nil, nil
		}
		profile := GetProfileForFeature(request.GetFeatureName())
		confirmation := profile.ParseConfirmation(request.GetFeatureName(), arr[2])
		delete(PendingRequests, message.UniqueId)
		callResult := CallResult{
			Message: message,
			Payload: confirmation,
		}
		return nil, callResult
	} else if typeId == CALL_ERROR {
		//TODO: handle error for pending request
		delete(PendingRequests, message.UniqueId)
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

func GetProfileForFeature(featureName string) *Profile {
	for _, p := range Profiles {
		if p.SupportsFeature(featureName) {
			return p
		}
	}
	return nil
}