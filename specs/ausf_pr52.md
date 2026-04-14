# 3GPP Specification Reference: ausf_pr52

## TS 29.509 — AUSF Services

### Section 6.1.3.2: UE Authentication (Nausf_UEAuthentication)

The AUSF provides authentication services via the Nausf_UEAuthentication service.

### Section 6.1.3.2.3: Resynchronization

When the UE sends a ResynchronizationInfo (containing AUTS and RAND), the AUSF SHALL:
1. Look up the existing SUCI-SUPI mapping for the UE
2. Retrieve the AUSF UE context
3. Forward the resynchronization parameters to the UDM

If no prior authentication context exists (no SUCI-SUPI mapping), the AUSF SHALL return an appropriate error response (404 Not Found with cause "USER_NOT_FOUND") rather than proceeding with the resynchronization.

### Section 6.1.3: Error Handling

The AUSF SHALL validate that required context exists before performing operations. Missing UE context SHALL result in a 404 response, not a server crash.

## Key Implication

A resynchronization request is only valid when a prior authentication context exists for the UE. If no SUCI-SUPI mapping has been established, the AUSF has no context to work with and must return a 404 error. Implementations must verify that the required authentication context exists before attempting to retrieve or operate on it.
