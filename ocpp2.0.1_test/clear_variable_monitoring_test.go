package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
)

// Test
func (suite *OcppV2TestSuite) TestClearVariableMonitoringRequestValidation() {
	var requestTable = []GenericTestEntry{
		{diagnostics.ClearVariableMonitoringRequest{ID: []int{0, 2, 15}}, true},
		{diagnostics.ClearVariableMonitoringRequest{ID: []int{0}}, true},
		{diagnostics.ClearVariableMonitoringRequest{ID: []int{}}, false},
		{diagnostics.ClearVariableMonitoringRequest{}, false},
		{diagnostics.ClearVariableMonitoringRequest{ID: []int{-1}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestClearVariableMonitoringConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{diagnostics.ClearVariableMonitoringResponse{ClearMonitoringResult: []diagnostics.ClearMonitoringResult{{ID: 2, Status: diagnostics.ClearMonitoringStatusAccepted}}}, true},
		{diagnostics.ClearVariableMonitoringResponse{ClearMonitoringResult: []diagnostics.ClearMonitoringResult{{ID: 2}}}, false},
		{diagnostics.ClearVariableMonitoringResponse{ClearMonitoringResult: []diagnostics.ClearMonitoringResult{}}, false},
		{diagnostics.ClearVariableMonitoringResponse{}, false},
		{diagnostics.ClearVariableMonitoringResponse{ClearMonitoringResult: []diagnostics.ClearMonitoringResult{{ID: -1, Status: diagnostics.ClearMonitoringStatusAccepted}}}, false},
		{diagnostics.ClearVariableMonitoringResponse{ClearMonitoringResult: []diagnostics.ClearMonitoringResult{{ID: 2, Status: "invalidClearMonitoringStatus"}}}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearVariableMonitoringE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	ids := []int{1, 2}
	result1 := diagnostics.ClearMonitoringResult{ID: 1, Status: diagnostics.ClearMonitoringStatusAccepted}
	result2 := diagnostics.ClearMonitoringResult{ID: 2, Status: diagnostics.ClearMonitoringStatusNotFound}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":[%v,%v]}]`, messageId, diagnostics.ClearVariableMonitoringFeatureName, ids[0], ids[1])
	responseJson := fmt.Sprintf(`[3,"%v",{"clearMonitoringResult":[{"id":%v,"status":"%v"},{"id":%v,"status":"%v"}]}]`, messageId, result1.ID, result1.Status, result2.ID, result2.Status)
	clearVariableMonitoringConfirmation := diagnostics.NewClearVariableMonitoringResponse([]diagnostics.ClearMonitoringResult{result1, result2})
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationDiagnosticsHandler{}
	handler.On("OnClearVariableMonitoring", mock.Anything).Return(clearVariableMonitoringConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*diagnostics.ClearVariableMonitoringRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Require().Len(request.ID, 2)
		suite.Equal(ids[0], request.ID[0])
		suite.Equal(ids[1], request.ID[1])
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearVariableMonitoring(wsId, func(confirmation *diagnostics.ClearVariableMonitoringResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Require().Len(confirmation.ClearMonitoringResult, 2)
		suite.Equal(result1.ID, confirmation.ClearMonitoringResult[0].ID)
		suite.Equal(result1.Status, confirmation.ClearMonitoringResult[0].Status)
		suite.Equal(result2.ID, confirmation.ClearMonitoringResult[1].ID)
		suite.Equal(result2.Status, confirmation.ClearMonitoringResult[1].Status)
		resultChannel <- true
	}, ids)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestClearVariableMonitoringInvalidEndpoint() {
	messageId := defaultMessageId
	ids := []int{1, 2}
	clearVariableMonitoringRequest := diagnostics.NewClearVariableMonitoringRequest(ids)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":[%v,%v]}]`, messageId, diagnostics.ClearVariableMonitoringFeatureName, ids[0], ids[1])
	testUnsupportedRequestFromChargingStation(suite, clearVariableMonitoringRequest, requestJson, messageId)
}
