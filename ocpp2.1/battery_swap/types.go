package battery_swap

import (
	"github.com/xBlaz3kx/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
)

type BatterSwapEvent string

const (
	BatterSwapEventBatteryIn         BatterSwapEvent = "BatteryIn"
	BatterSwapEventBatteryOut        BatterSwapEvent = "BatteryOut"
	BatterSwapEventBatteryOutTimeout BatterSwapEvent = "BatteryOutTimeout"
)

func isValidBatterSwapEvent(fl validator.FieldLevel) bool {
	event := fl.Field().String()
	switch event {
	case string(BatterSwapEventBatteryIn), string(BatterSwapEventBatteryOut), string(BatterSwapEventBatteryOutTimeout):
		return true
	default:
		return false
	}
}

func init() {
	_ = ocppj.Validate.RegisterValidation("batterySwapEvent", isValidBatterSwapEvent)
}
