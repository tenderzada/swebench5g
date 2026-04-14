# Gemini API Drawing Prompts

## Fig. 1: System Overview (Pipeline Diagram)

Reference style: BeyondSWE/SearchSWE architecture diagram (two-phase, left-right layout, blue color scheme, clean academic style).

### Gemini Prompt

```
Draw a two-phase evaluation pipeline diagram for a research paper titled "SWE-Bench 5G".
The diagram should be clean, professional, and styled for an IEEE conference paper.
Use a blue color scheme on a white background. No 3D effects, no drop shadows.

LAYOUT: Left-right, two phases separated visually.

PHASE 1 (left side, light blue background #E8F0FE, label at top: "Phase 1: Agent Coding"):

Top-left corner: a small scroll/document icon labeled "Prevent cheating blocklist" 
with a dashed arrow pointing down toward the input.

Left input: a dashed-border box "Issue / Problem Statement" in light tan.
Below it or beside it: a dashed-border box "3GPP Spec Excerpt (Optional)" in light tan.
Both have arrows pointing right into the central element.

Center: a large rounded box "SWE-Bench 5G Agent" in white with a bold border. 
This is the main element of Phase 1.

Below the agent: a rounded box "Docker Container (Local Context)" containing 
an icon of a whale (Docker logo) and text "NF Source Code". 
The agent and the Docker container are connected with bidirectional arrows:
  - Down arrow labeled "Exec Commands"
  - Up arrow labeled "Exec Outputs"

The agent also has bidirectional arrows going up to a region labeled 
"Iterative Reasoning & Feedback" (shown as a curved return arrow or loop icon),
representing the multi-turn feedback loop (up to K turns).

Right output of the agent: arrow pointing right to a dashed-border box 
"Proposed Patch" in light tan.

PHASE 2 (right side, light green background #E8F5E9, label at top: "Phase 2: Rigorous Evaluation"):

Arranged vertically, top to bottom, connected by downward arrows:
1. "Fresh Docker Container (Clean Env)" with a Docker whale icon
2. "Apply Patch"
3. "Run Test Suite (P2P & F2P)"  
4. "Verified Result" with a checkmark icon

An arrow connects "Proposed Patch" from Phase 1 to the top of Phase 2.

STYLE:
- Rounded rectangles with thin borders (1pt)
- Font: clean sans-serif (Helvetica or similar), 9-10pt
- Arrows: solid blue (#1565C0) with filled arrowheads, 1.5pt width
- Dashed elements use gray (#9E9E9E) dashed lines
- Phase background colors are subtle washes, not solid fills
- The central "SWE-Bench 5G Agent" box should be the visually dominant element
- Overall dimensions: approximately 7 inches wide x 3 inches tall (IEEE double-column width)
- The diagram should closely follow the layout of the SearchSWE architecture diagram 
  from BeyondSWE (arxiv 2603.03194), but adapted for the 5G domain:
  replace "Search Tool / Browser Tool" with "3GPP Spec Excerpt"
  replace "SearchSWE (Agent)" with "SWE-Bench 5G Agent"
  keep the Docker container and two-phase structure
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
