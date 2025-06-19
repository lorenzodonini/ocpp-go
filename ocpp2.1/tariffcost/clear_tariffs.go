package tariffcost

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

type TariffClearStatus string

const (
	TariffClearStatusAccepted     TariffClearStatus = "Accepted"
	TariffClearStatusRejected     TariffClearStatus = "Rejected"
	TariffClearStatusNotSupported TariffClearStatus = "NotSupported"
)

func isValidTariffClearStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch TariffClearStatus(status) {
	case TariffClearStatusAccepted, TariffClearStatusRejected, TariffClearStatusNotSupported:
		return true
	default:
		return false
	}
}

func init() {
	_ = ocppj.Validate.RegisterValidation("tariffClearStatus21", isValidTariffClearStatus)
}

// -------------------- Clear Tariffs (CSMS -> CS) --------------------

const ClearTariffs = "ClearTariffs"

// The field definition of the ClearTariffsRequest request payload sent by the CSMS to the Charging Station.
type ClearTariffsRequest struct {
	TariffIds []string `json:"tariffIds,omitempty" validate:"omitempty"`
	EvseId    int      `json:"evseId,omitempty" validate:"omitempty,gte=0"`
}

// This field definition of the ClearTariffs response payload, sent by the Charging Station to the CSMS in response to a ClearTariffsRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ClearTariffsResponse struct {
	ClearTariffsResult []ClearTariffsResult `json:"clearTariffsResult" validate:"required,min=1,dive"`
}

type ClearTariffsFeature struct{}

func (f ClearTariffsFeature) GetFeatureName() string {
	return ClearTariffs
}

func (f ClearTariffsFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ClearTariffsRequest{})
}

func (f ClearTariffsFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ClearTariffsResponse{})
}

func (r ClearTariffsRequest) GetFeatureName() string {
	return ClearTariffs
}

func (c ClearTariffsResponse) GetFeatureName() string {
	return ClearTariffs
}

// Creates a new NewClearTariffsRequest, containing all required fields. There are no optional fields for this message.
func NewClearTariffsRequest(tariffIds []string) *ClearTariffsRequest {
	return &ClearTariffsRequest{
		TariffIds: tariffIds,
	}
}

// Creates a new NewClearTariffsResponse, which doesn't contain any required or optional fields.
func NewClearTariffsResponse(results []ClearTariffsResult) *ClearTariffsResponse {
	return &ClearTariffsResponse{
		results,
	}
}
