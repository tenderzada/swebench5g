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
| [SWE-Bench Mobile](https://arxiv.org/abs/2602.09540) | iOS 开发，多模态输入（PRD+Figma） | 多模态思路 + **diff-based intent test** |
| [BeyondSWE](https://arxiv.org/abs/2603.03194) | 500 实例，Docker 镜像打包 | Docker 环境：每个 bug 一个完整的可复现环境 |

### 1.4 当前规模

| 指标 | 数值 |
|------|------|
| 总实例数 | **21**（含 1 个 pilot + 20 个新实例） |
| 已验证并发布 | **10** 个（已上传 HuggingFace） |
| 覆盖 NF | 7 个（AMF, PCF, SMF, UDM, NRF, NSSF, AUSF） |
| 难度分布 | 10 easy 已验证 + 2 medium 待验证 |
| 候选池 | 280 个（来自 16 个子仓库） |
| HuggingFace | [tenderzada/SWEBench5G](https://huggingface.co/datasets/tenderzada/SWEBench5G) v0.2 |

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

### 2.3 两种测试策略

在实践中，我们发现 FAIL_TO_PASS 测试需要根据 bug 的可测试性选择不同策略：

#### 策略 A：直接函数调用（首选）

直接调用有 bug 的函数，传入触发 bug 的输入。

```go
// 适用于：函数可以用简单参数直接调用
func TestProvisioningOfTrafficRoutingInfo_NilRouteReq(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Fatalf("BUG: panic: %v", r)
        }
    }()
    // 直接调用有 bug 的函数
    result := provisioningOfTrafficRoutingInfo(smPolicy, "app1", nil, "")
    _ = result
}
```

**优点**：直接验证函数行为。修复函数后，测试自动通过。
**适用**：pilot_pcf_879、amf_pr118、amf_pr157 等。

#### 策略 B：Diff-Based Intent Test（复杂函数的备选）

灵感来自 SWE-Bench Mobile。当 buggy 函数依赖复杂上下文（NGAP 连接、数据库等）无法直接调用时，**检查源代码中是否包含修复模式**。

```go
// 适用于：函数依赖太多上下文，无法在单元测试中调用
func TestHandlerHasGNbIdNilCheck(t *testing.T) {
    data, err := os.ReadFile("handler.go")
    if err != nil {
        t.Fatalf("cannot read handler.go: %v", err)
    }
    src := string(data)
    // 在 buggy 版本中，这个模式不存在 → FAIL
    // 修复后，这个模式存在 → PASS
    if !strings.Contains(src, "GNbId != nil") {
        t.Fatal("BUG: handler.go accesses GNbId.GNBValue without nil check")
    }
}
```

**优点**：不需要复杂的 mock。Agent 必须修改源码才能通过。
**适用**：amf_pr161、nssf_pr39、udm_pr45 等。

#### 如何选择？

```
能直接调用有 bug 的函数？
  ├── 是 → 策略 A（直接函数调用）
  └── 否 → 函数需要复杂上下文？
        ├── 是 → 策略 B（Diff-Based Intent Test）
        └── 否 → 尝试 mock 最小依赖后用策略 A
```

### 2.4 判定标准

一个任务被判定为 **Resolved（已解决）**，当且仅当：
- 所有 FAIL_TO_PASS 测试通过（bug 被修复）
- 所有 PASS_TO_PASS 测试仍然通过（没有回归）

### 2.5 Docker 镜像结构

每个任务被打包成一个 Docker 镜像：

```
镜像内部
├── /opt/free5gc-<nf>/             ← NF 源码（停在有 bug 的版本）
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

- Docker 20.10+（需要已拉取 `golang:1.25-bookworm` 镜像）
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

### 3.4 构建更多实例（模板化方式）

```bash
# 构建 AMF 实例
bash scripts/build_instance.sh instances/amf_pr161/instance.json

# 验证
bash scripts/validate_instance.sh instances/amf_pr161/instance.json
```

### 3.5 手动体验 Agent 视角

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

```bash
# 1. 创建实例目录
mkdir -p instances/amf_pr179

# 2. 编写配置文件 instances/amf_pr179/instance.json
# 3. 编写任务描述 instances/amf_pr179/problem_statement.md
# 4. 编写测试（选择策略 A 或 B）
#    instances/amf_pr179/existing_test.go
#    instances/amf_pr179/fail_test.go

# 5. 构建并验证
bash scripts/build_instance.sh instances/amf_pr179/instance.json
bash scripts/validate_instance.sh instances/amf_pr179/instance.json
```

也可以用批量生成脚本：

```bash
# 在 scripts/generate_instances.py 中添加任务定义，然后运行
python scripts/generate_instances.py
```

### 4.4 测试编写指南

**fail_test.go（最关键）：**

核心原则：**测试必须执行 buggy 代码路径本身，不能自己写安全逻辑。**

```go
// ✅ 正确：复现 buggy 代码的访问模式
if targetRanNodeID.GNbId.GNBValue != "" {  // 会 panic

// ❌ 错误：自己写了安全检查，永远不会 panic
if gnbId != nil && gnbId.GNBValue != "" {  // 不会 panic
```

**策略 A 模板（直接函数调用）：**

```go
package 目标包名

import "testing"

func TestXxx_触发Bug的场景(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Fatalf("BUG PRESENT: %v", r)
        }
    }()
    // 直接调用有 bug 的函数
    result := buggyFunction(恶意输入)
    _ = result
}
```

**策略 B 模板（Diff-Based Intent Test）：**

```go
package 目标包名

import (
    "os"
    "strings"
    "testing"
)

func TestSourceHasFixPattern(t *testing.T) {
    data, err := os.ReadFile("有bug的文件.go")
    if err != nil {
        t.Fatalf("cannot read file: %v", err)
    }
    if !strings.Contains(string(data), "修复后应存在的代码模式") {
        t.Fatal("BUG: source code missing fix pattern")
    }
}
```

**existing_test.go：**

```go
func TestXxx_正常场景(t *testing.T) {
    // 验证正常输入下函数行为正确
    // 这些测试在 buggy 版本上也必须 PASS
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

### 5.5 初步评测结果：Qwen3.5-Flash

我们在 pilot 任务（PCF Issue #879）上使用 **Qwen3.5-Flash**（通过 DashScope API，单轮非 Agent 模式）进行了 7 轮评测。

#### 结果

| 模型 | 模式 | Resolve Rate | 能否定位 Bug | 平均耗时 |
|------|------|-------------|-------------|---------|
| Qwen3.5-Flash | 单轮 API | **0%** (0/1) | **能** | 17.8s |

#### 关键发现

**Qwen3.5-Flash 能正确定位 bug，但无法精确修复。**

在 7 次尝试中，Qwen **每次都正确识别了 bug**——建议在 `provisioningOfTrafficRoutingInfo()` 函数中添加 `if routeReq == nil` 检查，这与官方 PR 的修复完全一致。但没有一次成功产出可正确应用的补丁。

#### 失败模式分析

| 失败模式 | 说明 | 次数 |
|---------|------|------|
| 空响应 | thinking 模式消耗所有 token，content 为 null | 1 |
| 文件截断 | 完整文件输出超过 max_tokens | 1 |
| 插入错误位置 | 上下文行匹配到文件中其他位置 | 1 |
| 替换正确代码 | 添加修复时删除了不该删的代码 | 1 |
| 格式不匹配 | 输出格式与解析器不一致 | 3 |

#### 启示

> **Agent 能力（而非模型智能）是解决 5G 软件工程任务的关键。**

单轮 API 调用模式下，模型缺乏：
- 迭代地读取和浏览代码的能力
- 精确编辑文件特定行的能力
- 编辑后自行验证的能力

这与 SWE-Bench Mobile 的发现一致："同一模型在不同 Agent 中表现差异高达 6 倍"。后续将评测 Claude Code、Cursor 等 Agent 模式工具。

#### 运行 Qwen 评测

```bash
# 设置环境
unset http_proxy https_proxy HTTP_PROXY HTTPS_PROXY all_proxy ALL_PROXY
export DASHSCOPE_API_KEY=sk-xxxx

# 评测单个实例
python eval/run_evaluation.py \
  --agent qwen \
  --model qwen3.5-flash \
  --instances free5gc_pcf_issue879

# 评测所有实例
python eval/run_evaluation.py --agent qwen --model qwen3.5-flash

# 结果保存在
eval/results/<timestamp>/results.json
eval/results/<timestamp>/summary.txt
```

---

## 六、项目结构

```
swebench5g/
│
├── pilot_pcf_879/              ← 已验证的 pilot 任务（独立结构）
│   ├── Dockerfile
│   ├── problem_statement.md
│   ├── task_metadata.json
│   ├── test-suite/
│   └── scripts/
│
├── instances/                  ← 所有任务实例（模板化构建）
│   ├── amf_pr118/              ← 每个实例包含 4 个文件
│   ├── amf_pr157/
│   ├── amf_pr161/              ← 已验证
│   ├── amf_pr181/
│   ├── amf_pr191/
│   ├── amf_pr192/              ← medium 难度
│   ├── amf_pr196/
│   ├── ausf_pr52/
│   ├── nrf_pr78/
│   ├── nrf_pr79/
│   ├── nssf_pr39/
│   ├── nssf_pr44/
│   ├── pcf_pr57/
│   ├── pcf_pr62/
│   ├── smf_pr125/
│   ├── smf_pr128/
│   ├── smf_pr189/              ← medium 难度
│   ├── udm_pr45/
│   ├── udm_pr66/
│   └── udm_pr77/
│
├── templates/                  ← Docker/测试模板
│   ├── Dockerfile.template
│   └── run_tests.sh.template
│
├── scripts/                    ← 工具脚本
│   ├── mine_issues.py          ← Issue 挖掘（已挖出 280 个候选）
│   ├── generate_instances.py   ← 批量生成实例
│   ├── build_instance.sh       ← 构建镜像
│   └── validate_instance.sh    ← 三步验证
│
├── eval/                       ← 评估框架
│   └── run_evaluation.py
│
├── dataset/                    ← HuggingFace 数据集
│   ├── swebench5g.jsonl
│   ├── README.md
│   └── upload_to_hf.py
│
├── paper/                      ← NeurIPS D&B 论文
│   ├── main.tex
│   └── references.bib
│
├── candidates.json             ← 280 个候选（mine_issues.py 输出）
├── ROADMAP.md                  ← 项目路线图
└── tutorial_zh.md              ← 本文件
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

# 基础镜像固定为 golang:1.25-bookworm（只需拉取一次）
```

### Q: Docker build 时找不到 golang 镜像？

模板已固定使用 `golang:1.25-bookworm`。如果服务器之前没拉过：
```bash
docker pull golang:1.25-bookworm
```
之后所有实例都复用同一个基础镜像。

### Q: 测试为什么要放在同一个 Go 包内？

Go 的未导出函数（小写字母开头）只能在同一个包内访问。由于 free5GC 的很多内部函数是未导出的，测试必须声明为 `package processor`（而不是 `package processor_test`）才能直接调用。

### Q: fail_test 在 buggy 版本上也 PASS 了？

这说明测试**没有真正复现 bug**。常见错误：

```go
// ❌ 测试自己写了安全逻辑（永远不会 panic）
if ptr != nil && ptr.Field != "" { ... }

// ✅ 应该复现 buggy 代码的不安全访问
if ptr.Field != "" { ... }  // ptr 为 nil 时会 panic
```

如果函数太复杂无法直接调用，使用 **策略 B（Diff-Based Intent Test）**。

### Q: validate 的 Step 3 失败（apply fix 后 fail test 仍然 FAIL）？

可能原因：
1. 测试是模拟代码（不调用实际函数），cherry-pick 不会改变测试逻辑
2. 解决方案：改用 Diff-Based Intent Test（检查源码而非执行函数）

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
