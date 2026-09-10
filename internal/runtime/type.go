package runtime

import (
	"fmt"
	"slices"
	"strings"
	"unsafe"

	"github.com/ncode/portugol-go/internal/ast"
	"github.com/ncode/portugol-go/internal/token"
)

// TypeKind is a Portugol runtime type category.
type TypeKind int

const (
	InvalidType TypeKind = iota
	IntegerType
	RealType
	StringType
	BoolType
	VectorType
	VoidType
	// NumericType is an analysis-only union; runtime values remain integer or real.
	NumericType
	RecordType
)

// Range is one vector dimension bound. Dynamic endpoints are unresolved
// declaration expressions; concrete allocated vectors never retain these flags.
type Range struct {
	Low         int64
	High        int64
	LowDynamic  bool
	HighDynamic bool
}

// Type describes a Portugol value type.
type Type struct {
	Kind     TypeKind
	Elem     *Type
	Ranges   []Range
	RecordID token.Pos
	Fields   []Field
}

// Field describes one scalar field in declaration order.
type Field struct {
	Name string
	Type Type
}

// Clone returns an independent copy of the type and its layout.
func (t Type) Clone() Type {
	t.Ranges = slices.Clone(t.Ranges)
	t.Fields = slices.Clone(t.Fields)
	for n := range t.Fields {
		t.Fields[n].Type = t.Fields[n].Type.Clone()
	}
	if t.Elem != nil {
		elem := t.Elem.Clone()
		t.Elem = &elem
	}
	return t
}

// Slots returns the checked number of scalar storage slots in this type.
// Reference-specific storage quotas are separate from representability.
func (t Type) Slots() (int, error) {
	if t.Kind == RecordType {
		for _, field := range t.Fields {
			if field.Type.Kind < IntegerType || field.Type.Kind > BoolType {
				return 0, fmt.Errorf("invalid record field layout")
			}
		}
		if len(t.Fields) > MaxVectorSlots {
			return 0, ErrVectorSize
		}
		// Empty records still occupy one cell when used as vector elements.
		return max(1, len(t.Fields)), nil
	}
	if t.Kind != VectorType {
		return 1, nil
	}
	if t.Elem == nil || len(t.Ranges) == 0 {
		return 0, fmt.Errorf("invalid vector layout")
	}
	n, err := t.Elem.Slots()
	if err != nil {
		return 0, err
	}
	limit := uint64(^uint(0)>>1) / uint64(unsafe.Sizeof(Cell{}))
	for _, r := range t.Ranges {
		if r.LowDynamic || r.HighDynamic {
			return 0, fmt.Errorf("vector bounds require declaration initialization")
		}
		width := uint64(r.High) - uint64(r.Low) + 1
		if r.High < r.Low || width == 0 || width > limit/uint64(n) {
			return 0, fmt.Errorf("vector layout exceeds addressable storage")
		}
		n *= int(width)
		if n > MaxVectorSlots {
			return 0, ErrVectorSize
		}
	}
	return n, nil
}

// TypeFromSpec converts a parsed type, marking named bounds for initialization.
func TypeFromSpec(spec ast.TypeSpec) Type {
	switch strings.ToLower(spec.Name) {
	case "inteiro":
		return Type{Kind: IntegerType}
	case "real":
		return Type{Kind: RealType}
	case "caractere":
		return Type{Kind: StringType}
	case "logico":
		return Type{Kind: BoolType}
	case "vetor":
		elem := Type{Kind: InvalidType}
		if spec.Elem != nil {
			elem = TypeFromSpec(*spec.Elem)
		}
		ranges := make([]Range, len(spec.Ranges))
		for i, r := range spec.Ranges {
			low, lowDynamic := boundFromSpec(r.Low)
			high, highDynamic := boundFromSpec(r.High)
			ranges[i] = Range{Low: low, High: high, LowDynamic: lowDynamic, HighDynamic: highDynamic}
		}
		return Type{Kind: VectorType, Elem: &elem, Ranges: ranges}
	default:
		return Type{Kind: InvalidType}
	}
}

func boundFromSpec(expr ast.Expr) (int64, bool) {
	if lit, ok := expr.(*ast.LiteralExpr); ok && lit.Kind == ast.IntLiteral {
		return lit.Int, false
	}
	return 0, true
}

// DynamicBounds reports whether any dimension needs declaration-time values.
func (t Type) DynamicBounds() bool {
	for _, r := range t.Ranges {
		if r.LowDynamic || r.HighDynamic {
			return true
		}
	}
	return t.Elem != nil && t.Elem.DynamicBounds()
}

// Equal reports whether two types are identical.
func (t Type) Equal(o Type) bool {
	return t.equal(o, false)
}

func (t Type) equal(o Type, allowDynamic bool) bool {
	if t.Kind != o.Kind || t.RecordID != o.RecordID || len(t.Fields) != len(o.Fields) || len(t.Ranges) != len(o.Ranges) {
		return false
	}
	for n, field := range t.Fields {
		if field.Name != o.Fields[n].Name || !field.Type.equal(o.Fields[n].Type, allowDynamic) {
			return false
		}
	}
	for i := range t.Ranges {
		a, b := t.Ranges[i], o.Ranges[i]
		if allowDynamic {
			if !a.LowDynamic && !b.LowDynamic && a.Low != b.Low || !a.HighDynamic && !b.HighDynamic && a.High != b.High {
				return false
			}
		} else if a != b {
			return false
		}
	}
	if t.Elem == nil || o.Elem == nil {
		return t.Elem == nil && o.Elem == nil
	}
	return t.Elem.equal(*o.Elem, allowDynamic)
}

// String returns a source-like type name.
func (t Type) String() string {
	switch t.Kind {
	case IntegerType:
		return "inteiro"
	case RealType:
		return "real"
	case NumericType:
		return "inteiro ou real"
	case StringType:
		return "caractere"
	case BoolType:
		return "logico"
	case RecordType:
		return "registro"
	case VectorType:
		parts := make([]string, len(t.Ranges))
		for i, r := range t.Ranges {
			low, high := fmt.Sprint(r.Low), fmt.Sprint(r.High)
			if r.LowDynamic {
				low = "?"
			}
			if r.HighDynamic {
				high = "?"
			}
			parts[i] = low + ".." + high
		}
		elem := "invalido"
		if t.Elem != nil {
			elem = t.Elem.String()
		}
		return "vetor[" + strings.Join(parts, ", ") + "] de " + elem
	case VoidType:
		return "vazio"
	default:
		return "invalido"
	}
}

// Assignable reports whether src satisfies dst, deferring dynamic endpoints.
// Runtime assignment uses concrete layouts, so it checks every endpoint.
func Assignable(dst, src Type) bool {
	if src.Kind == NumericType && (dst.Kind == IntegerType || dst.Kind == RealType) {
		return true
	}
	if dst.equal(src, true) {
		return true
	}
	return dst.Kind == RealType && src.Kind == IntegerType
}
