package processor_test

import (
	"os"
	"strings"
	"testing"
)

func TestN1N2MessageTransfer_ReturnOnError(t *testing.T) {
	srcPath := "association.go"
	data, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}
	src := string(data)

	// The bug: After N1N2MessageTransfer returns an error, the code logs a warning
	// but does not return. It falls through to `switch *statusCode` which panics
	// because statusCode is nil.
	//
	// The fix adds `return false, true` inside the `if err != nil` block.

	// Find the error handling block for N1N2MessageTransfer
	idx := strings.Index(src, "N1N2MessageTransfer")
	if idx == -1 {
		t.Fatal("BUG: could not find N1N2MessageTransfer in source")
	}

	// Look at the code after N1N2MessageTransfer for the error handling block
	after := src[idx:]

	// Find the `if err != nil` block after N1N2MessageTransfer
	errIdx := strings.Index(after, "if err != nil")
	if errIdx == -1 {
		t.Fatal("BUG: could not find 'if err != nil' after N1N2MessageTransfer")
	}

	// Get the error block content (between `if err != nil {` and the next `}`)
	errBlock := after[errIdx:]
	// Find the closing brace of the if block
	braceCount := 0
	blockEnd := -1
	for i, ch := range errBlock {
		if ch == '{' {
			braceCount++
		} else if ch == '}' {
			braceCount--
			if braceCount == 0 {
				blockEnd = i
				break
			}
		}
	}

	if blockEnd == -1 {
		t.Fatal("BUG: could not parse error handling block")
	}

	errorHandlingBlock := errBlock[:blockEnd]

	if !strings.Contains(errorHandlingBlock, "return") {
		t.Fatal("BUG: missing 'return' in error handling block after N1N2MessageTransfer. " +
			"When N1N2MessageTransfer fails, code falls through to 'switch *statusCode' " +
			"causing nil pointer dereference.")
	}

	t.Log("PASS: 'return' found in N1N2MessageTransfer error handling block")
}
