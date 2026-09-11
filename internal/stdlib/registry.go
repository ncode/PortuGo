package stdlib

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ncode/portugol-go/internal/runtime"
)

// CallForm identifies the syntax accepted by a built-in.
type CallForm uint8

const (
	// Parenthesized requires an argument list, which may be empty.
	Parenthesized CallForm = iota + 1
	// Bare rejects parentheses.
	Bare
)

// CallRule identifies an argument-validation and evaluation-order rule.
type CallRule uint8

const (
	// ConstantCall evaluates a bare zero-argument value.
	ConstantCall CallRule = iota + 1
	// UnaryNumericCall validates one optional numeric argument.
	UnaryNumericCall
	// PowerCall handles conditional arity and left-to-right absence propagation.
	PowerCall
	// NumberTextCall converts an optional number to text.
	NumberTextCall
	// TextCall requires every declared argument, beginning with text arguments.
	TextCall
	// CharacterCall accepts an optional integer code, treating absence as zero.
	CharacterCall
	// RandomIntegerCall accepts an optional integer bound, rejecting absence.
	RandomIntegerCall
)

// ParameterMode records how arguments are passed.
type ParameterMode uint8

// ByValue is the only parameter mode in the recorded built-in inventory.
const ByValue ParameterMode = 1

// Parameter describes one argument's type, mode, and absent-value conversion.
type Parameter struct {
	Type          runtime.TypeKind
	Mode          ParameterMode
	AbsenceAsZero bool
}

// Signature describes a call without owning mutable slices or runtime state.
// Arity is the ordinary argument count; Rule and Empty describe optional and
// absent arguments, including power's conditional evaluation of a third value.
// Empty is InvalidType for TextCall, whose arguments are all required.
type Signature struct {
	Form             CallForm
	Rule             CallRule
	Arity            int
	Parameters       [3]Parameter
	Result           runtime.TypeKind
	Empty            runtime.TypeKind
	PreserveNumeric  bool
	ArityFirst       bool
	NonNumericAbsent bool
}

// ClearsNumericAbsence reports whether a call discards an absence's origin.
func (s Signature) ClearsNumericAbsence() bool {
	return s.Rule == PowerCall || s.Rule == NumberTextCall
}

// AbsenceRule describes an operation's possible domain-dependent absence.
type AbsenceRule uint8

const (
	// NoDomainAbsence adds no possible absence beyond the arguments' own values.
	NoDomainAbsence AbsenceRule = iota
	// GenericDomainAbsence may produce an untyped absent result.
	GenericDomainAbsence
	// NumericDomainAbsence may produce absence that retains its numeric origin.
	NumericDomainAbsence
)

// Domain documents the operation's domain and possible absent results.
type Domain struct {
	Description string
	Absence     AbsenceRule
}

// Descriptor is immutable outside stdlib. Metadata accessors return copies;
// evaluators receive the calling library, never shared random state.
type Descriptor struct {
	name      string
	aliases   []string
	signature Signature
	domain    Domain
	evaluate  func(*Library, []runtime.Value) (runtime.Value, bool, error)
}

// Name returns the canonical lowercase name.
func (d Descriptor) Name() string { return d.name }

// Names returns independent storage containing the canonical name and aliases.
func (d Descriptor) Names() []string { return append([]string{d.name}, d.aliases...) }

// Signature returns the call's signature and evaluation rule.
func (d Descriptor) Signature() Signature { return d.signature }

// Domain returns the operation's documented domain.
func (d Descriptor) Domain() Domain { return d.domain }

type registry struct {
	entries []Descriptor
	byName  map[string]Descriptor
}

func newRegistry(entries []Descriptor) (registry, error) {
	r := registry{entries: slices.Clone(entries), byName: make(map[string]Descriptor)}
	for index, d := range r.entries {
		if err := d.validate(); err != nil {
			return registry{}, fmt.Errorf("builtin %q: %w", d.name, err)
		}
		d.aliases = slices.Clone(d.aliases)
		for _, name := range d.Names() {
			if !canonicalName(name) {
				return registry{}, fmt.Errorf("invalid builtin name %q", name)
			}
			if _, exists := r.byName[name]; exists {
				return registry{}, fmt.Errorf("duplicate builtin name %q", name)
			}
			r.byName[name] = d
		}
		r.entries[index] = d
	}
	return r, nil
}

