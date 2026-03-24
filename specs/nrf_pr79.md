# 3GPP Specification Reference: nrf_pr79

## TS 29.510 — NRF Services

### Section 5.4.2.2: Access Token Request (Nnrf_AccessToken)

The NRF provides OAuth2 access tokens via `POST /oauth2/token` with `application/x-www-form-urlencoded` content type.

### Section 5.4.2.2.3: Request Parameters

The defined form parameters are:
- `grant_type` (REQUIRED): must be "client_credentials"
- `nfInstanceId` (REQUIRED): NF instance identifier
- `nfType` (OPTIONAL): NF type
- `targetNfType` (OPTIONAL): target NF type
- `scope` (REQUIRED): requested scope
- `targetNfInstanceId` (OPTIONAL): target NF instance

### Section 5.4.2.2.4: Error Handling

Per RFC 6749 Section 5.2, if the request contains an unsupported or unknown parameter, the authorization server SHOULD ignore the unknown parameter. It SHALL NOT crash.

For malformed requests, the NRF SHALL return:
- HTTP 400 Bad Request
- Error response body: `{"error": "invalid_request"}`

## Key Implication for This Bug

The code uses reflection to map form keys to struct fields via YAML tags. When an unknown key is submitted, `FieldByName("")` returns a zero `reflect.Value`, and calling `.Set()` on it panics. The fix adds `if name == ""` check to return 400 for unknown keys, conforming to RFC 6749.
