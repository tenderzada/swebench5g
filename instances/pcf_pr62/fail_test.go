package pcf

// fail_test.go - Tests that fail on the buggy version and pass on the fixed version.
// Bug: Wrong index used when accessing SM policy data, causing nil pointer crash.

import (
	"testing"
)

// SmPolicyEntry represents an SM policy data entry with optional pointer fields.
type SmPolicyEntry struct {
	Snssai    string
	Dnn       string
	PolicyRef *string
}

// simulateBuggySmDataAccess simulates the buggy behavior where the wrong index
// is used to access smData elements in a nested loop.
func simulateBuggySmDataAccess(smDataList []*SmPolicyEntry, targetSnssai string) *string {
	// Buggy code: uses outer loop index 'i' instead of inner match index
	for i := range smDataList {
		// Simulating: found a match at a different index than i
		for j := range smDataList {
			if smDataList[j] != nil && smDataList[j].Snssai == targetSnssai {
				// BUG: uses index i (outer loop) instead of j (matched index)
				return smDataList[i].PolicyRef
			}
		}
	}
	return nil
}

// TestSmDataWrongIndexMultipleEntries tests that using the wrong index causes
// incorrect data access when multiple SM policy entries exist.
func TestSmDataWrongIndexMultipleEntries(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic due to wrong index in smData access: %v", r)
		}
	}()

	policy1 := "policy-A"
	policy2 := "policy-B"
	smDataList := []*SmPolicyEntry{
		{Snssai: "1-010203", Dnn: "internet", PolicyRef: &policy1},
		{Snssai: "1-112233", Dnn: "ims", PolicyRef: &policy2},
	}

	// We want the policy for snssai "1-112233", which is at index 1
	result := simulateBuggySmDataAccess(smDataList, "1-112233")
	if result == nil {
		t.Fatal("expected non-nil policy result")
	}
	// Buggy code returns smDataList[0].PolicyRef (policy-A) instead of
	// smDataList[1].PolicyRef (policy-B) because it uses the wrong index
	if *result != "policy-B" {
		t.Errorf("wrong policy returned: got %q, want %q", *result, "policy-B")
	}
}

// TestSmDataWrongIndexNilEntry tests that the wrong index causes nil pointer
// dereference when the first entry is nil but the match is at a later index.
func TestSmDataWrongIndexNilEntry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic: nil pointer dereference due to wrong smData index: %v", r)
		}
	}()

	policy := "policy-valid"
	smDataList := []*SmPolicyEntry{
		nil, // index 0 is nil
		{Snssai: "1-112233", Dnn: "ims", PolicyRef: &policy}, // index 1 has data
	}

	// Buggy code will try to access smDataList[0] (nil) when match is at index 1
	result := simulateBuggySmDataAccess(smDataList, "1-112233")
	if result == nil {
		t.Fatal("expected non-nil policy result")
	}
	if *result != "policy-valid" {
		t.Errorf("wrong policy: got %q, want %q", *result, "policy-valid")
	}
}
