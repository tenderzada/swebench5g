# SWE-Bench 5G

**Benchmarking AI Coding Agents on Telecom Network Engineering Tasks**

[Dataset](https://huggingface.co/datasets/tenderzada/SWEBench5G) | [Project Page](https://tenderzada.github.io/swebench5g/) | [Paper](paper/main.tex) | [Chinese Tutorial](tutorial_zh.md)

---

## Overview

SWE-Bench 5G is the first benchmark for evaluating AI coding agents in the **telecommunications domain**. Built on [free5GC](https://github.com/free5gc/free5gc), each task packages a real 5G core network bug as a Docker image with automated tests.

**Core Research Question** (inspired by [SWE-Skills-Bench](https://arxiv.org/abs/2603.15401)):

> *Does injecting 3GPP protocol specifications as "skill documents" improve AI agents' ability to fix real 5G bugs?*

## Key Features

- **10 validated task instances** across 7 Network Functions (AMF, PCF, SMF, UDM, NRF, NSSF, AUSF)
- **280 candidate pool** mined from 16 free5GC repositories
- **Dual test strategy**: direct function call + diff-based intent testing
- **Specification-as-Skill A/B framework**: evaluate with/without 3GPP spec injection
- **Docker-based reproducibility** following [BeyondSWE](https://arxiv.org/abs/2603.03194) methodology
- **Automated pipeline**: mine → build → validate → evaluate

## Quick Start

```bash
# Clone
git clone git@github.com:tenderzada/swebench5g.git && cd swebench5g

# Build pilot instance
cd pilot_pcf_879 && bash scripts/build_image.sh

# Validate (3-step check)
bash scripts/validate_image.sh
# Step 1: Existing tests     -> PASS
# Step 2: Fail-to-pass tests -> FAIL (bug confirmed)
# Step 3: After fix           -> ALL PASS

# Experience the agent workflow
docker run -it swebench5g/free5gc:pcf_issue_879
cat /opt/task/problem_statement.md
/opt/test-suite/run_tests.sh fail    # confirm bug
# ... fix the code ...
/opt/test-suite/run_tests.sh all     # verify fix
```

## Evaluate an Agent

```bash
# Qwen3.5-Flash (single-turn API)
export DASHSCOPE_API_KEY=sk-xxx
python eval/run_evaluation.py --agent qwen --model qwen3.5-flash

# A/B test: with vs without 3GPP specification
python eval/run_evaluation.py --agent qwen --model qwen3.5-flash --ab-test

# Batch build & validate all instances
bash scripts/batch_build_validate.sh --clean
```

## Preliminary Results

| Model | Mode | Resolve Rate | Bug Located? |
|-------|------|-------------|-------------|
| Qwen3.5-Flash | Single-turn API | 0% (0/1) | Yes (every attempt) |
| Claude Code | Agentic CLI | TBD | TBD |
| Cursor | Agentic IDE | TBD | TBD |

**Key finding**: Qwen3.5-Flash correctly identifies the bug in every attempt but cannot produce applicable patches in single-turn mode. Agentic capabilities — not domain knowledge alone — appear to be the bottleneck.

## Dataset

```python
from datasets import load_dataset
ds = load_dataset("tenderzada/SWEBench5G", split="test")
print(ds[0]["problem_statement"])
```

## Project Structure

```
swebench5g/
├── pilot_pcf_879/          # Validated pilot instance
├── instances/              # 21 task instances (10 validated)
├── specs/                  # 3GPP specification excerpts (10 files)
├── templates/              # Dockerfile & test templates
├── scripts/                # Mine, build, validate, batch tools
├── eval/                   # Evaluation harness (Qwen, Claude, Codex)
├── dataset/                # HuggingFace dataset files
├── docs/                   # GitHub Pages showcase
└── paper/                  # NeurIPS D&B paper draft
```

## Inspired By

| Project | What We Borrowed |
|---------|-----------------|
| [SWE-Bench](https://arxiv.org/abs/2310.06770) | Task format: issue → patch → fail-to-pass test |
| [SWE-Bench Mobile](https://arxiv.org/abs/2602.09540) | Diff-based intent testing |
| [BeyondSWE](https://arxiv.org/abs/2603.03194) | Docker environment per instance |
| [SWE-Skills-Bench](https://arxiv.org/abs/2603.15401) | Specification-as-Skill A/B framework |

## Citation

```bibtex
@misc{swebench5g2026,
  title={SWE-Bench 5G: Evaluating AI Coding Agents on 5G Core Network Tasks},
  author={tenderzada},
  year={2026},
  url={https://huggingface.co/datasets/tenderzada/SWEBench5G}
}
```

## License

Dataset: [CC-BY-4.0](https://creativecommons.org/licenses/by/4.0/) | Code: MIT | free5GC source: [Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)
