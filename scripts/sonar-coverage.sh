#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="$ROOT_DIR/out/sonar"

mkdir -p "$OUT_DIR"

run_go_coverage() {
  local service_dir="$1"
  local cover_file="$2"
  echo "==> Sonar coverage: $(basename "$service_dir")"
  (
    cd "$service_dir"
    go test ./... -coverprofile "$cover_file"
  )
}

run_go_coverage "$ROOT_DIR/services/gateway" "$OUT_DIR/gateway.coverage.out"
run_go_coverage "$ROOT_DIR/services/messages" "$OUT_DIR/messages.coverage.out"
