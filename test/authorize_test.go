package test

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6/core"
	"time"
)

func (suite *OcppTestSuite) TestAuthorizeRequestValidation() {
	t := suite.T()
	var requestTable = []RequestTestEntry{
		{core.AuthorizeRequest{IdTag: "12345"}, true},
		{core.AuthorizeRequest{}, false},
		{core.AuthorizeRequest{IdTag: ">20.................."}, false},
	}
	executeRequestTestTable(t, requestTable)
}

func (suite *OcppTestSuite) TestAuthorizeConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []ConfirmationTestEntry {
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), ParentIdTag: "00000", Status: ocpp.AuthorizationStatusAccepted}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ParentIdTag: "00000", Status: ocpp.AuthorizationStatusAccepted}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * 8), Status: ocpp.AuthorizationStatusAccepted}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusAccepted}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusBlocked}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusExpired}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusInvalid}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{Status: ocpp.AuthorizationStatusConcurrentTx}}, true},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ParentIdTag: ">20..................", Status: ocpp.AuthorizationStatusAccepted}}, false},
		{core.AuthorizeConfirmation{IdTagInfo: ocpp.IdTagInfo{ExpiryDate: time.Now().Add(time.Hour * -8), Status: ocpp.AuthorizationStatusAccepted}}, false},
	}
	executeConfirmationTestTable(t, confirmationTable)
}
