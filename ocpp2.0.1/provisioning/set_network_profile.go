package provisioning

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Set Network Profile (CSMS -> CS) --------------------

const SetNetworkProfileFeatureName = "SetNetworkProfile"

// Enumeration of OCPP versions.
type OCPPVersion string

// OCPP transport mechanisms. SOAP is currently not a valid value for OCPP 2.0 (and is unsupported by this library).
type OCPPTransport string

// SetNetworkProfileType indicates the type of reset that the charging station or EVSE should perform,
// as requested by the CSMS in a SetNetworkProfileRequest.
type SetNetworkProfileType string

// Network interface.
type OCPPInterface string

// Type of VPN.
type VPNType string

// APN Authentication method.
type APNAuthentication string

// Result of a SetNetworkProfileRequest.
type SetNetworkProfileStatus string

const (
	OCPPVersion12 OCPPVersion = "OCPP12" // 1.2
	OCPPVersion15 OCPPVersion = "OCPP15" // 1.5
	OCPPVersion16 OCPPVersion = "OCPP16" // 1.6
	OCPPVersion20 OCPPVersion = "OCPP20" // 2.0

	OCPPTransportJSON OCPPTransport = "JSON" // Use JSON over WebSockets for transport of OCPP PDU’s
	OCPPTransportSOAP OCPPTransport = "SOAP" // Use SOAP for transport of OCPP PDU’s

	OCPPInterfaceWired0    OCPPInterface = "Wired0"
	OCPPInterfaceWired1    OCPPInterface = "Wired1"
	OCPPInterfaceWired2    OCPPInterface = "Wired2"
	OCPPInterfaceWired3    OCPPInterface = "Wired3"
	OCPPInterfaceWireless0 OCPPInterface = "Wireless0"
	OCPPInterfaceWireless1 OCPPInterface = "Wireless1"
	OCPPInterfaceWireless2 OCPPInterface = "Wireless2"
	OCPPInterfaceWireless3 OCPPInterface = "Wireless3"

	VPNTypeIKEv2 VPNType = "IKEv2"
	VPNTypeIPSec VPNType = "IPSec"
	VPNTypeL2TP  VPNType = "L2TP"
	VPNTypePPTP  VPNType = "PPTP"

	APNAuthenticationCHAP APNAuthentication = "CHAP"
	APNAuthenticationNone APNAuthentication = "NONE"
	APNAuthenticationPAP  APNAuthentication = "PAP"
	APNAuthenticationAuto APNAuthentication = "AUTO" // Sequentially try CHAP, PAP, NONE.

	SetNetworkProfileStatusAccepted SetNetworkProfileStatus = "Accepted"
	SetNetworkProfileStatusRejected SetNetworkProfileStatus = "Rejected"
	SetNetworkProfileStatusFailed   SetNetworkProfileStatus = "Failed"
)

func isValidOCPPVersion(fl validator.FieldLevel) bool {
	v := OCPPVersion(fl.Field().String())
	switch v {
	case OCPPVersion12, OCPPVersion15, OCPPVersion16, OCPPVersion20:
		return true
	default:
		return false
	}
}

func isValidOCPPTransport(fl validator.FieldLevel) bool {
	t := OCPPTransport(fl.Field().String())
	switch t {
	case OCPPTransportJSON, OCPPTransportSOAP:
		return true
	default:
		return false
	}
}

func isValidOCPPInterface(fl validator.FieldLevel) bool {
	i := OCPPInterface(fl.Field().String())
	switch i {
	case OCPPInterfaceWired0, OCPPInterfaceWired1, OCPPInterfaceWired2, OCPPInterfaceWired3,
		OCPPInterfaceWireless0, OCPPInterfaceWireless1, OCPPInterfaceWireless2, OCPPInterfaceWireless3:
		return true
	default:
		return false
	}
}

func isValidVPNType(fl validator.FieldLevel) bool {
	t := VPNType(fl.Field().String())
	switch t {
	case VPNTypeIKEv2, VPNTypeIPSec, VPNTypeL2TP, VPNTypePPTP:
		return true
	default:
		return false
	}
}

func isValidAPNAuthentication(fl validator.FieldLevel) bool {
	a := APNAuthentication(fl.Field().String())
	switch a {
	case APNAuthenticationAuto, APNAuthenticationCHAP, APNAuthenticationPAP, APNAuthenticationNone:
		return true
	default:
		return false
	}
}

func isValidSetNetworkProfileStatus(fl validator.FieldLevel) bool {
	status := SetNetworkProfileStatus(fl.Field().String())
	switch status {
	case SetNetworkProfileStatusAccepted, SetNetworkProfileStatusRejected, SetNetworkProfileStatusFailed:
		return true
	default:
		return false
	}
}

