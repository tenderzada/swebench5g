# Figure Drawing Prompts

## Fig. 1: Overview Pipeline (TikZ in paper)

Already drawn in LaTeX using TikZ. Styled after BeyondSWE's main figure:
- Phase 1 (blue): Agent receives issue + optional 3GPP spec, interacts with Docker container
- Phase 2 (green): Fresh container → apply patch → run tests → verified result

## Fig. 2: Bar Chart (needs external tool — matplotlib/pgfplots)

### Drawing Prompt

Create a grouped bar chart with the following specifications:

**Data:**

| Model | Bug Diagnosed | Patch Applied | Resolved |
|-------|--------------|---------------|----------|
| Qwen3.5-Flash | 90% | 60% | 10% |
| Kimi-128k | 80% | 50% | 10% |
| Claude Sonnet 4 | 100% | 80% | 30% |
| GPT-4.1 | 90% | 70% | 20% |

**Style:**
- X-axis: 4 model groups
- Y-axis: 0% to 100%
- 3 bars per group:
  - "Bug Diagnosed" — light blue (#AEC6CF)
  - "Patch Applied" — medium blue (#6495ED)
  - "Resolved" — dark blue (#003366)
- Percentage labels on top of each bar
- Legend at top-right
- Clean academic style, no 3D, thin black borders
- Font: serif (Times), 9pt
- Figure width: 3.5 inches (IEEE single column)

**Optional enhancement:**
- Add a single hatched bar for "Qwen Single-Turn" (0% resolved, 30% patch, 90% diagnosed)
  next to the Qwen multi-turn group, to visualize the single→multi improvement.

### Python code (matplotlib):

```python
import matplotlib.pyplot as plt
import numpy as np

models = ['Qwen3.5\nFlash', 'Kimi\n128k', 'Claude\nSonnet 4', 'GPT\n4.1']
diagnosed = [90, 80, 100, 90]
applied   = [60, 50, 80, 70]
resolved  = [10, 10, 30, 20]

x = np.arange(len(models))
w = 0.22

fig, ax = plt.subplots(figsize=(3.5, 2.5))
b1 = ax.bar(x - w, diagnosed, w, label='Bug Diagnosed', color='#AEC6CF', edgecolor='black', linewidth=0.5)
b2 = ax.bar(x,     applied,   w, label='Patch Applied',  color='#6495ED', edgecolor='black', linewidth=0.5)
b3 = ax.bar(x + w, resolved,  w, label='Resolved',       color='#003366', edgecolor='black', linewidth=0.5)

ax.set_ylabel('Percentage (%)', fontsize=8)
ax.set_xticks(x)
ax.set_xticklabels(models, fontsize=7)
ax.set_ylim(0, 115)
ax.legend(fontsize=6, loc='upper right')
ax.tick_params(axis='y', labelsize=7)

for bars in [b1, b2, b3]:
    for bar in bars:
        h = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2, h + 1, f'{int(h)}%', ha='center', va='bottom', fontsize=5.5)

plt.tight_layout()
plt.savefig('fig_results.pdf', dpi=300, bbox_inches='tight')
plt.savefig('fig_results.png', dpi=300, bbox_inches='tight')
plt.show()
```

## Fig. 3 (optional): A/B Spec-as-Skill comparison

### Drawing Prompt

Side-by-side bar chart:

| Bug Category | -spec | +spec |
|-------------|-------|-------|
| Generic nil-check (6) | 33% | 33% |
| Spec-dependent (4) | 0% | 25% |
| Overall (10) | 20% | 30% |

Two bars per group: gray (-spec) vs blue (+spec).
Highlight the +25% delta on spec-dependent with an annotation arrow.
