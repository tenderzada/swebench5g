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



````
A clean, professional two-phase evaluation pipeline diagram for an IEEE conference paper. Flat 2D vector style, white background, no shadows, no gradients, no 3D. Horizontal layout (~7in × 3in), split into two side-by-side panels with subtle background washes and a thin divider. Every box has a small monochrome line icon (#1565C0, ~16px) on its left side, vertically centered with the text.
Left panel — "Phase 1: Agent Coding" (light blue wash #E8F1FB):
On the far left, two vertically stacked input boxes with dashed gray borders (#9E9E9E):

📄 "Issue / Problem Statement" — icon: a document with a small bug/exclamation mark
📘 "3GPP Spec Excerpt (Optional)" — icon: an open book / spec sheet

Solid blue arrows (#1565C0) flow right into a large, visually dominant rounded rectangle labeled "AI Coding Agent" (bold 1.5pt border, centered) — icon: a robot head or sparkle + code brackets </> symbol. A curved self-loop arrow sits above this box with a small 🔄 circular-arrow icon to indicate the multi-turn iterative feedback loop.
Directly below the agent, a rounded box labeled "Docker Container (NF Source Code)" with the Docker whale icon on the left; connected to the agent by a bidirectional vertical arrow.
A solid arrow exits the agent to the right into a dashed-border box "Proposed Patch" — icon: a diff/patch symbol (two overlapping horizontal lines with +/−) or a bandage.
Right panel — "Phase 2: Evaluation" (light green wash #E9F5EC):
A top-to-bottom vertical flow of four rounded rectangles connected by downward solid blue arrows:

"Fresh Docker Container (Clean Env)" — Docker whale icon + a small ✨ sparkle to hint "clean/fresh"
"Apply Patch" — icon: a wrench or a patch/bandage symbol
"Run Test Suite (P2P & F2P)" — icon: a checklist with checkmarks, or a beaker/flask
"Verified Result" — icon: a shield with a checkmark ✅, slightly green-tinted (#2E7D32) to signal success

A single horizontal arrow connects "Proposed Patch" (end of Phase 1) to the top of Phase 2, crossing the panel divider.
Style constraints: all boxes are rounded rectangles with 1pt thin borders; sans-serif font (Inter or Helvetica), 9–10pt body text, 11pt bold phase labels; all icons are line-style, monochrome #1565C0 (except the final success shield in #2E7D32), ~16px, consistent stroke weight (1.5px); arrows are solid #1565C0 with small filled arrowheads; dashed #9E9E9E reserved for optional or artifact boxes; generous whitespace, aligned grid, no decorative elements beyond the specified icons.
````



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
