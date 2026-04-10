# Pilot Experiment Summary: pcf_issue_879

**Date**: 2026-04-10
**Instance**: pcf_issue_879 — nil pointer dereference in `provisioningOfTrafficRoutingInfo()`
**Difficulty**: Easy (single file, <10 lines change)

---

## Results Overview

| # | Agent | Model | Mode | Resolve? | Patch Applied? | Key Failure |
|---|-------|-------|------|----------|---------------|-------------|
| 1 | Single-turn API | Qwen3.5-Flash | 1-shot | No | Format error | Correct diagnosis, unparseable patch |
| 2 | Aider | Qwen3.5-Flash | Agentic CLI | No | .gitignore only | Never found target source file |
| 3 | Multi-turn (5) | Kimi 128k | Feedback loop | No | 0/5 turns matched | SEARCH exact-match failure |
| 4 | Multi-turn (5) | Claude Sonnet 4.6 | Feedback loop | No | **2/5 turns applied** | Incomplete fix (tests fail) |
| 5 | Multi-turn (5) | GPT-4.1 | — | N/A | — | 403 region restriction |

---

## Failure Mode Analysis

### Stage 1: Bug Understanding
All models that saw the source code correctly identified the bug:
- **Qwen**: "routeReq can be nil when AfRoutReq is absent" ✅
- **Kimi**: Targeted correct function with nil check ✅  
- **Claude**: Targeted correct function with nil check ✅
- **Aider+Qwen**: Never reached this stage (file discovery failure) ❌

### Stage 2: Patch Generation
| Model | Generated correct fix idea? | Correct format? |
|-------|---------------------------|----------------|
| Qwen (single) | Yes | No (INSERT format parsing issues) |
| Kimi (multi) | Yes | Yes, but SEARCH content didn't match |
| Claude (multi) | Yes | **Yes (matched 2/5 turns)** |

### Stage 3: Patch Application
Only Claude Sonnet 4.6 succeeded in applying patches to the source code.

### Stage 4: Fix Completeness
Claude's applied patches compiled but didn't fully fix the bug — likely missing nil checks at multiple dereference points in the function.

---

## Key Findings

### Finding 1: Agentic Capability is the Bottleneck
> Single-turn Qwen correctly identifies the bug every time but cannot produce applicable patches. The gap is not domain knowledge — it's the ability to iteratively read, edit, compile, and test.

### Finding 2: Multi-turn Feedback Loops Help, but Aren't Sufficient
> Claude with a 5-turn feedback loop was the only model to successfully apply patches. The compile→test→feedback loop worked as designed (Turn 3→4 showed error-driven revision). But 5 turns wasn't enough to converge on a complete fix.

### Finding 3: SEARCH/REPLACE Exact Matching is a Bottleneck
> Kimi understood the bug and generated reasonable fixes, but failed 5/5 times at exact string matching. This is a harness limitation, not a model limitation. Fuzzy matching would likely unlock Kimi's potential.

### Finding 4: Agent Framework Choice Matters
> Aider + Qwen (same model as single-turn) performed worse because Aider's file discovery failed entirely. The agent framework's ability to identify relevant files is itself a critical capability.

---

## Infrastructure Notes

| API | Status on HK Server |
|-----|-------------------|
| DashScope (Qwen) | ✅ Working |
| Moonshot (Kimi) | ✅ Working |
| OpenRouter (Claude) | ⚠️ Requires SSH tunnel proxy |
| OpenRouter (GPT) | ❌ 403 region restriction |

---

## Next Steps

1. **Add fuzzy SEARCH matching** — unlock Kimi and Qwen multi-turn potential
2. **Increase max turns** — Claude might resolve with more turns (7-10)
3. **Qwen multi-turn** — run with same harness for complete comparison
4. **Scale to more instances** — validate findings beyond pilot

---

## Directory Structure

```
docs/results/
├── PILOT_SUMMARY.md              ← this file
├── aider_qwen_pilot/
│   ├── README.md
│   ├── aider_qwen_running.png
│   └── aider_qwen_result.png
├── multiturn_kimi_pilot/
│   ├── README.md
│   └── kimi_result.png
└── multiturn_claude_pilot/
    ├── README.md
    ├── claude_result.png
    └── gpt_403_error.png
```
