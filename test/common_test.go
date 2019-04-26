package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ParseCall(json string, t* testing.T) *ocpp.Call {
	parsedData := ocpp.ParseJsonMessage(json)
	err, result := ocpp.ParseMessage(parsedData)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	call, ok := result.(ocpp.Call)
	assert.Equal(t, true, ok)
	assert.NotNil(t, call)
	return &call
}

func CheckCall(call* ocpp.Call, t *testing.T, expectedAction string, expectedId string) {
	assert.Equal(t, ocpp.CALL, call.MessageTypeId)
	assert.Equal(t, expectedAction, call.Action)
	assert.Equal(t, expectedId, call.UniqueId)
	assert.NotNil(t, call.Payload)
}
