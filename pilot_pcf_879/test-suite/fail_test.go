package processor

// fail_test.go — Tests that FAIL on the buggy version (panic) and
// should PASS after the fix is applied.
// These are the "fail-to-pass" tests that verify the bug is fixed.

import (
	"testing"

	"github.com/free5gc/openapi/models"
	pcf_context "github.com/free5gc/pcf/internal/context"
)

// ---------------------------------------------------------------------------
// FAIL Test 1: provisioningOfTrafficRoutingInfo must NOT panic when
// routeReq is nil (this is the core bug in Issue #879)
// ---------------------------------------------------------------------------
func TestProvisioningOfTrafficRoutingInfo_NilRouteReq(t *testing.T) {
	smPolicy := &pcf_context.UeSmPolicyData{
		PolicyDecision: &models.SmPolicyDecision{
			PccRules:      map[string]*models.PccRule{},
			TraffContDecs: map[string]*models.TrafficControlData{},
		},
	}

	// On the buggy version, this panics with:
	//   "runtime error: invalid memory address or nil pointer dereference"
	// After fix, it should return gracefully (nil or a default PccRule).
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG PRESENT: provisioningOfTrafficRoutingInfo panicked "+
				"with nil routeReq: %v\n"+
				"Fix: add nil check for routeReq parameter in "+
				"internal/sbi/processor/policyauthorization.go", r)
		}
	}()

	result := provisioningOfTrafficRoutingInfo(smPolicy, "app1", nil, "")
	_ = result
	t.Log("PASS: provisioningOfTrafficRoutingInfo handled nil routeReq without panic")
}

// ---------------------------------------------------------------------------
// FAIL Test 2: provisioningOfTrafficRoutingInfo with nil routeReq and
// existing PccRule for the appID — should not crash, should return the
// existing rule or nil gracefully.
// ---------------------------------------------------------------------------
func TestProvisioningOfTrafficRoutingInfo_NilRouteReq_WithExistingRule(t *testing.T) {
	pccRuleID := "pcc-rule-app2"
	smPolicy := &pcf_context.UeSmPolicyData{
		PolicyDecision: &models.SmPolicyDecision{
			PccRules: map[string]*models.PccRule{
				pccRuleID: {
					PccRuleId: pccRuleID,
					AppId:     "app2",
				},
			},
			TraffContDecs: map[string]*models.TrafficControlData{},
		},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG PRESENT: panicked with nil routeReq even when "+
				"a PccRule exists for the appID: %v", r)
		}
	}()

	result := provisioningOfTrafficRoutingInfo(smPolicy, "app2", nil, "")
	_ = result
	t.Log("PASS: nil routeReq with existing PccRule handled gracefully")
}

// ---------------------------------------------------------------------------
// FAIL Test 3: Simulate the exact HTTP request scenario from Issue #879.
// When suppFeat=1 and medComponents has afAppId but no AfRoutReq,
// the traffic routing provisioning path is entered with nil routeReq.
// ---------------------------------------------------------------------------
func TestSuppFeatTrafficRouting_NoAfRoutReq(t *testing.T) {
	smPolicy := &pcf_context.UeSmPolicyData{
		PolicyDecision: &models.SmPolicyDecision{
			PccRules:      map[string]*models.PccRule{},
			TraffContDecs: map[string]*models.TrafficControlData{},
			SuppFeat:      "1", // traffic routing supported
		},
	}

	// This simulates the code path:
	//   traffRoutSupp = true (suppFeat bit 1 set)
	//   medComp.AfRoutReq = nil (not provided in request)
	//   → provisioningOfTrafficRoutingInfo(smPolicy, "app1", nil, "")
	//   → PANIC on buggy version

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("BUG PRESENT: traffic routing provisioning panicked "+
				"when AfRoutReq is absent (suppFeat=1 scenario): %v\n"+
				"This is the exact crash path from Issue #879", r)
		}
	}()

	result := provisioningOfTrafficRoutingInfo(smPolicy, "app1", nil, "")
	_ = result
	t.Log("PASS: suppFeat=1 without AfRoutReq handled correctly")
}
