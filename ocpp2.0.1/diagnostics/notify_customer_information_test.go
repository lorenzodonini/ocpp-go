package diagnostics

import (
	"reflect"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *diagnosticsTestSuite) TestNotifyCustomerInformationRequestValidation() {
	t := suite.T()
	var requestTable = []tests.GenericTestEntry{
		{NotifyCustomerInformationRequest{Data: "dummyData", Tbc: false, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{NotifyCustomerInformationRequest{Data: "dummyData", Tbc: true, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{NotifyCustomerInformationRequest{Data: "dummyData", SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{NotifyCustomerInformationRequest{Data: "dummyData", GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, true},
		{NotifyCustomerInformationRequest{Data: "dummyData", GeneratedAt: types.DateTime{Time: time.Now()}}, true},
		{NotifyCustomerInformationRequest{Data: "dummyData"}, true},
		{NotifyCustomerInformationRequest{}, false},
		{NotifyCustomerInformationRequest{Data: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", Tbc: false, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, false},
		{NotifyCustomerInformationRequest{Data: "dummyData", Tbc: false, SeqNo: -1, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: 42}, false},
		{NotifyCustomerInformationRequest{Data: "dummyData", Tbc: false, SeqNo: 0, GeneratedAt: types.DateTime{Time: time.Now()}, RequestID: -1}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *diagnosticsTestSuite) TestNotifyCustomerInformationConfirmationValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{NotifyCustomerInformationResponse{}, true},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *diagnosticsTestSuite) TestNotifyCustomerInformationFeature() {
	feature := NotifyCustomerInformationFeature{}
	suite.Equal(NotifyCustomerInformationFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(NotifyCustomerInformationRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(NotifyCustomerInformationResponse{}), feature.GetResponseType())
}

func (suite *diagnosticsTestSuite) TestNewNotifyCustomerInformationRequest() {
	ts := types.DateTime{Time: time.Now()}
	req := NewNotifyCustomerInformationRequest("data", 0, ts, 42)
	suite.NotNil(req)
	suite.Equal(NotifyCustomerInformationFeatureName, req.GetFeatureName())
	suite.Equal("data", req.Data)
	suite.Equal(0, req.SeqNo)
	suite.Equal(42, req.RequestID)
}

func (suite *diagnosticsTestSuite) TestNewNotifyCustomerInformationResponse() {
	resp := NewNotifyCustomerInformationResponse()
	suite.NotNil(resp)
	suite.Equal(NotifyCustomerInformationFeatureName, resp.GetFeatureName())
}
