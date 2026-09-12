package stdlib

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ncode/portugol-go/internal/runtime"
)

func caracpnum(args []runtime.Value) (runtime.Value, bool, error) {
	if len(args) != 1 || args[0].Kind != runtime.StringValue {
		return runtime.Value{}, true, fmt.Errorf("caracpnum expects one caractere argument")
	}
	text := strings.Trim(args[0].Str, " ")
	integer := runtime.Value{Kind: runtime.IntegerValue}
	prefix := text
	if strings.HasPrefix(prefix, "+") || strings.HasPrefix(prefix, "-") {
		prefix = prefix[1:]
	}
	hex := ""
	isHex := false
	if strings.HasPrefix(prefix, "$") {
		hex = prefix[1:]
		isHex = true
	} else if strings.HasPrefix(prefix, "0x") || strings.HasPrefix(prefix, "0X") {
		hex = prefix[2:]
		isHex = true
	}
	if isHex {
		if hex == "" || !strings.ContainsRune("0123456789abcdefABCDEF", rune(hex[0])) {
			return integer, true, nil
		}
		n, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return runtime.Value{}, true, fmt.Errorf("cannot convert text to a number")
		}
		if n <= math.MaxInt32 && !strings.HasPrefix(text, "-") {
			integer.Int = int64(n)
		}
		return integer, true, nil
	}
	// A zero integral prefix takes the integer-zero fallback before real parsing.
	prefix = strings.TrimLeft(prefix, "0")
	if prefix == "" || prefix[0] < '0' || prefix[0] > '9' {
		return integer, true, nil
	}
	if strings.ContainsRune(text, '_') {
		return runtime.Value{}, true, fmt.Errorf("cannot convert text to a number")
	}
	realSyntax := strings.ContainsAny(text, ".,eE")
	normalized := strings.ReplaceAll(text, ",", ".")
	if strings.HasSuffix(normalized, "e+") || strings.HasSuffix(normalized, "E+") || strings.HasSuffix(normalized, "e-") || strings.HasSuffix(normalized, "E-") {
		normalized = normalized[:len(normalized)-1]
	}
	if strings.HasSuffix(normalized, "e") || strings.HasSuffix(normalized, "E") {
		normalized = normalized[:len(normalized)-1]
	}
	n, err := strconv.ParseFloat(normalized, 64)
	if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
		return runtime.Value{}, true, fmt.Errorf("cannot convert text to a number")
	}
	// The recorded integer spelling switches to real at this decimal boundary.
	if !realSyntax && n >= math.MinInt32 && n < 2147483650 {
		if value, err := strconv.ParseInt(text, 10, 32); err == nil && value >= 0 {
			integer.Int = value
		}
		return integer, true, nil
	}
	return runtime.Value{Kind: runtime.RealValue, Real: n}, true, nil
}
