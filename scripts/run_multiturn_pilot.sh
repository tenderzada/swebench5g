#!/bin/bash
# run_multiturn_pilot.sh — Multi-turn agent on pilot instance
#
# Usage:
#   # Qwen via DashScope
#   export DASHSCOPE_API_KEY=sk-xxx
#   bash scripts/run_multiturn_pilot.sh qwen
#
#   # GPT via OpenRouter
#   export OPENROUTER_API_KEY=sk-or-xxx
#   bash scripts/run_multiturn_pilot.sh gpt
#
#   # Claude via OpenRouter
#   export OPENROUTER_API_KEY=sk-or-xxx
#   bash scripts/run_multiturn_pilot.sh claude
#
#   # Kimi (Moonshot) via Moonshot API
#   export MOONSHOT_API_KEY=sk-xxx
#   bash scripts/run_multiturn_pilot.sh kimi
#
#   # Custom: any OpenAI-compatible API
#   export AGENT_API_KEY=xxx AGENT_BASE_URL=https://... AGENT_MODEL=xxx
#   bash scripts/run_multiturn_pilot.sh custom

set -e

MODE="${1:-qwen}"
IMAGE="swebench5g/free5gc:pcf_issue_879"
WORKDIR="/opt/free5gc-pcf"
MAX_TURNS="${AGENT_MAX_TURNS:-5}"
RESULT_DIR="eval/results/multiturn_${MODE}_$(date +%Y%m%d_%H%M%S)"

mkdir -p "$RESULT_DIR"

# Configure API based on mode
case "$MODE" in
    qwen)
        if [ -z "$DASHSCOPE_API_KEY" ]; then echo "ERROR: DASHSCOPE_API_KEY not set"; exit 1; fi
        export AGENT_API_KEY="$DASHSCOPE_API_KEY"
        export AGENT_BASE_URL="https://dashscope.aliyuncs.com/compatible-mode/v1"
        export AGENT_MODEL="qwen3.5-flash"
        ;;
    gpt)
        if [ -z "$OPENROUTER_API_KEY" ]; then echo "ERROR: OPENROUTER_API_KEY not set"; exit 1; fi
        export AGENT_API_KEY="$OPENROUTER_API_KEY"
        export AGENT_BASE_URL="https://openrouter.ai/api/v1"
        export AGENT_MODEL="openai/gpt-4.1"
        ;;
    claude)
        if [ -z "$OPENROUTER_API_KEY" ]; then echo "ERROR: OPENROUTER_API_KEY not set"; exit 1; fi
        export AGENT_API_KEY="$OPENROUTER_API_KEY"
        export AGENT_BASE_URL="https://openrouter.ai/api/v1"
        export AGENT_MODEL="anthropic/claude-sonnet-4-6"
        ;;
    kimi)
        if [ -z "$MOONSHOT_API_KEY" ]; then echo "ERROR: MOONSHOT_API_KEY not set"; exit 1; fi
        export AGENT_API_KEY="$MOONSHOT_API_KEY"
        export AGENT_BASE_URL="https://api.moonshot.cn/v1"
        export AGENT_MODEL="moonshot-v1-128k"
        ;;
    custom)
        if [ -z "$AGENT_API_KEY" ] || [ -z "$AGENT_BASE_URL" ] || [ -z "$AGENT_MODEL" ]; then
            echo "ERROR: Set AGENT_API_KEY, AGENT_BASE_URL, AGENT_MODEL"; exit 1
        fi
        ;;
    *) echo "Usage: $0 {qwen|gpt|claude|kimi|custom}"; exit 1 ;;
esac

export AGENT_MAX_TURNS="$MAX_TURNS"

echo "============================================"
echo " SWE-Bench 5G — Multi-Turn Pilot"
echo " Model:     $AGENT_MODEL"
echo " Max turns: $MAX_TURNS"
echo "============================================"

# Step 1: Start container
echo "[1/6] Starting container..."
CID=$(docker run -d --rm "$IMAGE" sleep 7200)
echo "  Container: ${CID:0:12}"
cleanup() { echo "  Cleaning up..."; docker stop "$CID" >/dev/null 2>&1 || true; }
trap cleanup EXIT

