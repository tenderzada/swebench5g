---
license: cc-by-4.0
task_categories:
  - text-generation
language:
  - en
tags:
  - code
  - software-engineering
  - 5g
  - telecommunications
  - benchmark
  - agent-evaluation
size_categories:
  - n<1K
pretty_name: SWE-Bench 5G
---

# SWE-Bench 5G

A benchmark for evaluating AI coding agents on **5G core network** software engineering tasks, built on the open-source [free5GC](https://github.com/free5gc/free5gc) project.

## Overview

SWE-Bench 5G adapts the [SWE-Bench](https://www.swebench.com/) methodology to the **telecommunications domain**, specifically targeting the 5G core network (5GC) implementation based on 3GPP specifications. Each task instance is a real bug from the free5GC project, packaged as a Docker image with:

- Pre-installed Go toolchain and all dependencies
- Source code checked out at the **buggy commit**
- **Existing tests** that pass on the buggy version (no regression)
- **Fail-to-pass tests** that fail on the buggy version and pass after the correct fix

This design follows the [BeyondSWE](https://arxiv.org/abs/2603.03194) approach of providing fully reproducible Docker environments.

## What Makes 5G Tasks Unique

| Dimension | General SWE-Bench | SWE-Bench 5G |
|-----------|-------------------|--------------|
| Domain knowledge | General programming | 3GPP protocol specifications |
| Language | Python | Go |
| Architecture | Single application | Distributed microservices (NFs) |
| Complexity | Code logic | Protocol state machines + signaling |
| Specification | Issue description | Issue + 3GPP TS references |

## Dataset Schema

| Field | Type | Description |
|-------|------|-------------|
| `instance_id` | string | Unique task identifier |
| `dataset_id` | string | Always "SWEBench5G" |
| `task` | string | Task type: `SingleNF`, `CrossNF`, `Protocol`, `DataPlane` |
| `nf_type` | string | 5G Network Function: `AMF`, `SMF`, `PCF`, `UPF`, etc. |
| `repo` | string | Source repository |
| `language` | string | Programming language |
| `image_url` | string | Docker image for reproducible environment |
| `patch` | string | Ground-truth fix (unified diff) |
| `commit_id` | string | Fix commit hash |
| `parent_commit` | string | Buggy commit hash (base version) |
| `problem_statement` | string | Bug description with 3GPP spec references |
| `f2p_patch` | string | Fail-to-pass test code |
| `f2p_script` | string | Command to run fail-to-pass tests |
| `FAIL_TO_PASS` | string | JSON list of test names that should fail→pass |
| `PASS_TO_PASS` | string | JSON list of test names that should always pass |
| `github` | string | Link to original GitHub issue |
| `pre_commands` | string | Setup commands before testing |
| `spec_reference` | string | Related 3GPP specification sections |
| `difficulty` | string | `easy`, `medium`, `hard` |
| `affected_function` | string | Function containing the bug |
| `lines_changed` | int | Number of lines in the fix |
| `files_changed` | int | Number of files modified |
| `test_suite_num` | string | Total number of test cases |
| `split` | string | Data split (`test`) |

## Task Categories

- **SingleNF**: Bug within a single Network Function (e.g., PCF nil pointer)
- **CrossNF**: Bug requiring changes across multiple NFs (e.g., AMF-SMF signaling)
- **Protocol**: Protocol compliance issues (e.g., NAS/NGAP message handling)
- **DataPlane**: User plane (UPF) forwarding or PFCP issues

## 5G Network Functions Covered

| NF | Full Name | Role |
|----|-----------|------|
| AMF | Access and Mobility Management | UE registration, mobility |
| SMF | Session Management | PDU session lifecycle |
| PCF | Policy Control | QoS and charging policies |
| UPF | User Plane Function | Data forwarding |
| NRF | NF Repository | Service discovery |
| UDM | Unified Data Management | Subscriber data |
| AUSF | Authentication Server | Authentication |
| NSSF | Network Slice Selection | Slice management |

## Usage

### Load the dataset

```python
from datasets import load_dataset

ds = load_dataset("tenderzada/SWEBench5G", split="test")
print(ds[0]["problem_statement"])
```

### Run a task instance

```bash
# Pull the pre-built Docker image
docker pull swebench5g/free5gc:pcf_issue_879

# Start the environment
docker run -it swebench5g/free5gc:pcf_issue_879

# Inside the container:
cat /opt/task/problem_statement.md           # Read the task
/opt/test-suite/run_tests.sh existing        # Verify environment (should PASS)
/opt/test-suite/run_tests.sh fail            # Confirm bug exists (should FAIL)
# ... agent fixes the code ...
/opt/test-suite/run_tests.sh all             # Verify fix (should ALL PASS)
```

## Evaluation

A task is considered **resolved** if:
1. All `FAIL_TO_PASS` tests pass after the agent's patch
2. All `PASS_TO_PASS` tests still pass (no regression)

## Current Status

This benchmark is under active development. We are expanding the dataset by mining issues from the full free5GC ecosystem (20+ sub-repositories).

| Version | Instances | NFs Covered | Difficulty Distribution |
|---------|-----------|-------------|------------------------|
| v0.1 | 1 | PCF | 1 easy |
| v0.2 | 10 | AMF, PCF, SMF, UDM, NRF, NSSF, AUSF | 10 easy |
| v1.0 (planned) | 50+ | All NFs | easy/medium/hard |

## Related Work

- [SWE-Bench](https://www.swebench.com/) - The original Python software engineering benchmark
- [SWE-Bench Mobile](https://arxiv.org/abs/2602.09540) - iOS/Swift variant with multimodal inputs
- [BeyondSWE](https://arxiv.org/abs/2603.03194) - Extended scope with Docker environments
- [free5GC](https://github.com/free5gc/free5gc) - Open-source 5G core network

## Citation

```bibtex
@misc{swebench5g2026,
  title={SWE-Bench 5G: Evaluating AI Coding Agents on 5G Core Network Engineering Tasks},
  author={tenderzada},
  year={2026},
  url={https://huggingface.co/datasets/tenderzada/SWEBench5G}
}
```

## License

This dataset is released under [CC-BY-4.0](https://creativecommons.org/licenses/by/4.0/).
The underlying free5GC source code is under [Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0).
