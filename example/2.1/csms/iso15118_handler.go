package main

import (
	"github.com/xBlaz3kx/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ocpp-go/ocpp2.0.1/iso15118"
	"github.com/xBlaz3kx/ocpp-go/ocppj"
)

func (c *CSMSHandler) OnGet15118EVCertificate(chargingStationID string, request *iso15118.Get15118EVCertificateRequest) (response *iso15118.Get15118EVCertificateResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}

func (c *CSMSHandler) OnGetCertificateStatus(chargingStationID string, request *iso15118.GetCertificateStatusRequest) (response *iso15118.GetCertificateStatusResponse, err error) {
	logDefault(chargingStationID, request.GetFeatureName()).Warnf("Unsupported feature")
	return nil, ocpp.NewHandlerError(ocppj.NotSupported, "Not supported")
}
