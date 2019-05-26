package v16

import (
	"github.com/lorenzodonini/go-ocpp/ocpp"
)

const (
	BootNotificationFeatureName = "BootNotification"
	AuthorizeFeatureName = "Authorize"
	ChangeAvailabilityFeatureName = "ChangeAvailability"
)

type coreProfile struct {
	*ocpp.Profile
}

func (profile* coreProfile)CreateBootNotification(chargePointModel string, chargePointVendor string) *BootNotificationRequest {
	return &BootNotificationRequest{ChargePointModel: chargePointModel, ChargePointVendor: chargePointVendor}
}

func (profile* coreProfile)CreateAuthorization(idTag string) *AuthorizeRequest {
	return &AuthorizeRequest{IdTag: idTag}
}

func (profile* coreProfile)CreateChangeAvailability(connectorId int, availabilityType AvailabilityType) *ChangeAvailabilityRequest {
	return &ChangeAvailabilityRequest{ConnectorId: connectorId, Type: availabilityType}
}

var CoreProfile = coreProfile{
	ocpp.NewProfile("core", BootNotificationFeature{}),
}
