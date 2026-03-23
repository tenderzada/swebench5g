#!/bin/bash
# apply_fix.sh — Apply the official fix (PR #65) to verify tests pass
# This script is for validation only, NOT for agents to use.

set -e

PCF_DIR="/opt/free5gc-pcf"

echo "Applying fix from PR #65 (commit 508d70b)..."
cd "${PCF_DIR}"

# Cherry-pick the fix commits (use ghproxy mirror as fallback)
git remote set-url origin https://ghproxy.com/https://github.com/free5gc/pcf.git
git fetch origin || {
    git remote set-url origin https://github.com/free5gc/pcf.git
    git fetch origin
}
git cherry-pick --no-commit 581143c78a532506f89d7c9e51f8e5078dc23f56
git cherry-pick --no-commit 7c472a78542bddfc3afb1195207db618d927163e

echo "Fix applied. Run '/opt/test-suite/run_tests.sh all' to verify."
