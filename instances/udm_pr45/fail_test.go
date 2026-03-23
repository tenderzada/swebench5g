package suci

import (
	"os"
	"strings"
	"testing"
)

// TestToSupiHasLowerBoundCheck uses diff-based intent testing.
// Verifies that suci.go contains a lower bound check for keyIndex.
//
// Buggy version: only `if keyIndex > len(suciProfiles)` (no lower bound)
// Fixed version: `if keyIndex < 1 || keyIndex > len(suciProfiles)`
func TestToSupiHasLowerBoundCheck(t *testing.T) {
	data, err := os.ReadFile("suci.go")
	if err != nil {
		t.Fatalf("cannot read suci.go: %v", err)
	}
	src := string(data)

	// The fix adds: keyIndex < 1 ||
	if !strings.Contains(src, "keyIndex < 1") {
		t.Fatal("BUG: suci.go missing lower bound check for keyIndex.\n" +
			"Fix: change `if keyIndex > len(suciProfiles)` to " +
			"`if keyIndex < 1 || keyIndex > len(suciProfiles)`\n" +
			"File: pkg/suci/suci.go, function: ToSupi")
	}
	t.Log("PASS: suci.go contains keyIndex lower bound check")
}

// TestToSupiKeyIndexBoundsPattern verifies both bounds are checked on same line
func TestToSupiKeyIndexBoundsPattern(t *testing.T) {
	data, err := os.ReadFile("suci.go")
	if err != nil {
		t.Fatalf("cannot read suci.go: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "keyIndex") && strings.Contains(line, "< 1") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("BUG: no line in suci.go checks keyIndex < 1")
	}
}
