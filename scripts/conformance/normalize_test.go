package main

import (
	"bytes"
	"testing"
)

func TestNormalize(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, version, raw, want string
		invalid                  bool
	}{
		{"panel", "panel-v1", "Início da execução\r\n 1.50 FALSO \r\n\r\nFim da execução.\r\n", " 1.50 FALSO \n", false},
		{"literal notice", "panel-v1", "Início da execução\r\nFim da execução.\r\n\r\nFim da execução.\r\n", "Fim da execução.\n", false},
		{"text", "text-v1", "  ERRO linha 4, coluna 2\r\n", "  ERRO linha 4, coluna 2\n", false},
		{"file bytes", "bytes-v1", "\xff\r\n\x00 ", "\xff\r\n\x00 ", false},
		{"missing prefix", "panel-v1", " 1\r\n\r\nFim da execução.\r\n", "", true},
		{"unfinished run", "panel-v1", "Início da execução\r\n 1\r\n", "", true},
		{"unknown version", "future", "anything", "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := normalize(tt.version, []byte(tt.raw))
			if tt.invalid {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
