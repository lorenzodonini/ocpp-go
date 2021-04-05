// The security functional block contains OCPP 2.0 features aimed at providing E2E security between a CSMS and a Charging station.
package security

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Security profile.
type CSMSHandler interface {
	// OnSecurityEventNotification is called on the CSMS whenever a SecurityEventNotificationRequest is received from a charging station.
	OnSecurityEventNotification(chargingStationID string, request *SecurityEventNotificationRequest) (response *SecurityEventNotificationResponse, err error)
	// OnSignCertificate is called on the CSMS whenever a SignCertificateRequest is received from a charging station.
	OnSignCertificate(chargingStationID string, request *SignCertificateRequest) (response *SignCertificateResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Security profile.
type ChargingStationHandler interface {
	// OnCertificateSigned is called on a charging station whenever a CertificateSignedRequest is received from the CSMS.
	OnCertificateSigned(request *CertificateSignedRequest) (response *CertificateSignedResponse, err error)
}

const ProfileName = "security"

var Profile = ocpp.NewProfile(
	ProfileName,
	CertificateSignedFeature{},
	SecurityEventNotificationFeature{},
	SignCertificateFeature{},
)
