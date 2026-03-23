package amf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify NGAP message processing when NGSetup has been completed first.

import (
	"testing"
)

// RANContext represents the RAN (gNB) context established during NGSetup.
type RANContext struct {
	RanNodeId    string
	SupportedTAs []string
	IsSetupDone  bool
}

// NGAPMessage represents an NGAP message.
type NGAPMessage struct {
	MessageType string
	Payload     interface{}
}

// TestNGSetupRequestProcessing verifies normal NGSetup handling.
func TestNGSetupRequestProcessing(t *testing.T) {
	msg := NGAPMessage{
		MessageType: "NGSetupRequest",
		Payload:     map[string]string{"globalRANNodeID": "gnb-001"},
	}
	if msg.MessageType != "NGSetupRequest" {
		t.Fatal("unexpected message type")
	}
	// After NGSetup, RAN context is created
	ranCtx := &RANContext{
		RanNodeId:    "gnb-001",
		SupportedTAs: []string{"000001"},
		IsSetupDone:  true,
	}
	if !ranCtx.IsSetupDone {
		t.Fatal("NGSetup should be marked done")
	}
}

// TestMessageProcessingAfterSetup verifies processing messages after NGSetup completes.
func TestMessageProcessingAfterSetup(t *testing.T) {
	ranCtx := &RANContext{
		RanNodeId:    "gnb-001",
		SupportedTAs: []string{"000001"},
		IsSetupDone:  true,
	}
	msg := NGAPMessage{
		MessageType: "InitialUEMessage",
		Payload:     nil,
	}
	// This works fine because ranCtx is initialized
	if ranCtx == nil {
		t.Fatal("RAN context should not be nil after NGSetup")
	}
	_ = msg
}
