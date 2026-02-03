#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES=("$ROOT_DIR/services/gateway" "$ROOT_DIR/services/messages")

for svc in "${SERVICES[@]}"; do
  echo "==> Coverage: $(basename "$svc")"
  COVER_FILE="$(mktemp)"
  (cd "$svc" && go test ./... -coverprofile "$COVER_FILE")
  (cd "$svc" && go tool cover -func "$COVER_FILE" | tail -n 1)
  rm -f "$COVER_FILE"
  echo
done
