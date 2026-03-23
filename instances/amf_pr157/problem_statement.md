# AMF: Empty NAS PDU Bypasses Validation

## Bug Description

The function `DecodePlainNasNoIntegrityCheck` in `internal/nas/nas_security/security.go`
checks `if payload == nil` but fails to catch **empty (zero-length) payloads**.

A zero-length byte slice `[]byte{}` is not nil in Go, so it passes the nil check
but causes issues downstream when the code tries to access payload bytes.

## Root Cause

```go
// Buggy code:
if payload == nil {
    return nil, fmt.Errorf("nas payload is nil")
}
// A []byte{} passes this check but has no data
```

## Expected Fix

Change the nil check to a length check: `if len(payload) == 0`.

## 3GPP Reference

- TS 24.501: A valid NAS PDU must contain at least the Extended Protocol Discriminator octet.

## Task

Fix the payload validation to catch both nil and empty payloads.
