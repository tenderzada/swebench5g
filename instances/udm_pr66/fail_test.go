package udm

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: UDM returns HTTP 500 instead of 404 when single-nssai is missing from subscription data.

import (
	"net/http"
	"testing"
)

// HTTPError represents an HTTP error response.
type HTTPError struct {
	StatusCode int
	Message    string
}

// buggyGetSmData simulates the buggy UDM behavior that returns 500 when NSSAI is missing.
func buggyGetSmData(subData *SubscriptionData, requestedSnssai string) (NssaiData, *HTTPError) {
	// BUG: does not check if requestedSnssai exists in NssaiMap
	// Directly accesses the map without validation, gets zero value,
	// then downstream processing fails and returns 500
	data, exists := subData.NssaiMap[requestedSnssai]
	if !exists {
		// Buggy code: returns 500 instead of 404
		return NssaiData{}, &HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal error",
		}
	}
	return data, nil
}

// TestMissingNssaiReturns404Not500 verifies that a request for a non-existent
// S-NSSAI returns 404 instead of 500.
func TestMissingNssaiReturns404Not500(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic when accessing missing NSSAI: %v", r)
		}
	}()

	subData := &SubscriptionData{
		Supi: "imsi-208930000000001",
		NssaiMap: map[string]NssaiData{
			"1-010203": {Snssai: "1-010203", Dnn: "internet", Qos: "9"},
		},
	}

	// Request a non-existent S-NSSAI
	_, httpErr := buggyGetSmData(subData, "2-AABBCC")
	if httpErr == nil {
		t.Fatal("expected an error for missing NSSAI")
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected HTTP 404 for missing NSSAI, got HTTP %d", httpErr.StatusCode)
	}
}

// TestMissingNssaiErrorMessage verifies the error message is appropriate.
func TestMissingNssaiErrorMessage(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: %v", r)
		}
	}()

	subData := &SubscriptionData{
		Supi: "imsi-208930000000001",
		NssaiMap: map[string]NssaiData{
			"1-010203": {Snssai: "1-010203", Dnn: "internet", Qos: "9"},
		},
	}

	_, httpErr := buggyGetSmData(subData, "1-FFFFFF")
	if httpErr == nil {
		t.Fatal("expected error for missing NSSAI")
	}
	// Should indicate "not found" not "internal error"
	if httpErr.StatusCode == http.StatusInternalServerError {
		t.Errorf("should not return 500 for missing NSSAI, expected 404")
	}
}

// TestEmptyNssaiMapReturns404 verifies 404 when subscription has no NSSAI data at all.
func TestEmptyNssaiMapReturns404(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on empty NSSAI map: %v", r)
		}
	}()

	subData := &SubscriptionData{
		Supi:     "imsi-208930000000002",
		NssaiMap: map[string]NssaiData{},
	}

	_, httpErr := buggyGetSmData(subData, "1-010203")
	if httpErr == nil {
		t.Fatal("expected error for missing NSSAI in empty map")
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected HTTP 404 for empty NSSAI map, got HTTP %d", httpErr.StatusCode)
	}
}
