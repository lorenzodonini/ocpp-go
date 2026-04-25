package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/tests"
)

func (suite *provisioningTestSuite) TestVPNTypeValidation() {
	var requestTable = []tests.GenericTestEntry{
		{VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, true},
		{VPN{Server: "someServer", User: "user1", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, true},
		{VPN{Server: "someServer", User: "user1", Password: "deadc0de", Key: "deadbeef"}, false},
		{VPN{Server: "someServer", User: "user1", Password: "deadc0de", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", User: "user1", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{User: "user1", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{}, false},
		{VPN{Server: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", User: ">20..................", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", User: "user1", Group: ">20..................", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", User: "user1", Group: "group1", Password: ">20..................", Key: "deadbeef", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: ">255............................................................................................................................................................................................................................................................", Type: VPNTypeIPSec}, false},
		{VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: "invalidType"}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *provisioningTestSuite) TestAPNTypeValidation() {
	var requestTable = []tests.GenericTestEntry{
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: APNAuthenticationAuto}, true},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", APNAuthentication: APNAuthenticationAuto}, true},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), APNAuthentication: APNAuthenticationAuto}, true},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", APNAuthentication: APNAuthenticationAuto}, true},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNAuthentication: APNAuthenticationAuto}, true},
		{APN{APN: "internet.t-mobile", APNAuthentication: APNAuthenticationAuto}, true},
		{APN{APN: "internet.t-mobile"}, false},
		{APN{APNAuthentication: APNAuthenticationAuto}, false},
		{APN{APN: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: APNAuthenticationAuto}, false},
		{APN{APN: "internet.t-mobile", APNUsername: ">20..................", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: APNAuthenticationAuto}, false},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: ">20..................", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: APNAuthenticationAuto}, false},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(-1), PreferredNetwork: ">6.....", UseOnlyPreferredNetwork: true, APNAuthentication: APNAuthenticationAuto}, false},
		{APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: "invalidApnAuthentication"}, false},
	}
	tests.ExecuteGenericTestTable(suite.T(), requestTable)
}

func (suite *provisioningTestSuite) TestSetNetworkProfileRequestValidation() {
	t := suite.T()
	vpn := &VPN{Server: "someServer", User: "user1", Group: "group1", Password: "deadc0de", Key: "deadbeef", Type: VPNTypeIPSec}
	apn := &APN{APN: "internet.t-mobile", APNUsername: "user1", APNPassword: "deadc0de", SimPin: tests.NewInt(1234), PreferredNetwork: "26201", UseOnlyPreferredNetwork: true, APNAuthentication: APNAuthenticationAuto}
	var requestTable = []tests.GenericTestEntry{
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: apn}}, true},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn}}, true},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0}}, true},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, OCPPInterface: OCPPInterfaceWired0}}, true},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", OCPPInterface: OCPPInterfaceWired0}}, true},
		{SetNetworkProfileRequest{ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", OCPPInterface: OCPPInterfaceWired0}}, true},
		{SetNetworkProfileRequest{ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767"}}, false},
		{SetNetworkProfileRequest{ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, OCPPInterface: OCPPInterfaceWired0}}, false},
		{SetNetworkProfileRequest{ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, CSMSUrl: "http://someUrl:8767", OCPPInterface: OCPPInterfaceWired0}}, false},
		{SetNetworkProfileRequest{ConnectionData: NetworkConnectionProfile{OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", OCPPInterface: OCPPInterfaceWired0}}, false},
		{SetNetworkProfileRequest{}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: -1, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: "OCPP01", OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: "ProtoBuf", CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: ">512.............................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................................", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: -2, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: "invalidInterface", VPN: vpn, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: &VPN{}, APN: apn}}, false},
		{SetNetworkProfileRequest{ConfigurationSlot: 2, ConnectionData: NetworkConnectionProfile{OCPPVersion: OCPPVersion20, OCPPTransport: OCPPTransportJSON, CSMSUrl: "http://someUrl:8767", MessageTimeout: 30, SecurityProfile: 1, OCPPInterface: OCPPInterfaceWired0, VPN: vpn, APN: &APN{}}}, false},
	}
	tests.ExecuteGenericTestTable(t, requestTable)
}

func (suite *provisioningTestSuite) TestSetNetworkProfileResponseValidation() {
	t := suite.T()
	var confirmationTable = []tests.GenericTestEntry{
		{SetNetworkProfileResponse{Status: SetNetworkProfileStatusAccepted, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetNetworkProfileResponse{Status: SetNetworkProfileStatusRejected, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetNetworkProfileResponse{Status: SetNetworkProfileStatusFailed, StatusInfo: types.NewStatusInfo("200", "")}, true},
		{SetNetworkProfileResponse{Status: SetNetworkProfileStatusAccepted}, true},
		{SetNetworkProfileResponse{}, false},
		{SetNetworkProfileResponse{Status: SetNetworkProfileStatusAccepted, StatusInfo: types.NewStatusInfo("", "")}, false},
		{SetNetworkProfileResponse{Status: "invalidSetNetworkProfileStatus", StatusInfo: types.NewStatusInfo("200", "")}, false},
	}
	tests.ExecuteGenericTestTable(t, confirmationTable)
}

func (suite *provisioningTestSuite) TestSetNetworkProfileFeature() {
	feature := SetNetworkProfileFeature{}
	suite.Equal(SetNetworkProfileFeatureName, feature.GetFeatureName())
	suite.Equal(reflect.TypeOf(SetNetworkProfileRequest{}), feature.GetRequestType())
	suite.Equal(reflect.TypeOf(SetNetworkProfileResponse{}), feature.GetResponseType())
}

func (suite *provisioningTestSuite) TestNewSetNetworkProfileRequest() {
	connectionData := NetworkConnectionProfile{
		OCPPVersion:   OCPPVersion20,
		OCPPTransport: OCPPTransportJSON,
		CSMSUrl:       "http://csms:8080",
		OCPPInterface: OCPPInterfaceWired0,
	}
	req := NewSetNetworkProfileRequest(1, connectionData)
	suite.NotNil(req)
	suite.Equal(SetNetworkProfileFeatureName, req.GetFeatureName())
	suite.Equal(1, req.ConfigurationSlot)
	suite.Equal(connectionData, req.ConnectionData)
}

func (suite *provisioningTestSuite) TestNewSetNetworkProfileResponse() {
	resp := NewSetNetworkProfileResponse(SetNetworkProfileStatusAccepted)
	suite.NotNil(resp)
	suite.Equal(SetNetworkProfileFeatureName, resp.GetFeatureName())
	suite.Equal(SetNetworkProfileStatusAccepted, resp.Status)
}
