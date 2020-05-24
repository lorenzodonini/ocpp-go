// Contains the Basic Charge Point functionality comparable with OCPP 1.5.
package core

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by Central systems for handling messages part of the OCPP 1.6 Core profile.
type CentralSystemHandler interface {
	OnAuthorize(chargePointId string, request *AuthorizeRequest) (confirmation *AuthorizeConfirmation, err error)
	OnBootNotification(chargePointId string, request *BootNotificationRequest) (confirmation *BootNotificationConfirmation, err error)
	OnDataTransfer(chargePointId string, request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	OnHeartbeat(chargePointId string, request *HeartbeatRequest) (confirmation *HeartbeatConfirmation, err error)
	OnMeterValues(chargePointId string, request *MeterValuesRequest) (confirmation *MeterValuesConfirmation, err error)
	OnStatusNotification(chargePointId string, request *StatusNotificationRequest) (confirmation *StatusNotificationConfirmation, err error)
	OnStartTransaction(chargePointId string, request *StartTransactionRequest) (confirmation *StartTransactionConfirmation, err error)
	OnStopTransaction(chargePointId string, request *StopTransactionRequest) (confirmation *StopTransactionConfirmation, err error)
}

// Needs to be implemented by Charge points for handling messages part of the OCPP 1.6 Core profile.
type ChargePointHandler interface {
	OnChangeAvailability(request *ChangeAvailabilityRequest) (confirmation *ChangeAvailabilityConfirmation, err error)
	OnChangeConfiguration(request *ChangeConfigurationRequest) (confirmation *ChangeConfigurationConfirmation, err error)
	OnClearCache(request *ClearCacheRequest) (confirmation *ClearCacheConfirmation, err error)
	OnDataTransfer(request *DataTransferRequest) (confirmation *DataTransferConfirmation, err error)
	OnGetConfiguration(request *GetConfigurationRequest) (confirmation *GetConfigurationConfirmation, err error)
	OnRemoteStartTransaction(request *RemoteStartTransactionRequest) (confirmation *RemoteStartTransactionConfirmation, err error)
	OnRemoteStopTransaction(request *RemoteStopTransactionRequest) (confirmation *RemoteStopTransactionConfirmation, err error)
	OnReset(request *ResetRequest) (confirmation *ResetConfirmation, err error)
	OnUnlockConnector(request *UnlockConnectorRequest) (confirmation *UnlockConnectorConfirmation, err error)
}

// THe profile name
var ProfileName = "core"

// Provides support for Basic Charge Point functionality comparable with OCPP 1.5.
var Profile = ocpp.NewProfile(
	ProfileName,
	BootNotificationFeature{},
	AuthorizeFeature{},
	ChangeAvailabilityFeature{},
	ChangeConfigurationFeature{},
	ClearCacheFeature{},
	DataTransferFeature{},
	GetConfigurationFeature{},
	HeartbeatFeature{},
	MeterValuesFeature{},
	RemoteStartTransactionFeature{},
	RemoteStopTransactionFeature{},
	StartTransactionFeature{},
	StopTransactionFeature{},
	StatusNotificationFeature{},
	ResetFeature{},
	UnlockConnectorFeature{})
