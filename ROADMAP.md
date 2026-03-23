# SWE-Bench 5G Roadmap

## Current Status (v0.1)

- 1 validated task instance (PCF Issue #879)
- Docker image built and tested
- Dataset published on HuggingFace
- Paper draft (NeurIPS D&B format)
- No agent evaluation results yet
- No evaluation harness

---

## Phase 1: Scale Dataset to 30+ Instances (Priority: HIGH)

### Why
- SWE-Bench Mobile has 50 tasks, BeyondSWE has 500
- 1 instance is a proof-of-concept, not a benchmark
- Reviewers will expect at least 30-50 instances for a D&B paper

### How

#### 1.1 Mine Candidate Issues from All free5GC Sub-repos

Target repositories (each NF is a separate repo):

| Repo | NF | Priority |
|------|----|----------|
| free5gc/amf | AMF | High |
| free5gc/smf | SMF | High |
| free5gc/pcf | PCF | Done (1) |
| free5gc/upf | UPF | High |
| free5gc/udm | UDM | Medium |
| free5gc/udr | UDR | Medium |
| free5gc/ausf | AUSF | Medium |
| free5gc/nrf | NRF | Medium |
| free5gc/nssf | NSSF | Low |
| free5gc/n3iwf | N3IWF | Low |
| free5gc/nas | NAS lib | High |
| free5gc/ngap | NGAP lib | High |
| free5gc/pfcp | PFCP lib | High |
| free5gc/openapi | OpenAPI models | Medium |

#### 1.2 Automate Candidate Screening

Build a script (`scripts/mine_issues.py`) that:
1. Uses GitHub API to fetch closed issues with linked merged PRs
2. Filters: has code changes, not pure refactor/docs/deps
3. Identifies base_commit and fix_commit
4. Extracts affected files and functions
5. Outputs candidate list with metadata

Reference: BeyondSWE automated their pipeline across 246 repos.

#### 1.3 Semi-Automate Test Generation

The biggest bottleneck: free5GC has almost zero existing tests.
Strategy:
- Use LLM to draft fail-to-pass tests from issue description + diff
- Human review and refine
- Validate with the 3-step check (existing PASS, fail FAIL, fix ALL PASS)

#### 1.4 Templatize Docker Image Building

Create a reusable Dockerfile template + build script:
- Input: repo name, base_commit, test files
- Output: validated Docker image
- Reference: BeyondSWE's `pre_commands` pattern for environment setup

### Deliverable
- 30+ validated instances across AMF/SMF/PCF/UPF/NAS/NGAP
- Automated mining + semi-automated test generation pipeline
- All images built and pushed to Docker Hub

---

## Phase 2: Build Evaluation Harness (Priority: HIGH)

### Why
- BeyondSWE has a complete evaluation framework
- SWE-Bench Mobile has patch-to-task routing + automated scoring
- Without a harness, nobody can reproduce our results

### How

#### 2.1 Evaluation Runner (`eval/run_evaluation.py`)

Design following BeyondSWE's approach:
```
For each task instance:
  1. Pull Docker image
  2. Start container
  3. Inject agent (via docker exec or API)
  4. Agent reads problem_statement, writes patch
  5. Extract patch (git diff)
  6. Run test suite
  7. Record: resolved/not, time, tokens used
```

#### 2.2 Agent Adapters

Build adapters for each agent to standardize the interface:
- Claude Code: run inside container via CLI
- Cursor: use headless mode or API
- Codex CLI: run inside container
- OpenCode: run inside container
- Aider: run inside container

Reference: SWE-Bench Mobile tested 4 agents x 11 models = 22 configurations.

#### 2.3 Metrics

Following SWE-Bench conventions:
- **Resolve Rate**: % of instances where all FAIL_TO_PASS pass and all PASS_TO_PASS still pass
- **Test Pass Rate**: average % of individual tests passed across instances
- Breakdown by: NF type, difficulty, bug type, lines changed

#### 2.4 Result Dashboard

Generate a summary table + per-instance breakdown, similar to:
- BeyondSWE's per-task-type results
- SWE-Bench Mobile's agent x model matrix

### Deliverable
- `eval/` directory with runner, agent adapters, and scoring
- Reproducible evaluation on 3+ agent-model combinations
- Results table for the paper

---

## Phase 3: Run Baseline Experiments (Priority: HIGH)

### Why
- Paper Table 5 is currently TBD
- Reviewers need quantitative results to assess the benchmark's value

### How

#### 3.1 Agent-Model Configurations to Test

| Agent | Models | Total Configs |
|-------|--------|---------------|
| Claude Code | Opus 4.6, Sonnet 4.6, Haiku 4.5 | 3 |
| Cursor | Opus 4.6, Sonnet 4.6, GPT-5 | 3 |
| Codex CLI | GPT-4.1, GPT-5 | 2 |
| Aider | Opus 4.6, Sonnet 4.6 | 2 |
| **Total** | | **10** |

#### 3.2 Experiment Protocol

For each (agent, model, task):
- 3 independent runs (for variance estimation)
- Record: resolved (bool), time (s), tokens consumed, patch size
- Timeout: 30 min per instance

Reference: SWE-Bench Mobile ran 22 configurations; BeyondSWE reported per-task-type results.

#### 3.3 Analysis Dimensions

1. **Overall resolve rate** by agent and model
2. **By NF type**: Are PCF bugs easier than AMF bugs?
3. **By difficulty**: Easy vs Medium vs Hard
4. **By bug type**: nil pointer vs logic error vs protocol non-compliance
5. **Spec-dependence analysis**: Do tasks referencing 3GPP specs have lower resolve rates?
6. **Agent behavior analysis**: Do agents read the 3GPP spec references?

### Deliverable
- Complete results for paper Tables 5+
- Analysis figures for paper Section 6

---

## Phase 4: Add 5G-Specific Innovations (Priority: MEDIUM)

### Why
- Need differentiation from simply "SWE-Bench but in Go"
- SWE-Bench Mobile's key innovation was multi-modal; we need ours

### How

#### 4.1 3GPP Specification Augmentation (Multi-Modal)

Inspired by SWE-Bench Mobile's PRD + Figma:
- Extract relevant 3GPP TS clauses for each task
- Include protocol signaling flow diagrams (from TS 23.502) as images
- Test whether providing spec context improves resolve rate

This creates a unique "specification-grounded coding" evaluation axis.

#### 4.2 SearchSWE-5G: Spec-Aware Search Framework

Inspired by BeyondSWE's SearchSWE:
- Augment agent with a 3GPP specification search tool
- Agent can query: "What does TS 29.514 say about suppFeat bit 1?"
- Measure: does spec search improve resolve rate?

Implementation: Index 3GPP specs as a RAG knowledge base, expose via MCP tool.

#### 4.3 CrossNF Task Instances

BeyondSWE's CrossRepo is our CrossNF:
- Bugs that span AMF + SMF, or SMF + UPF
- Require multi-container Docker setup (docker-compose)
- Agent must trace signaling flow across NF boundaries

#### 4.4 Protocol Compliance Testing

Beyond "does the code not crash":
- Test that NF responses conform to 3GPP message formats
- Verify state machine transitions follow specification
- This is unique to telecom—no existing benchmark tests spec compliance

### Deliverable
- 3GPP spec RAG tool (MCP server)
- 5-10 CrossNF task instances
- Protocol compliance test layer

---

## Phase 5: Paper Revision (Priority: MEDIUM)

### What to Add to the Paper

#### Based on SWE-Bench Mobile's Structure
- [ ] Detailed failure analysis (what types of errors do agents make?)
- [ ] Prompt engineering experiments (defensive programming vs complex prompts)
- [ ] Agent x Model performance matrix (heatmap figure)
- [ ] Difficulty correlation analysis (resolve rate vs files/lines changed)
- [ ] Qualitative examples of agent successes and failures

#### Based on BeyondSWE's Structure
- [ ] Per-task-type breakdown (SingleNF vs CrossNF vs Protocol vs DataPlane)
- [ ] Comparison with SWE-Bench performance on similar bugs in Python
- [ ] SearchSWE-5G ablation: with vs without spec search
- [ ] Scaling analysis: how does performance change as dataset grows?

#### 5G-Specific Contributions to Highlight
- [ ] Specification-grounded evaluation (unique to telecom)
- [ ] Distributed NF coordination challenges
- [ ] Go language coverage (fills gap in existing benchmarks)
- [ ] Protocol state machine reasoning

### Deliverable
- Revised paper with full results and analysis
- Submission-ready for NeurIPS 2025 D&B track

---

## Timeline

| Phase | Duration | Deadline | Dependency |
|-------|----------|----------|------------|
| Phase 1: Scale to 30+ instances | 4-6 weeks | May 2026 | None |
| Phase 2: Evaluation harness | 2-3 weeks | May 2026 | Can parallel with Phase 1 |
| Phase 3: Baseline experiments | 2-3 weeks | June 2026 | Phase 1 + 2 |
| Phase 4: 5G innovations | 3-4 weeks | June 2026 | Phase 1 |
| Phase 5: Paper revision | 2-3 weeks | July 2026 | Phase 3 + 4 |

**NeurIPS 2025 D&B submission deadline: typically late May/early June.**
Check exact date and adjust timeline accordingly.

---

## Immediate Next Actions (This Week)

1. [ ] Build `scripts/mine_issues.py` to auto-scan all free5GC sub-repos
2. [ ] Create Issue #713 (AMF panic) and #794 (SMF crash) as next two instances
3. [ ] Design Dockerfile template for rapid instance creation
4. [ ] Set up eval harness skeleton
5. [ ] Run Claude Code on the pilot task to get first baseline number
