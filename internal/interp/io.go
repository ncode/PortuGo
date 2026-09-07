package interp

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (i *Interpreter) execRead(s *ast.ReadStmt) error {
	for _, target := range s.Targets {
		cell, err := i.lvalue(target)
		if err != nil {
			return err
		}
		if !i.in.Scan() {
			if err := i.in.Err(); err != nil {
				return err
			}
			return fmt.Errorf("not enough input for leia")
		}
		v, err := parseInput(i.in.Text(), cell.Type)
		if err != nil {
			return err
		}
		if err := assign(cell, v); err != nil {
			return err
		}
	}
	return nil
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
		if _, err := fmt.Fprint(i.out, formatValue(v, int(width), int(decimals))); err != nil {
			return err
		}
	}
	if s.Newline {
		_, err := fmt.Fprintln(i.out)
		return err
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
			return fmt.Sprintf("%-*.*s", width, width, v.Str)
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
		return fmt.Sprintf("%*s", width, s)
	}
	return s
}
