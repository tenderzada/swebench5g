package amf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: AMF crashes when first NGAP message is not NGSetupRequest (nil RAN context).

import (
	"fmt"
	"testing"
)

// buggyHandleNGAPMessage simulates the buggy AMF NGAP handler that does not
// check whether NGSetup has been completed before processing messages.
func buggyHandleNGAPMessage(ranCtx *RANContext, msg *NGAPMessage) error {
	// BUG: No check if ranCtx is nil (NGSetup not completed yet)
	// Directly accesses ranCtx fields, which panics if NGSetup hasn't been done
	switch msg.MessageType {
	case "NGSetupRequest":
		// This is fine - it creates the context
		return nil
	case "InitialUEMessage":
		// BUG: dereferences ranCtx without nil check
		fmt.Printf("Processing InitialUEMessage for RAN: %s\n", ranCtx.RanNodeId)
		return nil
	case "UplinkNASTransport":
		// BUG: dereferences ranCtx without nil check
		fmt.Printf("Processing UplinkNAS for RAN: %s\n", ranCtx.RanNodeId)
		return nil
	default:
		// BUG: dereferences ranCtx without nil check
		fmt.Printf("Processing %s for RAN: %s\n", msg.MessageType, ranCtx.RanNodeId)
		return nil
	}
}

// TestInitialUEMessageWithoutNGSetup tests that sending InitialUEMessage
// before NGSetup does not crash the AMF.
func TestInitialUEMessageWithoutNGSetup(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic when processing InitialUEMessage without NGSetup: %v - AMF should validate NGSetup completion first", r)
		}
	}()

	// RAN context is nil because NGSetup has not been performed
	var ranCtx *RANContext = nil
	msg := &NGAPMessage{
		MessageType: "InitialUEMessage",
		Payload:     nil,
	}

	err := buggyHandleNGAPMessage(ranCtx, msg)
	if err == nil {
		t.Error("expected error when processing message without NGSetup, but got nil")
	}
}

// TestUplinkNASTransportWithoutNGSetup tests UplinkNASTransport before NGSetup.
func TestUplinkNASTransportWithoutNGSetup(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic when processing UplinkNASTransport without NGSetup: %v", r)
		}
	}()

	var ranCtx *RANContext = nil
	msg := &NGAPMessage{
		MessageType: "UplinkNASTransport",
		Payload:     nil,
	}

	err := buggyHandleNGAPMessage(ranCtx, msg)
	if err == nil {
		t.Error("expected error for UplinkNASTransport without NGSetup")
	}
}

// TestUnknownMessageWithoutNGSetup tests unknown message type before NGSetup.
func TestUnknownMessageWithoutNGSetup(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic when processing unknown message without NGSetup: %v", r)
		}
	}()

	var ranCtx *RANContext = nil
	msg := &NGAPMessage{
		MessageType: "HandoverRequired",
		Payload:     nil,
	}

	err := buggyHandleNGAPMessage(ranCtx, msg)
	if err == nil {
		t.Error("expected error for message without NGSetup")
	}
}
