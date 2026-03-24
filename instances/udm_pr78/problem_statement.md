# UDM Nil Pointer in DataChangeNotification (Medium Difficulty)

## Bug Description

The UDM crashes with a nil pointer dereference when handling a
`POST /sdm-subscriptions` callback with a malformed or unexpected request.

Two issues in 2 files:

1. **DataChangeNotificationProcedure** (`internal/sbi/processor/notifier.go`):
   `UdmUeFindBySupi(supi)` returns nil when the UE context doesn't exist,
   but the code proceeds to use the nil UE pointer.

2. **HandleDataChangeNotificationToNF** (`internal/sbi/api_httpcallback.go`):
   No validation for empty `NotifyItems` or invalid SUPI format in the
   incoming request body.

## Panic Log

```
panic: runtime error: invalid memory address or nil pointer dereference
goroutine 240 [running]:
github.com/free5gc/udm/internal/sbi/processor.(*Processor).
  DataChangeNotificationProcedure(...)
  notifier.go:30 +0x101
```

## 3GPP Reference

- **TS 29.503**: UDM SDM - DataChangeNotification should return 204 when
  no action needed, not crash

## Why This Is Medium Difficulty

Requires fixes in **2 files**:
- `notifier.go`: check if UE context exists, return 204 if missing
- `api_httpcallback.go`: validate NotifyItems and SUPI format

The agent must understand the HTTP callback flow and return appropriate
status codes per 3GPP specification.

## Task

Fix both nil pointer dereferences. Return HTTP 204 when UE context is
missing, and validate input parameters in the callback handler.
