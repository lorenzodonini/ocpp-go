package ocpp16_test

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/extendedtriggermessage"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/securefirmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/security"
	"github.com/stretchr/testify/mock"
)

// ---------------------- MOCK CP SECURITY HANDLER ----------------------
type MockChargePointSecurityHandler struct {
	mock.Mock
}

func (m *MockChargePointSecurityHandler) OnCertificateSigned(request *security.CertificateSignedRequest) (response *security.CertificateSignedResponse, err error) {
	args := m.MethodCalled("OnCertificateSigned", request)
	conf := args.Get(0).(*security.CertificateSignedResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS SECURITY HANDLER ----------------------
type MockCentralSystemSecurityListener struct {
	mock.Mock
}

func (m *MockCentralSystemSecurityListener) OnSecurityEventNotification(chargingStationID string, request *security.SecurityEventNotificationRequest) (response *security.SecurityEventNotificationResponse, err error) {
	args := m.MethodCalled("OnSecurityEventNotification", request)
	conf := args.Get(0).(*security.SecurityEventNotificationResponse)
	return conf, args.Error(1)
}

func (m *MockCentralSystemSecurityListener) OnSignCertificate(chargingStationID string, request *security.SignCertificateRequest) (response *security.SignCertificateResponse, err error) {
	args := m.MethodCalled("OnSignCertificate", request)
	conf := args.Get(0).(*security.SignCertificateResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS CERTIFICATE HANDLER ----------------------

type MockCentralSystemCertificateListener struct {
	mock.Mock
}

// ---------------------- MOCK CP CERTIFICATE HANDLER ----------------------

type MockChargePointCertificateHandler struct {
	mock.Mock
}

func (m *MockChargePointCertificateHandler) OnDeleteCertificate(request *certificates.DeleteCertificateRequest) (response *certificates.DeleteCertificateResponse, err error) {
	args := m.MethodCalled("OnDeleteCertificate", request)
	conf := args.Get(0).(*certificates.DeleteCertificateResponse)
	return conf, args.Error(1)
}

func (m *MockChargePointCertificateHandler) OnGetInstalledCertificateIds(request *certificates.GetInstalledCertificateIdsRequest) (response *certificates.GetInstalledCertificateIdsResponse, err error) {
	args := m.MethodCalled("OnGetInstalledCertificateIds", request)
	conf := args.Get(0).(*certificates.GetInstalledCertificateIdsResponse)
	return conf, args.Error(1)
}

func (m *MockChargePointCertificateHandler) OnInstallCertificate(request *certificates.InstallCertificateRequest) (response *certificates.InstallCertificateResponse, err error) {
	args := m.MethodCalled("OnInstallCertificate", request)
	conf := args.Get(0).(*certificates.InstallCertificateResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS EXTENDED TRIGGER MESSAGE HANDLER ----------------------

type MockCentralSystemExtendedTriggerMessageListener struct {
	mock.Mock
}

// ---------------------- MOCK CP EXTENDED TRIGGER MESSAGE HANDLER ----------------------

type MockChargePointExtendedTriggerMessageHandler struct {
	mock.Mock
}

func (m *MockChargePointExtendedTriggerMessageHandler) OnExtendedTriggerMessage(request *extendedtriggermessage.ExtendedTriggerMessageRequest) (response *extendedtriggermessage.ExtendedTriggerMessageResponse, err error) {
	args := m.MethodCalled("OnExtendedTriggerMessage", request)
	conf := args.Get(0).(*extendedtriggermessage.ExtendedTriggerMessageResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CS SECURE FIRMWARE UPDATE HANDLER ----------------------

type MockCentralSystemSecureFirmwareUpdateListener struct {
	mock.Mock
}

func (m *MockCentralSystemSecureFirmwareUpdateListener) OnSignedFirmwareStatusNotification(chargingStationID string, request *securefirmware.SignedFirmwareStatusNotificationRequest) (response *securefirmware.SignedFirmwareStatusNotificationResponse, err error) {
	args := m.MethodCalled("OnSignedFirmwareStatusNotification", request)
	conf := args.Get(0).(*securefirmware.SignedFirmwareStatusNotificationResponse)
	return conf, args.Error(1)
}

// ---------------------- MOCK CP SECURE FIRMWARE UPDATE  HANDLER ----------------------

type MockChargePointSecureFirmwareUpdateHandler struct {
	mock.Mock
}

func (m *MockChargePointSecureFirmwareUpdateHandler) OnSignedUpdateFirmware(request *securefirmware.SignedUpdateFirmwareRequest) (response *securefirmware.SignedUpdateFirmwareResponse, err error) {
	args := m.MethodCalled("OnSignedUpdateFirmware", request)
	conf := args.Get(0).(*securefirmware.SignedUpdateFirmwareResponse)
	return conf, args.Error(1)
}
