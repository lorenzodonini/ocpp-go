package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test
func (suite *OcppV2TestSuite) TestGet15118EVCertificateRequestValidation() {
	var requestTable = []GenericTestEntry{
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: iso15118.CertificateActionInstall, ExiRequest: "deadbeef"}, true},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: iso15118.CertificateActionUpdate, ExiRequest: "deadbeef"}, true},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: iso15118.CertificateActionInstall}, false},
		{iso15118.Get15118EVCertificateRequest{ExiRequest: "deadbeef"}, false},
		{iso15118.Get15118EVCertificateRequest{}, false},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: ">50................................................", Action: iso15118.CertificateActionInstall, ExiRequest: "deadbeef"}, false},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: "invalidCertificateAction", ExiRequest: "deadbeef"}, false},
		{iso15118.Get15118EVCertificateRequest{SchemaVersion: "1.0", Action: iso15118.CertificateActionInstall, ExiRequest: newLongString(5601)}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateConfirmationValidation() {
	var confirmationTable = []GenericTestEntry{
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef"}, true},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted}, false},
		{iso15118.Get15118EVCertificateResponse{ExiResponse: "deadbeef"}, false},
		{iso15118.Get15118EVCertificateResponse{}, false},
		{iso15118.Get15118EVCertificateResponse{Status: "invalidCertificateStatus", ExiResponse: "deadbeef", StatusInfo: types.NewStatusInfo("200", "ok")}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: newLongString(5601), StatusInfo: types.NewStatusInfo("200", "ok")}, false},
		{iso15118.Get15118EVCertificateResponse{Status: types.Certificate15188EVStatusAccepted, ExiResponse: "deadbeef", StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	status := types.Certificate15188EVStatusAccepted
	schemaVersion := "1.0"
	action := iso15118.CertificateActionInstall
	exiRequest := "deadbeef"
	exiResponse := "deadbeef2"
	statusInfo := types.NewStatusInfo("200", "ok")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"iso15118SchemaVersion":"%v","action":"%v","exiRequest":"%v"}]`,
		messageId, iso15118.Get15118EVCertificateFeatureName, schemaVersion, action, exiRequest)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","exiResponse":"%v","statusInfo":{"reasonCode":"%v","additionalInfo":"%v"}}]`,
		messageId, status, exiResponse, statusInfo.ReasonCode, statusInfo.AdditionalInfo)
	get15118EVCertificateResponse := iso15118.NewGet15118EVCertificateResponse(status, exiResponse)
	get15118EVCertificateResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSIso15118Handler{}
	handler.On("OnGet15118EVCertificate", mock.AnythingOfType("string"), mock.Anything).Return(get15118EVCertificateResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*iso15118.Get15118EVCertificateRequest)
		suite.Require().True(ok)
		suite.Equal(schemaVersion, request.SchemaVersion)
		suite.Equal(action, request.Action)
		suite.Equal(exiRequest, request.ExiRequest)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	response, err := suite.chargingStation.Get15118EVCertificate(schemaVersion, action, exiRequest)
	suite.Require().Nil(err)
	suite.Require().NotNil(response)
	suite.Equal(status, response.Status)
	suite.Equal(exiResponse, response.ExiResponse)
	suite.Require().NotNil(response.StatusInfo)
	suite.Equal(statusInfo.ReasonCode, response.StatusInfo.ReasonCode)
	suite.Equal(statusInfo.AdditionalInfo, response.StatusInfo.AdditionalInfo)
}

func (suite *OcppV2TestSuite) TestGet15118EVCertificateInvalidEndpoint() {
	messageId := defaultMessageId
	schemaVersion := "1.0"
	action := iso15118.CertificateActionInstall
	exiRequest := "deadbeef"
	request := iso15118.NewGet15118EVCertificateRequest(schemaVersion, action, exiRequest)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"iso15118SchemaVersion":"%v","action":"%v","exiRequest":"%v"}]`,
		messageId, iso15118.Get15118EVCertificateFeatureName, schemaVersion, action, exiRequest)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
