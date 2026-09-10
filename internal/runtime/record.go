package runtime

import (
	"fmt"
	"strings"
)

// Record stores scalar fields in declaration order.
type Record struct {
	Type   Type
	Fields []Cell
}

// NewRecord creates zero-valued fields after the caller validates typ.Slots.
func NewRecord(typ Type) *Record {
	r := &Record{Type: typ, Fields: make([]Cell, len(typ.Fields))}
	for n, field := range typ.Fields {
		r.Fields[n] = Cell{Type: field.Type, Value: Zero(field.Type)}
	}
	return r
}

// AssignFrom copies values while preserving captured field locations.
func (r *Record) AssignFrom(source *Record) error {
	if err := r.validate(); err != nil {
		return err
	}
	if err := source.validate(); err != nil {
		return err
	}
	if !r.Type.Equal(source.Type) {
		return fmt.Errorf("incompatible record layouts")
	}
	value := Clone(Value{Kind: RecordValue, Rec: source})
	copy(r.Fields, value.Rec.Fields)
	return nil
}

// Cell returns the first matching field after checking its backing layout.
func (r *Record) Cell(name string) (*Cell, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}
	name = strings.ToLower(name)
	for n, field := range r.Type.Fields {
		if field.Name == name {
			return &r.Fields[n], nil
		}
	}
	return nil, fmt.Errorf("unknown record field %q", name)
}

func (r *Record) validate() error {
	if r == nil || r.Type.Kind != RecordType || len(r.Fields) != len(r.Type.Fields) {
		return fmt.Errorf("invalid record layout")
	}
	_, err := r.Type.Slots()
	return err
}
