package firmware

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// -------------------- Publish Firmware (CSMS -> CS) --------------------

const UpdateFirmwareFeatureName = "UpdateFirmware"

// Indicates whether the Charging Station was able to accept the request.
type UpdateFirmwareStatus string

const (
	UpdateFirmwareStatusAccepted           UpdateFirmwareStatus = "Accepted"
	UpdateFirmwareStatusRejected           UpdateFirmwareStatus = "Rejected"
	UpdateFirmwareStatusAcceptedCanceled   UpdateFirmwareStatus = "AcceptedCanceled"
	UpdateFirmwareStatusInvalidCertificate UpdateFirmwareStatus = "InvalidCertificate"
	UpdateFirmwareStatusRevokedCertificate UpdateFirmwareStatus = "RevokedCertificate"
)

func isValidUpdateFirmwareStatus(fl validator.FieldLevel) bool {
	status := UpdateFirmwareStatus(fl.Field().String())
	switch status {
	case UpdateFirmwareStatusAccepted,
		UpdateFirmwareStatusRejected,
		UpdateFirmwareStatusAcceptedCanceled,
		UpdateFirmwareStatusInvalidCertificate,
		UpdateFirmwareStatusRevokedCertificate:
		return true
	default:
		return false
	}
}

// Represents a copy of the firmware that can be loaded/updated on the Charging Station.
type Firmware struct {
	Location           string          `json:"location" validate:"required,max=512,uri"`         // URI defining the origin of the firmware.
	RetrieveDateTime   *types.DateTime `json:"retrieveDateTime" validate:"required"`             // Date and time at which the firmware shall be retrieved.
	InstallDateTime    *types.DateTime `json:"installDateTime,omitempty" validate:"omitempty"`   // Date and time at which the firmware shall be installed.
	SigningCertificate string          `json:"signingCertificate,omitempty" validate:"max=5500"` // Certificate with which the firmware was signed. PEM encoded X.509 certificate.
	Signature          string          `json:"signature,omitempty" validate:"max=800"`           // Base64 encoded firmware signature.
}

// The field definition of the UpdateFirmware request payload sent by the CSMS to the Charging Station.
type UpdateFirmwareRequest struct {
	Retries       *int     `json:"retries,omitempty" validate:"omitempty,gte=0"`       // This specifies how many times Charging Station must try to download the firmware before giving up. If this field is not present, it is left to Charging Station to decide how many times it wants to retry.
	RetryInterval *int     `json:"retryInterval,omitempty" validate:"omitempty,gte=0"` // The interval in seconds after which a retry may be attempted. If this field is not present, it is left to Charging Station to decide how long to wait between attempts.
	RequestID     int      `json:"requestId" validate:"gte=0"`                         // The Id of the request.
	Firmware      Firmware `json:"firmware" validate:"required"`                       // Specifies the firmware to be updated on the Charging Station.
}

// This field definition of the UpdateFirmware response payload, sent by the Charging Station to the CSMS in response to a UpdateFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type UpdateFirmwareResponse struct {
	Status     UpdateFirmwareStatus `json:"status" validate:"required,updateFirmwareStatus"`
	StatusInfo *types.StatusInfo    `json:"statusInfo,omitempty" validate:"omitempty"`
}

// A CSMS may instruct a Charging Station to update its firmware, by downloading and installing a new version.
// The CSMS sends an UpdateFirmwareRequest message that contains the location of the firmware,
// the time after which it should be retrieved, and information on how many times the
// Charging Station should retry downloading the firmware.
//
// The Charging station responds with an UpdateFirmwareResponse and then starts downloading the firmware.
// During the download/install procedure, the charging station shall notify the CSMS of its current status
// by sending FirmwareStatusNotification messages.
type UpdateFirmwareFeature struct{}

func (f UpdateFirmwareFeature) GetFeatureName() string {
	return UpdateFirmwareFeatureName
}

func (f UpdateFirmwareFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UpdateFirmwareRequest{})
}

func (f UpdateFirmwareFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(UpdateFirmwareResponse{})
}

func (r UpdateFirmwareRequest) GetFeatureName() string {
	return UpdateFirmwareFeatureName
}

func (c UpdateFirmwareResponse) GetFeatureName() string {
	return UpdateFirmwareFeatureName
}

// Creates a new UpdateFirmwareRequest,  containing all required fields. Optional fields may be set afterwards.
func NewUpdateFirmwareRequest(requestID int, firmware Firmware) *UpdateFirmwareRequest {
	return &UpdateFirmwareRequest{RequestID: requestID, Firmware: firmware}
}

// Creates a new UpdateFirmwareResponse, containing all required fields. Optional fields may be set afterwards.
func NewUpdateFirmwareResponse(status UpdateFirmwareStatus) *UpdateFirmwareResponse {
	return &UpdateFirmwareResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("updateFirmwareStatus", isValidUpdateFirmwareStatus)
}
