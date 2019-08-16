package ocpp16

import (
	"encoding/json"
	"github.com/lorenzodonini/go-ocpp/ocppj"
	"gopkg.in/go-playground/validator.v9"
	"strings"
	"time"
)

const (
	ISO8601 = "2006-01-02T15:04:05Z"
)

type DateTime struct {
	time.Time
}

func NewDateTime(time time.Time) *DateTime {
	return &DateTime{Time: time}
}

var DateTimeFormat = ISO8601

func (dt *DateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if DateTimeFormat == "" {
		defaultTime := time.Time{}
		err := json.Unmarshal(input, defaultTime)
		if err != nil {
			return err
		}
		dt.Time = defaultTime
	} else {
		newTime, err := time.Parse(DateTimeFormat, strInput)
		if err != nil {
			return err
		}
		dt.Time = newTime
	}
	return nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	if DateTimeFormat == "" {
		return json.Marshal(dt.Time)
	}
	timeStr := dt.Time.Format(DateTimeFormat)
	return json.Marshal(timeStr)
}

type PropertyViolation struct {
	error
	Property string
}

func (e *PropertyViolation) Error() string {
	return ""
}

type AuthorizationStatus string

const (
	AuthorizationStatusAccepted     AuthorizationStatus = "Accepted"
	AuthorizationStatusBlocked      AuthorizationStatus = "Blocked"
	AuthorizationStatusExpired      AuthorizationStatus = "Expired"
	AuthorizationStatusInvalid      AuthorizationStatus = "Invalid"
	AuthorizationStatusConcurrentTx AuthorizationStatus = "ConcurrentTx"
)

func isValidAuthorizationStatus(fl validator.FieldLevel) bool {
	status := AuthorizationStatus(fl.Field().String())
	switch status {
	case AuthorizationStatusAccepted, AuthorizationStatusBlocked, AuthorizationStatusExpired, AuthorizationStatusInvalid, AuthorizationStatusConcurrentTx:
		return true
	default:
		return false
	}
}

type IdTagInfo struct {
	ExpiryDate  DateTime            `json:"expiryDate" validate:"omitempty,gt"`
	ParentIdTag string              `json:"parentIdTag" validate:"omitempty,max=20"`
	Status      AuthorizationStatus `json:"status" validate:"required,authorizationStatus"`
}

func IdTagInfoStructLevelValidation(sl validator.StructLevel) {
	idTagInfo := sl.Current().Interface().(IdTagInfo)
	if !dateTimeIsNull(idTagInfo.ExpiryDate) && !validateDateTimeGt(idTagInfo.ExpiryDate, time.Now()) {
		sl.ReportError(idTagInfo.ExpiryDate, "ExpiryDate", "expiryDate", "gt", "")
	}
}

// Charging Profiles
type ChargingProfilePurposeType string
type ChargingProfileKindType string
type RecurrencyKindType string
type ChargingRateUnitType string

const (
	ChargingProfilePurposeChargePointMaxProfile ChargingProfilePurposeType = "ChargePointMaxProfile"
	ChargingProfilePurposeTxDefaultProfile      ChargingProfilePurposeType = "TxDefaultProfile"
	ChargingProfilePurposeTxProfile             ChargingProfilePurposeType = "TxProfile"
	ChargingProfileKindAbsolute                 ChargingProfileKindType    = "Absolute"
	ChargingProfileKindRecurring                ChargingProfileKindType    = "Recurring"
	ChargingProfileKindRelative                 ChargingProfileKindType    = "Relative"
	RecurrencyKindDaily                         RecurrencyKindType         = "Daily"
	RecurrencyKindWeekly                        RecurrencyKindType         = "Weekly"
	ChargingRateUnitWatts                       ChargingRateUnitType       = "W"
	ChargingRateUnitAmperes                     ChargingRateUnitType       = "A"
)

func isValidChargingProfilePurpose(fl validator.FieldLevel) bool {
	purposeType := ChargingProfilePurposeType(fl.Field().String())
	switch purposeType {
	case ChargingProfilePurposeChargePointMaxProfile, ChargingProfilePurposeTxDefaultProfile, ChargingProfilePurposeTxProfile:
		return true
	default:
		return false
	}
}

func isValidChargingProfileKind(fl validator.FieldLevel) bool {
	purposeType := ChargingProfileKindType(fl.Field().String())
	switch purposeType {
	case ChargingProfileKindAbsolute, ChargingProfileKindRecurring, ChargingProfileKindRelative:
		return true
	default:
		return false
	}
}

