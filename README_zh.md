# SWE-Bench 5G

**评测 AI 编程 Agent 在电信网络工程任务上的表现**

[数据集](https://huggingface.co/datasets/tenderzada/SWEBench5G) | [项目主页](https://tenderzada.github.io/swebench5g/) | [论文](paper/main.tex) | [中文教程](tutorial_zh.md)

---

## 概述

SWE-Bench 5G 是首个面向**电信领域**的 AI 编程 Agent 评测基准。基于开源 5G 核心网 [free5GC](https://github.com/free5gc/free5gc)，将真实 bug 打包为 Docker 镜像，配合自动化测试评估 Agent 的修复能力。

**核心研究问题**（受 [SWE-Skills-Bench](https://arxiv.org/abs/2603.15401) 启发）：

> *将 3GPP 协议规范作为"技能文档"注入 AI Agent，能否提升其修复 5G bug 的能力？*

## 核心特性

- **10 个已验证实例**，覆盖 7 个网络功能（AMF、PCF、SMF、UDM、NRF、NSSF、AUSF）
- **280 个候选池**，来自 16 个 free5GC 子仓库
- **双测试策略**：直接函数调用 + diff-based intent test
- **Specification-as-Skill A/B 框架**：对比有/无 3GPP 规范注入的效果
- **Docker 可复现环境**，参照 [BeyondSWE](https://arxiv.org/abs/2603.03194) 方法论
- **全自动化流水线**：挖掘 → 构建 → 验证 → 评测

## 快速开始

```bash
# 克隆项目
git clone git@github.com:tenderzada/swebench5g.git && cd swebench5g

# 构建 pilot 实例
cd pilot_pcf_879 && bash scripts/build_image.sh

# 验证（三步检查）
bash scripts/validate_image.sh
# Step 1: 已有测试     → PASS（环境正常）
# Step 2: 触发 bug     → FAIL（bug 存在）
# Step 3: 修复后       → ALL PASS

# 体验 Agent 视角
docker run -it swebench5g/free5gc:pcf_issue_879
cat /opt/task/problem_statement.md       # 阅读任务描述
/opt/test-suite/run_tests.sh fail        # 确认 bug
# ... 修复代码 ...
/opt/test-suite/run_tests.sh all         # 验证修复
```

## 评测 Agent

```bash
# Qwen3.5-Flash（单轮 API 模式）
export DASHSCOPE_API_KEY=sk-xxx
python eval/run_evaluation.py --agent qwen --model qwen3.5-flash

# A/B 测试：有 vs 无 3GPP 规范
python eval/run_evaluation.py --agent qwen --model qwen3.5-flash --ab-test

# 批量构建验证所有实例
bash scripts/batch_build_validate.sh --clean
```

## 初步评测结果

| 模型 | 模式 | Resolve Rate | 能否定位 Bug |
|------|------|-------------|-------------|
| Qwen3.5-Flash | 单轮 API | 0% (0/1) | 能（每次都正确） |
| Claude Code | Agent CLI | 待测 | 待测 |
| Cursor | Agent IDE | 待测 | 待测 |

**关键发现**：Qwen3.5-Flash 每次都能正确定位 bug 并提出合理修复方案，但无法在单轮模式下生成可正确应用的补丁。**Agent 能力（而非领域知识）是瓶颈。**

## 数据集

```python
from datasets import load_dataset
ds = load_dataset("tenderzada/SWEBench5G", split="test")
print(ds[0]["problem_statement"])
```

## 项目结构

```
swebench5g/
├── pilot_pcf_879/          # 已验证的 pilot 实例
├── instances/              # 21 个任务实例（10 个已验证）
├── specs/                  # 3GPP 规范摘录（10 个文件）
├── templates/              # Dockerfile 和测试模板
├── scripts/                # 挖掘、构建、验证、批量工具
├── eval/                   # 评估框架（支持 Qwen、Claude、Codex）
├── dataset/                # HuggingFace 数据集文件
├── docs/                   # GitHub Pages 展示页面
└── paper/                  # NeurIPS D&B 论文草稿
```

## 灵感来源

| 项目 | 我们借鉴了什么 |
|------|---------------|
| [SWE-Bench](https://arxiv.org/abs/2310.06770) | 任务格式：issue → patch → fail-to-pass 测试 |
| [SWE-Bench Mobile](https://arxiv.org/abs/2602.09540) | Diff-based intent testing |
| [BeyondSWE](https://arxiv.org/abs/2603.03194) | Docker 镜像打包每个 bug |
| [SWE-Skills-Bench](https://arxiv.org/abs/2603.15401) | Specification-as-Skill A/B 实验框架 |

## 引用

```bibtex
@misc{swebench5g2026,
  title={SWE-Bench 5G: Evaluating AI Coding Agents on 5G Core Network Tasks},
  author={tenderzada},
  year={2026},
  url={https://huggingface.co/datasets/tenderzada/SWEBench5G}
}
```

## 许可证

数据集：[CC-BY-4.0](https://creativecommons.org/licenses/by/4.0/) | 代码：MIT | free5GC 源码：[Apache-2.0](https://www.apache.org/licenses/LICENSE-2.0)

## 详细教程

请参阅 [中文教程 (tutorial_zh.md)](tutorial_zh.md)，包含完整的 11 章节详解。
