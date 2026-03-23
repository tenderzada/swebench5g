# AMF PR#192: Crash When First NGAP Message is Not NGSetupRequest

## Issue
https://github.com/free5gc/free5gc/issues/768

## Problem
The AMF crashes when a gNB sends an NGAP message other than NGSetupRequest as its first message on a new SCTP association. Per 3GPP TS 38.413, the NGSetup procedure must be the first NGAP procedure performed on a new SCTP association. The AMF does not validate this requirement and attempts to process the message, accessing uninitialized RAN context data which leads to a nil pointer crash.

## Expected Behavior
The AMF should validate that the first NGAP message from a gNB is an NGSetupRequest. If any other message type is received before NGSetup has been completed, the AMF should reject the message or close the SCTP association gracefully.

## Actual Behavior
The AMF crashes with a nil pointer dereference because the RAN context has not been initialized (it is normally set up during the NGSetup procedure), and the handler for the received message type expects the RAN context to exist.

## Root Cause
Missing validation that the NGSetup procedure has been completed for a given SCTP association before processing other NGAP messages. The RAN context is nil until NGSetup completes, and other handlers dereference it without checking.

## Specification Reference
- 3GPP TS 38.413: NGAP (NG Application Protocol) - Section 8.7.1: NGSetup must be the first procedure
