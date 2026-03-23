package pcf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify basic PCF SM policy functionality that is unaffected by the index bug.

import (
	"testing"
)

// SmPolicyData simulates SM policy subscription data.
type SmPolicyData struct {
	Snssai string
	Dnn    string
	Policy string
}

// TestSmPolicyDataCreation verifies that SM policy data structures can be created.
func TestSmPolicyDataCreation(t *testing.T) {
	smData := []SmPolicyData{
		{Snssai: "1-010203", Dnn: "internet", Policy: "default"},
	}
	if len(smData) != 1 {
		t.Fatal("expected 1 smData entry")
	}
	if smData[0].Snssai != "1-010203" {
		t.Fatal("unexpected snssai")
	}
}

// TestSmPolicyDataSingleElement verifies access with a single element.
// With only one element, the buggy index still happens to work.
func TestSmPolicyDataSingleElement(t *testing.T) {
	smData := []SmPolicyData{
		{Snssai: "1-010203", Dnn: "internet", Policy: "default"},
	}
	// When there's only one element, index 0 is always correct
	idx := 0
	if smData[idx].Policy != "default" {
		t.Fatal("unexpected policy")
	}
}
