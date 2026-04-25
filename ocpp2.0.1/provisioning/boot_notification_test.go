package provisioning

import (
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type provisioningTestSuite struct {
	suite.Suite
}

func (suite *provisioningTestSuite) TestBootNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ModemType{Iccid: "test", Imsi: "test"}}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ModemType{Iccid: "test"}}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ModemType{Imsi: "test"}}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version", Modem: &ModemType{}}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test", FirmwareVersion: "version"}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: "number", Model: "test", VendorName: "test"}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: "test", VendorName: "test"}}, true},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: "test"}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{VendorName: "test"}}, false},
		{BootNotificationRequest{ChargingStation: ChargingStationType{Model: "test", VendorName: "test"}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: ">20..................", VendorName: "test"}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: "test", VendorName: ">50................................................"}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{SerialNumber: ">25.......................", Model: "test", VendorName: "test"}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: "test", VendorName: "test", FirmwareVersion: ">50................................................"}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: "test", VendorName: "test", Modem: &ModemType{Iccid: ">20.................."}}}, false},
		{BootNotificationRequest{Reason: BootReasonPowerUp, ChargingStation: ChargingStationType{Model: "test", VendorName: "test", Modem: &ModemType{Imsi: ">20.................."}}}, false},
		{BootNotificationRequest{Reason: "invalidReason", ChargingStation: ChargingStationType{Model: "test", VendorName: "test"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestBootNotificationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: RegistrationStatusAccepted}, true},
		{BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Status: RegistrationStatusAccepted}, true},
		{BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: -1, Status: RegistrationStatusAccepted}, false},
		{BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: 60, Status: "invalidRegistrationStatus"}, false},
		{BootNotificationResponse{CurrentTime: types.NewDateTime(time.Now()), Interval: 60}, false},
		{BootNotificationResponse{Interval: 60, Status: RegistrationStatusAccepted}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestBootNotificationFeature() {
	feature := BootNotificationFeature{}
	suite.Equal(BootNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(BootNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(BootNotificationResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewBootNotificationRequest() {
	req := NewBootNotificationRequest(BootReasonPowerUp, "ModelX", "VendorY")
	suite.NotNil(req)
	suite.Equal(BootNotificationFeatureName, req.GetFeatureName())
	suite.Equal(BootReasonPowerUp, req.Reason)
	suite.Equal("ModelX", req.ChargingStation.Model)
	suite.Equal("VendorY", req.ChargingStation.VendorName)
}

func (suite *provisioningTestSuite) TestNewBootNotificationResponse() {
	now := types.NewDateTime(time.Now())
	resp := NewBootNotificationResponse(now, 60, RegistrationStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(BootNotificationFeatureName, resp.GetFeatureName())
	suite.Equal(now, resp.CurrentTime)
	suite.Equal(60, resp.Interval)
	suite.Equal(RegistrationStatusAccepted, resp.Status)
}

func TestProvisioningSuite(t *testing.T) {
	suite.Run(t, new(provisioningTestSuite))
}
