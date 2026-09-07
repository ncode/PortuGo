package interp

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode"
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
		text, err := i.readToken(target.Start())
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
	}
	return nil
}

func (i *Interpreter) readToken(pos token.Pos) (string, error) {
	var text strings.Builder
	for {
		r, _, err := i.in.ReadRune()
		if err == io.EOF && text.Len() != 0 {
			return text.String(), nil
		}
		if err != nil {
			return "", diag.Diagnostic{Code: diag.RInput, Pos: pos, Message: "cannot read input token", Cause: err}
		}
		if unicode.IsSpace(r) {
			if text.Len() != 0 {
				return text.String(), nil
			}
			continue
		}
		if text.Len() > maxTextBytes-utf8.RuneLen(r) {
			return "", failure(pos, diag.RStorage, fmt.Errorf("input token size limit exceeded"))
		}
		text.WriteRune(r)
	}
}

func (i *Interpreter) execWrite(s *ast.WriteStmt) error {
	for _, arg := range s.Args {
		v, err := i.eval(arg.Expr)
		if err != nil {
			return err
		}
		width := int64(0)
		if arg.Width != nil {
			width, err = i.evalInt(arg.Width)
			if err != nil {
				return err
			}
		}
		decimals := int64(-1)
		if arg.Decimals != nil {
			decimals, err = i.evalInt(arg.Decimals)
			if err != nil {
				return err
			}
		}
		if err := checkFormatSize(v, width, decimals); err != nil {
			return failure(arg.Expr.Start(), diag.RStorage, err)
		}
		if _, err := io.WriteString(i.out, formatValue(v, int(max(0, width)), int(max(-1, decimals)))); err != nil {
			return diag.Diagnostic{Code: diag.RHost, Pos: arg.Expr.Start(), Message: "cannot write output", Cause: err}
		}
	}
	if s.Newline {
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
			n = int64(len(strconv.FormatFloat(v.Real, 'f', -1, 64))) + 1
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
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return runtime.Value{}, fmt.Errorf("invalid inteiro input %q", text)
		}
		return runtime.Value{Kind: runtime.IntegerValue, Int: v}, nil
	case runtime.RealType:
		v, err := strconv.ParseFloat(strings.ReplaceAll(text, ",", "."), 64)
		if err != nil {
			return runtime.Value{}, fmt.Errorf("invalid real input %q", text)
		}
		return runtime.Value{Kind: runtime.RealValue, Real: v}, nil
	case runtime.StringType:
		return runtime.Value{Kind: runtime.StringValue, Str: text}, nil
	case runtime.BoolType:
		switch strings.ToLower(text) {
		case "verdadeiro", "true":
			return runtime.Value{Kind: runtime.BoolValue, Bool: true}, nil
		case "falso", "false":
			return runtime.Value{Kind: runtime.BoolValue, Bool: false}, nil
		default:
			return runtime.Value{}, fmt.Errorf("invalid logico input %q", text)
		}
	default:
		return runtime.Value{}, fmt.Errorf("cannot read %s", typ)
	}
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
			return " " + strconv.FormatFloat(v.Real, 'f', -1, 64)
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
