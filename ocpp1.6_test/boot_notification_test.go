package ocpp16_test

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
)

// Tests
func (suite *OcppV16TestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []GenericTestEntry{
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test"}, true},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, true},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointSerialNumber: "number", ChargePointVendor: "test", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: "test", ChargePointModel: "test", ChargePointSerialNumber: "number", FirmwareVersion: "version", Iccid: "test", Imsi: "test"}, false},
		{core.BootNotificationRequest{ChargeBoxSerialNumber: ">25.......................", ChargePointModel: "test", ChargePointVendor: "test"}, false},
		{core.BootNotificationRequest{ChargePointModel: ">20..................", ChargePointVendor: "test"}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointSerialNumber: ">25.......................", ChargePointVendor: "test"}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: ">20.................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", FirmwareVersion: ">50................................................"}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Iccid: ">20.................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", Imsi: ">20.................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterSerialNumber: ">25......................."}, false},
		{core.BootNotificationRequest{ChargePointModel: "test", ChargePointVendor: "test", MeterType: ">25......................."}, false},
	}
	ExecuteGenericTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []GenericTestEntry{
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: core.RegistrationStatusAccepted}, true},
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: core.RegistrationStatusPending}, true},
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: core.RegistrationStatusRejected}, true},
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Status: core.RegistrationStatusAccepted}, true},
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: -1, Status: core.RegistrationStatusRejected}, false},
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: "invalidRegistrationStatus"}, false},
		{core.BootNotificationConfirmation{CurrentTime: types.NewDateTime(time.Now()), Interval: 60}, false},
		{core.BootNotificationConfirmation{Interval: 60, Status: core.RegistrationStatusAccepted}, false},
	}
	ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *OcppV16TestSuite) TestBootNotificationE2EMocked() {
	t := suite.T()
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	interval := 60
	chargePointModel := "model1"
	chargePointVendor := "ABL"
	registrationStatus := core.RegistrationStatusAccepted
	currentTime := types.NewDateTime(time.Now())
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargePointModel":"%v","chargePointVendor":"%v"}]`, messageId, core.BootNotificationFeatureName, chargePointModel, chargePointVendor)
	responseJson := fmt.Sprintf(`[3,"%v",{"currentTime":"%v","interval":%v,"status":"%v"}]`, messageId, currentTime.FormatTimestamp(), interval, registrationStatus)
	bootNotificationConfirmation := core.NewBootNotificationConfirmation(currentTime, interval, registrationStatus)
	channel := NewMockWebSocket(wsId)

	coreListener := MockCentralSystemCoreListener{}
	coreListener.On("OnBootNotification", mock.AnythingOfType("string"), mock.Anything).Return(bootNotificationConfirmation, nil)
	setupDefaultCentralSystemHandlers(suite, coreListener, expectedCentralSystemOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true})
	setupDefaultChargePointHandlers(suite, nil, expectedChargePointOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.centralSystem.Start(8887, "somePath")
	err := suite.chargePoint.Start(wsUrl)
	assert.Nil(t, err)
	confirmation, err := suite.chargePoint.BootNotification(chargePointModel, chargePointVendor)
	assert.Nil(t, err)
	assert.NotNil(t, confirmation)
	assert.Equal(t, registrationStatus, confirmation.Status)
	assert.Equal(t, interval, confirmation.Interval)
	assertDateTimeEquality(t, *currentTime, *confirmation.CurrentTime)
}

func (suite *OcppV16TestSuite) TestBootNotificationInvalidEndpoint() {
	messageId := defaultMessageId
	chargePointModel := "model1"
	chargePointVendor := "ABL"
	bootNotificationRequest := core.NewBootNotificationRequest(chargePointModel, chargePointVendor)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"chargePointModel":"%v","chargePointVendor":"%v"}]`, messageId, core.BootNotificationFeatureName, chargePointModel, chargePointVendor)
	testUnsupportedRequestFromCentralSystem(suite, bootNotificationRequest, requestJson, messageId)
}
