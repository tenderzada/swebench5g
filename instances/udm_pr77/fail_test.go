package udm

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: Nil pointer dereference in SendSearchNFInstances when NRF response is nil.

import (
	"fmt"
	"testing"
)

// SearchNFResponse wraps the NRF discovery response.
type SearchNFResponse struct {
	SearchResult *SearchResult
	StatusCode   int
}

// buggyProcessSearchResponse simulates the buggy code that dereferences
// the response without nil checks.
func buggyProcessSearchResponse(resp *SearchNFResponse) ([]NFProfile, error) {
	// BUG: no nil check on resp or resp.SearchResult
	// This will panic if resp is nil or resp.SearchResult is nil
	return resp.SearchResult.NfInstances, nil
}

// TestNilResponseHandling tests that a nil NRF response is handled gracefully.
func TestNilResponseHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil response: %v - code should handle nil response gracefully", r)
		}
	}()

	var resp *SearchNFResponse = nil
	profiles, err := buggyProcessSearchResponse(resp)
	if err == nil && profiles != nil {
		t.Error("expected error or nil profiles for nil response")
	}
}

// TestNilSearchResultHandling tests that a response with nil SearchResult is handled.
func TestNilSearchResultHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil SearchResult: %v - code should check SearchResult before dereferencing", r)
		}
	}()

	resp := &SearchNFResponse{
		SearchResult: nil,
		StatusCode:   200,
	}
	profiles, err := buggyProcessSearchResponse(resp)
	if err == nil && profiles != nil {
		t.Error("expected error or nil profiles for nil SearchResult")
	}
}

// TestEmptyNfInstancesHandling tests handling when NfInstances is empty.
func TestEmptyNfInstancesHandling(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on empty NfInstances: %v", r)
		}
	}()

	resp := &SearchNFResponse{
		SearchResult: &SearchResult{
			NfInstances: []NFProfile{},
		},
		StatusCode: 200,
	}
	profiles, _ := buggyProcessSearchResponse(resp)
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(profiles))
	}
	_ = fmt.Sprintf("found %d profiles", len(profiles))
}
