package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/go-ocpp/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, Info: "mockInfo", Status: ocpp16.ChargePointStatusAvailable, Timestamp: ocpp16.DateTime{Time: time.Now().Add(-1 * time.Hour)}, VendorId: "mockId", VendorErrorCode: "mockErrorCode"}, true},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, Status: ocpp16.ChargePointStatusAvailable}, true},
		{ocpp16.StatusNotificationRequest{ErrorCode: ocpp16.NoError, Status: ocpp16.ChargePointStatusAvailable}, true},
		{ocpp16.StatusNotificationRequest{ConnectorId: -1, ErrorCode: ocpp16.NoError, Status: ocpp16.ChargePointStatusAvailable}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, Status: ocpp16.ChargePointStatusAvailable}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: "invalidErrorCode", Status: ocpp16.ChargePointStatusAvailable}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, Status: "invalidChargePointStatus"}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, Info: ">50................................................", Status: ocpp16.ChargePointStatusAvailable}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, VendorErrorCode: ">50................................................", Status: ocpp16.ChargePointStatusAvailable}, false},
		{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, VendorId: ">255............................................................................................................................................................................................................................................................", Status: ocpp16.ChargePointStatusAvailable}, false},
		//{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, Info: "mockInfo", Status: ocpp16.ChargePointStatusAvailable, Timestamp: ocpp16.DateTime{Time: time.Now().Add(1 * time.Hour)}, VendorId: "mockId", VendorErrorCode: "mockErrorCode"}, false},
	}
	ExecuteRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestStatusNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry{
		{ocpp16.StatusNotificationConfirmation{}, true},
	}
	ExecuteConfirmationTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestStatusNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	timestamp := ocpp16.DateTime{Time: time.Now().Add(-1 * time.Hour)}
	status := ocpp16.ChargePointStatusAvailable
	cpErrorCode := ocpp16.NoError
	info := "mockInfo"
	vendorId := "mockVendorId"
	vendorErrorCode := "mockErrorCode"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"errorCode":"%v","info":"%v","status":"%v","timestamp":"%v","vendorId":"%v","vendorErrorCode":"%v"}]`, messageId, ocpp16.StatusNotificationFeatureName, connectorId, cpErrorCode, info, status, timestamp.Format(ocpp16.ISO8601), vendorId, vendorErrorCode)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	statusNotificationConfirmation := ocpp16.NewStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(statusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*ocpp16.StatusNotificationRequest)
		assert.True(t, ok)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, cpErrorCode, request.ErrorCode)
		assert.Equal(t, status, request.Status)
		assert.Equal(t, info, request.Info)
		assert.Equal(t, vendorId, request.VendorId)
		assert.Equal(t, vendorErrorCode, request.VendorErrorCode)
		assertDateTimeEquality(t, timestamp, request.Timestamp)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.StatusNotification(connectorId, cpErrorCode, status, func(request *ocpp16.StatusNotificationRequest) {
		request.Timestamp = timestamp
		request.Info = info
		request.VendorId = vendorId
		request.VendorErrorCode = vendorErrorCode
	})
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
}

func (suite *OcppV16TestSuite) TestStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	timestamp := ocpp16.DateTime{Time: time.Now().Add(-1 * time.Hour)}
	status := ocpp16.ChargePointStatusAvailable
	cpErrorCode := ocpp16.NoError
	info := "mockInfo"
	vendorId := "mockVendorId"
	vendorErrorCode := "mockErrorCode"
	statusNotificationRequest := ocpp16.NewStatusNotificationRequest(connectorId, cpErrorCode, status)
	statusNotificationRequest.Info = info
	statusNotificationRequest.Timestamp = timestamp
	statusNotificationRequest.VendorId = vendorId
	statusNotificationRequest.VendorErrorCode = vendorErrorCode
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"errorCode":"%v","info":"%v","status":"%v","timestamp":"%v","vendorId":"%v","vendorErrorCode":"%v"}]`, messageId, ocpp16.StatusNotificationFeatureName, connectorId, cpErrorCode, info, status, timestamp.Format(ocpp16.ISO8601), vendorId, vendorErrorCode)
	testUnsupportedRequestFromCentralSystem(suite, statusNotificationRequest, requestJson, messageId)
}
