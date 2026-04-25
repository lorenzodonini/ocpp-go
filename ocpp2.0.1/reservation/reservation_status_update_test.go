package reservation

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *reservationTestSuite) TestReservationStatusUpdateRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{ReservationStatusUpdateRequest{ReservationID: 42, Status: ReservationUpdateStatusExpired}, true},
		{ReservationStatusUpdateRequest{ReservationID: 42, Status: ReservationUpdateStatusRemoved}, true},
		{ReservationStatusUpdateRequest{Status: ReservationUpdateStatusExpired}, true},
		{ReservationStatusUpdateRequest{}, false},
		{ReservationStatusUpdateRequest{ReservationID: 42}, false},
		{ReservationStatusUpdateRequest{ReservationID: -1, Status: ReservationUpdateStatusExpired}, false},
		{ReservationStatusUpdateRequest{ReservationID: 42, Status: "invalidReservationStatus"}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *reservationTestSuite) TestReservationStatusUpdateConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{ReservationStatusUpdateResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *reservationTestSuite) TestReservationStatusUpdateFeature() {
	feature := ReservationStatusUpdateFeature{}
	suite.Equal(ReservationStatusUpdateFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(ReservationStatusUpdateRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(ReservationStatusUpdateResponse{}), feature.GetResponseType())
}

func (suite *reservationTestSuite) TestNewReservationStatusUpdateRequest() {
	req := NewReservationStatusUpdateRequest(10, ReservationUpdateStatusExpired)
	suite.NotNil(req)
	suite.Equal(ReservationStatusUpdateFeatureName, req.GetFeatureName())
	suite.Equal(10, req.ReservationID)
	suite.Equal(ReservationUpdateStatusExpired, req.Status)
}

func (suite *reservationTestSuite) TestNewReservationStatusUpdateResponse() {
	resp := NewReservationStatusUpdateResponse()
	suite.NotNil(resp)
	suite.Equal(ReservationStatusUpdateFeatureName, resp.GetFeatureName())
}
