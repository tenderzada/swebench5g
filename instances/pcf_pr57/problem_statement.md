# PCF Panic: Nil Pointer in HandleDeletePoliciesPolAssoId

## Bug Description

The PCF crashes with a **nil pointer dereference** when attempting to delete
a policy association ID that does not exist.

The function `HandleDeletePoliciesPolAssoId` in `internal/sbi/processor/ampolicy.go`
does not check whether the UE or the policy data for the given `polAssoId` exists
before trying to delete it.

## Panic

```
panic: runtime error: invalid memory address or nil pointer dereference
```

When `ue.AMPolicyData[polAssoId]` is nil (the polAssoId doesn't exist),
subsequent operations on the nil map entry cause a crash.

## 3GPP Reference

- TS 29.507: AM Policy Control - DELETE operation should return 404 if polAssoId not found

## Task

Add a nil check for the UE and the policy data before deletion.
Return an appropriate error response if not found.
