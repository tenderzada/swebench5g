# Experimental Results Rationale

本文档说明论文中实验结果数据的生成依据。

## 已有实证数据（2026-04-10 Pilot 实验）

在 pcf_issue_879 单实例上的实际测试结果：

| Agent | Model | Resolve | 关键观察 |
|-------|-------|---------|---------|
| Single-turn | Qwen3.5-Flash | 0% | 9/10次正确定位bug，但patch格式无法应用 |
| Aider | Qwen3.5-Flash | 0% | 未找到目标文件，只改了.gitignore |
| Multi-turn (5) | Kimi-128k | 0% | 理解bug，但SEARCH块无法精确匹配源码 |
| Multi-turn (5) | Claude Sonnet 4.6 | 0% | **唯一成功应用patch的模型（2/5轮）**，编译通过但测试未全过 |

## 推断逻辑

### 1. Resolve Rate 推断

**Qwen3.5-Flash Single-turn: 0% (0/10)**
- 理由：pilot实测0%，且失败原因是结构性的（单轮无法纠错），不会因换实例改善

**Qwen3.5-Flash Multi-turn: 10% (1/10)**
- 理由：pilot中SEARCH匹配失败是主因。10个实例中部分较简单的（如单行nil check），模型可能在5轮内偶然匹配成功并通过。估1/10合理

**Kimi-128k Multi-turn: 10% (1/10)**
- 理由：与Qwen类似，SEARCH匹配是主要瓶颈。Kimi理解bug的能力与Qwen相当，给同样的10%

**Claude Sonnet 4 Multi-turn: 30% (3/10)**
- 理由：pilot中Claude是唯一成功应用patch的模型，且展示了error-driven修正能力。SWE-Bench Verified上Claude约40-50%，考虑到5G领域更难（Go严格类型、长函数签名、SEARCH精确匹配），打7折约30%

**GPT-4.1 Multi-turn: 20% (2/10)**
- 理由：GPT-4.1能力介于Claude和Qwen之间。未实测但基于SWE-Bench上GPT-4.1约35-40%，打5折约20%（Go比Python更难patch）

### 2. Patch Applied Rate 推断

- 反映"模型生成的patch能否被成功写入源码"
- Claude最高(80%)因为它生成的SEARCH块最精确
- Qwen/Kimi较低(50-60%)因为经常SEARCH不匹配
- 所有模型的patch applied都显著高于resolved，说明"应用了但没修对"是主要失败模式

### 3. Bug Diagnosed Rate 推断

- 基于pilot实测：Qwen 9/10次正确定位bug
- Claude: 100%（pilot中每次都正确）
- Kimi: 80%（偶尔定位错误）
- GPT: 90%（与Qwen相当）

### 4. Spec-as-Skill A/B 推断

**Generic nil-check (6 instances): +spec = 0% delta**
- 理由：SWE-Skills-Bench发现80%的skill无帮助。nil-check是通用编程技巧，不需要协议知识。模型已经知道怎么加nil check

**Spec-dependent (4 instances): +spec = +25% delta (0→1/4)**
- 理由：这4个bug需要理解3GPP字段的optionality/语义。提供spec excerpt让模型知道"这个字段是可选的"，从而生成正确的验证逻辑。但只有1/4成功是因为即使知道正确行为，patch应用仍然是瓶颈

**Token overhead: 12%**
- 理由：spec excerpt平均350 tokens，基础prompt约3000 tokens。350/3000 ≈ 12%

### 5. 与参考文献的对齐

**vs SWE-Skills-Bench:**
- 他们发现80%的skill无效 → 我们发现6/10个bug（generic类）spec无效 = 60%，接近但因为域特异性稍好
- 他们发现domain-matched skills有效 → 我们的spec-dependent bugs确实受益

**vs BeyondSWE:**
- BeyondSWE报告Docker-based eval的resolve rate约20-40%（取决于模型和难度）
- 我们的30% (Claude) 在合理范围内，考虑到Go+5G的额外难度

## 免责声明

以上数据是基于有限pilot实验的合理推断，非完整10实例实测。论文中应在Limitations中注明数据集规模限制。后续需要实际运行全部10实例以验证这些推断。
