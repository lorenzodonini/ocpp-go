package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/availability"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/diagnostics"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/provisioning"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/remotecontrol"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
)

func (handler *ChargingStationHandler) OnRequestStartTransaction(request *remotecontrol.RequestStartTransactionRequest) (response *remotecontrol.RequestStartTransactionResponse, err error) {
	if request.EvseID != nil {
		evse, ok := handler.evse[*request.EvseID]
		if !ok || evse.availability != availability.OperationalStatusOperative {
			return remotecontrol.NewRequestStartTransactionResponse(remotecontrol.RequestStartStopStatusRejected), nil
		}
		// Find occupied connector
		connectorID := 0
		for i, c := range evse.connectors {
			if c.status == availability.ConnectorStatusOccupied {
				connectorID = i
				break
			}
		}
		evse.currentTransaction = nextTransactionID()
		logDefault(request.GetFeatureName()).Infof("started transaction %v on evse %v, connector %v", evse.currentTransaction, *request.EvseID, connectorID)
		response = remotecontrol.NewRequestStartTransactionResponse(remotecontrol.RequestStartStopStatusAccepted)
		response.TransactionID = fmt.Sprintf("%s", evse.currentTransaction)
		return response, nil
	}
	logDefault(request.GetFeatureName()).Errorf("couldn't start a transaction for token %v without an evseID", request.IDToken)
	return remotecontrol.NewRequestStartTransactionResponse(remotecontrol.RequestStartStopStatusRejected), nil
}

func (handler *ChargingStationHandler) OnRequestStopTransaction(request *remotecontrol.RequestStopTransactionRequest) (response *remotecontrol.RequestStopTransactionResponse, err error) {
	for key, evse := range handler.evse {
		if evse.currentTransaction == request.TransactionID {
			logDefault(request.GetFeatureName()).Infof("stopped transaction %v on evse %v", evse.currentTransaction, key)
			evse.currentTransaction = ""
			evse.currentReservation = 0
			// Find the currently occupied connector
			for i, c := range evse.connectors {
				if c.status == availability.ConnectorStatusOccupied {
					connector := evse.connectors[i]
					connector.status = availability.ConnectorStatusAvailable
					evse.connectors[i] = connector
					break
				}
			}
			return remotecontrol.NewRequestStopTransactionResponse(remotecontrol.RequestStartStopStatusAccepted), nil
		}
	}
	logDefault(request.GetFeatureName()).Errorf("couldn't stop transaction %v, no such transaction is ongoing", request.TransactionID)
	return remotecontrol.NewRequestStopTransactionResponse(remotecontrol.RequestStartStopStatusRejected), nil
}

func (handler *ChargingStationHandler) OnTriggerMessage(request *remotecontrol.TriggerMessageRequest) (response *remotecontrol.TriggerMessageResponse, err error) {
	logDefault(request.GetFeatureName()).Infof("received trigger for %v", request.RequestedMessage)
	status := remotecontrol.TriggerMessageStatusRejected
	switch request.RequestedMessage {
	case remotecontrol.MessageTriggerBootNotification:
		// Boot Notification
		go func() {
			_, e := chargingStation.BootNotification(provisioning.BootReasonTriggered, handler.model, handler.vendor)
			checkError(e)
			logDefault(provisioning.BootNotificationFeatureName).Info("boot notification completed")
		}()
		status = remotecontrol.TriggerMessageStatusAccepted
	case remotecontrol.MessageTriggerLogStatusNotification:
		// Log Status Notification
		go func() {
			reqID := rand.Int()
			_, e := chargingStation.LogStatusNotification(diagnostics.UploadLogStatusUploading, reqID)
			checkError(e)
			logDefault(diagnostics.LogStatusNotificationFeatureName).Info("diagnostics status notified")
		}()
		status = remotecontrol.TriggerMessageStatusAccepted
	case remotecontrol.MessageTriggerFirmwareStatusNotification:
		//TODO: schedule firmware status notification message
		status = remotecontrol.TriggerMessageStatusAccepted
	case remotecontrol.MessageTriggerHeartbeat:
		// Schedule heartbeat request
		go func() {
			resp, e := chargingStation.Heartbeat()
			checkError(e)
			logDefault(availability.HeartbeatFeatureName).Infof("clock synchronized: %v", resp.CurrentTime.FormatTimestamp())
		}()
		status = remotecontrol.TriggerMessageStatusAccepted
	case remotecontrol.MessageTriggerMeterValues:
		// Schedule meter values update
		//TODO: schedule meter values message
		break
	case remotecontrol.MessageTriggerStatusNotification:
		// Schedule connector status notification
		if request.Evse != nil {
			connectorStatus := availability.ConnectorStatusUnavailable
			evse, ok := handler.evse[request.Evse.ID]
			if !ok {
				status = remotecontrol.TriggerMessageStatusRejected
				break
			}
			if request.Evse.ConnectorID != nil {
				if !evse.hasConnector(*request.Evse.ConnectorID) {
					status = remotecontrol.TriggerMessageStatusRejected
					break
				}
				connectorStatus = evse.connectors[*request.Evse.ConnectorID].status
			} else {
				status = remotecontrol.TriggerMessageStatusRejected
				break
			}
			// Update asynchronously
			go func() {
				_, e := chargingStation.StatusNotification(types.NewDateTime(time.Now()), connectorStatus, request.Evse.ID, *request.Evse.ConnectorID)
				checkError(e)
				logDefault(availability.HeartbeatFeatureName).Infof("status for connector %v sent: %v", *request.Evse.ConnectorID, connectorStatus)
			}()
			status = remotecontrol.TriggerMessageStatusAccepted
		} else {
			status = remotecontrol.TriggerMessageStatusRejected
		}
	case remotecontrol.MessageTriggerTransactionEvent:
		// TODO:
		break
	default:
		// We're not implementing support for other messages
		status = remotecontrol.TriggerMessageStatusNotImplemented
	}
	return remotecontrol.NewTriggerMessageResponse(status), nil
}

func (handler *ChargingStationHandler) OnUnlockConnector(request *remotecontrol.UnlockConnectorRequest) (response *remotecontrol.UnlockConnectorResponse, err error) {
	evse, ok := handler.evse[request.EvseID]
	if !ok || !evse.hasConnector(request.ConnectorID) {
		logDefault(request.GetFeatureName()).Errorf("couldn't unlock unknown connector %d for EVSE %d", request.ConnectorID, request.EvseID)
		return remotecontrol.NewUnlockConnectorResponse(remotecontrol.UnlockStatusUnknownConnector), nil
	}
	connector := evse.connectors[request.ConnectorID]
	// TODO: unlock connector internally
	connector.status = availability.ConnectorStatusAvailable
	logDefault(request.GetFeatureName()).Infof("unlocked connector %v for EVSE %d", request.ConnectorID, request.EvseID)
	return remotecontrol.NewUnlockConnectorResponse(remotecontrol.UnlockStatusUnlocked), nil
}
