#!/bin/bash
# validate_image.sh — End-to-end validation of the Docker image
# Run this OUTSIDE the container to verify the image is correctly built.
#
# Expected results:
#   1. Existing tests: PASS
#   2. Fail tests: FAIL (bug present)
#   3. After applying fix: all PASS

set -e

IMAGE_NAME="${1:-swebench5g/free5gc:pcf_issue_879}"

COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_NC='\033[0m'

echo "============================================"
echo " Validating image: ${IMAGE_NAME}"
echo "============================================"
echo ""

# Step 1: Existing tests should PASS
echo "Step 1: Running existing tests (expect PASS)..."
if docker run --rm "${IMAGE_NAME}" bash -c "/opt/test-suite/run_tests.sh existing" 2>&1; then
    echo -e "${COLOR_GREEN}Step 1 PASSED: existing tests pass on buggy version${COLOR_NC}"
else
    echo -e "${COLOR_RED}Step 1 FAILED: existing tests should pass!${COLOR_NC}"
    exit 1
fi

echo ""

# Step 2: Fail tests should FAIL
echo "Step 2: Running fail-to-pass tests (expect FAIL)..."
if docker run --rm "${IMAGE_NAME}" bash -c "/opt/test-suite/run_tests.sh fail" 2>&1; then
    echo -e "${COLOR_RED}Step 2 FAILED: fail tests should NOT pass on buggy version!${COLOR_NC}"
    exit 1
else
    echo -e "${COLOR_GREEN}Step 2 PASSED: fail tests correctly fail on buggy version${COLOR_NC}"
fi

echo ""

# Step 3: After fix, all tests should PASS
echo "Step 3: Applying fix and running all tests (expect PASS)..."
if docker run --rm "${IMAGE_NAME}" bash -c "/opt/scripts/apply_fix.sh && /opt/test-suite/run_tests.sh all" 2>&1; then
    echo -e "${COLOR_GREEN}Step 3 PASSED: all tests pass after fix${COLOR_NC}"
else
    echo -e "${COLOR_RED}Step 3 FAILED: tests should pass after fix!${COLOR_NC}"
    exit 1
fi

echo ""
echo "============================================"
echo -e "${COLOR_GREEN} IMAGE VALIDATION: ALL STEPS PASSED${COLOR_NC}"
echo "============================================"
echo ""
echo "The image is ready for SWE-Bench 5G evaluation."
echo "Agent workflow:"
echo "  1. docker run -it ${IMAGE_NAME}"
echo "  2. cat /opt/task/problem_statement.md"
echo "  3. /opt/test-suite/run_tests.sh existing   # verify env OK"
echo "  4. /opt/test-suite/run_tests.sh fail        # confirm bug exists"
echo "  5. <agent fixes the code>"
echo "  6. /opt/test-suite/run_tests.sh all         # verify fix"
