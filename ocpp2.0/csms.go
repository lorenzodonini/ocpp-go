package ocpp2

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

type csms struct {
	server       *ocppj.CentralSystem
	coreListener CSMSHandler
	//localAuthListListener CentralSystemLocalAuthListListener
	//firmwareListener      CentralSystemFirmwareManagementListener
	//reservationListener   CentralSystemReservationListener
	//remoteTriggerListener CentralSystemRemoteTriggerListener
	//smartChargingListener CentralSystemSmartChargingListener
	callbacks map[string]func(confirmation ocpp.Confirmation, err error)
}

// Cancels a previously reserved charge point or connector, given the reservation Id.
func (cs *csms) CancelReservation(clientId string, callback func(*CancelReservationConfirmation, error), reservationId int, props ...func(request *CancelReservationRequest)) error {
	request := NewCancelReservationRequest(reservationId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CancelReservationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Sends a new certificate (chain) to the charging station.
func (cs *csms) CertificateSigned(clientId string, callback func(*CertificateSignedConfirmation, error), certificate []string, props ...func(*CertificateSignedRequest)) error {
	request := NewCertificateSignedRequest(certificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CertificateSignedConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Instructs a charge point to change its availability. The target availability can be set for a single evse of for the whole charging station.
func (cs *csms) ChangeAvailability(clientId string, callback func(confirmation *ChangeAvailabilityConfirmation, err error), evseID int, operationalStatus OperationalStatus, props ...func(request *ChangeAvailabilityRequest)) error {
	request := NewChangeAvailabilityRequest(evseID, operationalStatus)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ChangeAvailabilityConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

//
//// Changes the configuration of a charge point, by setting a specific key-value pair.
//// The configuration key must be supported by the target charge point, in order for the configuration to be accepted.
//func (cs *server) ChangeConfiguration(clientId string, callback func(confirmation *ChangeConfigurationConfirmation, err error), key string, value string, props ...func(request *ChangeConfigurationRequest)) error {
//	request := NewChangeConfigurationRequest(key, value)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*ChangeConfigurationConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}

func (cs *csms) ClearCache(clientId string, callback func(confirmation *ClearCacheConfirmation, err error), props ...func(*ClearCacheRequest)) error {
	request := NewClearCacheRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearCacheConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Removes one or more charging profiles from a charging station.
func (cs *csms) ClearChargingProfile(clientId string, callback func(*ClearChargingProfileConfirmation, error), props ...func(request *ClearChargingProfileRequest)) error {
	request := NewClearChargingProfileRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearChargingProfileConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearDisplay(clientId string, callback func(*ClearDisplayConfirmation, error), id int, props ...func(*ClearDisplayRequest)) error {
	request := NewClearDisplayRequest(id)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearDisplayConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) ClearVariableMonitoring(clientId string, callback func(*ClearVariableMonitoringConfirmation, error), id []int, props ...func(*ClearVariableMonitoringRequest)) error {
	request := NewClearVariableMonitoringRequest(id)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*ClearVariableMonitoringConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) CostUpdated(clientId string, callback func(*CostUpdatedConfirmation, error), totalCost float64, transactionId string, props ...func(*CostUpdatedRequest)) error {
	request := NewCostUpdatedRequest(totalCost, transactionId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CostUpdatedConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) CustomerInformation(clientId string, callback func(*CustomerInformationConfirmation, error), requestId int, report bool, clear bool, props ...func(*CustomerInformationRequest)) error {
	request := NewCustomerInformationRequest(requestId, report, clear)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*CustomerInformationConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

// Starts a custom data transfer request. Every vendor may implement their own proprietary logic for this message.
func (cs *csms) DataTransfer(clientId string, callback func(confirmation *DataTransferConfirmation, err error), vendorId string, props ...func(request *DataTransferRequest)) error {
	request := NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*DataTransferConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) DeleteCertificate(clientId string, callback func(*DeleteCertificateConfirmation, error), data CertificateHashData, props ...func(*DeleteCertificateRequest)) error {
	request := NewDeleteCertificateRequest(data)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*DeleteCertificateConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetBaseReport(clientId string, callback func(*GetBaseReportConfirmation, error), requestId int, reportBase ReportBaseType, props ...func(*GetBaseReportRequest)) error {
	request := NewGetBaseReportRequest(requestId, reportBase)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetBaseReportConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetChargingProfiles(clientId string, callback func(*GetChargingProfilesConfirmation, error), chargingProfile ChargingProfileCriterion, props ...func(*GetChargingProfilesRequest)) error {
	request := NewGetChargingProfilesRequest(chargingProfile)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetChargingProfilesConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleConfirmation, error), duration int, evseId int, props ...func(*GetCompositeScheduleRequest)) error {
	request := NewGetCompositeScheduleRequest(duration, evseId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetCompositeScheduleConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetDisplayMessages(clientId string, callback func(*GetDisplayMessagesConfirmation, error), requestId int, props ...func(*GetDisplayMessagesRequest)) error {
	request := NewGetDisplayMessagesRequest(requestId)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetDisplayMessagesConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetInstalledCertificateIds(clientId string, callback func(*GetInstalledCertificateIdsConfirmation, error), typeOfCertificate CertificateUse, props ...func(*GetInstalledCertificateIdsRequest)) error {
	request := NewGetInstalledCertificateIdsRequest(typeOfCertificate)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetInstalledCertificateIdsConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(*GetLocalListVersionRequest)) error {
	request := NewGetLocalListVersionRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetLocalListVersionConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetLog(clientId string, callback func(*GetLogConfirmation, error), logType LogType, requestID int, logParameters LogParameters, props ...func(*GetLogRequest)) error {
	request := NewGetLogRequest(logType, requestID, logParameters)
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetLogConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

func (cs *csms) GetMonitoringReport(clientId string, callback func(*GetMonitoringReportConfirmation, error), props ...func(*GetMonitoringReportRequest)) error {
	request := NewGetMonitoringReportRequest()
	for _, fn := range props {
		fn(request)
	}
	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
		if confirmation != nil {
			callback(confirmation.(*GetMonitoringReportConfirmation), protoError)
		} else {
			callback(nil, protoError)
		}
	}
	return cs.SendRequestAsync(clientId, request, genericCallback)
}

//
//// Retrieves the configuration values for the provided configuration keys.
//func (cs *server) GetConfiguration(clientId string, callback func(confirmation *GetConfigurationConfirmation, err error), keys []string, props ...func(request *GetConfigurationRequest)) error {
//	request := NewGetConfigurationRequest(keys)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetConfigurationConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs a charge point to start a transaction for a specified client on a provided connector.
//// Depending on the configuration, an explicit authorization message may still be required, before the transaction can start.
//func (cs *server) RemoteStartTransaction(clientId string, callback func(*RemoteStartTransactionConfirmation, error), idTag string, props ...func(*RemoteStartTransactionRequest)) error {
//	request := NewRemoteStartTransactionRequest(idTag)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*RemoteStartTransactionConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs a charge point to stop an ongoing transaction, given the transaction's ID.
//func (cs *server) RemoteStopTransaction(clientId string, callback func(*RemoteStopTransactionConfirmation, error), transactionId int, props ...func(request *RemoteStopTransactionRequest)) error {
//	request := NewRemoteStopTransactionRequest(transactionId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*RemoteStopTransactionConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Forces a charge point to perform an internal hard or soft reset. In both cases, all ongoing transactions are stopped.
//func (cs *server) Reset(clientId string, callback func(*ResetConfirmation, error), resetType ResetType, props ...func(request *ResetRequest)) error {
//	request := NewResetRequest(resetType)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*ResetConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Attempts to unlock a specific connector on a charge point. Used for remote support purposes.
//func (cs *server) UnlockConnector(clientId string, callback func(*UnlockConnectorConfirmation, error), connectorId int, props ...func(*UnlockConnectorRequest)) error {
//	request := NewUnlockConnectorRequest(connectorId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*UnlockConnectorConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Queries the current version of the local authorization list from a charge point.
//func (cs *server) GetLocalListVersion(clientId string, callback func(*GetLocalListVersionConfirmation, error), props ...func(request *GetLocalListVersionRequest)) error {
//	request := NewGetLocalListVersionRequest()
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetLocalListVersionConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Sends or updates a local authorization list on a charge point. Versioning rules must be followed.
//func (cs *server) SendLocalList(clientId string, callback func(*SendLocalListConfirmation, error), version int, updateType UpdateType, props ...func(request *SendLocalListRequest)) error {
//	request := NewSendLocalListRequest(version, updateType)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*SendLocalListConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Requests diagnostics data from a charge point. The data will be uploaded out-of-band to the provided URL location.
//func (cs *server) GetDiagnostics(clientId string, callback func(*GetDiagnosticsConfirmation, error), location string, props ...func(request *GetDiagnosticsRequest)) error {
//	request := NewGetDiagnosticsRequest(location)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetDiagnosticsConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs the charge point to download and install a new firmware version. The firmware binary will be downloaded out-of-band from the provided URL location.
//func (cs *server) UpdateFirmware(clientId string, callback func(*UpdateFirmwareConfirmation, error), location string, retrieveDate *DateTime, props ...func(request *UpdateFirmwareRequest)) error {
//	request := NewUpdateFirmwareRequest(location, retrieveDate)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*UpdateFirmwareConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs the charge point to reserve a connector for a specific IdTag (client). The connector, or the entire charge point, will be reserved until the provided expiration time.
//func (cs *server) ReserveNow(clientId string, callback func(*ReserveNowConfirmation, error), connectorId int, expiryDate *DateTime, idTag string, reservationId int, props ...func(request *ReserveNowRequest)) error {
//	request := NewReserveNowRequest(connectorId, expiryDate, idTag, reservationId)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*ReserveNowConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Instructs a charge point to send a specific message to the central system. This is used for forcefully triggering status updates, when the last known state is either too old or not clear to the central system.
//func (cs *server) TriggerMessage(clientId string, callback func(*TriggerMessageConfirmation, error), requestedMessage MessageTrigger, props ...func(request *TriggerMessageRequest)) error {
//	request := NewTriggerMessageRequest(requestedMessage)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*TriggerMessageConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Sends a smart charging profile to a charge point. Refer to the smart charging documentation for more information.
//func (cs *server) SetChargingProfile(clientId string, callback func(*SetChargingProfileConfirmation, error), connectorId int, chargingProfile *ChargingProfile, props ...func(request *SetChargingProfileRequest)) error {
//	request := NewSetChargingProfileRequest(connectorId, chargingProfile)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*SetChargingProfileConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}
//
//// Queries a charge point to the composite smart charging schedules and rules for a specified time interval.
//func (cs *server) GetCompositeSchedule(clientId string, callback func(*GetCompositeScheduleConfirmation, error), connectorId int, duration int, props ...func(request *GetCompositeScheduleRequest)) error {
//	request := NewGetCompositeScheduleRequest(connectorId, duration)
//	for _, fn := range props {
//		fn(request)
//	}
//	genericCallback := func(confirmation ocpp.Confirmation, protoError error) {
//		if confirmation != nil {
//			callback(confirmation.(*GetCompositeScheduleConfirmation), protoError)
//		} else {
//			callback(nil, protoError)
//		}
//	}
//	return cs.SendRequestAsync(clientId, request, genericCallback)
//}

func (cs *csms) SetMessageHandler(handler CSMSHandler) {
	cs.coreListener = handler
}

// Registers a handler for incoming local authorization profile messages.
//func (cs *server) SetLocalAuthListListener(listener CentralSystemLocalAuthListListener) {
//	cs.localAuthListListener = listener
//}
//
//// Registers a handler for incoming firmware management profile messages.
//func (cs *server) SetFirmwareManagementListener(listener CentralSystemFirmwareManagementListener) {
//	cs.firmwareListener = listener
//}
//
//// Registers a handler for incoming reservation profile messages.
//func (cs *server) SetReservationListener(listener CentralSystemReservationListener) {
//	cs.reservationListener = listener
//}
//
//// Registers a handler for incoming remote trigger profile messages.
//func (cs *server) SetRemoteTriggerListener(listener CentralSystemRemoteTriggerListener) {
//	cs.remoteTriggerListener = listener
//}
//
//// Registers a handler for incoming smart charging profile messages.
//func (cs *server) SetSmartChargingListener(listener CentralSystemSmartChargingListener) {
//	cs.smartChargingListener = listener
//}

func (cs *csms) SetNewChargingStationHandler(handler func(chargePointId string)) {
	cs.server.SetNewChargePointHandler(handler)
}

func (cs *csms) SetChargingStationDisconnectedHandler(handler func(chargePointId string)) {
	cs.server.SetDisconnectedChargePointHandler(handler)
}

func (cs *csms) SendRequestAsync(clientId string, request ocpp.Request, callback func(confirmation ocpp.Confirmation, err error)) error {
	switch request.GetFeatureName() {
	case CancelReservationFeatureName, CertificateSignedFeatureName, ChangeAvailabilityFeatureName, ClearCacheFeatureName, ClearChargingProfileFeatureName, ClearDisplayFeatureName, ClearVariableMonitoringFeatureName, CostUpdatedFeatureName, CustomerInformationFeatureName, DataTransferFeatureName, DeleteCertificateFeatureName, GetBaseReportFeatureName, GetChargingProfilesFeatureName, GetCompositeScheduleFeatureName, GetDisplayMessagesFeatureName, GetInstalledCertificateIdsFeatureName, GetLocalListVersionFeatureName, GetLogFeatureName, GetMonitoringReportFeatureName:
		break
	//case ChangeConfigurationFeatureName, DataTransferFeatureName, GetConfigurationFeatureName, RemoteStartTransactionFeatureName, RemoteStopTransactionFeatureName, ResetFeatureName, UnlockConnectorFeatureName,
	//	GetLocalListVersionFeatureName, SendLocalListFeatureName,
	//	GetDiagnosticsFeatureName, UpdateFirmwareFeatureName,
	//	ReserveNowFeatureName,
	//	TriggerMessageFeatureName,
	//	SetChargingProfileFeatureName, ClearChargingProfileFeatureName, GetCompositeScheduleFeatureName:
	default:
		return fmt.Errorf("unsupported action %v on central system, cannot send request", request.GetFeatureName())
	}
	cs.callbacks[clientId] = callback
	err := cs.server.SendRequest(clientId, request)
	if err != nil {
		delete(cs.callbacks, clientId)
		return err
	}
	return nil
}

func (cs *csms) Start(listenPort int, listenPath string) {
	cs.server.Start(listenPort, listenPath)
}

func (cs *csms) sendResponse(chargePointId string, confirmation ocpp.Confirmation, err error, requestId string) {
	if confirmation != nil {
		err := cs.server.SendConfirmation(chargePointId, requestId, confirmation)
		if err != nil {
			//TODO: handle error somehow
			log.Print(err)
		}
	} else {
		err := cs.server.SendError(chargePointId, requestId, ocppj.ProtocolError, "Couldn't generate valid confirmation", nil)
		if err != nil {
			log.WithFields(log.Fields{
				"client":  chargePointId,
				"request": requestId,
			}).Errorf("unknown error %v while replying to message with CallError", err)
		}
	}
}

func (cs *csms) notImplementedError(chargePointId string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
	err := cs.server.SendError(chargePointId, requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *csms) notSupportedError(chargePointId string, requestId string, action string) {
	log.Warnf("Cannot handle call %v from charge point %v. Sending CallError instead", requestId, chargePointId)
	err := cs.server.SendError(chargePointId, requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on central system", action), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("unknown error %v while replying to message with CallError", err)
	}
}

func (cs *csms) handleIncomingRequest(chargePointId string, request ocpp.Request, requestId string, action string) {
	profile, found := cs.server.GetProfileForFeature(action)
	// Check whether action is supported and a listener for it exists
	if !found {
		cs.notImplementedError(chargePointId, requestId, action)
		return
	} else {
		switch profile.Name {
		case CoreProfileName:
			if cs.coreListener == nil {
				cs.notSupportedError(chargePointId, requestId, action)
				return
			}
			//case LocalAuthListProfileName:
			//	if cs.localAuthListListener == nil {
			//		cs.notSupportedError(chargePointId, requestId, action)
			//		return
			//	}
			//case FirmwareManagementProfileName:
			//	if cs.firmwareListener == nil {
			//		cs.notSupportedError(chargePointId, requestId, action)
			//		return
			//	}
			//case ReservationProfileName:
			//	if cs.reservationListener == nil {
			//		cs.notSupportedError(chargePointId, requestId, action)
			//		return
			//	}
			//case RemoteTriggerProfileName:
			//	if cs.remoteTriggerListener == nil {
			//		cs.notSupportedError(chargePointId, requestId, action)
			//		return
			//	}
			//case SmartChargingProfileName:
			//	if cs.smartChargingListener == nil {
			//		cs.notSupportedError(chargePointId, requestId, action)
			//		return
			//	}
		}
	}
	var confirmation ocpp.Confirmation = nil
	var err error = nil
	// Execute in separate goroutine, so the caller goroutine is available
	go func() {
		switch action {
		case BootNotificationFeatureName:
			confirmation, err = cs.coreListener.OnBootNotification(chargePointId, request.(*BootNotificationRequest))
		case AuthorizeFeatureName:
			confirmation, err = cs.coreListener.OnAuthorize(chargePointId, request.(*AuthorizeRequest))
		case ClearedChargingLimitFeatureName:
			confirmation, err = cs.coreListener.OnClearedChargingLimit(chargePointId, request.(*ClearedChargingLimitRequest))
		case DataTransferFeatureName:
			confirmation, err = cs.coreListener.OnDataTransfer(chargePointId, request.(*DataTransferRequest))
		case FirmwareStatusNotificationFeatureName:
			confirmation, err = cs.coreListener.OnFirmwareStatusNotification(chargePointId, request.(*FirmwareStatusNotificationRequest))
		case Get15118EVCertificateFeatureName:
			confirmation, err = cs.coreListener.OnGet15118EVCertificate(chargePointId, request.(*Get15118EVCertificateRequest))
		case GetCertificateStatusFeatureName:
			confirmation, err = cs.coreListener.OnGetCertificateStatus(chargePointId, request.(*GetCertificateStatusRequest))
		//case HeartbeatFeatureName:
		//	confirmation, err = cs.messageHandler.OnHeartbeat(chargePointId, request.(*HeartbeatRequest))
		//case MeterValuesFeatureName:
		//	confirmation, err = cs.messageHandler.OnMeterValues(chargePointId, request.(*MeterValuesRequest))
		//case StartTransactionFeatureName:
		//	confirmation, err = cs.messageHandler.OnStartTransaction(chargePointId, request.(*StartTransactionRequest))
		//case StopTransactionFeatureName:
		//	confirmation, err = cs.messageHandler.OnStopTransaction(chargePointId, request.(*StopTransactionRequest))
		//case StatusNotificationFeatureName:
		//	confirmation, err = cs.messageHandler.OnStatusNotification(chargePointId, request.(*StatusNotificationRequest))
		//case DiagnosticsStatusNotificationFeatureName:
		//	confirmation, err = cs.firmwareListener.OnDiagnosticsStatusNotification(chargePointId, request.(*DiagnosticsStatusNotificationRequest))
		//case FirmwareStatusNotificationFeatureName:
		//	confirmation, err = cs.firmwareListener.OnFirmwareStatusNotification(chargePointId, request.(*FirmwareStatusNotificationRequest))
		default:
			cs.notSupportedError(chargePointId, requestId, action)
			return
		}
		cs.sendResponse(chargePointId, confirmation, err, requestId)
	}()
}

func (cs *csms) handleIncomingConfirmation(chargePointId string, confirmation ocpp.Confirmation, requestId string) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(confirmation, nil)
	} else {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": requestId,
		}).Errorf("no handler available for Call Result of type %v", confirmation.GetFeatureName())
	}
}

func (cs *csms) handleIncomingError(chargePointId string, err *ocpp.Error, details interface{}) {
	if callback, ok := cs.callbacks[chargePointId]; ok {
		delete(cs.callbacks, chargePointId)
		callback(nil, err)
	} else {
		log.WithFields(log.Fields{
			"client":  chargePointId,
			"request": err.MessageId,
		}).Errorf("no handler available for Call Error %v", err.Code)
	}
}
