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
		text, err := i.readLine(target.Start())
		if err != nil {
			return err
		}
		v, err := parseInput(text, cell.Type)
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
		if err := checkFormatSize(v, width, decimals); err != nil {
			return failure(arg.Expr.Start(), diag.RStorage, err)
		}
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

func checkFormatSize(v runtime.Value, width, decimals int64) error {
	if width > maxTextBytes {
		return fmt.Errorf("formatted item size limit exceeded")
	}
	n := int64(0)
	switch v.Kind {
	case runtime.StringValue:
		n = int64(len(v.Str))
		if width > 0 {
			n = width
			count := int64(0)
			for _, r := range v.Str {
				if count == width {
					break
				}
				n += int64(utf8.RuneLen(r) - 1)
				count++
			}
		}
	case runtime.IntegerValue:
		n = int64(len(strconv.FormatInt(v.Int, 10)))
		if width <= 0 {
			n++
		} else if decimals > 0 {
			if decimals > maxTextBytes-n-1 {
				return fmt.Errorf("formatted item size limit exceeded")
			}
			n += decimals + 1
		}
	case runtime.RealValue:
		if width <= 0 {
			n = int64(len(runtime.FormatReal(v.Real))) + 1
		} else {
			n = int64(len(strconv.FormatFloat(v.Real, 'f', 0, 64)))
			if decimals > 0 {
				if decimals > maxTextBytes-n-1 {
					return fmt.Errorf("formatted item size limit exceeded")
				}
				n += decimals + 1
			}
		}
	}
	if n > maxTextBytes {
		return fmt.Errorf("formatted item size limit exceeded")
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
		return runtime.Value{Kind: runtime.StringValue, Str: text}, nil
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
		decimals = max(0, decimals)
		// The reference rounds decimal ties away from zero.
		scale := math.Pow10(decimals)
		rounded := v.Real
		if scaled := rounded * scale; !math.IsInf(scale, 0) && !math.IsInf(scaled, 0) {
			rounded = math.Round(scaled) / scale
		}
		s = strconv.FormatFloat(rounded, 'f', decimals, 64)
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
	case runtime.VoidValue:
		s = ""
	default:
		s = "<invalido>"
	}
	if width > 0 && len(s) < width {
		return strings.Repeat(" ", width-len(s)) + s
	}
	return s
}
