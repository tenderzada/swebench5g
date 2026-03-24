# 3GPP Specification Reference: amf_pr161

## TS 38.413 — NGAP Protocol

### Section 9.2.6.9: Uplink RAN Configuration Transfer

The UplinkRANConfigurationTransfer message is sent by the gNB to the AMF to transfer RAN configuration information.

### Section 9.3.1.16: Target RAN Node ID

The TargetRANNodeID IE contains:
- **GlobalRANNodeID**: identifies the target RAN node
  - **gNB-ID** (OPTIONAL): identifies a gNB
  - **ng-eNB-ID** (OPTIONAL): identifies an ng-eNB

The `gNB-ID` field is OPTIONAL within GlobalRANNodeID. A valid TargetRANNodeID may contain an ng-eNB-ID instead of a gNB-ID, or the ID may be absent if the information is not available.

### Section 8.4.6: Uplink RAN Configuration Transfer Procedure

The AMF uses the TargetRANNodeID to route configuration to the target node. If the target node ID type is not gNB (e.g., ng-eNB or absent), the AMF SHALL handle this appropriately without assuming gNB-ID is always present.

## Key Implication for This Bug

The code accesses `targetRanNodeID.GNbId.GNBValue` without checking if `GNbId` is nil. Since gNB-ID is optional in the specification, the AMF must check for nil before accessing its fields. The fix adds `targetRanNodeID.GNbId != nil &&` before the field access.
