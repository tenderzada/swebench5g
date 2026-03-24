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

## Key Implication for This Bug

The code calls `subscription.SubscriptionData.Expiry.IsZero()` without checking if `Expiry` is nil. Since `Expiry` is a `*time.Time` pointer and the field is optional per TS 29.531, it can be nil when not provided in the request. The fix adds `Expiry != nil &&` before calling `.IsZero()`.
