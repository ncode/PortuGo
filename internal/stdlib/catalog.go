package stdlib

import (
	"fmt"
	"math"
	"unicode"

	"github.com/ncode/portugol-go/internal/runtime"
)

// This catalog is constructed once and never mutated. Descriptors contain no
// execution state, and their accessors do not expose the registry's storage.
var builtins = buildCatalog()

func buildCatalog() registry {
	numericParameter := Parameter{Type: runtime.NumericType, Mode: ByValue}
	integerParameter := Parameter{Type: runtime.IntegerType, Mode: ByValue}
	textParameter := Parameter{Type: runtime.StringType, Mode: ByValue}
	numeric := Signature{Form: Parenthesized, Rule: UnaryNumericCall, Arity: 1,
		Parameters: [3]Parameter{numericParameter}, Result: runtime.RealType, Empty: runtime.RealType}
	preserving := numeric
	preserving.Result, preserving.Empty = runtime.NumericType, runtime.VoidType
	preserving.PreserveNumeric, preserving.ArityFirst, preserving.NonNumericAbsent = true, true, true
	root := numeric
	root.ArityFirst, root.NonNumericAbsent = true, true
	integer := numeric
	integer.Result, integer.Empty, integer.NonNumericAbsent = runtime.IntegerType, runtime.IntegerType, true
	power := Signature{Form: Parenthesized, Rule: PowerCall, Arity: 2,
		Parameters: [3]Parameter{numericParameter, numericParameter}, Result: runtime.RealType, Empty: runtime.VoidType}
	constant := Signature{Form: Bare, Rule: ConstantCall, Result: runtime.RealType, Empty: runtime.RealType}
	numberText := numeric
	numberText.Rule, numberText.Result, numberText.Empty = NumberTextCall, runtime.StringType, runtime.StringType
	text := Signature{Form: Parenthesized, Rule: TextCall, Arity: 1,
		Parameters: [3]Parameter{textParameter}, Result: runtime.StringType}
	textInteger := text
	textInteger.Result = runtime.IntegerType
	textNumber := text
	textNumber.Result = runtime.NumericType
	search := textInteger
	search.Arity, search.Parameters[1] = 2, textParameter
	copyText := text
	copyText.Arity = 3
	numericParameter.AbsenceAsZero = true
	copyText.Parameters[1], copyText.Parameters[2] = numericParameter, numericParameter
	random := Signature{Form: Parenthesized, Rule: RandomIntegerCall, Arity: 1,
		Parameters: [3]Parameter{integerParameter}, Result: runtime.IntegerType, Empty: runtime.IntegerType}
	character := random
	character.Rule, character.Result, character.Empty = CharacterCall, runtime.StringType, runtime.StringType
	character.Parameters[0].AbsenceAsZero = true
	entries := []Descriptor{
		builtin("abs", preserving, Domain{"signed 32-bit integer wrapping; real absolute value", NoDomainAbsence}, pure(abs)),
		builtin("arccos", numeric, Domain{"[-1, 1] to radians; outside the domain produces numeric absence", NumericDomainAbsence}, realOperation(math.Acos)),
		builtin("arcsen", numeric, Domain{"[-1, 1] to radians; outside the domain produces numeric absence", NumericDomainAbsence}, realOperation(math.Asin)),
		builtin("arctan", numeric, Domain{"numeric argument to radians", NoDomainAbsence}, realOperation(math.Atan)),
		builtin("asc", textInteger, Domain{"first Windows-1252 code; empty text produces absence; unsupported character receives R007", GenericDomainAbsence}, pure(asc)),
		builtin("carac", character, Domain{"recorded 0..255 character table; out-of-range integer produces absence", GenericDomainAbsence}, pure(carac)),
		builtin("caracpnum", textNumber, Domain{"recorded numeric text conversion and zero-prefix fallback; malformed nonzero suffix receives R007", NoDomainAbsence}, pure(caracpnum)),
		builtin("compr", textInteger, Domain{"decoded character count", NoDomainAbsence}, pure(compr)),
		builtin("copia", copyText, Domain{"1-based substring; signed 32-bit truncated bounds; absent bounds become zero", NoDomainAbsence}, pure(copia)),
		builtin("cos", numeric, Domain{"radians to cosine", NoDomainAbsence}, realOperation(math.Cos)),
		builtin("cotan", numeric, Domain{"radians to reciprocal tangent; nonfinite result produces numeric absence", NumericDomainAbsence}, realOperation(func(x float64) float64 { return 1 / math.Tan(x) })),
		builtin("exp", power, Domain{"finite power; invalid or nonfinite result receives R007; underflow yields zero", NoDomainAbsence}, pure(exp)),
		builtin("grauprad", numeric, Domain{"degrees to radians without premature overflow", NoDomainAbsence}, realOperation(func(x float64) float64 { return x * (math.Pi / 180) })),
		builtin("int", integer, Domain{"truncate through signed 64-bit then narrow to signed 32-bit; unsupported range receives R007", NoDomainAbsence}, pure(intval)),
		builtin("log", numeric, Domain{"positive argument to base-ten logarithm; invalid or nonfinite result receives R007", NoDomainAbsence}, pure(func(args []runtime.Value) (runtime.Value, bool, error) { return checkedNumeric1(args, math.Log10) })),
		builtin("logn", numeric, Domain{"positive argument to natural logarithm; invalid or nonfinite result receives R007", NoDomainAbsence}, pure(func(args []runtime.Value) (runtime.Value, bool, error) { return checkedNumeric1(args, math.Log) })),
		builtin("maiusc", text, Domain{"Windows-1252 uppercase; preserve micro sign and florin", NoDomainAbsence}, pure(func(args []runtime.Value) (runtime.Value, bool, error) {
			return string1(args, func(r rune) rune {
				if r == 'µ' || r == 'ƒ' {
					return r
				}
				return unicode.ToUpper(r)
			})
		})),
		builtin("minusc", text, Domain{"Windows-1252 lowercase", NoDomainAbsence}, pure(func(args []runtime.Value) (runtime.Value, bool, error) { return string1(args, unicode.ToLower) })),
		builtin("numpcarac", numberText, Domain{"finite number to 15 significant digits; nonnumeric input produces absence; nonfinite receives R007", NoDomainAbsence}, pure(numpcarac)),
		builtin("pi", constant, Domain{"real circle constant", NoDomainAbsence}, func(_ *Library, _ []runtime.Value) (runtime.Value, bool, error) {
			return runtime.Value{Kind: runtime.RealValue, Real: math.Pi}, true, nil
		}),
		builtin("pos", search, Domain{"1-based match position; empty or absent match returns zero", NoDomainAbsence}, pure(pos)),
		builtin("quad", preserving, Domain{"signed 32-bit integer wrapping; nonfinite real square receives R007", NoDomainAbsence}, pure(square)),
		builtin("radpgrau", numeric, Domain{"radians to degrees; nonfinite result produces numeric absence", NumericDomainAbsence}, realOperation(func(x float64) float64 { return x * (180 / math.Pi) })),
		builtin("raizq", root, Domain{"nonnegative argument; invalid or nonfinite result receives R007", NoDomainAbsence}, pure(squareRoot)),
		builtin("rand", constant, Domain{"injected fraction in [0, 1); missing or invalid source receives R007", NoDomainAbsence}, func(l *Library, args []runtime.Value) (runtime.Value, bool, error) {
			if len(args) != 0 {
				return runtime.Value{}, true, fmt.Errorf("rand takes no arguments")
			}
			return l.randomFraction()
		}),
		builtin("randi", random, Domain{"unsigned 32-bit bound width and signed result; empty or zero bound yields zero", NoDomainAbsence}, (*Library).randi),
		builtin("sen", numeric, Domain{"radians to sine", NoDomainAbsence}, realOperation(math.Sin)),
		builtin("tan", numeric, Domain{"radians to tangent", NoDomainAbsence}, realOperation(math.Tan)),
	}
	r, err := newRegistry(entries)
	if err != nil {
		panic(fmt.Errorf("internal: invalid builtin catalog: %w", err))
	}
	return r
}

func builtin(name string, signature Signature, domain Domain, evaluate func(*Library, []runtime.Value) (runtime.Value, bool, error)) Descriptor {
	return Descriptor{name: name, signature: signature, domain: domain, evaluate: evaluate}
}

func pure(evaluate func([]runtime.Value) (runtime.Value, bool, error)) func(*Library, []runtime.Value) (runtime.Value, bool, error) {
	return func(_ *Library, args []runtime.Value) (runtime.Value, bool, error) { return evaluate(args) }
}

func realOperation(operation func(float64) float64) func(*Library, []runtime.Value) (runtime.Value, bool, error) {
	return pure(func(args []runtime.Value) (runtime.Value, bool, error) { return numeric1(args, operation) })
}
