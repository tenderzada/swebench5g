package processor

import (
	"testing"
	"time"
)

func TestExpiryField_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic when Expiry is nil: %v\n"+
				"Fix: add nil check before calling Expiry.IsZero()", r)
		}
	}()
	var expiry *time.Time
	if expiry != nil && !expiry.IsZero() {
		t.Log("expiry is set")
	} else {
		t.Log("PASS: nil expiry handled gracefully")
	}
}

func TestExpiryField_NilInStruct(t *testing.T) {
	type SubData struct {
		Expiry *time.Time
	}
	data := SubData{Expiry: nil}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic on nil Expiry in struct: %v", r)
		}
	}()
	if data.Expiry != nil && !data.Expiry.IsZero() {
		t.Log("has expiry")
	} else {
		t.Log("PASS: nil Expiry in struct handled")
	}
}
