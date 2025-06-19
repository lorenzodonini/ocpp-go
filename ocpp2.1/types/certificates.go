package types

import "gopkg.in/go-playground/validator.v9"

// Hash Algorithms
type HashAlgorithmType string

const (
	SHA256 HashAlgorithmType = "SHA256"
	SHA384 HashAlgorithmType = "SHA384"
	SHA512 HashAlgorithmType = "SHA512"
)

func isValidHashAlgorithmType(fl validator.FieldLevel) bool {
	algorithm := HashAlgorithmType(fl.Field().String())
	switch algorithm {
	case SHA256, SHA384, SHA512:
		return true
	default:
		return false
	}
}

// OCSPRequestDataType
type OCSPRequestDataType struct {
	HashAlgorithm  HashAlgorithmType `json:"hashAlgorithm" validate:"required,hashAlgorithm"`
	IssuerNameHash string            `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string            `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string            `json:"serialNumber" validate:"required,max=40"`
	ResponderURL   string            `json:"responderURL,omitempty" validate:"max=512"`
}

// CertificateHashDataType
type CertificateHashData struct {
	HashAlgorithm  HashAlgorithmType `json:"hashAlgorithm" validate:"required,hashAlgorithm"`
	IssuerNameHash string            `json:"issuerNameHash" validate:"required,max=128"`
	IssuerKeyHash  string            `json:"issuerKeyHash" validate:"required,max=128"`
	SerialNumber   string            `json:"serialNumber" validate:"required,max=40"`
}

// CertificateHashDataChain
type CertificateHashDataChain struct {
	CertificateType          CertificateUse        `json:"certificateType" validate:"required,certificateUse"`
	CertificateHashData      CertificateHashData   `json:"certificateHashData" validate:"required"`
	ChildCertificateHashData []CertificateHashData `json:"childCertificateHashData,omitempty" validate:"omitempty,dive"`
}

// Certificate15118EVStatus
type Certificate15118EVStatus string

const (
	Certificate15188EVStatusAccepted Certificate15118EVStatus = "Accepted"
	Certificate15118EVStatusFailed   Certificate15118EVStatus = "Failed"
)

func isValidCertificate15118EVStatus(fl validator.FieldLevel) bool {
	status := Certificate15118EVStatus(fl.Field().String())
	switch status {
	case Certificate15188EVStatusAccepted, Certificate15118EVStatusFailed:
		return true
	default:
		return false
	}
}

// Indicates the type of the signed certificate that is returned.
// When omitted the certificate is used for both the 15118 connection (if implemented) and the Charging Station to CSMS connection.
// This field is required when a typeOfCertificate was included in the SignCertificateRequest that requested this certificate to be signed AND both the 15118 connection and the Charging Station connection are implemented.
type CertificateSigningUse string

const (
	ChargingStationCert CertificateSigningUse = "ChargingStationCertificate"
	V2GCertificate      CertificateSigningUse = "V2GCertificate"
)

func isValidCertificateSigningUse(fl validator.FieldLevel) bool {
	status := CertificateSigningUse(fl.Field().String())
	switch status {
	case ChargingStationCert, V2GCertificate:
		return true
	default:
		return false
	}
}

// Indicates the type of the requested certificate.
// It is used in GetInstalledCertificateIdsRequest and InstallCertificateRequest messages.
type CertificateUse string

const (
	V2GRootCertificate          CertificateUse = "V2GRootCertificate"
	MORootCertificate           CertificateUse = "MORootCertificate"
	CSOSubCA1                   CertificateUse = "CSOSubCA1"
	CSOSubCA2                   CertificateUse = "CSOSubCA2"
	CSMSRootCertificate         CertificateUse = "CSMSRootCertificate"
	V2GCertificateChain         CertificateUse = "V2GCertificateChain"
	ManufacturerRootCertificate CertificateUse = "ManufacturerRootCertificate"
	OEMRootCertificate          CertificateUse = "OEMRootCertificate"
)

func isValidCertificateUse(fl validator.FieldLevel) bool {
	use := CertificateUse(fl.Field().String())
	switch use {
	case V2GRootCertificate, MORootCertificate, CSOSubCA1,
		CSOSubCA2, CSMSRootCertificate, V2GCertificateChain, ManufacturerRootCertificate, OEMRootCertificate:
		return true
	default:
		return false
	}
}

// Enumeration of the cryptographic method used to create the digital signature.
// The list is expected to grow in future OCPP releases to allow other signature methods used by Smart Meters.
type SignatureMethod string

const (
	SignatureECDSAP256SHA256 SignatureMethod = "ECDSAP256SHA256" // The encoded data is hashed with the SHA-256 hash function, and the hash value is then signed with the ECDSA algorithm using the NIST P-256 elliptic curve.
	SignatureECDSAP384SHA384 SignatureMethod = "ECDSAP384SHA384" // The encoded data is hashed with the SHA-384 hash function, and the hash value is then signed with the ECDSA algorithm using the NIST P-384 elliptic curve.
	SignatureECDSA192SHA256  SignatureMethod = "ECDSA192SHA256"  // The encoded data is hashed with the SHA-256 hash function, and the hash value is then signed with the ECDSA algorithm using a 192-bit elliptic curve.
)

func isValidSignatureMethod(fl validator.FieldLevel) bool {
	signature := SignatureMethod(fl.Field().String())
	switch signature {
	case SignatureECDSA192SHA256, SignatureECDSAP256SHA256, SignatureECDSAP384SHA384:
		return true
	default:
		return false
	}
}
