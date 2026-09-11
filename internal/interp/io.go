package interp

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
	"github.com/ncode/portugol-go/internal/token"
)

func (i *Interpreter) execRead(s *ast.ReadStmt) error {
	for _, target := range s.Targets {
		cell, err := i.lvalue(target)
		if err != nil {
			return err
		}
		var text string
		var v runtime.Value
		if i.randomInput.active && cell.Type.Kind != runtime.BoolType {
			state := i.randomInput
			v, err = i.lib.RandomInput(cell.Type.Kind, state.low, state.high, state.decimals)
			text = v.Str
		} else {
			text, err = i.readLine(target.Start())
			if err != nil {
				return err
			}
			v, err = parseInput(text, cell.Type)
		}
		if err != nil {
			return failure(target.Start(), diag.RInput, err)
		}
		if err := assign(cell, v); err != nil {
			return err
		}
		switch v.Kind {
		case runtime.IntegerValue:
			text = strconv.FormatInt(v.Int, 10)
		case runtime.RealValue:
			text = strconv.FormatFloat(v.Real, 'f', 10, 64)
		case runtime.BoolValue:
			text = "Falso"
			if v.Bool {
				text = "Verdadeiro"
			}
		}
		if _, err := fmt.Fprintln(i.out, text); err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: target.Start(), Message: "cannot echo input", Cause: err}
		}
	}
	return nil
}

func (i *Interpreter) readLine(pos token.Pos) (string, error) {
	var text strings.Builder
	for {
		r, _, err := i.in.ReadRune()
		if err == io.EOF && text.Len() != 0 {
			return text.String(), nil
		}
		if err != nil {
			return "", diag.Diagnostic{Code: diag.RInput, Pos: pos, Message: "cannot read input line", Cause: err}
		}
		if r == '\n' {
			return strings.TrimSuffix(text.String(), "\r"), nil
		}
		if text.Len() > maxTextBytes-utf8.RuneLen(r) {
			return "", failure(pos, diag.RStorage, fmt.Errorf("input line size limit exceeded"))
		}
		text.WriteRune(r)
	}
}

func (i *Interpreter) execWrite(s *ast.WriteStmt) error {
	// Nested writes consume the newline requested by an outer escreval.
	i.writeNewline = i.writeNewline || s.Newline
	items := make([]string, len(s.Args))
	buffered := 0
	defer func() { i.writeBytes -= buffered }()
	for index, arg := range s.Args {
		v, err := i.eval(arg.Expr)
		if err != nil {
			return err
		}
		if v.Kind == runtime.VoidValue {
			// Discard this statement without consuming a pending newline.
			return nil
		}
		width := int64(0)
		if arg.Width != nil {
			if v.Kind == runtime.BoolValue {
				return failure(arg.Width.Start(), diag.EParse, fmt.Errorf("cannot format logico with a field width"))
			}
			width, err = i.evalInt(arg.Width, diag.EParse)
			if err != nil {
				return err
			}
		}
		decimals := int64(-1)
		if arg.Decimals != nil {
			decimals, err = i.evalInt(arg.Decimals, diag.EParse)
			if err != nil {
				return err
			}
		}
		width = min(width, runtime.MaxTextChars)
		// The reference caps requested decimal places before expanding the field.
		decimals = min(decimals, 216)
		text := formatValue(v, int(max(0, width)), int(max(-1, decimals)))
		if len(text) > maxTextBytes-i.writeBytes {
			return failure(arg.Expr.Start(), diag.RStorage, fmt.Errorf("pending output size limit exceeded"))
		}
		items[index] = text
		buffered += len(text)
		i.writeBytes += len(text)
	}
	for index, text := range items {
		if _, err := io.WriteString(i.out, text); err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: s.Args[index].Expr.Start(), Message: "cannot write output", Cause: err}
		}
	}
	if i.writeNewline {
		i.writeNewline = false
		_, err := fmt.Fprintln(i.out)
		if err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: s.Start(), Message: "cannot write output", Cause: err}
		}
	}
	return nil
}

func parseInput(text string, typ runtime.Type) (runtime.Value, error) {
	switch typ.Kind {
	case runtime.IntegerType:
		v, err := parseInputNumber(text)
		if err != nil {
			return runtime.Value{}, err
		}
		if v < -0x1p63 || v >= 0x1p63 {
			return runtime.Value{}, fmt.Errorf("inteiro input out of range")
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: int64(int32(int64(v)))}, nil
	case runtime.RealType:
		v, err := parseInputNumber(text)
		if err != nil {
			return runtime.Value{}, err
		}
		return runtime.Value{Kind: runtime.RealValue, Real: v}, nil
	case runtime.StringType:
		return runtime.Value{Kind: runtime.StringValue, Str: runtime.LimitText(text)}, nil
	case runtime.BoolType:
		return runtime.Value{Kind: runtime.BoolValue, Bool: len(text) != 0 && (text[0] == 'v' || text[0] == 'V')}, nil
	default:
		return runtime.Value{}, fmt.Errorf("cannot read %s", typ)
	}
}

