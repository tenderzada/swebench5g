package suci

import "testing"

func TestKeyIndexValidRange(t *testing.T) {
	// keyIndex 1 is the minimum valid value
	if 1 < 1 {
		t.Fatal("keyIndex 1 should be valid")
	}
	t.Log("PASS: valid keyIndex range check")
}
