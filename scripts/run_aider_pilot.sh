#!/bin/bash
# run_aider_pilot.sh — 用 Aider + Qwen 在 pilot 实例上做 agentic 验证
#
# 使用方法：
#   # 1. 先确保 pilot Docker 镜像已构建
#   cd pilot_pcf_879 && bash scripts/build_image.sh
#
#   # 2. 设置 API Key
#   export DASHSCOPE_API_KEY=sk-xxx
#
#   # 3. 运行
#   bash scripts/run_aider_pilot.sh
#
#   # 可选：用 OpenRouter Claude
#   export OPENROUTER_API_KEY=sk-or-xxx
#   bash scripts/run_aider_pilot.sh --claude

set -e

# ── 配置 ──────────────────────────────────────────────────
IMAGE="swebench5g/free5gc:pcf_issue_879"
WORKDIR="/opt/free5gc-pcf"
RESULT_DIR="eval/results/aider_pilot_$(date +%Y%m%d_%H%M%S)"
USE_CLAUDE=false

if [ "$1" == "--claude" ]; then
    USE_CLAUDE=true
fi

mkdir -p "$RESULT_DIR"

echo "============================================"
echo " SWE-Bench 5G — Aider Pilot Verification"
echo "============================================"

# ── Step 1: 启动容器 ─────────────────────────────────────
echo "[1/6] Starting container..."
CID=$(docker run -d --rm "$IMAGE" sleep 7200)
echo "  Container: ${CID:0:12}"

cleanup() {
    echo "  Cleaning up container..."
    docker stop "$CID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

# ── Step 2: 验证环境（existing tests 应该 PASS）─────────
echo "[2/6] Verifying environment (existing tests)..."
if ! docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh existing" > "$RESULT_DIR/step2_existing.log" 2>&1; then
    echo "  ERROR: Existing tests failed. Environment broken."
    cat "$RESULT_DIR/step2_existing.log"
    exit 1
fi
echo "  Existing tests: PASS"

# ── Step 3: 确认 bug 存在（fail tests 应该 FAIL）────────
echo "[3/6] Confirming bug exists (fail-to-pass tests)..."
if docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh fail" > "$RESULT_DIR/step3_fail.log" 2>&1; then
    echo "  WARNING: Fail-to-pass tests already pass! Bug may be pre-fixed."
    exit 1
fi
echo "  Bug confirmed: fail-to-pass tests FAIL as expected"

# ── Step 4: 安装 Aider ───────────────────────────────────
echo "[4/6] Installing Aider in container..."
docker exec "$CID" bash -c "
    apt-get update -qq && apt-get install -y -qq python3-pip python3-venv > /dev/null 2>&1
    python3 -m venv /opt/aider-env
    /opt/aider-env/bin/pip install -q aider-chat
" 2>&1 | tail -3
echo "  Aider installed"

# ── Step 5: 运行 Aider 修复 bug ──────────────────────────
echo "[5/6] Running Aider to fix the bug..."

# 读取 problem statement
PROBLEM=$(docker exec "$CID" cat /opt/task/problem_statement.md)

# 构建 Aider 命令
if [ "$USE_CLAUDE" = true ]; then
    # Claude via OpenRouter
    if [ -z "$OPENROUTER_API_KEY" ]; then
        echo "  ERROR: OPENROUTER_API_KEY not set"
        exit 1
    fi
    MODEL_NAME="openrouter/anthropic/claude-sonnet-4-6"
    API_KEY_VAR="OPENROUTER_API_KEY"
    API_KEY_VAL="$OPENROUTER_API_KEY"
    AGENT_LABEL="aider+claude-sonnet-4-6"
    echo "  Model: Claude Sonnet 4.6 via OpenRouter"
else
    # Qwen via DashScope
    if [ -z "$DASHSCOPE_API_KEY" ]; then
        echo "  ERROR: DASHSCOPE_API_KEY not set"
        exit 1
    fi
    MODEL_NAME="openai/qwen3.5-flash"
    API_KEY_VAR="DASHSCOPE_API_KEY"
    API_KEY_VAL="$DASHSCOPE_API_KEY"
    AGENT_LABEL="aider+qwen3.5-flash"
    echo "  Model: Qwen3.5-Flash via DashScope"
fi

# Aider 需要 git repo
docker exec "$CID" bash -c "cd $WORKDIR && git config user.email 'aider@test' && git config user.name 'aider'"

START_TIME=$(date +%s)

# 运行 Aider（全自动模式）
# --yes: 自动确认所有操作
# --no-auto-commits: 不自动 commit（我们用 git diff 提取 patch）
# --map-tokens 1024: 限制 repo map token 数
docker exec -e "$API_KEY_VAR=$API_KEY_VAL" "$CID" bash -c "
    cd $WORKDIR && /opt/aider-env/bin/aider \
        --model $MODEL_NAME \
        $([ '$USE_CLAUDE' = false ] && echo '--openai-api-base https://dashscope.aliyuncs.com/compatible-mode/v1') \
        --yes \
        --no-auto-commits \
        --map-tokens 1024 \
        --message '$PROBLEM'
" > "$RESULT_DIR/aider_output.log" 2>&1 || true

END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))
echo "  Aider finished in ${ELAPSED}s"

