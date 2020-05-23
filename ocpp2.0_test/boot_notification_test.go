package ocpp2_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"time"
)

// Tests
func (suite *OcppV2TestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &provisioning.ModemType{Iccid: "test", Imsi: "test"}}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &provisioning.ModemType{Iccid: "test"}}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &provisioning.ModemType{Imsi: "test"}}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &provisioning.ModemType{}}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version"}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test"}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: "test"}}, true},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: "test"}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{VendorName: "test"}}, false},
		{provisioning.BootNotificationRequest{ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: "test"}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: ">20..................", VendorName: "test"}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: ">50................................................"}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{SerialNumber: ">20..................", Model: "test", VendorName: "test"}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: "test", FirmwareVersion: ">50................................................"}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: "test", Modem: &provisioning.ModemType{Iccid: ">20.................."}}}, false},
		{provisioning.BootNotificationRequest{Reason: provisioning.BootReasonPowerUp, ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: "test", Modem: &provisioning.ModemType{Imsi: ">20.................."}}}, false},
		{provisioning.BootNotificationRequest{Reason: "invalidReason", ChargingStation: provisioning.ChargingStationType{Model: "test", VendorName: "test"}}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV2TestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{provisioning.BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: provisioning.RegistrationStatusAccepted}, true},
		{provisioning.BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Status: provisioning.RegistrationStatusAccepted}, true},
		{provisioning.BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: -1, Status: provisioning.RegistrationStatusAccepted}, false},
		{provisioning.BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: "invalidRegistrationStatus"}, false},
		{provisioning.BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: 60}, false},
		{provisioning.BootNotificationResponse{Interval: 60, Status: provisioning.RegistrationStatusAccepted}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV2TestSuite) TestBootNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	interval := 60
	reason := provisioning.BootReasonPowerUp
	chargePointModel := "model1"
	chargePointVendor := "ABL"
	registrationStatus := provisioning.RegistrationStatusAccepted
	currentTime := types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reason":"%v","chargingStation":{"model":"%v","vendorName":"%v"}}]`, messageId, provisioning.BootNotificationFeatureName, reason, chargePointModel, chargePointVendor)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v","interval":%v,"status":"%v"}]`, messageId, currentTime.FormatTimestamp(), interval, registrationStatus)
	fmt.Println(responseJson)
	bootNotificationConfirmation := provisioning.NewBootNotificationResponse(currentTime, interval, registrationStatus)
	channel := NewMockWebSocket(wsId)

	handler := MockCSMSProvisioningHandler{}
	handler.On("OnBootNotification", mock.AnythingOfType("string"), mock.Anything).Return(bootNotificationConfirmation, nil).Run(func(args mock.Arguments) {
		request := args.Get(1).(*provisioning.BootNotificationRequest)
		assert.Equal(t, reason, request.Reason)
		assert.Equal(t, chargePointVendor, request.ChargingStation.VendorName)
		assert.Equal(t, chargePointModel, request.ChargingStation.Model)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	require.Nil(t, err)
	confirmation, err := suite.chargingStation.BootNotification(reason, chargePointModel, chargePointVendor)
	require.Nil(t, err)
	require.NotNil(t, confirmation)
	assert.Equal(t, registrationStatus, confirmation.Status)
	assert.Equal(t, interval, confirmation.Interval)
	assertDateTimeEquality(t, currentTime, confirmation.CurrentTime)
}

func (suite *OcppV2TestSuite) TestBootNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	chargePointModel := "model1"
	chargePointVendor := "ABL"
	reason := provisioning.BootReasonPowerUp
	bootNotificationRequest := provisioning.NewBootNotificationRequest(reason, chargePointModel, chargePointVendor)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"reason":"%v","chargingStation":{"model":"%v","vendorName":"%v"}}]`, messageId, provisioning.BootNotificationFeatureName, reason, chargePointModel, chargePointVendor)
	testUnsupportedRequestFromCentralSystem(suite, bootNotificationRequest, requestJson, messageId)
}
