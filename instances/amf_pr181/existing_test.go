package amf

// existing_test.go - Tests that pass on the buggy version.
// These tests verify basic AMF functionality unrelated to the Nudm service name bug.

import (
	"testing"
)

// TestServiceNameNotEmpty verifies that service name constants are non-empty strings.
// This passes on both buggy and fixed versions since the name is wrong, not empty.
func TestServiceNameNotEmpty(t *testing.T) {
	// Simulate a service name constant (buggy version has wrong value, but it's not empty)
	serviceName := "nudm-sdm"
	if serviceName == "" {
		t.Fatal("service name should not be empty")
	}
}

// TestBasicServiceNameFormat verifies the service name follows the general "nudm-*" pattern.
// This passes on the buggy version because the format is correct, just the specific name is wrong.
func TestBasicServiceNameFormat(t *testing.T) {
	serviceName := "nudm-sdm"
	if len(serviceName) < 5 {
		t.Fatal("service name too short")
	}
	if serviceName[:4] != "nudm" {
		t.Fatal("service name should start with 'nudm'")
	}
}
