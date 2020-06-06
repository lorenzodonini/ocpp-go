// The authorization functional block contains OCPP 2.0 authorization-related features. It contains different ways of authorizing a user, online and/or offline .
package authorization

import "github.com/lorenzodonini/ocpp-go/ocpp"

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.0 Authorization profile.
type CSMSHandler interface {
	// OnAuthorize is called on the CSMS whenever an AuthorizeRequest is received from a charging station.
	OnAuthorize(chargingStationID string, request *AuthorizeRequest) (confirmation *AuthorizeResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.0 Authorization profile.
type ChargingStationHandler interface {
	// OnClearCache is called on a charging station whenever a ClearCacheRequest is received from the CSMS.
	OnClearCache(request *ClearCacheRequest) (confirmation *ClearCacheResponse, err error)
}

const ProfileName = "authorization"

var Profile = ocpp.NewProfile(
	ProfileName,
	AuthorizeFeature{},
	ClearCacheFeature{},
)
