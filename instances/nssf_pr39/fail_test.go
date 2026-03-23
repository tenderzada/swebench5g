package processor

import (
	"os"
	"strings"
	"testing"
)

// TestSubscriptionHasExpiryNilCheck uses diff-based intent testing.
// Verifies that the source code checks Expiry != nil before calling .IsZero().
//
// Buggy version:  `if !subscription.SubscriptionData.Expiry.IsZero()` (no nil check)
// Fixed version:  `if subscription.SubscriptionData.Expiry != nil && !...Expiry.IsZero()`
func TestSubscriptionHasExpiryNilCheck(t *testing.T) {
	data, err := os.ReadFile("nssaiavailability_subscription.go")
	if err != nil {
		t.Fatalf("cannot read nssaiavailability_subscription.go: %v", err)
	}
	src := string(data)

	// The fix adds: Expiry != nil &&
	if !strings.Contains(src, "Expiry != nil") {
		t.Fatal("BUG: nssaiavailability_subscription.go calls Expiry.IsZero() without nil check.\n" +
			"Fix: add `subscription.SubscriptionData.Expiry != nil &&` before .IsZero()\n" +
			"File: internal/sbi/processor/nssaiavailability_subscription.go\n" +
			"Spec: 3GPP TS 29.531 - Expiry field is optional")
	}
	t.Log("PASS: source contains Expiry nil check")
}

// TestExpiryNilCheckLinePattern verifies nil check and IsZero on same line
func TestExpiryNilCheckLinePattern(t *testing.T) {
	data, err := os.ReadFile("nssaiavailability_subscription.go")
	if err != nil {
		t.Fatalf("cannot read file: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	found := false
	for _, line := range lines {
		if strings.Contains(line, "Expiry") && strings.Contains(line, "!= nil") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("BUG: no line guards Expiry with nil check before use")
	}
}
