# SWE-Bench 5G Development Roadmap

## Current Status (v0.2) — March 2026

### Completed

- [x] 10 validated task instances across 7 NFs
- [x] Dataset published on HuggingFace (v0.2)
- [x] Docker-based evaluation harness with Qwen adapter
- [x] Specification-as-Skill A/B framework with 10 3GPP excerpts
- [x] Preliminary baseline: Qwen3.5-Flash (0% resolve, 100% bug detection)
- [x] Project showcase page (GitHub Pages)
- [x] Paper draft with Spec-as-Skill angle (NeurIPS D&B format)
- [x] Chinese tutorial (11 chapters)
- [x] Automated pipeline: mine (280 candidates) → build → validate → evaluate
- [x] Batch build/validate script with --clean mode
- [x] Dual test strategy: direct call + diff-based intent test

### In Progress

- [ ] Validate pcf_pr57, smf_pr112, udm_pr78
- [ ] Fix amf_pr191 (wrong parent commit)
- [ ] Validate 10 remaining instances (second batch)

---

## Phase 1: Complete Dataset & Validation (Priority: CRITICAL)

**Goal**: 20+ fully validated instances with 3GPP spec excerpts

| Task | Status | Notes |
|------|--------|-------|
| Validate pcf_pr57 | Pending | Last fix pushed, needs re-test |
| Validate smf_pr112 (medium) | Pending | Multi-file: chf_service.go + nrf_service.go |
| Validate udm_pr78 (medium) | Pending | Multi-file: notifier.go + api_httpcallback.go |
| Fix amf_pr191 parent commit | Blocked | Parent already contains fix |
| Validate 10 second-batch instances | Pending | Have PLACEHOLDER commits, need filling |
| Write specs for new validated instances | Ongoing | Template established |
| Update HuggingFace to v0.3 | After validation | |

---

## Phase 2: Agent Baseline Experiments (Priority: HIGH)

**Goal**: Fill paper Table 6 with real data across multiple agents

### 2.1 Fix Qwen Evaluation Pipeline

The single-turn API approach has fundamental limitations for code editing. Options:

| Approach | Effort | Expected Result |
|----------|--------|----------------|
| Fix output parsing (current) | Low | Still likely 0% due to edit precision |
| Use Qwen via Aider (agentic) | Medium | Should improve significantly |
| Use Qwen via OpenRouter + Claude Code style | Medium | Needs adapter |

**Decision**: Focus on agentic tools (Claude Code, Aider) rather than fixing single-turn pipeline.

### 2.2 Agentic Agent Evaluation

| Agent | Model | Status | Notes |
|-------|-------|--------|-------|
| Claude Code | Opus 4.6 | Not started | Install CLI in Docker image |
| Claude Code | Sonnet 4.6 | Not started | |
| Aider | Qwen3.5-Flash | Not started | Aider supports OpenAI-compatible APIs |
| Aider | Claude Sonnet 4.6 | Not started | |
| Codex CLI | GPT-4.1 | Not started | |

**Deliverable**: Resolve rate table for paper, broken down by agent, model, NF, difficulty.

### 2.3 Analysis Dimensions

- [ ] Per-NF resolve rate
- [ ] Per-bug-type resolve rate
- [ ] Agent behavior analysis (does it read spec references?)
- [ ] Failure mode classification

---

## Phase 3: Specification-as-Skill A/B Experiment (Priority: HIGH)

**Goal**: Answer the core research question with quantitative data

### 3.1 Run A/B Experiments

```bash
# For each agent that achieves >0% resolve rate:
python eval/run_evaluation.py --agent <agent> --model <model> --ab-test --runs 3
```

### 3.2 Expected Analysis

| Hypothesis | Metric | Expected |
|-----------|--------|----------|
| H1: Simple nil checks don't need spec | Δ resolve rate for nil-pointer bugs | ~0% |
| H2: Protocol semantics bugs need spec | Δ resolve rate for spec-dependent bugs | >0% |
| H3: Token overhead is modest | Avg token increase | <50% |

### 3.3 Paper Deliverable

- Table: A/B results per instance (with/without spec)
- Bar chart: Δ resolve rate by bug type
- Token overhead analysis

---

## Phase 4: Scale to 30+ Instances (Priority: MEDIUM)

### 4.1 Add Medium & Hard Difficulty

| Difficulty | Current | Target | Source |
|-----------|---------|--------|--------|
| Easy | 10 | 20 | candidates.json (22 remaining) |
| Medium | 2 (unvalidated) | 8 | candidates.json |
| Hard | 0 | 2+ | CrossNF candidates |

### 4.2 Expand NF Coverage

| NF | Current | Target |
|----|---------|--------|
| AMF | 3 | 5 |
| PCF | 2 | 3 |
| SMF | 1 | 3 |
| UDM | 1 | 3 |
| UPF | 0 | 2 |
| Others | 3 | 4 |

### 4.3 Add CrossNF Tasks

- AMF-SMF registration failure
- SMF-UPF PFCP signaling bug
- Requires docker-compose for multi-container setup

---

## Phase 5: Paper Finalization (Priority: MEDIUM)

### 5.1 Content Checklist

- [ ] Abstract: finalized with quantitative results
- [ ] Section 6 (Spec-as-Skill): fill A/B experiment data
- [ ] Section 7 (Baseline): fill agent comparison table
- [ ] Section 7.3 (Failure Analysis): add qualitative examples
- [ ] Figures: resolve rate bar charts, A/B comparison, token overhead
- [ ] Appendix: full instance table, spec excerpt examples

### 5.2 Target Venue

- Primary: **NeurIPS 2025 Datasets & Benchmarks Track**
- Backup: EMNLP 2025 Industry Track, ICSE 2026 NIER

### 5.3 Submission Timeline

| Milestone | Target Date | Status |
|-----------|-------------|--------|
| 20+ validated instances | April 2026 | In progress |
| Agent baselines complete | April 2026 | Not started |
| A/B experiment complete | May 2026 | Framework ready |
| Paper draft v2 | May 2026 | v1 done |
| Internal review | May 2026 | |
| Submission | June 2026 | |

---

## Phase 6: Future Directions (Priority: LOW)

### 6.1 SearchSWE-5G

Build a 3GPP specification search tool via MCP:
- Agent can query: "What does TS 29.514 say about suppFeat?"
- Evaluate: does RAG-based spec access improve resolve rate?
- Comparison: curated excerpt (Phase 3) vs RAG retrieval

### 6.2 Multi-Modal Inputs

- Add 3GPP signaling flow diagrams (from TS 23.502) as images
- Test vision-capable models (GPT-4o, Claude Sonnet) on diagram understanding
- Parallels SWE-Bench Mobile's Figma design inputs

### 6.3 Protocol Compliance Testing

Beyond "does the code not crash":
- Verify NF responses conform to 3GPP message formats
- Check state machine transitions follow specification
- Unique to telecom domain

### 6.4 Beyond Bug Fixing

- Feature implementation tasks (e.g., add 3GPP R16 feature)
- Protocol migration tasks (e.g., R15 → R16 upgrade)
- Configuration debugging tasks

---

## Immediate Next Actions

1. **Validate remaining instances** on server: pcf_pr57, smf_pr112, udm_pr78
2. **Install Aider in Docker** and test on pilot instance
3. **Run Claude Code** on pilot instance (if API key available)
4. **Upload v0.3 to HuggingFace** after more validations pass
5. **Enable GitHub Pages** (Settings → Pages → /docs)
