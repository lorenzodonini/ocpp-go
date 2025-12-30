// The Local authorization list functional block contains OCPP 2.1 features for synchronizing local authorization lists between CSMS and charging station.
// Local lists are used for offline and generally optimized authorization.
package localauth

import "github.com/xBlaz3kx/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 Local Authorization List profile.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.1 Local Authorization List profile.
type ChargingStationHandler interface {
	// OnGetLocalListVersion is called on a charging station whenever a GetLocalListVersionRequest is received from the CSMS.
	OnGetLocalListVersion(request *GetLocalListVersionRequest) (response *GetLocalListVersionResponse, err error)
	// OnSendLocalList is called on a charging station whenever a SendLocalListRequest is received from the CSMS.
	OnSendLocalList(request *SendLocalListRequest) (response *SendLocalListResponse, err error)
}

const ProfileName = "LocalAuthList"

var Profile = ocpp.NewProfile(
	ProfileName,
	GetLocalListVersionFeature{},
	SendLocalListFeature{},
)
