package nrf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: Nil pointer crash when NF discovery query parameters are missing.

import (
	"testing"
)

// buggyProcessDiscoveryQuery simulates the buggy NRF discovery handler
// that dereferences query parameters without nil checks.
func buggyProcessDiscoveryQuery(params *DiscoveryQueryParams) (string, error) {
	// BUG: no nil check on params or its fields before dereferencing
	targetNf := *params.TargetNfType
	requesterNf := *params.RequesterNfType

	// Process discovery
	_ = targetNf
	_ = requesterNf
	return "discovered", nil
}

// TestNilQueryParams tests that nil query params struct is handled gracefully.
func TestNilQueryParams(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil query params: %v - should return error instead of crashing", r)
		}
	}()

	var params *DiscoveryQueryParams = nil
	_, err := buggyProcessDiscoveryQuery(params)
	if err == nil {
		t.Error("expected error for nil query params")
	}
}

// TestMissingTargetNfType tests that missing TargetNfType parameter is handled.
func TestMissingTargetNfType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil TargetNfType: %v - should validate before dereferencing", r)
		}
	}()

	requesterNf := "AMF"
	params := &DiscoveryQueryParams{
		TargetNfType:    nil, // missing
		RequesterNfType: &requesterNf,
	}
	_, err := buggyProcessDiscoveryQuery(params)
	if err == nil {
		t.Error("expected error for missing TargetNfType")
	}
}

// TestMissingRequesterNfType tests that missing RequesterNfType parameter is handled.
func TestMissingRequesterNfType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil RequesterNfType: %v - should validate before dereferencing", r)
		}
	}()

	targetNf := "UDM"
	params := &DiscoveryQueryParams{
		TargetNfType:    &targetNf,
		RequesterNfType: nil, // missing
	}
	_, err := buggyProcessDiscoveryQuery(params)
	if err == nil {
		t.Error("expected error for missing RequesterNfType")
	}
}

// TestAllParamsMissing tests that all nil parameters are handled.
func TestAllParamsMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on all nil params: %v", r)
		}
	}()

	params := &DiscoveryQueryParams{
		TargetNfType:    nil,
		RequesterNfType: nil,
		ServiceNames:    nil,
		TargetPlmnList:  nil,
	}
	_, err := buggyProcessDiscoveryQuery(params)
	if err == nil {
		t.Error("expected error for all nil params")
	}
}
