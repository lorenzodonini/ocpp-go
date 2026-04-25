package tests

import (
	"strings"
	"testing"

	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/stretchr/testify/assert"
)

type GenericTestEntry struct {
	Element       interface{}
	ExpectedValid bool
}

// TODO: pass expected error value for improved validation and error message
func ExecuteGenericTestTable(t *testing.T, testTable []GenericTestEntry) {
	for _, testCase := range testTable {
		err := ocppj.Validate.Struct(testCase.Element)
		if err != nil {
			assert.Equal(t, testCase.ExpectedValid, false, err.Error())
		} else {
			assert.Equal(t, testCase.ExpectedValid, true, "%v is valid", testCase.Element)
		}
	}
}

// Generates a new dummy string of the specified length.
func NewLongString(length int) string {
	reps := length / 32
	s := strings.Repeat("................................", reps)
	for i := len(s); i < length; i++ {
		s += "."
	}
	return s
}

func NewInt(i int) *int {
	return &i
}

func NewFloat(f float64) *float64 {
	return &f
}

func NewBool(b bool) *bool {
	return &b
}
