package ngap

import (
	"os"
	"strings"
	"testing"
)

// TestHandlerHasGNbIdNilCheck uses diff-based intent testing (inspired by
// SWE-Bench Mobile). Instead of calling the complex handler function directly,
// we verify that the SOURCE CODE contains the necessary nil check.
//
// On buggy version: handler.go has `targetRanNodeID.GNbId.GNBValue`
//                   WITHOUT nil guard → FAIL
// After fix: handler.go has `targetRanNodeID.GNbId != nil &&` → PASS
func TestHandlerHasGNbIdNilCheck(t *testing.T) {
	// Read the actual source file
	data, err := os.ReadFile("handler.go")
	if err != nil {
		t.Fatalf("cannot read handler.go: %v", err)
	}
	src := string(data)

	// The fix adds: targetRanNodeID.GNbId != nil && targetRanNodeID.GNbId.GNBValue != ""
	// The buggy version has: targetRanNodeID.GNbId.GNBValue != ""  (without nil check)
	if !strings.Contains(src, "GNbId != nil") {
		t.Fatal("BUG: handler.go accesses GNbId.GNBValue without nil check.\n" +
			"Fix: add `targetRanNodeID.GNbId != nil &&` before accessing GNbId.GNBValue\n" +
			"File: internal/ngap/handler.go, function: handleUplinkRANConfigurationTransferMain")
	}
	t.Log("PASS: handler.go contains GNbId nil check")
}

// TestHandlerGNbIdLinePattern verifies the specific fix pattern exists
func TestHandlerGNbIdLinePattern(t *testing.T) {
	data, err := os.ReadFile("handler.go")
	if err != nil {
		t.Fatalf("cannot read handler.go: %v", err)
	}
	src := string(data)

	// The fixed line should contain both nil check and value check
	// Pattern: GNbId != nil && ... GNbId.GNBValue
	lines := strings.Split(src, "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "GNbId") && strings.Contains(line, "GNBValue") {
			if strings.Contains(line, "!= nil") {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatal("BUG: no line in handler.go combines GNbId nil check with GNBValue access.\n" +
			"Expected pattern: `if targetRanNodeID.GNbId != nil && targetRanNodeID.GNbId.GNBValue != \"\"`")
	}
	t.Log("PASS: handler.go has correct nil-guarded GNbId access pattern")
}
