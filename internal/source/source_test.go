package source

import "testing"

func TestDecodeCP1252(t *testing.T) {
	got := Decode([]byte{0x4f, 0x6c, 0xe1})
	if got != "Olá" {
		t.Fatalf("Decode() = %q, want %q", got, "Olá")
	}
}
