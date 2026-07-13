package types

import (
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"gopkg.in/go-playground/validator.v9"
)

type Tariff struct {
	TariffId         string           `json:"tariffId" validate:"required,max=60"` // Identifier used to identify one tariff.
	Currency         string           `json:"currency" validate:"required,max=3"`
	ValidFrom        *DateTime        `json:"validFrom,omitempty" validate:"omitempty"`
	Description      []MessageContent `json:"description,omitempty" validate:"omitempty,max=10,dive"`
	Energy           *TariffEnergy    `json:"energy,omitempty" validate:"omitempty,dive"`
	ChargingTime     *TariffTime      `json:"chargingTime,omitempty" validate:"omitempty,dive"`
	IdleTime         *TariffTime      `json:"idleTime,omitempty" validate:"omitempty,dive"`
	FixedFee         *TariffFixed     `json:"fixedFee,omitempty" validate:"omitempty,dive"`
	MinCost          *Price           `json:"minCost,omitempty" validate:"omitempty,dive"`
	MaxCost          *Price           `json:"maxCost,omitempty" validate:"omitempty,dive"`
	ReservationTime  *TariffTime      `json:"reservationTime,omitempty" validate:"omitempty,dive"`
	ReservationFixed *TariffFixed     `json:"reservationFixed,omitempty" validate:"omitempty,dive"`
}

type Price struct {
	ExclTax  *float64  `json:"exclTax,omitempty" validate:"omitempty"`
	InclTax  *float64  `json:"inclTax,omitempty" validate:"omitempty"`
	TaxRates []TaxRate `json:"taxRates,omitempty" validate:"omitempty,max=5,dive"`
}

type TaxRate struct {
	Type  string  `json:"type" validate:"required,max=20"` // Type of tax rate, e.g., VAT.
	Tax   float64 `json:"tax" validate:"required"`
	Stack *int    `json:"stack,omitempty" validate:"omitempty,gte=0"`
}

type TaxRule struct {
	TaxRuleId                   int            `json:"taxRuleId" validate:"required,gte=0"`
	TaxRuleName                 *string        `json:"taxRuleName,omitempty" validate:"omitempty,max=100"`
	TaxIncludedInPrice          bool           `json:"taxIncludedInPrice,omitempty" validate:"omitempty"`
	AppliesToEnergyFee          bool           `json:"appliesToEnergyFee" validate:"required"`
	AppliesToParkingFee         bool           `json:"appliesToParkingFee" validate:"required"`
	AppliesToOverstayFee        bool           `json:"appliesToOverstayFee" validate:"required"`
	AppliesToMinimumMaximumCost bool           `json:"appliesToMinimumMaximumCost" validate:"required"`
	TaxRate                     RationalNumber `json:"taxRate" validate:"required,dive"` // Tax rate as a rational number.
}

type RationalNumber struct {
	Exponent int `json:"exponent" validate:"required"`
	Value    int `json:"value" validate:"required"`
}

type TariffTime struct {
	Prices   []TariffTimePrice `json:"prices" validate:"required,min=1,max=5,dive"`
	TaxRates []TaxRate         `json:"taxRates,omitempty" validate:"omitempty,max=5,dive"`
}

type TariffTimePrice struct {
	PriceMinute float64           `json:"priceMinute" validate:"required"` // Price per minute.
	Conditions  *TariffConditions `json:"conditions,omitempty" validate:"omitempty,dive"`
}

type TariffFixed struct {
	Prices   []TariffFixedPrice `json:"prices" validate:"required,min=1,dive"` // Prices for fixed fees.
	TaxRates []TaxRate          `json:"taxRates,omitempty" validate:"omitempty,max=5,dive"`
}

type TariffFixedPrice struct {
	PriceFixed float64                `json:"priceFixed" validate:"required"` // Fixed price.
	Conditions *TariffFixedConditions `json:"conditions,omitempty" validate:"omitempty,dive"`
}

