#!/bin/bash
# build_image.sh — Build the Docker image for this pilot task
#
# Step 1: Clone PCF source on the HOST (where network is available)
# Step 2: Build Docker image using COPY instead of git clone inside Docker

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
PCF_SRC_DIR="${PROJECT_DIR}/pcf-src"

IMAGE_NAME="swebench5g/free5gc:pcf_issue_879"

# --- Step 1: Clone PCF repo on the host ---
if [ -d "${PCF_SRC_DIR}/.git" ]; then
    echo "PCF source already exists at ${PCF_SRC_DIR}, skipping clone."
else
    echo "Cloning free5gc/pcf to ${PCF_SRC_DIR} ..."
    git clone https://github.com/free5gc/pcf.git "${PCF_SRC_DIR}"
fi

echo ""

# --- Step 2: Build Docker image ---
echo "Building image: ${IMAGE_NAME}"
echo "Context: ${PROJECT_DIR}"
echo ""

docker build -t "${IMAGE_NAME}" "${PROJECT_DIR}"

echo ""
echo "Build complete: ${IMAGE_NAME}"
echo ""
echo "Quick start:"
echo "  docker run -it ${IMAGE_NAME}"
echo ""
echo "Validate:"
echo "  bash ${SCRIPT_DIR}/validate_image.sh ${IMAGE_NAME}"
