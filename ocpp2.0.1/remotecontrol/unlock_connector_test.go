package remotecontrol

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *remoteControlTestSuite) TestUnlockConnectorRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{UnlockConnectorRequest{EvseID: 2, ConnectorID: 1}, true},
		{UnlockConnectorRequest{EvseID: 2}, true},
		{UnlockConnectorRequest{}, true},
		{UnlockConnectorRequest{EvseID: -1, ConnectorID: 1}, false},
		{UnlockConnectorRequest{EvseID: 2, ConnectorID: -1}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *remoteControlTestSuite) TestUnlockConnectorResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{UnlockConnectorResponse{Status: UnlockStatusUnlocked, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{UnlockConnectorResponse{Status: UnlockStatusUnlocked}, true},
		{UnlockConnectorResponse{}, false},
		{UnlockConnectorResponse{Status: "invalidUnlockStatus", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{UnlockConnectorResponse{Status: UnlockStatusUnlocked, StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *remoteControlTestSuite) TestUnlockConnectorFeature() {
	feature := UnlockConnectorFeature{}
	suite.Equal(UnlockConnectorFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(UnlockConnectorRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(UnlockConnectorResponse{}), feature.GetResponseType())
}

func (suite *remoteControlTestSuite) TestNewUnlockConnectorRequest() {
	req := NewUnlockConnectorRequest(2, 1)
	suite.NotNil(req)
	suite.Equal(UnlockConnectorFeatureName, req.GetFeatureName())
	suite.Equal(2, req.EvseID)
	suite.Equal(1, req.ConnectorID)
}

func (suite *remoteControlTestSuite) TestNewUnlockConnectorResponse() {
	resp := NewUnlockConnectorResponse(UnlockStatusUnlocked)
	suite.NotNil(resp)
	suite.Equal(UnlockConnectorFeatureName, resp.GetFeatureName())
	suite.Equal(UnlockStatusUnlocked, resp.Status)
}
