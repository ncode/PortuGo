package runtime

import "fmt"

// ValueKind is a concrete runtime value category.
type ValueKind int

const (
	InvalidValue ValueKind = iota
	IntegerValue
	RealValue
	StringValue
	BoolValue
	VectorValue
	VoidValue
	RecordValue
)

// Cell is an assignable storage location.
type Cell struct {
	Type  Type
	Value Value
}

// Value is a Portugol runtime value.
type Value struct {
	Kind ValueKind
	Int  int64
	Real float64
	Str  string
	Bool bool
	// Comparison keeps the logical assignment/condition category of a comparison
	// without replacing an unmatched operand's concrete value.
	Comparison bool
	// RealFallback retains the real expression category of mixed division
	// without replacing the right operand's concrete output value.
	RealFallback bool
	// NumericAbsence retains the domain origin of a VoidValue.
	NumericAbsence bool
	Vec            *Vector
	Rec            *Record
}

// Zero returns the zero value for a type.
func Zero(t Type) Value {
	switch t.Kind {
	case IntegerType:
		return Value{Kind: IntegerValue}
	case RealType:
		return Value{Kind: RealValue}
	case StringType:
		return Value{Kind: StringValue}
	case BoolType:
		return Value{Kind: BoolValue}
	case VectorType:
		return Value{Kind: VectorValue, Vec: NewVector(t)}
	case RecordType:
		return Value{Kind: RecordValue, Rec: NewRecord(t)}
	case VoidType:
		return Value{Kind: VoidValue}
	default:
		return Value{Kind: InvalidValue}
	}
}

// Type returns the runtime type of a value.
func (v Value) Type() Type {
	switch v.Kind {
	case IntegerValue:
		return Type{Kind: IntegerType}
	case RealValue:
		return Type{Kind: RealType}
	case StringValue:
		return Type{Kind: StringType}
	case BoolValue:
		return Type{Kind: BoolType}
	case VectorValue:
		if v.Vec != nil {
			return v.Vec.Type
		}
		return Type{Kind: VectorType}
	case RecordValue:
		if v.Rec != nil {
			return v.Rec.Type
		}
		return Type{Kind: RecordType}
	case VoidValue:
		return Type{Kind: VoidType}
	default:
		return Type{Kind: InvalidType}
	}
}

// ConvertForAssign converts integer to real when assigning to a real cell.
func ConvertForAssign(dst Type, v Value) (Value, error) {
	v.RealFallback = false
	if v.Kind == RecordValue {
		if err := v.Rec.validate(); err != nil {
			return Value{}, err
		}
	}
	if v.Comparison {
		if dst.Kind == BoolType {
			if v.Kind == RecordValue {
				// A logical cell has no declared record layout to retain.
				return Zero(Type{Kind: RecordType}), nil
			}
			v.Comparison = false
			return v, nil
		}
		return Value{}, fmt.Errorf("cannot assign logico to %s", dst)
	}
	if dst.Kind == RealType && v.Kind == IntegerValue {
		return Value{Kind: RealValue, Real: float64(v.Int)}, nil
	}
	if Assignable(dst, v.Type()) {
		return v, nil
	}
	return Value{}, fmt.Errorf("cannot assign %s to %s", v.Type(), dst)
}

// Clone returns a deep copy of v where mutation would otherwise be observable.
func Clone(v Value) Value {
	if v.Kind == RecordValue && v.Rec != nil {
		record := &Record{Type: v.Rec.Type.Clone(), Fields: make([]Cell, len(v.Rec.Fields))}
		for n, cell := range v.Rec.Fields {
			record.Fields[n] = Cell{Type: cell.Type.Clone(), Value: Clone(cell.Value)}
		}
		v.Rec = record
		return v
	}
	if v.Kind != VectorValue || v.Vec == nil {
		return v
	}
	vec := &Vector{
		Type:     v.Vec.Type.Clone(),
		Elements: make([]Cell, len(v.Vec.Elements)),
	}
	for i, cell := range v.Vec.Elements {
		vec.Elements[i] = Cell{Type: cell.Type.Clone(), Value: Clone(cell.Value)}
	}
	v.Vec = vec
	return v
}

// Truth returns a boolean value or an error.
func Truth(v Value) (bool, error) {
	if v.Kind != BoolValue && !v.Comparison {
		return false, fmt.Errorf("expected logico, got %s", v.Type())
	}
	return v.Bool, nil
}
