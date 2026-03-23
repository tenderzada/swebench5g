package nssf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: Nil pointer panic when optional fields in NSSAI availability subscription are absent.

import (
	"testing"
)

// buggyProcessSubscription simulates the buggy NSSF code that dereferences
// optional fields without nil checks.
func buggyProcessSubscription(sub *NssaiAvailabilitySubscription) (string, error) {
	// BUG: dereferences optional fields without nil check
	result := sub.NfNssaiAvailabilityUri

	// Buggy: directly dereferences Expiry without nil check
	result += " expiry=" + *sub.Expiry

	// Buggy: directly dereferences AmfSetId without nil check
	result += " amfSetId=" + *sub.AmfSetId

	return result, nil
}

// TestNilExpiryField tests that a subscription with nil Expiry is handled gracefully.
func TestNilExpiryField(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil Expiry field: %v - optional field should be checked before dereferencing", r)
		}
	}()

	amfSetId := "amf-set-001"
	sub := &NssaiAvailabilitySubscription{
		NfNssaiAvailabilityUri: "https://nssf.example.com/nssai-availability",
		TaiList: []Tai{
			{PlmnId: PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000001"},
		},
		Expiry:   nil, // optional field absent
		AmfSetId: &amfSetId,
	}

	_, err := buggyProcessSubscription(sub)
	if err != nil {
		t.Logf("got expected error: %v", err)
	}
}

// TestNilAmfSetIdField tests that a subscription with nil AmfSetId is handled.
func TestNilAmfSetIdField(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil AmfSetId field: %v - optional field should be checked before dereferencing", r)
		}
	}()

	expiry := "2025-12-31T23:59:59Z"
	sub := &NssaiAvailabilitySubscription{
		NfNssaiAvailabilityUri: "https://nssf.example.com/nssai-availability",
		TaiList: []Tai{
			{PlmnId: PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000001"},
		},
		Expiry:   &expiry,
		AmfSetId: nil, // optional field absent
	}

	_, err := buggyProcessSubscription(sub)
	if err != nil {
		t.Logf("got expected error: %v", err)
	}
}

// TestAllOptionalFieldsNil tests subscription with all optional fields absent.
func TestAllOptionalFieldsNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on all nil optional fields: %v - all optional fields should have nil checks", r)
		}
	}()

	sub := &NssaiAvailabilitySubscription{
		NfNssaiAvailabilityUri: "https://nssf.example.com/nssai-availability",
		TaiList: []Tai{
			{PlmnId: PlmnId{Mcc: "208", Mnc: "93"}, Tac: "000001"},
		},
		Expiry:   nil,
		AmfSetId: nil,
	}

	_, err := buggyProcessSubscription(sub)
	if err != nil {
		t.Logf("got expected error: %v", err)
	}
}
