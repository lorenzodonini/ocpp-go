package test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"github.com/lorenzodonini/go-ocpp/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
	"testing"
	"time"
)

// Utility functions
func GetBootNotificationRequest(t* testing.T, request ocpp.Request) *v16.BootNotificationRequest {
	assert.NotNil(t, request)
	result := request.(*v16.BootNotificationRequest)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.BootNotificationRequest{}, result)
	return result
}

func GetBootNotificationConfirmation(t* testing.T, confirmation ocpp.Confirmation) *v16.BootNotificationConfirmation {
	assert.NotNil(t, confirmation)
	result := confirmation.(*v16.BootNotificationConfirmation)
	assert.NotNil(t, result)
	assert.IsType(t, &v16.BootNotificationConfirmation{}, result)
	return result
}

var validate = validator.New()

// Tests
type CoreTestSuite struct {
	suite.Suite
	chargePoint *ocpp.ChargePoint
	centralSystem *ocpp.CentralSystem
	mockServer *MockWebsocketServer
	mockClient *MockWebsocketClient
}

func (suite *CoreTestSuite) SetupTest() {
	coreProfile := ocpp.NewProfile("core",  v16.BootNotificationFeature{})
	mockClient := MockWebsocketClient{}
	mockServer := MockWebsocketServer{}
	suite.mockClient = &mockClient
	suite.mockServer = &mockServer
	suite.chargePoint = ocpp.NewChargePoint("test_id", suite.mockClient, coreProfile)
	suite.centralSystem = ocpp.NewCentralSystem(suite.mockServer, coreProfile)
}

func (suite *CoreTestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []struct {
		request ocpp.Request
		expectedValid bool
	} {
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test"}, true},
		{v16.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, true},
		{v16.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{v16.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{v16.BootNotificationRequest{ChargeBoxSerialNumber: ">25.......................", ChargePointModel: "test", ChargePointVendor: "test"}, false},
		{v16.BootNotificationRequest{ChargePointModel: ">20..................", ChargePointVendor: "test"}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointSerialNumber: ">25.......................", ChargePointVendor: "test"}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: ">20.................."}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", FirmwareVersion: ">50................................................"}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Iccid: ">20.................."}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Imsi: ">20.................."}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterSerialNumber: ">25......................."}, false},
		{v16.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterType: ">25......................."}, false},
	}
	for _, testCase := range requestTable {
		 err := validate.Struct(testCase.request)
		 assert.Equal(t, testCase.expectedValid, err == nil)
	}
}

func (suite *CoreTestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []struct {
		confirmation ocpp.Confirmation
		expectedValid bool
	} {
		{v16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp.RegistrationStatusAccepted}, true},
		{v16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp.RegistrationStatusPending}, true},
		{v16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60, Status: ocpp.RegistrationStatusRejected}, true},
		{v16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: 60}, false},
		{v16.BootNotificationConfirmation{CurrentTime: time.Now(), Status: ocpp.RegistrationStatusAccepted}, false},
		{v16.BootNotificationConfirmation{Interval: 60, Status: ocpp.RegistrationStatusAccepted}, false},
		{v16.BootNotificationConfirmation{CurrentTime: time.Now(), Interval: -1, Status: ocpp.RegistrationStatusAccepted}, false},
		//TODO: incomplete list, see core.go
	}
	for _, testCase := range confirmationTable {
		err := validate.Struct(testCase.confirmation)
		assert.Equal(t, testCase.expectedValid, err == nil)
	}
}

func (suite *CoreTestSuite) TestBootNotificationRequestFromJson() {
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

func (suite *CoreTestSuite) TestBootNotificationRequestToJson() {
	t := suite.T()
	modelId := "model1"
	vendor := "ABL"
	request := v16.BootNotificationRequest{ChargePointModel: modelId, ChargePointVendor: vendor}
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

func (suite *CoreTestSuite) TestBootNotificationConfirmationFromJson() {
	t := suite.T()
	uniqueId := "5678"
	rawTime := time.Now().Format(ocpp.ISO8601)
	currentTime, err := time.Parse(ocpp.ISO8601, rawTime)
	assert.Nil(t, err)
	interval := 60
	status := ocpp.RegistrationStatusAccepted
	dummyRequest := v16.BootNotificationRequest{}
	dataJson := fmt.Sprintf(`[3,"%v",{"currentTime": "%v", "interval": 60, "status": "%v"}]`, uniqueId, currentTime.Format(ocpp.ISO8601), status)
	suite.chargePoint.Endpoint.AddPendingRequest(uniqueId, dummyRequest)
	callResult := ParseCallResult(&suite.chargePoint.Endpoint, dataJson, t)
	CheckCallResult(callResult, t, uniqueId)
	confirmation := GetBootNotificationConfirmation(t, callResult.Payload)
	assert.Equal(t, status, string(confirmation.Status))
	assert.Equal(t, interval, confirmation.Interval)
	assert.Equal(t, currentTime, confirmation.CurrentTime)
}

func (suite *CoreTestSuite) TestBootNotificationConfirmationToJson() {
	t := suite.T()
	uniqueId := "1234"
	now := time.Now()
	interval := 60
	status := ocpp.RegistrationStatusAccepted
	confirmation := v16.BootNotificationConfirmation{CurrentTime: now, Interval: interval, Status: ocpp.RegistrationStatus(status)}
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

func (suite *CoreTestSuite) TestBootNotificationInvalidMessage() {
	//TODO: implement
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

func TestBootNotification(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}
