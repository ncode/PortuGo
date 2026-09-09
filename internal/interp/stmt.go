package interp

import (
	"fmt"
	"math"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/diag"
	"github.com/ncode/portugol-go/internal/runtime"
)

func (i *Interpreter) execStmts(stmts []ast.Stmt) (control, error) {
	for _, stmt := range stmts {
		ctrl, err := i.execStmt(stmt)
		if err != nil || ctrl.kind != noControl {
			return ctrl, err
		}
	}
	return control{}, nil
}

func (i *Interpreter) execStmt(stmt ast.Stmt) (ctrl control, err error) {
	if err := i.charge(stmt.Start()); err != nil {
		return ctrl, err
	}
	defer func() { err = failure(stmt.Start(), diag.RType, err) }()
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		cell, err := i.lvalue(s.Target)
		if err != nil {
			return control{}, err
		}
		v, err := i.eval(s.Value)
		if err != nil {
			return control{}, err
		}
		return control{}, assign(cell, v)
	case *ast.CallStmt:
		return control{}, i.callProcedure(s.Call)
	case *ast.IfStmt:
		cond, err := i.evalBool(s.Cond)
		if err != nil {
			return control{}, err
		}
		if cond {
			return i.execStmts(s.Then)
		}
		return i.execStmts(s.Else)
	case *ast.SwitchStmt:
		x, err := i.eval(s.X)
		if err != nil {
			return control{}, err
		}
		if x.Kind == runtime.RealValue {
			x.Real = math.Trunc(x.Real)
			if x.Real < math.MinInt32 || x.Real > math.MaxInt32 {
				x = runtime.Value{Kind: runtime.VoidValue}
			}
		}
		for _, cc := range s.Cases {
			for _, label := range cc.Labels {
				match, err := i.matchesCase(x, label)
				if err != nil {
					return control{}, err
				}
				if match {
					return i.execStmts(cc.Body)
				}
			}
		}
		return i.execStmts(s.Default)
	case *ast.WhileStmt:
		for {
			if err := i.charge(s.Start()); err != nil {
				return control{}, err
			}
			cond, err := i.evalBool(s.Cond)
			if err != nil {
				return control{}, err
			}
			if !cond {
				return control{}, nil
			}
			ctrl, err := i.execStmts(s.Body)
			if err != nil {
				return control{}, err
			}
			if ctrl.kind == breakControl {
				return control{}, nil
			}
		}
	case *ast.RepeatStmt:
		for {
			if err := i.charge(s.Start()); err != nil {
				return control{}, err
			}
			ctrl, err := i.execStmts(s.Body)
			if err != nil {
				return control{}, err
			}
			if ctrl.kind == breakControl {
				return control{}, nil
			}
			cond, err := i.evalBool(s.Cond)
			if err != nil {
				return control{}, err
			}
			if cond {
				return control{}, nil
			}
		}
	case *ast.ForStmt:
		return i.execFor(s)
	case *ast.BreakStmt:
		return control{kind: breakControl, at: s.Start()}, nil
	case *ast.ReturnStmt:
		if i.result == nil {
			return control{}, failure(s.Start(), diag.RCall, fmt.Errorf("retorne outside function"))
		}
		v, err := i.eval(s.Value)
		if err != nil {
			return control{}, err
		}
		return control{}, failure(s.Start(), diag.RCall, assign(i.result, v))
	case *ast.ReadStmt:
		return control{}, i.execRead(s)
	case *ast.WriteStmt:
		return control{}, i.execWrite(s)
	default:
		return control{}, fmt.Errorf("unsupported statement %T", stmt)
	}
}

func (i *Interpreter) execFor(s *ast.ForStmt) (ctrl control, err error) {
	defer func() { err = failure(s.Start(), diag.RLoop, err) }()
	cell, err := i.lookupCell(s.Name)
	if err != nil {
		return control{}, err
	}
	from, err := i.evalInt(s.From)
	if err != nil {
		return control{}, err
	}
	to, err := i.evalInt(s.To)
	if err != nil {
		return control{}, err
	}
	step := int64(1)
	if s.Step != nil {
		step, err = i.evalInt(s.Step)
		if err != nil {
			return control{}, err
		}
	}
	if step == 0 {
		return control{}, fmt.Errorf("para passo cannot be zero")
	}
	final := from
	for cur := from; (step > 0 && cur <= to) || (step < 0 && cur >= to); cur += step {
		if err := i.charge(s.Start()); err != nil {
			return control{}, err
		}
		if err := assign(cell, runtime.Value{Kind: runtime.IntegerValue, Int: cur}); err != nil {
			return control{}, err
		}
		ctrl, err := i.execStmts(s.Body)
		if err != nil {
			return control{}, err
		}
		if ctrl.kind == breakControl {
			final = min(cell.Value.Int, to)
			break
		}
		// VisuAlg caps the exposed exit value at the terminal bound, even
		// for descending loops. Body assignments do not change progression.
		final = min(cur+step, to)
	}
	return control{}, assign(cell, runtime.Value{Kind: runtime.IntegerValue, Int: final})
}
