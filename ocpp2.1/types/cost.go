package types

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/types"
	"gopkg.in/go-playground/validator.v9"
)

// Sales Tariff

// The kind of cost referred to in a CostType.
type CostKind string

const (
	CostKindCarbonDioxideEmission         CostKind = "CarbonDioxideEmission"         // Carbon Dioxide emissions, in grams per kWh.
	CostKindRelativePricePercentage       CostKind = "RelativePricePercentage"       // Price per kWh, as percentage relative to the maximum price stated in any of all tariffs indicated to the EV.
	CostKindRenewableGenerationPercentage CostKind = "RenewableGenerationPercentage" // Percentage of renewable generation within total generation.
)

func isValidCostKind(fl validator.FieldLevel) bool {
	purposeType := CostKind(fl.Field().String())
	switch purposeType {
	case CostKindCarbonDioxideEmission, CostKindRelativePricePercentage, CostKindRenewableGenerationPercentage:
		return true
	default:
		return false
	}
}

// Defines the time interval the SalesTariffEntry is valid for, based upon relative times.
type RelativeTimeInterval struct {
	Start    int  `json:"start"`                                         // Start of the interval, in seconds from NOW.
	Duration *int `json:"duration,omitempty" validate:"omitempty,gte=0"` // Duration of the interval, in seconds.
}

// Cost details.
type CostType struct {
	CostKind         CostKind `json:"costKind" validate:"required,costKind21"`                      // The kind of cost referred to in the message element amount.
	Amount           int      `json:"amount" validate:"gte=0"`                                      // The estimated or actual cost per kWh.
	AmountMultiplier *int     `json:"amountMultiplier,omitempty" validate:"omitempty,min=-3,max=3"` // The exponent to base 10 (dec).
}

// Contains price information and/or alternative costs.
type ConsumptionCost struct {
	StartValue float64    `json:"startValue"`                          // The lowest level of consumption that defines the starting point of this consumption block
	Cost       []CostType `json:"cost" validate:"required,max=3,dive"` // Contains the cost details.
}

// NewConsumptionCost instantiates a new ConsumptionCost struct. No additional parameters need to be set.
func NewConsumptionCost(startValue float64, cost []CostType) ConsumptionCost {
	return ConsumptionCost{
		StartValue: startValue,
		Cost:       cost,
	}
}

// Element describing all relevant details for one time interval of the SalesTariff.
type SalesTariffEntry struct {
	EPriceLevel          *int                 `json:"ePriceLevel,omitempty" validate:"omitempty,gte=0"`          // The price level of this SalesTariffEntry (referring to NumEPriceLevels). Small values for the EPriceLevel represent a cheaper TariffEntry.
	RelativeTimeInterval RelativeTimeInterval `json:"relativeTimeInterval" validate:"required"`                  // The time interval the SalesTariffEntry is valid for, based upon relative times.
	ConsumptionCost      []ConsumptionCost    `json:"consumptionCost,omitempty" validate:"omitempty,max=3,dive"` // Additional means for further relative price information and/or alternative costs.
}

// Sales tariff associated with this charging schedule.
type SalesTariff struct {
	ID                     int                `json:"id"`                                                           // Identifier used to identify one sales tariff.
	SalesTariffDescription string             `json:"salesTariffDescription,omitempty" validate:"omitempty,max=32"` // A human readable title/short description of the sales tariff e.g. for HMI display purposes.
	NumEPriceLevels        *int               `json:"numEPriceLevels,omitempty" validate:"omitempty"`               // Defines the overall number of distinct price levels used across all provided SalesTariff elements.
	SalesTariffEntry       []SalesTariffEntry `json:"salesTariffEntry" validate:"required,min=1,max=1024,dive"`     // Encapsulates elements describing all relevant details for one time interval of the SalesTariff.
}

// NewSalesTariff instantiates a new SalesTariff struct. Only required fields are passed as parameters.
func NewSalesTariff(id int, salesTariffEntries []SalesTariffEntry) *SalesTariff {
	return &SalesTariff{
		ID:               id,
		SalesTariffEntry: salesTariffEntries,
	}
}

type TariffCost string

const (
	TariffCostNormal  = "NormalCost"
	TariffCostMinimum = "MinCost"
	TariffCostMaximum = "MaxCost"
)

