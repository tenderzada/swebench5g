package processor

import (
	"os"
	"strings"
	"testing"
)

func TestDeletePolicies_MissingReturnAfterErrorResponse(t *testing.T) {
	data, err := os.ReadFile("ampolicy.go")
	if err != nil {
		t.Fatalf("failed to read ampolicy.go: %v", err)
	}
	src := string(data)

	// The bug: after sending error response for "polAssoId not found",
	// there is no `return`, so code falls through to delete(ue.AMPolicyData, ...)
	// causing nil pointer dereference.
	//
	// Buggy code:
	//   c.JSON(int(problemDetails.Status), problemDetails)
	//   }                    ← missing return here
	//   delete(ue.AMPolicyData, polAssoId)  ← crashes
	//
	// Fixed code:
	//   c.JSON(int(problemDetails.Status), problemDetails)
	//   return               ← added
	//   }

	// Find the "polAssoId not found" block
	notFoundIdx := strings.Index(src, "polAssoId not found")
	if notFoundIdx == -1 {
		t.Fatal("could not find 'polAssoId not found' in ampolicy.go")
	}

	// Look at the 500 chars after "polAssoId not found"
	// In fixed version: there must be a `return` before `delete(`
	after := src[notFoundIdx:]
	if len(after) > 500 {
		after = after[:500]
	}

	deleteIdx := strings.Index(after, "delete(")
	returnIdx := strings.Index(after, "return")

	if returnIdx == -1 || (deleteIdx != -1 && returnIdx > deleteIdx) {
		t.Fatal("BUG: missing 'return' after error response for 'polAssoId not found'.\n" +
			"Code falls through to delete(ue.AMPolicyData, polAssoId) causing nil pointer dereference.\n" +
			"Fix: add 'return' after c.JSON error response in HandleDeletePoliciesPolAssoId")
	}

	t.Log("PASS: 'return' found between error response and delete() in ampolicy.go")
}
