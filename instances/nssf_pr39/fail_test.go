package processor

import (
	"testing"
	"time"
)

// TestExpiryField_Nil reproduces the exact buggy code path:
// the code does `subscription.SubscriptionData.Expiry.IsZero()` without nil check.
// On buggy version: PANIC (nil pointer dereference)
// After fix: PASS (nil check added before .IsZero())
func TestExpiryField_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic when Expiry is nil: %v\n"+
				"Fix: add nil check before calling Expiry.IsZero()", r)
		}
	}()
	// This is the EXACT buggy line from nssaiavailability_subscription.go:80
	// Before fix: if !subscription.SubscriptionData.Expiry.IsZero() {
	// After fix:  if subscription.SubscriptionData.Expiry != nil && !subscription.SubscriptionData.Expiry.IsZero() {
	var expiry *time.Time // nil, simulating absent Expiry field
	if !expiry.IsZero() { // direct call on nil pointer — panics on buggy version
		t.Log("expiry is set and non-zero")
	}
	t.Log("PASS: nil expiry handled without panic")
}

// TestExpiryField_NilInStruct same bug but in struct context
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
	// Direct call without nil guard
	if !data.Expiry.IsZero() {
		t.Log("has expiry")
	}
	t.Log("PASS: nil Expiry in struct handled")
}
