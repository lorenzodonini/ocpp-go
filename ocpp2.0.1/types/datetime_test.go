package types

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/relvacode/iso8601"
	"github.com/stretchr/testify/suite"
)

type dateTimeSuite struct {
	suite.Suite
}

func (suite *dateTimeSuite) TestUnmarshalDateTime() {
	testTable := []struct {
		RawDateTime   string
		ExpectedValid bool
		ExpectedTime  time.Time
		ExpectedError error
	}{
		{"\"2019-03-01T10:00:00Z\"", true, time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00+01:00\"", true, time.Date(2019, 3, 1, 9, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00.000Z\"", true, time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00.000+01:00\"", true, time.Date(2019, 3, 1, 9, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00\"", true, time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00+01\"", true, time.Date(2019, 3, 1, 9, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00.000\"", true, time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01T10:00:00.000+01\"", true, time.Date(2019, 3, 1, 9, 0, 0, 0, time.UTC), nil},
		{"\"2019-03-01 10:00:00+00:00\"", false, time.Time{}, &iso8601.UnexpectedCharacterError{Character: ' '}},
		{"\"null\"", false, time.Time{}, &iso8601.UnexpectedCharacterError{Character: 110}},
		{"\"\"", false, time.Time{}, &iso8601.RangeError{Element: "month", Min: 1, Max: 12}},
		{"null", true, time.Time{}, nil},
	}
	for _, dt := range testTable {
		jsonStr := []byte(dt.RawDateTime)
		var dateTime DateTime
		err := json.Unmarshal(jsonStr, &dateTime)
		if dt.ExpectedValid {
			suite.NoError(err)
			suite.NotNil(dateTime)
			suite.True(dt.ExpectedTime.Equal(dateTime.Time))
		} else {
			suite.Error(err)
			suite.ErrorAs(err, &dt.ExpectedError)
		}
	}
}

func (suite *dateTimeSuite) TestMarshalDateTime() {
	testTable := []struct {
		Time                    time.Time
		Format                  string
		ExpectedFormattedString string
	}{
		{time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), "", "2019-03-01T10:00:00Z"},
		{time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), time.RFC3339, "2019-03-01T10:00:00Z"},
		{time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), time.RFC822, "01 Mar 19 10:00 UTC"},
		{time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), time.RFC1123, "Fri, 01 Mar 2019 10:00:00 UTC"},
		{time.Date(2019, 3, 1, 10, 0, 0, 0, time.UTC), "invalidFormat", "invalidFormat"},
	}
	for _, dt := range testTable {
		dateTime := NewDateTime(dt.Time)
		DateTimeFormat = dt.Format
		rawJson, err := dateTime.MarshalJSON()
		suite.NoError(err)
		formatted := strings.Trim(string(rawJson), "\"")
		suite.Equal(dt.ExpectedFormattedString, formatted)
	}
}

func (suite *dateTimeSuite) TestNowDateTime() {
	now := Now()
	suite.NotNil(now)
	suite.True(time.Now().Sub(now.Time) < 1*time.Second)
}

func TestDateTime(t *testing.T) {
	suite.Run(t, new(dateTimeSuite))
}
