# UDM Panic: Index Out of Range in SUCI-to-SUPI Conversion

## Bug Description

The UDM crashes with an **index out of range** error when the Home Network
Public Key Identifier in a SUCI (Subscription Concealed Identifier) is 0.

The function `ToSupi()` in `pkg/suci/suci.go` checks
`if keyIndex > len(suciProfiles)` but does NOT check `keyIndex < 1`.
When `keyIndex` is 0, accessing `suciProfiles[keyIndex-1]` causes
`suciProfiles[-1]` which panics.

## Root Cause

```go
// Buggy:
if keyIndex > len(suciProfiles) {
    return "", fmt.Errorf("keyIndex(%d) out of range", keyIndex)
}
profile := suciProfiles[keyIndex-1]  // panics if keyIndex == 0
```

## 3GPP Reference

- TS 29.509: Home Network Public Key Identifier must be >= 1

## Task

Add a lower bound check for `keyIndex` to prevent index-out-of-range panic.
