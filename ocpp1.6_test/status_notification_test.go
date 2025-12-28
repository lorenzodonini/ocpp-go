package ocpp16_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// Test
func (suite *OcppV16TestSuite) TestStatusNotificationRequestValidation() {
	var requestTable = []GenericTestEntry{
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError, Info: "mockInfo", Status: core.ChargePointStatusAvailable, Timestamp: types.NewDateTime(time.Now()), VendorId: "mockId", VendorErrorCode: "mockErrorCode"}, true},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError, Status: core.ChargePointStatusAvailable}, true},
		{core.StatusNotificationRequest{ErrorCode: core.NoError, Status: core.ChargePointStatusAvailable}, true},
		{core.StatusNotificationRequest{ConnectorId: -1, ErrorCode: core.NoError, Status: core.ChargePointStatusAvailable}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, Status: core.ChargePointStatusAvailable}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: "invalidErrorCode", Status: core.ChargePointStatusAvailable}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError, Status: "invalidChargePointStatus"}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError, Info: ">50................................................", Status: core.ChargePointStatusAvailable}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError, VendorErrorCode: ">50................................................", Status: core.ChargePointStatusAvailable}, false},
		{core.StatusNotificationRequest{ConnectorId: 0, ErrorCode: core.NoError, VendorId: ">255............................................................................................................................................................................................................................................................", Status: core.ChargePointStatusAvailable}, false},
		//{ocpp16.StatusNotificationRequest{ConnectorId: 0, ErrorCode: ocpp16.NoError, Info: "mockInfo", Status: ocpp16.ChargePointStatusAvailable, Timestamp: ocpp16.DateTime{Time: time.Now().Add(1 * time.Hour)}, VendorId: "mockId", VendorErrorCode: "mockErrorCode"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV16TestSuite) TestStatusNotificationConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{core.StatusNotificationConfirmation{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV16TestSuite) TestStatusNotificationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	status := core.ChargePointStatusAvailable
	cpErrorCode := core.NoError
	info := "mockInfo"
	vendorId := "mockVendorId"
	vendorErrorCode := "mockErrorCode"
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"errorCode":"%v","info":"%v","status":"%v","timestamp":"%v","vendorId":"%v","vendorErrorCode":"%v"}]`, messageId, core.StatusNotificationFeatureName, connectorId, cpErrorCode, info, status, timestamp.FormatTimestamp(), vendorId, vendorErrorCode)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	statusNotificationConfirmation := core.NewStatusNotificationConfirmation()
	channel := NewMockWebSocket(wsId)

	coreListener := &MockCentralSystemCoreListener{}
	coreListener.On("OnStatusNotification", mock.AnythingOfType("string"), mock.Anything).Return(statusNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*core.StatusNotificationRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(connectorId, request.ConnectorId)
		suite.Equal(cpErrorCode, request.ErrorCode)
		suite.Equal(status, request.Status)
		suite.Equal(info, request.Info)
		suite.Equal(vendorId, request.VendorId)
		suite.Equal(vendorErrorCode, request.VendorErrorCode)
		assertDateTimeEquality(suite, *timestamp, *request.Timestamp)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	confirmation, err := suite.chargePoint.StatusNotification(connectorId, cpErrorCode, status, func(request *core.StatusNotificationRequest) {
		request.Timestamp = timestamp
		request.Info = info
		request.VendorId = vendorId
		request.VendorErrorCode = vendorErrorCode
	})
	suite.Require().Nil(err)
	suite.Require().NotNil(confirmation)
}

func (suite *OcppV16TestSuite) TestStatusNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	connectorId := 1
	timestamp := types.NewDateTime(time.Now())
	status := core.ChargePointStatusAvailable
	cpErrorCode := core.NoError
	info := "mockInfo"
	vendorId := "mockVendorId"
	vendorErrorCode := "mockErrorCode"
	statusNotificationRequest := core.NewStatusNotificationRequest(connectorId, cpErrorCode, status)
	statusNotificationRequest.Info = info
	statusNotificationRequest.Timestamp = timestamp
	statusNotificationRequest.VendorId = vendorId
	statusNotificationRequest.VendorErrorCode = vendorErrorCode
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"errorCode":"%v","info":"%v","status":"%v","timestamp":"%v","vendorId":"%v","vendorErrorCode":"%v"}]`, messageId, core.StatusNotificationFeatureName, connectorId, cpErrorCode, info, status, timestamp.FormatTimestamp(), vendorId, vendorErrorCode)
	testUnsupportedRequestFromCentralSystem(suite, statusNotificationRequest, requestJson, messageId)
}
