package test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CoreTestSuite struct {
	suite.Suite
}

func (suite *CoreTestSuite) SetupTest() {
	coreProfile := ocpp.Profile{Features: make(map[string]ocpp.Feature)}
	feature := v16.BootNotificationFeature{}
	coreProfile.AddFeature(feature)
	ocpp.AddProfile(&coreProfile)
}

func GetBootNotificationRequest(t* testing.T, request ocpp.Request) *v16.BootNotificationRequest {
	assert.NotNil(t, request)
	result := request.(*v16.BootNotificationRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.BootNotificationRequest{}, result)
	return result
}

func (suite *CoreTestSuite) TestBootNotificationValid() {
	t := suite.T()
	dataJson := `[2,"1234","BootNotification",{"chargePointModel": "model1", "chargePointVendor": "ABL"}]`
	call := ParseCall(dataJson, t)
	CheckCall(call, t, v16.BootNotificationFeatureName, "1234")
	request := GetBootNotificationRequest(t, call.Payload)
	assert.Equal(t, "model1", request.ChargePointModel)
	assert.Equal(t, "ABL", request.ChargePointVendor)
}

func (suite *CoreTestSuite) TestBootNotificationInvalid() {

}

func (suite *CoreTestSuite) TestBootNotificationMessage() {

}

func (suite *CoreTestSuite) TestBootNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargePointModel": "model1", "chargePointVendor": "ABL"}]`, messageId, v16.BootNotificationFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime": "%v", "interval": 60, "status": "%v"}]`, messageId, time.Now().Format(ocpp.ISO8601), ocpp.RegistrationStatusAccepted)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	server := MockWebsocketServer{}
	client := MockWebsocketClient{}
	channel := MockWebSocket{id: wsId}
	// Setting server handlers
	server.SetNewClientHandler(func(ws ws.Channel) {
		assert.NotNil(t, ws)
		assert.Equal(t, wsId, ws.GetId())
	})
	server.SetMessageHandler(func(ws ws.Channel, data []byte) error {
		assert.Equal(t, requestRaw, data)
		jsonData := string(data)
		assert.Equal(t, requestJson, jsonData)
		call := ParseCall(jsonData, t)
		CheckCall(call, t, v16.BootNotificationFeatureName, messageId)
		ocpp.PendingRequests[messageId] = call.Payload
		// TODO: generate the response dynamically
		err := client.messageHandler(responseRaw)
		assert.Nil(t, err)
		return nil
	})
	// Setting client handlers
	client.On("Start", mock.AnythingOfType("string")).Return().Run(func(args mock.Arguments) {
		u := args.String(0)
		assert.Equal(t, wsUrl, u)
		server.newClientHandler(channel)
	})
	client.SetMessageHandler(func(data []byte) error {
		assert.Equal(t, responseRaw, data)
		jsonData := string(data)
		assert.Equal(t, responseJson, jsonData)
		callResult := ParseCallResult(jsonData, t)
		CheckCallResult(callResult, t, messageId)
		return nil
	})
	client.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		data := args.Get(0)
		bytes := data.([]byte)
		assert.NotNil(t, bytes)
		err := server.messageHandler(channel, bytes)
		assert.Nil(t, err)
	})
	// Test Run
	client.Start(wsUrl)
	client.Write(requestRaw)
}

func TestBootNotification(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}
