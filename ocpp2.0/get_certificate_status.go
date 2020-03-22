package ocpp2

import (
	"reflect"
)

// -------------------- Get Certificate Status (CS -> CSMS) --------------------

// The field definition of the GetCertificateStatus request payload sent by the Charging Station to the CSMS.
type GetCertificateStatusRequest struct {
	OcspRequestData OCSPRequestDataType `json:"ocspRequestData" validate:"required"`
}

// This field definition of the GetCertificateStatus confirmation payload, sent by the CSMS to the Charging Station in response to a GetCertificateStatusRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetCertificateStatusConfirmation struct {
	Status     GenericStatus `json:"status" validate:"required,genericStatus"`
	OcspResult string        `json:"ocspResult,omitempty" validate:"omitempty,max=5500"`
}

// For 15118 certificate installation on EVs, the Charging Station requests the CSMS to provide the OCSP certificate
// status for its 15118 certificates.
// The CSMS responds with a GetCertificateStatusResponse, containing the OCSP certificate status.
// The status indicator in the GetCertificateStatusResponse indicates whether or not the CSMS was successful in retrieving the certificate status.
// It does NOT indicate the validity of the certificate.
type GetCertificateStatusFeature struct{}

func (f GetCertificateStatusFeature) GetFeatureName() string {
	return GetCertificateStatusFeatureName
}

func (f GetCertificateStatusFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetCertificateStatusRequest{})
}

func (f GetCertificateStatusFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(GetCertificateStatusConfirmation{})
}

func (r GetCertificateStatusRequest) GetFeatureName() string {
	return GetCertificateStatusFeatureName
}

func (c GetCertificateStatusConfirmation) GetFeatureName() string {
	return GetCertificateStatusFeatureName
}

// Creates a new GetCertificateStatusRequest, containing all required fields. There are no optional fields for this message.
func NewGetCertificateStatusRequest(ocspRequestData OCSPRequestDataType) *GetCertificateStatusRequest {
	return &GetCertificateStatusRequest{OcspRequestData: ocspRequestData}
}

// Creates a new GetCertificateStatusConfirmation, containing all required fields. Optional fields may be set afterwards.
func NewGetCertificateStatusConfirmation(status GenericStatus) *GetCertificateStatusConfirmation {
	return &GetCertificateStatusConfirmation{Status: status}
}
