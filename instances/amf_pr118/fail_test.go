package nas_security

import "testing"

func TestDecodePlainNas_ShortPayload(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic on short NAS payload: %v\n"+
				"Fix: add len(payload) < 7 check in DecodePlainNasNoIntegrityCheck", r)
		}
	}()
	shortPayload := []byte{0x7e, 0x00, 0x01}  // only 3 bytes
	_, err := DecodePlainNasNoIntegrityCheck(shortPayload)
	if err == nil {
		t.Fatal("BUG: short payload should return error")
	}
	t.Log("PASS: short payload correctly rejected")
}

func TestDecodePlainNas_SingleBytePayload(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic on 1-byte NAS payload: %v", r)
		}
	}()
	_, err := DecodePlainNasNoIntegrityCheck([]byte{0x7e})
	if err == nil {
		t.Fatal("BUG: 1-byte payload should return error")
	}
}
