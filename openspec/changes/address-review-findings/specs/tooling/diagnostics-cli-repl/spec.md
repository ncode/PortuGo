## MODIFIED Requirements

### Requirement: Formatter contract
`portugol fmt` SHALL use the canonical printer, preserve source behavior and all accepted comments and reference-ignored suffixes, be idempotent, and reject malformed or resource-limited input without overwriting the original. Comment preservation SHALL retain each comment's association with its declaration, statement, or closing delimiter across every supported global declaration order. Check, standard-output, and in-place modes SHALL have documented and testable exit statuses.

#### Scenario: Format a valid file in place
- **WHEN** in-place formatting succeeds
- **THEN** the file contains canonical source that parses equivalently and a second formatting pass makes no change

#### Scenario: Format a malformed file
- **WHEN** in-place formatting encounters syntax diagnostics
- **THEN** it exits with status `1`, reports the diagnostics, and leaves the original bytes unchanged

#### Scenario: Preserve comments with trailing global declarations
- **WHEN** a supported program declares a procedure or function before its single global `var` section and contains leading, body, closing-delimiter, and variable comments
- **THEN** formatting retains those comments with their original constructs instead of moving subprogram comments into the global section or variable comments onto a subprogram terminator
- **AND** the formatted program preserves execution, and repeated formatting is byte-identical

#### Scenario: Preserve formatter modes for either declaration order
- **WHEN** a commented program places global variables before or after its subprograms
- **THEN** standard-output and in-place modes produce the same canonical bytes, and check mode succeeds on those bytes without executing the program

## ADDED Requirements

### Requirement: Live runtime storage safeguard
Each interpreter run SHALL enforce a fixed budget of 2,097,152 live storage units independently of its execution-step budget and existing per-aggregate scalar-slot limit. Each allocated variable, constant, parameter, function-result, record-field, or vector-element storage cell SHALL count as one unit; each allocated vector or record container SHALL count as one additional unit. Nested aggregates and temporary aggregate copies created during execution SHALL include their complete storage recursively. Shared references SHALL not count an existing cell or container a second time.

The implementation SHALL compute required units with checked arithmetic before allocating storage. An allocation that would exceed the remaining budget SHALL return positioned `R003` before that allocation, preserve preceding language output, and unwind normally. Released call-frame storage and completed temporary copies SHALL release their units; starting a new run or REPL submission SHALL start fresh accounting. Retained storage SHALL remain charged for its actual lifetime.

The budget SHALL be documented as a project storage safeguard, not a measured VisuAlg limit or a bound on total process memory. Source buffers, immutable semantic metadata, string-byte allocations, existing operand/output guards, and caller-requested state snapshots retain their separate contracts. The change SHALL add no command-line memory option. A safeguard hit SHALL remain an acceptance failure for a previously verified reference program.

#### Scenario: Reject cumulative declarations before allocation
- **WHEN** individually valid aggregate declarations collectively request more live storage than the run's budget
- **THEN** execution returns `R003` at the first declaration whose allocation would exceed the remaining budget before allocating that variable
- **AND** a small or unlimited execution-step budget does not bypass the storage safeguard

#### Scenario: Count nested aggregate containers
- **WHEN** a nested vector layout has an allowed scalar-slot count but its cells and containers exceed the live-storage budget
- **THEN** execution returns positioned `R003` before expanding the nested allocation
- **AND** accounting does not enumerate or allocate the layout's individual elements to discover the excess

#### Scenario: Share the budget across recursive calls
- **WHEN** a callee's parameters, result, local variables, or execution-time copies would exceed the budget shared with its active callers
- **THEN** the failing allocation reports positioned `R003` and existing frames unwind without a panic or leaked reservation

#### Scenario: Reuse released frame storage
- **WHEN** sequential completed calls cumulatively allocate more units than the limit over time but their simultaneously live storage stays within the budget
- **THEN** those calls continue without a cumulative-allocation failure

#### Scenario: Bound execution-time aggregate copies
- **WHEN** a record assignment or another execution-time aggregate copy would exceed the remaining budget
- **THEN** it reports positioned `R003` before allocating that copy and leaves the assignment destination unchanged
- **AND** a successful copy releases temporary reservations while preserving captured field references and independent stored values

#### Scenario: Reset storage accounting after failure or halt
- **WHEN** a run ends through success, diagnostic, or halt and a new interpreter run or REPL submission starts
- **THEN** the new run receives a fresh budget without retained charges from the preceding program

#### Scenario: Retain earlier builtin arguments during later calls
- **WHEN** a builtin's earlier argument holds an aggregate and a later argument calls a subprogram that replaces its original owning value
- **THEN** the aggregate stays charged until the builtin invocation releases that argument
- **AND** storage-root tracking preserves the language's existing operand and argument-evaluation behavior
