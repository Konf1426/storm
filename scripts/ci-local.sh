#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

log() { printf "\n==> %s\n" "$1"; }

log "Go tests + coverage"
bash "$ROOT_DIR/scripts/coverage.sh"

log "Security scan (optional)"
if bash "$ROOT_DIR/scripts/security-scan.sh"; then
  echo "Security scan done"
else
  echo "Security scan skipped or failed (tools missing)"
fi

log "Done"
