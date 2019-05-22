package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	v16 "github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"time"
)

func (suite *CoreTestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{v16.AuthorizeRequest{IdTag: "12345"}, true},
		{v16.AuthorizeRequest{}, false},
		{v16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	executeRequestTestTable(t, requestTable)
}

func (suite *CoreTestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry {
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), ParentIdTag: "00000", Status: ocpp.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ParentIdTag: "00000", Status: ocpp.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), Status: ocpp.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusBlocked}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusExpired}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusInvalid}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusConcurrentTx}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ParentIdTag: ">20..................", Status: ocpp.AuthorizationStatusAccepted}}, false},
		{v16.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * -8), Status: ocpp.AuthorizationStatusAccepted}}, false},
	}
	executeConfirmationTestTable(t, confirmationTable)
}
