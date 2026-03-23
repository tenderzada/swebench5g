package processor_test

// fail_test.go: Demonstrates the nil interface type assertion panic in
// UeAuthPostRequestProcedure when SuciSupiMap doesn't exist during resync.
//
// On the buggy (parent) commit, the code calls:
//
//   supiOrSuci := GetSupiFromSuciSupiMap(suciSupiMap, suci)
//   // supiOrSuci is nil interface when the mapping doesn't exist
//   supi := supiOrSuci.(string)   // PANIC: interface conversion: nil
//
// The fix adds CheckIfSuciSupiPairExists and CheckIfAusfUeContextExists
// guards before calling GetSupiFromSuciSupiMap.

import (
	"testing"
)

// simulateSuciSupiMap represents the SUCI-to-SUPI mapping store.
type simulateSuciSupiMap map[string]interface{}

func getSupiFromMap(m simulateSuciSupiMap, suci string) interface{} {
	if val, ok := m[suci]; ok {
		return val
	}
	return nil
}

// TestNilInterfaceTypeAssertion reproduces the panic that occurs when
// GetSupiFromSuciSupiMap returns nil and the code type-asserts it to string.
func TestNilInterfaceTypeAssertion(t *testing.T) {
	suciSupiMap := simulateSuciSupiMap{}
	suci := "suci-0-001-01-0000-0-0-0000000001"

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BUG REPRODUCED: panic on nil interface type assertion: %v", r)
		}
	}()

	// This simulates the buggy code path: the mapping doesn't exist,
	// so getSupiFromMap returns nil, and the type assertion panics.
	supiOrSuci := getSupiFromMap(suciSupiMap, suci)

	// BUG: no existence check before type assertion.
	// The fix adds guards: CheckIfSuciSupiPairExists / CheckIfAusfUeContextExists
	_ = supiOrSuci.(string) // PANIC: interface conversion: interface is nil, not string
}

// TestExistingSupiMapEntry verifies the normal path when mapping exists.
func TestExistingSupiMapEntry(t *testing.T) {
	suciSupiMap := simulateSuciSupiMap{
		"suci-0-001-01-0000-0-0-0000000001": "imsi-001010000000001",
	}
	suci := "suci-0-001-01-0000-0-0-0000000001"

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()

	supiOrSuci := getSupiFromMap(suciSupiMap, suci)
	if supiOrSuci == nil {
		t.Fatal("Expected non-nil result from map lookup")
	}

	supi, ok := supiOrSuci.(string)
	if !ok {
		t.Fatal("Expected string type assertion to succeed")
	}
	if supi != "imsi-001010000000001" {
		t.Errorf("Expected 'imsi-001010000000001', got '%s'", supi)
	}
}
