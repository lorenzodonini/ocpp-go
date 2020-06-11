package types

import (
	"encoding/json"
	"strings"
	"time"
)

// DateTimeFormat to be used for all OCPP messages.
//
// The default dateTime format is RFC3339.
// Change this if another format is desired.
var DateTimeFormat = time.RFC3339

// DateTime wraps a time.Time struct, allowing for improved dateTime JSON compatibility.
type DateTime struct {
	time.Time
}

// Creates a new DateTime struct, embedding a time.Time struct.
func NewDateTime(time time.Time) *DateTime {
	return &DateTime{Time: time}
}

func (dt *DateTime) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	if DateTimeFormat == "" {
		defaultTime := time.Time{}
		err := json.Unmarshal(input, defaultTime)
		if err != nil {
			return err
		}
		dt.Time = defaultTime.Local()
	} else {
		newTime, err := time.Parse(DateTimeFormat, strInput)
		if err != nil {
			return err
		}
		dt.Time = newTime.Local()
	}
	return nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	if DateTimeFormat == "" {
		return json.Marshal(dt.Time)
	}
	timeStr := FormatTimestamp(dt.Time)
	return json.Marshal(timeStr)
}

// Formats the UTC timestamp using the DateTimeFormat setting.
// This function is used during JSON marshaling as well.
func (dt *DateTime) FormatTimestamp() string {
	return dt.UTC().Format(DateTimeFormat)
}

func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(DateTimeFormat)
}

// DateTime Validation

func DateTimeIsNull(dateTime *DateTime) bool {
	return dateTime != nil && dateTime.IsZero()
}

func validateDateTimeGt(dateTime *DateTime, than time.Time) bool {
	return dateTime != nil && dateTime.After(than)
}

func validateDateTimeNow(dateTime DateTime) bool {
	dur := time.Now().Sub(dateTime.Time).Minutes()
	return dur < 1
}

func validateDateTimeLt(dateTime DateTime, than time.Time) bool {
	return dateTime.Before(than)
}
