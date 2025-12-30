package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/authorization"
)

func (handler *ChargingStationHandler) OnClearCache(request *authorization.ClearCacheRequest) (response *authorization.ClearCacheResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("cleared mocked cache")
	return authorization.NewClearCacheResponse(authorization.ClearCacheStatusAccepted), nil
}
