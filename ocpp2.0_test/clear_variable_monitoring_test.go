package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test
func (suite *OcppV2TestSuite) TestClearVariableMonitoringRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.ClearVariableMonitoringRequest{ID: []int{0, 2, 15}}, true},
		{ocpp2.ClearVariableMonitoringRequest{ID: []int{0}}, true},
		{ocpp2.ClearVariableMonitoringRequest{ID: []int{}}, false},
		{ocpp2.ClearVariableMonitoringRequest{}, false},
		{ocpp2.ClearVariableMonitoringRequest{ID: []int{-1}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestClearVariableMonitoringConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.ClearVariableMonitoringConfirmation{ClearMonitoringResult: []ocpp2.ClearMonitoringResult{{ID: 2, Status: ocpp2.ClearMonitoringStatusAccepted}}}, true},
		{ocpp2.ClearVariableMonitoringConfirmation{ClearMonitoringResult: []ocpp2.ClearMonitoringResult{{ID: 2}}}, false},
		{ocpp2.ClearVariableMonitoringConfirmation{ClearMonitoringResult: []ocpp2.ClearMonitoringResult{}}, false},
		{ocpp2.ClearVariableMonitoringConfirmation{}, false},
		{ocpp2.ClearVariableMonitoringConfirmation{ClearMonitoringResult: []ocpp2.ClearMonitoringResult{{ID: -1, Status: ocpp2.ClearMonitoringStatusAccepted}}}, false},
		{ocpp2.ClearVariableMonitoringConfirmation{ClearMonitoringResult: []ocpp2.ClearMonitoringResult{{ID: 2, Status: "invalidClearMonitoringStatus"}}}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestClearVariableMonitoringE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	ids := []int{1,2}
	result1 := ocpp2.ClearMonitoringResult{ID: 1, Status: ocpp2.ClearMonitoringStatusAccepted}
	result2 := ocpp2.ClearMonitoringResult{ID: 2, Status: ocpp2.ClearMonitoringStatusNotFound}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":[%v,%v]}]`, messageId, ocpp2.ClearVariableMonitoringFeatureName, ids[0], ids[1])
	responseJson := fmt.Sprintf(`[3,"%v",{"clearMonitoringResult":[{"id":%v,"status":"%v"},{"id":%v,"status":"%v"}]}]`, messageId, result1.ID, result1.Status, result2.ID, result2.Status)
	clearVariableMonitoringConfirmation := ocpp2.NewClearVariableMonitoringConfirmation([]ocpp2.ClearMonitoringResult{result1, result2})
	channel := NewMockWebSocket(wsId)

	coreListener := MockChargePointCoreListener{}
	coreListener.On("OnClearVariableMonitoring", mock.Anything).Return(clearVariableMonitoringConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp2.ClearVariableMonitoringRequest)
		require.True(t, ok)
		require.NotNil(t, request)
		require.Len(t, request.ID, 2)
		assert.Equal(t, ids[0], request.ID[0])
		assert.Equal(t, ids[1], request.ID[1])
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, coreListener, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.ClearVariableMonitoring(wsId, func(confirmation *ocpp2.ClearVariableMonitoringConfirmation, err error) {
		require.Nil(t, err)
		require.NotNil(t, confirmation)
		require.Len(t, confirmation.ClearMonitoringResult, 2)
		assert.Equal(t, result1.ID, confirmation.ClearMonitoringResult[0].ID)
		assert.Equal(t, result1.Status, confirmation.ClearMonitoringResult[0].Status)
		assert.Equal(t, result2.ID, confirmation.ClearMonitoringResult[1].ID)
		assert.Equal(t, result2.Status, confirmation.ClearMonitoringResult[1].Status)
		resultChannel <- true
	}, ids)
	require.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV2TestSuite) TestClearVariableMonitoringInvalidEndpoint() {
	messageId := defaultMessageId
	ids := []int{1,2}
	ClearVariableMonitoringRequest := ocpp2.NewClearVariableMonitoringRequest(ids)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"id":[%v,%v]}]`, messageId, ocpp2.ClearVariableMonitoringFeatureName, ids[0], ids[1])
	testUnsupportedRequestFromChargePoint(suite, ClearVariableMonitoringRequest, requestJson, messageId)
}
