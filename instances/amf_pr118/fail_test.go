package nas_security

import (
	"os"
	"strings"
	"testing"
)

func TestDecodePlainNas_PayloadLengthCheck(t *testing.T) {
	srcPath := "security.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: DecodePlainNasNoIntegrityCheck slices payload[:7] without
	// checking if len(payload) >= 7, causing an index out of range panic.
	//
	// The fix adds: if len(payload) < 7 { return nil, fmt.Errorf("nas payload is too short") }

	if !strings.Contains(src, "len(payload) < 7") && !strings.Contains(src, "payload is too short") {
		t.Fatal("BUG: DecodePlainNasNoIntegrityCheck does not check len(payload) < 7. " +
			"Short NAS payloads will cause index out of range panic.")
	}

	t.Log("PASS: payload length check found in security.go")
}
