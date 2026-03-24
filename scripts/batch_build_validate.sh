#!/bin/bash
# batch_build_validate.sh — Build and validate all instances in one go.
#
# Usage:
#   bash scripts/batch_build_validate.sh                    # all instances
#   bash scripts/batch_build_validate.sh amf_pr118 pcf_pr57 # specific ones
#   bash scripts/batch_build_validate.sh --clean            # build, validate, remove image
#
# Options:
#   --clean     Remove Docker image after successful validation (save disk)
#   --build-only   Only build, skip validation
#   --validate-only  Only validate (image must exist)

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
INSTANCES_DIR="${PROJECT_ROOT}/instances"
LOG_DIR="${PROJECT_ROOT}/tmp"
mkdir -p "${LOG_DIR}"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

CLEAN=false
BUILD_ONLY=false
VALIDATE_ONLY=false
TARGETS=()

# Parse args
for arg in "$@"; do
    case "$arg" in
        --clean) CLEAN=true ;;
        --build-only) BUILD_ONLY=true ;;
        --validate-only) VALIDATE_ONLY=true ;;
        *) TARGETS+=("$arg") ;;
    esac
done

# If no targets specified, find all instances with instance.json
if [ ${#TARGETS[@]} -eq 0 ]; then
    for dir in "${INSTANCES_DIR}"/*/; do
        if [ -f "${dir}instance.json" ]; then
            TARGETS+=("$(basename "$dir")")
        fi
    done
fi

TOTAL=${#TARGETS[@]}
PASS=0
FAIL=0
SKIP=0
RESULTS=()

echo "============================================"
echo " SWE-Bench 5G Batch Build & Validate"
echo " Instances: ${TOTAL}"
echo " Clean mode: ${CLEAN}"
echo "============================================"
echo ""

for i in "${!TARGETS[@]}"; do
    ID="${TARGETS[$i]}"
    NUM=$((i + 1))
    CONFIG="${INSTANCES_DIR}/${ID}/instance.json"
    IMAGE="swebench5g/free5gc:${ID}"

    echo -e "${YELLOW}[${NUM}/${TOTAL}] ${ID}${NC}"

    # Check config exists
    if [ ! -f "$CONFIG" ]; then
        echo -e "  ${RED}SKIP: instance.json not found${NC}"
        SKIP=$((SKIP + 1))
        RESULTS+=("${ID}: SKIP (no config)")
        echo ""
        continue
    fi

    # Check for PLACEHOLDER commits
    if grep -q "PLACEHOLDER" "$CONFIG" 2>/dev/null; then
        echo -e "  ${RED}SKIP: commit SHAs not yet filled in${NC}"
        SKIP=$((SKIP + 1))
        RESULTS+=("${ID}: SKIP (placeholder commits)")
        echo ""
        continue
    fi

    # Build
    if [ "$VALIDATE_ONLY" = false ]; then
        echo "  Building..."
        rm -rf "${INSTANCES_DIR}/${ID}/.build"
        if bash "${SCRIPT_DIR}/build_instance.sh" "$CONFIG" > "${LOG_DIR}/swebench5g_build_${ID}.log" 2>&1; then
            echo -e "  ${GREEN}Build: OK${NC}"
        else
            echo -e "  ${RED}Build: FAILED${NC}"
            echo "  Log: ${LOG_DIR}/swebench5g_build_${ID}.log"
            tail -5 "${LOG_DIR}/swebench5g_build_${ID}.log" | sed 's/^/  /'
            FAIL=$((FAIL + 1))
            RESULTS+=("${ID}: BUILD FAILED")
            echo ""
            continue
        fi
    fi

    # Validate
    if [ "$BUILD_ONLY" = false ]; then
        echo "  Validating..."
        if bash "${SCRIPT_DIR}/validate_instance.sh" "$CONFIG" > "${LOG_DIR}/swebench5g_validate_${ID}.log" 2>&1; then
            echo -e "  ${GREEN}Validate: PASSED${NC}"
            PASS=$((PASS + 1))
            RESULTS+=("${ID}: PASSED")
        else
            echo -e "  ${RED}Validate: FAILED${NC}"
            echo "  Log: ${LOG_DIR}/swebench5g_validate_${ID}.log"
            grep -E "Step [0-9]|FAILED|PASSED|panic|Error" "${LOG_DIR}/swebench5g_validate_${ID}.log" | tail -5 | sed 's/^/  /'
            FAIL=$((FAIL + 1))
            RESULTS+=("${ID}: VALIDATE FAILED")
        fi
    else
        PASS=$((PASS + 1))
        RESULTS+=("${ID}: BUILT (not validated)")
    fi

    # Clean up image if requested
    if [ "$CLEAN" = true ]; then
        docker rmi "$IMAGE" > /dev/null 2>&1 && echo "  Cleaned image"
    fi

    # Clean up build context
    rm -rf "${INSTANCES_DIR}/${ID}/.build"

    echo ""
done

# Summary
echo "============================================"
echo " SUMMARY"
echo "============================================"
echo -e " Total:   ${TOTAL}"
echo -e " ${GREEN}Passed:  ${PASS}${NC}"
echo -e " ${RED}Failed:  ${FAIL}${NC}"
echo -e " Skipped: ${SKIP}"
echo ""
echo " Results:"
for r in "${RESULTS[@]}"; do
    if echo "$r" | grep -q "PASSED"; then
        echo -e "   ${GREEN}${r}${NC}"
    elif echo "$r" | grep -q "FAILED"; then
        echo -e "   ${RED}${r}${NC}"
    else
        echo -e "   ${YELLOW}${r}${NC}"
    fi
done
echo "============================================"
echo ""
echo " Build logs: ${LOG_DIR}/swebench5g_build_*.log"
echo " Validate logs: ${LOG_DIR}/swebench5g_validate_*.log"

# Exit with failure if any failed
[ $FAIL -eq 0 ]
