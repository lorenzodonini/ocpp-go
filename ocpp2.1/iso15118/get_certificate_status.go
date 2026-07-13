package iso15118

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Certificate Status (CS -> CSMS) --------------------

const GetCertificateStatusFeatureName = "GetCertificateStatus"

// The field definition of the GetCertificateStatus request payload sent by the Charging Station to the CSMS.
type GetCertificateStatusRequest struct {
	OcspRequestData types.OCSPRequestDataType `json:"ocspRequestData" validate:"required"`
}

// This field definition of the GetCertificateStatus response payload, sent by the CSMS to the Charging Station in response to a GetCertificateStatusRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetCertificateStatusResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	OcspResult string              `json:"ocspResult,omitempty" validate:"omitempty,max=18000"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty"`
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

func (f GetCertificateStatusFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetCertificateStatusResponse{})
}

func (r GetCertificateStatusRequest) GetFeatureName() string {
	return GetCertificateStatusFeatureName
}

func (c GetCertificateStatusResponse) GetFeatureName() string {
	return GetCertificateStatusFeatureName
}

// Creates a new GetCertificateStatusRequest, containing all required fields. There are no optional fields for this message.
func NewGetCertificateStatusRequest(ocspRequestData types.OCSPRequestDataType) *GetCertificateStatusRequest {
	return &GetCertificateStatusRequest{OcspRequestData: ocspRequestData}
}

// Creates a new GetCertificateStatusResponse, containing all required fields. Optional fields may be set afterwards.
func NewGetCertificateStatusResponse(status types.GenericStatus) *GetCertificateStatusResponse {
	return &GetCertificateStatusResponse{Status: status}
}
