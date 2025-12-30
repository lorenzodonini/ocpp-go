package tariffcost

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

type TariffSetStatus string

const (
	TariffSetStatusAccepted              TariffSetStatus = "Accepted"
	TariffSetStatusRejected              TariffSetStatus = "Rejected"
	TariffSetStatusTooManyElements       TariffSetStatus = "TooManyElements"
	TariffSetStatusDuplicateTariffId     TariffSetStatus = "DuplicateTariffId"
	TariffSetStatusConditionNotSupported TariffSetStatus = "ConditionNotSupported"
)

func init() {
	_ = ocppj.Validate.RegisterValidation("tariffSetStatus", isValidTariffSetStatus)
}

func isValidTariffSetStatus(level validator.FieldLevel) bool {
	switch TariffSetStatus(level.Field().String()) {
	case TariffSetStatusAccepted,
		TariffSetStatusRejected,
		TariffSetStatusTooManyElements,
		TariffSetStatusDuplicateTariffId,
		TariffSetStatusConditionNotSupported:
		return true
	default:
		return false
	}
}

// -------------------- Set Default Tariff (CSMS -> CS) --------------------

const SetDefaultTariffFeatureName = "SetDefaultTariff"

// The field definition of the SetDefaultTariff request payload sent by the CSMS to the Charging Station.
type SetDefaultTariffRequest struct {
	EvseId int          `json:"evseId" validate:"required,gte=0"`
	Tariff types.Tariff `json:"tariff" validate:"required,dive"`
}

// This field definition of the SetDefaultTariff response payload, sent by the Charging Station to the CSMS in response to a SetDefaultTariffRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type SetDefaultTariffResponse struct {
	Status     TariffSetStatus   `json:"status" validate:"required,tariffSetStatus"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

// The driver wants to know how much the running total cost is, updated at a relevant interval, while a transaction is ongoing.
// To fulfill this requirement, the CSMS sends a SetDefaultTariffRequest to the Charging Station to update the current total cost, every Y seconds.
// Upon receipt of the SetDefaultTariffRequest, the Charging Station responds with a SetDefaultTariffResponse, then shows the updated cost to the driver.
type SetDefaultTariffFeature struct{}

func (f SetDefaultTariffFeature) GetFeatureName() string {
	return SetDefaultTariffFeatureName
}

func (f SetDefaultTariffFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(SetDefaultTariffRequest{})
}

func (f SetDefaultTariffFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(SetDefaultTariffResponse{})
}

func (r SetDefaultTariffRequest) GetFeatureName() string {
	return SetDefaultTariffFeatureName
}

func (c SetDefaultTariffResponse) GetFeatureName() string {
	return SetDefaultTariffFeatureName
}

// Creates a new SetDefaultTariffRequest, containing all required fields. There are no optional fields for this message.
func NewSetDefaultTariffRequest(evseId int, tariff types.Tariff) *SetDefaultTariffRequest {
	return &SetDefaultTariffRequest{
		EvseId: evseId,
		Tariff: tariff,
	}
}

// Creates a new SetDefaultTariffResponse, which doesn't contain any required or optional fields.
func NewSetDefaultTariffResponse(status TariffSetStatus) *SetDefaultTariffResponse {
	return &SetDefaultTariffResponse{
		Status: status,
	}
}
