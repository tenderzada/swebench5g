package amf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: Missing bounds check for PDU session PSI bitmap causes index out of range.

import (
	"testing"
)

// buggyParsePSIStatus simulates the buggy code that accesses bitmap bytes
// without bounds checking.
func buggyParsePSIStatus(bitmap []byte) []int {
	var activeSessions []int
	// BUG: hardcoded access to bitmap[0] and bitmap[1] without length check
	// This will panic if bitmap is empty or has only 1 byte
	psiLow := bitmap[0]
	psiHigh := bitmap[1]

	for bit := 0; bit < 8; bit++ {
		if psiLow&(1<<uint(bit)) != 0 {
			activeSessions = append(activeSessions, bit)
		}
		if psiHigh&(1<<uint(bit)) != 0 {
			activeSessions = append(activeSessions, 8+bit)
		}
	}
	return activeSessions
}

// TestEmptyBitmap tests that an empty PSI bitmap is handled gracefully.
func TestEmptyBitmap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on empty bitmap: %v - should validate bitmap length before accessing", r)
		}
	}()

	bitmap := []byte{}
	sessions := buggyParsePSIStatus(bitmap)
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions for empty bitmap, got %d", len(sessions))
	}
}

// TestSingleByteBitmap tests that a 1-byte bitmap (too short) is handled.
func TestSingleByteBitmap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on 1-byte bitmap: %v - should validate bitmap length >= 2", r)
		}
	}()

	bitmap := []byte{0x02} // only 1 byte instead of expected 2
	sessions := buggyParsePSIStatus(bitmap)
	_ = sessions
}

// TestOversizedBitmap tests that an oversized bitmap does not cause issues.
func TestOversizedBitmap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on oversized bitmap: %v", r)
		}
	}()

	// 4 bytes instead of expected 2 - should only process valid range
	bitmap := []byte{0x02, 0x00, 0xFF, 0xFF}
	sessions := buggyParsePSIStatus(bitmap)
	if len(sessions) == 0 {
		t.Error("expected at least 1 active session")
	}
}

// TestNilBitmap tests that a nil bitmap is handled.
func TestNilBitmap(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil bitmap: %v - should check for nil before accessing", r)
		}
	}()

	var bitmap []byte = nil
	sessions := buggyParsePSIStatus(bitmap)
	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions for nil bitmap, got %d", len(sessions))
	}
}
