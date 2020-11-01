package firmware

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

// -------------------- Publish Firmware Status Notification (CS -> CSMS) --------------------

const PublishFirmwareStatusNotificationFeatureName = "PublishFirmwareStatusNotification"

// Status reported in PublishFirmwareStatusNotificationRequest.
type PublishFirmwareStatus string

const (
	PublishFirmwareStatusIdle              PublishFirmwareStatus = "Idle"
	PublishFirmwareStatusDownloadScheduled PublishFirmwareStatus = "DownloadScheduled"
	PublishFirmwareStatusDownloading       PublishFirmwareStatus = "Downloading"
	PublishFirmwareStatusDownloaded        PublishFirmwareStatus = "Downloaded"
	PublishFirmwareStatusPublished         PublishFirmwareStatus = "Published"
	PublishFirmwareStatusDownloadFailed    PublishFirmwareStatus = "DownloadFailed"
	PublishFirmwareStatusDownloadPaused    PublishFirmwareStatus = "DownloadPaused"
	PublishFirmwareStatusInvalidChecksum   PublishFirmwareStatus = "InvalidChecksum"
	PublishFirmwareStatusChecksumVerified  PublishFirmwareStatus = "ChecksumVerified"
	PublishFirmwareStatusPublishFailed     PublishFirmwareStatus = "PublishFailed"
)

func isValidPublishFirmwareStatus(fl validator.FieldLevel) bool {
	status := PublishFirmwareStatus(fl.Field().String())
	switch status {
	case PublishFirmwareStatusIdle, PublishFirmwareStatusDownloadScheduled, PublishFirmwareStatusDownloading, PublishFirmwareStatusDownloaded, PublishFirmwareStatusPublished, PublishFirmwareStatusDownloadFailed, PublishFirmwareStatusDownloadPaused, PublishFirmwareStatusInvalidChecksum, PublishFirmwareStatusChecksumVerified, PublishFirmwareStatusPublishFailed:
		return true
	default:
		return false
	}
}

// The field definition of the PublishFirmwareStatusNotification request payload sent by the Charging Station to the CSMS.
type PublishFirmwareStatusNotificationRequest struct {
	Status PublishFirmwareStatus `json:"status" validate:"required,publishFirmwareStatus"` // This contains the progress status of the publishfirmware installation.
	//TODO: add required_if validation tag after upgrade to govalidator v10
	Location  []string `json:"location,omitempty" validate:"omitempty,dive,max=512"` // Can be multiple URI’s, if the Local Controller supports e.g. HTTP, HTTPS, and FTP.
	RequestID *int     `json:"requestId,omitempty" validate:"omitempty,gte=0"`       // The request id that was provided in the PublishFirmwareRequest which triggered this action.
}

// This field definition of the PublishFirmwareStatusNotification response payload, sent by the CSMS to the Charging Station in response to a PublishFirmwareStatusNotificationRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type PublishFirmwareStatusNotificationResponse struct {
}

// The local controller sends a PublishFirmwareStatusNotificationRequest to inform the CSMS about the current PublishFirmware status.
// If the firmware was published correctly, the request will contain the location(s) URI(s) where the firmware was published at.
//
// The CSMS responds to each request with a PublishFirmwareStatusNotificationResponse.
type PublishFirmwareStatusNotificationFeature struct{}

func (f PublishFirmwareStatusNotificationFeature) GetFeatureName() string {
	return PublishFirmwareStatusNotificationFeatureName
}

func (f PublishFirmwareStatusNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(PublishFirmwareStatusNotificationRequest{})
}

func (f PublishFirmwareStatusNotificationFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(PublishFirmwareStatusNotificationResponse{})
}

func (r PublishFirmwareStatusNotificationRequest) GetFeatureName() string {
	return PublishFirmwareStatusNotificationFeatureName
}

func (c PublishFirmwareStatusNotificationResponse) GetFeatureName() string {
	return PublishFirmwareStatusNotificationFeatureName
}

// Creates a new PublishFirmwareStatusNotificationRequest, containing all required fields. Optional fields may be set afterwards.
func NewPublishFirmwareStatusNotificationRequest(status PublishFirmwareStatus) *PublishFirmwareStatusNotificationRequest {
	return &PublishFirmwareStatusNotificationRequest{Status: status}
}

// Creates a new PublishFirmwareStatusNotificationResponse, which doesn't contain any required or optional fields.
func NewPublishFirmwareStatusNotificationResponse() *PublishFirmwareStatusNotificationResponse {
	return &PublishFirmwareStatusNotificationResponse{}
}

func init() {
	_ = types.Validate.RegisterValidation("publishFirmwareStatus", isValidPublishFirmwareStatus)
}
