package amf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify PDU session bitmap processing with valid-sized bitmaps.

import (
	"testing"
)

// parsePSIBitmap extracts active PDU session IDs from a PSI bitmap.
// The bitmap is expected to be 2 bytes (16 bits for PSI 0-15).
func parsePSIBitmap(bitmap []byte) []int {
	var activeSessions []int
	for i := 0; i < len(bitmap); i++ {
		for bit := 0; bit < 8; bit++ {
			if bitmap[i]&(1<<uint(bit)) != 0 {
				psi := i*8 + bit
				activeSessions = append(activeSessions, psi)
			}
		}
	}
	return activeSessions
}

// TestValidPSIBitmap verifies processing of a valid 2-byte PSI bitmap.
func TestValidPSIBitmap(t *testing.T) {
	// Bitmap: 0x02 0x00 = PSI 1 is active
	bitmap := []byte{0x02, 0x00}
	sessions := parsePSIBitmap(bitmap)
	if len(sessions) != 1 {
		t.Fatalf("expected 1 active session, got %d", len(sessions))
	}
	if sessions[0] != 1 {
		t.Fatalf("expected PSI 1, got PSI %d", sessions[0])
	}
}

// TestMultipleActiveSessions verifies multiple active PDU sessions.
func TestMultipleActiveSessions(t *testing.T) {
	// Bitmap: 0x06 0x00 = PSI 1 and PSI 2 are active
	bitmap := []byte{0x06, 0x00}
	sessions := parsePSIBitmap(bitmap)
	if len(sessions) != 2 {
		t.Fatalf("expected 2 active sessions, got %d", len(sessions))
	}
}
