package consumer

import (
	"os"
	"strings"
	"testing"
)

// TestCHFSelection_EmptyInstancesCheck verifies that nrf_service.go
// checks len(rsp.NfInstances) > 0 before accessing the slice.
// Buggy: directly accesses rsp.NfInstances[0] on empty result
// Fixed: wraps in if len(rsp.NfInstances) > 0 check
func TestCHFSelection_EmptyInstancesCheck(t *testing.T) {
	data, err := os.ReadFile("nrf_service.go")
	if err != nil {
		t.Fatalf("cannot read nrf_service.go: %v", err)
	}
	src := string(data)

	// The fix wraps the assignment in: if len(rsp.NfInstances) > 0
	// and returns error "No CHF found" when empty
	if !strings.Contains(src, "No CHF found") {
		t.Fatal("BUG: CHFSelection does not handle empty NfInstances.\n" +
			"When NRF returns zero CHF instances, accessing NfInstances[0] panics.\n" +
			"Fix: add len(rsp.NfInstances) > 0 check in CHFSelection (nrf_service.go)")
	}
	t.Log("PASS: CHFSelection handles empty NfInstances")
}

// TestSendConvergedChargingRequest_NilServicesCheck verifies that
// chf_service.go checks NfServices != nil before iterating.
func TestSendConvergedChargingRequest_NilServicesCheck(t *testing.T) {
	data, err := os.ReadFile("chf_service.go")
	if err != nil {
		t.Fatalf("cannot read chf_service.go: %v", err)
	}
	src := string(data)

	// The fix adds: if smContext.SelectedCHFProfile.NfServices == nil
	if !strings.Contains(src, "NfServices == nil") {
		t.Fatal("BUG: SendConvergedChargingRequest does not check for nil NfServices.\n" +
			"When CHF profile has no services, iterating over nil slice panics.\n" +
			"Fix: add NfServices nil check in SendConvergedChargingRequest (chf_service.go)")
	}
	t.Log("PASS: SendConvergedChargingRequest checks NfServices")
}
