package main

import "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/localauth"

func (handler *ChargingStationHandler) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (response *localauth.GetLocalListVersionResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("returning current local list version: %v", handler.localAuthListVersion)
	return localauth.NewGetLocalListVersionResponse(handler.localAuthListVersion), nil
}

func (handler *ChargingStationHandler) OnSendLocalList(request *localauth.SendLocalListRequest) (response *localauth.SendLocalListResponse, err error) {
	if request.VersionNumber <= handler.localAuthListVersion {
		logDefault(request.GetFeatureName()).
			Errorf("requested listVersion %v is lower/equal than the current list version %v", request.VersionNumber, handler.localAuthListVersion)
		return localauth.NewSendLocalListResponse(localauth.SendLocalListStatusVersionMismatch), nil
	}
	if request.UpdateType == localauth.UpdateTypeFull {
		handler.localAuthList = request.LocalAuthorizationList
		handler.localAuthListVersion = request.VersionNumber
	} else if request.UpdateType == localauth.UpdateTypeDifferential {
		handler.localAuthList = append(handler.localAuthList, request.LocalAuthorizationList...)
		handler.localAuthListVersion = request.VersionNumber
	}
	logDefault(request.GetFeatureName()).Errorf("accepted new local authorization list %v, %v",
		request.VersionNumber, request.UpdateType)
	return localauth.NewSendLocalListResponse(localauth.SendLocalListStatusAccepted), nil
}
