package iso15118

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- GetCertificateChainStatus (CS -> CSMS) --------------------

const GetCertificateChainStatusFeatureName = "GetCertificateChainStatus"

type CertificateStatusSourceEnumType string

const (
	CertificateStatusSourceCRL  CertificateStatusSourceEnumType = "CRL"
	CertificateStatusSourceOCSP CertificateStatusSourceEnumType = "OCSP"
)

func isValidCertificateStatusSource(fl validator.FieldLevel) bool {
	v := CertificateStatusSourceEnumType(fl.Field().String())
	switch v {
	case CertificateStatusSourceCRL, CertificateStatusSourceOCSP:
		return true
	default:
		return false
	}
}

type CertificateStatusEnumType string

const (
	CertChainStatusGood    CertificateStatusEnumType = "Good"
	CertChainStatusRevoked CertificateStatusEnumType = "Revoked"
	CertChainStatusUnknown CertificateStatusEnumType = "Unknown"
	CertChainStatusFailed  CertificateStatusEnumType = "Failed"
)

func isValidCertificateStatus(fl validator.FieldLevel) bool {
	v := CertificateStatusEnumType(fl.Field().String())
	switch v {
	case CertChainStatusGood, CertChainStatusRevoked, CertChainStatusUnknown, CertChainStatusFailed:
		return true
	default:
		return false
	}
}

type CertificateStatusRequestInfoType struct {
	RequestId      int                             `json:"requestId" validate:"required"`
	Source         CertificateStatusSourceEnumType `json:"source" validate:"required,certificateStatusSource21"`
	OcspRequestData *types.OCSPRequestDataType     `json:"ocspRequestData,omitempty" validate:"omitempty,dive"`
	Certificate    string                          `json:"certificate,omitempty" validate:"omitempty,max=5500"`
}

type CertificateStatusType struct {
	RequestId  int                             `json:"requestId" validate:"required"`
	Status     CertificateStatusEnumType       `json:"status" validate:"required,certificateStatus21"`
	Source     CertificateStatusSourceEnumType `json:"source" validate:"required,certificateStatusSource21"`
	OcspResult string                          `json:"ocspResult,omitempty" validate:"omitempty,max=5500"`
	NextUpdate *types.DateTime                 `json:"nextUpdate,omitempty" validate:"omitempty"`
	StatusInfo *types.StatusInfo               `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

// The field definition of the GetCertificateChainStatusRequest request payload sent by the Charging Station to the CSMS.
type GetCertificateChainStatusRequest struct {
	CertificateStatusRequestInfo []CertificateStatusRequestInfoType `json:"certificateStatusRequestInfo" validate:"required,min=1,dive"`
}

// This field definition of the GetCertificateChainStatusResponse response payload, sent by the CSMS to the Charging Station.
type GetCertificateChainStatusResponse struct {
	CertificateStatus []CertificateStatusType `json:"certificateStatus,omitempty" validate:"omitempty,dive"`
}

type GetCertificateChainStatusFeature struct{}

func (f GetCertificateChainStatusFeature) GetFeatureName() string {
	return GetCertificateChainStatusFeatureName
}

func (f GetCertificateChainStatusFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetCertificateChainStatusRequest{})
}

func (f GetCertificateChainStatusFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetCertificateChainStatusResponse{})
}

func (r GetCertificateChainStatusRequest) GetFeatureName() string {
	return GetCertificateChainStatusFeatureName
}

func (c GetCertificateChainStatusResponse) GetFeatureName() string {
	return GetCertificateChainStatusFeatureName
}

// Creates a new GetCertificateChainStatusRequest, containing all required fields.
func NewGetCertificateChainStatusRequest(info []CertificateStatusRequestInfoType) *GetCertificateChainStatusRequest {
	return &GetCertificateChainStatusRequest{CertificateStatusRequestInfo: info}
}

// Creates a new GetCertificateChainStatusResponse, containing all required fields. Optional fields may be set afterwards.
func NewGetCertificateChainStatusResponse() *GetCertificateChainStatusResponse {
	return &GetCertificateChainStatusResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("certificateStatusSource21", isValidCertificateStatusSource)
	_ = types.Validate.RegisterValidation("certificateStatus21", isValidCertificateStatus)
}