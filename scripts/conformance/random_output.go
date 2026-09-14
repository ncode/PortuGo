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
	Decimals  int    `json:"decimals,omitempty"`
	FixedTail int64  `json:"fixedTail"`
	Review    review `json:"review"`
}

func (c randomOutputExpectation) compare(output []byte) error {
	switch c.Kind {
	case "randi-lines":
		if c.Lines < 1 || c.Lines > 4096 || c.Minimum != 0 || c.Decimals != 0 || c.Bound < 1 || c.Bound > math.MaxInt32 || c.FixedTail < 0 || c.FixedTail >= c.Bound {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-int-text-lines":
		if c.Lines < 2 || c.Lines > 4096 || c.Lines%2 != 0 || c.Minimum != 0 || c.Decimals != 0 || c.Bound < 1 || c.Bound > math.MaxInt32 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-int-sort-lines":
		if c.Lines < 2 || c.Lines > 4096 || c.Lines%2 != 0 || c.Minimum < math.MinInt32 || c.Minimum >= c.Bound || c.Bound < 1 || c.Bound > math.MaxInt32 || c.Decimals != 0 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-real-sort-lines":
		if c.Lines < 2 || c.Lines > 4096 || c.Lines%2 != 0 || c.Decimals < 1 || c.Decimals > 5 || c.Minimum < math.MinInt32 || c.Minimum >= c.Bound || c.Bound < 1 || c.Bound > math.MaxInt32 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-int-search-table-lines":
		if c.Lines != 12 || c.Minimum < math.MinInt32 || c.Minimum >= c.Bound || c.Bound < 1 || c.Bound > math.MaxInt32 || c.Decimals != 0 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-int-search-lines":
		if c.Lines != 21 || c.Minimum < math.MinInt32 || c.Minimum >= c.Bound || c.Bound < 1 || c.Bound > math.MaxInt32 || c.Decimals != 0 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-randi-repeat-lines":
		if c.Lines != 11 || c.Minimum < 1 || c.Minimum >= c.Bound || c.Bound < 2 || c.Bound > math.MaxInt32 || c.Decimals != 0 || c.FixedTail != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	case "random-record-sort-lines":
		if c.Lines != 42 || c.Minimum < 0 || c.Minimum >= c.Bound || c.Bound > math.MaxInt32 || c.FixedTail != 0 || c.Decimals != 0 {
			return fmt.Errorf("invalid random-output contract")
		}
	default:
		return fmt.Errorf("invalid random-output contract")
	}
	lines := strings.Split(string(output), "\n")
	if c.Kind == "random-randi-repeat-lines" {
		if len(lines) != c.Lines {
			return fmt.Errorf("random-output line count mismatch")
		}
	} else if len(lines) != c.Lines+1 || lines[len(lines)-1] != "" {
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
	if c.Kind == "random-real-sort-lines" {
		input := make([]int64, c.Lines/2)
		for index, line := range lines[:c.Lines/2] {
			value, err := parseFixedDecimal(line, c.Decimals, 10)
			if err != nil || value < c.Minimum || value >= c.Bound {
				return fmt.Errorf("random-output real line %d outside domain", index+1)
			}
			input[index] = value
		}
		expected := slices.Clone(input)
		slices.Sort(expected)
		for index, line := range lines[c.Lines/2 : c.Lines] {
			prefix := fmt.Sprintf("%3d - ", index+1)
			if !strings.HasPrefix(line, prefix) {
				return fmt.Errorf("random-output sorted real line %d framing mismatch", index+1)
			}
			valueText := strings.TrimSpace(line[len(prefix):])
			value, err := parseFixedDecimal(valueText, c.Decimals, c.Decimals)
			if err != nil || value < c.Minimum || value >= c.Bound || line != fmt.Sprintf("%s%10s", prefix, valueText) {
				return fmt.Errorf("random-output sorted real line %d formatting mismatch", index+1)
			}
			if value != expected[index] {
				return fmt.Errorf("random-output sorted real line %d is not the input permutation", index+1)
			}
		}
		return nil
	}
	if c.Kind == "random-randi-repeat-lines" {
		if lines[0] != " " || lines[1] != " ============================================== " || lines[2] != "QUANTOS NUMEROS (1-10): 1" || lines[3] != "Digite o destaque: 0" || lines[4] != " " || lines[5] != "A SEQUENCIA É " || lines[8] != " " || lines[9] != "RESTOU A SEQUENCIA: " {
			return fmt.Errorf("random-output repeated-value framing mismatch")
		}
		valueText := strings.TrimSuffix(strings.TrimPrefix(lines[6], " "), " ")
		value, err := strconv.ParseInt(valueText, 10, 64)
		if err != nil || value < c.Minimum || value >= c.Bound || lines[6] != fmt.Sprintf(" %d ", value) || lines[10] != fmt.Sprintf(" %d", value) {
			return fmt.Errorf("random-output repeated value mismatch")
		}
		if lines[7] != "O NUMERO DE REPETIÇÕES FOI:  0" {
			return fmt.Errorf("random-output repeat count mismatch")
		}
		return nil
	}
	if c.Kind == "random-int-search-table-lines" {
		for index := range 10 {
			line := lines[index]
			if len(line) != 10 {
				return fmt.Errorf("random-output search row %d width mismatch", index+1)
			}
			valueText := strings.TrimSpace(line[5:])
			value, err := strconv.ParseInt(valueText, 10, 64)
			if err != nil || value < c.Minimum || value >= c.Bound || line != fmt.Sprintf("%5d%5d", index+1, value) {
				return fmt.Errorf("random-output search row %d value mismatch", index+1)
			}
		}
		if lines[10] != "Entre com o valor de busca (ESC termina) :-1" || lines[11] != "Nao achei." {
			return fmt.Errorf("random-output search result mismatch")
		}
		return nil
	}
	if c.Kind == "random-int-search-lines" {
		for index := range 20 {
			line := lines[index]
			value, err := strconv.ParseInt(line, 10, 64)
			if err != nil || line != strconv.FormatInt(value, 10) || value < c.Minimum || value >= c.Bound {
				return fmt.Errorf("random-output search value %d mismatch", index+1)
			}
		}
		if lines[20] != "Valor para busca (ESC ou menor que 0 termina) : -1" {
			return fmt.Errorf("random-output search prompt mismatch")
		}
		return nil
	}
	if c.Kind == "random-record-sort-lines" {
		type record struct {
			code int64
			name string
		}
		input := make([]record, 10)
		for index := range input {
			codePrefix := fmt.Sprintf("Digite o codigo do  %do registro:", index+1)
			namePrefix := fmt.Sprintf("Digite o nome do  %do registro:", index+1)
			codeLine, nameLine := lines[2*index], lines[2*index+1]
			if !strings.HasPrefix(codeLine, codePrefix) || !strings.HasPrefix(nameLine, namePrefix) {
				return fmt.Errorf("random-output record prompt %d mismatch", index+1)
			}
			code, err := strconv.ParseInt(strings.TrimPrefix(codeLine, codePrefix), 10, 64)
			name := strings.TrimPrefix(nameLine, namePrefix)
			if err != nil || strings.TrimPrefix(codeLine, codePrefix) != strconv.FormatInt(code, 10) || code < c.Minimum || code >= c.Bound || len(name) != 5 || strings.Trim(name, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") != "" {
				return fmt.Errorf("random-output record %d input mismatch", index+1)
			}
			input[index] = record{code: code, name: name}
		}
		if lines[20] != "Item - Codigo Nome" || lines[31] != "Item - Codigo Nome" {
			return fmt.Errorf("random-output record headers mismatch")
		}
		canonical := slices.Clone(input)
		slices.SortFunc(canonical, func(left, right record) int {
			if left.code < right.code {
				return -1
			}
			if left.code > right.code {
				return 1
			}
			return strings.Compare(left.name, right.name)
		})
		for offset := range [2]struct{}{} {
			rows := make([]record, 10)
			for index := range rows {
				line := lines[21+offset*11+index]
				prefix := fmt.Sprintf("%4d - ", index+1)
				if !strings.HasPrefix(line, prefix) {
					return fmt.Errorf("random-output record row %d framing mismatch", index+1)
				}
				value := strings.TrimPrefix(line, prefix)
				if len(value) != 12 {
					return fmt.Errorf("random-output record row %d width mismatch", index+1)
				}
				code, err := strconv.ParseInt(strings.TrimSpace(value[:6]), 10, 64)
				name := value[7:]
				got := record{code: code, name: name}
				if err != nil || len(name) != 5 || strings.Trim(name, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") != "" || code < c.Minimum || code >= c.Bound || line != fmt.Sprintf("%s%6d %s", prefix, code, name) {
					return fmt.Errorf("random-output record row %d value mismatch", index+1)
				}
				if index > 0 && ((offset == 0 && rows[index-1].name > got.name) || (offset == 1 && rows[index-1].code > got.code)) {
					return fmt.Errorf("random-output record row %d sort order mismatch", index+1)
				}
				rows[index] = got
			}
			slices.SortFunc(rows, func(left, right record) int {
				if left.code < right.code {
					return -1
				}
				if left.code > right.code {
					return 1
				}
				return strings.Compare(left.name, right.name)
			})
			if !slices.Equal(rows, canonical) {
				return fmt.Errorf("random-output record rows are not an input permutation")
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

func parseFixedDecimal(text string, valueDecimals, textDecimals int) (int64, error) {
	if valueDecimals < 0 || textDecimals < valueDecimals || textDecimals > 10 {
		return 0, fmt.Errorf("invalid decimal precision")
	}
	value, ok := new(big.Rat).SetString(text)
	if !ok {
		return 0, fmt.Errorf("invalid decimal")
	}
	valueScale := int64(math.Pow10(valueDecimals))
	scaled := new(big.Rat).Mul(value, new(big.Rat).SetInt64(valueScale))
	if !scaled.IsInt() || !scaled.Num().IsInt64() {
		return 0, fmt.Errorf("decimal is off grid")
	}
	ticks := scaled.Num().Int64()
	textTicks := ticks
	for precision := valueDecimals; precision < textDecimals; precision++ {
		if textTicks > math.MaxInt64/10 || textTicks < math.MinInt64/10 {
			return 0, fmt.Errorf("decimal is out of range")
		}
		textTicks *= 10
	}
	if formatFixedDecimal(textTicks, textDecimals) != text {
		return 0, fmt.Errorf("non-canonical decimal")
	}
	return ticks, nil
}

func formatFixedDecimal(ticks int64, decimals int) string {
	if decimals == 0 {
		return strconv.FormatInt(ticks, 10)
	}
	scale := int64(math.Pow10(decimals))
	negative := ticks < 0
	if negative {
		if ticks == math.MinInt64 {
			return ""
		}
		ticks = -ticks
	}
	whole, fraction := ticks/scale, ticks%scale
	formatted := fmt.Sprintf("%d.%0*d", whole, decimals, fraction)
	if negative {
		return "-" + formatted
	}
	return formatted
}
