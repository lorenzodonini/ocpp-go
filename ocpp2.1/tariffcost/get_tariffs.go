package tariffcost

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Get Tariffs (CSMS -> CS) --------------------

const GetTariffsFeatureName = "GetTariffs"

// The field definition of the GetTariffsRequest request payload sent by the CSMS to the Charging Station.
type GetTariffsRequest struct {
	EvseId int `json:"evseId" validate:"required,gte=0"`
}

// This field definition of the GetTariffs response payload, sent by the Charging Station to the CSMS in response to a CostUpdatedRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type GetTariffsResponse struct {
	Status            TariffGetStatus    `json:"status" validate:"required,tariffGetStatus21"`
	TariffAssignments []TariffAssignment `json:"tariffAssignments,omitempty" validate:"omitempty,dive"`
	StatusInfo        *types.StatusInfo  `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

type TariffAssignment struct {
	TariffId   string          `json:"tariffId" validate:"required"` // The ID of the tariff.
	TariffKind string          `json:"tariffKind" validate:"required,tariffKind21"`
	ValidFrom  *types.DateTime `json:"validFrom,omitempty" validate:"omitempty"`
	EvseIds    []int           `json:"evseIds,omitempty" validate:"omitempty"`
	IdTokens   []string        `json:"idTokens,omitempty" validate:"omitempty"`
}

type GetTariffsFeature struct{}

func (f GetTariffsFeature) GetFeatureName() string {
	return GetTariffsFeatureName
}

func (f GetTariffsFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(GetTariffsRequest{})
}

func (f GetTariffsFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(GetTariffsResponse{})
}

func (r GetTariffsRequest) GetFeatureName() string {
	return GetTariffsFeatureName
}

func (c GetTariffsResponse) GetFeatureName() string {
	return GetTariffsFeatureName
}

// Creates a new GetTariffsRequest, containing all required fields. There are no optional fields for this message.
func NewGetTariffsRequest(evseId int) *GetTariffsRequest {
	return &GetTariffsRequest{
		evseId,
	}
}

// Creates a new NewGetTariffsResponse, which doesn't contain any required or optional fields.
func NewGetTariffsResponse(status TariffGetStatus) *GetTariffsResponse {
	return &GetTariffsResponse{
		Status: status,
	}
}
