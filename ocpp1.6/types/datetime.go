package types

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/relvacode/iso8601"
)

// DateTimeFormat to be used when serializing all OCPP messages.
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

// Creates a new DateTime struct, containing a time.Now() value.
func Now() *DateTime {
	return &DateTime{Time: time.Now()}
}

func null(b []byte) bool {
	if len(b) != 4 {
		return false
	}
	if b[0] != 'n' && b[1] != 'u' && b[2] != 'l' && b[3] != 'l' {
		return false
	}
	return true
}

func (dt *DateTime) UnmarshalJSON(input []byte) error {
	// Do not parse null timestamps
	if null(input) {
		return nil
	}
	// Assert that timestamp is a string
	if len(input) > 0 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	} else {
		return errors.New("timestamp not enclosed in double quotes")
	}
	// Parse ISO8601
	var err error
	dt.Time, err = iso8601.Parse(input)
	return err
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
