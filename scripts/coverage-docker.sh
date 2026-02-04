#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES=("$ROOT_DIR/services/gateway" "$ROOT_DIR/services/messages")
MIN_COVERAGE="${MIN_COVERAGE:-80.0}"
GO_IMAGE="${GO_IMAGE:-golang:1.24-alpine}"

log() { printf "\n==> %s\n" "$1"; }

for svc in "${SERVICES[@]}"; do
  log "Coverage (docker): $(basename "$svc")"
  docker run --rm \
    -v "$svc:/app" \
    -w /app \
    "$GO_IMAGE" sh -lc \
    "/usr/local/go/bin/go test ./... -coverprofile /tmp/cover.out >/tmp/test.out && \
     /usr/local/go/bin/go tool cover -func /tmp/cover.out | tail -n 1" | tee /tmp/cover_line.txt

  TOTAL_PCT=$(awk '{print $3}' /tmp/cover_line.txt | tr -d '%')
  if [ -z "$TOTAL_PCT" ]; then
    echo "[WARN] unable to parse coverage percentage" >&2
  else
    awk -v got="$TOTAL_PCT" -v min="$MIN_COVERAGE" 'BEGIN { exit (got+0 < min+0) }' || {
      echo "[FAIL] coverage ${TOTAL_PCT}% < ${MIN_COVERAGE}% for $(basename "$svc")" >&2
      exit 1
    }
  fi
  rm -f /tmp/cover_line.txt
done
