# 3GPP Specification Reference: udm_pr45

## TS 29.509 — AUSF/UDM Authentication Services

### Section 6.2.3: SUCI De-concealment

The UDM performs SUCI-to-SUPI conversion using Home Network Public Key Identifier and protection scheme.

### Section 6.2.3.2: Home Network Public Key Identifier

The `HomeNetworkPublicKeyId` identifies which public key was used for SUCI concealment:
- Valid range: **1 to N** (where N is the number of configured key profiles)
- Value **0** is reserved and SHALL NOT be used as a valid key identifier
- Negative values are invalid

### Section 6.2.3.3: Processing

The UDM SHALL:
1. Validate the Home Network Public Key Identifier
2. If the identifier is out of range (< 1 or > number of profiles): return error
3. If valid: use the corresponding key profile for de-concealment

## Key Implication

The Home Network Public Key Identifier has a valid range starting at 1. Value 0 is reserved per the specification and must not be accepted as a valid identifier. Implementations must validate the full range of the identifier, including the lower bound, before using it to index into configured key profiles.
