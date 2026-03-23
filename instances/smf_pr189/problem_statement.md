# SMF PR#189: Nil Pointer Crash on PFCP Messages Missing Mandatory IEs

## Issue
https://github.com/free5gc/free5gc/issues/814

## Problem
The SMF crashes when processing PFCP (Packet Forwarding Control Protocol) messages that are missing mandatory Information Elements (IEs). Multiple PFCP message handlers dereference IE pointers without first checking whether they are nil, leading to nil pointer panics.

## Expected Behavior
The SMF should validate that mandatory IEs are present in PFCP messages before processing them. If mandatory IEs are missing, the SMF should return a PFCP error response with an appropriate cause code.

## Actual Behavior
The SMF panics with nil pointer dereferences when mandatory IEs are missing from PFCP messages such as Session Establishment Response, Session Modification Response, or Session Deletion Response.

## Root Cause
Missing nil checks on mandatory IE fields in PFCP message structures. The PFCP message handlers directly dereference IE pointers (e.g., Cause IE, NodeID IE, SEID IE) without verifying they are present.

## Specification Reference
- 3GPP TS 29.244: Interface between the control plane and the user plane nodes (PFCP)
