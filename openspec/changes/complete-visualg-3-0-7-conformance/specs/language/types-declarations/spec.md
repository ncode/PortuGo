## Purpose

Defines VisuAlg 3.0.7 declarations, types, storage accounting, designators, initialization, assignment, copying, and reference behavior.

## ADDED Requirements

### Requirement: Declaration feature scope
Acceptance requirements and scenarios for constants, named types, aliases, records, fields, and their use inside aggregates SHALL apply only to forms accepted by recorded VisuAlg 3.0.7 evidence. Rejected forms SHALL produce positioned diagnostics and SHALL NOT acquire positive support through dependent vector, assignment, or subprogram requirements. Before dependent implementation, positive requirements and scenarios disproved by the reference SHALL be replaced by explicit rejection requirements with traceable dispositions.

#### Scenario: Reference rejects record declarations
- **WHEN** a program uses a record declaration form rejected by the recorded reference
- **THEN** the implementation diagnoses that form and does not need a record layout, field-access runtime, or successful record example to satisfy conformance

#### Scenario: Reference accepts a disputed declaration form
- **WHEN** reduced probes confirm a candidate declaration form and its use contexts
- **THEN** only those confirmed forms and contexts become positive declaration and dependent aggregate requirements

### Requirement: Ordered declaration sections
The language SHALL support the complete oracle-confirmed declaration section vocabulary for constants, named types, variables, records, and subprogram-local declarations. It SHALL enforce the reference ordering, repetition, visibility, and duplicate-name rules while preserving declarations in source order.

#### Scenario: Analyze a valid declaration sequence
- **WHEN** a program uses multiple declaration section kinds in an order accepted by VisuAlg 3.0.7
- **THEN** every declaration is bound in the reference-defined scope and order without diagnostics

#### Scenario: Reject an invalid declaration sequence
- **WHEN** a declaration section is repeated or placed where the reference rejects it
- **THEN** semantic or syntax analysis reports the positioned traced diagnostic without losing later independent declarations

### Requirement: Constants and declaration expressions
If constant declarations are accepted by the reference, they SHALL be immutable, case-insensitive named values evaluated according to the recorded declaration-time expression rules. Every context that the reference permits to use a constant expression, including vector bounds, case labels, and other declarations, SHALL resolve it before execution and SHALL diagnose cycles, non-constant dependencies, overflow, and invalid types. Rejection of constant declarations SHALL NOT prohibit literal bounds or other expression forms independently accepted by the reference.

#### Scenario: Use a constant in a vector bound
- **WHEN** an integer constant is referenced by a vector bound accepted by the oracle
- **THEN** the declared bound is resolved before runtime and is reflected in indexing and storage accounting

#### Scenario: Detect a cyclic constant
- **WHEN** constant declarations depend on each other cyclically
- **THEN** analysis returns a positioned diagnostic for the cycle and does not recurse indefinitely or panic

### Requirement: Named types and aliases
The language SHALL support oracle-confirmed named-type and alias declarations, resolve alias chains, preserve name identity wherever the reference distinguishes it, and detect unknown or cyclic type definitions. Assignment, parameter matching, and comparison SHALL use the resulting reference-compatible identity and compatibility rules.

#### Scenario: Resolve an alias chain
- **WHEN** a variable is declared through multiple valid type aliases
- **THEN** it receives the final resolved representation while retaining every identity distinction observable in assignment and call compatibility

#### Scenario: Detect a cyclic type alias
- **WHEN** named types form a cycle with no concrete base type accepted by the reference
- **THEN** analysis emits a positioned diagnostic and terminates safely

### Requirement: Record types and fields
The language SHALL support oracle-confirmed record declarations, nested records, record-typed variables, field selection, and fields whose types include named types and vectors. Field names SHALL follow the reference's case and duplicate rules, and field designators SHALL be valid wherever the selected value is valid.

#### Scenario: Assign a nested record field
- **WHEN** a program selects an existing field through a valid chain and assigns a compatible value
- **THEN** only the selected storage location changes

#### Scenario: Select an unknown field
- **WHEN** a designator names a field absent from the resolved record layout
- **THEN** analysis reports a diagnostic at that field selection and execution is not attempted

### Requirement: Vector types and declared bounds
The language SHALL support the recorded one- and two-dimensional vector declarations and reject a third dimension. Literal bounds SHALL be unsigned integers in nondecreasing order, including zero and positive non-one lower bounds. Signed, fractional, parenthesized, and arithmetic literal-bound forms rejected by the reference SHALL receive positioned syntax diagnostics. Independently accepted named-constant bounds remain required. Each dimension SHALL retain its declared lower and upper bounds, and indexing SHALL use those bounds rather than Go slice indices.

