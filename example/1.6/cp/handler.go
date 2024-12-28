package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/certificates"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/extendedtriggermessage"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/logging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/securefirmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/security"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// ConnectorInfo contains some simple state about a single connector.
type ConnectorInfo struct {
	status             core.ChargePointStatus
	availability       core.AvailabilityType
	currentTransaction int
	currentReservation int
}

// ChargePointHandler contains some simple state that a charge point needs to keep.
// In production this will typically be replaced by database/API calls.
type ChargePointHandler struct {
	status               core.ChargePointStatus
	connectors           map[int]*ConnectorInfo
	errorCode            core.ChargePointErrorCode
	configuration        ConfigMap
	meterValue           int
	localAuthList        []localauth.AuthorizationData
	localAuthListVersion int
}

var chargePoint ocpp16.ChargePoint

func (handler *ChargePointHandler) isValidConnectorID(ID int) bool {
	_, ok := handler.connectors[ID]
	return ok || ID == 0
}

// ------------- Core profile callbacks -------------

func (handler *ChargePointHandler) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	if _, ok := handler.connectors[request.ConnectorId]; !ok {
		logDefault(request.GetFeatureName()).Errorf("cannot change availability for invalid connector %v", request.ConnectorId)
		return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusRejected), nil
	}
	handler.connectors[request.ConnectorId].availability = request.Type
	if request.Type == core.AvailabilityTypeInoperative {
		// TODO: stop ongoing transactions
		handler.connectors[request.ConnectorId].status = core.ChargePointStatusUnavailable
	} else {
		handler.connectors[request.ConnectorId].status = core.ChargePointStatusAvailable
	}
	logDefault(request.GetFeatureName()).Infof("change availability for connector %v", request.ConnectorId)
	go updateStatus(handler, request.ConnectorId, handler.connectors[request.ConnectorId].status)
	return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	configKey, ok := handler.configuration[request.Key]
	if !ok {
		logDefault(request.GetFeatureName()).Errorf("couldn't change configuration for unsupported parameter %v", configKey.Key)
		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusNotSupported), nil
	} else if configKey.Readonly {
		logDefault(request.GetFeatureName()).Errorf("couldn't change configuration for readonly parameter %v", configKey.Key)
		return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusRejected), nil
	}
	configKey.Value = &request.Value
	handler.configuration[request.Key] = configKey
	logDefault(request.GetFeatureName()).Infof("changed configuration for parameter %v to %v", configKey.Key, configKey.Value)
	return core.NewChangeConfigurationConfirmation(core.ConfigurationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("cleared mocked cache")
	return core.NewClearCacheConfirmation(core.ClearCacheStatusAccepted), nil
}

func (handler *ChargePointHandler) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("data transfer [Vendor: %v Message: %v]: %v", request.VendorId, request.MessageId, request.Data)
	return core.NewDataTransferConfirmation(core.DataTransferStatusAccepted), nil
}

func (handler *ChargePointHandler) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	var resultKeys []core.ConfigurationKey
	var unknownKeys []string
	for _, key := range request.Key {
		configKey, ok := handler.configuration[key]
		if !ok {
			unknownKeys = append(unknownKeys, key)
		} else {
			resultKeys = append(resultKeys, configKey)
		}
	}
	if len(request.Key) == 0 {
		// Return config for all keys
		for _, v := range handler.configuration {
			resultKeys = append(resultKeys, v)
		}
	}
	logDefault(request.GetFeatureName()).Infof("returning configuration for requested keys: %v", request.Key)
	conf := core.NewGetConfigurationConfirmation(resultKeys)
	conf.UnknownKey = unknownKeys
	return conf, nil
}

func (handler *ChargePointHandler) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	if request.ConnectorId != nil {
		connector, ok := handler.connectors[*request.ConnectorId]
		if !ok {
			return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
		} else if connector.availability != core.AvailabilityTypeOperative || connector.status != core.ChargePointStatusAvailable || connector.currentTransaction > 0 {
			return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
		}
		logDefault(request.GetFeatureName()).Infof("started transaction %v on connector %v", connector.currentTransaction, request.ConnectorId)
		connector.currentTransaction = *request.ConnectorId
		return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusAccepted), nil
	}
	logDefault(request.GetFeatureName()).Errorf("couldn't start a transaction for %v without a connectorID", request.IdTag)
	return core.NewRemoteStartTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
}

