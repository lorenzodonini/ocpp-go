package tariffcost

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- NotifySettlement (CS -> CSMS) --------------------

const NotifySettlementFeatureName = "NotifySettlement"

type PaymentStatusEnumType string

const (
	PaymentStatusSettled   PaymentStatusEnumType = "Settled"
	PaymentStatusCanceled  PaymentStatusEnumType = "Canceled"
	PaymentStatusRejected  PaymentStatusEnumType = "Rejected"
	PaymentStatusFailed    PaymentStatusEnumType = "Failed"
)

func isValidPaymentStatus(fl validator.FieldLevel) bool {
	v := PaymentStatusEnumType(fl.Field().String())
	switch v {
	case PaymentStatusSettled, PaymentStatusCanceled, PaymentStatusRejected, PaymentStatusFailed:
		return true
	default:
		return false
	}
}

// AddressType contains a physical address. Used by NotifySettlement and VatNumberValidation.
type AddressType struct {
	Name       string `json:"name" validate:"required,max=50"`
	Address1   string `json:"address1" validate:"required,max=100"`
	Address2   string `json:"address2,omitempty" validate:"omitempty,max=100"`
	City       string `json:"city" validate:"required,max=100"`
	PostalCode string `json:"postalCode,omitempty" validate:"omitempty,max=20"`
	Country    string `json:"country" validate:"required,max=50"`
}

// The field definition of the NotifySettlementRequest request payload sent by the Charging Station to the CSMS.
type NotifySettlementRequest struct {
	TransactionId    string                `json:"transactionId" validate:"required,max=36"`
	PaymentReference string                `json:"paymentReference" validate:"required,max=40"`
	Status           PaymentStatusEnumType `json:"status" validate:"required,paymentStatus21"`
	StatusInfo       string                `json:"statusInfo,omitempty" validate:"omitempty,max=500"`
	SettlementAmount float64               `json:"settlementAmount" validate:"required"`
	SettlementTime   types.DateTime        `json:"settlementTime" validate:"required"`
	ReceiptId        string                `json:"receiptId,omitempty" validate:"omitempty,max=36"`
	ReceiptUrl       string                `json:"receiptUrl,omitempty" validate:"omitempty,max=2000"`
	VatCompany       *AddressType          `json:"vatCompany,omitempty" validate:"omitempty,dive"`
}

// This field definition of the NotifySettlementResponse response payload, sent by the CSMS to the Charging Station.
type NotifySettlementResponse struct {
}

type NotifySettlementFeature struct{}

func (f NotifySettlementFeature) GetFeatureName() string {
	return NotifySettlementFeatureName
}

func (f NotifySettlementFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifySettlementRequest{})
}

func (f NotifySettlementFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifySettlementResponse{})
}

func (r NotifySettlementRequest) GetFeatureName() string {
	return NotifySettlementFeatureName
}

func (c NotifySettlementResponse) GetFeatureName() string {
	return NotifySettlementFeatureName
}

// Creates a new NotifySettlementRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifySettlementRequest(transactionId string, paymentReference string, status PaymentStatusEnumType, settlementAmount float64, settlementTime types.DateTime) *NotifySettlementRequest {
	return &NotifySettlementRequest{
		TransactionId:    transactionId,
		PaymentReference: paymentReference,
		Status:           status,
		SettlementAmount: settlementAmount,
		SettlementTime:   settlementTime,
	}
}

// Creates a new NotifySettlementResponse, containing all required fields.
func NewNotifySettlementResponse() *NotifySettlementResponse {
	return &NotifySettlementResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("paymentStatus21", isValidPaymentStatus)
}