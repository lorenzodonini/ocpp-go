package availability

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type availabilityTestSuite struct {
	suite.Suite
}

func (suite *availabilityTestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ChangeAvailabilityRequest{OperationalStatus: OperationalStatusOperative, Evse: &types.EVSE{ID: 1, ConnectorID: tests.NewInt(1)}}, true},
		{ChangeAvailabilityRequest{OperationalStatus: OperationalStatusInoperative, Evse: &types.EVSE{ID: 1}}, true},
		{ChangeAvailabilityRequest{OperationalStatus: OperationalStatusInoperative}, true},
		{ChangeAvailabilityRequest{OperationalStatus: OperationalStatusOperative}, true},
		{ChangeAvailabilityRequest{}, false},
		{ChangeAvailabilityRequest{OperationalStatus: "invalidAvailabilityType"}, false},
		{ChangeAvailabilityRequest{OperationalStatus: OperationalStatusOperative, Evse: &types.EVSE{ID: -1}}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *availabilityTestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []tests.GenericTestEntry{
		{ChangeAvailabilityResponse{Status: ChangeAvailabilityStatusAccepted}, true},
		{ChangeAvailabilityResponse{Status: ChangeAvailabilityStatusRejected}, true},
		{ChangeAvailabilityResponse{Status: ChangeAvailabilityStatusScheduled}, true},
		{ChangeAvailabilityResponse{Status: "invalidAvailabilityStatus"}, false},
		{ChangeAvailabilityResponse{}, false},
	}
	tests.ExecuteGenericTestTable(t, testTable)
}

func (suite *availabilityTestSuite) TestChangeAvailabilityFeature() {
	feature := ChangeAvailabilityFeature{}
	suite.Equal(ChangeAvailabilityFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ChangeAvailabilityRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ChangeAvailabilityResponse{}), feature.GetResponseType())
}

func (suite *availabilityTestSuite) TestNewChangeAvailabilityRequest() {
	req := NewChangeAvailabilityRequest(OperationalStatusOperative)
	suite.NotNil(req)
	suite.Equal(ChangeAvailabilityFeatureName, req.GetFeatureName())
	suite.Equal(OperationalStatusOperative, req.OperationalStatus)
}

func (suite *availabilityTestSuite) TestNewChangeAvailabilityResponse() {
	resp := NewChangeAvailabilityResponse(ChangeAvailabilityStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(ChangeAvailabilityFeatureName, resp.GetFeatureName())
	suite.Equal(ChangeAvailabilityStatusAccepted, resp.Status)
}

func TestAvailabilitySuite(t *testing.T) {
	suite.Run(t, new(availabilityTestSuite))
}
