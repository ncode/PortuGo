package source

import "testing"

func TestDecodeCP1252(t *testing.T) {
	got, err := Decode([]byte{0x4f, 0x6c, 0xe1})
	if err != nil {
		t.Fatal(err)
	}
	if got != "Olá" {
		t.Fatalf("Decode() = %q, want %q", got, "Olá")
	}
}
