package remotecontrol

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type remoteControlTestSuite struct {
	suite.Suite
}

func (suite *remoteControlTestSuite) TestRequestStartTransactionRequestValidation() {
	t := suite.T()
	chargingProfile := types.ChargingProfile{
		ID:                     1,
		StackLevel:             0,
		ChargingProfilePurpose: types.ChargingProfilePurposeTxProfile,
		ChargingProfileKind:    types.ChargingProfileKindAbsolute,
		ChargingSchedule: []types.ChargingSchedule{
			{
				ChargingRateUnit: types.ChargingRateUnitAmperes,
				ChargingSchedulePeriod: []types.ChargingSchedulePeriod{
					{
						StartPeriod: 0,
						Limit:       16.0,
					},
				},
			},
		},
	}
	var requestTable = []tests.GenericTestEntry{
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, ChargingProfile: &chargingProfile, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, true},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, ChargingProfile: &chargingProfile}, true},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{RequestStartTransactionRequest{RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{RequestStartTransactionRequest{IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}}, true},
		{RequestStartTransactionRequest{}, false},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(0), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, ChargingProfile: &chargingProfile, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: -1, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, ChargingProfile: &chargingProfile, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: "invalidIdToken"}, ChargingProfile: &chargingProfile, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, ChargingProfile: &types.ChargingProfile{}, GroupIdToken: &types.IdToken{IdToken: "1234", Type: types.IdTokenTypeISO15693}}, false},
		{RequestStartTransactionRequest{EvseID: tests.NewInt(1), RemoteStartID: 42, IDToken: types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}, ChargingProfile: &chargingProfile, GroupIdToken: &types.IdToken{IdToken: "1234", Type: "invalidGroupIdToken"}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *remoteControlTestSuite) TestRequestStartTransactionConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{RequestStartTransactionResponse{Status: RequestStartStopStatusAccepted, TransactionID: "12345", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, true},
		{RequestStartTransactionResponse{Status: RequestStartStopStatusAccepted, TransactionID: "12345"}, true},
		{RequestStartTransactionResponse{Status: RequestStartStopStatusAccepted}, true},
		{RequestStartTransactionResponse{Status: RequestStartStopStatusRejected}, true},
		{RequestStartTransactionResponse{}, false},
		{RequestStartTransactionResponse{Status: "invalidRequestStartStopStatus", TransactionID: "12345", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{RequestStartTransactionResponse{Status: RequestStartStopStatusAccepted, TransactionID: ">36..................................", StatusInfo: &types.StatusInfo{ReasonCode: "200"}}, false},
		{RequestStartTransactionResponse{Status: RequestStartStopStatusAccepted, TransactionID: "12345", StatusInfo: &types.StatusInfo{}}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *remoteControlTestSuite) TestRequestStartTransactionFeature() {
	feature := RequestStartTransactionFeature{}
	suite.Equal(RequestStartTransactionFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(RequestStartTransactionRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(RequestStartTransactionResponse{}), feature.GetResponseType())
}

func (suite *remoteControlTestSuite) TestNewRequestStartTransactionRequest() {
	idToken := types.IdToken{IdToken: "1234", Type: types.IdTokenTypeKeyCode}
	req := NewRequestStartTransactionRequest(42, idToken)
	suite.NotNil(req)
	suite.Equal(RequestStartTransactionFeatureName, req.GetFeatureName())
	suite.Equal(42, req.RemoteStartID)
	suite.Equal(idToken, req.IDToken)
}

func (suite *remoteControlTestSuite) TestNewRequestStartTransactionResponse() {
	resp := NewRequestStartTransactionResponse(RequestStartStopStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(RequestStartTransactionFeatureName, resp.GetFeatureName())
	suite.Equal(RequestStartStopStatusAccepted, resp.Status)
}

func TestRemoteControlSuite(t *testing.T) {
	suite.Run(t, new(remoteControlTestSuite))
}