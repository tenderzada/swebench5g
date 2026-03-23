package smf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify PFCP message processing when all mandatory IEs are present.

import (
	"testing"
)

// PFCPCause represents a PFCP Cause IE.
type PFCPCause struct {
	CauseValue uint8
}

// PFCPNodeID represents a PFCP Node ID IE.
type PFCPNodeID struct {
	NodeIDType  uint8
	NodeIDValue string
}

// SessionEstablishmentResponse represents a PFCP Session Establishment Response.
type SessionEstablishmentResponse struct {
	Cause  *PFCPCause
	NodeID *PFCPNodeID
	SEID   *uint64
}

// SessionModificationResponse represents a PFCP Session Modification Response.
type SessionModificationResponse struct {
	Cause *PFCPCause
	SEID  *uint64
}

// SessionDeletionResponse represents a PFCP Session Deletion Response.
type SessionDeletionResponse struct {
	Cause *PFCPCause
	SEID  *uint64
}

// TestValidSessionEstablishmentResponse verifies handling with all IEs present.
func TestValidSessionEstablishmentResponse(t *testing.T) {
	seid := uint64(12345)
	resp := SessionEstablishmentResponse{
		Cause:  &PFCPCause{CauseValue: 1}, // Request accepted
		NodeID: &PFCPNodeID{NodeIDType: 0, NodeIDValue: "10.0.0.1"},
		SEID:   &seid,
	}
	if resp.Cause.CauseValue != 1 {
		t.Fatal("expected cause value 1 (accepted)")
	}
	if resp.NodeID.NodeIDValue != "10.0.0.1" {
		t.Fatal("unexpected node ID")
	}
	if *resp.SEID != 12345 {
		t.Fatal("unexpected SEID")
	}
}

// TestValidSessionModificationResponse verifies modification response with all IEs.
func TestValidSessionModificationResponse(t *testing.T) {
	seid := uint64(67890)
	resp := SessionModificationResponse{
		Cause: &PFCPCause{CauseValue: 1},
		SEID:  &seid,
	}
	if resp.Cause.CauseValue != 1 {
		t.Fatal("expected cause value 1")
	}
}
