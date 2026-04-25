package firmware

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *firmwareTestSuite) TestUpdateFirmwareRequestValidation() {
	t := suite.T()
	fw := Firmware{
		Location:           "https://someurl",
		RetrieveDateTime:   types.NewDateTime(time.Now()),
		InstallDateTime:    types.NewDateTime(time.Now()),
		SigningCertificate: "1337c0de",
		Signature:          "deadc0de",
	}
	var requestTable = []tests.GenericTestEntry{
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: fw}, true},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RequestID: 42, Firmware: fw}, true},
		{UpdateFirmwareRequest{RequestID: 42, Firmware: fw}, true},
		{UpdateFirmwareRequest{Firmware: fw}, true},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{Location: "https://someurl", RetrieveDateTime: types.NewDateTime(time.Now())}}, true},
		{UpdateFirmwareRequest{}, false},
		{UpdateFirmwareRequest{Retries: tests.NewInt(-1), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: fw}, false},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(-1), RequestID: 42, Firmware: fw}, false},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: -1, Firmware: fw}, false},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{Location: "https://someurl", InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
		{UpdateFirmwareRequest{Retries: tests.NewInt(5), RetryInterval: tests.NewInt(300), RequestID: 42, Firmware: Firmware{Location: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", RetrieveDateTime: types.NewDateTime(time.Now()), InstallDateTime: types.NewDateTime(time.Now()), SigningCertificate: "1337c0de", Signature: "deadc0de"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *firmwareTestSuite) TestUpdateFirmwareResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{UpdateFirmwareResponse{Status: UpdateFirmwareStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "ok", AdditionalInfo: "someInfo"}}, true},
		{UpdateFirmwareResponse{Status: UpdateFirmwareStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "ok"}}, true},
		{UpdateFirmwareResponse{Status: UpdateFirmwareStatusAccepted}, true},
		{UpdateFirmwareResponse{}, false},
		{UpdateFirmwareResponse{Status: "invalidFirmwareUpdateStatus"}, false},
		{UpdateFirmwareResponse{Status: UpdateFirmwareStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *firmwareTestSuite) TestUpdateFirmwareFeature() {
	feature := UpdateFirmwareFeature{}
	suite.Equal(UpdateFirmwareFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(UpdateFirmwareRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(UpdateFirmwareResponse{}), feature.GetResponseType())
}

func (suite *firmwareTestSuite) TestNewUpdateFirmwareRequest() {
	fw := Firmware{
		Location:         "https://example.com/fw",
		RetrieveDateTime: types.NewDateTime(time.Now()),
	}
	req := NewUpdateFirmwareRequest(1, fw)
	suite.NotNil(req)
	suite.Equal(UpdateFirmwareFeatureName, req.GetFeatureName())
	suite.Equal(1, req.RequestID)
	suite.Equal(fw, req.Firmware)
}

func (suite *firmwareTestSuite) TestNewUpdateFirmwareResponse() {
	resp := NewUpdateFirmwareResponse(UpdateFirmwareStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(UpdateFirmwareFeatureName, resp.GetFeatureName())
	suite.Equal(UpdateFirmwareStatusAccepted, resp.Status)
}
