#!/bin/bash
# build_instance.sh — Build a SWE-Bench 5G task instance from a config file.
#
# Usage:
#   bash scripts/build_instance.sh instances/pcf_pr65/instance.json
#
# The instance.json specifies everything needed:
#   repo, parent_commit, fix_commit, test files, etc.
#
# Steps:
#   1. Clone repo on host (if not cached)
#   2. Generate Dockerfile from template
#   3. Copy test files
#   4. Build Docker image
#   5. Validate (existing PASS, fail FAIL, fix→ALL PASS)

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 <instance.json>"
    exit 1
fi

CONFIG="$(cd "$(dirname "$1")" && pwd)/$(basename "$1")"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TEMPLATE_DIR="${PROJECT_ROOT}/templates"
CACHE_DIR="${PROJECT_ROOT}/.cache/repos"

# --- Parse config ---
parse() { python3 -c "import json,sys; d=json.load(open('$CONFIG')); print(d.get('$1','$2'))"; }

INSTANCE_ID=$(parse instance_id "")
REPO=$(parse repo "")                        # e.g. free5gc/amf
REPO_NAME=$(echo "$REPO" | cut -d/ -f2)      # e.g. amf
NF_TYPE=$(parse nf_type "")
PARENT_COMMIT=$(parse parent_commit "")
FIX_COMMIT=$(parse fix_commit "")
DIFFICULTY=$(parse difficulty "easy")
GO_VERSION=$(parse go_version "1.21")
TEST_PACKAGE_DIR=$(parse test_package_dir "")  # e.g. internal/sbi/processor
TEST_PACKAGE_PATH=$(parse test_package_path "") # e.g. internal/sbi/processor
EXISTING_TEST_PATTERN=$(parse existing_test_pattern "Test")
FAIL_TEST_PATTERN=$(parse fail_test_pattern "Test")
IMAGE_NAME="swebench5g/free5gc:${INSTANCE_ID}"

INSTANCE_DIR="$(cd "$(dirname "$CONFIG")" && pwd)"

echo "============================================"
echo " Building: ${INSTANCE_ID}"
echo " Repo: ${REPO}"
echo " NF: ${NF_TYPE} | Difficulty: ${DIFFICULTY}"
echo " Parent: ${PARENT_COMMIT} | Fix: ${FIX_COMMIT}"
echo "============================================"
echo ""

# --- Step 1: Clone repo on host ---
mkdir -p "$CACHE_DIR"
SRC_CACHE="${CACHE_DIR}/${REPO_NAME}"

if [ -d "${SRC_CACHE}/.git" ]; then
    echo "[1/5] Repo cached at ${SRC_CACHE}, fetching updates..."
    cd "$SRC_CACHE" && git fetch --all && cd "$PROJECT_ROOT"
else
    echo "[1/5] Cloning https://github.com/${REPO}.git ..."
    git clone "https://github.com/${REPO}.git" "$SRC_CACHE"
fi

# --- Step 2: Prepare build context ---
BUILD_DIR="${INSTANCE_DIR}/.build"
rm -rf "$BUILD_DIR"
mkdir -p "${BUILD_DIR}/src" "${BUILD_DIR}/test-suite" "${BUILD_DIR}/task"

echo "[2/5] Copying source to build context..."
# Use git archive to get a clean copy, then init a git repo for checkout
cd "$SRC_CACHE"
git archive HEAD | tar -x -C "${BUILD_DIR}/src"
# We need .git for checkout, so copy it
cp -r "${SRC_CACHE}/.git" "${BUILD_DIR}/src/.git"
cd "$PROJECT_ROOT"

# --- Step 3: Generate Dockerfile from template ---
echo "[3/5] Generating Dockerfile..."
WORKDIR_NAME="free5gc-${REPO_NAME}"

sed -e "s|{{GO_VERSION}}|${GO_VERSION}|g" \
    -e "s|{{INSTANCE_ID}}|${INSTANCE_ID}|g" \
    -e "s|{{NF_TYPE}}|${NF_TYPE}|g" \
    -e "s|{{DIFFICULTY}}|${DIFFICULTY}|g" \
    -e "s|{{REPO}}|${REPO}|g" \
    -e "s|{{PARENT_COMMIT}}|${PARENT_COMMIT}|g" \
    -e "s|{{WORKDIR_NAME}}|${WORKDIR_NAME}|g" \
    "${TEMPLATE_DIR}/Dockerfile.template" > "${BUILD_DIR}/Dockerfile"

# --- Step 4: Copy tests and generate run_tests.sh ---
echo "[4/5] Copying tests and task description..."
cp "${INSTANCE_DIR}/existing_test.go" "${BUILD_DIR}/test-suite/" 2>/dev/null || true
cp "${INSTANCE_DIR}/fail_test.go" "${BUILD_DIR}/test-suite/" 2>/dev/null || true
cp "${INSTANCE_DIR}/problem_statement.md" "${BUILD_DIR}/task/" 2>/dev/null || true

# Generate run_tests.sh from template
sed -e "s|{{INSTANCE_ID}}|${INSTANCE_ID}|g" \
    -e "s|{{WORKDIR_NAME}}|${WORKDIR_NAME}|g" \
    -e "s|{{TEST_PACKAGE_DIR}}|${TEST_PACKAGE_DIR}|g" \
    -e "s|{{TEST_PACKAGE_PATH}}|${TEST_PACKAGE_PATH}|g" \
    -e "s|{{EXISTING_TEST_PATTERN}}|${EXISTING_TEST_PATTERN}|g" \
    -e "s|{{FAIL_TEST_PATTERN}}|${FAIL_TEST_PATTERN}|g" \
    "${TEMPLATE_DIR}/run_tests.sh.template" > "${BUILD_DIR}/test-suite/run_tests.sh"

# --- Step 5: Build Docker image ---
echo "[5/5] Building Docker image: ${IMAGE_NAME}..."
DOCKER_BUILDKIT=0 docker build --network=host -t "${IMAGE_NAME}" "${BUILD_DIR}"

echo ""
echo "============================================"
echo " BUILD SUCCESS: ${IMAGE_NAME}"
echo "============================================"
echo ""
echo "Quick start:"
echo "  docker run -it ${IMAGE_NAME}"
echo ""
echo "Validate:"
echo "  bash scripts/validate_instance.sh ${CONFIG}"
