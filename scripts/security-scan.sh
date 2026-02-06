#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES=("$ROOT_DIR/services/gateway" "$ROOT_DIR/services/messages")

log() { printf "\n==> %s\n" "$1"; }
warn() { printf "[WARN] %s\n" "$1"; }
pass() { printf "[OK] %s\n" "$1"; }

if ! command -v go >/dev/null 2>&1; then
  echo "go is not installed or not in PATH" >&2
  exit 1
fi

HAS_GOVULN=0
if command -v govulncheck >/dev/null 2>&1; then
  HAS_GOVULN=1
else
  warn "govulncheck not found (install: go install golang.org/x/vuln/cmd/govulncheck@latest)"
fi

HAS_GOSEC=0
if command -v gosec >/dev/null 2>&1; then
  HAS_GOSEC=1
else
  warn "gosec not found (install: go install github.com/securego/gosec/v2/cmd/gosec@latest)"
fi

GOSEC_EXCLUDE="${GOSEC_EXCLUDE:-G402}"
GOSEC_STRICT="${GOSEC_STRICT:-1}"
GOSEC_SKIP="${GOSEC_SKIP:-0}"
GOVULN_SKIP="${GOVULN_SKIP:-0}"

for svc in "${SERVICES[@]}"; do
  log "Dependency scan: $(basename "$svc")"
  if [ "$GOVULN_SKIP" -eq 1 ]; then
    warn "govulncheck skipped by GOVULN_SKIP"
  elif [ "$HAS_GOVULN" -eq 1 ]; then
    (cd "$svc" && govulncheck ./...)
    pass "govulncheck complete"
  else
    warn "govulncheck skipped"
  fi

  if [ "$GOSEC_SKIP" -eq 1 ]; then
    warn "gosec skipped by GOSEC_SKIP"
  elif [ "$HAS_GOSEC" -eq 1 ]; then
    if (cd "$svc" && gosec -exclude "$GOSEC_EXCLUDE" ./...); then
      pass "gosec complete"
    else
      warn "gosec found issues"
      if [ "$GOSEC_STRICT" -eq 1 ]; then
        exit 1
      fi
    fi
  else
    warn "gosec skipped"
  fi

done

log "Done"
