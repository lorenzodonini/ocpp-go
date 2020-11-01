package types

import (
	"encoding/json"
	"github.com/araddon/dateparse"
	"strings"
	"time"
)

// ISO8601 time format, assuming Zulu timestamp.
const ISO8601 = "2006-01-02T15:04:05Z"

// The default dateTime format is RFC3339.
var DefaultTimeFormat = time.RFC3339

// DateTimeFormat to be used for all OCPP messages.
// If not specified DefaultTimeFormat is used
var DateTimeFormat = ""

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
		var stringTime string
		err := json.Unmarshal(input, &stringTime)
		if err != nil {
			return err
		}
		defaultTime, err := dateparse.ParseAny(stringTime)
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
	timeStr := dt.FormatTimestamp()
	return json.Marshal(timeStr)
}

// Formats the UTC timestamp using the DateTimeFormat setting.
// This function is used during JSON marshaling as well.
func (dt *DateTime) FormatTimestamp() string {
	if DateTimeFormat != "" {
		return dt.UTC().Format(DateTimeFormat)
	}
	return dt.UTC().Format(DefaultTimeFormat)

}

func FormatTimestamp(t time.Time) string {
	if DateTimeFormat != "" {
		return t.UTC().Format(DateTimeFormat)
	}
	return t.UTC().Format(DefaultTimeFormat)
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
