package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/iso15118"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

func (handler *ChargingStationHandler) OnDeleteCertificate(request *iso15118.DeleteCertificateRequest) (response *iso15118.DeleteCertificateResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnGetInstalledCertificateIds(request *iso15118.GetInstalledCertificateIdsRequest) (response *iso15118.GetInstalledCertificateIdsResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (handler *ChargingStationHandler) OnInstallCertificate(request *iso15118.InstallCertificateRequest) (response *iso15118.InstallCertificateResponse, err error) {
	logDefault(request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}
