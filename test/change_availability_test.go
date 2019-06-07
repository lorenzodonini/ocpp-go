package test

import v16 "github.com/lorenzodonini/go-ocpp/ocpp/1.6"

func (suite *OcppV16TestSuite) TestChangeAvailabilityRequestValidation() {
	t := suite.T()
	var testTable = []RequestTestEntry{
		{v16.ChangeAvailabilityRequest{ConnectorId: 0, Type: v16.AvailabilityTypeOperative}, true},
		{v16.ChangeAvailabilityRequest{ConnectorId: 0, Type: v16.AvailabilityTypeInoperative}, true},
		{v16.ChangeAvailabilityRequest{ConnectorId: 0}, false},
		{v16.ChangeAvailabilityRequest{Type: v16.AvailabilityTypeOperative}, true},
		{v16.ChangeAvailabilityRequest{ConnectorId: -1, Type: v16.AvailabilityTypeOperative}, false},
	}
	executeRequestTestTable(t, testTable)
}

func (suite *OcppV16TestSuite) TestChangeAvailabilityConfirmationValidation() {
	t := suite.T()
	var testTable = []ConfirmationTestEntry{
		{v16.ChangeAvailabilityConfirmation{Status: v16.AvailabilityStatusAccepted}, true},
		{v16.ChangeAvailabilityConfirmation{Status: v16.AvailabilityStatusRejected}, true},
		{v16.ChangeAvailabilityConfirmation{Status: v16.AvailabilityStatusScheduled}, true},
		{v16.ChangeAvailabilityConfirmation{}, false},
	}
	executeConfirmationTestTable(t, testTable)
}
