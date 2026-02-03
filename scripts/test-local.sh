#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker/docker-compose.yml"

log() { printf "\n==> %s\n" "$1"; }
pass() { printf "[OK] %s\n" "$1"; }
warn() { printf "[WARN] %s\n" "$1"; }
note() { printf "[INFO] %s\n" "$1"; }

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
"${COMPOSE[@]}" -f "$COMPOSE_FILE" up -d --build

log "Waiting for gateway"
for i in {1..30}; do
  if curl -fsS http://localhost:8080/healthz >/dev/null; then
    break
  fi
  sleep 1
  if [ "$i" -eq 30 ]; then
    echo "gateway not responding on :8080" >&2
    exit 1
  fi
done

log "Register/login"
AUTH_PAYLOAD='{"user_id":"test-user","password":"test-pass","display_name":"Test User"}'
if ! curl -fsS -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d "$AUTH_PAYLOAD" >/dev/null; then
  note "register skipped (user may already exist)"
fi

COOKIE_JAR=$(mktemp)
LOGIN_RES=$(curl -fsS -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d "$AUTH_PAYLOAD" -c "$COOKIE_JAR")
pass "login ok"

log "Health checks"
curl -fsS http://localhost:8080/healthz && echo
curl -fsS http://localhost:8080/ping-nats && echo
pass "Gateway health + NATS OK"

log "Channels"
CHANNEL_RES=$(curl -fsS http://localhost:8080/channels -b "$COOKIE_JAR")
pass "list channels"

log "Create channel"
CREATE_RES=$(curl -fsS -X POST http://localhost:8080/channels -H "Content-Type: application/json" -d '{"name":"general"}' -b "$COOKIE_JAR")
CHANNEL_ID=$(echo "$CREATE_RES" | python3 -c "import json,sys; print(json.load(sys.stdin).get('id',''))" 2>/dev/null || true)
if [ -z "$CHANNEL_ID" ]; then
  warn "channel create failed (maybe exists)"
  CHANNEL_ID=$(echo "$CHANNEL_RES" | python3 -c "import json,sys; data=json.load(sys.stdin); print(data[0]['id'] if data else '')" 2>/dev/null || true)
fi

log "Publish test message"
if [ -n "$CHANNEL_ID" ]; then
  curl -fsS -X POST "http://localhost:8080/channels/${CHANNEL_ID}/messages" \
    -H "Content-Type: application/json" \
    -d '{"payload":"hello from test script"}' \
    -b "$COOKIE_JAR" >/dev/null
  pass "channel message published"
else
  curl -fsS -X POST "http://localhost:8080/publish?subject=storm.events" \
    -H "Content-Type: application/json" \
    -d '{"type":"test","msg":"hello from test script"}' \
    -b "$COOKIE_JAR" >/dev/null
  pass "subject message published"
fi

log "SSE stream check"
SSE_OUT=$(mktemp)
curl -fsS --no-buffer --max-time 4 "http://localhost:8080/events?subject=storm.events" \
  -b "$COOKIE_JAR" >"$SSE_OUT" 2>/dev/null &
CURL_PID=$!
sleep 1
curl -fsS -X POST "http://localhost:8080/publish?subject=storm.events" \
  -H "Content-Type: application/json" \
  -d '{"type":"sse","msg":"stream check"}' \
  -b "$COOKIE_JAR" >/dev/null
wait "$CURL_PID" || true
if grep -q "data:" "$SSE_OUT"; then
  pass "SSE stream received data"
else
  warn "SSE stream did not receive data (check /events endpoint)"
fi
rm -f "$SSE_OUT" "$COOKIE_JAR"

log "Prometheus targets"
if curl -fsS http://localhost:9090/-/ready >/dev/null; then
  echo "Prometheus is ready: http://localhost:9090"
else
  echo "Prometheus not ready yet (check logs if needed)"
fi

log "Frontend dev server (auto)"
if curl -fsS http://localhost:5173 >/dev/null; then
  pass "Frontend reachable on http://localhost:5173"
elif command -v npm >/dev/null 2>&1; then
  if [ -d "$ROOT_DIR/frontend/node_modules" ]; then
    FRONT_LOG=$(mktemp)
    (cd "$ROOT_DIR/frontend" && npm run dev -- --host 127.0.0.1 --port 5173 >"$FRONT_LOG" 2>&1 &) 
    FRONT_PID=$!
    READY=0
    for i in {1..15}; do
      if curl -fsS http://localhost:5173 >/dev/null; then
        READY=1
        break
      fi
      sleep 1
    done
    if [ "$READY" -eq 1 ]; then
      pass "Frontend started and reachable on http://localhost:5173"
    else
      warn "Frontend did not start (check frontend logs)"
    fi
    kill "$FRONT_PID" >/dev/null 2>&1 || true
  else
    warn "Frontend deps missing (run: cd frontend && npm install)"
  fi
else
  warn "npm not found (install Node.js to run frontend)"
fi

log "Done"