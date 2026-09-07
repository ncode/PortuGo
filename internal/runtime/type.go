package runtime

import (
	"fmt"
	"slices"
	"strings"
	"unsafe"

	"github.com/ncode/portugol-go/internal/ast"
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
)

// Range is one vector dimension bound.
type Range struct {
	Low  int64
	High int64
}

// Type describes a Portugol value type.
type Type struct {
	Kind   TypeKind
	Elem   *Type
	Ranges []Range
}

// Clone returns an independent copy of the type and its layout.
func (t Type) Clone() Type {
	t.Ranges = slices.Clone(t.Ranges)
	if t.Elem != nil {
		elem := t.Elem.Clone()
		t.Elem = &elem
	}
	return t
}

// Slots returns the checked number of scalar storage slots in this type.
// Reference-specific storage quotas are separate from representability.
func (t Type) Slots() (int, error) {
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
		width := uint64(r.High) - uint64(r.Low) + 1
		if r.High < r.Low || width == 0 || width > limit/uint64(n) {
			return 0, fmt.Errorf("vector layout exceeds addressable storage")
		}
		n *= int(width)
	}
	return n, nil
}

// TypeFromSpec converts a parsed type into a runtime type.
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
			ranges[i] = Range{Low: r.Low, High: r.High}
		}
		return Type{Kind: VectorType, Elem: &elem, Ranges: ranges}
	default:
		return Type{Kind: InvalidType}
	}
}

// Equal reports whether two types are identical.
func (t Type) Equal(o Type) bool {
	if t.Kind != o.Kind || len(t.Ranges) != len(o.Ranges) {
		return false
	}
	for i := range t.Ranges {
		if t.Ranges[i] != o.Ranges[i] {
			return false
		}
	}
	if t.Elem == nil || o.Elem == nil {
		return t.Elem == nil && o.Elem == nil
	}
	return t.Elem.Equal(*o.Elem)
}

// String returns a source-like type name.
func (t Type) String() string {
	switch t.Kind {
	case IntegerType:
		return "inteiro"
	case RealType:
		return "real"
	case StringType:
		return "caractere"
	case BoolType:
		return "logico"
	case VectorType:
		parts := make([]string, len(t.Ranges))
		for i, r := range t.Ranges {
			parts[i] = fmt.Sprintf("%d..%d", r.Low, r.High)
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

// Assignable reports whether a value of src can be assigned to dst.
func Assignable(dst, src Type) bool {
	if dst.Equal(src) {
		return true
	}
	return dst.Kind == RealType && src.Kind == IntegerType
}
