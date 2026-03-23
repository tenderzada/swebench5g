# AUSF Panic on Nil Interface During Authentication Resync

## Issue

The AUSF crashes with a nil interface type assertion panic during authentication resynchronization when the SUCI-SUPI mapping does not exist in the context.

In `internal/sbi/processor/ue_authentication.go`, the function `UeAuthPostRequestProcedure` calls `GetSupiFromSuciSupiMap` without first verifying that the SUCI-SUPI pair exists in the map and that the AUSF UE context is present. When the mapping does not exist (e.g., during a resync scenario), the returned nil interface is type-asserted, causing a panic.

## Root Cause

The code calls `GetSupiFromSuciSupiMap` directly without guarding with `CheckIfSuciSupiPairExists` and `CheckIfAusfUeContextExists`. When the SUCI-SUPI pair is not in the map, the function returns a nil interface value. A subsequent type assertion on this nil interface causes a runtime panic.

## Expected Behavior

Before calling `GetSupiFromSuciSupiMap`, the code should verify the existence of the SUCI-SUPI pair using `CheckIfSuciSupiPairExists` and verify the AUSF UE context using `CheckIfAusfUeContextExists`. If either check fails, the function should handle the error gracefully instead of panicking.

## References

- Issue: https://github.com/free5gc/free5gc/issues/778
- PR: https://github.com/free5gc/ausf/pull/52
- Spec: 3GPP TS 29.509 (AUSF UE Authentication)
