package processor_test

import (
	"os"
	"strings"
	"testing"
)

func TestSuciSupiMap_ResyncExistenceCheck(t *testing.T) {
	srcPath := "ue_authentication.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: In the resynchronization code path, GetSupiFromSuciSupiMap
	// is called without checking if the mapping exists.
	// The fix adds: CheckIfSuciSupiPairExists(supiOrSuci)
	//
	// Note: CheckIfSuciSupiPairExists already exists for OTHER code paths
	// (eapSessionID, ConfirmationDataResponseID), so we must check for
	// the specific resync pattern: CheckIfSuciSupiPairExists(supiOrSuci)
	if !strings.Contains(src, "CheckIfSuciSupiPairExists(supiOrSuci)") {
		t.Fatal("BUG: missing CheckIfSuciSupiPairExists(supiOrSuci) guard in resync path.\n" +
			"When resynchronization request arrives without prior context,\n" +
			"GetSupiFromSuciSupiMap returns nil causing interface assertion panic.\n" +
			"Fix: add CheckIfSuciSupiPairExists(supiOrSuci) before GetSupiFromSuciSupiMap")
	}

	t.Log("PASS: resync path has CheckIfSuciSupiPairExists(supiOrSuci) guard")
}
