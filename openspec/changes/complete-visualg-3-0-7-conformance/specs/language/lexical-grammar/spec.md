## Purpose

Defines how VisuAlg 3.0.7 source bytes become positioned tokens and recoverable syntax while preserving physical-line and canonical-printing behavior.

## ADDED Requirements

### Requirement: Source decoding
The language SHALL accept valid UTF-8 source, with or without a UTF-8 BOM, and Windows-1252 source as used by VisuAlg 3.0.7. Decoding SHALL occur before lexing, SHALL preserve a mapping from decoded text to original byte positions, and SHALL diagnose byte sequences that cannot be accepted under the selected rules rather than passing invalid bytes into later stages.

#### Scenario: Load a Windows-1252 program
- **WHEN** a source file contains valid Windows-1252 identifiers, keywords, comments, or string contents and is not valid UTF-8
- **THEN** it is transcoded to the same Unicode text observed in the reference while diagnostics retain correct source positions

#### Scenario: Load a UTF-8 BOM program
- **WHEN** a valid program begins with a UTF-8 BOM
- **THEN** the BOM does not become a token and all later token positions refer to the original file accurately

### Requirement: Physical newline tokens
Lexing SHALL preserve physical line boundaries after normalizing supported CRLF, LF, and reference-confirmed line endings. Newlines SHALL remain available to grammar rules that distinguish same-line and next-line forms, including when the preceding line ends with whitespace or a comment.

#### Scenario: Comment ends at a physical newline
- **WHEN** a line comment is followed by a statement on the next physical line
- **THEN** the comment does not consume the newline and the next statement is parsed at its own source position

#### Scenario: Compare line-ending encodings
- **WHEN** otherwise identical accepted programs use CRLF and LF
- **THEN** they produce equivalent token kinds and syntax trees with line and column positions appropriate to each file

### Requirement: Strict identifiers and case handling
The lexer SHALL accept exactly the oracle-confirmed VisuAlg 3.0.7 identifier character set and SHALL reject or split unsupported forms consistently with the reference. Identifier spelling SHALL be preserved for diagnostics and printing, while keyword recognition and name resolution SHALL be case-insensitive using one locale-independent canonical form.

#### Scenario: Resolve mixed-case spelling
- **WHEN** a declaration and its uses differ only in letter case
- **THEN** they denote the same symbol while diagnostics and canonical output can retain source spelling

#### Scenario: Reject an unsupported identifier form
- **WHEN** an identifier contains a leading character, continuation character, or accent form rejected by the oracle corpus
- **THEN** the implementation emits the traced lexical or syntax diagnostic and continues safely

### Requirement: Complete keyword vocabulary
The lexer SHALL recognize the complete VisuAlg 3.0.7 keyword and command vocabulary, including oracle-confirmed accented spellings, unaccented spellings, and aliases, and SHALL not reserve words that the reference accepts as identifiers. Recognition SHALL be case-insensitive without rewriting token text.

#### Scenario: Recognize an accented command spelling
- **WHEN** a program uses an oracle-confirmed accented spelling in any letter case
- **THEN** it is tokenized as the corresponding keyword or command at the correct position

#### Scenario: Preserve a non-keyword identifier
- **WHEN** a word absent from the reference keyword vocabulary is used as an identifier
- **THEN** it remains an identifier even if the current implementation previously treated it as reserved

#### Scenario: Recognize the division word alias
- **WHEN** a program uses `div` in any letter case
- **THEN** it uses the `\` operator's precedence, operand types, and left-to-right evaluation, formats as `\`, and is rejected as a variable name with `P001` on its declaration line

### Requirement: Comments and literals
The lexer SHALL accept only the reference-confirmed comment forms and string, integer, and real literal forms. It SHALL preserve literal source text and decoded value separately, apply the reference rules for escapes, delimiters, decimal syntax, and line termination, and emit positioned diagnostics for malformed or unterminated constructs.

#### Scenario: Unterminated string at end of line
- **WHEN** a string reaches a physical newline or end of file without the reference-required terminator
- **THEN** a positioned diagnostic is produced and lexing resumes or terminates according to the recorded recovery behavior without panicking

#### Scenario: Ambiguous numeric text
- **WHEN** numeric source could be interpreted as a real literal, an integer followed by punctuation, or a range boundary
- **THEN** tokenization matches the reduced oracle probe for that exact form

#### Scenario: Classify a large whole-number literal
- **WHEN** decimal digit text exceeds `2147483647` but is representable as a finite real
- **THEN** its expression has real type, including under a unary minus
- **AND** smaller digit-only literals have integer type, including forms with leading zeros

#### Scenario: Preserve an integral real literal through formatting
- **WHEN** a real literal with an integral value is formatted and parsed again
- **THEN** its real type and value are preserved

### Requirement: Complete program grammar
The parser SHALL recognize the oracle-confirmed VisuAlg 3.0.7 program header, declaration sections, subprograms, executable body, statements, and `fimalgoritmo` termination rules. Grammar variants not accepted by the reference SHALL produce diagnostics even if the earlier implementation accepted them.

#### Scenario: Parse an accepted complete program
- **WHEN** a program uses only reference-confirmed header, declaration, body, and termination forms
- **THEN** parsing returns a positioned syntax tree and no syntax diagnostics

#### Scenario: Detect text after program termination
- **WHEN** non-comment, non-whitespace source appears after `fimalgoritmo`
- **THEN** the parser accepts or diagnoses it exactly as established by the post-termination oracle probes

### Requirement: Ordered and positioned syntax
The syntax tree SHALL preserve source order for declarations, subprograms, statements, case labels, arguments, dimensions, and record fields. Every node and designator SHALL carry a source position sufficient to report an error at the construct responsible for it.

#### Scenario: Retain interleaved declaration order
- **WHEN** an accepted program contains multiple declaration section kinds in an oracle-confirmed order
- **THEN** the syntax tree exposes those declarations in source order with positions at their introducing tokens

### Requirement: Recoverable parsing
Malformed input SHALL produce collected diagnostics rather than a user-visible panic. Recovery SHALL synchronize at physical newlines and oracle-confirmed structural delimiters so that independent later errors can be reported without manufacturing valid meaning for the malformed construct.

#### Scenario: Report multiple syntax errors
- **WHEN** a file contains independent malformed declarations and statements on later lines
- **THEN** parsing returns diagnostics for each recoverable error with stable positions and does not panic or loop indefinitely

### Requirement: Canonical source printing
Canonical printing SHALL emit valid VisuAlg 3.0.7 source for every representable valid syntax tree, preserve behavior and all accepted comments, be idempotent, and support parse-print-parse equivalence modulo documented normalization of indentation, line endings, and spelling. Comment contents, order, and association with surrounding constructs SHALL survive decoding, lexing, parsing, and printing. Any reference-ignored suffix after program termination SHALL remain ignored and SHALL be retained when formatting.

#### Scenario: Canonical round trip
- **WHEN** a valid program is parsed, canonically printed, parsed again, and printed again
- **THEN** the second syntax tree is structurally equivalent to the first and both printed results are byte-identical

#### Scenario: Preserve comments around declarations and blocks
- **WHEN** a valid program contains leading, same-line, comment-only-block, pre-terminator, and trailing comments, including CP1252 text
- **THEN** formatting retains each decoded comment exactly once in the same order and associated location, including on a second formatting pass

#### Scenario: Preserve ignored trailing notes
- **WHEN** the reference ignores text following `fimalgoritmo`
- **THEN** formatting retains that suffix after the terminator without interpreting it as executable syntax