func isValidRecurrencyKind(fl validator.FieldLevel) bool {
	purposeType := RecurrencyKindType(fl.Field().String())
	switch purposeType {
	case RecurrencyKindDaily, RecurrencyKindWeekly:
		return true
	default:
		return false
	}
}

func isValidChargingRateUnit(fl validator.FieldLevel) bool {
	purposeType := ChargingRateUnitType(fl.Field().String())
	switch purposeType {
	case ChargingRateUnitWatts, ChargingRateUnitAmperes:
		return true
	default:
		return false
	}
}

type ChargingSchedulePeriod struct {
	StartPeriod  int     `json:"startPeriod" validate:"gte=0"`
	Limit        float64 `json:"limit" validate:"gte=0"`
	NumberPhases int     `json:"numberPhases,omitempty" validate:"gte=0"`
}

func NewChargingSchedulePeriod(startPeriod int, limit float64) ChargingSchedulePeriod {
	return ChargingSchedulePeriod{StartPeriod: startPeriod, Limit: limit}
}

type ChargingSchedule struct {
	Duration               int                      `json:"duration,omitempty" validate:"gte=0"`
	StartSchedule          *DateTime                `json:"startSchedule,omitempty"`
	ChargingRateUnit       ChargingRateUnitType     `json:"chargingRateUnit" validate:"required,chargingRateUnit"`
	ChargingSchedulePeriod []ChargingSchedulePeriod `json:"chargingSchedulePeriod" validate:"required,min=1"`
	MinChargingRate        float64                  `json:"minChargingRate,omitempty" validate:"gte=0"`
}

func NewChargingSchedule(chargingRateUnit ChargingRateUnitType, schedulePeriod ...ChargingSchedulePeriod) *ChargingSchedule {
	return &ChargingSchedule{ChargingRateUnit: chargingRateUnit, ChargingSchedulePeriod: schedulePeriod}
}

type ChargingProfile struct {
	ChargingProfileId      int                        `json:"chargingProfileId" validate:"gte=0"`
	TransactionId          int                        `json:"transactionId,omitempty"`
	StackLevel             int                        `json:"stackLevel" validate:"gt=0"`
	ChargingProfilePurpose ChargingProfilePurposeType `json:"chargingProfilePurpose" validate:"required,chargingProfilePurpose"`
	ChargingProfileKind    ChargingProfileKindType    `json:"chargingProfileKind" validate:"required,chargingProfileKind"`
	RecurrencyKind         RecurrencyKindType         `json:"recurrencyKind,omitempty" validate:"omitempty,recurrencyKind"`
	ValidFrom              *DateTime                  `json:"validFrom,omitempty"`
	ValidTo                *DateTime                  `json:"validTo,omitempty"`
	ChargingSchedule       *ChargingSchedule          `json:"chargingSchedule" validate:"required"`
}

func NewChargingProfile(chargingProfileId int, stackLevel int, chargingProfilePurpose ChargingProfilePurposeType, chargingProfileKind ChargingProfileKindType, schedule *ChargingSchedule) *ChargingProfile {
	return &ChargingProfile{ChargingProfileId: chargingProfileId, StackLevel: stackLevel, ChargingProfilePurpose: chargingProfilePurpose, ChargingProfileKind: chargingProfileKind, ChargingSchedule: schedule}
}

// DateTime Validation
func dateTimeIsNull(dateTime DateTime) bool {
	return dateTime.IsZero()
}

func validateDateTimeGt(dateTime DateTime, than time.Time) bool {
	return dateTime.After(than)
}

func validateDateTimeNow(dateTime DateTime) bool {
	dur := time.Now().Sub(dateTime.Time).Minutes()
	return dur < 1
}

func validateDateTimeLt(dateTime DateTime, than time.Time) bool {
	return dateTime.Before(than)
}

// Initialize validator
var Validate = ocppj.Validate

func init() {
	_ = Validate.RegisterValidation("authorizationStatus", isValidAuthorizationStatus)
	_ = Validate.RegisterValidation("chargingProfilePurpose", isValidChargingProfilePurpose)
	_ = Validate.RegisterValidation("chargingProfileKind", isValidChargingProfileKind)
	_ = Validate.RegisterValidation("recurrencyKind", isValidRecurrencyKind)
	_ = Validate.RegisterValidation("chargingRateUnit", isValidChargingRateUnit)
	Validate.RegisterStructValidation(IdTagInfoStructLevelValidation, IdTagInfo{})
}
