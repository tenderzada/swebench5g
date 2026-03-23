package nssf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify NSSAI availability subscription with all fields present.

import (
	"testing"
)

// NssaiAvailabilitySubscription represents an NSSAI availability subscription.
type NssaiAvailabilitySubscription struct {
	NfNssaiAvailabilityUri string
	TaiList                []Tai
	Expiry                 *string
	AmfSetId               *string
}

// Tai represents a Tracking Area Identity.
type Tai struct {
	PlmnId PlmnId
	Tac    string
}

// PlmnId represents a PLMN Identity.
type PlmnId struct {
	Mcc string
	Mnc string
}

// TestSubscriptionWithAllFields verifies creation with all fields populated.
func TestSubscriptionWithAllFields(t *testing.T) {
	expiry := "2025-12-31T23:59:59Z"
	amfSetId := "amf-set-001"
	sub := NssaiAvailabilitySubscription{
		NfNssaiAvailabilityUri: "https://nssf.example.com/nssai-availability",
		TaiList: []Tai{
			{PlmnId: PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000001"},
		},
		Expiry:   &expiry,
		AmfSetId: &amfSetId,
	}
	if sub.NfNssaiAvailabilityUri == "" {
		t.Fatal("URI should not be empty")
	}
	if len(sub.TaiList) != 1 {
		t.Fatal("expected 1 TAI")
	}
	if *sub.Expiry != expiry {
		t.Fatal("unexpected expiry")
	}
}

// TestSubscriptionTaiList verifies TAI list handling.
func TestSubscriptionTaiList(t *testing.T) {
	sub := NssaiAvailabilitySubscription{
		NfNssaiAvailabilityUri: "https://nssf.example.com/nssai-availability",
		TaiList: []Tai{
			{PlmnId: PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000001"},
			{PlmnId: PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000002"},
		},
	}
	if len(sub.TaiList) != 2 {
		t.Fatalf("expected 2 TAIs, got %d", len(sub.TaiList))
	}
}
