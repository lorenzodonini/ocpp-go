package meter

import (
	"reflect"
	"testing"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
	"github.com/stretchr/testify/suite"
)

type meterTestSuite struct {
	suite.Suite
}

func (suite *meterTestSuite) TestMeterValuesRequestValidation() {
	var requestTable = []tests.GenericTestEntry{
		{MeterValuesRequest{EvseID: 1, MeterValue: []types.MeterValue{{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}}}, true},
		{MeterValuesRequest{MeterValue: []types.MeterValue{{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}}}, true},
		{MeterValuesRequest{EvseID: 1, MeterValue: []types.MeterValue{}}, false},
		{MeterValuesRequest{EvseID: 1}, false},
		{MeterValuesRequest{EvseID: 1, MeterValue: []types.MeterValue{{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{{Value: 3.14, Context: "invalidContext", Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}}}, false},
		{MeterValuesRequest{EvseID: -1, MeterValue: []types.MeterValue{{Timestamp: types.DateTime{Time: time.Now()}, SampledValue: []types.SampledValue{{Value: 3.14, Context: types.ReadingContextTransactionEnd, Measurand: types.MeasurandPowerActiveExport, Phase: types.PhaseL2, Location: types.LocationBody}}}}}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *meterTestSuite) TestMeterValuesConfirmationValidation() {
	var responseTable = []tests.GenericTestEntry{
		{MeterValuesResponse{}, true},
	}
	tests.ExecuteGenericTestTable(suite.T(), responseTable)
}

func (suite *meterTestSuite) TestMeterValuesFeature() {
	feature := MeterValuesFeature{}
	suite.Equal(MeterValuesFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(MeterValuesRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(MeterValuesResponse{}), feature.GetResponseType())
}

func (suite *meterTestSuite) TestNewMeterValuesRequest() {
	mv := []types.MeterValue{{
		Timestamp:    types.DateTime{Time: time.Now()},
		SampledValue: []types.SampledValue{{Value: 3.14}},
	}}
	req := NewMeterValuesRequest(1, mv)
	suite.NotNil(req)
	suite.Equal(MeterValuesFeatureName, req.GetFeatureName())
	suite.Equal(1, req.EvseID)
	suite.Equal(mv, req.MeterValue)
}

func (suite *meterTestSuite) TestNewMeterValuesResponse() {
	resp := NewMeterValuesResponse()
	suite.NotNil(resp)
	suite.Equal(MeterValuesFeatureName, resp.GetFeatureName())
}

func TestMeterSuite(t *testing.T) {
	suite.Run(t, new(meterTestSuite))
}