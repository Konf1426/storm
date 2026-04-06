#!/usr/bin/env bash
# ──────────────────────────────────────────────────────
# local-to-azure-20k.sh – Run 20k VUs locally against Azure
# ──────────────────────────────────────────────────────
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT="$ROOT_DIR/scripts/k6/storm-azure.js"
IP="40.119.146.38"

echo -e "\n\033[0;36m═══ STORM Local-to-Azure Load Test (20,000 VUs) ═══\033[0m\n"

if ! command -v docker >/dev/null 2>&1; then
  echo "❌ Error: docker is not installed locally. Please install it first." >&2
  exit 1
fi

echo -e "\033[0;33m🔧  Tuning local OS limits (ulimit) for 20k concurrent connections...\033[0m"
# Increase max open files to avoid "too many open files" errors locally
ulimit -n 65535 || echo "⚠️  Could not increase ulimit. Test might fail if limit is too low."

export BASE_URL="http://$IP"
export WS_URL="ws://$IP/ws"
export TARGET_VUS="20000"
export RAMP_UP="3m"
export HOLD="5m"
export RAMP_DOWN="1m"

echo -e "\033[0;33m🚀  Launching k6 with 20,000 VUs via Docker...\033[0m"
echo -e "    Target URL: $BASE_URL"
echo -e "    Target WS : $WS_URL"
echo ""

docker run -i --rm \
  -e BASE_URL="$BASE_URL" \
  -e WS_URL="$WS_URL" \
  -e TARGET_VUS="$TARGET_VUS" \
  -e RAMP_UP="$RAMP_UP" \
  -e HOLD="$HOLD" \
  -e RAMP_DOWN="$RAMP_DOWN" \
  -v "$ROOT_DIR/scripts/k6":/scripts \
  grafana/k6:latest run /scripts/storm-azure.js