func (handler *ChargePointHandler) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	for key, val := range handler.connectors {
		if val.currentTransaction == request.TransactionId {
			logDefault(request.GetFeatureName()).Infof("stopped transaction %v on connector %v", val.currentTransaction, key)
			val.currentTransaction = 0
			val.currentReservation = 0
			val.status = core.ChargePointStatusAvailable
			return core.NewRemoteStopTransactionConfirmation(types.RemoteStartStopStatusAccepted), nil
		}
	}
	logDefault(request.GetFeatureName()).Errorf("couldn't stop transaction %v, no such transaction is ongoing", request.TransactionId)
	return core.NewRemoteStopTransactionConfirmation(types.RemoteStartStopStatusRejected), nil
}

func (handler *ChargePointHandler) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	// TODO: stop all ongoing transactions
	logDefault(request.GetFeatureName()).Warn("no reset logic implemented yet")
	return core.NewResetConfirmation(core.ResetStatusAccepted), nil
}

func (handler *ChargePointHandler) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	_, ok := handler.connectors[request.ConnectorId]
	if !ok {
		logDefault(request.GetFeatureName()).Errorf("couldn't unlock invalid connector %v", request.ConnectorId)
		return core.NewUnlockConnectorConfirmation(core.UnlockStatusNotSupported), nil
	}
	logDefault(request.GetFeatureName()).Infof("unlocked connector %v", request.ConnectorId)
	return core.NewUnlockConnectorConfirmation(core.UnlockStatusUnlocked), nil
}

// ------------- Local authorization list profile callbacks -------------

func (handler *ChargePointHandler) OnGetLocalListVersion(request *localauth.GetLocalListVersionRequest) (confirmation *localauth.GetLocalListVersionConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("returning current local list version: %v", handler.localAuthListVersion)
	return localauth.NewGetLocalListVersionConfirmation(handler.localAuthListVersion), nil
}

func (handler *ChargePointHandler) OnSendLocalList(request *localauth.SendLocalListRequest) (confirmation *localauth.SendLocalListConfirmation, err error) {
	if request.ListVersion <= handler.localAuthListVersion {
		logDefault(request.GetFeatureName()).Errorf("requested listVersion %v is lower/equal than the current list version %v", request.ListVersion, handler.localAuthListVersion)
		return localauth.NewSendLocalListConfirmation(localauth.UpdateStatusVersionMismatch), nil
	}
	if request.UpdateType == localauth.UpdateTypeFull {
		handler.localAuthList = request.LocalAuthorizationList
		handler.localAuthListVersion = request.ListVersion
	} else if request.UpdateType == localauth.UpdateTypeDifferential {
		handler.localAuthList = append(handler.localAuthList, request.LocalAuthorizationList...)
		handler.localAuthListVersion = request.ListVersion
	}
	logDefault(request.GetFeatureName()).Errorf("accepted new local authorization list %v, %v", request.ListVersion, request.UpdateType)
	return localauth.NewSendLocalListConfirmation(localauth.UpdateStatusAccepted), nil
}

// ------------- Firmware management profile callbacks -------------

func (handler *ChargePointHandler) OnGetDiagnostics(request *firmware.GetDiagnosticsRequest) (confirmation *firmware.GetDiagnosticsConfirmation, err error) {
	// TODO: perform diagnostics upload out-of-band
	logDefault(request.GetFeatureName()).Warn("no diagnostics upload logic implemented yet")
	return firmware.NewGetDiagnosticsConfirmation(), nil
}

func (handler *ChargePointHandler) OnUpdateFirmware(request *firmware.UpdateFirmwareRequest) (confirmation *firmware.UpdateFirmwareConfirmation, err error) {
	retries := 0
	retryInterval := 30
	if request.Retries != nil {
		retries = *request.Retries
	}
	if request.RetryInterval != nil {
		retryInterval = *request.RetryInterval
	}
	logDefault(request.GetFeatureName()).Infof("starting update firmware procedure")
	go updateFirmware(request.Location, request.RetrieveDate, retries, retryInterval)
	return firmware.NewUpdateFirmwareConfirmation(), nil
}

// ------------- Remote trigger profile callbacks -------------