# ── Step 6: 提取 patch 并验证 ────────────────────────────
echo "[6/6] Extracting patch and running tests..."

# 提取 patch
PATCH=$(docker exec "$CID" bash -c "cd $WORKDIR && git diff")
echo "$PATCH" > "$RESULT_DIR/patch.diff"

if [ -z "$PATCH" ]; then
    echo "  WARNING: Aider produced no changes"
    RESOLVED=false
else
    echo "  Patch extracted ($(echo "$PATCH" | wc -l) lines)"
    echo ""
    echo "--- PATCH ---"
    echo "$PATCH"
    echo "--- END PATCH ---"
    echo ""
fi

# 跑 existing tests
echo "  Running existing tests..."
if docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh existing" > "$RESULT_DIR/step6_existing.log" 2>&1; then
    EXISTING_PASS=true
    echo "  Existing tests: PASS"
else
    EXISTING_PASS=false
    echo "  Existing tests: FAIL (regression!)"
fi

# 跑 fail-to-pass tests
echo "  Running fail-to-pass tests..."
if docker exec "$CID" bash -c "/opt/test-suite/run_tests.sh fail" > "$RESULT_DIR/step6_fail.log" 2>&1; then
    FAIL_PASS=true
    echo "  Fail-to-pass tests: PASS (bug fixed!)"
else
    FAIL_PASS=false
    echo "  Fail-to-pass tests: FAIL (bug still present)"
fi

# 判断是否 resolved
if [ "$EXISTING_PASS" = true ] && [ "$FAIL_PASS" = true ]; then
    RESOLVED=true
else
    RESOLVED=false
fi

# ── 输出结果 ──────────────────────────────────────────────
echo ""
echo "============================================"
echo " RESULT"
echo "============================================"
if [ "$RESOLVED" = true ]; then
    echo "  Status:       RESOLVED"
else
    echo "  Status:       NOT RESOLVED"
fi
echo "  Agent:        $AGENT_LABEL"
echo "  Instance:     pcf_issue_879"
echo "  Time:         ${ELAPSED}s"
echo "  Existing:     $([ '$EXISTING_PASS' = true ] && echo 'PASS' || echo 'FAIL')"
echo "  Fail-to-pass: $([ '$FAIL_PASS' = true ] && echo 'PASS' || echo 'FAIL')"
echo "  Patch lines:  $(echo "$PATCH" | wc -l)"
echo "  Results dir:  $RESULT_DIR/"
echo "============================================"

# 保存 JSON 结果
cat > "$RESULT_DIR/result.json" << ENDJSON
{
    "instance_id": "pcf_issue_879",
    "agent": "aider",
    "model": "$([ '$USE_CLAUDE' = true ] && echo 'claude-sonnet-4-6' || echo 'qwen3.5-flash')",
    "resolved": $RESOLVED,
    "existing_tests_pass": $EXISTING_PASS,
    "fail_tests_pass": $FAIL_PASS,
    "time_seconds": $ELAPSED,
    "timestamp": "$(date -Iseconds)"
}
ENDJSON

echo "  result.json saved"
