## ADDED Requirements

### Requirement: Four rejected-example completion gate
Acceptance of `close-rejected-program-gaps` SHALL require exact original-source replay and existing regression-test mappings for `bundled-24919c4a2e4c`, `bundled-49eb7f47065e`, `bundled-f251f2cd0c22`, and `bundled-f4fa0c514cf7`. Their classification as unusable examples SHALL NOT exempt them from this gate. All previously verified probes SHALL continue to pass; the two reviewed nonterminating exclusions SHALL retain their recorded scope.

#### Scenario: Reject an incomplete four-example fix
- **WHEN** any of the four named probes remains pending, lacks its implementation-test mapping, or differs in exit status, diagnostic sequence or established location, or standard-output bytes
- **THEN** acceptance of this change fails even if the accepted-example sweep and general release acceptance pass

#### Scenario: Complete all four rejection mappings
- **WHEN** all four named probes replay exactly, their regression and reduced-control tests pass, and all previously verified probes replay successfully
- **THEN** the probes can be marked verified and current progress documentation can report zero pending rejected-example mappings

### Requirement: Rejection evidence integrity
Investigation of the four rejected examples SHALL preserve committed source, stdin, raw observations, normalized evidence, and their hashes. New or conflicting reference observations SHALL be added separately with pinned provenance, exact inputs, fresh capture time, sufficient diagnostic and output evidence, and an explicit reviewed disposition. Partial GUI output SHALL NOT be substituted for a complete standard-output contract. A new observation SHALL NOT silently replace or weaken an existing expectation to match current implementation behavior.

#### Scenario: Find a conflicting repeat observation
- **WHEN** a fresh reference run differs from an existing rejection recording
- **THEN** the original evidence remains available, the conflict and new observation are reviewed explicitly, and the affected implementation remains pending until a reproducible contract is established

#### Scenario: Reduce a rejection rule
- **WHEN** a rejected example suggests a new diagnostic-order, lazy-body, call-state, or output rule
- **THEN** reduced triggering and neighboring non-triggering programs are recorded before dependent implementation
- **AND** evidence publication excludes private operational details and reference binaries
