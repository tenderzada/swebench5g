# SWE-Bench 5G 中文教程

## 一、项目简介

### 1.1 什么是 SWE-Bench 5G？

SWE-Bench 5G 是一个用于**评估 AI 编程 Agent 在 5G 核心网软件工程任务上表现**的基准测试（Benchmark）。

简单来说：我们从开源 5G 核心网项目 [free5GC](https://github.com/free5gc/free5gc) 中收集真实的 bug，把每个 bug 打包成一个 Docker 镜像，然后让 AI Agent（如 Claude Code、Cursor、Codex）去修复。通过自动化测试来判断 Agent 是否成功。

### 1.2 为什么要做这件事？

现有的编程 Agent 评测基准（SWE-Bench、SWE-Bench Mobile、BeyondSWE）主要针对 Python、Swift 等通用软件。**没有任何基准覆盖电信领域**。

5G 核心网软件有三个独特挑战：

1. **规范驱动**：代码正确性由 3GPP 技术规范定义，不仅仅是"不崩溃"
2. **分布式架构**：多个网络功能（NF）协同工作，bug 可能跨服务传播
3. **协议状态机**：注册、会话建立等流程涉及复杂的多步状态机

### 1.3 项目灵感来源

| 项目 | 特点 | 我们借鉴了什么 |
|------|------|---------------|
| [SWE-Bench](https://arxiv.org/abs/2310.06770) | Python bug 修复，2294 个实例 | 任务格式：issue → patch → fail-to-pass 测试 |
| [SWE-Bench Mobile](https://arxiv.org/abs/2602.09540) | iOS 开发，多模态输入（PRD+Figma） | 多模态思路：3GPP 规范 + 信令流程图 |
| [BeyondSWE](https://arxiv.org/abs/2603.03194) | 500 实例，Docker 镜像打包 | Docker 环境：每个 bug 一个完整的可复现环境 |

---

## 二、核心概念

### 2.1 Task Instance（任务实例）

每个任务实例代表一个真实的 bug，包含：

```
T = (I, O, E)

I（输入）：问题描述 + 3GPP 规范引用 + 源代码
O（输出）：Agent 生成的代码补丁（unified diff）
E（评估）：自动化测试套件
```

### 2.2 两类测试

| 类型 | 名称 | 在 buggy 版本上 | 在修复后 | 作用 |
|------|------|----------------|---------|------|
| PASS_TO_PASS | 已有测试 | PASS | PASS | 确保修复不引入回归 |
| FAIL_TO_PASS | 新增测试 | FAIL | PASS | 验证 bug 确实被修复 |

### 2.3 判定标准

一个任务被判定为 **Resolved（已解决）**，当且仅当：
- 所有 FAIL_TO_PASS 测试通过（bug 被修复）
- 所有 PASS_TO_PASS 测试仍然通过（没有回归）

### 2.4 Docker 镜像结构

每个任务被打包成一个 Docker 镜像：

```
镜像内部
├── /opt/free5gc-pcf/              ← NF 源码（停在有 bug 的版本）
│   ├── internal/sbi/processor/    ← bug 所在位置
│   ├── go.mod, go.sum             ← 依赖（已预下载）
│   └── .git/                      ← 完整 git 历史
├── /opt/test-suite/
│   ├── existing_test.go           ← PASS_TO_PASS 测试
│   ├── fail_test.go               ← FAIL_TO_PASS 测试
│   └── run_tests.sh               ← 一键测试脚本
└── /opt/task/
    └── problem_statement.md       ← Agent 看到的任务描述
```

---

## 三、快速开始

### 3.1 环境要求

- Docker 20.10+
- Python 3.8+（用于挖掘脚本和评估 harness）
- Git

### 3.2 获取项目

```bash
git clone git@github.com:tenderzada/swebench5g.git
cd swebench5g
```

### 3.3 体验 Pilot 任务

我们已经有一个完整验证过的任务：**PCF Issue #879**（策略控制函数空指针崩溃）。

```bash
# 进入 pilot 目录
cd pilot_pcf_879

# 构建镜像（首次需要在宿主机 clone 源码）
bash scripts/build_image.sh

# 一键验证
bash scripts/validate_image.sh
```

验证过程会依次执行三步：

```
Step 1: 已有测试 → PASS（环境正常）           ✅
Step 2: 触发 bug 的测试 → FAIL（bug 存在）     ✅
Step 3: 应用修复后 → ALL PASS（修复有效）       ✅
```

### 3.4 手动体验 Agent 视角

模拟 AI Agent 的工作流程：

```bash
# 启动容器
docker run -it swebench5g/free5gc:pcf_issue_879

# === 以下在容器内操作 ===

# 1. 阅读任务描述
cat /opt/task/problem_statement.md

# 2. 确认环境正常
/opt/test-suite/run_tests.sh existing
# → 2/2 PASS

# 3. 确认 bug 存在
/opt/test-suite/run_tests.sh fail
# → 3/3 FAIL（panic: nil pointer dereference）

# 4. 查看有 bug 的代码
# 提示：bug 在 internal/sbi/processor/policyauthorization.go
# 函数 provisioningOfTrafficRoutingInfo 没有检查 routeReq 是否为 nil

# 5. 修复 bug（加一个 nil check）
# ... 编辑代码 ...

# 6. 验证修复
/opt/test-suite/run_tests.sh all
# → 5/5 PASS（修复成功！）

# 7. 退出
exit
```

---

## 四、扩充数据集

### 4.1 挖掘候选 Issues

使用自动化脚本扫描 free5GC 的所有子仓库：

```bash
# 设置 GitHub Token（可选但强烈推荐，否则限流 60 次/小时）
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx

# 扫描所有仓库
python scripts/mine_issues.py --output candidates.json

# 或只扫描核心 NF
python scripts/mine_issues.py --repo amf smf pcf upf --output candidates.json
```

脚本会自动：
- 获取所有 merged PR
- 过滤出 bug 修复类 PR（包含 fix/crash/panic 等关键词）
- 排除纯重构、文档、依赖升级
- 提取 base_commit（buggy）和 fix_commit
- 分类难度（easy/medium/hard）和 bug 类型
- 输出结构化的 `candidates.json`

### 4.2 候选筛选标准

从 `candidates.json` 中手动挑选优质候选，标准：

| 条件 | 说明 |
|------|------|
| 有清晰的 issue 描述 | Agent 需要理解问题 |
| 有明确的复现步骤 | 便于编写 fail-to-pass 测试 |
| 修改了源码文件 | 排除纯配置/文档变更 |
| 单 NF 内可测试 | 避免需要启动整个 5GC |
| 有 3GPP 规范关联 | 体现电信领域特色 |

### 4.3 创建新的任务实例

以 AMF Issue #713 为例：

```bash
# 1. 创建实例目录
mkdir -p instances/amf_pr179

# 2. 编写配置文件
# instances/amf_pr179/instance.json
{
  "instance_id": "amf_pr179",
  "nf_type": "AMF",
  "repo": "free5gc/amf",
  "parent_commit": "xxxx",        ← bug 版本的 commit
  "fix_commit": "yyyy",            ← 修复后的 commit
  "go_version": "1.21",
  "test_package_dir": "internal/nas",
  "test_package_path": "internal/nas",
  "difficulty": "easy",
  ...
}

# 3. 编写任务描述
# instances/amf_pr179/problem_statement.md

# 4. 编写测试
# instances/amf_pr179/existing_test.go  ← 验证已有功能正常
# instances/amf_pr179/fail_test.go      ← 触发 bug 的测试

# 5. 构建并验证
bash scripts/build_instance.sh instances/amf_pr179/instance.json
bash scripts/validate_instance.sh instances/amf_pr179/instance.json
```

### 4.4 测试编写指南

**fail_test.go（最关键）：**

```go
package 目标包名

import "testing"

func TestXxx_触发Bug的场景(t *testing.T) {
    // 构造最小化的输入，能触发 bug
    // 使用 defer recover 捕获 panic
    defer func() {
        if r := recover(); r != nil {
            t.Fatalf("BUG PRESENT: %v", r)
        }
    }()

    // 调用有 bug 的函数
    result := buggyFunction(恶意输入)
    _ = result
}
```

**existing_test.go：**

```go
func TestXxx_正常场景(t *testing.T) {
    // 验证正常输入下函数行为正确
    // 这些测试在 buggy 版本上也必须 PASS
    result := buggyFunction(正常输入)
    if result != expected {
        t.Fatalf("expected %v, got %v", expected, result)
    }
}
```

**关键原则：**
- 测试必须在目标函数的**同一个 Go 包**内（才能访问未导出的函数）
- fail_test.go 在 buggy 版本上**必须 FAIL**
- existing_test.go 在 buggy 版本上**必须 PASS**
- 修复后**两组都必须 PASS**

---

## 五、评估 AI Agent

### 5.1 评估框架架构

```
┌─────────────────┐
│  run_evaluation  │  评估入口
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌──────────────────┐
│  Docker 容器     │ ←── │  任务镜像         │
│  (隔离环境)      │     │  (源码+测试+描述)  │
└────────┬────────┘     └──────────────────┘
         │
         ▼
┌─────────────────┐
│  Agent 适配器    │  Claude Code / Aider / Codex
│  (读描述→改代码) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  提取 patch      │  git diff
│  运行测试        │  run_tests.sh all
│  记录结果        │  resolved / not
└─────────────────┘
```

### 5.2 运行评估

```bash
# 评估单个 Agent
python eval/run_evaluation.py \
  --agent claude-code \
  --model opus-4.6 \
  --dataset dataset/swebench5g.jsonl

# 评估所有 Agent
python eval/run_evaluation.py --agent all

# 多次运行（统计方差）
python eval/run_evaluation.py \
  --agent claude-code \
  --model opus-4.6 \
  --runs 3

# 指定特定实例
python eval/run_evaluation.py \
  --agent claude-code \
  --instances free5gc_pcf_issue879
```

### 5.3 结果输出

评估完成后在 `eval/results/<timestamp>/` 生成：

**results.json**（每条评测的详细数据）：
```json
{
  "instance_id": "free5gc_pcf_issue879",
  "agent": "claude-code",
  "model": "opus-4.6",
  "resolved": true,
  "existing_tests_pass": true,
  "fail_tests_pass": true,
  "patch": "diff --git a/...",
  "time_seconds": 45.2,
  "timestamp": "2026-03-23T..."
}
```

**summary.txt**（汇总报告）：
```
Total: 10 | Resolved: 4 (40.0%)
  claude-code+opus-4.6:  3/5 (60.0%)
  aider+opus-4.6:        1/5 (20.0%)
```

### 5.4 核心指标

| 指标 | 定义 |
|------|------|
| **Resolve Rate** | 完全解决的实例比例（FAIL_TO_PASS 全过 + PASS_TO_PASS 无回归） |
| **Test Pass Rate** | 单个测试用例的平均通过率 |

分析维度：
- 按 NF 类型（PCF 比 AMF 容易吗？）
- 按难度（easy / medium / hard）
- 按 bug 类型（空指针 vs 逻辑错误 vs 协议违规）
- 按是否涉及 3GPP 规范知识

---

## 六、项目结构

```
swebench5g/
│
├── pilot_pcf_879/              ← 已验证的 pilot 任务
│   ├── Dockerfile
│   ├── problem_statement.md
│   ├── task_metadata.json
│   ├── test-suite/
│   │   ├── existing_test.go
│   │   ├── fail_test.go
│   │   └── run_tests.sh
│   └── scripts/
│
├── instances/                  ← 新任务实例（模板化构建）
│   └── amf_pr179/
│       ├── instance.json
│       └── problem_statement.md
│
├── templates/                  ← Docker/测试模板
│   ├── Dockerfile.template
│   └── run_tests.sh.template
│
├── scripts/                    ← 工具脚本
│   ├── mine_issues.py          ← Issue 挖掘
│   ├── build_instance.sh       ← 构建镜像
│   └── validate_instance.sh    ← 验证镜像
│
├── eval/                       ← 评估框架
│   └── run_evaluation.py
│
├── dataset/                    ← HuggingFace 数据集
│   ├── swebench5g.jsonl
│   ├── README.md
│   └── upload_to_hf.py
│
├── paper/                      ← 论文
│   ├── main.tex
│   └── references.bib
│
└── ROADMAP.md                  ← 项目路线图
```

---

## 七、5G 背景知识（给不熟悉电信的读者）

### 7.1 什么是 5G 核心网？

5G 核心网（5GC）是 5G 网络的"大脑"，负责用户认证、会话管理、策略控制等。它由多个**网络功能（NF）**组成，每个 NF 是一个独立的微服务：

```
手机 (UE)
  │
  ▼
基站 (gNB)
  │
  ▼
┌─────────────────────────────────────────┐
│  5G 核心网                               │
│                                          │
│  AMF ──── SMF ──── UPF ──→ 互联网       │
│   │        │                             │
│  AUSF    PCF                             │
│   │                                      │
│  UDM ── UDR     NRF    NSSF             │
└─────────────────────────────────────────┘
```

### 7.2 各 NF 的职责

| NF | 全称 | 职责 | 类比 |
|----|------|------|------|
| AMF | 接入和移动管理 | 用户注册、移动性 | 前台接待 |
| SMF | 会话管理 | 建立数据通道 | 交换机 |
| UPF | 用户面功能 | 转发用户数据 | 路由器 |
| PCF | 策略控制 | QoS、计费策略 | 规则引擎 |
| UDM | 用户数据管理 | 订阅信息 | 用户数据库 |
| AUSF | 认证服务 | 身份验证 | 门禁系统 |
| NRF | NF 注册发现 | 服务发现 | DNS |
| NSSF | 切片选择 | 网络切片 | VLAN 管理器 |

### 7.3 什么是 3GPP 规范？

3GPP（第三代合作伙伴计划）是制定移动通信标准的国际组织。5G 核心网的行为由一系列**技术规范（TS）**定义：

| 规范 | 内容 |
|------|------|
| TS 23.501 | 5G 系统架构 |
| TS 23.502 | 5G 系统流程（注册、会话等） |
| TS 29.5xx | 各 NF 的服务接口定义 |
| TS 24.501 | 非接入层（NAS）协议 |

在我们的 benchmark 中，部分 bug 的修复需要理解 3GPP 规范。例如 pilot 任务中，`suppFeat` 字段的含义定义在 TS 29.514 中。

### 7.4 free5GC 简介

[free5GC](https://github.com/free5gc/free5gc) 是用 **Go 语言**编写的开源 5G 核心网实现：

- 基于 3GPP Release 15+ 规范
- Apache 2.0 开源许可
- 2.3K GitHub Stars
- 实现了完整的控制面和用户面 NF
- 每个 NF 是独立的 Git 仓库（方便我们逐个挖掘 bug）

---

## 八、常见问题

### Q: Docker build 时 GitHub/Go 模块下载超时？

**原因**：国内网络访问 GitHub 和 Go 模块代理不稳定。

**解决**：
```bash
# Dockerfile 中已设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 源码在宿主机 clone（利用宿主机代理/镜像）
# build_instance.sh 会自动处理
```

### Q: 测试为什么要放在同一个 Go 包内？

Go 的未导出函数（小写字母开头）只能在同一个包内访问。由于 free5GC 的很多内部函数是未导出的，测试必须声明为 `package processor`（而不是 `package processor_test`）才能直接调用。

### Q: 如何判断一个 issue 适不适合做 task instance？

理想的候选：
- 有清晰的 crash/panic 日志
- 能用简单的输入触发
- 不需要启动整个 5GC（单元/函数级可测试）
- PR 修改在 1-3 个文件内

不适合的：
- 需要真实 UE/gNB 设备
- 纯性能问题（难以用测试判定）
- 涉及数据库 schema 迁移

### Q: 数据集在哪里？

HuggingFace: [tenderzada/SWEBench5G](https://huggingface.co/datasets/tenderzada/SWEBench5G)

```python
from datasets import load_dataset
ds = load_dataset("tenderzada/SWEBench5G", split="test")
print(ds[0]["problem_statement"])
```

---

## 九、参与贡献

欢迎贡献新的 task instance！流程：

1. Fork 本仓库
2. 从 `candidates.json` 中选择一个候选
3. 在 `instances/<id>/` 下创建测试和描述
4. 用 `build_instance.sh` 构建并用 `validate_instance.sh` 验证
5. 提交 PR

---

## 十、引用

```bibtex
@misc{swebench5g2026,
  title={SWE-Bench 5G: Evaluating AI Coding Agents on 5G Core Network Tasks},
  author={tenderzada},
  year={2026},
  url={https://huggingface.co/datasets/tenderzada/SWEBench5G}
}
```