func parseInputNumber(text string) (float64, error) {
	text = strings.ReplaceAll(strings.TrimLeft(text, " "), ",", ".")
	start := 0
	if len(text) != 0 && (text[0] == '+' || text[0] == '-') {
		start++
	}
	end := start
	for end < len(text) && text[end] >= '0' && text[end] <= '9' {
		end++
	}
	if end < len(text) && text[end] == '.' {
		end++
		for end < len(text) && text[end] >= '0' && text[end] <= '9' {
			end++
		}
	}
	mantissa := text[start:end]
	valid := mantissa != "" && mantissa != "."
	if end < len(text) && (text[end] == 'e' || text[end] == 'E') {
		end++
		if end < len(text) && (text[end] == '+' || text[end] == '-') {
			end++
		}
		exponent := end
		for end < len(text) && text[end] >= '0' && text[end] <= '9' {
			end++
		}
		valid = valid && end > exponent
	}
	if !valid || end != len(text) {
		// Malformed input retains the unsigned mantissa before decimal/exponent scaling.
		text = strings.ReplaceAll(mantissa, ".", "")
		if text == "" {
			return 0, nil
		}
	}
	v, err := strconv.ParseFloat(text, 64)
	if err != nil || math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, fmt.Errorf("numeric input out of range")
	}
	if v == 0 {
		return 0, nil // Input conversion discards the sign of zero.
	}
	return v, nil
}

func formatValue(v runtime.Value, width, decimals int) string {
	var s string
	switch v.Kind {
	case runtime.IntegerValue:
		s = strconv.FormatInt(v.Int, 10)
		if width <= 0 {
			return " " + s
		}
		if decimals > 0 {
			s += "." + strings.Repeat("0", decimals)
		}
	case runtime.RealValue:
		if width <= 0 {
			return " " + runtime.FormatReal(v.Real)
		}
		if math.Abs(v.Real) >= 0x1p120 {
			width = max(width, 10)
			s = strconv.FormatFloat(v.Real, 'E', min(width-9, 17), 64)
			if mantissa, exponent, ok := strings.Cut(s, "E"); ok {
				s = mantissa + "E" + exponent[:1] + strings.Repeat("0", 5-len(exponent)) + exponent[1:]
			}
		} else {
			s = formatRealFixed(v.Real, max(0, decimals))
		}
	case runtime.StringValue:
		if width > 0 {
			count, end := 0, len(v.Str)
			for offset := range v.Str {
				if count == width {
					end = offset
					break
				}
				count++
			}
			return v.Str[:end] + strings.Repeat(" ", width-count)
		}
		return v.Str
	case runtime.BoolValue:
		if v.Bool {
			return " VERDADEIRO"
		}
		return " FALSO"
	case runtime.VoidValue, runtime.RecordValue:
		s = ""
	default:
		s = "<invalido>"
	}
	if width > 0 && len(s) < width {
		return strings.Repeat(" ", width-len(s)) + s
	}
	return s
}

func formatRealFixed(value float64, decimals int) string {
	_, binaryExponent := math.Frexp(value)
	// Recorded fixed conversion budgets digits from the binary exponent, then
	// pads requested places beyond that budget with zeros.
	calculated := min(decimals, 17, 18-int(math.Ceil(float64(binaryExponent-1)*math.Log10(2))))
	var text string
	if calculated < 0 {
		scientific := strconv.FormatFloat(value, 'e', 18, 64)
		_, exponent, _ := strings.Cut(scientific, "e")
		exp, _ := strconv.Atoi(exponent) // Generated by FormatFloat, not source input.
		scientific = strconv.FormatFloat(value, 'e', exp+calculated, 64)
		mantissa, _, _ := strings.Cut(scientific, "e")
		text = strings.ReplaceAll(mantissa, ".", "") + strings.Repeat("0", -calculated)
		calculated = 0
	} else {
		// Power-of-two scaling detects exact decimal halves without creating
		// the false ties caused by multiplying a binary float by a power of ten.
		if math.Abs(math.Mod(math.Ldexp(value, calculated+1), 2)) == 1 {
			value = math.Nextafter(value, math.Copysign(math.Inf(1), value))
		}
		text = strconv.FormatFloat(value, 'f', calculated, 64)
	}
	if decimals > calculated {
		if calculated == 0 {
			text += "."
		}
		text += strings.Repeat("0", decimals-calculated)
	}
	return text
}
