package smf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: SMF crashes on PFCP messages missing mandatory IEs due to nil pointer dereference.

import (
	"testing"
)

// buggyHandleSessionEstResponse simulates the buggy handler that dereferences
// mandatory IEs without nil checks.
func buggyHandleSessionEstResponse(resp *SessionEstablishmentResponse) (uint8, error) {
	// BUG: no nil checks on mandatory IEs before dereferencing
	causeValue := resp.Cause.CauseValue
	_ = resp.NodeID.NodeIDValue
	_ = *resp.SEID
	return causeValue, nil
}

// buggyHandleSessionModResponse simulates the buggy modification handler.
func buggyHandleSessionModResponse(resp *SessionModificationResponse) (uint8, error) {
	// BUG: no nil check on Cause IE
	causeValue := resp.Cause.CauseValue
	_ = *resp.SEID
	return causeValue, nil
}

// buggyHandleSessionDelResponse simulates the buggy deletion handler.
func buggyHandleSessionDelResponse(resp *SessionDeletionResponse) (uint8, error) {
	// BUG: no nil check on Cause IE
	causeValue := resp.Cause.CauseValue
	return causeValue, nil
}

// TestNilCauseInEstablishmentResponse tests missing Cause IE in establishment response.
func TestNilCauseInEstablishmentResponse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil Cause IE in Session Establishment Response: %v - should check mandatory IEs before dereferencing", r)
		}
	}()

	seid := uint64(12345)
	resp := &SessionEstablishmentResponse{
		Cause:  nil, // mandatory IE missing
		NodeID: &PFCPNodeID{NodeIDType: 0, NodeIDValue: "10.0.0.1"},
		SEID:   &seid,
	}
	_, err := buggyHandleSessionEstResponse(resp)
	if err == nil {
		t.Error("expected error for missing Cause IE")
	}
}

// TestNilNodeIDInEstablishmentResponse tests missing NodeID IE.
func TestNilNodeIDInEstablishmentResponse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil NodeID IE: %v - should validate mandatory IEs", r)
		}
	}()

	seid := uint64(12345)
	resp := &SessionEstablishmentResponse{
		Cause:  &PFCPCause{CauseValue: 1},
		NodeID: nil, // mandatory IE missing
		SEID:   &seid,
	}
	_, err := buggyHandleSessionEstResponse(resp)
	if err == nil {
		t.Error("expected error for missing NodeID IE")
	}
}

// TestNilSEIDInEstablishmentResponse tests missing SEID IE.
func TestNilSEIDInEstablishmentResponse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil SEID IE: %v - should validate mandatory IEs", r)
		}
	}()

	resp := &SessionEstablishmentResponse{
		Cause:  &PFCPCause{CauseValue: 1},
		NodeID: &PFCPNodeID{NodeIDType: 0, NodeIDValue: "10.0.0.1"},
		SEID:   nil, // mandatory IE missing
	}
	_, err := buggyHandleSessionEstResponse(resp)
	if err == nil {
		t.Error("expected error for missing SEID IE")
	}
}

// TestAllIEsMissingInEstablishment tests completely empty establishment response.
func TestAllIEsMissingInEstablishment(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on all nil IEs: %v - should validate all mandatory IEs", r)
		}
	}()

	resp := &SessionEstablishmentResponse{
		Cause:  nil,
		NodeID: nil,
		SEID:   nil,
	}
	_, err := buggyHandleSessionEstResponse(resp)
	if err == nil {
		t.Error("expected error for all missing IEs")
	}
}

// TestNilCauseInModificationResponse tests missing Cause IE in modification response.
func TestNilCauseInModificationResponse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil Cause IE in Session Modification Response: %v", r)
		}
	}()

	seid := uint64(67890)
	resp := &SessionModificationResponse{
		Cause: nil,
		SEID:  &seid,
	}
	_, err := buggyHandleSessionModResponse(resp)
	if err == nil {
		t.Error("expected error for missing Cause IE in modification response")
	}
}

// TestNilCauseInDeletionResponse tests missing Cause IE in deletion response.
func TestNilCauseInDeletionResponse(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on nil Cause IE in Session Deletion Response: %v", r)
		}
	}()

	resp := &SessionDeletionResponse{
		Cause: nil,
		SEID:  nil,
	}
	_, err := buggyHandleSessionDelResponse(resp)
	if err == nil {
		t.Error("expected error for missing Cause IE in deletion response")
	}
}
