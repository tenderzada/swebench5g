# NRF PR#78: Nil Pointer Crash on Missing Discovery Query Parameters

## Issue
https://github.com/free5gc/free5gc/issues/757

## Problem
The NRF crashes with a nil pointer dereference when handling NF discovery requests that have missing or invalid query parameters. The discovery handler does not validate that required query parameters are present before attempting to dereference pointer values derived from them.

## Expected Behavior
The NRF should validate query parameters and return an appropriate HTTP error (e.g., 400 Bad Request) when required parameters are missing or invalid.

## Actual Behavior
The NRF panics with a nil pointer dereference when query parameters are missing, causing the entire NRF process to crash.

## Root Cause
Missing nil checks on query parameter values before dereferencing. When optional or required query parameters are absent, the parsed values are nil pointers, and the code dereferences them without checking.

## Specification Reference
- 3GPP TS 29.510: NRF NF Discovery service
