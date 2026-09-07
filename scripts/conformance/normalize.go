package main

import (
	"bytes"
	"fmt"
)

func normalize(version string, raw []byte) ([]byte, error) {
	switch version {
	case "bytes-v1":
		return bytes.Clone(raw), nil
	case "panel-v1":
		prefix, suffix := []byte("Início da execução\r\n"), []byte("\r\nFim da execução.\r\n")
		if !bytes.HasPrefix(raw, prefix) || !bytes.HasSuffix(raw, suffix) || len(raw) < len(prefix)+len(suffix) {
			return nil, fmt.Errorf("reference panel lacks complete execution notices")
		}
		raw = raw[len(prefix) : len(raw)-len(suffix)]
	case "text-v1":
	default:
		return nil, fmt.Errorf("unknown normalizer %q", version)
	}
	return bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n")), nil
}
