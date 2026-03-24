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

## Key Implication for This Bug

The code checks `if keyIndex > len(suciProfiles)` but NOT `keyIndex < 1`. When `keyIndex` is 0 (reserved value per spec), `suciProfiles[keyIndex-1]` accesses index -1, causing a panic. The fix adds `keyIndex < 1 ||` to the bounds check, enforcing the specification's requirement that the identifier starts at 1.
