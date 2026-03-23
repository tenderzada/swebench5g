package sbi_test

// fail_test.go: Demonstrates the reflect panic in HTTPAccessTokenRequest when
// an unknown form key is submitted to the /oauth2/token endpoint.
//
// On the buggy (parent) commit, the code does:
//
//   for key, values := range r.Form {
//       name := formKeyToFieldName[key]   // name is "" for unknown keys
//       reflect.ValueOf(&req).Elem().FieldByName(name).Set(...)  // PANIC
//   }
//
// reflect.ValueOf.FieldByName("") returns a zero reflect.Value, and calling
// .Set() on a zero Value panics with:
//   "reflect: call of reflect.Value.Set on zero Value"
//
// The fix adds: if name == "" { return 400 Bad Request }

import (
	"reflect"
	"testing"
)

// AccessTokenReq is a simplified version of the NRF access token request struct.
type AccessTokenReq struct {
	GrantType string
	NfType    string
	Scope     string
}

// formKeyToFieldName maps known OAuth2 form keys to struct field names.
var formKeyToFieldName = map[string]string{
	"grant_type": "GrantType",
	"nf_type":    "NfType",
	"scope":      "Scope",
}

// TestReflectFieldByNameEmptyString reproduces the panic when an unknown form
// key results in FieldByName("") being called.
func TestReflectFieldByNameEmptyString(t *testing.T) {
	req := AccessTokenReq{}

	// Simulate form data with an unknown key.
	formData := map[string]string{
		"grant_type":  "client_credentials",
		"unknown_key": "some_value", // This key is not in the mapping.
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BUG REPRODUCED: panic on reflect.FieldByName('').Set: %v", r)
		}
	}()

	for key, value := range formData {
		name := formKeyToFieldName[key] // Returns "" for "unknown_key"

		// BUG: no check for name == "" before FieldByName.
		// The fix adds: if name == "" { return 400 }
		reflect.ValueOf(&req).Elem().FieldByName(name).SetString(value)
	}
}

// TestReflectFieldByNameKnownKeys verifies that known form keys work correctly.
func TestReflectFieldByNameKnownKeys(t *testing.T) {
	req := AccessTokenReq{}

	formData := map[string]string{
		"grant_type": "client_credentials",
		"nf_type":    "AMF",
		"scope":      "nsmf-pdusession",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic on known form keys: %v", r)
		}
	}()

	for key, value := range formData {
		name := formKeyToFieldName[key]
		if name == "" {
			t.Errorf("Unexpected empty name for key '%s'", key)
			continue
		}
		reflect.ValueOf(&req).Elem().FieldByName(name).SetString(value)
	}

	if req.GrantType != "client_credentials" {
		t.Errorf("Expected GrantType 'client_credentials', got '%s'", req.GrantType)
	}
	if req.NfType != "AMF" {
		t.Errorf("Expected NfType 'AMF', got '%s'", req.NfType)
	}
	if req.Scope != "nsmf-pdusession" {
		t.Errorf("Expected Scope 'nsmf-pdusession', got '%s'", req.Scope)
	}
}
