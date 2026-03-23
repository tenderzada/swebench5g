package nas_security

import "testing"

func TestDecodePlainNas_NilPayload(t *testing.T) {
	_, err := DecodePlainNasNoIntegrityCheck(nil)
	if err == nil {
		t.Fatal("expected error for nil payload")
	}
}
