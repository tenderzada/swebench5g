# NSSF Panic: Nil Expiry in NSSAIAvailability POST

## Bug Description

The NSSF crashes with a **nil pointer dereference** when handling a
`POST /nssai-availability/subscriptions` request where the `Expiry` field
is absent (nil).

The function `NssaiAvailabilitySubscriptionCreate` in
`internal/sbi/processor/nssaiavailability_subscription.go` calls
`subscription.SubscriptionData.Expiry.IsZero()` without checking if `Expiry` is nil.

## Root Cause

```go
// Buggy: directly calls method on potentially nil pointer
if !subscription.SubscriptionData.Expiry.IsZero() {
```

## 3GPP Reference

- TS 29.531: The `Expiry` field is **optional** in NSSAIAvailability subscriptions

## Task

Add a nil check for `Expiry` before calling `.IsZero()`.
