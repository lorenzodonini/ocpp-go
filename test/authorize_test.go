package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6"
	"time"
)

func (suite *OcppV16TestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{v16.AuthorizeRequest{IdTag: "12345"}, true},
		{v16.AuthorizeRequest{}, false},
		{v16.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	executeRequestTestTable(t, requestTable)
}

func (suite *OcppV16TestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry {
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), ParentIdTag: "00000", Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ParentIdTag: "00000", Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusAccepted}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusBlocked}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusExpired}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusInvalid}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{Status: v16.AuthorizationStatusConcurrentTx}}, true},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ParentIdTag: ">20..................", Status: v16.AuthorizationStatusAccepted}}, false},
		{v16.AuthorizeConfirmation{IdTagInfo: v16.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * -8), Status: v16.AuthorizationStatusAccepted}}, false},
	}
	executeConfirmationTestTable(t, confirmationTable)
}