// VPN Configuration settings.
type VPN struct {
	Server   string  `json:"server" validate:"required,max=512"`          // VPN Server Address.
	User     string  `json:"user" validate:"required,max=20"`             // VPN User.
	Group    string  `json:"group,omitempty" validate:"omitempty,max=20"` // VPN group.
	Password string  `json:"password" validate:"required,max=20"`         // VPN Password.
	Key      string  `json:"key" validate:"required,max=255"`             // VPN shared secret.
	Type     VPNType `json:"type" validate:"required,vpnType"`            // Type of VPN.
}

type APN struct {
	APN                     string            `json:"apn" validate:"required,max=512"`                         // The Access Point Name as an URL.
	APNUsername             string            `json:"apnUserName,omitempty" validate:"omitempty,max=20"`       // APN username.
	APNPassword             string            `json:"apnPassword,omitempty" validate:"omitempty,max=20"`       // APN password.
	SimPin                  *int              `json:"simPin,omitempty" validate:"omitempty,gte=0"`             // SIM card pin code.
	PreferredNetwork        string            `json:"preferredNetwork,omitempty" validate:"omitempty,max=6"`   // Preferred network, written as MCC and MNC concatenated.
	UseOnlyPreferredNetwork bool              `json:"useOnlyPreferredNetwork,omitempty"`                       // Use only the preferred Network, do not dial in when not available.
	APNAuthentication       APNAuthentication `json:"apnAuthentication" validate:"required,apnAuthentication"` // Authentication method.
}

// NetworkConnectionProfile defines the functional and technical parameters of a communication link.
type NetworkConnectionProfile struct {
	OCPPVersion     OCPPVersion   `json:"ocppVersion" validate:"required,ocppVersion"`     // The OCPP version used for this communication function.
	OCPPTransport   OCPPTransport `json:"ocppTransport" validate:"required,ocppTransport"` // Defines the transport protocol (only OCPP-J is supported by this library).
	CSMSUrl         string        `json:"ocppCsmsUrl" validate:"required,max=512,url"`     // URL of the CSMS(s) that this Charging Station communicates with.
	MessageTimeout  int           `json:"messageTimeout" validate:"gte=-1"`                // Duration in seconds before a message send by the Charging Station via this network connection times out.
	SecurityProfile int           `json:"securityProfile"`                                 // The security profile used when connecting to the CSMS with this NetworkConnectionProfile.
	OCPPInterface   OCPPInterface `json:"ocppInterface" validate:"required,ocppInterface"` // Applicable Network Interface.
	VPN             *VPN          `json:"vpn,omitempty" validate:"omitempty"`              // Settings to be used to set up the VPN connection.
	APN             *APN          `json:"apn,omitempty" validate:"omitempty"`              // Collection of configuration data needed to make a data-connection over a cellular network.
}

// The field definition of the SetNetworkProfile request payload sent by the CSMS to the Charging Station.
type SetNetworkProfileRequest struct {
	ConfigurationSlot int                      `json:"configurationSlot" validate:"gte=0"` // Slot in which the configuration should be stored.
	ConnectionData    NetworkConnectionProfile `json:"connectionData" validate:"required"` // Connection details.
}

// Field definition of the SetNetworkProfile response payload, sent by the Charging Station to the CSMS in response to a SetNetworkProfileRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetNetworkProfileResponse struct {
	Status     SetNetworkProfileStatus `json:"status" validate:"required,setNetworkProfileStatus"`
	StatusInfo *types.StatusInfo       `json:"statusInfo" validate:"omitempty"`
}

// The CSMS may update the connection details on the Charging Station.
// For instance in preparation of a migration to a new CSMS. In order to achieve this,
// the CSMS sends a SetNetworkProfileRequest PDU containing an updated connection profile.
//
// The Charging station validates the content and stores the new data,
// eventually responding with a SetNetworkProfileResponse PDU.
// After completion of this use case, the Charging Station to CSMS connection data has been updated.
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

// Creates a new SetNetworkProfileRequest, containing all required fields. There are no optional fields for this message.
func NewSetNetworkProfileRequest(configurationSlot int, connectionData NetworkConnectionProfile) *SetNetworkProfileRequest {
	return &SetNetworkProfileRequest{ConfigurationSlot: configurationSlot, ConnectionData: connectionData}
}

// Creates a new SetNetworkProfileResponse, containing all required fields. Optional fields may be set afterwards.
func NewSetNetworkProfileResponse(status SetNetworkProfileStatus) *SetNetworkProfileResponse {
	return &SetNetworkProfileResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("ocppVersion", isValidOCPPVersion)
	_ = types.Validate.RegisterValidation("ocppTransport", isValidOCPPTransport)
	_ = types.Validate.RegisterValidation("ocppInterface", isValidOCPPInterface)
	_ = types.Validate.RegisterValidation("vpnType", isValidVPNType)
	_ = types.Validate.RegisterValidation("apnAuthentication", isValidAPNAuthentication)
	_ = types.Validate.RegisterValidation("setNetworkProfileStatus", isValidSetNetworkProfileStatus)
}
