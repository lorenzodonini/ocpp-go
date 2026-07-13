package smartcharging

import (
	"reflect"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- UsePriorityCharging (CSMS -> CS) --------------------

const UsePriorityChargingFeatureName = "UsePriorityCharging"

type PriorityChargingStatusEnumType string

const (
	PriorityChargingStatusAccepted  PriorityChargingStatusEnumType = "Accepted"
	PriorityChargingStatusRejected  PriorityChargingStatusEnumType = "Rejected"
	PriorityChargingStatusNoProfile PriorityChargingStatusEnumType = "NoProfile"
)

func isValidPriorityChargingStatus(fl validator.FieldLevel) bool {
	v := PriorityChargingStatusEnumType(fl.Field().String())
	switch v {
	case PriorityChargingStatusAccepted, PriorityChargingStatusRejected, PriorityChargingStatusNoProfile:
		return true
	default:
		return false
	}
}

// The field definition of the UsePriorityChargingRequest request payload sent by the CSMS to the Charging Station.
type UsePriorityChargingRequest struct {
	TransactionId string `json:"transactionId" validate:"required,max=36"`
	Activate      bool   `json:"activate" validate:"required"`
}

// This field definition of the UsePriorityChargingResponse response payload, sent by the Charging Station to the CSMS.
type UsePriorityChargingResponse struct {
	Status PriorityChargingStatusEnumType `json:"status" validate:"required,priorityChargingStatus21"`
}

type UsePriorityChargingFeature struct{}

func (f UsePriorityChargingFeature) GetFeatureName() string {
	return UsePriorityChargingFeatureName
}

func (f UsePriorityChargingFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UsePriorityChargingRequest{})
}

func (f UsePriorityChargingFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(UsePriorityChargingResponse{})
}

func (r UsePriorityChargingRequest) GetFeatureName() string {
	return UsePriorityChargingFeatureName
}

func (c UsePriorityChargingResponse) GetFeatureName() string {
	return UsePriorityChargingFeatureName
}

// Creates a new UsePriorityChargingRequest, containing all required fields.
func NewUsePriorityChargingRequest(transactionId string, activate bool) *UsePriorityChargingRequest {
	return &UsePriorityChargingRequest{TransactionId: transactionId, Activate: activate}
}

// Creates a new UsePriorityChargingResponse, containing all required fields.
func NewUsePriorityChargingResponse(status PriorityChargingStatusEnumType) *UsePriorityChargingResponse {
	return &UsePriorityChargingResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("priorityChargingStatus21", isValidPriorityChargingStatus)
}