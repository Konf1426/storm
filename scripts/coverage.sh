#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES=("$ROOT_DIR/services/gateway" "$ROOT_DIR/services/messages")
MIN_COVERAGE="${MIN_COVERAGE:-80.0}"
GO_CMD="${GO_CMD:-go}"

if ! command -v "$GO_CMD" >/dev/null 2>&1; then
  echo "go is not installed or not in PATH (set GO_CMD or use scripts/coverage-docker.sh)" >&2
  exit 1
fi

for svc in "${SERVICES[@]}"; do
  echo "==> Coverage: $(basename "$svc")"
  COVER_FILE="$(mktemp)"
  (cd "$svc" && "$GO_CMD" test ./... -coverprofile "$COVER_FILE")
  TOTAL_LINE=$((cd "$svc" && "$GO_CMD" tool cover -func "$COVER_FILE" | tail -n 1) | tee /dev/stderr)
  TOTAL_PCT=$(printf "%s" "$TOTAL_LINE" | awk '{print $3}' | tr -d '%')
  if [ -z "$TOTAL_PCT" ]; then
    echo "[WARN] unable to parse coverage percentage" >&2
  else
    awk -v got="$TOTAL_PCT" -v min="$MIN_COVERAGE" 'BEGIN { exit (got+0 < min+0) }' || {
      echo "[FAIL] coverage ${TOTAL_PCT}% < ${MIN_COVERAGE}% for $(basename "$svc")" >&2
      rm -f "$COVER_FILE"
      exit 1
    }
  fi
  rm -f "$COVER_FILE"
  echo
done
