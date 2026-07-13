package der

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
)

// Needs to be implemented by a CSMS for handling messages part of the OCPP 2.1 DER profile.
type CSMSHandler interface {
	OnNotifyDERStartStop(chargingStationId string, req *NotifyDERStartStopRequest) (res *NotifyDERStartStopResponse, err error)
	OnNotifyDERAlarm(chargingStationId string, req *NotifyDERAlarmRequest) (res *NotifyDERAlarmResponse, err error)
	OnReportDERControl(chargingStationId string, req *ReportDERControlRequest) (res *ReportDERControlResponse, err error)
}

// Needs to be implemented by Charging stations for handling messages part of the OCPP 2.1 DER profile.
type ChargingStationHandler interface {
	OnGetDERControl(req *GetDERControlRequest) (res *GetDERControlResponse, err error)
	OnSetDERControl(req *SetDERControlRequest) (res *SetDERControlResponse, err error)
	OnClearDERControl(req *ClearDERControlRequest) (res *ClearDERControlResponse, err error)
}

const ProfileName = "DERControl"

var Profile = ocpp.NewProfile(
	ProfileName,
	GetDERControlFeature{},
	SetDERControlFeature{},
	ClearDERControlFeature{},
	NotifyDERStartStopFeature{},
	NotifyDERAlarmFeature{},
	ReportDERControlFeature{},
)
