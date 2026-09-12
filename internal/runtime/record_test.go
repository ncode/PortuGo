package runtime

import "testing"

func TestRecordCellRejectsInvalidStorage(t *testing.T) {
	typ := Type{Kind: RecordType, RecordID: 1, Fields: []Field{{Name: "x", Type: Type{Kind: IntegerType}}}}
	for _, tt := range []struct {
		name   string
		record *Record
		field  string
	}{
		{"nil", nil, "x"},
		{"wrong kind", &Record{Type: Type{Kind: IntegerType}}, "x"},
		{"short backing", &Record{Type: typ}, "x"},
		{"long backing", &Record{Type: typ, Fields: make([]Cell, 2)}, "x"},
		{"non-scalar field", &Record{Type: Type{Kind: RecordType, Fields: []Field{{Name: "x", Type: Type{Kind: RecordType}}}}, Fields: make([]Cell, 1)}, "x"},
		{"missing field", NewRecord(typ), "missing"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if cell, err := tt.record.Cell(tt.field); err == nil || cell != nil {
				t.Fatalf("Cell(%q) = %v, %v; want nil and an error", tt.field, cell, err)
			}
			if tt.name != "missing field" {
				value, err := ConvertForAssign(typ, Value{Kind: RecordValue, Rec: tt.record})
				if err == nil || value.Kind != InvalidValue {
					t.Fatalf("invalid record copied: %v, %v", value, err)
				}
			}
		})
	}
}

func TestRecordCopyAndIdentity(t *testing.T) {
	typ := Type{Kind: RecordType, RecordID: 1, Fields: []Field{
		{Name: "x", Type: Type{Kind: IntegerType}},
		{Name: "x", Type: Type{Kind: RealType}},
	}}
	r := NewRecord(typ)
	cell, err := r.Cell("X")
	if err != nil || cell != &r.Fields[0] || cell.Value.Kind != IntegerValue {
		t.Fatalf("first field = %v, %v", cell, err)
	}
	cell.Value.Int = 7
	copy := Clone(Value{Kind: RecordValue, Rec: r})
	cell.Value.Int = 8
	copy.Rec.Type.Fields[0].Name = "changed"
	if copy.Rec.Fields[0].Value.Int != 7 || r.Type.Fields[0].Name != "x" {
		t.Fatal("record copy shares values or type metadata")
	}
	other := typ.Clone()
	other.RecordID = 2
	if Assignable(typ, other) {
		t.Fatal("distinct record declarations share assignment identity")
	}
}

func TestRecordVectorSlotLimit(t *testing.T) {
	for _, fields := range [][]Field{nil, {{Name: "x", Type: Type{Kind: IntegerType}}, {Name: "y", Type: Type{Kind: RealType}}}} {
		elem := Type{Kind: RecordType, RecordID: 1, Fields: fields}
		slots := max(1, len(fields))
		for _, extra := range []int{0, 1} {
			typ := Type{Kind: VectorType, Elem: &elem, Ranges: []Range{{Low: 1, High: int64(MaxVectorSlots/slots + extra)}}}
			n, err := typ.Slots()
			if (err == nil) != (extra == 0) || err == nil && n != MaxVectorSlots {
				t.Fatalf("record vector slots = %d, %v; extra=%d", n, err, extra)
			}
		}
	}
}

func TestRecordAssignmentPreservesLocations(t *testing.T) {
	typ := Type{Kind: RecordType, RecordID: 1, Fields: []Field{{Name: "x", Type: Type{Kind: IntegerType}}}}
	destination, source := NewRecord(typ), NewRecord(typ)
	field := &destination.Fields[0]
	source.Fields[0].Value.Int = 7
	if err := destination.AssignFrom(source); err != nil {
		t.Fatal(err)
	}
	source.Fields[0].Value.Int = 8
	if field != &destination.Fields[0] || field.Value.Int != 7 {
		t.Fatal("assignment replaced a field location or shared its source value")
	}
	if err := destination.AssignFrom(destination); err != nil || field != &destination.Fields[0] || field.Value.Int != 7 {
		t.Fatalf("self-assignment replaced a captured location: %v", err)
	}
	other := typ.Clone()
	other.RecordID = 2
	for _, invalid := range []*Record{nil, {Type: typ}, {Type: typ, Fields: make([]Cell, 2)}, NewRecord(other)} {
		if err := destination.AssignFrom(invalid); err == nil || field.Value.Int != 7 {
			t.Fatalf("invalid assignment changed destination: %v", err)
		}
	}
	invalid := &Record{Type: typ, Fields: make([]Cell, 2)}
	invalid.Fields[0].Value.Int = 9
	if err := invalid.AssignFrom(source); err == nil || invalid.Fields[0].Value.Int != 9 {
		t.Fatalf("invalid destination was partly changed: %v", err)
	}
}
