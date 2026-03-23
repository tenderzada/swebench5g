#!/bin/bash
# build_image.sh — Build the Docker image for this pilot task

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

IMAGE_NAME="swebench5g/free5gc:pcf_issue_879"

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
