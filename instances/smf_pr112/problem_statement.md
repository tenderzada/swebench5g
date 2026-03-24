# SMF Crash When No CHF Running (Medium Difficulty)

## Bug Description

The SMF crashes when it attempts to negotiate with a CHF (Charging Function) that
is registered in the NRF but not actively running. Two separate nil pointer
dereferences occur:

1. **CHFSelection** (`internal/sbi/consumer/nrf_service.go`): When NRF returns
   zero CHF instances, the code accesses `rsp.NfInstances[0]` without checking
   the slice length.

2. **SendConvergedChargingRequest** (`internal/sbi/consumer/chf_service.go`):
   When `SelectedCHFProfile.NfServices` is nil (CHF has no services), the code
   panics trying to iterate over nil services.

## Reproduction

1. Register a CHF in the NRF but do not start the CHF process
2. Create a PDU session that requires charging
3. SMF attempts CHF selection and crashes

## Panic Log

```
panic: runtime error: index out of range [0] with length 0
goroutine ... [running]:
github.com/free5gc/smf/internal/sbi/consumer.CHFSelection(...)
```

## 3GPP Reference

- **TS 32.291**: Converged Charging - the SMF should gracefully handle
  absence of CHF
- **TS 29.510**: NRF Discovery - empty search results are valid responses

## Why This Is Medium Difficulty

This bug requires fixes in **2 files** across the same package:
- `chf_service.go`: add nil check for NfServices + service type validation
- `nrf_service.go`: add length check for NfInstances slice

The agent must understand the relationship between NRF discovery and CHF
service selection to produce a complete fix.

## Task

Fix both nil pointer dereferences so the SMF handles missing/unavailable
CHF gracefully without crashing.
