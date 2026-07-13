package tariffcost

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- VatNumberValidation (CS -> CSMS) --------------------

const VatNumberValidationFeatureName = "VatNumberValidation"

// The field definition of the VatNumberValidationRequest request payload sent by the Charging Station to the CSMS.
type VatNumberValidationRequest struct {
	VatNumber string `json:"vatNumber" validate:"required,max=20"`
	EvseId    *int   `json:"evseId,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the VatNumberValidationResponse response payload, sent by the CSMS to the Charging Station.
type VatNumberValidationResponse struct {
	VatNumber  string              `json:"vatNumber" validate:"required,max=20"`
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty,dive"`
	Company    *AddressType        `json:"company,omitempty" validate:"omitempty,dive"`
}

type VatNumberValidationFeature struct{}

func (f VatNumberValidationFeature) GetFeatureName() string {
	return VatNumberValidationFeatureName
}

func (f VatNumberValidationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(VatNumberValidationRequest{})
}

func (f VatNumberValidationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(VatNumberValidationResponse{})
}

func (r VatNumberValidationRequest) GetFeatureName() string {
	return VatNumberValidationFeatureName
}

func (c VatNumberValidationResponse) GetFeatureName() string {
	return VatNumberValidationFeatureName
}

// Creates a new VatNumberValidationRequest, containing all required fields. Optional fields may be set afterwards.
func NewVatNumberValidationRequest(vatNumber string) *VatNumberValidationRequest {
	return &VatNumberValidationRequest{VatNumber: vatNumber}
}

// Creates a new VatNumberValidationResponse, containing all required fields. Optional fields may be set afterwards.
func NewVatNumberValidationResponse(vatNumber string, status types.GenericStatus) *VatNumberValidationResponse {
	return &VatNumberValidationResponse{VatNumber: vatNumber, Status: status}
}