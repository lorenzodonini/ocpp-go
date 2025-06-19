package v2x

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- NotifyAllowedEnergyTransfer (CSMS -> CS) --------------------

const NotifyAllowedEnergyTransfer = "NotifyAllowedEnergyTransfer"

// The field definition of the NotifyAllowedEnergyTransferRequest request payload sent by the CSMS to the Charging Station.
type NotifyAllowedEnergyTransferRequest struct {
	TransactionId         string                     `json:"transactionId" validate:"required,max=36"`
	AllowedEnergyTransfer []types.EnergyTransferMode `json:"allowedEnergyTransfer" validate:"required,energyTransferMode21"`
}

// This field definition of the NotifyAllowedEnergyTransferResponse
type NotifyAllowedEnergyTransferResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus21"`
	StatusInfo *types.StatusInfo   `json:"statusInfo" validate:"omitempty,dive"`
}

type NotifyAllowedEnergyTransferFeature struct{}

func (f NotifyAllowedEnergyTransferFeature) GetFeatureName() string {
	return NotifyAllowedEnergyTransfer
}

func (f NotifyAllowedEnergyTransferFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyAllowedEnergyTransferRequest{})
}

func (f NotifyAllowedEnergyTransferFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyAllowedEnergyTransferResponse{})
}

func (r NotifyAllowedEnergyTransferRequest) GetFeatureName() string {
	return NotifyAllowedEnergyTransfer
}

func (c NotifyAllowedEnergyTransferResponse) GetFeatureName() string {
	return NotifyAllowedEnergyTransfer
}

// Creates a new NotifyAllowedEnergyTransferRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyAllowedEnergyTransferRequest(transactionId string, allowedEnergyTransfers ...types.EnergyTransferMode) *NotifyAllowedEnergyTransferRequest {
	return &NotifyAllowedEnergyTransferRequest{
		TransactionId:         transactionId,
		AllowedEnergyTransfer: allowedEnergyTransfers,
	}
}

// Creates a new NotifyAllowedEnergyTransferResponse, containing all required fields. Optional fields may be set afterwards.
func NewNotifyAllowedEnergyTransferResponse(status types.GenericStatus) *NotifyAllowedEnergyTransferResponse {
	return &NotifyAllowedEnergyTransferResponse{Status: status}
}
