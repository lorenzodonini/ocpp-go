package ocpp16_test

import (
	"fmt"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Test
func (suite *OcppV16TestSuite) TestReserveNowRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: "9999"}, true},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{ocpp16.ReserveNowRequest{ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, true},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345"}, true},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now())}, false},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, IdTag: "12345"}, false},
		{ocpp16.ReserveNowRequest{ConnectorId: -1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42}, false},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345", ReservationId: -1}, false},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: ">20.................."}, false},
		{ocpp16.ReserveNowRequest{ConnectorId: 1, ExpiryDate: ocpp16.NewDateTime(time.Now()), IdTag: "12345", ReservationId: 42, ParentIdTag: ">20.................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestReserveNowConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp16.ReserveNowConfirmation{Status: ocpp16.ReservationStatusAccepted}, true},
		{ocpp16.ReserveNowConfirmation{Status: "invalidReserveNowStatus"}, false},
		{ocpp16.ReserveNowConfirmation{}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestReserveNowE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	idTag := "12345"
	parentIdTag := "00000"
	connectorId := 1
	reservationId := 42
	expiryDate := ocpp16.NewDateTime(time.Now())
	status := ocpp16.ReservationStatusAccepted
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"expiryDate":"%v","idTag":"%v","parentIdTag":"%v","reservationId":%v}]`,
		messageId, ocpp16.ReserveNowFeatureName, connectorId, expiryDate.Format(ocpp16.ISO8601), idTag, parentIdTag, reservationId)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v"}]`, messageId, status)
	ReserveNowConfirmation := ocpp16.NewReserveNowConfirmation(status)
	channel := NewMockWebSocket(wsId)

	reservationListener := MockChargePointReservationListener{}
	reservationListener.On("OnReserveNow", mock.Anything).Return(ReserveNowConfirmation, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*ocpp16.ReserveNowRequest)
		assert.True(t, ok)
		assert.Equal(t, connectorId, request.ConnectorId)
		assert.Equal(t, expiryDate.Format(ocpp16.ISO8601), request.ExpiryDate.Format(ocpp16.ISO8601))
		assert.Equal(t, idTag, request.IdTag)
		assert.Equal(t, parentIdTag, request.ParentIdTag)
		assert.Equal(t, reservationId, request.ReservationId)
	})
	setupDefaultCentralSystemHandlers(suite, nil, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	suite.chargePoint.SetReservationListener(reservationListener)
	// Run Test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	resultChannel := make(chan bool, 1)
	err = suite.centralSystem.ReserveNow(wsId, func(confirmation *ocpp16.ReserveNowConfirmation, err error) {
		assert.Nil(t, err)
		assert.NotNil(t, confirmation)
		assert.Equal(t, status, confirmation.Status)
		resultChannel <- true
	}, connectorId, expiryDate, idTag, reservationId, func(request *ocpp16.ReserveNowRequest) {
		request.ParentIdTag = parentIdTag
	})
	assert.Nil(t, err)
	result := <-resultChannel
	assert.True(t, result)
}

func (suite *OcppV16TestSuite) TestReserveNowInvalidEndpoint() {
	messageId := defaultMessageId
	idTag := "12345"
	parentIdTag := "00000"
	connectorId := 1
	reservationId := 42
	expiryDate := ocpp16.NewDateTime(time.Now())
	reserveNowRequest := ocpp16.NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
	reserveNowRequest.ParentIdTag = parentIdTag
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"connectorId":%v,"expiryDate":"%v","idTag":"%v","parentIdTag":"%v","reservationId":%v}]`,
		messageId, ocpp16.ReserveNowFeatureName, connectorId, expiryDate.Format(ocpp16.ISO8601), idTag, parentIdTag, reservationId)
	testUnsupportedRequestFromChargePoint(suite, reserveNowRequest, requestJson, messageId)
}
