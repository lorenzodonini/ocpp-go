package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"reflect"
)

// -------------------- Publish Firmware (CSMS -> CS) --------------------

const PublishFirmwareFeatureName = "PublishFirmware"

// The field definition of the PublishFirmware request payload sent by the CSMS to the Charging Station.
type PublishFirmwareRequest struct {
	Location      string `json:"location" validate:"required,max=512"`               // This contains a string containing a URI pointing to a location from which to retrieve the firmware.
	Retries       *int   `json:"retries,omitempty" validate:"omitempty,gte=0"`       // This specifies how many times Charging Station must try to download the firmware before giving up. If this field is not present, it is left to Charging Station to decide how many times it wants to retry.
	Checksum      string `json:"checksum" validate:"required,max=32"`                // The MD5 checksum over the entire firmware file as a hexadecimal string of length 32.
	RequestID     int    `json:"requestId" validate:"gte=0"`                         // The Id of the request.
	RetryInterval *int   `json:"retryInterval,omitempty" validate:"omitempty,gte=0"` // The interval in seconds after which a retry may be attempted. If this field is not present, it is left to Charging Station to decide how long to wait between attempts.
}

// This field definition of the PublishFirmware response payload, sent by the Charging Station to the CSMS in response to a PublishFirmwareRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type PublishFirmwareResponse struct {
	Status     types.GenericStatus `json:"status" validate:"required,genericStatus"`
	StatusInfo *types.StatusInfo   `json:"statusInfo,omitempty" validate:"omitempty"`
}

// The CSMS sends a PublishFirmwareRequest to instruct the Local Controller to download and publish the firmware,
// including an MD5 checksum of the firmware file.
// Upon receipt of PublishFirmwareRequest, the Local Controller responds with PublishFirmwareResponse.
//
// The local controller will download the firmware out-of-band and publish the URI of the updated firmware to
// the CSMS via a PublishFirmwareStatusNotificationRequest.
//
// Whenever the CSMS instructs charging stations to update their firmware, it will instruct to download the
// firmware form the local controller instead of from the CSMS, saving data and bandwidth on the WAN interface.
type PublishFirmwareFeature struct{}

func (f PublishFirmwareFeature) GetFeatureName() string {
	return PublishFirmwareFeatureName
}

func (f PublishFirmwareFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(PublishFirmwareRequest{})
}

func (f PublishFirmwareFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(PublishFirmwareResponse{})
}

func (r PublishFirmwareRequest) GetFeatureName() string {
	return PublishFirmwareFeatureName
}

func (c PublishFirmwareResponse) GetFeatureName() string {
	return PublishFirmwareFeatureName
}

// Creates a new PublishFirmwareRequest,  containing all required fields. Optional fields may be set afterwards.
func NewPublishFirmwareRequest(location string, checksum string, requestID int) *PublishFirmwareRequest {
	return &PublishFirmwareRequest{Location: location, Checksum: checksum, RequestID: requestID}
}

// Creates a new PublishFirmwareResponse, containing all required fields. Optional fields may be set afterwards.
func NewPublishFirmwareResponse(status types.GenericStatus) *PublishFirmwareResponse {
	return &PublishFirmwareResponse{Status: status}
}
