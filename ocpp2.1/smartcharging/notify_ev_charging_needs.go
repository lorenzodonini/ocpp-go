package smartcharging

import (
	"reflect"

	"github.com/xBlaz3kx/ocpp-go/ocpp2.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// -------------------- Notify EV Charging Needs (CS -> CSMS) --------------------

const NotifyEVChargingNeedsFeatureName = "NotifyEVChargingNeeds"

func isValidEnergyTransferMode(fl validator.FieldLevel) bool {
	status := types.EnergyTransferMode(fl.Field().String())
	switch status {
	case types.EnergyTransferModeAC1Phase, types.EnergyTransferModeAC2Phase, types.EnergyTransferModeAC3Phase, types.EnergyTransferModeDC:
		return true
	default:
		return false
	}
}

// EVChargingNeedsStatus contains the status returned by the CSMS.
type EVChargingNeedsStatus string

const (
	EVChargingNeedsStatusAccepted          EVChargingNeedsStatus = "Accepted"
	EVChargingNeedsStatusRejected          EVChargingNeedsStatus = "Rejected"
	EVChargingNeedsStatusNoChargingProfile EVChargingNeedsStatus = "NoChargingProfile"
	EVChargingNeedsStatusProcessing        EVChargingNeedsStatus = "Processing"
)

func isValidEVChargingNeedsStatus(fl validator.FieldLevel) bool {
	status := EVChargingNeedsStatus(fl.Field().String())
	switch status {
	case EVChargingNeedsStatusAccepted, EVChargingNeedsStatusRejected, EVChargingNeedsStatusProcessing, EVChargingNeedsStatusNoChargingProfile:
		return true
	default:
		return false
	}
}

// ACChargingParameters contains EV AC charging parameters. Used by ChargingNeeds.
type ACChargingParameters struct {
	EnergyAmount int `json:"energyAmount" validate:"gte=0"` // Amount of energy requested (in Wh). This includes energy required for preconditioning.
	EVMinCurrent int `json:"evMinCurrent" validate:"gte=0"` // Minimum current (amps) supported by the electric vehicle (per phase).
	EVMaxCurrent int `json:"evMaxCurrent" validate:"gte=0"` // Maximum current (amps) supported by the electric vehicle (per phase). Includes cable capacity.
	EVMaxVoltage int `json:"evMaxVoltage" validate:"gte=0"` // Maximum voltage supported by the electric vehicle.
}

// DCChargingParameters contains EV DC charging parameters. Used by ChargingNeeds.
type DCChargingParameters struct {
	EVMaxCurrent     int  `json:"evMaxCurrent" validate:"gte=0"`                              // Maximum current (amps) supported by the electric vehicle (per phase). Includes cable capacity.
	EVMaxVoltage     int  `json:"evMaxVoltage" validate:"gte=0"`                              // Maximum voltage supported by the electric vehicle.
	EnergyAmount     *int `json:"energyAmount,omitempty" validate:"omitempty,gte=0"`          // Amount of energy requested (in Wh). This includes energy required for preconditioning.
	EVMaxPower       *int `json:"evMaxPower,omitempty" validate:"omitempty,gte=0"`            // Maximum power (in W) supported by the electric vehicle. Required for DC charging.
	StateOfCharge    *int `json:"stateOfCharge,omitempty" validate:"omitempty,gte=0,lte=100"` // Energy available in the battery (in percent of the battery capacity).
	EVEnergyCapacity *int `json:"evEnergyCapacity,omitempty" validate:"omitempty,gte=0"`      // Capacity of the electric vehicle battery (in Wh)
	FullSoC          *int `json:"fullSoC,omitempty" validate:"omitempty,gte=0,lte=100"`       // Percentage of SoC at which the EV considers the battery fully charged. (possible values: 0 - 100)
	BulkSoC          *int `json:"bulkSoC,omitempty" validate:"omitempty,gte=0,lte=100"`       // Percentage of SoC at which the EV considers a fast charging process to end. (possible values: 0 - 100)
}

// ChargingNeeds contains the characteristics of the energy delivery required. Used by NotifyEVChargingNeedsRequest.
type ChargingNeeds struct {
	RequestedEnergyTransfer types.EnergyTransferMode   `json:"requestedEnergyTransfer" validate:"required,energyTransferMode21"`            // Mode of energy transfer requested by the EV.
	DepartureTime           *types.DateTime            `json:"departureTime,omitempty" validate:"omitempty"`                                // Estimated departure time of the EV.
	ACChargingParameters    *ACChargingParameters      `json:"acChargingParameters,omitempty" validate:"omitempty,dive"`                    // AC charging parameters.
	DCChargingParameters    *DCChargingParameters      `json:"dcChargingParameters,omitempty" validate:"omitempty,dive"`                    // DC charging parameters.
	AvailableEnergyTransfer []types.EnergyTransferMode `json:"availableEnergyTransfer,omitempty" validate:"omitempty,energyTransferMode21"` // Energy transfer modes supported by the Charging Station.
	ControlMode             *ControlMode               `json:"controlMode,omitempty" validate:"omitempty,controlMode"`
	MobilityNeedsMode       *MobilityNeedsMode         `json:"mobilityNeedsMode,omitempty" validate:"omitempty,mobilityNeedsMode21"` // Mobility needs mode requested by the EV.
	V2XCharginParameters    *V2XChargingParameters     `json:"v2xChargingParameters,omitempty" validate:"omitempty,dive"`            // V2X charging parameters.
	EvEnergyOffer           *EvEnergyOffer             `json:"evEnergyOffer,omitempty" validate:"omitempty,dive"`                    // EV energy offer.
	DERChargingParameters   *DERChargingParameters     `json:"derChargingParameters,omitempty" validate:"omitempty,dive"`            // DER charging parameters.
}

type DERChargingParameters struct {
	EVSupportedDERControl           []string `json:"evSupportedDERControl,omitempty" validate:"omitempty"` // Supported DER control modes by the EV.
	EVOverExcitedMaxDischargePower  *float64 `json:"evOverExcitedMaxDischargePower,omitempty" validate:"omitempty"`
	EVOverExcitedPowerFactor        *float64 `json:"evOverExcitedPowerFactor,omitempty" validate:"omitempty"`
	EVUnderExcitedMaxDischargePower *float64 `json:"evUnderExcitedMaxDischargePower,omitempty" validate:"omitempty"`
	EVUnderExcitedPowerFactor       *float64 `json:"evUnderExcitedPowerFactor,omitempty" validate:"omitempty"`
	MaxApparentPower                *float64 `json:"maxApparentPower,omitempty" validate:"omitempty"` // Maximum apparent power (in VA) supported by the EV.
	MaxChargeApparentPower          *float64 `json:"maxChargeApparentPower,omitempty" validate:"omitempty"`
	MaxChargeApparentPowerL2        *float64 `json:"maxChargeApparentPower_L2,omitempty" validate:"omitempty"`
	MaxChargeApparentPowerL3        *float64 `json:"maxChargeApparentPower_L3,omitempty" validate:"omitempty"`
	MaxDischargeApparentPower       *float64 `json:"maxDischargeApparentPower,omitempty" validate:"omitempty"`
	MaxDischargeApparentPowerL2     *float64 `json:"maxDischargeApparentPower_L2,omitempty" validate:"omitempty"`
	MaxDischargeApparentPowerL3     *float64 `json:"maxDischargeApparentPower_L3,omitempty" validate:"omitempty"`
	MaxChargeReactivePower          *float64 `json:"maxChargeReactivePower,omitempty" validate:"omitempty"`
	MaxChargeReactivePowerL2        *float64 `json:"maxChargeReactivePower_L2,omitempty" validate:"omitempty"`
	MaxChargeReactivePowerL3        *float64 `json:"maxChargeReactivePower_L3,omitempty" validate:"omitempty"`
}

type V2XChargingParameters struct {
}

type EvEnergyOffer struct {
}

type ControlMode string

const ()

type MobilityNeedsMode string

const ()

// The field definition of the NotifyEVChargingNeeds request payload sent by the Charging Station to the CSMS.
type NotifyEVChargingNeedsRequest struct {
	MaxScheduleTuples *int            `json:"maxScheduleTuples,omitempty" validate:"omitempty,gte=0"`
	EvseID            int             `json:"evseId" validate:"gt=0"`
	Timestamp         *types.DateTime `json:"timestamp,omitempty" validate:"omitempty"`
	ChargingNeeds     ChargingNeeds   `json:"chargingNeeds" validate:"required"`
}

// This field definition of the NotifyEVChargingNeeds response payload, sent by the CSMS to the Charging Station in response to a NotifyEVChargingNeedsRequest.
// In case the request was invalid, or couldn't be processed, an error will be sent instead.
type NotifyEVChargingNeedsResponse struct {
	// Returns whether the CSMS has been able to process the message successfully.
	// It does not imply that the evChargingNeeds can be met with the current charging profile.
	Status     EVChargingNeedsStatus `json:"status" validate:"required,evChargingNeedsStatus21"`
	StatusInfo *types.StatusInfo     `json:"statusInfo,omitempty" validate:"omitempty,dive"` // Detailed status information.
}

// When an EV sends a ChargeParameterDiscoveryReq with with charging needs parameters,
// the Charging Station sends this information in a NotifyEVChargingNeedsRequest to the CSMS.
// The CSMS replies to the Charging Station with a NotifyEVChargingNeedsResponse message.
//
// The CSMS will re-calculate a new charging schedule, trying to accomodate the EV needs,
// then asynchronously send a SetChargingProfileRequest with the new schedule to the Charging Station.
type NotifyEVChargingNeedsFeature struct{}

func (f NotifyEVChargingNeedsFeature) GetFeatureName() string {
	return NotifyEVChargingNeedsFeatureName
}

func (f NotifyEVChargingNeedsFeature) GetRequestType() reflect.Type {
	return reflect.TypeOf(NotifyEVChargingNeedsRequest{})
}

func (f NotifyEVChargingNeedsFeature) GetResponseType() reflect.Type {
	return reflect.TypeOf(NotifyEVChargingNeedsResponse{})
}

func (r NotifyEVChargingNeedsRequest) GetFeatureName() string {
	return NotifyEVChargingNeedsFeatureName
}

func (c NotifyEVChargingNeedsResponse) GetFeatureName() string {
	return NotifyEVChargingNeedsFeatureName
}

// Creates a new NotifyEVChargingNeedsRequest, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEVChargingNeedsRequest(evseID int, chargingNeeds ChargingNeeds) *NotifyEVChargingNeedsRequest {
	return &NotifyEVChargingNeedsRequest{EvseID: evseID, ChargingNeeds: chargingNeeds}
}

// Creates a new NotifyEVChargingNeedsResponse, containing all required fields. Optional fields may be set afterwards.
func NewNotifyEVChargingNeedsResponse(status EVChargingNeedsStatus) *NotifyEVChargingNeedsResponse {
	return &NotifyEVChargingNeedsResponse{Status: status}
}

func init() {
	_ = types.Validate.RegisterValidation("energyTransferMode21", isValidEnergyTransferMode)
	_ = types.Validate.RegisterValidation("evChargingNeedsStatus21", isValidEVChargingNeedsStatus)
}
