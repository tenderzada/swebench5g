# UDM PR#66: Returns 500 Instead of 404 for Missing NSSAI

## Issue
https://github.com/free5gc/free5gc/issues/701

## Problem
When a subscriber's subscription data does not contain a particular single-NSSAI (S-NSSAI), the UDM returns HTTP 500 (Internal Server Error) instead of HTTP 404 (Not Found). This is because the code does not validate whether the requested S-NSSAI exists in the subscription data before attempting to process it.

## Expected Behavior
When the requested S-NSSAI is not found in the subscriber's subscription data, the UDM should return HTTP 404 with an appropriate error message indicating the resource was not found.

## Actual Behavior
The UDM encounters an error (likely nil pointer or empty data access) when the S-NSSAI is missing, which bubbles up as an HTTP 500 Internal Server Error.

## Root Cause
Missing validation to check whether the requested S-NSSAI exists in the subscription data. The code proceeds to process the data without checking for its existence, leading to an unhandled error condition that results in a 500 response.

## Specification Reference
- 3GPP TS 29.503: UDM Subscriber Data Management (SDM) service
