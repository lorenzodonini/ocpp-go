package diagnostics

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0/types"
	"gopkg.in/go-playground/validator.v9"
)

// MonitorType specifies the type of this monitor.
type MonitorType string

const (
	MonitorUpperThreshold       MonitorType = "UpperThreshold"       // Triggers an event notice when the actual value of the Variable rises above monitorValue.
	MonitorLowerThreshold       MonitorType = "LowerThreshold"       // Triggers an event notice when the actual value of the Variable drops below monitorValue.
	MonitorDelta                MonitorType = "Delta"                // Triggers an event notice when the actual value has changed more than plus or minus monitorValue since the time that this monitor was set or since the last time this event notice was sent, whichever was last.
	MonitorPeriodic             MonitorType = "Periodic"             // Triggers an event notice every monitorValue seconds interval, starting from the time that this monitor was set.
	MonitorPeriodicClockAligned MonitorType = "PeriodicClockAligned" // Triggers an event notice every monitorValue seconds interval, starting from the nearest clock-aligned interval after this monitor was set.
)

func isValidMonitorType(fl validator.FieldLevel) bool {
	status := MonitorType(fl.Field().String())
	switch status {
	case MonitorUpperThreshold, MonitorLowerThreshold, MonitorDelta, MonitorPeriodic, MonitorPeriodicClockAligned:
		return true
	default:
		return false
	}
}

func init() {
	_ = types.Validate.RegisterValidation("monitorType", isValidMonitorType)
}
