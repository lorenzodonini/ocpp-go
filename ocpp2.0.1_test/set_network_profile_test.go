package ocpp2_test

import (
	"fmt"

	"github.com/stretchr/testify/mock"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// Test

func (suite *OcppV2TestSuite) TestVPNTypeValidation() {
	var requestTable = []GenericTestEntry{
		{provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, true},
		{provisioning.VPN{Server: "someServer", User: "user1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, true},
		{provisioning.VPN{Server: "someServer", User: "user1", Password: "deadc0de", Key: "deadbeef"}, false},
		{provisioning.VPN{Server: "someServer", User: "user1", Password: "deadc0de", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", User: "user1", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{User: "user1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{}, false},
		{provisioning.VPN{Server: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", User: ">20..................", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", User: "user1", Group: ">20..................", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: ">20..................", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: ">255............................................................................................................................................................................................................................................................", Type: provisioning.VPNTypeIPSec}, false},
		{provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: "invalidType"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestAPNTypeValidation() {
	var requestTable = []GenericTestEntry{
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: provisioning.APNAuthenticationAuto}, true},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", APNAuthentication: provisioning.APNAuthenticationAuto}, true},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), APNAuthentication: provisioning.APNAuthenticationAuto}, true},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", APNAuthentication: provisioning.APNAuthenticationAuto}, true},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNAuthentication: provisioning.APNAuthenticationAuto}, true},
		{provisioning.APN{APN: "internet.t-mobile", APNAuthentication: provisioning.APNAuthenticationAuto}, true},
		{provisioning.APN{APN: "internet.t-mobile"}, false},
		{provisioning.APN{APNAuthentication: provisioning.APNAuthenticationAuto}, false},
		{provisioning.APN{APN: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: provisioning.APNAuthenticationAuto}, false},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: ">20..................", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: provisioning.APNAuthenticationAuto}, false},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: ">20..................", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: provisioning.APNAuthenticationAuto}, false},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(-1), PreferredNetwork: ">6.....", UseOnlyPreferredNetwork: true, APNAuthentication: provisioning.APNAuthenticationAuto}, false},
		{provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: "invalidApnAuthentication"}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetNetworkProfileRequestValidation() {
	vpn := &provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}
	apn := &provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: provisioning.APNAuthenticationAuto}
	var requestTable = []GenericTestEntry{
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: apn}}, true},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn}}, true},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0}}, true},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, OCPPInterface: provisioning.OCPPInterfaceWired0}}, true},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", OCPPInterface: provisioning.OCPPInterfaceWired0}}, true},
		{provisioning.SetNetworkProfileRequest{ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", OCPPInterface: provisioning.OCPPInterfaceWired0}}, true},
		{provisioning.SetNetworkProfileRequest{ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767"}}, false},
		{provisioning.SetNetworkProfileRequest{ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, OCPPInterface: provisioning.OCPPInterfaceWired0}}, false},
		{provisioning.SetNetworkProfileRequest{ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, CSMSUrl: "http://someUrl:8767", OCPPInterface: provisioning.OCPPInterfaceWired0}}, false},
		{provisioning.SetNetworkProfileRequest{ConnectionData: provisioning.NetworkConnectionProfile{OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", OCPPInterface: provisioning.OCPPInterfaceWired0}}, false},
		{provisioning.SetNetworkProfileRequest{}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: -1, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: "OCPP01", OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: "ProtoBuf", CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: -2, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: "invalidInterface", VPN: vpn, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: &provisioning.VPN{}, APN: apn}}, false},
		{provisioning.SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: vpn, APN: &provisioning.APN{}}}, false},
	}
	ExecuteGenericTestTable(suite, requestTable)
}

