#!/bin/bash
# run_tests.sh — SWE-Bench 5G test runner
# Usage:
#   ./run_tests.sh existing    — Run only existing tests (should always PASS)
#   ./run_tests.sh fail        — Run only fail-to-pass tests (FAIL before fix, PASS after)
#   ./run_tests.sh all         — Run all tests

set -e

PCF_DIR="/opt/free5gc-pcf"
TEST_DIR="/opt/test-suite"
PROCESSOR_DIR="${PCF_DIR}/internal/sbi/processor"

COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[1;33m'
COLOR_NC='\033[0m'

MODE="${1:-all}"

echo "============================================"
echo " SWE-Bench 5G — PCF Issue #879 Test Runner"
echo " Mode: ${MODE}"
echo "============================================"
echo ""

# Copy test files to the processor package directory
# (Tests must be in the same package to access unexported functions)
copy_tests() {
    local test_type=$1
    if [ "$test_type" == "existing" ]; then
        cp "${TEST_DIR}/existing_test.go" "${PROCESSOR_DIR}/existing_test.go"
    elif [ "$test_type" == "fail" ]; then
        cp "${TEST_DIR}/fail_test.go" "${PROCESSOR_DIR}/fail_test.go"
    elif [ "$test_type" == "all" ]; then
        cp "${TEST_DIR}/existing_test.go" "${PROCESSOR_DIR}/existing_test.go"
        cp "${TEST_DIR}/fail_test.go" "${PROCESSOR_DIR}/fail_test.go"
    fi
}

# Clean up test files after run
cleanup_tests() {
    rm -f "${PROCESSOR_DIR}/existing_test.go"
    rm -f "${PROCESSOR_DIR}/fail_test.go"
}

run_existing() {
    echo -e "${COLOR_YELLOW}--- Running EXISTING tests (should PASS) ---${COLOR_NC}"
    echo ""
    copy_tests "existing"
    cd "${PCF_DIR}"
    if go test -v -run "TestProvisioningOfTrafficRoutingInfo_ValidRouteReq|TestPccRuleLookupByAppId" ./internal/sbi/processor/ 2>&1; then
        echo ""
        echo -e "${COLOR_GREEN}EXISTING tests: ALL PASSED${COLOR_NC}"
        EXISTING_RESULT=0
    else
        echo ""
        echo -e "${COLOR_RED}EXISTING tests: FAILED (environment issue!)${COLOR_NC}"
        EXISTING_RESULT=1
    fi
    cleanup_tests
    return $EXISTING_RESULT
}

run_fail() {
    echo -e "${COLOR_YELLOW}--- Running FAIL-TO-PASS tests ---${COLOR_NC}"
    echo -e "${COLOR_YELLOW}    (These should FAIL on buggy version, PASS after fix)${COLOR_NC}"
    echo ""
    copy_tests "fail"
    cd "${PCF_DIR}"
    if go test -v -run "TestProvisioningOfTrafficRoutingInfo_NilRouteReq|TestSuppFeatTrafficRouting_NoAfRoutReq" ./internal/sbi/processor/ 2>&1; then
        echo ""
        echo -e "${COLOR_GREEN}FAIL-TO-PASS tests: ALL PASSED (bug is fixed!)${COLOR_NC}"
        FAIL_RESULT=0
    else
        echo ""
        echo -e "${COLOR_RED}FAIL-TO-PASS tests: FAILED (bug still present)${COLOR_NC}"
        FAIL_RESULT=1
    fi
    cleanup_tests
    return $FAIL_RESULT
}

case "$MODE" in
    existing)
        run_existing
        ;;
    fail)
        run_fail
        ;;
    all)
        run_existing
        E_RESULT=$?
        echo ""
        echo "--------------------------------------------"
        echo ""
        run_fail
        F_RESULT=$?
        echo ""
        echo "============================================"
        echo " SUMMARY"
        echo "============================================"
        if [ $E_RESULT -eq 0 ] && [ $F_RESULT -eq 0 ]; then
            echo -e "${COLOR_GREEN} ALL TESTS PASSED — Bug is fixed!${COLOR_NC}"
        elif [ $E_RESULT -eq 0 ] && [ $F_RESULT -ne 0 ]; then
            echo -e "${COLOR_RED} Existing: PASS | Fail-to-pass: FAIL — Bug still present${COLOR_NC}"
        else
            echo -e "${COLOR_RED} Some tests failed unexpectedly${COLOR_NC}"
        fi
        echo "============================================"
        # Exit with combined result
        exit $(( E_RESULT + F_RESULT ))
        ;;
    *)
        echo "Usage: $0 {existing|fail|all}"
        exit 1
        ;;
esac
