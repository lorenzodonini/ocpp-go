package reservation

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *reservationTestSuite) TestReserveNowRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, EvseID: tests.NewInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, true},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, EvseID: tests.NewInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{ReserveNowRequest{ExpiryDateTime: types.NewDateTime(time.Now()), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now())}, false},
		{ReserveNowRequest{ID: 42, IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, false},
		{ReserveNowRequest{}, false},
		{ReserveNowRequest{ID: -1, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, EvseID: tests.NewInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: "invalidConnectorType", EvseID: tests.NewInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, EvseID: tests.NewInt(-1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, EvseID: tests.NewInt(1), IdToken: types.IdToken{IdToken: "1234", Type: "invalidIdToken"}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{ReserveNowRequest{ID: 42, ExpiryDateTime: types.NewDateTime(time.Now()), ConnectorType: ConnectorTypeCCS1, EvseID: tests.NewInt(1), IdToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: "invalidIdToken"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *reservationTestSuite) TestReserveNowConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ReserveNowResponse{Status: ReserveNowStatusAccepted, StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{ReserveNowResponse{Status: ReserveNowStatusAccepted}, true},
		{ReserveNowResponse{}, false},
		{ReserveNowResponse{Status: "invalidReserveNowStatus"}, false},
		{ReserveNowResponse{Status: ReserveNowStatusAccepted, StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *reservationTestSuite) TestReserveNowFeature() {
	feature := ReserveNowFeature{}
	suite.Equal(ReserveNowFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ReserveNowRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ReserveNowResponse{}), feature.GetResponseType())
}

func (suite *reservationTestSuite) TestNewReserveNowRequest() {
	expiry := types.NewDateTime(time.Now())
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	req := NewReserveNowRequest(7, expiry, idToken)
	suite.NotNil(req)
	suite.Equal(ReserveNowFeatureName, req.GetFeatureName())
	suite.Equal(7, req.ID)
	suite.Equal(expiry, req.ExpiryDateTime)
	suite.Equal(idToken, req.IdToken)
}

func (suite *reservationTestSuite) TestNewReserveNowResponse() {
	resp := NewReserveNowResponse(ReserveNowStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(ReserveNowFeatureName, resp.GetFeatureName())
	suite.Equal(ReserveNowStatusAccepted, resp.Status)
}
