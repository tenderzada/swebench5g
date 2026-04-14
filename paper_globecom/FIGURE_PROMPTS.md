# Gemini API Drawing Prompts

## Fig. 1: System Overview (Pipeline Diagram)

### Gemini Prompt

```
Draw a two-phase system architecture diagram for an AI coding agent evaluation pipeline,
styled for an IEEE academic paper. Use a clean, professional design with a blue color scheme.

Phase 1 (left side, light blue background #E8F0FE, labeled "Phase 1: Agent Coding"):
- A rounded box "Issue / Problem Statement" in light orange (#FFF3E0)
  with a solid blue arrow pointing right to
- A central rounded box "AI Agent (LLM)" in light yellow (#FFF9C4)
- Above the agent, a dashed-border rounded box "3GPP Spec (Optional)" in light blue (#E3F2FD)
  with a gray dashed arrow pointing down to the agent
- Below the agent, a rounded box "Docker Container (NF Source Code)" in light blue (#E3F2FD)
  with bidirectional solid blue arrows connecting to the agent
- A solid blue arrow from the agent pointing right to a rounded box "Proposed Patch" in light orange

Phase 2 (right side, light green background #E8F5E9, labeled "Phase 2: Rigorous Evaluation"):
- Arrow from "Proposed Patch" to "Fresh Docker Container (Clean Env)"
- Arrow down to "Apply Patch"
- Arrow down to "Run Test Suite (P2P and F2P)"
- Arrow down to "Verified Result" with a small checkmark icon

Style requirements:
- Rounded rectangles with thin borders
- Clean sans-serif font (similar to Helvetica)
- No 3D effects, no drop shadows
- Solid blue arrows (#1565C0) with arrowheads
- Dashed gray arrow (#9E9E9E) for the optional spec connection
- White overall background
- Phase labels in bold at the top of each phase region
- Similar to the BeyondSWE/SearchSWE architecture diagram style from academic papers
- Figure should be landscape orientation, approximately 3.5 inches wide
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
