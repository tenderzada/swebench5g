# AMF Panic on Empty RestrictedRatList

## Issue

The AMF crashes with an index-out-of-range panic when processing a UE Context Create request that contains an empty `RestrictedRatList`.

In `internal/sbi/processor/ue_context.go`, the function `CreateUEContextProcedure` accesses `RestrictedRatList[0]` without first checking whether the slice has any elements. When the list is present but empty, this causes a runtime panic.

## Root Cause

The code assigns `ue.RatType = ...RestrictedRatList[0]` unconditionally. There is no length check on the `RestrictedRatList` slice before accessing index 0. If the slice is empty (length 0), Go panics with "index out of range [0] with length 0".

## Expected Behavior

The code should check `if len(...RestrictedRatList) > 0` before accessing `RestrictedRatList[0]`. When the list is empty, the assignment should be skipped.

## References

- Issue: https://github.com/free5gc/free5gc/issues/756
- PR: https://github.com/free5gc/amf/pull/191
- Spec: 3GPP TS 29.518 (UE Context Create)
