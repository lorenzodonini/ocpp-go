package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Tests
func (suite *OcppV2TestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ocpp2.ModemType{Iccid: "test", Imsi: "test"}}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ocpp2.ModemType{Iccid: "test"}}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ocpp2.ModemType{Imsi: "test"}}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ocpp2.ModemType{}}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version"}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test"}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: "test"}}, true},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: "test"}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{VendorName: "test"}}, false},
		{ocpp2.BootNotificationRequest{ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: "test"}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: ">20..................", VendorName: "test"}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: ">50................................................"}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{SerialNumber: ">20..................", Model: "test", VendorName: "test"}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: "test", FirmwareVersion: ">50................................................"}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: "test", Modem: &ocpp2.ModemType{Iccid: ">20.................."}}}, false},
		{ocpp2.BootNotificationRequest{Reason: ocpp2.BootReasonPowerUp, ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: "test", Modem: &ocpp2.ModemType{Imsi: ">20.................."}}}, false},
		{ocpp2.BootNotificationRequest{Reason: "invalidReason", ChargingStation: ocpp2.ChargingStationType{Model: "test", VendorName: "test"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{ocpp2.BootNotificationConfirmation{CurrentTime: ocpp2.NewDateTime(time.Now()), Interval: 60, Status: ocpp2.RegistrationStatusAccepted}, true},
		{ocpp2.BootNotificationConfirmation{CurrentTime: ocpp2.NewDateTime(time.Now()), Status: ocpp2.RegistrationStatusAccepted}, true},
		{ocpp2.BootNotificationConfirmation{CurrentTime: ocpp2.NewDateTime(time.Now()), Interval: -1, Status: ocpp2.RegistrationStatusAccepted}, false},
		{ocpp2.BootNotificationConfirmation{CurrentTime: ocpp2.NewDateTime(time.Now()), Interval: 60, Status: "invalidRegistrationStatus"}, false},
		{ocpp2.BootNotificationConfirmation{CurrentTime: ocpp2.NewDateTime(time.Now()), Interval: 60}, false},
		{ocpp2.BootNotificationConfirmation{Interval: 60, Status: ocpp2.RegistrationStatusAccepted}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestBootNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	interval := 60
	reason := ocpp2.BootReasonPowerUp
	chargePointModel := "model1"
	chargePointVendor := "ABL"
	registrationStatus := ocpp2.RegistrationStatusAccepted
	currentTime := ocpp2.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reason":"%v","chargingStation":{"model":"%v","vendorName":"%v"}}]`, messageId, ocpp2.BootNotificationFeatureName, reason, chargePointModel, chargePointVendor)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v","interval":%v,"status":"%v"}]`, messageId, ocpp2.FormatTimestamp(currentTime.Time), interval, registrationStatus)
	bootNotificationConfirmation := ocpp2.NewBootNotificationConfirmation(currentTime, interval, registrationStatus)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnBootNotification", mock.AnythingOfType("string"), mock.Anything).Return(bootNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*ocpp2.BootNotificationRequest)
		assert.Equal(t, reason, request.Reason)
		assert.Equal(t, chargePointVendor, request.ChargingStation.VendorName)
		assert.Equal(t, chargePointModel, request.ChargingStation.Model)
	})
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargePoint.BootNotification(reason, chargePointModel, chargePointVendor)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, registrationStatus, confirmation.Status)
	assert.Equal(t, interval, confirmation.Interval)
	assertDateTimeEquality(t, *currentTime, *confirmation.CurrentTime)
}

func (suite *OcppV2TestSuite) TestBootNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	chargePointModel := "model1"
	chargePointVendor := "ABL"
	reason := ocpp2.BootReasonPowerUp
	bootNotificationRequest := ocpp2.NewBootNotificationRequest(reason, chargePointModel, chargePointVendor)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reason":"%v","chargingStation":{"model":"%v","vendorName":"%v"}}]`, messageId, ocpp2.BootNotificationFeatureName, reason, chargePointModel, chargePointVendor)
	testUnsupportedRequestFromCentralSystem(suite, bootNotificationRequest, requestJson, messageId)
}
