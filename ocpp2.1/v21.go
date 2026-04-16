// Package ocpp21 contains an implementation of the OCPP 2.1 communication protocol
// between a Charging Station and a Charging Station Management System (CSMS)
// in an EV charging infrastructure.
package ocpp21

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/data"
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
)

// Profiles
var (
	AuthorizationProfile = ocpp.NewProfile(
		"authorization",
		authorization.AuthorizeFeature{},
		authorization.ClearCacheFeature{},
	)
	AvailabilityProfile = ocpp.NewProfile(
		"availability",
		availability.ChangeAvailabilityFeature{},
		availability.HeartbeatFeature{},
		availability.StatusNotificationFeature{},
	)
	DataProfile = ocpp.NewProfile(
		"data",
		data.DataTransferFeature{},
	)
	ProvisioningProfile = ocpp.NewProfile(
		"provisioning",
		provisioning.BootNotificationFeature{},
		provisioning.ResetFeature{},
		provisioning.GetBaseReportFeature{},
		provisioning.GetReportFeature{},
		provisioning.NotifyReportFeature{},
		provisioning.GetVariablesFeature{},
		provisioning.SetVariablesFeature{},
		provisioning.SetNetworkProfileFeature{},
	)
)

// Subprotocol returns the OCPP 2.1 websocket subprotocol.
func Subprotocol() string {
	return types.V21Subprotocol
}
