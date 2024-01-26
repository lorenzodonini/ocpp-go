package ocpp16

import (
	"fmt"
	"reflect"

	"github.com/lorenzodonini/ocpp-go/internal/callbackqueue"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/smartcharging"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
)

type chargePoint struct {
	client               *ocppj.Client
	coreHandler          core.ChargePointHandler
	localAuthListHandler localauth.ChargePointHandler
	firmwareHandler      firmware.ChargePointHandler
	reservationHandler   reservation.ChargePointHandler
	remoteTriggerHandler remotetrigger.ChargePointHandler
	smartChargingHandler smartcharging.ChargePointHandler
	confirmationHandler  chan ocpp.Response
	errorHandler         chan error
	callbacks            callbackqueue.CallbackQueue
	stopC                chan struct{}
	errC                 chan error // external error channel
}

func (cp *chargePoint) error(err error) {
	if cp.errC != nil {
		cp.errC <- err
	}
}

// Callback invoked whenever a queued request is canceled, due to timeout.
// By default, the callback returns a GenericError to the caller, who sent the original request.
func (cp *chargePoint) onRequestTimeout(_ string, _ ocpp.Request, err *ocpp.Error) {
	cp.errorHandler <- err
}

// Errors returns a channel for error messages. If it doesn't exist it es created.
func (cp *chargePoint) Errors() <-chan error {
	if cp.errC == nil {
		cp.errC = make(chan error, 1)
	}
	return cp.errC
}

func (cp *chargePoint) BootNotification(chargePointModel string, chargePointVendor string, props ...func(request *core.BootNotificationRequest)) (*core.BootNotificationConfirmation, error) {
	request := core.NewBootNotificationRequest(chargePointModel, chargePointVendor)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.BootNotificationConfirmation), err
	}
}

func (cp *chargePoint) Authorize(idTag string, props ...func(request *core.AuthorizeRequest)) (*core.AuthorizeConfirmation, error) {
	request := core.NewAuthorizationRequest(idTag)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.AuthorizeConfirmation), err
	}
}

func (cp *chargePoint) DataTransfer(vendorId string, props ...func(request *core.DataTransferRequest)) (*core.DataTransferConfirmation, error) {
	request := core.NewDataTransferRequest(vendorId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.DataTransferConfirmation), err
	}
}

func (cp *chargePoint) Heartbeat(props ...func(request *core.HeartbeatRequest)) (*core.HeartbeatConfirmation, error) {
	request := core.NewHeartbeatRequest()
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.HeartbeatConfirmation), err
	}
}

func (cp *chargePoint) MeterValues(connectorId int, meterValues []types.MeterValue, props ...func(request *core.MeterValuesRequest)) (*core.MeterValuesConfirmation, error) {
	request := core.NewMeterValuesRequest(connectorId, meterValues)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.MeterValuesConfirmation), err
	}
}

func (cp *chargePoint) StartTransaction(connectorId int, idTag string, meterStart int, timestamp *types.DateTime, props ...func(request *core.StartTransactionRequest)) (*core.StartTransactionConfirmation, error) {
	request := core.NewStartTransactionRequest(connectorId, idTag, meterStart, timestamp)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.StartTransactionConfirmation), err
	}
}

func (cp *chargePoint) StopTransaction(meterStop int, timestamp *types.DateTime, transactionId int, props ...func(request *core.StopTransactionRequest)) (*core.StopTransactionConfirmation, error) {
	request := core.NewStopTransactionRequest(meterStop, timestamp, transactionId)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.StopTransactionConfirmation), err
	}
}

