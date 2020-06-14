// The ISO 15118 functional block contains OCPP 2.0 features that allow:
//
// - communication between EV and an EVSE
//
// - support for certificate-based authentication and authorization at the charging station, i.e. plug and charge
package iso15118

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 ISO 15118 profile.
type CSMSHandler interface {
	// OnGet15118EVCertificate is called on the CSMS whenever a Get15118EVCertificateRequest is received from a charging station.
	OnGet15118EVCertificate(chargingStationID string, request *Get15118EVCertificateRequest) (response *Get15118EVCertificateResponse, err error)
	// OnGetCertificateStatus is called on the CSMS whenever a GetCertificateStatusRequest is received from a charging station.
	OnGetCertificateStatus(chargingStationID string, request *GetCertificateStatusRequest) (response *GetCertificateStatusResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 ISO 15118 profile.
type ChargingStationHandler interface {
	// OnDeleteCertificate is called on a charging station whenever a DeleteCertificateRequest is received from the CSMS.
	OnDeleteCertificate(request *DeleteCertificateRequest) (response *DeleteCertificateResponse, err error)
	// OnGetInstalledCertificateIds is called on a charging station whenever a GetInstalledCertificateIdsRequest is received from the CSMS.
	OnGetInstalledCertificateIds(request *GetInstalledCertificateIdsRequest) (response *GetInstalledCertificateIdsResponse, err error)
	// OnInstallCertificate is called on a charging station whenever an InstallCertificateRequest is received from the CSMS.
	OnInstallCertificate(request *InstallCertificateRequest) (response *InstallCertificateResponse, err error)
}

const ProfileName = "iso15118"

var Profile = ocpp.NewProfile(
	ProfileName,
	DeleteCertificateFeature{},
	Get15118EVCertificateFeature{},
	GetCertificateStatusFeature{},
	GetInstalledCertificateIdsFeature{},
	InstallCertificateFeature{},
)
