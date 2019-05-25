package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
	"github.com/lorenzodonini/go-ocpp/ocpp/1.6/core"
)

const (
	BootNotificationFeatureName = "BootNotification"
	AuthorizeFeatureName = "Authorize"
)

type coreProfile struct {
	*ocpp.Profile
}

func (profile* coreProfile)CreateBootNotification(chargePointModel string, chargePointVendor string) *core.BootNotificationRequest {
	return &core.BootNotificationRequest{ChargePointModel: chargePointModel, ChargePointVendor: chargePointVendor}
}

func (profile* coreProfile)CreateAuthorization(idTag string) *core.AuthorizeRequest {
	return &core.AuthorizeRequest{IdTag: idTag}
}

var CoreProfile = coreProfile{
	ocpp.NewProfile("core", core.BootNotificationFeature{}),
}
