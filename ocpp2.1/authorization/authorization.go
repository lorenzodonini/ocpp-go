// Package authorization contains the authorization functional block for OCPP 2.1.
// It contains different ways of authorizing a user, online and/or offline.
package authorization

import "github.com/lorenzodonini/ocpp-go/ocpp"

// ProfileName is the name of the authorization functional block.
const ProfileName = "authorization"

// ChargingStationHandler contains callbacks for handling incoming CSMS requests.
type ChargingStationHandler interface {
	// OnClearCache handles a ClearCacheRequest from the CSMS.
	OnClearCache(request *ClearCacheRequest) (response *ClearCacheResponse, err error)
}

// CSMSHandler contains callbacks for handling incoming Charging Station requests.
type CSMSHandler interface {
	// OnAuthorize handles an AuthorizeRequest from a Charging Station.
	OnAuthorize(chargingStationID string, request *AuthorizeRequest) (response *AuthorizeResponse, err error)
}

// Profile returns the authorization profile with all features.
func Profile() *ocpp.Profile {
	return ocpp.NewProfile(
		ProfileName,
		AuthorizeFeature{},
		ClearCacheFeature{},
	)
}
