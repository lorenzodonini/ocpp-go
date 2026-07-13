package ocppj

import (
	"bytes"
	"encoding/json"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"reflect"
)

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

// Marshals data by manipulating EscapeHTML property of encoder
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(EscapeHTML)
	err := encoder.Encode(t)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

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
