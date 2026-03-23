# AMF Panic: Nil Pointer in Uplink RAN Configuration Transfer

## Bug Description

The AMF crashes with a **nil pointer dereference** when handling an NGAP
UplinkRANConfigurationTransfer message where `targetRanNodeID.GNbId` is nil.

The function `handleUplinkRANConfigurationTransferMain` in `internal/ngap/handler.go`
accesses `targetRanNodeID.GNbId.GNBValue` without checking if `GNbId` is nil first.

## Panic

```
panic: runtime error: invalid memory address or nil pointer dereference
```

at the line: `if targetRanNodeID.GNbId.GNBValue != "" {`

## 3GPP Reference

- TS 38.413: NGAP protocol - GNB-ID is an optional IE in Target RAN Node ID

## Task

Add a nil check for `targetRanNodeID.GNbId` before accessing its fields.
