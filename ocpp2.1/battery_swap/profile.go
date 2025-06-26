package battery_swap

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 Battery Swap.
type CSMSHandler interface {
	OnBatterySwap(chargingStationID string, request *BatterySwapRequest) (*BatterySwapResponse, error)
}

// Needs to be implemented by Charging stations for handling messages part of the Battery Swap.
type ChargingStationHandler interface {
	OnRequestBatterySwap(request *RequestBatterySwapRequest) (*RequestBatterySwapResponse, error)
}

const ProfileName = "BatterySwap"

var Profile = ocpp.NewProfile(
	ProfileName,
	BatterySwapFeature{},
	RequestBatterySwapFeature{},
)
