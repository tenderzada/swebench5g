# SMF Segfault on N1N2MessageTransfer After gNB Removal

## Issue

When a gNB is removed and the SMF attempts to call N1N2MessageTransfer to release PDU session resources via the AMF, the SMF crashes with a nil pointer dereference.

In `internal/sbi/processor/association.go`, the function `requestAMFToReleasePDUResources` calls the AMF N1N2MessageTransfer API. When the call returns an error, `statusCode` may be nil. However, the code proceeds to execute `switch *statusCode`, which dereferences the nil pointer and causes a segfault.

## Root Cause

The `if err != nil` block after the N1N2MessageTransfer call logs the error but does not return. Execution falls through to the `switch *statusCode` statement, which panics because `statusCode` is nil when the HTTP call itself fails (as opposed to returning an HTTP error status).

## Expected Behavior

When the N1N2MessageTransfer call returns an error with a nil statusCode, the function should return early (`return false, true`) inside the `if err != nil` block, preventing the nil pointer dereference on `*statusCode`.

## References

- Issue: https://github.com/free5gc/smf/issues/121
- PR: https://github.com/free5gc/smf/pull/125
- Spec: 3GPP TS 29.518 (Namf_Communication N1N2MessageTransfer)
