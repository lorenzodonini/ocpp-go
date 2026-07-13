package main

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/authorization"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (c *CSMSHandler) OnAuthorize(chargingStationID string, request *authorization.AuthorizeRequest) (response *authorization.AuthorizeResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Infof("client with token %v authorized", request.IdToken)
	response = authorization.NewAuthorizationResponse(*types.NewIdTokenInfo(types.AuthorizationStatusAccepted))
	return
}
