# PCF PR#62: smData Mispositioned - Wrong Index Causes Nil Pointer Crash

## Issue
https://github.com/free5gc/free5gc/issues/803

## Problem
When the PCF processes SM policy creation requests, it accesses SM policy subscription data using the wrong array index. This causes the code to read from an incorrect position in the array, resulting in a nil pointer dereference and crash.

## Expected Behavior
The PCF should correctly index into the SM policy data array to retrieve the appropriate subscription data for the given S-NSSAI and DNN combination.

## Actual Behavior
The PCF uses a wrong index (e.g., using the outer loop variable instead of the inner loop variable, or an off-by-one error), causing it to access the wrong element or nil element in the SM policy data array, leading to a nil pointer panic.

## Root Cause
In the SM policy processing code (likely in smpolicy.go), the wrong loop index variable is used to access the smData slice. For example, using index `i` from an outer loop when `j` from the inner loop should be used, or vice versa.

## Specification Reference
- 3GPP TS 29.512: Session Management Policy Control Service
