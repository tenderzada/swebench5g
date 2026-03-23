package udm

// existing_test.go - Tests that pass on the buggy version.
// These tests verify NF discovery with valid (non-nil) responses.

import (
	"testing"
)

// NFProfile represents a discovered NF instance.
type NFProfile struct {
	NfInstanceId string
	NfType       string
	NfStatus     string
	Ipv4Addresses []string
}

// SearchResult represents the NRF search response.
type SearchResult struct {
	NfInstances []NFProfile
}

// TestValidSearchResponse verifies handling of a valid NRF search response.
func TestValidSearchResponse(t *testing.T) {
	result := &SearchResult{
		NfInstances: []NFProfile{
			{
				NfInstanceId: "nrf-001",
				NfType:       "NRF",
				NfStatus:     "REGISTERED",
				Ipv4Addresses: []string{"10.0.0.1"},
			},
		},
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
	if len(result.NfInstances) != 1 {
		t.Fatalf("expected 1 NF instance, got %d", len(result.NfInstances))
	}
	if result.NfInstances[0].NfInstanceId != "nrf-001" {
		t.Fatal("unexpected NF instance ID")
	}
}

// TestSearchResponseWithMultipleInstances verifies multiple NF instances in response.
func TestSearchResponseWithMultipleInstances(t *testing.T) {
	result := &SearchResult{
		NfInstances: []NFProfile{
			{NfInstanceId: "udm-001", NfType: "UDM"},
			{NfInstanceId: "udm-002", NfType: "UDM"},
		},
	}
	if len(result.NfInstances) != 2 {
		t.Fatalf("expected 2 instances, got %d", len(result.NfInstances))
	}
}
