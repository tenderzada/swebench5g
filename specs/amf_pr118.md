# 3GPP Specification Reference: amf_pr118

## TS 24.501 — NAS Protocol for 5GS

### Section 9.2: General Message Format

A NAS message consists of:
- **Octet 1**: Extended Protocol Discriminator (EPD)
- **Octet 2**: Security Header Type (4 bits) + Spare Half Octet
- For security-protected messages: **Octets 3-6**: Message Authentication Code (MAC), **Octet 7**: Sequence Number

### Section 9.3: NAS Security Protected Message

A security-protected NAS message has the following structure (minimum 7 octets):

| Octet | Field | Length |
|-------|-------|--------|
| 1 | Extended Protocol Discriminator | 1 byte |
| 2 | Security Header Type | 1 byte |
| 3-6 | Message Authentication Code | 4 bytes |
| 7 | Sequence Number | 1 byte |
| 8+ | Plain NAS Message | variable |

### Section 4.4.4: Integrity Protection

The AMF SHALL verify the integrity of NAS messages. Before stripping the security header (7 bytes), the AMF MUST verify that the NAS PDU contains at least 7 octets. A NAS PDU shorter than 7 octets cannot be a valid security-protected message.

## Key Implication

A security-protected NAS message has a fixed 7-byte header structure. Any NAS PDU shorter than this minimum cannot be a valid security-protected message. The specification defines the header as mandatory and fixed-length, requiring implementations to validate the PDU length against this minimum before attempting to parse or strip the header.
