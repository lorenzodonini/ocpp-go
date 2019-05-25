package test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6/core"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Utility functions
func GetBootNotificationRequest(t* testing.T, request ocpp.Request) *core.BootNotificationRequest {
	assert.NotNil(t, request)
	result := request.(*core.BootNotificationRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &core.BootNotificationRequest{}, result)
	return result
}

func GetBootNotificationConfirmation(t* testing.T, confirmation ocpp.Confirmation) *core.BootNotificationConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*core.BootNotificationConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &core.BootNotificationConfirmation{}, result)
	return result
}

// Tests
func (suite *OcppTestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test"}, true},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, true},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: ">25.......................", ChargePointModel: "test", ChargePointVendor: "test"}, false},
		{core.BootNotificationRequest{ChargePointModel: ">20..................", ChargePointVendor: "test"}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointSerialNumber: ">25.......................", ChargePointVendor: "test"}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: ">20.................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", FirmwareVersion: ">50................................................"}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Iccid: ">20.................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Imsi: ">20.................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterSerialNumber: ">25......................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterType: ">25......................."}, false},
	}
	executeRequestTestTable(t, requestTable)
}

func (suite *OcppTestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{core.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp.RegistrationStatusAccepted}, true},
		{core.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp.RegistrationStatusPending}, true},
		{core.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp.RegistrationStatusRejected}, true},
		{core.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60}, false},
		{core.BootNotificationConfirmation{CurrentTime: time.Now(), Status: ocpp.RegistrationStatusAccepted}, false},
		{core.BootNotificationConfirmation{Interval: 60, Status: ocpp.RegistrationStatusAccepted}, false},
		{core.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: -1, Status: ocpp.RegistrationStatusAccepted}, false},
		//TODO: incomplete list, see core.go
	}
	executeConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppTestSuite) TestBootNotificationRequestFromJson() {
	t := suite.T()
	uniqueId := "1234"
	modelId := "model1"
	vendor := "ABL"
	dataJson := fmt.Sprintf(`[2,"%v","BootNotification",{"chargePointModel": "%v", "chargePointVendor": "%v"}]`, uniqueId, modelId, vendor)
	call := ParseCall(&suite.centralSystem.Endpoint, dataJson, t)
	CheckCall(call, t, v16.BootNotificationFeatureName, uniqueId)
	request := GetBootNotificationRequest(t, call.Payload)
	assert.Equal(t, modelId, request.ChargePointModel)
	assert.Equal(t, vendor, request.ChargePointVendor)
}

func (suite *OcppTestSuite) TestBootNotificationRequestToJson() {
	t := suite.T()
	modelId := "model1"
	vendor := "ABL"
	request := core.BootNotificationRequest{ChargePointModel: modelId, ChargePointVendor: vendor}
	call, err := suite.chargePoint.CreateCall(request)
	uniqueId := call.GetUniqueId()
	assert.Nil(t, err)
	assert.NotNil(t, call)
	err = validate.Struct(call)
	assert.Nil(t, err)
	jsonData, err := call.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[2,"%v","BootNotification",{"chargePointModel":"%v","chargePointVendor":"%v"}]`, uniqueId, modelId, vendor)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppTestSuite) TestBootNotificationConfirmationFromJson() {
	t := suite.T()
	uniqueId := "5678"
	rawTime := time.Now().Format(ocpp.ISO8601)
	currentTime, err := time.Parse(ocpp.ISO8601, rawTime)
	assert.Nil(t, err)
	interval := 60
	status := ocpp.RegistrationStatusAccepted
	dummyRequest := core.BootNotificationRequest{}
	dataJson := fmt.Sprintf(`[3,"%v",{"currentTime": "%v", "interval": 60, "status": "%v"}]`, uniqueId, currentTime.Format(ocpp.ISO8601), status)
	suite.chargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	callResult := ParseCallResult(&suite.chargePoint.Endpoint, dataJson, t)
	CheckCallResult(callResult, t, uniqueId)
	confirmation := GetBootNotificationConfirmation(t, callResult.Payload)
	assert.Equal(t, status, string(confirmation.Status))
	assert.Equal(t, interval, confirmation.Interval)
	assert.Equal(t, currentTime, confirmation.CurrentTime)
}

func (suite *OcppTestSuite) TestBootNotificationConfirmationToJson() {
	t := suite.T()
	uniqueId := "1234"
	now := time.Now()
	interval := 60
	status := ocpp.RegistrationStatusAccepted
	confirmation := core.BootNotificationConfirmation{CurrentTime: now, Interval: interval, Status: ocpp.RegistrationStatus(status)}
	callResult, err := suite.centralSystem.CreateCallResult(confirmation, uniqueId)
	assert.Nil(t, err)
	assert.NotNil(t, callResult)
	err = validate.Struct(callResult)
	assert.Nil(t, err)
	jsonData, err := callResult.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v","interval":60,"status":"%v"}]`, uniqueId, now.Format(time.RFC3339Nano), status)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppTestSuite) TestBootNotificationInvalidMessage() {
	//TODO: implement
}

func (suite *OcppTestSuite) TestBootNotificationE2EMocked() {
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
		call := ParseCall(&suite.chargePoint.Endpoint, jsonData, t)
		CheckCall(call, t, v16.BootNotificationFeatureName, messageId)
		suite.chargePoint.AddPendingRequest(messageId, call.Payload)
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
		callResult := ParseCallResult(&suite.chargePoint.Endpoint, jsonData, t)
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
	err := client.Start(wsUrl)
	assert.Nil(t, err)
	client.Write(requestRaw)
}
