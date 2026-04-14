# 3GPP Specification Reference: amf_pr157

## TS 24.501 — NAS Protocol for 5GS

### Section 9.2: General Message Format

Every NAS message MUST contain at least the Extended Protocol Discriminator (1 octet). A zero-length NAS PDU is invalid per the specification.

### Section 5.4.4.1: NAS Message Processing at the AMF

Upon receiving a NAS PDU from the AN (via NGAP InitialUEMessage), the AMF SHALL:
1. Check that the NAS PDU is not empty
2. Decode the Extended Protocol Discriminator
3. Determine the Security Header Type

If the NAS PDU is empty or malformed, the AMF SHALL discard the message and MAY log an error. It SHALL NOT crash or become unavailable.

## Key Implication

Every valid NAS message must contain at least one octet for the Extended Protocol Discriminator. A zero-length NAS PDU is explicitly invalid per the specification. Implementations must ensure that empty payloads are rejected before any field access, regardless of how emptiness is represented at the language level.
