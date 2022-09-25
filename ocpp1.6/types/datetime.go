package types

import (
	"encoding/json"
	"strings"
	"time"
)

// ISO8601 time format, assuming Zulu timestamp.
const ISO8601 = "2006-01-02T15:04:05Z"

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
		var defaultTime time.Time
		err := json.Unmarshal(input, &defaultTime)
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
	return dt.UTC().Format(DateTimeFormat)
}

func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(DateTimeFormat)
}

// DateTime Validation

func DateTimeIsNull(dateTime *DateTime) bool {
	return dateTime != nil && dateTime.IsZero()
}
