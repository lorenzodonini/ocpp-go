package firmware

import (
	"reflect"

	"gopkg.in/go-playground/validator.v9"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

// -------------------- Publish Firmware (CSMS -> CS) --------------------

const UnpublishFirmwareFeatureName = "UnpublishFirmware"

// Status for when stopping to publish a Firmware.
type UnpublishFirmwareStatus string

const (
	UnpublishFirmwareStatusDownloadOngoing UnpublishFirmwareStatus = "DownloadOngoing" // Intermediate state. Firmware is being downloaded.
	UnpublishFirmwareStatusNoFirmware      UnpublishFirmwareStatus = "NoFirmware"      // There is no published file.
	UnpublishFirmwareStatusUnpublished     UnpublishFirmwareStatus = "Unpublished"     // Successful end state. Firmware file no longer being published.
)

func isValidUnpublishFirmwareStatus(fl validator.FieldLevel) bool {
	status := UnpublishFirmwareStatus(fl.Field().String())
	switch status {
	case UnpublishFirmwareStatusDownloadOngoing, UnpublishFirmwareStatusNoFirmware, UnpublishFirmwareStatusUnpublished:
		return true
	default:
		return false
	}
}

// The field definition of the UnpublishFirmware request payload sent by the CSMS to the Charging Station.
type UnpublishFirmwareRequest struct {
	Checksum string `json:"checksum" validate:"required,max=32"` // The MD5 checksum over the entire firmware file as a hexadecimal string of length 32.
}

// This field definition of the UnpublishFirmware response payload, sent by the Charging Station to the CSMS in response to a UnpublishFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type UnpublishFirmwareResponse struct {
	Status UnpublishFirmwareStatus `json:"status" validate:"required,unpublishFirmwareStatus"`
}

// Allows to stop a Local Controller from publishing a firmware update to connected Charging Stations.
// The CSMS sends an UnpublishFirmwareRequest to instruct the local controller to unpublish the firmware.
// The local controller unpublishes the firmware, then responds with an UnpublishFirmwareResponse.
type UnpublishFirmwareFeature struct{}

func (f UnpublishFirmwareFeature) GetFeatureName() string {
	return UnpublishFirmwareFeatureName
}

func (f UnpublishFirmwareFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(UnpublishFirmwareRequest{})
}

func (f UnpublishFirmwareFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(UnpublishFirmwareResponse{})
}

func (r UnpublishFirmwareRequest) GetFeatureName() string {
	return UnpublishFirmwareFeatureName
}

func (c UnpublishFirmwareResponse) GetFeatureName() string {
	return UnpublishFirmwareFeatureName
}

// Creates a new UnpublishFirmwareRequest, containing all required fields. There are no optional fields for this message.
func NewUnpublishFirmwareRequest(checksum string) *UnpublishFirmwareRequest {
	return &UnpublishFirmwareRequest{Checksum: checksum}
}

// Creates a new UnpublishFirmwareResponse, containing all required fields. There are no optional fields for this message.
func NewUnpublishFirmwareResponse(status UnpublishFirmwareStatus) *UnpublishFirmwareResponse {
	return &UnpublishFirmwareResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("unpublishFirmwareStatus", isValidUnpublishFirmwareStatus)
}
