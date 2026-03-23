# NRF Crash on Unknown Form Key in OAuth2 Token Endpoint

## Issue

The NRF crashes when an unknown form key is sent to the `/oauth2/token` endpoint. The crash occurs in `reflect.ValueOf.FieldByName("").Set` because the code attempts to set a struct field using an empty string field name derived from an unrecognized form key.

In `internal/sbi/api_accesstoken.go`, the function `HTTPAccessTokenRequest` iterates over form values and maps form keys to struct field names. When an unknown key is encountered, the mapped name resolves to an empty string. Calling `FieldByName("")` on a reflect.Value returns an invalid (zero) Value, and calling `.Set()` on it causes a panic.

## Root Cause

The code does not check whether the mapped field name is empty before calling `FieldByName(name).Set(...)`. When a form key is not recognized in the mapping, `name` becomes an empty string. `reflect.Value.FieldByName("")` returns a zero `reflect.Value`, and `.Set()` on a zero Value panics.

## Expected Behavior

The code should check `if name == ""` before calling `FieldByName`, and return a 400 Bad Request error for unrecognized form keys instead of panicking.

## References

- Issue: https://github.com/free5gc/free5gc/issues/770
- PR: https://github.com/free5gc/nrf/pull/79
- Spec: RFC 6749 (OAuth2), 3GPP TS 29.510 (NRF Access Token)