func canonicalName(name string) bool {
	if name == "" {
		return false
	}
	for index, c := range name {
		if (c < 'a' || c > 'z') && c != '_' && (index == 0 || c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func (d Descriptor) validate() error {
	s := d.signature
	if s.Form != Parenthesized && s.Form != Bare || s.Rule < ConstantCall || s.Rule > RandomIntegerCall {
		return fmt.Errorf("missing call form or rule")
	}
	if s.Arity < 0 || s.Arity > len(s.Parameters) {
		return fmt.Errorf("unsupported arity")
	}
	for index, parameter := range s.Parameters {
		if index >= s.Arity {
			if parameter != (Parameter{}) {
				return fmt.Errorf("parameter exceeds arity")
			}
			continue
		}
		if parameter.Mode != ByValue || parameter.Type != runtime.NumericType && parameter.Type != runtime.IntegerType && parameter.Type != runtime.StringType {
			return fmt.Errorf("missing or unsupported parameter type or mode")
		}
		if parameter.AbsenceAsZero && parameter.Type == runtime.StringType {
			return fmt.Errorf("text parameter cannot convert absence to zero")
		}
	}
	if s.Result != runtime.RealType && s.Result != runtime.IntegerType && s.Result != runtime.StringType && s.Result != runtime.NumericType {
		return fmt.Errorf("missing or unsupported result rule")
	}
	if s.Rule == TextCall && s.Empty != runtime.InvalidType || s.Rule != TextCall && s.Empty != runtime.RealType && s.Empty != runtime.IntegerType && s.Empty != runtime.StringType && s.Empty != runtime.VoidType {
		return fmt.Errorf("missing or unsupported empty-call result")
	}
	if s.PreserveNumeric && (s.Rule != UnaryNumericCall || s.Result != runtime.NumericType) {
		return fmt.Errorf("inconsistent argument-preserving result")
	}
	if s.Rule != UnaryNumericCall && (s.PreserveNumeric || s.ArityFirst || s.NonNumericAbsent) {
		return fmt.Errorf("numeric argument policy on a nonnumeric call")
	}
	if s.Form == Bare != (s.Rule == ConstantCall) {
		return fmt.Errorf("call form contradicts evaluation rule")
	}
	valid := false
	switch s.Rule {
	case ConstantCall:
		valid = s.Arity == 0 && s.Result == runtime.RealType && s.Empty == s.Result
	case UnaryNumericCall:
		valid = s.Arity == 1 && s.Parameters[0].Type == runtime.NumericType &&
			(s.Result == runtime.RealType || s.Result == runtime.IntegerType || s.Result == runtime.NumericType && s.PreserveNumeric)
		if s.PreserveNumeric {
			valid = valid && s.Empty == runtime.VoidType
		} else {
			valid = valid && s.Empty == s.Result
		}
	case PowerCall:
		valid = s.Arity == 2 && s.Parameters[0].Type == runtime.NumericType && s.Parameters[1].Type == runtime.NumericType && s.Result == runtime.RealType && s.Empty == runtime.VoidType
	case NumberTextCall:
		valid = s.Arity == 1 && s.Parameters[0].Type == runtime.NumericType && s.Result == runtime.StringType && s.Empty == s.Result
	case CharacterCall:
		valid = s.Arity == 1 && s.Parameters[0].Type == runtime.IntegerType && s.Result == runtime.StringType && s.Empty == s.Result
	case RandomIntegerCall:
		valid = s.Arity == 1 && s.Parameters[0].Type == runtime.IntegerType && s.Result == runtime.IntegerType && s.Empty == s.Result
	case TextCall:
		valid = s.Arity > 0 && s.Parameters[0].Type == runtime.StringType
		numeric := false
		for _, parameter := range s.Parameters[:s.Arity] {
			numeric = numeric || parameter.Type == runtime.NumericType
			if parameter.Type != runtime.StringType && parameter.Type != runtime.NumericType || numeric && parameter.Type == runtime.StringType {
				valid = false
			}
		}
	}
	for _, parameter := range s.Parameters[:s.Arity] {
		wantZero := s.Rule == CharacterCall || s.Rule == TextCall && parameter.Type == runtime.NumericType
		if parameter.AbsenceAsZero != wantZero {
			valid = false
		}
	}
	if !valid {
		return fmt.Errorf("signature contradicts evaluation rule")
	}
	if strings.TrimSpace(d.domain.Description) == "" || d.domain.Absence > NumericDomainAbsence {
		return fmt.Errorf("missing or unsupported domain")
	}
	if d.evaluate == nil {
		return fmt.Errorf("missing evaluator")
	}
	return nil
}

// Catalog returns independent descriptor storage in declaration order.
func Catalog() []Descriptor { return slices.Clone(builtins.entries) }

// Lookup resolves a canonical lowercase name or alias.
func Lookup(name string) (Descriptor, bool) {
	d, ok := builtins.byName[name]
	return d, ok
}
