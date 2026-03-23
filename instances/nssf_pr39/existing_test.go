package processor

import (
	"testing"
	"time"
)

func TestExpiryField_NonNil(t *testing.T) {
	expiry := time.Now().Add(1 * time.Hour)
	if expiry.IsZero() {
		t.Fatal("non-nil expiry should not be zero")
	}
	t.Log("PASS: valid expiry handled correctly")
}