# Step 2: Verify environment
echo "[2/6] Verifying environment..."
if ! docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh existing" > "$RESULT_DIR/step2.log" 2>&1; then
    echo "  ERROR: Existing tests failed"; exit 1
fi
echo "  Existing tests: PASS"

# Step 3: Confirm bug
echo "[3/6] Confirming bug..."
if docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh fail" > "$RESULT_DIR/step3.log" 2>&1; then
    echo "  WARNING: Bug already fixed?"; exit 1
fi
echo "  Bug confirmed (fail-to-pass tests FAIL)"

# Step 4: Install openai SDK in container
echo "[4/6] Installing dependencies in container..."
docker exec "$CID" bash -c "
    apt-get update -qq && apt-get install -y -qq python3-pip > /dev/null 2>&1
    pip3 install -q openai > /dev/null 2>&1
" 2>&1 | tail -2
echo "  Done"

# Step 5: Run multi-turn agent via evaluation harness
echo "[5/6] Running multi-turn agent ($AGENT_MODEL, max $MAX_TURNS turns)..."
START_TIME=$(date +%s)

# The multi-turn agent runs on the HOST (not in container), calling Docker exec
# So we just invoke the Python evaluation directly
python3 -c "
import os, sys, json, time
sys.path.insert(0, 'eval')
from run_evaluation import start_container, stop_container, _run_multiturn, run_tests_in_container, extract_patch

container_id = '$CID'
workdir = '$WORKDIR'

# Read problem statement
import subprocess
r = subprocess.run(['docker', 'exec', container_id, 'cat', '/opt/task/problem_statement.md'],
                   capture_output=True, text=True)
problem = r.stdout

print('  Starting multi-turn loop...')
result = _run_multiturn(container_id, '$AGENT_MODEL', workdir, problem, 1800)

if len(result) == 3:
    success, error, token_info = result
else:
    success, error = result
    token_info = {}

# Extract patch
r2 = subprocess.run(['docker', 'exec', container_id, 'bash', '-c', f'cd {workdir} && git diff'],
                     capture_output=True, text=True)
patch = r2.stdout

print(f'  Success: {success}')
print(f'  Error: {error}')
print(f'  Tokens: {token_info}')

# Save result
result_data = {
    'success': success,
    'error': error,
    'tokens': token_info,
    'patch': patch,
    'model': '$AGENT_MODEL',
    'max_turns': $MAX_TURNS,
}
with open('$RESULT_DIR/agent_result.json', 'w') as f:
    json.dump(result_data, f, indent=2)

with open('$RESULT_DIR/patch.diff', 'w') as f:
    f.write(patch)
" 2>&1 | tee "$RESULT_DIR/agent_output.log"

END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))

# Step 6: Final verification
echo "[6/6] Final verification..."
PATCH=$(docker exec "$CID" bash -c "cd $WORKDIR && git diff")

if [ -z "$PATCH" ]; then
    echo "  No changes made"
    EXISTING_PASS=false
    FAIL_PASS=false
else
    echo "  Patch:"
    echo "$PATCH" | head -30
    echo ""

    docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh existing" > "$RESULT_DIR/step6_existing.log" 2>&1 && EXISTING_PASS=true || EXISTING_PASS=false
    docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh fail" > "$RESULT_DIR/step6_fail.log" 2>&1 && FAIL_PASS=true || FAIL_PASS=false
fi

if [ "$EXISTING_PASS" = true ] && [ "$FAIL_PASS" = true ]; then RESOLVED=true; else RESOLVED=false; fi

echo ""
echo "============================================"
echo " RESULT"
echo "============================================"
echo "  Status:       $([ "$RESOLVED" = true ] && echo 'RESOLVED' || echo 'NOT RESOLVED')"
echo "  Agent:        multiturn"
echo "  Model:        $AGENT_MODEL"
echo "  Max turns:    $MAX_TURNS"
echo "  Instance:     pcf_issue_879"
echo "  Time:         ${ELAPSED}s"
echo "  Existing:     $([ "$EXISTING_PASS" = true ] && echo 'PASS' || echo 'FAIL')"
echo "  Fail-to-pass: $([ "$FAIL_PASS" = true ] && echo 'PASS' || echo 'FAIL')"
echo "  Results dir:  $RESULT_DIR/"
echo "============================================"