func (cp *chargePoint) StatusNotification(connectorId int, errorCode core.ChargePointErrorCode, status core.ChargePointStatus, props ...func(request *core.StatusNotificationRequest)) (*core.StatusNotificationConfirmation, error) {
	request := core.NewStatusNotificationRequest(connectorId, errorCode, status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*core.StatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) DiagnosticsStatusNotification(status firmware.DiagnosticsStatus, props ...func(request *firmware.DiagnosticsStatusNotificationRequest)) (*firmware.DiagnosticsStatusNotificationConfirmation, error) {
	request := firmware.NewDiagnosticsStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*firmware.DiagnosticsStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) FirmwareStatusNotification(status firmware.FirmwareStatus, props ...func(request *firmware.FirmwareStatusNotificationRequest)) (*firmware.FirmwareStatusNotificationConfirmation, error) {
	request := firmware.NewFirmwareStatusNotificationRequest(status)
	for _, fn := range props {
		fn(request)
	}
	confirmation, err := cp.SendRequest(request)
	if err != nil {
		return nil, err
	} else {
		return confirmation.(*firmware.FirmwareStatusNotificationConfirmation), err
	}
}

func (cp *chargePoint) SetCoreHandler(handler core.ChargePointHandler) {
	cp.coreHandler = handler
}

func (cp *chargePoint) SetLocalAuthListHandler(handler localauth.ChargePointHandler) {
	cp.localAuthListHandler = handler
}

func (cp *chargePoint) SetFirmwareManagementHandler(handler firmware.ChargePointHandler) {
	cp.firmwareHandler = handler
}

func (cp *chargePoint) SetReservationHandler(handler reservation.ChargePointHandler) {
	cp.reservationHandler = handler
}

func (cp *chargePoint) SetRemoteTriggerHandler(handler remotetrigger.ChargePointHandler) {
	cp.remoteTriggerHandler = handler
}

func (cp *chargePoint) SetSmartChargingHandler(handler smartcharging.ChargePointHandler) {
	cp.smartChargingHandler = handler
}

func (cp *chargePoint) SendRequest(request ocpp.Request) (ocpp.Response, error) {
	featureName := request.GetFeatureName()
	if _, found := cp.client.GetProfileForFeature(featureName); !found {
		return nil, fmt.Errorf("feature %v is unsupported on charge point (missing profile), cannot send request", featureName)
	}

	// Wraps an asynchronous response
	type asyncResponse struct {
		r ocpp.Response
		e error
	}
	// Create channel and pass it to a callback function, for retrieving asynchronous response
	asyncResponseC := make(chan asyncResponse, 1)
	send := func() error {
		return cp.client.SendRequest(request)
	}
	err := cp.callbacks.TryQueue("main", send, func(confirmation ocpp.Response, err error) {
		asyncResponseC <- asyncResponse{r: confirmation, e: err}
	})
	if err != nil {
		return nil, err
	}
	select {
	case asyncResult, ok := <-asyncResponseC:
		if !ok {
			return nil, fmt.Errorf("internal error while receiving result for %v request", request.GetFeatureName())
		}
		return asyncResult.r, asyncResult.e
	case <-cp.stopC:
		return nil, fmt.Errorf("client stopped while waiting for response to %v", request.GetFeatureName())
	}
}

func (cp *chargePoint) SendRequestAsync(request ocpp.Request, callback func(confirmation ocpp.Response, err error)) error {
	featureName := request.GetFeatureName()
	if _, found := cp.client.GetProfileForFeature(featureName); !found {
		return fmt.Errorf("feature %v is unsupported on charge point (missing profile), cannot send request", featureName)
	}
	switch featureName {
	case core.AuthorizeFeatureName, core.BootNotificationFeatureName, core.DataTransferFeatureName, core.HeartbeatFeatureName, core.MeterValuesFeatureName, core.StartTransactionFeatureName, core.StopTransactionFeatureName, core.StatusNotificationFeatureName,
		firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		break
	default:
		return fmt.Errorf("unsupported action %v on charge point, cannot send request", featureName)
	}
	// Response will be retrieved asynchronously via asyncHandler
	send := func() error {
		return cp.client.SendRequest(request)
	}
	err := cp.callbacks.TryQueue("main", send, callback)
	return err
}

func (cp *chargePoint) asyncCallbackHandler() {
	for {
		select {
		case confirmation := <-cp.confirmationHandler:
			// Get and invoke callback
			if callback, ok := cp.callbacks.Dequeue("main"); ok {
				callback(confirmation, nil)
			} else {
				err := fmt.Errorf("no handler available for incoming response %v", confirmation.GetFeatureName())
				cp.error(err)
			}
		case protoError := <-cp.errorHandler:
			// Get and invoke callback
			if callback, ok := cp.callbacks.Dequeue("main"); ok {
				callback(nil, protoError)
			} else {
				err := fmt.Errorf("no handler available for error %v", protoError.Error())
				cp.error(err)
			}
		case <-cp.stopC:
			// Handler stopped, cleanup callbacks.
			// No callback invocation, since the user manually stopped the client.
			cp.clearCallbacks(false)
			return
		}
	}
}

func (cp *chargePoint) clearCallbacks(invokeCallback bool) {
	for cb, ok := cp.callbacks.Dequeue("main"); ok; cb, ok = cp.callbacks.Dequeue("main") {
		if invokeCallback {
			err := ocpp.NewError(ocppj.GenericError, "client stopped, no response received from server", "")
			cb(nil, err)
		}
	}
}

func (cp *chargePoint) sendResponse(confirmation ocpp.Response, err error, requestId string) {
	if err != nil {
		// Send error response
		if ocppError, ok := err.(*ocpp.Error); ok {
			err = cp.client.SendError(requestId, ocppError.Code, ocppError.Description, nil)
		} else {
			err = cp.client.SendError(requestId, ocppj.InternalError, err.Error(), nil)
		}
		if err != nil {
			// Error while sending an error. Will attempt to send a default error instead
			cp.client.HandleFailedResponseError(requestId, err, "")
			// Notify client implementation
			err = fmt.Errorf("replying to request %s with 'internal error' failed: %w", requestId, err)
			cp.error(err)
		}
		return
	}

	if confirmation == nil || reflect.ValueOf(confirmation).IsNil() {
		err = fmt.Errorf("empty confirmation to request %s", requestId)
		// Sending a dummy error to server instead, then notify client implementation
		_ = cp.client.SendError(requestId, ocppj.GenericError, err.Error(), nil)
		cp.error(err)
		return
	}

	// send confirmation response
	err = cp.client.SendResponse(requestId, confirmation)
	if err != nil {
		// Error while sending an error. Will attempt to send a default error instead
		cp.client.HandleFailedResponseError(requestId, err, confirmation.GetFeatureName())
		// Notify client implementation
		err = fmt.Errorf("failed responding to request %s: %w", requestId, err)
		cp.error(err)
	}
}

func (cp *chargePoint) Start(centralSystemUrl string) error {
	// Start client
	cp.stopC = make(chan struct{}, 1)
	err := cp.client.Start(centralSystemUrl)
	// Async response handler receives incoming responses/errors and triggers callbacks
	if err == nil {
		go cp.asyncCallbackHandler()
	}
	return err
}

func (cp *chargePoint) Stop() {
	cp.client.Stop()
	close(cp.stopC)

	if cp.errC != nil {
		close(cp.errC)
		cp.errC = nil
	}
}

func (cp *chargePoint) IsConnected() bool {
	return cp.client.IsConnected()
}

func (cp *chargePoint) notImplementedError(requestId string, action string) {
	err := cp.client.SendError(requestId, ocppj.NotImplemented, fmt.Sprintf("no handler for action %v implemented", action), nil)
	if err != nil {
		err = fmt.Errorf("replying cs to request %s with 'not implemented': %w", requestId, err)
		cp.error(err)
	}
}

func (cp *chargePoint) notSupportedError(requestId string, action string) {
	err := cp.client.SendError(requestId, ocppj.NotSupported, fmt.Sprintf("unsupported action %v on charge point", action), nil)
	if err != nil {
		err = fmt.Errorf("replying cs to request %s with 'not supported': %w", requestId, err)
		cp.error(err)
	}
}

func (cp *chargePoint) handleIncomingRequest(request ocpp.Request, requestId string, action string) {
	profile, found := cp.client.GetProfileForFeature(action)
	// Check whether action is supported and a handler for it exists
	if !found {
		cp.notImplementedError(requestId, action)
		return
	} else {
		switch profile.Name {
		case core.ProfileName:
			if cp.coreHandler == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case localauth.ProfileName:
			if cp.localAuthListHandler == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case firmware.ProfileName:
			if cp.firmwareHandler == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case reservation.ProfileName:
			if cp.reservationHandler == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case remotetrigger.ProfileName:
			if cp.remoteTriggerHandler == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		case smartcharging.ProfileName:
			if cp.smartChargingHandler == nil {
				cp.notSupportedError(requestId, action)
				return
			}
		}
	}
	// Process request
	var confirmation ocpp.Response
	cp.client.GetProfileForFeature(action)
	var err error
	switch action {
	case core.ChangeAvailabilityFeatureName:
		confirmation, err = cp.coreHandler.OnChangeAvailability(request.(*core.ChangeAvailabilityRequest))
	case core.ChangeConfigurationFeatureName:
		confirmation, err = cp.coreHandler.OnChangeConfiguration(request.(*core.ChangeConfigurationRequest))
	case core.ClearCacheFeatureName:
		confirmation, err = cp.coreHandler.OnClearCache(request.(*core.ClearCacheRequest))
	case core.DataTransferFeatureName:
		confirmation, err = cp.coreHandler.OnDataTransfer(request.(*core.DataTransferRequest))
	case core.GetConfigurationFeatureName:
		confirmation, err = cp.coreHandler.OnGetConfiguration(request.(*core.GetConfigurationRequest))
	case core.RemoteStartTransactionFeatureName:
		confirmation, err = cp.coreHandler.OnRemoteStartTransaction(request.(*core.RemoteStartTransactionRequest))
	case core.RemoteStopTransactionFeatureName:
		confirmation, err = cp.coreHandler.OnRemoteStopTransaction(request.(*core.RemoteStopTransactionRequest))
	case core.ResetFeatureName:
		confirmation, err = cp.coreHandler.OnReset(request.(*core.ResetRequest))
	case core.UnlockConnectorFeatureName:
		confirmation, err = cp.coreHandler.OnUnlockConnector(request.(*core.UnlockConnectorRequest))
	case localauth.GetLocalListVersionFeatureName:
		confirmation, err = cp.localAuthListHandler.OnGetLocalListVersion(request.(*localauth.GetLocalListVersionRequest))
	case localauth.SendLocalListFeatureName:
		confirmation, err = cp.localAuthListHandler.OnSendLocalList(request.(*localauth.SendLocalListRequest))
	case firmware.GetDiagnosticsFeatureName:
		confirmation, err = cp.firmwareHandler.OnGetDiagnostics(request.(*firmware.GetDiagnosticsRequest))
	case firmware.UpdateFirmwareFeatureName:
		confirmation, err = cp.firmwareHandler.OnUpdateFirmware(request.(*firmware.UpdateFirmwareRequest))
	case reservation.ReserveNowFeatureName:
		confirmation, err = cp.reservationHandler.OnReserveNow(request.(*reservation.ReserveNowRequest))
	case reservation.CancelReservationFeatureName:
		confirmation, err = cp.reservationHandler.OnCancelReservation(request.(*reservation.CancelReservationRequest))
	case remotetrigger.TriggerMessageFeatureName:
		confirmation, err = cp.remoteTriggerHandler.OnTriggerMessage(request.(*remotetrigger.TriggerMessageRequest))
	case smartcharging.SetChargingProfileFeatureName:
		confirmation, err = cp.smartChargingHandler.OnSetChargingProfile(request.(*smartcharging.SetChargingProfileRequest))
	case smartcharging.ClearChargingProfileFeatureName:
		confirmation, err = cp.smartChargingHandler.OnClearChargingProfile(request.(*smartcharging.ClearChargingProfileRequest))
	case smartcharging.GetCompositeScheduleFeatureName:
		confirmation, err = cp.smartChargingHandler.OnGetCompositeSchedule(request.(*smartcharging.GetCompositeScheduleRequest))
	default:
		cp.notSupportedError(requestId, action)
		return
	}
	cp.sendResponse(confirmation, err, requestId)
}
