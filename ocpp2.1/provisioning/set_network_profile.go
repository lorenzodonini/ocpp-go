package provisioning

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Set Network Profile (CSMS -> CS) --------------------

const SetNetworkProfileFeatureName = "SetNetworkProfile"

// SetNetworkProfileStatus represents the status of setting a network profile.
type SetNetworkProfileStatus string

const (
	SetNetworkProfileStatusAccepted SetNetworkProfileStatus = "Accepted"
	SetNetworkProfileStatusRejected SetNetworkProfileStatus = "Rejected"
	SetNetworkProfileStatusFailed   SetNetworkProfileStatus = "Failed"
)

func isValidSetNetworkProfileStatus(fl validator.FieldLevel) bool {
	s := SetNetworkProfileStatus(fl.Field().String())
	switch s {
	case SetNetworkProfileStatusAccepted, SetNetworkProfileStatusRejected, SetNetworkProfileStatusFailed:
		return true
	default:
		return false
	}
}

// OCPPVersion represents OCPP protocol versions.
type OCPPVersion string

const (
	OCPPVersion12  OCPPVersion = "OCPP12"
	OCPPVersion15  OCPPVersion = "OCPP15"
	OCPPVersion16  OCPPVersion = "OCPP16"
	OCPPVersion20  OCPPVersion = "OCPP20"
	OCPPVersion201 OCPPVersion = "OCPP201"
	OCPPVersion21  OCPPVersion = "OCPP21"
)

// OCPPTransport represents the transport protocol.
type OCPPTransport string

const (
	OCPPTransportJSON OCPPTransport = "JSON"
	OCPPTransportSOAP OCPPTransport = "SOAP"
)

// OCPPInterface represents the interface for OCPP communication.
type OCPPInterface string

const (
	OCPPInterfaceWired0    OCPPInterface = "Wired0"
	OCPPInterfaceWired1    OCPPInterface = "Wired1"
	OCPPInterfaceWired2    OCPPInterface = "Wired2"
	OCPPInterfaceWired3    OCPPInterface = "Wired3"
	OCPPInterfaceWireless0 OCPPInterface = "Wireless0"
	OCPPInterfaceWireless1 OCPPInterface = "Wireless1"
	OCPPInterfaceWireless2 OCPPInterface = "Wireless2"
	OCPPInterfaceWireless3 OCPPInterface = "Wireless3"
	OCPPInterfaceAny       OCPPInterface = "Any"
)

// APNAuthentication represents APN authentication types.
type APNAuthentication string

const (
	APNAuthenticationPAP  APNAuthentication = "PAP"
	APNAuthenticationCHAP APNAuthentication = "CHAP"
	APNAuthenticationNone APNAuthentication = "NONE"
	APNAuthenticationAuto APNAuthentication = "AUTO"
)

// VPNType represents VPN types.
type VPNType string

const (
	VPNTypeIKEv2 VPNType = "IKEv2"
	VPNTypeIPSec VPNType = "IPSec"
	VPNTypeL2TP  VPNType = "L2TP"
	VPNTypePPTP  VPNType = "PPTP"
)

// APN represents Access Point Name configuration.
type APN struct {
	APN                     string                `json:"apn" validate:"required,max=2000"`
	APNUserName             string                `json:"apnUserName,omitempty" validate:"omitempty,max=50"`
	APNPassword             string                `json:"apnPassword,omitempty" validate:"omitempty,max=64"`
	SimPin                  *int                  `json:"simPin,omitempty"`
	PreferredNetwork        string                `json:"preferredNetwork,omitempty" validate:"omitempty,max=6"`
	UseOnlyPreferredNetwork bool                  `json:"useOnlyPreferredNetwork,omitempty"`
	APNAuthentication       APNAuthentication     `json:"apnAuthentication" validate:"required"`
	CustomData              *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// VPN represents VPN configuration.
type VPN struct {
	Server     string                `json:"server" validate:"required,max=2000"`
	User       string                `json:"user" validate:"required,max=50"`
	Password   string                `json:"password" validate:"required,max=64"`
	Key        string                `json:"key" validate:"required,max=255"`
	Type       VPNType               `json:"type" validate:"required"`
	Group      string                `json:"group,omitempty" validate:"omitempty,max=50"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// NetworkConnectionProfile contains network connection configuration.
type NetworkConnectionProfile struct {
	OCPPVersion       OCPPVersion           `json:"ocppVersion,omitempty" validate:"omitempty"`
	OCPPTransport     OCPPTransport         `json:"ocppTransport" validate:"required"`
	OCPPCsmsUrl       string                `json:"ocppCsmsUrl" validate:"required,max=2000"`
	MessageTimeout    int                   `json:"messageTimeout" validate:"required,gte=0"`
	SecurityProfile   int                   `json:"securityProfile" validate:"gte=0,lte=3"`
	OCPPInterface     OCPPInterface         `json:"ocppInterface" validate:"required"`
	Identity          string                `json:"identity,omitempty" validate:"omitempty,max=48"`
	BasicAuthPassword string                `json:"basicAuthPassword,omitempty" validate:"omitempty,max=64"`
	VPN               *VPN                  `json:"vpn,omitempty" validate:"omitempty"`
	APN               *APN                  `json:"apn,omitempty" validate:"omitempty"`
	CustomData        *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// SetNetworkProfileRequest is the payload for a SetNetworkProfile request from CSMS to CS.
type SetNetworkProfileRequest struct {
	ConfigurationSlot int                      `json:"configurationSlot" validate:"gte=0"`
	ConnectionData    NetworkConnectionProfile `json:"connectionData" validate:"required"`
	CustomData        *types.CustomDataType    `json:"customData,omitempty" validate:"omitempty"`
}

// SetNetworkProfileResponse is the payload for a SetNetworkProfile response from CS to CSMS.
type SetNetworkProfileResponse struct {
	Status     SetNetworkProfileStatus `json:"status" validate:"required,setNetworkProfileStatus21"`
	StatusInfo *types.StatusInfo       `json:"statusInfo,omitempty" validate:"omitempty"`
	CustomData *types.CustomDataType   `json:"customData,omitempty" validate:"omitempty"`
}

// SetNetworkProfileFeature represents the SetNetworkProfile feature.
type SetNetworkProfileFeature struct{}

func (f SetNetworkProfileFeature) GetFeatureName() string {
	return SetNetworkProfileFeatureName
}

func (f SetNetworkProfileFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetNetworkProfileRequest{})
}

func (f SetNetworkProfileFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetNetworkProfileResponse{})
}

func (r SetNetworkProfileRequest) GetFeatureName() string {
	return SetNetworkProfileFeatureName
}

func (c SetNetworkProfileResponse) GetFeatureName() string {
	return SetNetworkProfileFeatureName
}

// NewSetNetworkProfileRequest creates a new SetNetworkProfileRequest.
func NewSetNetworkProfileRequest(configurationSlot int, connectionData NetworkConnectionProfile) *SetNetworkProfileRequest {
	return &SetNetworkProfileRequest{ConfigurationSlot: configurationSlot, ConnectionData: connectionData}
}

// NewSetNetworkProfileResponse creates a new SetNetworkProfileResponse.
func NewSetNetworkProfileResponse(status SetNetworkProfileStatus) *SetNetworkProfileResponse {
	return &SetNetworkProfileResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("setNetworkProfileStatus21", isValidSetNetworkProfileStatus)
}
