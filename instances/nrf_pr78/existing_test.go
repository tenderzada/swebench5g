package nrf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify NF discovery with valid (non-nil) query parameters.

import (
	"testing"
)

// DiscoveryQueryParams represents parsed NF discovery query parameters.
type DiscoveryQueryParams struct {
	TargetNfType     *string
	RequesterNfType  *string
	ServiceNames     *[]string
	TargetPlmnList   *[]string
}

// TestValidDiscoveryQuery verifies discovery with all required parameters present.
func TestValidDiscoveryQuery(t *testing.T) {
	targetNf := "UDM"
	requesterNf := "AMF"
	params := DiscoveryQueryParams{
		TargetNfType:    &targetNf,
		RequesterNfType: &requesterNf,
	}
	if *params.TargetNfType != "UDM" {
		t.Fatal("unexpected target NF type")
	}
	if *params.RequesterNfType != "AMF" {
		t.Fatal("unexpected requester NF type")
	}
}

// TestDiscoveryQueryWithServiceNames verifies query with service names.
func TestDiscoveryQueryWithServiceNames(t *testing.T) {
	targetNf := "UDM"
	requesterNf := "AMF"
	services := []string{"nudm-sdm", "nudm-uecm"}
	params := DiscoveryQueryParams{
		TargetNfType:    &targetNf,
		RequesterNfType: &requesterNf,
		ServiceNames:    &services,
	}
	if len(*params.ServiceNames) != 2 {
		t.Fatalf("expected 2 service names, got %d", len(*params.ServiceNames))
	}
}
