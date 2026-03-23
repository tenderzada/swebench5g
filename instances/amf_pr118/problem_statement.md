# AMF Crash on Malformed NAS Message (Short Payload)

## Bug Description

The AMF crashes when processing a NAS message with a payload shorter than 7 bytes.
The function `DecodePlainNasNoIntegrityCheck` in `internal/nas/nas_security/security.go`
attempts to strip the 7-byte security header without checking if the payload is long enough.

## Root Cause

```go
// Buggy: no length check before slicing
payload = payload[7:]  // panics if len(payload) < 7
```

## Panic

```
panic: runtime error: slice bounds out of range [7:3]
```

## 3GPP Reference

- TS 24.501: NAS security header is 7 bytes (1 EPD + 1 Security Header + 4 MAC + 1 SQN)
- A valid secured NAS PDU must have at least 7 bytes for the header

## Task

Add a length check before stripping the security header:
if `len(payload) < 7`, return an error instead of panicking.
