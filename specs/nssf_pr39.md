# 3GPP Specification Reference: nssf_pr39

## TS 29.531 — NSSF Services

### Section 5.3.2.3: NSSAIAvailability Subscription

The `POST /nssai-availability/subscriptions` API creates a subscription for NSSAI availability notifications.

### Section 5.3.2.3.2: Request Body (NssfEventSubscriptionCreateData)

| Field | Type | Required |
|-------|------|----------|
| nfNssaiAvailabilityUri | string | REQUIRED |
| taiList | array | REQUIRED |
| event | enum | REQUIRED |
| expiry | DateTime | **OPTIONAL** |

The `expiry` field specifies when the subscription expires. It is **OPTIONAL** per the specification. If absent, the subscription does not expire (or a server-defined default applies).

### Section 5.3.2.3.3: Processing

The NSSF SHALL:
1. Validate required fields
2. If `expiry` is present AND not zero, set a timer for subscription expiry
3. If `expiry` is absent (null/nil), treat the subscription as non-expiring
4. Return 201 Created with the subscription resource

## Key Implication

The expiry field is explicitly optional in the subscription request. When a client omits this field, the NSSF must treat the subscription as non-expiring or apply a server-defined default. Implementations must account for the absence of optional fields before performing any operations on them.
