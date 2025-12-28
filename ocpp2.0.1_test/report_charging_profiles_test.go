package ocpp2_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Tests
func (suite *OcppV2TestSuite) TestReportChargingProfilesRequestValidation() {
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	chargingProfile := types.NewChargingProfile(1, 0, types.ChargingProfilePurposeTxDefaultProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{chargingSchedule})
	var requestTable = []GenericTestEntry{
		{smartcharging.ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{smartcharging.ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{smartcharging.ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{smartcharging.ReportChargingProfilesRequest{ChargingLimitSource: types.ChargingLimitSourceCSO, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, true},
		{smartcharging.ReportChargingProfilesRequest{ChargingLimitSource: types.ChargingLimitSourceCSO, ChargingProfile: []types.ChargingProfile{}}, false},
		{smartcharging.ReportChargingProfilesRequest{ChargingLimitSource: types.ChargingLimitSourceCSO}, false},
		{smartcharging.ReportChargingProfilesRequest{ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{smartcharging.ReportChargingProfilesRequest{}, false},
		{smartcharging.ReportChargingProfilesRequest{RequestID: -1, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{smartcharging.ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: -1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{smartcharging.ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: "invalidChargingLimitSource", Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{*chargingProfile}}, false},
		{smartcharging.ReportChargingProfilesRequest{RequestID: 42, ChargingLimitSource: types.ChargingLimitSourceCSO, Tbc: true, EvseID: 1, ChargingProfile: []types.ChargingProfile{
			*types.NewChargingProfile(1, -1, types.ChargingProfilePurposeTxDefaultProfile, types.ChargingProfileKindAbsolute, []types.ChargingSchedule{chargingSchedule})}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestReportChargingProfilesResponseValidation() {
	var responseTable = []GenericTestEntry{
		{smartcharging.ReportChargingProfilesResponse{}, true},
	}
	ExecuteGenericTestTable(suite, responseTable)
}

func (suite *OcppV2TestSuite) TestReportChargingProfilesE2EMocked() {
	wsId := "test_id"
	messageId := "1234"
	wsUrl := "someUrl"
	requestID := 42
	chargingLimitSource := types.ChargingLimitSourceEMS
	evseID := 1
	tbc := false
	chargingSchedule := types.ChargingSchedule{
		ID:                     1,
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	chargingProfile := types.ChargingProfile{
		ID:                     1,
		StackLevel:             0,
		ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile,
		ChargingProfileKind:    types.ChargingProfileKindAbsolute,
		ChargingSchedule:       []types.ChargingSchedule{chargingSchedule},
	}
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"chargingLimitSource":"%v","evseId":%v,"chargingProfile":[{"id":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":[{"id":%v,"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]}]}]`,
		messageId, smartcharging.ReportChargingProfilesFeatureName, requestID, chargingLimitSource, evseID, chargingProfile.ID, chargingProfile.StackLevel, chargingProfile.ChargingProfilePurpose, chargingProfile.ChargingProfileKind, chargingSchedule.ID, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	responseJson := fmt.Sprintf(`[3,"%v",{}]`, messageId)
	response := smartcharging.NewReportChargingProfilesResponse()
	channel := NewMockWebSocket(wsId)

	handler := &MockCSMSSmartChargingHandler{}
	handler.On("OnReportChargingProfiles", mock.AnythingOfType("string"), mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(1).(*smartcharging.ReportChargingProfilesRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(requestID, request.RequestID)
		suite.Equal(chargingLimitSource, request.ChargingLimitSource)
		suite.Equal(evseID, request.EvseID)
		suite.Equal(tbc, request.Tbc)
		suite.Require().Len(request.ChargingProfile, 1)
		suite.Equal(chargingProfile.ID, request.ChargingProfile[0].ID)
		suite.Equal(chargingProfile.StackLevel, request.ChargingProfile[0].StackLevel)
		suite.Equal(chargingProfile.ChargingProfilePurpose, request.ChargingProfile[0].ChargingProfilePurpose)
		suite.Equal(chargingProfile.ChargingProfileKind, request.ChargingProfile[0].ChargingProfileKind)
		suite.Require().Len(request.ChargingProfile[0].ChargingSchedule, 1)
		suite.Equal(chargingSchedule.ID, request.ChargingProfile[0].ChargingSchedule[0].ID)
		suite.Equal(chargingSchedule.StartSchedule.FormatTimestamp(), request.ChargingProfile[0].ChargingSchedule[0].StartSchedule.FormatTimestamp())
		suite.Equal(*chargingSchedule.Duration, *request.ChargingProfile[0].ChargingSchedule[0].Duration)
		suite.Equal(chargingSchedule.ChargingRateUnit, request.ChargingProfile[0].ChargingSchedule[0].ChargingRateUnit)
		suite.Equal(*chargingSchedule.MinChargingRate, *request.ChargingProfile[0].ChargingSchedule[0].MinChargingRate)
		suite.Require().Len(request.ChargingProfile[0].ChargingSchedule[0].ChargingSchedulePeriod, 1)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, request.ChargingProfile[0].ChargingSchedule[0].ChargingSchedulePeriod[0].StartPeriod)
		suite.Equal(chargingSchedule.ChargingSchedulePeriod[0].Limit, request.ChargingProfile[0].ChargingSchedule[0].ChargingSchedulePeriod[0].Limit)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	// Run test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	r, err := suite.chargingStation.ReportChargingProfiles(requestID, chargingLimitSource, evseID, []types.ChargingProfile{chargingProfile}, func(request *smartcharging.ReportChargingProfilesRequest) {
		request.Tbc = tbc
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(r)
}

func (suite *OcppV2TestSuite) TestReportChargingProfilesInvalidEndpoint() {
	messageId := defaultMessageId
	requestID := 42
	chargingLimitSource := types.ChargingLimitSourceEMS
	evseID := 1
	tbc := false
	chargingSchedule := types.ChargingSchedule{
		StartSchedule:          types.NewDateTime(time.Now()),
		Duration:               newInt(600),
		ChargingRateUnit:       types.ChargingRateUnitWatts,
		MinChargingRate:        newFloat(6.0),
		ChargingSchedulePeriod: []types.ChargingSchedulePeriod{types.NewChargingSchedulePeriod(0, 10.0)},
	}
	chargingProfile := types.ChargingProfile{
		ID:                     1,
		StackLevel:             0,
		ChargingProfilePurpose: types.ChargingProfilePurposeTxDefaultProfile,
		ChargingProfileKind:    types.ChargingProfileKindAbsolute,
		ChargingSchedule:       []types.ChargingSchedule{chargingSchedule},
	}
	request := smartcharging.NewReportChargingProfilesRequest(requestID, chargingLimitSource, evseID, []types.ChargingProfile{chargingProfile})
	request.Tbc = tbc
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"requestId":%v,"chargingLimitSource":"%v","evseId":%v,"chargingProfile":[{"id":%v,"stackLevel":%v,"chargingProfilePurpose":"%v","chargingProfileKind":"%v","chargingSchedule":[{"startSchedule":"%v","duration":%v,"chargingRateUnit":"%v","minChargingRate":%v,"chargingSchedulePeriod":[{"startPeriod":%v,"limit":%v}]}]}]}]`,
		messageId, smartcharging.ReportChargingProfilesFeatureName, requestID, chargingLimitSource, evseID, chargingProfile.ID, chargingProfile.StackLevel, chargingProfile.ChargingProfilePurpose, chargingProfile.ChargingProfileKind, chargingSchedule.StartSchedule.FormatTimestamp(), *chargingSchedule.Duration, chargingSchedule.ChargingRateUnit, *chargingSchedule.MinChargingRate, chargingSchedule.ChargingSchedulePeriod[0].StartPeriod, chargingSchedule.ChargingSchedulePeriod[0].Limit)
	testUnsupportedRequestFromCentralSystem(suite, request, requestJson, messageId)
}
