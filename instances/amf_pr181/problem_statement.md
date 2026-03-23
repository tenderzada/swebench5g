# AMF PR#181: Wrong Nudm Service Name String

## Issue
https://github.com/free5gc/free5gc/issues/722

## Problem
The AMF uses an incorrect service name string when calling UDM services. According to 3GPP TS 29.503, the Nudm service names must follow a specific naming convention. The AMF has a typo or incorrect string for one of the Nudm service names, which causes NRF service discovery to fail when the AMF tries to locate the UDM instance.

## Expected Behavior
The AMF should use the correct Nudm service name as defined in 3GPP TS 29.503 so that NRF service discovery succeeds and the AMF can properly communicate with the UDM.

## Actual Behavior
Service discovery fails because the AMF sends an incorrect service name string to the NRF, and no matching UDM service is found.

## Root Cause
A logic error in the service name constant or string used when constructing the Nudm service request. The wrong service name string is used, likely a typo such as using "nudm-sdm" where "nudm-uecm" is needed, or vice versa.

## Specification Reference
- 3GPP TS 29.503: Nudm service names and operations
- 3GPP TS 29.510: NRF NFDiscovery service
