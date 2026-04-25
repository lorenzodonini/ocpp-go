package reservation

import (
	"reflect"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type reservationTestSuite struct {
	suite.Suite
}

func (suite *reservationTestSuite) TestCancelReservationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{CancelReservationRequest{ReservationID: 42}, true},
		{CancelReservationRequest{}, true},
		{CancelReservationRequest{ReservationID: -1}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *reservationTestSuite) TestCancelReservationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{CancelReservationResponse{Status: CancelReservationStatusAccepted, StatusInfo: types.NewStatusInfo("200", "ok")}, true},
		{CancelReservationResponse{Status: CancelReservationStatusAccepted}, true},
		{CancelReservationResponse{Status: "invalidCancelReservationStatus"}, false},
		{CancelReservationResponse{}, false},
		{CancelReservationResponse{Status: CancelReservationStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *reservationTestSuite) TestCancelReservationFeature() {
	feature := CancelReservationFeature{}
	suite.Equal(CancelReservationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(CancelReservationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(CancelReservationResponse{}), feature.GetResponseType())
}

func (suite *reservationTestSuite) TestNewCancelReservationRequest() {
	req := NewCancelReservationRequest(42)
	suite.NotNil(req)
	suite.Equal(CancelReservationFeatureName, req.GetFeatureName())
	suite.Equal(42, req.ReservationID)
}

func (suite *reservationTestSuite) TestNewCancelReservationResponse() {
	resp := NewCancelReservationResponse(CancelReservationStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(CancelReservationFeatureName, resp.GetFeatureName())
	suite.Equal(CancelReservationStatusAccepted, resp.Status)
}

func TestReservationTestSuite(t *testing.T) {
	suite.Run(t, new(reservationTestSuite))
}
