package processor_test

import (
	"os"
	"strings"
	"testing"
)

func TestSuciSupiMap_ExistenceCheck(t *testing.T) {
	srcPath := "ue_authentication.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: GetSupiFromSuciSupiMap is called without first checking if the
	// SUCI-SUPI pair exists, leading to nil interface type assertion panic.
	//
	// The fix adds a CheckIfSuciSupiPairExists guard before the call.

	if !strings.Contains(src, "CheckIfSuciSupiPairExists") {
		t.Fatal("BUG: no CheckIfSuciSupiPairExists guard before GetSupiFromSuciSupiMap. " +
			"Missing SUCI-SUPI mapping will cause nil interface type assertion panic.")
	}

	t.Log("PASS: CheckIfSuciSupiPairExists guard found in ue_authentication.go")
}
