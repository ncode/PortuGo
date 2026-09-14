package main

import (
	"fmt"
	"math"
	"math/big"
	"slices"
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

// randomOutputExpectation qualifies the fixed framing around generated lines
// while leaving the reference generator sequence unspecified.
type randomOutputExpectation struct {
	Kind      string `json:"kind"`
	Lines     int    `json:"lines"`
	Minimum   int64  `json:"minimum,omitempty"`
	Bound     int64  `json:"bound"`
	FixedTail int64  `json:"fixedTail"`
	Review    review `json:"review"`
}

func (c randomOutputExpectation) compare(output []byte) error {
	switch c.Kind {
	case "randi-lines":
		if c.Lines < 1 || c.Lines > 4096 || c.Minimum != 0 || c.Bound < 1 || c.Bound > math.MaxInt32 || c.FixedTail < 0 || c.FixedTail >= c.Bound {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-int-text-lines":
		if c.Lines < 2 || c.Lines > 4096 || c.Lines%2 != 0 || c.Minimum != 0 || c.Bound < 1 || c.Bound > math.MaxInt32 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-int-sort-lines":
		if c.Lines < 2 || c.Lines > 4096 || c.Lines%2 != 0 || c.Minimum < math.MinInt32 || c.Minimum >= c.Bound || c.Bound < 1 || c.Bound > math.MaxInt32 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	default:
		return fmt.Errorf("invalid random-output contract")
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) != c.Lines+1 || lines[len(lines)-1] != "" {
		return fmt.Errorf("random-output line count mismatch")
	}
	if c.Kind == "random-int-text-lines" {
		for index, line := range lines[:c.Lines] {
			if index%2 == 1 {
				if len(line) != 5 || strings.Trim(line, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") != "" {
					return fmt.Errorf("random-output text line %d mismatch", index+1)
				}
				continue
			}
			value, err := strconv.ParseInt(line, 10, 64)
			if err != nil || line != strconv.FormatInt(value, 10) || value < 0 || value >= c.Bound {
				return fmt.Errorf("random-output integer line %d outside domain", index+1)
			}
		}
		return nil
	}
	if c.Kind == "random-int-sort-lines" {
		input := make([]int64, c.Lines/2)
		for index, line := range lines[:c.Lines/2] {
			value, err := strconv.ParseInt(line, 10, 64)
			if err != nil || line != strconv.FormatInt(value, 10) || value < c.Minimum || value >= c.Bound {
				return fmt.Errorf("random-output integer line %d outside domain", index+1)
			}
			input[index] = value
		}
		expected := slices.Clone(input)
		slices.Sort(expected)
		for index, line := range lines[c.Lines/2 : c.Lines] {
			if !strings.HasPrefix(line, " ") {
				return fmt.Errorf("random-output sorted line %d spacing mismatch", index+1)
			}
			value, err := strconv.ParseInt(line[1:], 10, 64)
			if err != nil || line[1:] != strconv.FormatInt(value, 10) || value < c.Minimum || value >= c.Bound {
				return fmt.Errorf("random-output sorted line %d outside domain", index+1)
			}
			if value != expected[index] {
				return fmt.Errorf("random-output sorted line %d is not the input permutation", index+1)
			}
		}
		return nil
	}
	for index, line := range lines[:c.Lines] {
		if !strings.HasPrefix(line, " ") {
			return fmt.Errorf("random-output line %d spacing mismatch", index+1)
		}
		value, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil || line[1:] != strconv.FormatInt(value, 10) || value < 0 || value >= c.Bound {
			return fmt.Errorf("random-output line %d outside domain", index+1)
		}
		if index == c.Lines-1 && value != c.FixedTail {
			return fmt.Errorf("random-output fixed tail mismatch")
		}
	}
	return nil
}
