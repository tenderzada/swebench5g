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

## Key Implication

Since gNB-ID is optional within GlobalRANNodeID, a valid TargetRANNodeID may identify an ng-eNB instead or may omit the node-type identifier entirely. Implementations must not assume that gNB-ID is always present when processing Uplink RAN Configuration Transfer messages.
