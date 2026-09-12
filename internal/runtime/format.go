package runtime

import (
	"strconv"
	"strings"
)

// FormatReal renders the recorded 15-significant-digit numeric text profile.
// It has no leading space and discards the sign of zero.
func FormatReal(value float64) string {
	if value == 0 {
		return "0"
	}
	text := strconv.FormatFloat(value, 'G', 15, 64)
	if mantissa, exponent, ok := strings.Cut(text, "E"); ok {
		n, _ := strconv.Atoi(exponent) // FormatFloat always emits a valid exponent.
		text = mantissa + "E" + strconv.Itoa(n)
	}
	return text
}
