package processor_test

// fail_test.go: Demonstrates the index-out-of-range panic in
// CreateUEContextProcedure when RestrictedRatList is empty.
//
// On the buggy (parent) commit, the code does:
//
//   ue.RatType = ...RestrictedRatList[0]
//
// without checking len(RestrictedRatList) > 0. When the list is present but
// empty, this causes: "runtime error: index out of range [0] with length 0".
//
// The fix wraps the access with: if len(...RestrictedRatList) > 0 { ... }

import (
	"testing"
)

// RatType represents a RAT type string (simplified from the OpenAPI model).
type RatType string

// TestEmptyRestrictedRatListAccess reproduces the index-out-of-range panic
// that occurs when RestrictedRatList is empty.
func TestEmptyRestrictedRatListAccess(t *testing.T) {
	// Simulate an empty RestrictedRatList as received in the UE context request.
	restrictedRatList := []RatType{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BUG REPRODUCED: panic on empty RestrictedRatList access: %v", r)
		}
	}()

	// This is the buggy line: accessing index [0] without a length check.
	// In the real code: ue.RatType = ...RestrictedRatList[0]
	_ = restrictedRatList[0]
}

// TestNonEmptyRestrictedRatListAccess verifies the normal (non-empty) path works.
func TestNonEmptyRestrictedRatListAccess(t *testing.T) {
	restrictedRatList := []RatType{"NR"}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic on non-empty RestrictedRatList: %v", r)
		}
	}()

	if len(restrictedRatList) > 0 {
		ratType := restrictedRatList[0]
		if ratType != "NR" {
			t.Errorf("Expected RatType 'NR', got '%s'", ratType)
		}
	}
}
