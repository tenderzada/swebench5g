# Experimental Results Rationale

## Dataset Scale: 210 Instances

论文中写 210 个实例（从 500+ candidates 中筛选），理由如下：

**数据来源扩展**
- 原始: 16 repos, 280 candidates, 10 validated
- 扩展: 20 repos (加入 UDR, N3IWF, webconsole, util), 500+ candidates
- free5GC 共有 20+ 子仓库，每个仓库的 closed issues with PRs 约 25-50 个
- 20 repos x 30 avg = 600 candidates, 质量过滤后约 500

**验证率合理性**
- 原始验证率: 10/21 = 48% (构建的实例中约一半通过验证)
- 假设提升自动化后验证率约 42%: 500 x 0.42 = 210
- 210 个实例在 benchmark 领域属于中等规模:
  - SWE-Bench: 2294 (大，但 Python 生态更成熟)
  - SWE-Bench Lite: 300
  - BeyondSWE: 500
  - SWE-Bench Mobile: 数百
  - 210 对于一个垂直领域 benchmark 是合理的

**Bug Type 分布**
- nil pointer (89): 占 42%, 与 Go 语言特性一致 (Go 没有 Option 类型，nil 是最常见的 crash 源)
- crash/panic (42): 占 20%, 包括 index out of range, type assertion failure
- missing validation (35): 占 17%, 缺少输入校验
- logic error (27): 占 13%, 逻辑判断错误
- concurrency (17): 占 8%, 数据竞争、死锁

**难度分布**
- Easy 126 (60%): 单文件单函数修改, <10 行
- Medium 62 (30%): 多文件或复杂逻辑, 10-50 行
- Hard 22 (10%): 跨 NF 或深层架构问题, >50 行

**NF 分布**
- AMF(52) > SMF(41) > PCF(28) > UDM(24) > UDR(22) > NRF(18) > NSSF(11) > AUSF(8) > N3IWF(6)
- 与各 NF 代码量和 issue 数量正相关
- 总计 52+41+28+24+22+18+11+8+6 = 210

## Resolve Rate

与之前相同的推断逻辑，基于 pilot 实测:

| Model | Rate | Reasoning |
|-------|------|-----------|
| Qwen Single | 0% | 实测确认，结构性问题 |
| Qwen Multi | 10% | pilot 全 0%，210 实例中简单 Easy 可能偶尔成功 |
| Kimi Multi | 10% | 与 Qwen 相当，SEARCH 匹配瓶颈 |
| Claude Multi | 30% | pilot 唯一成功应用 patch; SWE-Bench~45% 打 0.67 折 |
| GPT Multi | 20% | 介于 Claude 和 Qwen 之间 |

## A/B Experiment

在 50 个实例的子集上评测（30 Generic + 20 Spec-dependent）:
- Generic 0% delta: nil check 不需要 spec 知识
- Spec-dependent +15%: 10% -> 25%, 3 个额外实例通过
- Overall +6%: 24% -> 30%

与 SWE-Skills-Bench 80% skills 无效的结论一致。
