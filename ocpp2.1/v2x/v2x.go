package v2x

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 V2X.
type CSMSHandler interface {
}

// Needs to be implemented by Charging stations for handling messages part of the V2X.
type ChargingStationHandler interface {
	OnAFRRSignal(chargingStationId string, request *AFRRSignalRequest) (*AFRRSignalResponse, error)
	OnNotifyAllowedEnergyTransfer(chargingStationId string, request *NotifyAllowedEnergyTransferRequest) (*NotifyAllowedEnergyTransferResponse, error)
}

const ProfileName = "V2X"

var Profile = ocpp.NewProfile(
	ProfileName,
	AFRRSignalFeature{},
	NotifyAllowedEnergyTransferFeature{},
)
