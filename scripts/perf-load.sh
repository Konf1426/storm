#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT="$ROOT_DIR/scripts/k6-load.js"

if ! command -v k6 >/dev/null 2>&1; then
  echo "k6 is not installed. Install: https://k6.io/docs/get-started/installation/" >&2
  exit 1
fi

: "${GATEWAY_URL:=http://localhost:8080}"
: "${SUBJECT:=storm.events}"
: "${DURATION:=30s}"
: "${PUB_VUS:=5}"
: "${WS_VUS:=5}"
: "${PUB_RATE:=10}"
: "${MODE:=baseline}"

if [ "$MODE" = "chaos" ]; then
  : "${THRESH_FAIL:=0.10}"
  : "${THRESH_P95:=1500}"
else
  : "${THRESH_FAIL:=0.02}"
  : "${THRESH_P95:=500}"
fi

if [ -z "${ACCESS_TOKEN:-}" ]; then
  echo "ACCESS_TOKEN is empty. Generate one with: bash scripts/gen-jwt.sh user-1" >&2
  echo "Then run: ACCESS_TOKEN=... bash scripts/perf-load.sh" >&2
  exit 1
fi

export GATEWAY_URL SUBJECT DURATION PUB_VUS WS_VUS PUB_RATE ACCESS_TOKEN CHANNEL_ID THRESH_FAIL THRESH_P95

k6 run "$SCRIPT"
