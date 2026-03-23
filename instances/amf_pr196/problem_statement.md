# AMF PR#196: Missing Bounds Check for PDU Session PSI Bitmap

## Issue
https://github.com/free5gc/free5gc/issues/795

## Problem
The AMF does not validate the bounds of the PDU Session Status IE (PSI bitmap) received in NAS messages. A malformed bitmap (e.g., oversized or empty) can cause an index out of range panic, crashing the AMF.

## Expected Behavior
The AMF should validate the PSI bitmap length before processing. If the bitmap is malformed (too large, too small, or empty), the AMF should reject the message or handle it gracefully.

## Actual Behavior
The AMF panics with an "index out of range" error when processing a malformed PSI bitmap because there is no bounds checking before accessing bitmap elements.

## Root Cause
Missing bounds validation on the PDU Session Status IE bitmap. The code accesses bitmap byte positions without verifying the bitmap length is within the expected range (per TS 24.501, the PSI bitmap is typically 2 bytes for 16 PDU sessions).

## Specification Reference
- 3GPP TS 24.501: NAS protocol for 5GS, PDU Session Status IE
