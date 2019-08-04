package ocpp16

import "reflect"

// -------------------- Data Transfer --------------------
type DataTransferStatus string

const (
	DataTransferStatusAccepted         DataTransferStatus = "Accepted"
	DataTransferStatusRejected         DataTransferStatus = "Rejected"
	DataTransferStatusUnknownMessageId DataTransferStatus = "UnknownMessageId"
	DataTransferStatusUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

type DataTransferRequest struct {
	VendorId  string      `json:"vendorId" validate:"required,max=255"`
	MessageId string      `json:"messageId,omitempty" validate:"max=50"`
	Data      interface{} `json:"data,omitempty"`
}

type DataTransferConfirmation struct {
	Status DataTransferStatus `json:"status" validate:"required"`
	Data   interface{}        `json:"data,omitempty"`
}

type DataTransferFeature struct{}

func (f DataTransferFeature) GetFeatureName() string {
	return DataTransferFeatureName
}

func (f DataTransferFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(DataTransferRequest{})
}

func (f DataTransferFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(DataTransferConfirmation{})
}

func (r DataTransferRequest) GetFeatureName() string {
	return DataTransferFeatureName
}

func (c DataTransferConfirmation) GetFeatureName() string {
	return DataTransferFeatureName
}

func NewDataTransferRequest(vendorId string) *DataTransferRequest {
	return &DataTransferRequest{VendorId: vendorId}
}

func NewDataTransferConfirmation(status DataTransferStatus) *DataTransferConfirmation {
	return &DataTransferConfirmation{Status: status}
}
