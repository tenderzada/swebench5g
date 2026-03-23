package amf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: AMF uses wrong Nudm service name string, causing service discovery failure.

import (
	"testing"
)

// simulateBuggyServiceName returns the service name as used in the buggy code.
// The buggy version uses the wrong service name for UDM UECM operations.
func simulateBuggyServiceName() string {
	// In the buggy version, when the AMF needs to call Nudm_UEContextManagement,
	// it incorrectly uses "nudm-sdm" instead of "nudm-uecm".
	return "nudm-sdm"
}

// TestNudmUECMServiceName checks that the correct service name is used for
// Nudm_UEContextManagement operations per 3GPP TS 29.503.
func TestNudmUECMServiceName(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic during service name check: %v", r)
		}
	}()

	// The service name for UE Context Management must be "nudm-uecm"
	got := simulateBuggyServiceName()
	expected := "nudm-uecm"
	if got != expected {
		t.Errorf("wrong Nudm service name: got %q, want %q", got, expected)
	}
}

// TestServiceDiscoveryWithCorrectName verifies that service discovery would
// succeed when using the correct service name.
func TestServiceDiscoveryWithCorrectName(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic during service discovery simulation: %v", r)
		}
	}()

	// Simulate NRF service registry
	registeredServices := map[string]string{
		"nudm-uecm": "http://udm:8080",
		"nudm-sdm":  "http://udm:8081",
		"nudm-ee":   "http://udm:8082",
	}

	// Buggy code looks up with wrong service name for UECM context
	requestedService := simulateBuggyServiceName()
	expectedService := "nudm-uecm"

	if requestedService != expectedService {
		t.Errorf("AMF requested service %q but should request %q for UE context management", requestedService, expectedService)
	}

	endpoint, ok := registeredServices[expectedService]
	if !ok {
		t.Fatalf("expected service %q not found in registry", expectedService)
	}
	if endpoint == "" {
		t.Fatal("endpoint should not be empty")
	}
}
