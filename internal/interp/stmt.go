package interp

import (
	"fmt"

	"github.com/ncode/portugol-go/internal/ast"
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

func (i *Interpreter) execStmt(stmt ast.Stmt) (control, error) {
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
		for _, cc := range s.Cases {
			for _, expr := range cc.Values {
				v, err := i.eval(expr)
				if err != nil {
					return control{}, err
				}
				eq, err := equalValues(x, v)
				if err != nil {
					return control{}, err
				}
				if eq {
					return i.execStmts(cc.Body)
				}
			}
		}
		return i.execStmts(s.Default)
	case *ast.WhileStmt:
		for {
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
			if ctrl.kind != noControl {
				return ctrl, nil
			}
		}
	case *ast.RepeatStmt:
		for {
			ctrl, err := i.execStmts(s.Body)
			if err != nil {
				return control{}, err
			}
			if ctrl.kind == breakControl {
				return control{}, nil
			}
			if ctrl.kind != noControl {
				return ctrl, nil
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
		return control{kind: breakControl}, nil
	case *ast.ReturnStmt:
		v, err := i.eval(s.Value)
		if err != nil {
			return control{}, err
		}
		return control{kind: returnControl, value: v}, nil
	case *ast.ReadStmt:
		return control{}, i.execRead(s)
	case *ast.WriteStmt:
		return control{}, i.execWrite(s)
	default:
		return control{}, fmt.Errorf("unsupported statement %T", stmt)
	}
}

func (i *Interpreter) execFor(s *ast.ForStmt) (control, error) {
	cell, err := lookupCell(i.env, s.Name.Text)
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
		if ctrl.kind != noControl {
			return ctrl, nil
		}
		// VisuAlg caps the exposed exit value at the terminal bound, even
		// for descending loops. Body assignments do not change progression.
		final = min(cur+step, to)
	}
	return control{}, assign(cell, runtime.Value{Kind: runtime.IntegerValue, Int: final})
}
