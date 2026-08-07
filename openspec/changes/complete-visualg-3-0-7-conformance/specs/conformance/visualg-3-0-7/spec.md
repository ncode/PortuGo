## Purpose

Defines the evidence, authority, traceability, and release gates used to claim observable compatibility with the official VisuAlg 3.0.7 release.

## ADDED Requirements

### Requirement: Reference release authority
The project SHALL treat observable behavior of the official VisuAlg 3.0.7 release as authoritative for language and environment compatibility. When current code, `AGENTS.md`, `docs/language.md`, another written source, or an assumption conflicts with a reproducible reference observation, the recorded reference observation SHALL take precedence and the conflicting project material SHALL be corrected in the same implementation group.

#### Scenario: Current behavior conflicts with the reference
- **WHEN** a reduced program produces one deterministic result in the current implementation and a different result in the recorded VisuAlg 3.0.7 reference run
- **THEN** the conformance expectation uses the reference result and traces the required code, test, and documentation corrections

### Requirement: Reference provenance manifest
The conformance corpus SHALL include a machine-validated manifest that identifies the reference product and version, the acquisition source and date, cryptographic hashes of the uncommitted installer or archive and executable used for recording, the Windows environment, recorder version, normalization version, and evidence files associated with each probe. The repository MUST NOT contain the VisuAlg executable, installer, or distribution archive.

#### Scenario: Validate reference provenance
- **WHEN** conformance validation reads the committed manifest
- **THEN** every probe references complete provenance and evidence metadata while prohibited reference binaries and archives are absent from the repository

### Requirement: Interactive oracle recording
Behavior that is disputed, undocumented, marked `[VERIFICAR]`, dependent on an environment command, or otherwise not established by committed evidence SHALL be recorded on an interactive Windows reference host before the corresponding behavior is implemented. Each probe SHALL be reduced to isolate one behavior and SHALL preserve its source, input, observed output or error, and relevant environment state.

#### Scenario: Encounter an unresolved semantic question
- **WHEN** an implementation group reaches behavior without decisive committed VisuAlg 3.0.7 evidence
- **THEN** that behavior remains unimplemented until a reduced probe has been recorded and added to the manifest

### Requirement: Complete oracle evidence
Each deterministic probe SHALL store the source bytes, input bytes, normalized output or normalized error, hashes of all raw observations, and any generated file bytes. A GUI-only observation SHALL additionally store a screenshot and a textual transcription sufficient for review. Normalization SHALL be versioned, reproducible, and SHALL preserve semantically significant whitespace, casing, decimal syntax, diagnostic location, and file bytes.

#### Scenario: Reproduce a deterministic probe
- **WHEN** the corpus runner processes a committed probe using its declared normalization version
- **THEN** it reproduces the committed normalized evidence and verifies every declared raw and generated-file hash

#### Scenario: Record a GUI-only result
- **WHEN** the reference behavior can only be observed through the VisuAlg interface
- **THEN** the probe contains both a screenshot and a reviewable transcription linked from the manifest

### Requirement: Requirement traceability
The conformance manifest SHALL trace every OpenSpec requirement, every `[VERIFICAR]` item, every item in the original compatibility checklist, every audited defect, every feature discovered in the official release, and every bundled `.alg` example to one or more probes and implementation tests or to an explicit, reviewed non-applicable rationale. Validation SHALL fail when a required trace is absent or stale.

#### Scenario: Detect an untraced requirement
- **WHEN** a specification requirement or required reference artifact has no current probe and test mapping
- **THEN** traceability validation fails and identifies the missing mapping

### Requirement: Bundled example acceptance
Every program bundled with the official VisuAlg 3.0.7 release SHALL be inventoried and classified as accepted, dependent on an explicit non-goal, or unusable with a recorded reason. Every accepted example SHALL parse, analyze, and run with deterministic behavior matching its recorded oracle evidence.

#### Scenario: Run the accepted example sweep
- **WHEN** final conformance validation executes all examples classified as accepted
- **THEN** every example completes with zero deterministic output, error, and generated-file mismatches

### Requirement: Conformance release gate
A release SHALL NOT claim VisuAlg 3.0.7 conformance while any deterministic oracle mismatch, untraced requirement, unexpected panic on malformed or oversized input, unpositioned runtime diagnostic, unsupported accepted example, incomplete OpenSpec task, or required quality-gate failure remains. Random compatibility SHALL cover domains, bounds, modes, and observable command effects but SHALL NOT require reproducing the reference generator's exact sequence.

#### Scenario: Qualify version 0.1.0
- **WHEN** the project prepares the `v0.1.0` release
- **THEN** all deterministic conformance, traceability, robustness, example, task, and CI gates pass and random tests verify semantics without comparing exact sequences

#### Scenario: Reject a premature release
- **WHEN** any required gate has a known failure or missing result
- **THEN** conformance sign-off and the `v0.1.0` release are blocked with the failing evidence reported
