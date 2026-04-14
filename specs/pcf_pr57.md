# 3GPP Specification Reference: pcf_pr57

## TS 29.507 — AM Policy Control Service

### Section 5.3.2.4: Delete AM Policy Association

`DELETE /policies/{polAssoId}` removes an AM policy association.

### Section 5.3.2.4.2: Processing

The PCF SHALL:
1. Look up the policy association by `polAssoId`
2. If not found: return **404 Not Found** with ProblemDetails (cause: CONTEXT_NOT_FOUND)
3. If found: delete the association and return **204 No Content**

### Section 5.7: Error Handling

Per 3GPP TS 29.500 (common API framework), when a resource is not found, the NF SHALL:
- Return the appropriate HTTP error status code
- Include a ProblemDetails structure in the response body
- **Stop processing** the request (return immediately after sending the error response)

## Key Implication

When a requested resource is not found, the NF must return the appropriate error response and immediately stop processing the request. The 3GPP common API framework in TS 29.500 requires that error handling be terminal: once an error response is sent, no further operations on the requested resource should be attempted.