func (handler *ChargePointHandler) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	logDefault(request.GetFeatureName()).Infof("received trigger for %v", request.RequestedMessage)
	status := remotetrigger.TriggerMessageStatusRejected
	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:
		// TODO: schedule boot notification message
		break
	case firmware.DiagnosticsStatusNotificationFeatureName:
		// Schedule diagnostics status notification request
		go func() {
			_, e := chargePoint.DiagnosticsStatusNotification(firmware.DiagnosticsStatusIdle)
			checkError(e)
			logDefault(firmware.DiagnosticsStatusNotificationFeatureName).Info("diagnostics status notified")
		}()
		status = remotetrigger.TriggerMessageStatusAccepted
	case firmware.FirmwareStatusNotificationFeatureName:
		// TODO: schedule firmware status notification message
		break
	case core.HeartbeatFeatureName:
		// Schedule heartbeat request
		go func() {
			conf, e := chargePoint.Heartbeat()
			checkError(e)
			logDefault(core.HeartbeatFeatureName).Infof("clock synchronized: %v", conf.CurrentTime.FormatTimestamp())
		}()
		status = remotetrigger.TriggerMessageStatusAccepted
	case core.MeterValuesFeatureName:
		// TODO: schedule meter values message
		break
	case core.StatusNotificationFeatureName:
		connectorID := *request.ConnectorId
		// Check if requested connector is valid and status can be retrieved
		if !handler.isValidConnectorID(connectorID) {
			logDefault(request.GetFeatureName()).Errorf("cannot trigger %v: requested invalid connector %v", request.RequestedMessage, connectorID)
			return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusRejected), nil
		}
		// Schedule status notification request
		go func() {
			status := handler.status
			if c, ok := handler.connectors[connectorID]; ok {
				status = c.status
			}
			statusConfirmation, err := chargePoint.StatusNotification(connectorID, handler.errorCode, status)
			checkError(err)
			logDefault(statusConfirmation.GetFeatureName()).Infof("status for connector %v sent: %v", connectorID, status)
		}()
		status = remotetrigger.TriggerMessageStatusAccepted
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}
	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}

// ------------- Reservation profile callbacks -------------

func (handler *ChargePointHandler) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	connector := handler.connectors[request.ConnectorId]
	if connector == nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	} else if connector.status != core.ChargePointStatusAvailable {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	}
	connector.currentReservation = request.ReservationId
	logDefault(request.GetFeatureName()).Infof("reservation %v for connector %v accepted", request.ReservationId, request.ConnectorId)
	go updateStatus(handler, request.ConnectorId, core.ChargePointStatusReserved)
	// TODO: automatically remove reservation after expiryDate
	return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), nil
}

func (handler *ChargePointHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	for k, v := range handler.connectors {
		if v.currentReservation == request.ReservationId {
			v.currentReservation = 0
			if v.status == core.ChargePointStatusReserved {
				go updateStatus(handler, k, core.ChargePointStatusAvailable)
			}
			logDefault(request.GetFeatureName()).Infof("reservation %v for connector %v canceled", request.ReservationId, k)
			return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusAccepted), nil
		}
	}
	logDefault(request.GetFeatureName()).Infof("couldn't cancel reservation %v: reservation not found!", request.ReservationId)
	return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
}

// ------------- Smart charging profile callbacks -------------

func (handler *ChargePointHandler) OnSetChargingProfile(request *smartcharging.SetChargingProfileRequest) (confirmation *smartcharging.SetChargingProfileConfirmation, err error) {
	// TODO: handle logic
	logDefault(request.GetFeatureName()).Warn("no set charging profile logic implemented yet")
	return smartcharging.NewSetChargingProfileConfirmation(smartcharging.ChargingProfileStatusNotSupported), nil
}

func (handler *ChargePointHandler) OnClearChargingProfile(request *smartcharging.ClearChargingProfileRequest) (confirmation *smartcharging.ClearChargingProfileConfirmation, err error) {
	// TODO: handle logic
	logDefault(request.GetFeatureName()).Warn("no clear charging profile logic implemented yet")
	return smartcharging.NewClearChargingProfileConfirmation(smartcharging.ClearChargingProfileStatusUnknown), nil
}

