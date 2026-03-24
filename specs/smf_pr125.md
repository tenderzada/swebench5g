# 3GPP Specification Reference: smf_pr125

## TS 29.518 — AMF Communication Service

### Section 5.2.2.3: N1N2 Message Transfer

The `POST /namf-comm/v1/ue-contexts/{ueContextId}/n1-n2-messages` API transfers N1/N2 messages between SMF and AMF.

### Section 5.2.2.3.3: Error Handling

If the N1N2MessageTransfer fails (e.g., AMF unreachable, UE context not found, gNB disconnected), the response includes:
- An HTTP error status code
- A ProblemDetails body

The caller (SMF) SHALL handle the error gracefully. When the transfer fails, the SMF SHOULD:
1. Log the error
2. Take appropriate recovery action (e.g., mark the PDU session for local release)
3. NOT proceed with processing that depends on a successful transfer result

### Section 5.2.2.7: PDU Session Resource Release

When releasing PDU session resources, if the notification to the AMF fails, the SMF SHALL proceed with local resource cleanup. It SHALL NOT crash.

## Key Implication for This Bug

When `N1N2MessageTransfer` returns an error, `statusCode` is nil. The code logs the error but does NOT return — it falls through to `switch *statusCode` which dereferences the nil pointer. The fix adds `return false, true` inside the error handling block, conforming to the specification's requirement that the SMF handle transfer failures gracefully.
