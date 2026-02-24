#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

log() { printf "\n==> %s\n" "$1"; }

log "Health"
curl -fsS "${BASE_URL}/healthz" >/dev/null
curl -fsS "${BASE_URL}/ping-nats" >/dev/null

log "Auth"
PAYLOAD='{"user_id":"smoke-user","password":"smoke-pass","display_name":"Smoke User"}'
curl -fsS -X POST "${BASE_URL}/auth/register" -H "Content-Type: application/json" -d "$PAYLOAD" >/dev/null || true
COOKIE_JAR=$(mktemp)
curl -fsS -X POST "${BASE_URL}/auth/login" -H "Content-Type: application/json" -d "$PAYLOAD" -c "$COOKIE_JAR" >/dev/null

log "Channels"
curl -fsS "${BASE_URL}/channels" -b "$COOKIE_JAR" >/dev/null
CHAN=$(curl -fsS -X POST "${BASE_URL}/channels" -H "Content-Type: application/json" -d '{"name":"smoke"}' -b "$COOKIE_JAR" 2>/dev/null | python3 -c "import json,sys; print(json.load(sys.stdin).get('id',''))" || true)

log "Publish"
if [ -n "$CHAN" ]; then
  curl -fsS -X POST "${BASE_URL}/channels/${CHAN}/messages" -H "Content-Type: application/json" -d '{"payload":"smoke"}' -b "$COOKIE_JAR" >/dev/null
else
  curl -fsS -X POST "${BASE_URL}/publish?subject=storm.events" -H "Content-Type: application/json" -d '{"msg":"smoke"}' -b "$COOKIE_JAR" >/dev/null
fi

rm -f "$COOKIE_JAR"
log "OK"
