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

	// The fix adds exactly: if len(payload) < 7
	// We must search for this specific pattern, NOT "payload is too short"
	// because that string already exists for other checks (line 120, 138).
	if !strings.Contains(src, "len(payload) < 7") {
		t.Fatal("BUG: DecodePlainNasNoIntegrityCheck does not check len(payload) < 7.\n" +
			"Short NAS payloads will cause slice bounds out of range panic.\n" +
			"Fix: add `if len(payload) < 7 { return nil, fmt.Errorf(...) }` before payload[7:]")
	}

	t.Log("PASS: len(payload) < 7 check found in security.go")
}
