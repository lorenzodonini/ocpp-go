package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"reflect"
	"time"
)

const (
	BootNotificationFeatureName = "BootNotification"
	AuthorizeFeatureName = "Authorize"
)

type coreProfile struct {
	*ocpp.Profile
}

func (profile* coreProfile)CreateBootNotification(chargePointModel string, chargePointVendor string) *BootNotificationRequest {
	return &BootNotificationRequest{ChargePointModel: chargePointModel, ChargePointVendor: chargePointVendor}
}

var CoreProfile = coreProfile{
	ocpp.NewProfile("core", BootNotificationFeature{}),
}

// -------------------- Boot Notification --------------------
type BootNotificationRequest struct {
	ocpp.Request					`json:"-"`
	ChargeBoxSerialNumber string 	`json:"chargeBoxSerialNumber,omitempty" validate:"max=25"`
	ChargePointModel string			`json:"chargePointModel" validate:"required,max=20"`
	ChargePointSerialNumber string	`json:"chargePointSerialNumber,omitempty" validate:"max=25"`
	ChargePointVendor string		`json:"chargePointVendor" validate:"required,max=20"`
	FirmwareVersion string			`json:"firmwareVersion,omitempty" validate:"max=50"`
	Iccid string					`json:"iccid,omitempty" validate:"max=20"`
	Imsi string						`json:"imsi,omitempty" validate:"max=20"`
	MeterSerialNumber string		`json:"meterSerialNumber,omitempty" validate:"max=25"`
	MeterType string				`json:"meterType,omitempty" validate:"max=25"`
}

//TODO: add custom validator for registration status & interval
type BootNotificationConfirmation struct {
	ocpp.Confirmation				`json:"-"`
	CurrentTime time.Time			`json:"currentTime" validate:"required"`
	Interval int					`json:"interval" validate:"required,gte=0"`
	Status ocpp.RegistrationStatus	`json:"status" validate:"required"`
}

type BootNotificationFeature struct {}

func (f BootNotificationFeature) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (f BootNotificationFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(BootNotificationRequest{})
}

func (f BootNotificationFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(BootNotificationConfirmation{})
}

func (r BootNotificationRequest) GetFeatureName() string {
	return BootNotificationFeatureName
}

func (c BootNotificationConfirmation) GetFeatureName() string {
	return BootNotificationFeatureName
}

// -------------------- Authorize --------------------
type AuthorizeRequest struct {
	IdTag string				`json:"idTag" validate:"required,max=20"`
}

type AuthorizeConfirmation struct {
	IdTagInfo ocpp.IdTagInfo	`json:"idTagInfo" validate:"required"`
}

type AuthorizeFeature struct {}

func (f AuthorizeFeature) GetFeatureName() string {
	return AuthorizeFeatureName
}

func (f AuthorizeFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(AuthorizeRequest{})
}

func (f AuthorizeFeature) GetConfirmationType() reflect.Type {
	return reflect.TypeOf(AuthorizeConfirmation{})
}

func (r AuthorizeRequest) GetFeatureName() string {
	return AuthorizeFeatureName
}

func (c AuthorizeConfirmation) GetFeatureName() string {
	return AuthorizeFeatureName
}
