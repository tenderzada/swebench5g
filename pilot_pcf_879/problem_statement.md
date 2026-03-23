# SWE-Bench 5G Pilot Task: PCF PolicyAuthorization Nil Pointer Panic

## Source
- **Project**: free5gc/pcf (5G Policy Control Function)
- **Issue**: [free5gc/free5gc#879](https://github.com/free5gc/free5gc/issues/879)
- **Fix PR**: [free5gc/pcf#65](https://github.com/free5gc/pcf/pull/65)
- **Difficulty**: Easy (single file, <10 lines change)

## Bug Description

The PCF's `POST /npcf-policyauthorization/v1/app-sessions` endpoint crashes with a
**nil pointer dereference** when receiving a valid authenticated request that:
- Sets `suppFeat="1"` (enabling traffic-routing support)
- But omits the `AfRoutReq` field in `medComponents`

The function `provisioningOfTrafficRoutingInfo()` in
`internal/sbi/processor/policyauthorization.go` dereferences `routeReq.RouteToLocs`,
`routeReq.UpPathChgSub`, and `routeReq.AppReloc` without checking if `routeReq` is nil.

## 3GPP Specification Reference

- **3GPP TS 29.514**: Policy Authorization Service (Section on suppFeat bit 1 = InfluenceOnTrafficRouting)
- **3GPP TS 29.512**: Session Management Policy Control (suppFeat bit 1 = Traffic Steering Control)

When both PCF and SMF negotiate traffic-routing support via `suppFeat`, the PCF should
handle the case where `AfRoutingRequirement` is absent gracefully, not crash.

## Reproduction

```bash
# Step 1: Obtain OAuth token
curl -sS -X POST 'http://<nrf-ip>:8000/oauth2/token' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data 'grant_type=client_credentials&nfType=NEF&nfInstanceId=<uuid>&targetNfType=PCF&scope=npcf-policyauthorization'

# Step 2: Trigger crash (suppFeat=1, no AfRoutReq)
curl -i -X POST 'http://<pcf-ip>:8000/npcf-policyauthorization/v1/app-sessions' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <TOKEN>' \
  --data '{
    "ascReqData": {
      "suppFeat": "1",
      "notifUri": "http://127.0.0.1:9999/appsess",
      "ueIpv4": "10.60.0.3",
      "dnn": "internet",
      "medComponents": {
        "1": {
          "medCompN": 1,
          "afAppId": "app1"
        }
      }
    }
  }'
```

## Expected Behavior
- `suppFeat="0"` (same payload): Returns `201 Created` (works correctly)
- `suppFeat="1"` (no AfRoutReq): Should return `201 Created` or appropriate error, NOT crash

## Actual Behavior
```
panic: runtime error: invalid memory address or nil pointer dereference
github.com/free5gc/pcf/internal/sbi/processor.provisioningOfTrafficRoutingInfo
  policyauthorization.go:1740
```

## Affected Function

```
File: internal/sbi/processor/policyauthorization.go
Function: provisioningOfTrafficRoutingInfo(smPolicy, appID, routeReq, fStatus)
```

The `routeReq` parameter (type `*models.AfRoutingRequirement`) is nil when
`medComponents` does not include `AfRoutReq`, but the function unconditionally
accesses its fields.

## Task for Agent

Fix the nil pointer dereference so that the PCF handles `suppFeat=1` requests
gracefully when `AfRoutReq` is absent. Ensure existing functionality is preserved.
