#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

: "${GATEWAY_URL:=http://localhost:8080}"
: "${SUBJECT:=storm.events}"
: "${DURATION:=60s}"
: "${PUB_VUS:=200}"
: "${WS_VUS:=200}"
: "${PUB_RATE:=200}"

TOKEN="$(bash "$ROOT_DIR/scripts/gen-jwt.sh" test-runner)"

export GATEWAY_URL SUBJECT DURATION PUB_VUS WS_VUS PUB_RATE

docker run --rm --network host \
  -e ACCESS_TOKEN="$TOKEN" \
  -e GATEWAY_URL \
  -e SUBJECT \
  -e DURATION \
  -e PUB_VUS \
  -e WS_VUS \
  -e PUB_RATE \
  -v "$ROOT_DIR:/work" \
  -w /work \
  grafana/k6 run scripts/k6-load.js
