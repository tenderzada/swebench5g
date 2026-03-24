package processor

import (
	"os"
	"strings"
	"testing"
)

func TestDeletePolicies_MissingReturnAfterErrorResponse(t *testing.T) {
	srcPath := "ampolicy.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: HandleDeletePoliciesPolAssoId is missing a `return` after
	// c.JSON(int(problemDetails.Status), problemDetails) in the "polAssoId not found"
	// error block. This causes a nil pointer dereference by falling through.
	//
	// The fix adds a `return` statement after that c.JSON call.
	// We check that the error handling block has a return after c.JSON.

	// Find the block where problemDetails is sent as JSON response for not-found.
	// After the fix, there must be a `return` between that c.JSON line and the
	// next meaningful code. We verify by checking that within a window after
	// "polAssoId not found" or the c.JSON line, a `return` exists.

	idx := strings.Index(src, "problemDetails.Status")
	if idx == -1 {
		t.Fatal("BUG: could not find problemDetails.Status in source")
	}

	// Look at the code after the c.JSON call (within the next 200 chars).
	// In the fixed version, there should be a `return` before the next function logic.
	window := src[idx:]
	if len(window) > 300 {
		window = window[:300]
	}

	// Find c.JSON line, then check for return after it
	jsonIdx := strings.Index(window, "c.JSON")
	if jsonIdx == -1 {
		// c.JSON might be before our index, search a broader area
		broadStart := idx - 100
		if broadStart < 0 {
			broadStart = 0
		}
		window = src[broadStart : idx+300]
	}

	// After the c.JSON(...) call in the error block, the fix adds `return`.
	// Look for `return` appearing after `c.JSON` in this region.
	cjsonIdx := strings.Index(window, "c.JSON")
	if cjsonIdx == -1 {
		t.Fatal("BUG: could not find c.JSON near problemDetails.Status")
	}

	afterJSON := window[cjsonIdx:]
	// Find the end of the c.JSON line
	newlineIdx := strings.Index(afterJSON, "\n")
	if newlineIdx == -1 {
		t.Fatal("BUG: unexpected source format")
	}

	// Check the next few lines after c.JSON for a `return`
	remaining := afterJSON[newlineIdx:]
	if len(remaining) > 150 {
		remaining = remaining[:150]
	}

	if !strings.Contains(remaining, "return") {
		t.Fatal("BUG: missing 'return' after c.JSON error response in HandleDeletePoliciesPolAssoId. " +
			"Without return, code falls through and causes nil pointer dereference.")
	}

	t.Log("PASS: 'return' found after c.JSON error response")
}