type TariffFixedConditions struct {
	StartTimeOfDay     *string     `json:"startTimeOfDay,omitempty" validate:"omitempty"`
	EndTimeOfDay       *string     `json:"endTimeOfDay,omitempty" validate:"omitempty"`
	DayOfWeek          []DayOfWeek `json:"dayOfWeek,omitempty" validate:"omitempty,dayOfWeek"`
	ValidFromDate      string      `json:"validFromDate,omitempty" validate:"omitempty"`
	ValidToDate        string      `json:"validToDate,omitempty" validate:"omitempty"`
	EvseKind           *EvseKind   `json:"evseKind,omitempty" validate:"omitempty,evseKind"`
	PaymentBrand       *string     `json:"paymentBrand,omitempty" validate:"omitempty,max=20"`
	PaymentRecognition *string     `json:"paymentRecognition,omitempty" validate:"omitempty,max=20"`
}

type TariffEnergy struct {
	TaxRates []TaxRate           `json:"taxRates,omitempty" validate:"omitempty,max=5,dive"`
	Prices   []TariffEnergyPrice `json:"prices" validate:"required,min=1,dive"` // Prices for energy in kWh.
}

type TariffEnergyPrice struct {
	PriceKwh   float64           `json:"priceKWh" validate:"required"` // Price per kWh.
	Conditions *TariffConditions `json:"conditions,omitempty" validate:"omitempty,dive"`
}

type EvseKind string

const (
	EvseKindAC EvseKind = "AC" // Alternating Current
	EvseKindDC EvseKind = "DC" // Direct Current
)

func isValidEvseKind(fl validator.FieldLevel) bool {
	switch EvseKind(fl.Field().String()) {
	case EvseKindAC, EvseKindDC:
		return true
	default:
		return false
	}
}

type DayOfWeek string

const (
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
)

func isValidDayOfWeek(fl validator.FieldLevel) bool {
	switch DayOfWeek(fl.Field().String()) {
	case DayOfWeekMonday, DayOfWeekTuesday, DayOfWeekWednesday,
		DayOfWeekThursday, DayOfWeekFriday, DayOfWeekSaturday, DayOfWeekSunday:
		return true
	default:
		return false
	}
}

func init() {
	_ = ocppj.Validate.RegisterValidation("evseKind", isValidEvseKind)
	_ = ocppj.Validate.RegisterValidation("dayOfWeek", isValidDayOfWeek)
}

type TariffConditions struct {
	StartTimeOfDay  *string     `json:"startTimeOfDay,omitempty" validate:"omitempty"`
	EndTimeOfDay    *string     `json:"endTimeOfDay,omitempty" validate:"omitempty"`
	DayOfWeek       []DayOfWeek `json:"dayOfWeek,omitempty" validate:"omitempty,dayOfWeek"`
	ValidFromDate   string      `json:"validFromDate,omitempty" validate:"omitempty"`
	ValidToDate     string      `json:"validToDate,omitempty" validate:"omitempty"`
	EvseKind        *EvseKind   `json:"evseKind,omitempty" validate:"omitempty,evseKind"`
	MinEnergy       *float64    `json:"minEnergy,omitempty" validate:"omitempty"`
	MaxEnergy       *float64    `json:"maxEnergy,omitempty" validate:"omitempty"`
	MinCurrent      *float64    `json:"minCurrent,omitempty" validate:"omitempty"`
	MaxCurrent      *float64    `json:"maxCurrent,omitempty" validate:"omitempty"`
	MinPower        *float64    `json:"minPower,omitempty" validate:"omitempty"`
	MaxPower        *float64    `json:"maxPower,omitempty" validate:"omitempty"`
	MinTime         *int        `json:"minTime,omitempty" validate:"omitempty"` // Minimum time in seconds.
	MaxTime         *int        `json:"maxTime,omitempty" validate:"omitempty"` // Maximum time in seconds.
	MinChargingTime *int        `json:"minChargingTime,omitempty" validate:"omitempty"`
	MaxChargingTime *int        `json:"maxChargingTime,omitempty" validate:"omitempty"` // Maximum charging time in seconds.
	MinIdleTime     *int        `json:"minIdleTime,omitempty" validate:"omitempty"`     // Minimum idle time in seconds.
	MaxIdleTime     *int        `json:"maxIdleTime,omitempty" validate:"omitempty"`     // Maximum idle time in seconds.
}