func isValidTariffCost(fl validator.FieldLevel) bool {
	costType := TariffCost(fl.Field().String())
	switch costType {
	case TariffCostNormal, TariffCostMinimum, TariffCostMaximum:
		return true
	default:
		return false
	}
}

type CostDimensionType string

const (
	CostDimensionTypeEnergy       CostDimensionType = "Energy"       // Cost dimension for energy consumed.
	CostDimensionTypeChargingTime CostDimensionType = "ChargingTime" // Cost dimension for the time the EV was charging.
	CostDimensionTypeIdleTime     CostDimensionType = "IdleTime"     // Cost dimension for the time the EV was idle.
	CostDimensionMinCurrent       CostDimensionType = "MinCurrent"   // Cost dimension for the minimum current used during the charging session.
	CostDimensionMaxCurrent       CostDimensionType = "MaxCurrent"   // Cost dimension for the maximum current used during the charging session.
	CostDimensionMinPower         CostDimensionType = "MinPower"     // Cost dimension for the minimum power used during the charging session.
	CostDimensionMaxPower         CostDimensionType = "MaxPower"     // Cost dimension for the maximum power used during the charging session.
)

func isValidCostDimensionType(fl validator.FieldLevel) bool {
	dimensionType := CostDimensionType(fl.Field().String())
	switch dimensionType {
	case CostDimensionTypeEnergy, CostDimensionTypeChargingTime, CostDimensionTypeIdleTime,
		CostDimensionMinCurrent, CostDimensionMaxCurrent, CostDimensionMinPower, CostDimensionMaxPower:
		return true
	default:
		return false
	}
}

type CostDetails struct {
	FailureToCalculate *bool            `json:"failureToCalculate,omitempty"`
	FailureReason      *string          `json:"failureReason,omitempty" validate:"omitempty,max=500"`
	ChargingPeriods    []ChargingPeriod `json:"chargingPeriods,omitempty" validate:"omitempty,dive"`
	TotalCost          TotalCost        `json:"totalCost" validate:"required,dive"`
	TotalUsage         TotalUsage       `json:"totalUsage" validate:"required,dive"` // Total usage of the EV during the charging session.
}

type TotalCost struct {
	Currency         string     `json:"currency" validate:"required,max=3"`          // The currency of the total cost.
	TypeOfCost       TariffCost `json:"typeOfCost" validate:"required,tariffCost21"` // The type of cost, e.g. NormalCost, MinCost, MaxCost.
	Fixed            *Price     `json:"fixedPrice,omitempty" validate:"omitempty,dive"`
	Energy           *Price     `json:"energy,omitempty" validate:"omitempty,dive"`
	ChargingTime     *Price     `json:"chargingTime,omitempty" validate:"omitempty,dive"`
	IdleTime         *Price     `json:"idleTime,omitempty" validate:"omitempty,dive"`
	ReservationTime  *Price     `json:"reservationTime,omitempty" validate:"omitempty,dive"`
	Total            TotalPrice `json:"total" validate:"required,dive"` // Total cost of the charging session.
	ReservationFixed *Price     `json:"reservationFixed,omitempty" validate:"omitempty,dive"`
}

type TotalPrice struct {
	ExclTax *float64 `json:"exclTax,omitempty" validate:"omitempty"` // Total price excluding tax.
	InclTax *float64 `json:"inclTax,omitempty" validate:"omitempty"` // Total price including tax.
}

type TotalUsage struct {
	Energy          float64 `json:"energy" validate:"required"`
	ChargingTime    int     `json:"chargingTime" validate:"required"`
	IdleTime        int     `json:"idleTime" validate:"required"`                   // Total idle time of the EV during the charging session, in seconds.
	ReservationTime *int    `json:"reservationTime,omitempty" validate:"omitempty"` // Total reservation time of the EV during the charging session, in seconds.
}

type ChargingPeriod struct {
	TariffId    *string         `json:"tariffId,omitempty" validate:"omitempty,max=60"` // The ID of the tariff used for this charging period.
	StartPeriod types.DateTime  `json:"startPeriod" validate:"required"`                // The start of the charging period.
	Dimensions  []CostDimension `json:"dimensions,omitempty" validate:"omitempty,dive"`
}

type CostDimension struct {
	Type   CostDimensionType `json:"type" validate:"required,costDimension21"`
	Volume float64           `json:"volume" validate:"required"`
}