#### Scenario: Index a multidimensional vector
- **WHEN** a vector has multiple resolved dimensions and every supplied index is within its declared dimension
- **THEN** the selected element corresponds to the reference row and dimension ordering

#### Scenario: Omit a second index
- **WHEN** a two-dimensional vector is accessed with one index
- **THEN** the second index is its declared lower bound, independent of prior explicit index values

#### Scenario: Reject a third dimension
- **WHEN** a vector declaration specifies three dimensions
- **THEN** parsing returns a positioned syntax diagnostic and execution is not attempted

#### Scenario: Reject an invalid dimension
- **WHEN** a vector bound is unresolved, has an invalid type or order, or overflows during size calculation
- **THEN** a positioned syntax or semantic diagnostic is returned at the corresponding analysis stage and no backing storage is allocated

### Requirement: Reference-confirmed storage limits
The implementation SHALL reproduce storage restrictions established by VisuAlg 3.0.7 recordings, including their accounting unit, declaration scope, record treatment, and vector multiplication. It SHALL NOT reject a program solely for exceeding the previously assumed 500-slot limit: recorded vectors with 500, 501, 5000, and 5001 elements are accepted. These observations do not establish the upper limit or accounting rules for every declaration context. Slot totals SHALL be computed with overflow-safe arithmetic before allocation; independent project resource guards SHALL be identified as such.

#### Scenario: Preserve recorded accepted sizes
- **WHEN** a vector uses one of the recorded accepted sizes of 500, 501, 5000, or 5001 elements
- **THEN** analysis and allocation succeed and the recorded final element can be assigned and read

#### Scenario: Reject unsafe or confirmed excessive storage
- **WHEN** a dimension product overflows, exceeds a documented project resource guard, or violates a separately recorded reference restriction
- **THEN** the implementation emits a positioned storage diagnostic before allocation and identifies the applicable restriction

### Requirement: Zero initialization
Every variable, record field, and vector element SHALL begin with the oracle-confirmed zero value of its resolved type. Reading a valid, unassigned storage location SHALL return that value rather than an uninitialized-memory error.

#### Scenario: Read an unassigned aggregate
- **WHEN** a program reads scalar fields and vector elements before any assignment
- **THEN** it observes the reference zero values for their resolved types

### Requirement: Assignment compatibility and coercion
Assignment SHALL require the exact reference-compatible source and destination types and SHALL apply only oracle-confirmed implicit coercions. Incompatible scalar, aggregate, named-type, field, and vector assignments SHALL be rejected before execution when statically knowable, and the assigned expression SHALL be evaluated once.

#### Scenario: Apply an accepted numeric coercion
- **WHEN** a value is assigned across a numeric type boundary accepted by the reference
- **THEN** the destination receives the same converted value and type as the oracle

#### Scenario: Reject an incompatible aggregate assignment
- **WHEN** a record or vector value is assigned to a destination that is not reference-compatible
- **THEN** analysis emits a positioned type diagnostic and no partial assignment occurs

### Requirement: Aggregate copy and reference semantics
Whole-vector assignment, including self-assignment and assignment to a scalar, SHALL be rejected as recorded. Element assignment remains supported. Record assignment and aggregate parameters SHALL be supported only in independently accepted forms, with copying, conversion, visibility, and copy-back matching their recordings. Nested aggregates SHALL follow the same confirmed rules without accidental sharing or copying; scalar parameter behavior SHALL NOT establish unrecorded aggregate alias behavior.

#### Scenario: Reject whole-vector assignment
- **WHEN** an assignment uses a whole vector as its destination or assigns a whole vector to a scalar
- **THEN** analysis returns a positioned type diagnostic and no partial assignment occurs

#### Scenario: Mutate a value copy
- **WHEN** a subprogram mutates a field or element received through a value parameter
- **THEN** the caller observes or does not observe that mutation exactly as established by the aggregate-copy oracle probes

#### Scenario: Mutate through an alias
- **WHEN** a compatible aggregate designator is passed through a `var` parameter and mutated
- **THEN** the caller's selected storage reflects the mutation

### Requirement: Defensive storage access
Every record and vector access SHALL validate the resolved layout, dimensionality, and bounds before reading or writing storage. Invalid or corrupted access state SHALL produce a positioned diagnostic instead of a Go panic, integer wraparound, or partial mutation.

#### Scenario: Index outside a declared bound
- **WHEN** any evaluated vector index lies below its lower bound or above its upper bound
- **THEN** execution returns the stable positioned storage diagnostic and does not access backing storage
