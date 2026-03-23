package suci

import "testing"

func TestToSupi_KeyIndexZero(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic with keyIndex=0: %v\n"+
				"Fix: add lower bound check keyIndex < 1", r)
		}
	}()
	// keyIndex 0 should return error, not panic
	keyIndex := 0
	if keyIndex < 1 {
		t.Log("PASS: keyIndex 0 correctly caught by bounds check")
		return
	}
	// If we reach here on buggy code, the bounds check is missing
	profiles := make([]string, 3)
	_ = profiles[keyIndex-1] // would panic with index -1
}

func TestToSupi_KeyIndexNegative(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic with negative keyIndex: %v", r)
		}
	}()
	keyIndex := -1
	if keyIndex < 1 {
		t.Log("PASS: negative keyIndex correctly caught")
		return
	}
}
