package data

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
)

// -------------------- Data Transfer (CS -> CSMS / CSMS -> CS) --------------------

const DataTransferFeatureName = "DataTransfer"

// Status in DataTransferResponse messages.
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

// The field definition of the DataTransfer request payload sent by an endpoint to ther other endpoint.
type DataTransferRequest struct {
	MessageID string      `json:"messageId,omitempty" validate:"max=50"`
	Data      interface{} `json:"data,omitempty"`
	VendorID  string      `json:"vendorId" validate:"required,max=255"`
}

// This field definition of the DataTransfer response payload, sent by an endpoint in response to a DataTransferRequest, coming from the other endpoint.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type DataTransferResponse struct {
	Status     DataTransferStatus `json:"status" validate:"required,dataTransferStatus"`
	Data       interface{}        `json:"data,omitempty"`
	StatusInfo *types.StatusInfo  `json:"statusInfo,omitempty" validate:"omitempty"`
}

// If a CS needs to send information to the CSMS for a function not supported by OCPP, it SHALL use a DataTransfer message.
// The same functionality may also be offered the other way around, allowing a CSMS to send arbitrary custom commands to a CS.
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

// Creates a new DataTransferRequest, containing all required fields. Optional fields may be set afterwards.
func NewDataTransferRequest(vendorId string) *DataTransferRequest {
	return &DataTransferRequest{VendorID: vendorId}
}

// Creates a new DataTransferResponse. Optional fields may be set afterwards.
func NewDataTransferResponse(status DataTransferStatus) *DataTransferResponse {
	return &DataTransferResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("dataTransferStatus", isValidDataTransferStatus)
}
