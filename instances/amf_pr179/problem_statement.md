# AMF Panic: Index Out of Range in NAS GetSecurityHeaderType

## Bug Description

The AMF (Access and Mobility Management Function) crashes when receiving an
NGAP InitialUEMessage with a **malformed NAS PDU**. The function
`nas.GetSecurityHeaderType()` attempts to access `nasPdu[1]` without validating
that the NAS payload contains at least 2 bytes, causing an **index out of range** panic.

## Reproduction Steps

1. Start the 5G core network
2. Send an NGAP NGSetupRequest to the AMF via SCTP
3. Send an NGAP InitialUEMessage containing an undersized NAS PDU (< 2 bytes)

## Panic Stack Trace

```
panic: runtime error: index out of range [1] with length 1
goroutine X [running]:
github.com/free5gc/nas.GetSecurityHeaderType(...)
    nas.go:106
github.com/free5gc/amf/internal/nas/nas_security.DecodePlainNasNoIntegrityCheck(...)
    security.go:365
github.com/free5gc/amf/internal/ngap.handleInitialUEMessageMain(...)
    handler.go:453
```

## 3GPP Reference

- **3GPP TS 24.501**: NAS protocol for 5GS - Security header type is in octet 2 of NAS PDU
- A valid NAS PDU must have at least 2 octets (Extended Protocol Discriminator + Security Header Type)

## Expected Behavior

The AMF should validate the NAS PDU length before accessing individual octets.
If the PDU is too short, it should reject the message gracefully without crashing.

## Task

Fix the index out of range panic in the NAS security header type extraction.
Ensure that malformed (undersized) NAS PDUs are handled gracefully.
