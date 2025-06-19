package battey_swap

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.1/types"
	"reflect"
)

// -------------------- BatterySwap (CS -> CSMS) --------------------

const BatterySwap = "BatterySwap"

// The field definition of the BatterySwapRequest request payload sent by the CSMS to the Charging Station.
type BatterySwapRequest struct {
	EventType   BatterSwapEvent   `json:"eventType" validate:"required,batterySwapEvent"`
	RequestId   int               `json:"requestId" validate:"required"`        // Unique identifier for the request
	IdToken     types.IdTokenType `json:"idToken" validate:"required,dive"`     // Optional field for the ID token of the user
	BatteryData BatteryData       `json:"batteryData" validate:"required,dive"` // Contains information about the battery to be swapped
}

type BatteryData struct {
	EvseId         int             `json:"evseId" validate:"required,gte=1"`        // The ID of the EVSE where the battery swap is taking place
	SerialNumber   string          `json:"serialNumber" validate:"required,max=50"` // The serial number of the battery being swapped
	SoC            float64         `json:"soC" validate:"required,gte=0,lte=100"`
	SoH            float64         `json:"soH" validate:"required,gte=0,lte=100"`
	ProductionDate *types.DateTime `json:"productionDate" validate:"omitempty"`
	VendorInfo     string          `json:"vendorInfo" validate:"omitempty,max=500"`
}

// This field definition of the BatterySwapResponse
type BatterySwapResponse struct {
}

type BatterySwapFeature struct{}

func (f BatterySwapFeature) GetFeatureName() string {
	return BatterySwap
}

func (f BatterySwapFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(BatterySwapRequest{})
}

func (f BatterySwapFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(BatterySwapResponse{})
}

func (r BatterySwapRequest) GetFeatureName() string {
	return BatterySwap
}

func (c BatterySwapResponse) GetFeatureName() string {
	return BatterySwap
}

func NewBatterySwapRequest(eventType BatterSwapEvent, requestId int, idToken types.IdTokenType, batteryData BatteryData) BatterySwapRequest {
	return BatterySwapRequest{
		EventType:   eventType,
		RequestId:   requestId,
		IdToken:     idToken,
		BatteryData: batteryData,
	}
}

func NewBatteryData(evseId int, serialNumber string, soc, soh float64) BatteryData {
	return BatteryData{
		EvseId:       evseId,
		SerialNumber: serialNumber,
		SoC:          soc,
		SoH:          soh,
	}
}

func NewBatterySwapResponse() BatterySwapResponse {
	return BatterySwapResponse{}
}
