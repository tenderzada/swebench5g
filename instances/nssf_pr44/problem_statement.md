# NSSF PR#44: Nil Pointer Panic in NSSAI Availability Subscription Handling

## Issue
https://github.com/free5gc/free5gc/issues/764

## Problem
The NSSF crashes with a nil pointer panic when handling NSSAI availability subscriptions where optional fields are absent. The subscription handling code dereferences optional pointer fields without first checking whether they are nil.

## Expected Behavior
The NSSF should handle NSSAI availability subscription requests gracefully when optional fields are not provided, skipping or defaulting the missing optional fields.

## Actual Behavior
The NSSF panics with a nil pointer dereference when optional fields in the NSSAI availability subscription request are absent (nil).

## Root Cause
Missing nil checks on optional fields in the NSSAI availability subscription data structure. The code assumes all fields are populated and directly dereferences pointer fields that may be nil.

## Specification Reference
- 3GPP TS 29.531: NSSF NSSAI Availability service
