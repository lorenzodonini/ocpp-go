package data

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// -------------------- Data Transfer (bidirectional) --------------------

const DataTransferFeatureName = "DataTransfer"

// DataTransferStatus represents the status of a data transfer.
type DataTransferStatus string

const (
	DataTransferStatusAccepted         DataTransferStatus = "Accepted"
	DataTransferStatusRejected         DataTransferStatus = "Rejected"
	DataTransferStatusUnknownMessageId DataTransferStatus = "UnknownMessageId"
	DataTransferStatusUnknownVendorId  DataTransferStatus = "UnknownVendorId"
)

func isValidDataTransferStatus(fl validator.FieldLevel) bool {
	s := DataTransferStatus(fl.Field().String())
	switch s {
	case DataTransferStatusAccepted, DataTransferStatusRejected,
		DataTransferStatusUnknownMessageId, DataTransferStatusUnknownVendorId:
		return true
	default:
		return false
	}
}

// DataTransferRequest is the payload for a DataTransfer request.
type DataTransferRequest struct {
	VendorId   string                `json:"vendorId" validate:"required,max=255"`
	MessageId  string                `json:"messageId,omitempty" validate:"omitempty,max=50"`
	Data       interface{}           `json:"data,omitempty"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// DataTransferResponse is the payload for a DataTransfer response.
type DataTransferResponse struct {
	Status     DataTransferStatus    `json:"status" validate:"required,dataTransferStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty"`
	Data       interface{}           `json:"data,omitempty"`
	CustomData *types.CustomDataType `json:"customData,omitempty" validate:"omitempty"`
}

// DataTransferFeature represents the DataTransfer feature.
type DataTransferFeature struct{}

func (f DataTransferFeature) GetFeatureName() string {
	return DataTransferFeatureName
}

func (f DataTransferFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(DataTransferRequest{})
}

func (f DataTransferFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(DataTransferResponse{})
}

func (r DataTransferRequest) GetFeatureName() string {
	return DataTransferFeatureName
}

func (c DataTransferResponse) GetFeatureName() string {
	return DataTransferFeatureName
}

// NewDataTransferRequest creates a new DataTransferRequest.
func NewDataTransferRequest(vendorId string) *DataTransferRequest {
	return &DataTransferRequest{VendorId: vendorId}
}

// NewDataTransferResponse creates a new DataTransferResponse.
func NewDataTransferResponse(status DataTransferStatus) *DataTransferResponse {
	return &DataTransferResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("dataTransferStatus21", isValidDataTransferStatus)
}
