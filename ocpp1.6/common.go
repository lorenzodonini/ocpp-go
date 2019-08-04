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

func dateTimeIsNull(dateTime DateTime) bool {
	return dateTime.IsZero()
}

func validateDateTimeGt(dateTime DateTime, than time.Time) bool {
	return dateTime.After(than)
}

func validateDateTimeNow(dateTime DateTime) bool {
	dur := time.Now().Sub(dateTime.Time).Seconds()
	return dur < 1
}

func validateDateTimeLt(dateTime DateTime, than time.Time) bool {
	return dateTime.Before(than)
}

var Validate = ocppj.Validate

func init() {
	_ = Validate.RegisterValidation("authorizationStatus", isValidAuthorizationStatus)
	Validate.RegisterStructValidation(IdTagInfoStructLevelValidation, IdTagInfo{})
}
