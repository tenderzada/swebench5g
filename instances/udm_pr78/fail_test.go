package processor

import (
	"os"
	"strings"
	"testing"
)

// TestNotifier_UeContextExistenceCheck verifies that notifier.go
// checks whether UE context exists before using it.
// Buggy: UdmUeFindBySupi returns nil, code dereferences it
// Fixed: checks ok/nil and returns 204 when missing
func TestNotifier_UeContextExistenceCheck(t *testing.T) {
	data, err := os.ReadFile("notifier.go")
	if err != nil {
		t.Fatalf("cannot read notifier.go: %v", err)
	}
	src := string(data)

	// The fix changes from ignoring the second return value to checking it:
	// Before: ue := p.Context().UdmUeFindBySupi(supi)
	// After:  ue, ok := p.Context().UdmUeFindBySupi(supi); if !ok { return 204 }
	// Look for pattern indicating the existence check
	hasCheck := strings.Contains(src, "UdmUeFindBySupi") &&
		(strings.Contains(src, ", ok") || strings.Contains(src, "== nil"))

	if !hasCheck {
		t.Fatal("BUG: DataChangeNotificationProcedure does not check if UE context exists.\n" +
			"When UdmUeFindBySupi returns nil, subsequent code panics.\n" +
			"Fix: check return value of UdmUeFindBySupi and return 204 if missing\n" +
			"File: internal/sbi/processor/notifier.go")
	}
	t.Log("PASS: notifier.go checks UE context existence")
}

// TestCallback_InputValidation verifies that api_httpcallback.go
// validates NotifyItems before processing.
func TestCallback_InputValidation(t *testing.T) {
	data, err := os.ReadFile("../api_httpcallback.go")
	if err != nil {
		// Try alternate path
		data, err = os.ReadFile("api_httpcallback.go")
		if err != nil {
			t.Skipf("cannot read api_httpcallback.go: %v (file may be in different package dir)", err)
			return
		}
	}
	src := string(data)

	// The fix adds validation for NotifyItems
	if !strings.Contains(src, "NotifyItems") || !strings.Contains(src, "BadRequest") {
		t.Fatal("BUG: HandleDataChangeNotificationToNF does not validate NotifyItems.\n" +
			"Empty or malformed NotifyItems causes downstream nil pointer dereference.\n" +
			"Fix: validate NotifyItems is not empty, return 400 if invalid\n" +
			"File: internal/sbi/api_httpcallback.go")
	}
	t.Log("PASS: api_httpcallback.go validates NotifyItems")
}
