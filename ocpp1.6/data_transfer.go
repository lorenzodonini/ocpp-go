package ocpp16

import (
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Data Transfer (CP -> CS / CS -> CP) --------------------
type DataTransferStatus string

const (
	DataTransferStatusAccepted         DataTransferStatus = "Accepted"
	DataTransferStatusRejected         DataTransferStatus = "Rejected"
	DataTransferStatusUnknownMessageId DataTransferStatus = "UnknownMessageId"
	DataTransferStatusUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

func isValidDataTransferStatus(fl validator.FieldLevel) bool {
	status := DataTransferStatus(fl.Field().String())
	switch status {
	case DataTransferStatusAccepted, DataTransferStatusRejected, DataTransferStatusUnknownMessageId, DataTransferStatusUnknownVendorId:
		return true
	default:
		return false
	}
}

type DataTransferRequest struct {
	VendorId  string      `json:"vendorId" validate:"required,max=255"`
	MessageId string      `json:"messageId,omitempty" validate:"max=50"`
	Data      interface{} `json:"data,omitempty"`
}

type DataTransferConfirmation struct {
	Status DataTransferStatus `json:"status" validate:"required,dataTransferStatus"`
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

func init() {
	_ = Validate.RegisterValidation("dataTransferStatus", isValidDataTransferStatus)
}
