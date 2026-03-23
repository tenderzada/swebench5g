package processor

// existing_test.go — Tests that PASS on the buggy version.
// These verify that existing functionality works and must not regress after fix.

import (
	"testing"

	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
)

// ---------------------------------------------------------------------------
// Helper: build a minimal UeSmPolicyData with PccRules so the function can
// look up a PccRule by AfAppId.
// ---------------------------------------------------------------------------
func newSmPolicyWithAppRule(appID string) *pcf_context.UeSmPolicyData {
	pccRuleID := "pcc-rule-" + appID
	return &pcf_context.UeSmPolicyData{
		PolicyDecision: &models.SmPolicyDecision{
			PccRules: map[string]*models.PccRule{
				pccRuleID: {
					PccRuleId: pccRuleID,
					AppId:     appID,
				},
			},
			TraffContDecs: map[string]*models.TrafficControlData{},
		},
	}
}

// ---------------------------------------------------------------------------
// Test 1: provisioningOfTrafficRoutingInfo works when routeReq is valid
// ---------------------------------------------------------------------------
func TestProvisioningOfTrafficRoutingInfo_ValidRouteReq(t *testing.T) {
	smPolicy := newSmPolicyWithAppRule("app1")
	routeReq := &models.AfRoutingRequirement{}

	// Should not panic with a valid (non-nil) routeReq
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic with valid routeReq: %v", r)
		}
	}()

	result := provisioningOfTrafficRoutingInfo(smPolicy, "app1", routeReq, "")
	// We only verify it doesn't crash; the returned rule may or may not be nil
	// depending on internal logic, but no panic means existing behavior is OK.
	_ = result
}

// ---------------------------------------------------------------------------
// Test 2: PccRule lookup by AppId works (util function sanity check)
// ---------------------------------------------------------------------------
func TestPccRuleLookupByAppId(t *testing.T) {
	smPolicy := newSmPolicyWithAppRule("testApp")
	rules := smPolicy.PolicyDecision.PccRules

	found := false
	for _, rule := range rules {
		if rule.AppId == "testApp" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find PccRule with AppId 'testApp'")
	}
}
