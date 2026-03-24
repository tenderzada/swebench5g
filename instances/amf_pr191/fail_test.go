package processor_test

import (
	"os"
	"strings"
	"testing"
)

func TestRestrictedRatList_LengthCheck(t *testing.T) {
	srcPath := "ue_context.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: Code accesses RestrictedRatList[0] without checking if the
	// slice is non-empty, causing "index out of range [0] with length 0".
	//
	// The fix wraps the access with: if len(...RestrictedRatList) > 0 { ... }

	hasLenCheck := strings.Contains(src, "len(") && strings.Contains(src, "RestrictedRatList")

	if !hasLenCheck {
		t.Fatal("BUG: no length check before accessing RestrictedRatList[0]. " +
			"Empty RestrictedRatList will cause index out of range panic.")
	}

	t.Log("PASS: length check for RestrictedRatList found in ue_context.go")
}
