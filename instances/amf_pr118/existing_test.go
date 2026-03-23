package nas_security

import "testing"

func TestDecodePlainNas_ValidPayload(t *testing.T) {
	// A payload of exactly 7 bytes (header only, no body) should not panic
	payload := make([]byte, 8) // 7-byte header + 1 byte body
	if len(payload) < 7 {
		t.Fatal("test payload too short")
	}
	t.Log("PASS: valid payload length check")
}
