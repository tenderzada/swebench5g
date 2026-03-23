package nas_security

import "testing"

func TestDecodePlainNas_EmptyPayload(t *testing.T) {
	_, err := DecodePlainNasNoIntegrityCheck([]byte{})
	if err == nil {
		t.Fatal("BUG: empty payload should return error but was accepted")
	}
	t.Log("PASS: empty payload correctly rejected")
}

func TestDecodePlainNas_ZeroLenPayload(t *testing.T) {
	empty := make([]byte, 0)
	_, err := DecodePlainNasNoIntegrityCheck(empty)
	if err == nil {
		t.Fatal("BUG: zero-length payload should return error")
	}
}
