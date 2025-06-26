package tariffcost

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
)

type TariffGetStatus string

type TariffKind string

const (
	TariffGetStatusAccepted TariffGetStatus = "Accepted"
	TariffGetStatusRejected TariffGetStatus = "Rejected"
	TariffGetStatusNoTariff TariffGetStatus = "NoTariff"

	TariffKindDefaultTariff TariffKind = "DefaultTariff"
	TariffKindDriverTariff  TariffKind = "DriverTariff"
)

func isValidTariffKind(fl validator.FieldLevel) bool {
	switch TariffKind(fl.Field().String()) {
	case TariffKindDefaultTariff, TariffKindDriverTariff:
		return true
	default:
		return false
	}
}

func isValidTariffGetStatus(fl validator.FieldLevel) bool {
	switch TariffGetStatus(fl.Field().String()) {
	case TariffGetStatusAccepted, TariffGetStatusRejected, TariffGetStatusNoTariff:
		return true
	default:
		return false
	}
}

func init() {
	_ = ocppj.Validate.RegisterValidation("tariffKind21", isValidTariffKind)
	_ = ocppj.Validate.RegisterValidation("tariffGetStatus21", isValidTariffGetStatus)
}

type ClearTariffsResult struct {
	TariffId   string            `json:"tariffId,omitempty" validate:"omitempty,max=60"`
	Status     TariffClearStatus `json:"status" validate:"required,tariffClearStatus21"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}
