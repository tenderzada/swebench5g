#!/bin/bash
# validate_instance.sh — Validate a built SWE-Bench 5G Docker image.
#
# Three-step validation:
#   1. Existing tests PASS on buggy commit
#   2. Fail-to-pass tests FAIL on buggy commit
#   3. After applying fix: ALL tests PASS
#
# Usage:
#   bash scripts/validate_instance.sh instances/pcf_pr65/instance.json

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <instance.json>"
    exit 1
fi

CONFIG="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"

parse() { python3 -c "import json,sys; d=json.load(open('$CONFIG')); print(d.get('$1','$2'))"; }

INSTANCE_ID=$(parse instance_id "")
FIX_COMMIT=$(parse full_fix_commit "")
IMAGE_NAME="swebench5g/free5gc:${INSTANCE_ID}"
REPO_NAME=$(parse repo "" | cut -d/ -f2)
WORKDIR="/opt/free5gc-${REPO_NAME}"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "============================================"
echo " Validating: ${IMAGE_NAME}"
echo "============================================"
echo ""

# Step 1: Existing tests should PASS
echo "Step 1/3: Existing tests (expect PASS)..."
if docker run --rm "${IMAGE_NAME}" /opt/test-suite/run_tests.sh existing 2>&1; then
    echo -e "${GREEN}Step 1 PASSED${NC}"
else
    echo -e "${RED}Step 1 FAILED: existing tests should pass on buggy version${NC}"
    exit 1
fi
echo ""

# Step 2: Fail tests should FAIL
echo "Step 2/3: Fail-to-pass tests (expect FAIL)..."
if docker run --rm "${IMAGE_NAME}" /opt/test-suite/run_tests.sh fail 2>&1; then
    echo -e "${RED}Step 2 FAILED: fail tests should NOT pass on buggy version${NC}"
    exit 1
else
    echo -e "${GREEN}Step 2 PASSED: fail tests correctly fail${NC}"
fi
echo ""

# Step 3: After fix, all should PASS
echo "Step 3/3: Apply fix and run all tests (expect ALL PASS)..."
FIX_CMD="cd ${WORKDIR} && git cherry-pick --no-commit ${FIX_COMMIT}"
if docker run --rm "${IMAGE_NAME}" bash -c "${FIX_CMD} && /opt/test-suite/run_tests.sh all" 2>&1; then
    echo -e "${GREEN}Step 3 PASSED${NC}"
else
    echo -e "${RED}Step 3 FAILED: all tests should pass after fix${NC}"
    exit 1
fi

echo ""
echo "============================================"
echo -e "${GREEN} VALIDATION PASSED: ${INSTANCE_ID}${NC}"
echo "============================================"
