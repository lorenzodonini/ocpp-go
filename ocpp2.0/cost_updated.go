package ocpp2

import (
	"reflect"
)

// -------------------- Cost Updated (CSMS -> CS) --------------------

// The field definition of the CostUpdated request payload sent by the CSMS to the Charging Station.
type CostUpdatedRequest struct {
	TotalCost     float64 `json:"totalCost" validate:"required"`
	TransactionID string  `json:"transactionId" validate:"required,max=36"`
}

// This field definition of the CostUpdated confirmation payload, sent by the Charging Station to the CSMS in response to a CostUpdatedRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type CostUpdatedConfirmation struct {
}

// The driver wants to know how much the running total cost is, updated at a relevant interval, while a transaction is ongoing.
// To fulfill this requirement, the CSMS sends a CostUpdatedRequest to the Charging Station to update the current total cost, every Y seconds.
// Upon receipt of the CostUpdatedRequest, the Charging Station responds with a CostUpdatedResponse, then shows the updated cost to the driver.
type CostUpdatedFeature struct{}

func (f CostUpdatedFeature) GetFeatureName() string {
	return CostUpdatedFeatureName
}

func (f CostUpdatedFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(CostUpdatedRequest{})
}

func (f CostUpdatedFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(CostUpdatedConfirmation{})
}

func (r CostUpdatedRequest) GetFeatureName() string {
	return CostUpdatedFeatureName
}

func (c CostUpdatedConfirmation) GetFeatureName() string {
	return CostUpdatedFeatureName
}

// Creates a new CostUpdatedRequest, containing all required fields. There are no optional fields for this message.
func NewCostUpdatedRequest(totalCost float64, transactionID string) *CostUpdatedRequest {
	return &CostUpdatedRequest{TotalCost: totalCost, TransactionID: transactionID}
}

// Creates a new CostUpdatedConfirmation, which doesn't contain any required or optional fields.
func NewCostUpdatedConfirmation() *CostUpdatedConfirmation {
	return &CostUpdatedConfirmation{}
}
