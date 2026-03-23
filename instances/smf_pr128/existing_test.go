package smf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify scenarios where pointer comparison happens to work correctly.

import (
	"testing"
)

// TestSamePointerComparison verifies that comparing the same pointer to itself works.
// This passes on both buggy and fixed versions.
func TestSamePointerComparison(t *testing.T) {
	val := int32(5)
	a := &val
	b := a // same pointer

	// Even buggy pointer comparison works when both point to same address
	if a != b {
		t.Fatal("same pointer should be equal")
	}
}

// TestNilPointerComparison verifies nil pointer comparison.
func TestNilPointerComparison(t *testing.T) {
	var a *int32
	var b *int32
	// Both nil pointers are equal even with pointer comparison
	if a != b {
		t.Fatal("both nil pointers should be equal")
	}
}

// TestValueComparisonDirect verifies direct value comparison works.
func TestValueComparisonDirect(t *testing.T) {
	a := int32(10)
	b := int32(10)
	if a != b {
		t.Fatal("same values should be equal")
	}
}
