package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type stagedRecording struct {
	ID        string   `json:"id"`
	Source    artifact `json:"source"`
	Input     artifact `json:"input"`
	Generated []string `json:"generated,omitempty"`
}

func writeNew(name string, data []byte) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(data)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

func prepareRecording(root string, p probe, stage string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	stage, err = filepath.Abs(stage)
	if err != nil {
		return err
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(stage))
	if err != nil {
		return err
	}
	stage = filepath.Join(parent, filepath.Base(stage))
	rel, err := filepath.Rel(root, stage)
	if err != nil {
		return err
	}
	if rel == "." || rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("raw recordings must be staged outside the repository")
	}
	source, err := readArtifact(root, p.Source)
	if err != nil {
		return err
	}
	input, err := readArtifact(root, p.Input)
	if err != nil {
		return err
	}
	if err := os.Mkdir(stage, 0o700); err != nil {
		return err
	}
	for name, data := range map[string][]byte{"source.alg": source, "input.txt": input} {
		if err := writeNew(filepath.Join(stage, name), data); err != nil {
			return err
		}
	}
	for _, file := range p.Files {
		data, err := readArtifact(root, file.Content)
		if err != nil {
			return err
		}
		name, err := safePath(stage, file.Path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(name), 0o700); err != nil {
			return err
		}
		if err := writeNew(name, data); err != nil {
			return err
		}
	}
	staged := stagedRecording{ID: p.ID, Source: artifact{Path: "source.alg", SHA256: hashBytes(source)}, Input: artifact{Path: "input.txt", SHA256: hashBytes(input)}}
	for _, file := range p.Implementation.Expected.Generated {
		if _, err := safePath(stage, file.Path); err != nil {
			return err
		}
		staged.Generated = append(staged.Generated, file.Path)
	}
	metadata, err := json.MarshalIndent(staged, "", "  ")
	if err != nil {
		return err
	}
	if err := writeNew(filepath.Join(stage, "staged.json"), append(metadata, '\n')); err != nil {
		return err
	}
	return writeNew(filepath.Join(stage, "instructions.txt"), []byte(`1. Open source.alg in the official VisuAlg 3.0.7 Windows application.
2. Verify the editor contains exactly this source; supply input.txt values in order.
3. Run the program. Record acceptance/rejection and the capture time in UTC.
4. Save unchanged output/error control text as UTF-8 raw.txt. Do not trim spaces.
5. Preserve generated files as bytes. For GUI-only results, retain screenshot.png
   and a UTF-8 transcription.txt. Record error locations without inventing columns.
6. Run capture with the recorded acceptance, time, normalizer and GUI-only flag.
7. Keep this directory private. Review text, images and metadata before copying
   necessary evidence to the corpus. Remove private names and machine paths,
   document any redaction, then recompute hashes over the published bytes.
8. Update the manifest from evidence.json, keeping implementation pending until
   a linked test and replay succeed. Never copy an executable or archive to Git.
`))
}

func captureRecording(stage string, accepted bool, capturedAt, normalizer string, guiOnly bool) (evidence, error) {
	var e evidence
	if _, err := time.Parse(time.RFC3339Nano, capturedAt); err != nil {
		return e, fmt.Errorf("invalid capture date: %w", err)
	}
	data, err := readFile(stage, "staged.json")
	if err != nil {
		return e, err
	}
	var staged stagedRecording
	if err := json.Unmarshal(data, &staged); err != nil {
		return e, err
	}
	for _, a := range []artifact{staged.Source, staged.Input} {
		if _, err := readArtifact(stage, a); err != nil {
			return e, err
		}
	}
	raw, err := readFile(stage, "raw.txt")
	if err != nil {
		return e, err
	}
	normalized, err := normalize(normalizer, raw)
	if err != nil {
		return e, err
	}
	e = evidence{State: "recorded", Accepted: &accepted, CapturedAt: capturedAt, Raw: artifact{Path: "raw.txt", SHA256: hashBytes(raw)}, Normalized: artifact{Path: "normalized.txt", SHA256: hashBytes(normalized)}, Normalizer: normalizer, GUIOnly: guiOnly}
	for _, name := range staged.Generated {
		b, err := readFile(stage, name)
		if err != nil {
			return evidence{}, err
		}
		e.Generated = append(e.Generated, generatedFile{Path: name, Content: artifact{Path: name, SHA256: hashBytes(b)}})
	}
	if guiOnly {
		for name, target := range map[string]**artifact{"screenshot.png": &e.Screenshot, "transcription.txt": &e.Transcription} {
			b, err := readFile(stage, name)
			if err != nil {
				return evidence{}, fmt.Errorf("GUI evidence requires %s", name)
			}
			*target = &artifact{Path: name, SHA256: hashBytes(b)}
		}
	}
	if err := writeNew(filepath.Join(stage, "normalized.txt"), normalized); err != nil {
		return evidence{}, err
	}
	metadata, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return evidence{}, err
	}
	if err := writeNew(filepath.Join(stage, "evidence.json"), append(metadata, '\n')); err != nil {
		return evidence{}, err
	}
	return e, nil
}
