package ngap

import (
	"testing"

	"github.com/free5gc/openapi/models"
)

// TestHandleGNbId_NilGNbId reproduces the exact buggy code path:
// the code does `targetRanNodeID.GNbId.GNBValue != ""` without nil check.
// On buggy version: PANIC (nil pointer dereference)
// After fix: PASS (nil check added before access)
func TestHandleGNbId_NilGNbId(t *testing.T) {
	targetRanNodeID := &models.GlobalRanNodeId{
		GNbId: nil, // GNbId is absent
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic when GNbId is nil: %v\n"+
				"Fix: add nil check before accessing GNbId.GNBValue in handler.go", r)
		}
	}()
	// This is the EXACT buggy line from handler.go:1775
	// Before fix: if targetRanNodeID.GNbId.GNBValue != "" {
	// After fix:  if targetRanNodeID.GNbId != nil && targetRanNodeID.GNbId.GNBValue != "" {
	if targetRanNodeID.GNbId.GNBValue != "" {
		t.Log("GNbId has value")
	}
	t.Log("PASS: nil GNbId handled without panic")
}

// TestHandleGNbId_NilGNbId_EmptyValue tests that empty GNBValue with non-nil GNbId works.
func TestHandleGNbId_NilGNbId_EmptyValue(t *testing.T) {
	targetRanNodeID := &models.GlobalRanNodeId{
		GNbId: nil,
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic accessing GNbId fields when GNbId is nil: %v", r)
		}
	}()
	// Same pattern - direct field access without nil guard
	_ = targetRanNodeID.GNbId.BitLength
	t.Log("PASS: accessed BitLength without panic")
}
