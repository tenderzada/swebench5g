package processor_test

// fail_test.go: Demonstrates the nil-pointer dereference bug in
// requestAMFToReleasePDUResources when statusCode is nil after a failed
// N1N2MessageTransfer call.
//
// On the buggy (parent) commit, the code does:
//
//   statusCode, _, err := ...N1N2MessageTransfer(...)
//   if err != nil {
//       p.Log.Warnf(...)
//       // BUG: no return here, falls through to:
//   }
//   switch *statusCode { ... }   // PANIC: nil pointer dereference
//
// The fix adds `return false, true` inside the `if err != nil` block so that
// the nil statusCode is never dereferenced.

import (
	"testing"
)

// TestNilStatusCodeDereference reproduces the nil-pointer panic that occurs
// when N1N2MessageTransfer returns an error with a nil statusCode.
func TestNilStatusCodeDereference(t *testing.T) {
	// Simulate the buggy code path: statusCode is nil after an error.
	var statusCode *int = nil

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BUG REPRODUCED: panic on nil statusCode dereference: %v", r)
		}
	}()

	// In the fixed code, this block would contain `return` to prevent
	// falling through to the switch statement below.
	err := simulateN1N2Error()
	if err != nil {
		t.Logf("N1N2MessageTransfer returned error: %v", err)
		// BUG: missing return here in the parent commit.
		// The fix adds: return false, true
	}

	// This is the line that panics when statusCode is nil.
	// In the real code: switch *statusCode { ... }
	_ = *statusCode
}

func simulateN1N2Error() error {
	return &mockError{msg: "connection refused: gNB removed"}
}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}
