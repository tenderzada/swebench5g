package sbi_test

import (
	"os"
	"strings"
	"testing"
)

func TestAccessToken_EmptyFieldNameCheck(t *testing.T) {
	srcPath := "api_accesstoken.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: When an unknown form key is submitted to /oauth2/token,
	// formKeyToFieldName[key] returns "" and reflect.FieldByName("").Set()
	// panics with "reflect: call of reflect.Value.Set on zero Value".
	//
	// The fix adds: if name == "" { return 400 Bad Request }

	if !strings.Contains(src, `name == ""`) {
		t.Fatal("BUG: no check for empty field name (name == \"\") in HTTPAccessTokenRequest. " +
			"Unknown form keys will cause reflect.FieldByName('').Set panic.")
	}

	t.Log("PASS: empty field name check found in api_accesstoken.go")
}
