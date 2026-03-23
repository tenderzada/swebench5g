package ngap

import (
	"testing"

	"github.com/free5gc/openapi/models"
)

func TestHandleGNbId_Valid(t *testing.T) {
	targetRanNodeID := &models.NgRanTargetId{
		RanNodeId: &models.GlobalRanNodeId{
			GNbId: &models.GNbId{
				GNBValue: "000102",
				BitLength: 24,
			},
		},
	}
	gnbId := targetRanNodeID.RanNodeId.GNbId
	if gnbId == nil {
		t.Fatal("GNbId should not be nil")
	}
	if gnbId.GNBValue != "000102" {
		t.Fatalf("expected GNBValue '000102', got '%s'", gnbId.GNBValue)
	}
}
