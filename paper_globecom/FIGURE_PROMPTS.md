# Gemini API Drawing Prompts

## Fig. 1: System Overview (Pipeline Diagram)

Reference style: BeyondSWE/SearchSWE architecture diagram (two-phase, left-right layout, blue color scheme, clean academic style).

### Gemini Prompt

```
Draw a two-phase evaluation pipeline diagram for an IEEE conference paper.
Clean, professional, blue color scheme, white background. No 3D, no shadows.

LAYOUT: Left-right, two phases.

PHASE 1 (left, light blue background, label: "Phase 1: Agent Coding"):

Left inputs, two dashed-border boxes stacked vertically:
  - "Issue / Problem Statement"
  - "3GPP Spec Excerpt (Optional)"
Both have arrows pointing right into the central element.

Center: a large rounded box "AI Coding Agent" with a bold border.
This is the visually dominant element. A curved return arrow on top
represents the multi-turn iterative feedback loop.

Below the agent: a rounded box with a Docker whale icon,
labeled "Docker Container (NF Source Code)".
Connected to the agent with bidirectional arrows.

Right output: arrow from the agent to a dashed-border box "Proposed Patch".

PHASE 2 (right, light green background, label: "Phase 2: Evaluation"):

Vertical flow, top to bottom, connected by downward arrows:
1. "Fresh Docker Container (Clean Env)" with Docker icon
2. "Apply Patch"
3. "Run Test Suite (P2P & F2P)"
4. "Verified Result"

Arrow from "Proposed Patch" connects Phase 1 to Phase 2.

STYLE:
- Rounded rectangles, thin borders (1pt)
- Sans-serif font, 9-10pt
- Blue arrows (#1565C0) with arrowheads
- Dashed gray (#9E9E9E) for optional elements
- Subtle background washes for phases
- Approximately 7 inches wide, 3 inches tall
```

## Fig. 2: Bar Chart (Already generated as fig_results.pdf)

Already generated via matplotlib. If regeneration needed:

### Gemini Prompt

```
Create a grouped bar chart for an IEEE academic paper with the following data:

Four model groups on X-axis: Qwen3.5-Flash, Kimi-128k, Claude Sonnet 4, GPT-4.1

Three bars per group with these values:
- "Bug Diagnosed" (light steel blue #AEC6CF): 90%, 80%, 100%, 90%
- "Patch Applied" (cornflower blue #6495ED): 60%, 50%, 80%, 70%
- "Resolved" (dark navy #003366): 10%, 10%, 30%, 20%

Style:
- Y-axis from 0% to 100%, labeled "Percentage (%)"
- Percentage labels on top of each bar
- Legend at top-right corner
- Serif font (Times New Roman), 8-9pt size
- Clean academic style with no 3D effects
- Thin black borders on bars (0.4pt)
- Remove top and right spines
- Figure size approximately 3.5 x 2.3 inches
- White background
```

## Fig. 3 (Optional): A/B Spec-as-Skill Comparison

### Gemini Prompt

```
Create a grouped bar chart comparing specification injection results for an IEEE paper.

Three groups on X-axis: "Generic Bugs (30)", "Spec-Dependent (20)", "Overall (50)"

Two bars per group:
- "Without Spec" (gray #9E9E9E): 33%, 10%, 24%
- "With Spec" (blue #1565C0): 33%, 25%, 30%

Style:
- Y-axis from 0% to 40%, labeled "Resolve Rate (%)"
- Percentage labels on top of each bar
- Add an annotation arrow or bracket on the "Spec-Dependent" group
  showing "+15%" delta between the two bars
- Serif font (Times New Roman)
- Clean academic style, no 3D
- Figure size approximately 3.5 x 2.0 inches
- White background
```
