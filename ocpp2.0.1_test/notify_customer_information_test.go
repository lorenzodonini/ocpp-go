package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestNotifyCustomerInformationRequestValidation() {
	var requestTable = []GenericTestEntry{
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", Tbc: false, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", Tbc: true, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", GeneratedAt: types.DateTime{Time: time.Now()}}, true},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData"}, true},
		{diagnostics.NotifyCustomerInformationRequest{}, false},
		{diagnostics.NotifyCustomerInformationRequest{Data: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Tbc: false, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, false},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", Tbc: false, SeqNo: -1, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, false},
		{diagnostics.NotifyCustomerInformationRequest{Data: "dummyData", Tbc: false, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: -1}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestNotifyCustomerInformationConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{diagnostics.NotifyCustomerInformationResponse{}, true},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestNotifyCustomerInformationE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	data := "dummyData"
	tbc := false
	seqNo := 0
	generatedAt := types.DateTime{Time: time.Now()}
	requestID := 42
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"data":"%v","seqNo":%v,"generatedAt":"%v","requestId":%v}]`,
		messageId, diagnostics.NotifyCustomerInformationFeatureName, data, seqNo, generatedAt.FormatTimestamp(), requestID)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := diagnostics.NewNotifyCustomerInformationResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSDiagnosticsHandler{}
	handler.On("OnNotifyCustomerInformation", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*diagnostics.NotifyCustomerInformationRequest)
		suite.Require().True(ok)
		suite.Equal(data, request.Data)
		suite.Equal(tbc, request.Tbc)
		suite.Equal(seqNo, request.SeqNo)
		assertDateTimeEquality(suite, &generatedAt, &request.GeneratedAt)
		suite.Equal(requestID, request.RequestID)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	r, err := suite.chargingStation.NotifyCustomerInformation(data, seqNo, generatedAt, requestID)
	suite.Nil(err)
	suite.NotNil(r)
}

func (suite *OcppV2TestSuite) TestNotifyCustomerInformationInvalidEndpoint() {
	messageId := defaultMessageId
	data := "dummyData"
	seqNo := 0
	generatedAt := types.DateTime{Time: time.Now()}
	requestID := 42
	req := diagnostics.NewNotifyCustomerInformationRequest(data, seqNo, generatedAt, requestID)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"data":"%v","seqNo":%v,"generatedAt":"%v","requestId":%v}]`,
		messageId, diagnostics.NotifyCustomerInformationFeatureName, data, seqNo, generatedAt.FormatTimestamp(), requestID)
	testUnsupportedRequestFromCentralSystem(suite, req, requestJson, messageId)
}
