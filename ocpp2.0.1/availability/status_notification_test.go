package availability

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *availabilityTestSuite) TestStatusNotificationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: ConnectorStatusAvailable, EvseID: 1, ConnectorID: 1}, true},
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: ConnectorStatusAvailable, EvseID: 1}, true},
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: ConnectorStatusAvailable}, true},
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now())}, false},
		{StatusNotificationRequest{ConnectorStatus: ConnectorStatusAvailable}, false},
		{StatusNotificationRequest{}, false},
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: "invalidConnectorStatus", EvseID: 1, ConnectorID: 1}, false},
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: ConnectorStatusAvailable, EvseID: -1, ConnectorID: 1}, false},
		{StatusNotificationRequest{Timestamp: types.NewDateTime(time.Now()), ConnectorStatus: ConnectorStatusAvailable, EvseID: 1, ConnectorID: -1}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *availabilityTestSuite) TestStatusNotificationResponseValidation() {
	t := suite.T()
	var responseTable = []tests.GenericTestEntry{
		{StatusNotificationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, responseTable)
}

func (suite *availabilityTestSuite) TestStatusNotificationFeature() {
	feature := StatusNotificationFeature{}
	suite.Equal(StatusNotificationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(StatusNotificationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(StatusNotificationResponse{}), feature.GetResponseType())
}

func (suite *availabilityTestSuite) TestNewStatusNotificationRequest() {
	ts := types.NewDateTime(time.Now())
	req := NewStatusNotificationRequest(ts, ConnectorStatusAvailable, 1, 2)
	suite.NotNil(req)
	suite.Equal(StatusNotificationFeatureName, req.GetFeatureName())
	suite.Equal(ts, req.Timestamp)
	suite.Equal(ConnectorStatusAvailable, req.ConnectorStatus)
	suite.Equal(1, req.EvseID)
	suite.Equal(2, req.ConnectorID)
}

func (suite *availabilityTestSuite) TestNewStatusNotificationResponse() {
	resp := NewStatusNotificationResponse()
	suite.NotNil(resp)
	suite.Equal(StatusNotificationFeatureName, resp.GetFeatureName())
}