func (suite *OcppV2TestSuite) TestSetNetworkProfileResponseValidation() {
	var confirmationTable = []GenericTestEntry{
		{provisioning.SetNetworkProfileResponse{Status: provisioning.SetNetworkProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{provisioning.SetNetworkProfileResponse{Status: provisioning.SetNetworkProfileStatusRejected, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{provisioning.SetNetworkProfileResponse{Status: provisioning.SetNetworkProfileStatusFailed, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{provisioning.SetNetworkProfileResponse{Status: provisioning.SetNetworkProfileStatusAccepted}, true},
		{provisioning.SetNetworkProfileResponse{}, false},
		{provisioning.SetNetworkProfileResponse{Status: provisioning.SetNetworkProfileStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
		{provisioning.SetNetworkProfileResponse{Status: "invalidSetNetworkProfileStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
	}
	ExecuteGenericTestTable(suite, confirmationTable)
}

func (suite *OcppV2TestSuite) TestSetNetworkProfileE2EMocked() {
	wsId := "test_id"
	messageId := defaultMessageId
	wsUrl := "someUrl"
	configurationSlot := 2
	vpn := provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}
	apn := provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: false, APNAuthentication: provisioning.APNAuthenticationAuto}
	data := provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: &vpn, APN: &apn}
	status := provisioning.SetNetworkProfileStatusAccepted
	statusInfo := types.NewStatusInfo("200", "")
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"configurationSlot":%v,"connectionData":{"ocppVersion":"%v","ocppTransport":"%v","ocppCsmsUrl":"%v","messageTimeout":%v,"securityProfile":%v,"ocppInterface":"%v","vpn":{"server":"%v","user":"%v","group":"%v","password":"%v","key":"%v","type":"%v"},"apn":{"apn":"%v","apnUserName":"%v","apnPassword":"%v","simPin":%v,"preferredNetwork":"%v","apnAuthentication":"%v"}}}]`,
		messageId, provisioning.SetNetworkProfileFeatureName, configurationSlot, data.OCPPVersion, data.OCPPTransport, data.CSMSUrl, data.MessageTimeout, data.SecurityProfile, data.OCPPInterface, vpn.Server, vpn.User, vpn.Group, vpn.Password, vpn.Key, vpn.Type, apn.APN, apn.APNUsername, apn.APNPassword, *apn.SimPin, apn.PreferredNetwork, apn.APNAuthentication)
	responseJson := fmt.Sprintf(`[3,"%v",{"status":"%v","statusInfo":{"reasonCode":"%v"}}]`, messageId, status, statusInfo.ReasonCode)
	resetResponse := provisioning.NewSetNetworkProfileResponse(status)
	resetResponse.StatusInfo = statusInfo
	channel := NewMockWebSocket(wsId)

	handler := &MockChargingStationProvisioningHandler{}
	handler.On("OnSetNetworkProfile", mock.Anything).Return(resetResponse, nil).Run(func(args mock.Arguments) {
		request, ok := args.Get(0).(*provisioning.SetNetworkProfileRequest)
		suite.Require().True(ok)
		suite.Require().NotNil(request)
		suite.Equal(configurationSlot, request.ConfigurationSlot)
		suite.Equal(data.OCPPVersion, request.ConnectionData.OCPPVersion)
		suite.Equal(data.OCPPTransport, request.ConnectionData.OCPPTransport)
		suite.Equal(data.CSMSUrl, request.ConnectionData.CSMSUrl)
		suite.Equal(data.MessageTimeout, request.ConnectionData.MessageTimeout)
		suite.Equal(data.SecurityProfile, request.ConnectionData.SecurityProfile)
		suite.Equal(data.OCPPInterface, request.ConnectionData.OCPPInterface)
		suite.Require().NotNil(request.ConnectionData.VPN)
		suite.Equal(vpn.Server, request.ConnectionData.VPN.Server)
		suite.Equal(vpn.User, request.ConnectionData.VPN.User)
		suite.Equal(vpn.Group, request.ConnectionData.VPN.Group)
		suite.Equal(vpn.Password, request.ConnectionData.VPN.Password)
		suite.Equal(vpn.Key, request.ConnectionData.VPN.Key)
		suite.Equal(vpn.Type, request.ConnectionData.VPN.Type)
		suite.Require().NotNil(request.ConnectionData.APN)
		suite.Equal(apn.APN, request.ConnectionData.APN.APN)
		suite.Equal(apn.APNUsername, request.ConnectionData.APN.APNUsername)
		suite.Equal(apn.APNPassword, request.ConnectionData.APN.APNPassword)
		suite.Equal(*apn.SimPin, *request.ConnectionData.APN.SimPin)
		suite.Equal(apn.PreferredNetwork, request.ConnectionData.APN.PreferredNetwork)
		suite.Equal(apn.UseOnlyPreferredNetwork, request.ConnectionData.APN.UseOnlyPreferredNetwork)
		suite.Equal(apn.APNAuthentication, request.ConnectionData.APN.APNAuthentication)
	})
	setupDefaultCSMSHandlers(suite, expectedCSMSOptions{clientId: wsId, rawWrittenMessage: []byte(requestJson), forwardWrittenMessage: true})
	setupDefaultChargingStationHandlers(suite, expectedChargingStationOptions{serverUrl: wsUrl, clientId: wsId, createChannelOnStart: true, channel: channel, rawWrittenMessage: []byte(responseJson), forwardWrittenMessage: true}, handler)
	// Run Test
	suite.csms.Start(8887, "somePath")
	err := suite.chargingStation.Start(wsUrl)
	suite.Require().Nil(err)
	resultChannel := make(chan bool, 1)
	err = suite.csms.SetNetworkProfile(wsId, func(resp *provisioning.SetNetworkProfileResponse, err error) {
		suite.Require().Nil(err)
		suite.Require().NotNil(resp)
		suite.Equal(status, resp.Status)
		suite.Equal(statusInfo.ReasonCode, resp.StatusInfo.ReasonCode)
		resultChannel <- true
	}, configurationSlot, data)
	suite.Require().Nil(err)
	result := <-resultChannel
	suite.True(result)
}

func (suite *OcppV2TestSuite) TestSetNetworkProfileInvalidEndpoint() {
	messageId := defaultMessageId
	configurationSlot := 2
	vpn := provisioning.VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: provisioning.VPNTypeIPSec}
	apn := provisioning.APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: newInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: false, APNAuthentication: provisioning.APNAuthenticationAuto}
	data := provisioning.NetworkConnectionProfile{OCPPVersion: provisioning.OCPPVersion20, OCPPTransport: provisioning.OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: provisioning.OCPPInterfaceWired0, VPN: &vpn, APN: &apn}
	request := provisioning.NewSetNetworkProfileRequest(configurationSlot, data)
	requestJson := fmt.Sprintf(`[2,"%v","%v",{"configurationSlot":%v,"connectionData":{"ocppVersion":"%v","ocppTransport":"%v","ocppCsmsUrl":"%v","messageTimeout":%v,"securityProfile":%v,"ocppInterface":"%v","vpn":{"server":"%v","user":"%v","group":"%v","password":"%v","key":"%v","type":"%v"},"apn":{"apn":"%v","apnUserName":"%v","apnPassword":"%v","simPin":%v,"preferredNetwork":"%v","apnAuthentication":"%v"}}}]`,
		messageId, provisioning.SetNetworkProfileFeatureName, configurationSlot, data.OCPPVersion, data.OCPPTransport, data.CSMSUrl, data.MessageTimeout, data.SecurityProfile, data.OCPPInterface, vpn.Server, vpn.User, vpn.Group, vpn.Password, vpn.Key, vpn.Type, apn.APN, apn.APNUsername, apn.APNPassword, *apn.SimPin, apn.PreferredNetwork, apn.APNAuthentication)
	testUnsupportedRequestFromChargingStation(suite, request, requestJson, messageId)
}
