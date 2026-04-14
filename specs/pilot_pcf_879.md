# 3GPP Specification Reference: pilot_pcf_879

## TS 29.514 — Policy Authorization Service

### Section 5.6.2.2: Supported Features

The `suppFeat` attribute is used to negotiate optional features between the NF Service Consumer and the PCF.

- **Bit 1 (InfluenceOnTrafficRouting)**: When set, the PCF supports traffic routing influence from the AF (Application Function).

When this feature is negotiated (bit 1 = 1), the PCF SHALL process `AfRoutingRequirement` if present in the `MediaComponent`. However, the `AfRoutingRequirement` (AfRoutReq) field is **OPTIONAL** — its absence does not constitute an error.

### Section 5.6.7: Traffic Routing Information Provisioning

The PCF provisions traffic routing information based on the `AfRoutingRequirement` received from the AF. The relevant fields are:
- `RouteToLocs`: routing destinations
- `UpPathChgSub`: UP path change subscription
- `AppReloc`: application relocation indication

If `AfRoutingRequirement` is not provided in the request, the PCF SHALL NOT attempt to provision traffic routing information for that media component.

## TS 29.512 — Session Management Policy Control Service

### Section 5.8: Feature Negotiation

- **Bit 1 (TrafficSteeringControl)**: Indicates SMF support for traffic steering control.

The PCF and SMF negotiate traffic steering support via `suppFeat`. When both support it (`suppFeat` bit 1 set on both sides), the PCF MAY provision traffic control data. The provisioning is conditional on the AF providing routing requirements — it is not mandatory.

## Key Implication

When traffic routing support is negotiated via suppFeat bit 1, the AfRoutingRequirement field remains optional per the specification. Its absence does not constitute an error. The PCF must handle requests where suppFeat indicates traffic routing support but no AfRoutingRequirement is provided in the media components.
