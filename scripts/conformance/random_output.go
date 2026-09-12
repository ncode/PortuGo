package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// randomInputExpectation qualifies only repeated leia/escreval pairs. Numeric
// endpoints are inclusive integer ticks at 10^-Decimals, not sampled extrema.
// Every other output channel and every unqualified probe remains byte-exact.
type randomInputExpectation struct {
	Kind     string `json:"kind"`
	Minimum  int64  `json:"minimum"`
	Maximum  int64  `json:"maximum"`
	Decimals int    `json:"decimals"`
	Samples  int    `json:"samples"`
	Review   review `json:"review"`
}

func (c randomInputExpectation) compare(output []byte) error {
	if c.Samples < 1 || c.Samples > 4096 || c.Decimals < 0 || c.Decimals > 5 || c.Minimum > c.Maximum {
		return fmt.Errorf("invalid random-input contract")
	}
	scale := int64(math.Pow10(c.Decimals))
	if c.Minimum < math.MinInt32*scale || c.Maximum > math.MaxInt32*scale {
		return fmt.Errorf("random-input contract exceeds qualified numeric range")
	}
	switch c.Kind {
	case "integer":
		if c.Decimals != 0 {
			return fmt.Errorf("integer random-input contract has fractional ticks")
		}
	case "real":
	case "text":
		if c.Decimals != 0 || c.Minimum != 0 || c.Maximum != 0 {
			return fmt.Errorf("text random-input contract has numeric bounds")
		}
	default:
		return fmt.Errorf("unknown random-input contract kind")
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) != 2*c.Samples+1 || lines[len(lines)-1] != "" {
		return fmt.Errorf("random-input echo/output pair count mismatch")
	}
	for sample := range c.Samples {
		echo, printed := lines[2*sample], lines[2*sample+1]
		if c.Kind == "text" {
			if len(echo) != 5 || strings.Trim(echo, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") != "" || printed != echo {
				return fmt.Errorf("random-input text sample %d mismatch", sample+1)
			}
			continue
		}
		// Parse and check canonical spelling before using exact rational arithmetic;
		// this also rejects fractions and unbounded exponent spellings from output.
		n, err := strconv.ParseFloat(echo, 64)
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) || n < math.MinInt32 || n > math.MaxInt32 {
			return fmt.Errorf("invalid numeric random-input sample %d", sample+1)
		}
		if n == 0 {
			n = 0 // The recorded input profile discards the sign of zero.
		}
		precision := 10
		if c.Kind == "integer" {
			precision = 0
		}
		if echo != strconv.FormatFloat(n, 'f', precision, 64) {
			return fmt.Errorf("random-input sample %d echo formatting mismatch", sample+1)
		}
		value, ok := new(big.Rat).SetString(echo)
		if !ok {
			return fmt.Errorf("invalid numeric random-input sample %d", sample+1)
		}
		ticks := new(big.Rat).Mul(value, new(big.Rat).SetInt64(scale))
		if !ticks.IsInt() || !ticks.Num().IsInt64() || ticks.Num().Int64() < c.Minimum || ticks.Num().Int64() > c.Maximum {
			return fmt.Errorf("random-input sample %d outside domain or precision grid", sample+1)
		}
		if c.Kind == "integer" {
			canonical := strconv.FormatInt(ticks.Num().Int64(), 10)
			if echo != canonical || printed != " "+canonical {
				return fmt.Errorf("random-input integer sample %d formatting mismatch", sample+1)
			}
			continue
		}
		canonical := strconv.FormatFloat(n, 'g', 15, 64)
		if mantissa, exponent, found := strings.Cut(canonical, "e"); found {
			power, err := strconv.Atoi(exponent)
			if err != nil {
				return err
			}
			canonical = mantissa + "E" + strconv.Itoa(power)
		}
		if printed != " "+canonical {
			return fmt.Errorf("random-input real sample %d formatting mismatch", sample+1)
		}
	}
	return nil
}
