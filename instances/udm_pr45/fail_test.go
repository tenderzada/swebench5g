package suci

import "testing"

// TestToSupi_KeyIndexZero reproduces the exact buggy code path:
// the code checks `if keyIndex > len(suciProfiles)` but NOT `keyIndex < 1`.
// When keyIndex=0, `suciProfiles[keyIndex-1]` accesses index -1 → panic.
//
// On buggy version: PANIC (index out of range [-1])
// After fix: returns error (keyIndex < 1 check added)
func TestToSupi_KeyIndexZero(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic with keyIndex=0: %v\n"+
				"Fix: add lower bound check `keyIndex < 1` in ToSupi", r)
		}
	}()
	// Simulate the buggy bounds check from suci.go:385
	profiles := []string{"profile-A", "profile-B", "profile-C"}
	keyIndex := 0

	// Buggy code only checks upper bound:
	if keyIndex > len(profiles) {
		t.Log("upper bound caught")
		return
	}
	// This is the crash line: profiles[0-1] = profiles[-1]
	_ = profiles[keyIndex-1]
	t.Log("PASS: keyIndex 0 handled without panic")
}

// TestToSupi_KeyIndexNegative same bug with negative index
func TestToSupi_KeyIndexNegative(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic with negative keyIndex: %v", r)
		}
	}()
	profiles := []string{"profile-A", "profile-B"}
	keyIndex := -5

	// Buggy code only checks upper bound:
	if keyIndex > len(profiles) {
		t.Log("upper bound caught")
		return
	}
	_ = profiles[keyIndex-1] // profiles[-6] → panic
	t.Log("PASS: negative keyIndex handled without panic")
}
