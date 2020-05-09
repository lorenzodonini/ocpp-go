package ocpp2

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Get Installed Certificate IDs (CSMS -> CS) --------------------

// Status returned in response to GetInstalledCertificateIdsRequest, that indicates whether certificate signing has been accepted or rejected.
type GetInstalledCertificateStatus string

const (
	GetInstalledCertificateStatusAccepted GetInstalledCertificateStatus = "Accepted" // Normal successful completion (no errors).
	GetInstalledCertificateStatusNotFound GetInstalledCertificateStatus = "NotFound" // Requested resource not found
)

func isValidGetInstalledCertificateStatus(fl validator.FieldLevel) bool {
	status := GetInstalledCertificateStatus(fl.Field().String())
	switch status {
	case GetInstalledCertificateStatusAccepted, GetInstalledCertificateStatusNotFound:
		return true
	default:
		return false
	}
}

// The field definition of the GetInstalledCertificateIdsRequest PDU sent by the CSMS to the Charging Station.
type GetInstalledCertificateIdsRequest struct {
	TypeOfCertificate types.CertificateUse `json:"typeOfCertificate" validate:"required,certificateUse"`
}

// The field definition of the GetInstalledCertificateIdsResponse payload sent by the Charging Station to the CSMS in response to a GetInstalledCertificateIdsRequest.
type GetInstalledCertificateIdsConfirmation struct {
	Status              GetInstalledCertificateStatus `json:"status" validate:"required,getInstalledCertificateStatus"`
	CertificateHashData []types.CertificateHashData   `json:"certificateHashData,omitempty" validate:"omitempty,dive"`
}

// To facilitate the management of the Charging Station’s installed certificates, a method of retrieving the installed certificates is provided.
// The CSMS requests the Charging Station to send a list of installed certificates by sending a GetInstalledCertificateIdsRequest.
// The Charging Station responds with a GetInstalledCertificateIdsResponse.
type GetInstalledCertificateIdsFeature struct{}

func (f GetInstalledCertificateIdsFeature) GetFeatureName() string {
	return GetInstalledCertificateIdsFeatureName
}

func (f GetInstalledCertificateIdsFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetInstalledCertificateIdsRequest{})
}

func (f GetInstalledCertificateIdsFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetInstalledCertificateIdsConfirmation{})
}

func (r GetInstalledCertificateIdsRequest) GetFeatureName() string {
	return GetInstalledCertificateIdsFeatureName
}

func (c GetInstalledCertificateIdsConfirmation) GetFeatureName() string {
	return GetInstalledCertificateIdsFeatureName
}

// Creates a new GetInstalledCertificateIdsRequest, containing all required fields. There are no optional fields for this message.
func NewGetInstalledCertificateIdsRequest(typeOfCertificate types.CertificateUse) *GetInstalledCertificateIdsRequest {
	return &GetInstalledCertificateIdsRequest{TypeOfCertificate: typeOfCertificate}
}

// Creates a new ChangeAvailabilityConfirmation, containing all required fields. Additional optional fields may be set afterwards.
func NewGetInstalledCertificateIdsConfirmation(status GetInstalledCertificateStatus) *GetInstalledCertificateIdsConfirmation {
	return &GetInstalledCertificateIdsConfirmation{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("getInstalledCertificateStatus", isValidGetInstalledCertificateStatus)
}
