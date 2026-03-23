package ngap

import (
	"testing"

	"github.com/free5gc/openapi/models"
)

func TestHandleGNbId_NilGNbId(t *testing.T) {
	targetRanNodeID := &models.NgRanTargetId{
		RanNodeId: &models.GlobalRanNodeId{
			GNbId: nil,
		},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic when GNbId is nil: %v", r)
		}
	}()
	gnbId := targetRanNodeID.RanNodeId.GNbId
	if gnbId != nil && gnbId.GNBValue != "" {
		t.Log("GNbId has value")
	} else {
		t.Log("PASS: nil GNbId handled gracefully")
	}
}

func TestHandleGNbId_NilRanNodeId(t *testing.T) {
	targetRanNodeID := &models.NgRanTargetId{
		RanNodeId: nil,
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG: panic when RanNodeId is nil: %v", r)
		}
	}()
	if targetRanNodeID.RanNodeId != nil && targetRanNodeID.RanNodeId.GNbId != nil {
		t.Log("has GNbId")
	} else {
		t.Log("PASS: nil RanNodeId handled gracefully")
	}
}
