#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker/docker-compose.yml"
OUT_DIR="${OUT_DIR:-$ROOT_DIR/out/storm-day}"
RUN_ID="${RUN_ID:-$(date +%Y%m%d-%H%M%S)}"
RUN_DIR="$OUT_DIR/$RUN_ID"

log() { printf "\n==> %s\n" "$1"; }
warn() { printf "[WARN] %s\n" "$1"; }

mkdir -p "$RUN_DIR"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is not installed or not in PATH" >&2
  exit 1
fi

if docker compose version >/dev/null 2>&1; then
  COMPOSE=(docker compose)
elif command -v docker-compose >/dev/null 2>&1; then
  COMPOSE=(docker-compose)
else
  echo "docker compose not available" >&2
  exit 1
fi

log "Starting stack"
"${COMPOSE[@]}" -f "$COMPOSE_FILE" up -d --build | tee "$RUN_DIR/compose-up.log"

log "Waiting for gateway"
for i in {1..60}; do
  if curl -fsS http://localhost:8080/healthz >/dev/null; then
    break
  fi
  sleep 1
  if [ "$i" -eq 60 ]; then
    echo "gateway not responding on :8080" >&2
    exit 1
  fi
done

if ! command -v k6 >/dev/null 2>&1; then
  warn "k6 not installed; load tests will be skipped"
  HAS_K6=0
else
  HAS_K6=1
fi

log "Generating access token"
ACCESS_TOKEN="$(JWT_SECRET=dev-secret bash "$ROOT_DIR/scripts/gen-jwt.sh" storm-day-user)"

export ACCESS_TOKEN
export GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"
export SUBJECT="${SUBJECT:-storm.events}"

log "Phase 1 - Warm-up"
if [ "$HAS_K6" -eq 1 ]; then
  DURATION=2m PUB_VUS=50 WS_VUS=50 PUB_RATE=20 MODE=baseline bash "$ROOT_DIR/scripts/perf-load.sh" | tee "$RUN_DIR/warmup.log"
else
  warn "Skipping k6 warm-up"
fi

log "Phase 2 - Growth spike"
if [ "$HAS_K6" -eq 1 ]; then
  DURATION=3m PUB_VUS=150 WS_VUS=150 PUB_RATE=50 MODE=baseline bash "$ROOT_DIR/scripts/perf-load.sh" | tee "$RUN_DIR/spike1.log"
else
  warn "Skipping k6 spike"
fi

log "Phase 3 - Incident 1 (NATS pause)"
bash "$ROOT_DIR/scripts/chaos.sh" nats-pause | tee "$RUN_DIR/incident1.log"

log "Phase 4 - Recovery"
sleep 30

log "Phase 5 - Incident 2 (Gateway latency)"
bash "$ROOT_DIR/scripts/chaos.sh" gateway-latency | tee "$RUN_DIR/incident2.log"

log "Phase 6 - Growth spike 2"
if [ "$HAS_K6" -eq 1 ]; then
  DURATION=3m PUB_VUS=200 WS_VUS=200 PUB_RATE=60 MODE=chaos bash "$ROOT_DIR/scripts/perf-load.sh" | tee "$RUN_DIR/spike2.log"
else
  warn "Skipping k6 spike 2"
fi

log "Capturing pprof snapshots (best effort)"
if curl -fsS "http://localhost:6060/debug/pprof/profile?seconds=10" -o "$RUN_DIR/gateway-cpu.pprof"; then
  echo "gateway cpu profile saved"
else
  warn "gateway cpu profile not captured"
fi
if curl -fsS "http://localhost:6061/debug/pprof/profile?seconds=10" -o "$RUN_DIR/messages-cpu.pprof"; then
  echo "messages cpu profile saved"
else
  warn "messages cpu profile not captured"
fi

log "Done"
echo "Artifacts: $RUN_DIR"
