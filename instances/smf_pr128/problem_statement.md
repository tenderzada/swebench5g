# SMF PR#128: Pointer Comparison Instead of Value Comparison

## Issue
https://github.com/free5gc/free5gc/issues/610

## Problem
The SMF compares two pointers directly instead of comparing the values they point to. In Go, comparing two pointers checks if they point to the same memory address, not whether the underlying values are equal. This causes the comparison to return false even when the values are identical, leading to incorrect logic flow.

## Expected Behavior
The SMF should compare the dereferenced values of the pointers (i.e., `*a == *b`) so that two different pointers to the same value are treated as equal.

## Actual Behavior
The SMF compares pointers directly (i.e., `a == b`), which only returns true if they point to the exact same memory address. Two separately allocated pointers with identical values are treated as not equal.

## Root Cause
A logic error where pointer variables are compared with `==` instead of dereferencing them first. This is a common Go mistake when working with pointer types.

## Specification Reference
- 3GPP TS 29.502: Session Management services
