package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestDateTime struct {
	Timestamp DateTime `json:"timestamp"`
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	jsonMessage := json.RawMessage(`{"timestamp": "8/8/1965 13:00:00 PM"}`)
	var data TestDateTime
	json.Unmarshal(jsonMessage, &data)
	date := time.Date(1965, 8, 8, 13, 0, 0, 0, time.UTC)
	assert.True(t, data.Timestamp.Equal(date))

}
