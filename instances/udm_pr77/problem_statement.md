# UDM PR#77: Nil Pointer Dereference in SendSearchNFInstances

## Issue
https://github.com/free5gc/free5gc/issues/810

## Problem
The UDM crashes with a nil pointer dereference when calling SendSearchNFInstances and the NRF response is nil or contains unexpected data. The code does not perform nil checks on the response before attempting to access its fields.

## Expected Behavior
The UDM should gracefully handle nil or unexpected NRF responses by checking for nil before accessing response fields and returning an appropriate error.

## Actual Behavior
The UDM panics with a nil pointer dereference when the NRF returns a nil response or when the response structure is missing expected fields.

## Root Cause
Missing nil check on the NRF discovery response. The code directly dereferences the response pointer without verifying it is non-nil.

## Specification Reference
- 3GPP TS 29.510: NRF NFDiscovery service
