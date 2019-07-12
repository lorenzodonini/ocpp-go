package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Utility functions
func GetBootNotificationRequest(t *testing.T, request ocppj.Request) *ocpp16.BootNotificationRequest {
	assert.NotNil(t, request)
	result := request.(*ocpp16.BootNotificationRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.BootNotificationRequest{}, result)
	return result
}

func GetBootNotificationConfirmation(t *testing.T, confirmation ocppj.Confirmation) *ocpp16.BootNotificationConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*ocpp16.BootNotificationConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &ocpp16.BootNotificationConfirmation{}, result)
	return result
}

// Tests
func (suite *OcppV16TestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test"}, true},
		{ocpp16.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, true},
		{ocpp16.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{ocpp16.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{ocpp16.BootNotificationRequest{ChargeBoxSerialNumber: ">25.......................", ChargePointModel: "test", ChargePointVendor: "test"}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: ">20..................", ChargePointVendor: "test"}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointSerialNumber: ">25.......................", ChargePointVendor: "test"}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: ">20.................."}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", FirmwareVersion: ">50................................................"}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Iccid: ">20.................."}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Imsi: ">20.................."}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterSerialNumber: ">25......................."}, false},
		{ocpp16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterType: ">25......................."}, false},
	}
	ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp16.RegistrationStatusAccepted}, true},
		{ocpp16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp16.RegistrationStatusPending}, true},
		{ocpp16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp16.RegistrationStatusRejected}, true},
		{ocpp16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60}, false},
		{ocpp16.BootNotificationConfirmation{CurrentTime: time.Now(), Status: ocpp16.RegistrationStatusAccepted}, false},
		{ocpp16.BootNotificationConfirmation{Interval: 60, Status: ocpp16.RegistrationStatusAccepted}, false},
		{ocpp16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: -1, Status: ocpp16.RegistrationStatusAccepted}, false},
		//TODO: incomplete list, see core.go
	}
	ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestBootNotificationRequestFromJson() {
	t := suite.T()
	uniqueId := "1234"
	modelId := "model1"
	vendor := "ABL"
	dataJson := fmt.Sprintf(`[2,"%v","BootNotification",{"chargePointModel": "%v", "chargePointVendor": "%v"}]`, uniqueId, modelId, vendor)
	call := ParseCall(&suite.centralSystem.Endpoint, dataJson, t)
	CheckCall(call, t, ocpp16.BootNotificationFeatureName, uniqueId)
	request := GetBootNotificationRequest(t, call.Payload)
	assert.Equal(t, modelId, request.ChargePointModel)
	assert.Equal(t, vendor, request.ChargePointVendor)
}

