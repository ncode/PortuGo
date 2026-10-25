package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestObservationAdapterRecords(t *testing.T) {
	state := filepath.Join(t.TempDir(), "state.json")
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--state", state, "testdata/records.alg"}, strings.NewReader(""), &out, &stderr)
	if status != 0 || out.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, &out, &stderr)
	}
	got, err := os.ReadFile(state)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/records.state.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("state=%s, want %s", got, want)
	}
}

func TestObservationAdapterBoundsRecordState(t *testing.T) {
	dir := t.TempDir()
	src, state := filepath.Join(dir, "source.alg"), filepath.Join(dir, "state.json")
	program := `algoritmo "large record state"
tipo Entry = registro
name: caractere
fimregistro
var rows: vetor[1..5000] de Entry
index: inteiro
inicio
para index de 1 ate 5000 faca
rows[index].name <- "` + strings.Repeat("x", 255) + `"
fimpara
fimalgoritmo`
	if err := os.WriteFile(src, []byte(program), 0600); err != nil {
		t.Fatal(err)
	}
	var out, stderr bytes.Buffer
	status := executeProbe([]string{"--max-steps", "100000", "--state", state, src}, strings.NewReader(""), &out, &stderr)
	if status != 1 || !strings.Contains(stderr.String(), errObservationSize.Error()) || out.Len() != 0 {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, &out, &stderr)
	}
	if _, err := os.Stat(state); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("oversized state was written: %v", err)
	}
}
