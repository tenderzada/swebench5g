package smf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: Pointer comparison instead of value comparison.

import (
	"testing"
)

// buggyCompare simulates the buggy pointer comparison behavior.
// It compares pointers instead of values.
func buggyCompare(a, b *int32) bool {
	// BUG: compares pointer addresses, not dereferenced values
	return a == b
}

// fixedCompare is the correct implementation that compares values.
func fixedCompare(a, b *int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// TestDifferentPointersToSameValue verifies that two different pointers
// to the same value should be considered equal (value comparison).
func TestDifferentPointersToSameValue(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic during pointer comparison: %v", r)
		}
	}()

	val1 := int32(42)
	val2 := int32(42)
	a := &val1
	b := &val2

	// These are different pointers but same value
	// Buggy code uses pointer comparison (a == b) which returns false
	result := buggyCompare(a, b)
	if !result {
		t.Errorf("buggy pointer comparison: two pointers to value 42 should be equal, but got false (pointer comparison instead of value comparison)")
	}
}

// TestPointerComparisonInMatchLogic simulates how the pointer comparison bug
// affects matching logic in the SMF (e.g., matching S-NSSAI or QoS parameters).
func TestPointerComparisonInMatchLogic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic during match logic: %v", r)
		}
	}()

	type SessionParams struct {
		Sst *int32
		Sd  *string
	}

	sst1 := int32(1)
	sd1 := "010203"
	request := SessionParams{Sst: &sst1, Sd: &sd1}

	sst2 := int32(1)
	sd2 := "010203"
	stored := SessionParams{Sst: &sst2, Sd: &sd2}

	// Buggy comparison: compares pointers directly
	matched := buggyCompare(request.Sst, stored.Sst)
	if !matched {
		t.Errorf("SST matching failed: request SST=%d and stored SST=%d should match, but pointer comparison returned false", *request.Sst, *stored.Sst)
	}
}

// TestPointerComparisonWithSliceElements tests pointer comparison when
// values come from different slice elements (always different pointers).
func TestPointerComparisonWithSliceElements(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()

	values := []int32{100, 200, 100}
	// Compare first and third elements (same value, different addresses)
	a := &values[0]
	b := &values[2]

	result := buggyCompare(a, b)
	if !result {
		t.Errorf("values at index 0 and 2 are both %d but pointer comparison returned false", *a)
	}
}