func (handler *ChargePointHandler) OnGetCompositeSchedule(request *smartcharging.GetCompositeScheduleRequest) (confirmation *smartcharging.GetCompositeScheduleConfirmation, err error) {
	// TODO: handle logic
	logDefault(request.GetFeatureName()).Warn("no get composite schedule logic implemented yet")
	return smartcharging.NewGetCompositeScheduleConfirmation(smartcharging.GetCompositeScheduleStatusRejected), nil
}

func (handler *ChargePointHandler) OnDeleteCertificate(request *certificates.DeleteCertificateRequest) (response *certificates.DeleteCertificateResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("certificate %v deleted", request.CertificateHashData)
	return certificates.NewDeleteCertificateResponse(certificates.DeleteCertificateStatusAccepted), nil
}

func (handler *ChargePointHandler) OnGetInstalledCertificateIds(request *certificates.GetInstalledCertificateIdsRequest) (response *certificates.GetInstalledCertificateIdsResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("returning installed certificate ids")
	return certificates.NewGetInstalledCertificateIdsResponse(certificates.GetInstalledCertificateStatusAccepted), nil
}

func (handler *ChargePointHandler) OnInstallCertificate(request *certificates.InstallCertificateRequest) (response *certificates.InstallCertificateResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("certificate installed")
	return certificates.NewInstallCertificateResponse(certificates.CertificateStatusAccepted), nil
}

func (handler *ChargePointHandler) OnGetLog(request *logging.GetLogRequest) (response *logging.GetLogResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("returning log")
	return logging.NewGetLogResponse(logging.LogStatusAccepted), nil
}

func (handler *ChargePointHandler) OnSignedUpdateFirmware(request *securefirmware.SignedUpdateFirmwareRequest) (response *securefirmware.SignedUpdateFirmwareResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("signed update firmware request received")
	return securefirmware.NewSignedUpdateFirmwareResponse(securefirmware.UpdateFirmwareStatusAccepted), nil
}

func (handler *ChargePointHandler) OnExtendedTriggerMessage(request *extendedtriggermessage.ExtendedTriggerMessageRequest) (response *extendedtriggermessage.ExtendedTriggerMessageResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("extended trigger message received")
	return extendedtriggermessage.NewExtendedTriggerMessageResponse(extendedtriggermessage.ExtendedTriggerMessageStatusAccepted), nil
}

func (handler *ChargePointHandler) OnCertificateSigned(request *security.CertificateSignedRequest) (response *security.CertificateSignedResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("certificate signed")
	return security.NewCertificateSignedResponse(security.CertificateSignedStatusAccepted), nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getExpiryDate(info *types.IdTagInfo) string {
	if info.ExpiryDate != nil {
		return fmt.Sprintf("authorized until %v", info.ExpiryDate.String())
	}
	return ""
}

func updateStatus(stateHandler *ChargePointHandler, connector int, status core.ChargePointStatus, props ...func(request *core.StatusNotificationRequest)) {
	if connector == 0 {
		stateHandler.status = status
	} else {
		stateHandler.connectors[connector].status = status
	}
	statusConfirmation, err := chargePoint.StatusNotification(connector, stateHandler.errorCode, status, props...)
	checkError(err)
	if connector == 0 {
		logDefault(statusConfirmation.GetFeatureName()).Infof("status for all connectors updated to %v", status)
	} else {
		logDefault(statusConfirmation.GetFeatureName()).Infof("status for connector %v updated to %v", connector, status)
	}
}

func updateFirmwareStatus(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) {
	statusConfirmation, err := chargePoint.FirmwareStatusNotification(status, props...)
	checkError(err)
	logDefault(statusConfirmation.GetFeatureName()).Infof("firmware status updated to %v", status)
}

func updateFirmware(location string, retrieveDate *types.DateTime, retries int, retryInterval int) {
	updateFirmwareStatus(firmware.FirmwareStatusDownloading)
	err := downloadFile("/tmp/out.bin", location)
	if err != nil {
		logDefault(firmware.UpdateFirmwareFeatureName).Errorf("error while downloading file %v", err)
		updateFirmwareStatus(firmware.FirmwareStatusDownloadFailed)
		return
	}
	updateFirmwareStatus(firmware.FirmwareStatusDownloaded)
	// Simulate installation
	updateFirmwareStatus(firmware.FirmwareStatusInstalling)
	time.Sleep(time.Second * 5)
	// Notify completion
	updateFirmwareStatus(firmware.FirmwareStatusInstalled)
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
