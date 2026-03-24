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

## Key Implication for This Bug

The code calls `GetSupiFromSuciSupiMap()` without first checking `CheckIfSuciSupiPairExists()`. When no mapping exists, the function returns nil, and the subsequent type assertion `nil.(string)` panics. The fix adds existence checks before accessing the mapping and UE context.
