package tariffcost

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

type TariffChangeStatus string

const (
	TariffChangeStatusAccepted              TariffChangeStatus = "Accepted"
	TariffChangeStatusRejected              TariffChangeStatus = "Rejected"
	TariffChangeStatusTooManyElements       TariffChangeStatus = "TooManyElements"
	TariffChangeStatusConditionNotSupported TariffChangeStatus = "ConditionNotSupported"
	TariffChangeStatusTxNotFound            TariffChangeStatus = "TxNotFound"
	TariffChangeStatusNoCurrencyChange      TariffChangeStatus = "NoCurrencyChange"
)

func isValidTariffChangeStatus(fl validator.FieldLevel) bool {
	switch TariffChangeStatus(fl.Field().String()) {
	case TariffChangeStatusAccepted,
		TariffChangeStatusRejected,
		TariffChangeStatusTooManyElements,
		TariffChangeStatusConditionNotSupported,
		TariffChangeStatusTxNotFound,
		TariffChangeStatusNoCurrencyChange:
		return true
	default:
		return false
	}
}

func init() {
	_ = ocppj.Validate.RegisterValidation("tariffChangeStatus", isValidTariffChangeStatus)
}

// -------------------- Change Transaction Tariff (CSMS -> CS) --------------------

const ChangeTransactionTariff = "ChangeTransactionTariff"

// The field definition of the ChangeTransactionTariff request payload sent by the CSMS to the Charging Station.
type ChangeTransactionTariffRequest struct {
	TransactionId string       `json:"transactionId" validate:"required,max=36"`
	Tariff        types.Tariff `json:"tariff" validate:"required,dive"`
}

// This field definition of the ChangeTransactionTariff response payload, sent by the Charging Station to the CSMS in response to a ChangeTransactionTariffRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type ChangeTransactionTariffResponse struct {
	Status     TariffSetStatus   `json:"status" validate:"required,tariffChangeStatus"`
	StatusInfo *types.StatusInfo `json:"statusInfo,omitempty" validate:"omitempty,dive"`
}

// The driver wants to know how much the running total cost is, updated at a relevant interval, while a transaction is ongoing.
// To fulfill this requirement, the CSMS sends a ChangeTransactionTariffRequest to the Charging Station to update the current total cost, every Y seconds.
// Upon receipt of the ChangeTransactionTariffRequest, the Charging Station responds with a ChangeTransactionTariffResponse, then shows the updated cost to the driver.
type ChangeTransactionTariffFeature struct{}

func (f ChangeTransactionTariffFeature) GetFeatureName() string {
	return ChangeTransactionTariff
}

func (f ChangeTransactionTariffFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(ChangeTransactionTariffRequest{})
}

func (f ChangeTransactionTariffFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(ChangeTransactionTariffResponse{})
}

func (r ChangeTransactionTariffRequest) GetFeatureName() string {
	return ChangeTransactionTariff
}

func (c ChangeTransactionTariffResponse) GetFeatureName() string {
	return ChangeTransactionTariff
}

// Creates a new ChangeTransactionTariffRequest, containing all required fields. There are no optional fields for this message.
func NewChangeTransactionTariffRequest(transactionId string, tariff types.Tariff) *ChangeTransactionTariffRequest {
	return &ChangeTransactionTariffRequest{
		TransactionId: transactionId,
		Tariff:        tariff,
	}
}

// Creates a new ChangeTransactionTariffResponse, which doesn't contain any required or optional fields.
func NewChangeTransactionTariffResponse(status TariffSetStatus) *ChangeTransactionTariffResponse {
	return &ChangeTransactionTariffResponse{
		Status: status,
	}
}
