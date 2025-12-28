package ocpp16_test

import (
	"fmt"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6_test/mocks"
	"github.com/stretchr/testify/mock"
)

func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsRequestValidation() {
	var testTable = []GenericTestEntry{
		{certificates.GetInstalledCertificateIdsRequest{CertificateType: types.CentralSystemRootCertificate}, true},
		{certificates.GetInstalledCertificateIdsRequest{}, false},
		{certificates.GetInstalledCertificateIdsRequest{CertificateType: "invalidCertificateUse"}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsConfirmationValidation() {
	var testTable = []GenericTestEntry{
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusAccepted}, true},
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusNotFound}, true},
		{certificates.GetInstalledCertificateIdsResponse{Status: certificates.GetInstalledCertificateStatusAccepted}, true},
		{certificates.GetInstalledCertificateIdsResponse{}, false},
		{certificates.GetInstalledCertificateIdsResponse{Status: "invalidGetInstalledCertificateStatus"}, false},
	}
	ExecuteGenericTestTable(suite, testTable)
}

// Test
func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	certificateType := types.CentralSystemRootCertificate
	status := certificates.GetInstalledCertificateStatusAccepted

	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v"}]`, messageId, certificates.GetInstalledCertificateIdsFeatureName, certificateType)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`,
		messageId, status)
	getInstalledCertificateIdsConfirmation := certificates.NewGetInstalledCertificateIdsResponse(status)
	channel := NewMockWebSocket(wsId)

	// Setting handlers
	handler := mocks.NewMockCertificatesChargePointHandler(t)
	handler.EXPECT().OnGetInstalledCertificateIds(mock.Anything).RunAndReturn(func(request *certificates.GetInstalledCertificateIdsRequest) (*certificates.GetInstalledCertificateIdsResponse, error) {
		suite.Equal(certificateType, request.CertificateType)
		return getInstalledCertificateIdsConfirmation, nil
	})

	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetCertificateHandler(handler)

	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	suite.chargePoint.SetCertificateHandler(handler)
	err := suite.chargePoint.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.GetInstalledCertificateIds(wsId, func(confirmation *certificates.GetInstalledCertificateIdsResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(confirmation)
		suite.Equal(status, confirmation.Status)
		resultChannel <- true
	}, certificateType)

	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV16TestSuite) TestGetInstalledCertificateIdsInvalidEndpoint() {
	messageId := defaultMessageId
	certificateType := types.CentralSystemRootCertificate
	GetInstalledCertificateIdsRequest := certificates.NewGetInstalledCertificateIdsRequest(certificateType)
	GetInstalledCertificateIdsRequest.CertificateType = certificateType
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"certificateType":"%v"}]`, messageId, certificates.GetInstalledCertificateIdsFeatureName, certificateType)
	testUnsupportedRequestFromChargePoint(suite, GetInstalledCertificateIdsRequest, requestJson, messageId)
}