func (suite *OcppV16TestSuite) TestBootNotificationRequestToJson() {
	t := suite.T()
	modelId := "model1"
	vendor := "ABL"
	request := ocpp16.BootNotificationRequest{ChargePointModel: modelId, ChargePointVendor: vendor}
	call, err := suite.chargePoint.CreateCall(request)
	uniqueId := call.GetUniqueId()
	assert.Nil(t, err)
	assert.NotNil(t, call)
	err = Validate.Struct(call)
	assert.Nil(t, err)
	jsonData, err := call.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[2,"%v","BootNotification",{"chargePointModel":"%v","chargePointVendor":"%v"}]`, uniqueId, modelId, vendor)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppV16TestSuite) TestBootNotificationConfirmationFromJson() {
	t := suite.T()
	uniqueId := "5678"
	rawTime := time.Now().Format(ocpp16.ISO8601)
	currentTime, err := time.Parse(ocpp16.ISO8601, rawTime)
	assert.Nil(t, err)
	interval := 60
	status := ocpp16.RegistrationStatusAccepted
	dummyRequest := ocpp16.BootNotificationRequest{}
	dataJson := fmt.Sprintf(`[3,"%v",{"currentTime": "%v", "interval": 60, "status": "%v"}]`, uniqueId, currentTime.Format(ocpp16.ISO8601), status)
	suite.chargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	callResult := ParseCallResult(&suite.chargePoint.Endpoint, dataJson, t)
	CheckCallResult(callResult, t, uniqueId)
	confirmation := GetBootNotificationConfirmation(t, callResult.Payload)
	assert.Equal(t, status, confirmation.Status)
	assert.Equal(t, interval, confirmation.Interval)
	assert.Equal(t, currentTime, confirmation.CurrentTime)
}

func (suite *OcppV16TestSuite) TestBootNotificationConfirmationToJson() {
	t := suite.T()
	uniqueId := "1234"
	now := time.Now()
	interval := 60
	status := ocpp16.RegistrationStatusAccepted
	confirmation := ocpp16.BootNotificationConfirmation{CurrentTime: now, Interval: interval, Status: ocpp16.RegistrationStatus(status)}
	callResult, err := suite.centralSystem.CreateCallResult(confirmation, uniqueId)
	assert.Nil(t, err)
	assert.NotNil(t, callResult)
	err = Validate.Struct(callResult)
	assert.Nil(t, err)
	jsonData, err := callResult.MarshalJSON()
	assert.Nil(t, err)
	assert.NotNil(t, jsonData)
	expectedJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v","interval":60,"status":"%v"}]`, uniqueId, now.Format(time.RFC3339Nano), status)
	assert.Equal(t, []byte(expectedJson), jsonData)
}

func (suite *OcppV16TestSuite) TestBootNotificationInvalidMessage() {
	//TODO: implement
}

func (suite *OcppV16TestSuite) TestBootNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargePointModel": "model1", "chargePointVendor": "ABL"}]`, messageId, ocpp16.BootNotificationFeatureName)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime": "%v", "interval": 60, "status": "%v"}]`, messageId, time.Now().Format(ocpp16.ISO8601), ocpp16.RegistrationStatusAccepted)
	requestRaw := []byte(requestJson)
	responseRaw := []byte(responseJson)
	channel := NewMockWebSocket(wsId)
	// Setting server handlers
	suite.mockServer.SetNewClientHandler(func(ws ws.Channel) {
		assert.NotNil(t, ws)
		assert.Equal(t, wsId, ws.GetId())
	})
	suite.mockServer.SetMessageHandler(func(ws ws.Channel, data []byte) error {
		assert.Equal(t, requestRaw, data)
		jsonData := string(data)
		assert.Equal(t, requestJson, jsonData)
		call := ParseCall(&suite.chargePoint.Endpoint, jsonData, t)
		CheckCall(call, t, ocpp16.BootNotificationFeatureName, messageId)
		suite.chargePoint.AddPendingRequest(messageId, call.Payload)
		// TODO: generate the response dynamically
		err := suite.mockClient.MessageHandler(responseRaw)
		assert.Nil(t, err)
		return nil
	})
	// Setting client handlers
	suite.mockClient.On("Start", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		u := args.String(0)
		assert.Equal(t, wsUrl, u)
		suite.mockServer.NewClientHandler(channel)
	})
	suite.mockClient.SetMessageHandler(func(data []byte) error {
		assert.Equal(t, responseRaw, data)
		jsonData := string(data)
		assert.Equal(t, responseJson, jsonData)
		callResult := ParseCallResult(&suite.chargePoint.Endpoint, jsonData, t)
		CheckCallResult(callResult, t, messageId)
		return nil
	})
	suite.mockClient.On("Write", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(0)
		bytes := data.([]byte)
		assert.NotNil(t, bytes)
		err := suite.mockServer.MessageHandler(channel, bytes)
		assert.Nil(t, err)
	})
	// Test Run
	err := suite.mockClient.Start(wsUrl)
	assert.Nil(t, err)
	err = suite.mockClient.Write(requestRaw)
	assert.Nil(t, err)
}
