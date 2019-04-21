package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"reflect"
	"time"
)

// -------------------- Boot Notification --------------------
type BootNotificationRequest struct {
	ChargeBotSerialNumber string 	`json:"chargeBoxSerialNumber,omitempty" valid:"stringlength(0|25)"`
	ChargePointModel string			`json:"chargePointModel" valid:"stringlength(1|20)"`
	ChargePointSerialNumber string	`json:"chargePointSerialNumber,omitempty" valid:"stringlength(0|25)"`
	ChargePointVendor string		`json:"chargePointVendor" valid:"stringlength(1|20)"`
	FirmwareVersion string			`json:"chargePointVendor,omitempty" valid:"stringlength(0|50)"`
	Iccid string					`json:"iccid,omitempty" valid:"stringlength(0|20)"`
	Imsi string						`json:"imsi,omitempty" valid:"stringlength(0|20)"`
	MeterSerialNumber string		`json:"meterSerialNumber,omitempty" valid:"stringlength(0|25)"`
	MeterType string				`json:"meterType,omitempty" valid:"stringlength(0|25)"`
}

type BootNotificationConfirmation struct {
	CurrentTime time.Time			`json:"currentTime" valid:"time"`
	Interval int					`json:"interval" valid:"numeric"`
	Status ocpp.RegistrationStatus	`json:"status" valid:"registration"`
}

type BootNotificationFeature struct {
	ocpp.Feature
}

func (f* BootNotificationFeature) GetFeatureName() string {
	return "BootNotification"
}

func (f* BootNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(BootNotificationRequest{})
}

func (f * BootNotificationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(BootNotificationConfirmation{})
}