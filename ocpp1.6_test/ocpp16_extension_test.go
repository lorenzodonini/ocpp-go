package ocpp16_test

import (
	"strings"
)

// Generates a new dummy string of the specified length.
func newLongString(length int) string {
	reps := length / 32
	s := strings.Repeat("................................", reps)
	for i := len(s); i < length; i++ {
		s += "."
	}
	return s
}
