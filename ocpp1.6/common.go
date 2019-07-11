package ocpp16

import (
	"encoding/json"
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

type IdTagInfo struct {
	ExpiryDate  time.Time           `json:"expiryDate" validate:"omitempty,gt"`
	ParentIdTag string              `json:"parentIdTag" validate:"omitempty,max=20"`
	Status      AuthorizationStatus `json:"status" validate:"required"`
}
