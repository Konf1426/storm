#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/infra/docker/docker-compose.yml"

log() { printf "\n==> %s\n" "$1"; }
warn() { printf "[WARN] %s\n" "$1"; }

svc_id() {
  docker compose -f "$COMPOSE_FILE" ps -q "$1"
}

require_running() {
  local id
  id="$(svc_id "$1")"
  if [ -z "$id" ]; then
    warn "service $1 is not running"
    return 1
  fi
  echo "$id"
}

pause_service() {
  local svc="$1" seconds="${2:-10}"
  local id
  id="$(require_running "$svc")" || return 0
  log "Pause $svc for ${seconds}s"
  docker pause "$id" >/dev/null
  sleep "$seconds"
  docker unpause "$id" >/dev/null
}

restart_service() {
  local svc="$1" seconds="${2:-5}"
  local id
  id="$(require_running "$svc")" || return 0
  log "Restart $svc"
  docker stop "$id" >/dev/null
  sleep "$seconds"
  docker start "$id" >/dev/null
}

kill_service() {
  local svc="$1" seconds="${2:-5}"
  local id
  id="$(require_running "$svc")" || return 0
  log "Kill $svc"
  docker kill "$id" >/dev/null
  sleep "$seconds"
  docker start "$id" >/dev/null
}

add_latency() {
  local svc="$1" delay="${2:-100ms}" jitter="${3:-20ms}" seconds="${4:-10}"
  local id
  id="$(require_running "$svc")" || return 0
  log "Add latency to $svc (${delay} +/- ${jitter}) for ${seconds}s"
  if ! docker exec "$id" sh -c "command -v tc >/dev/null 2>&1"; then
    warn "tc not found in $svc container; installing iproute2 (ephemeral)"
    docker exec "$id" sh -c "apk add --no-cache iproute2 >/dev/null 2>&1" || {
      warn "failed to install iproute2 in $svc; skipping latency"
      return 0
    }
  fi
  docker exec "$id" sh -c "tc qdisc add dev eth0 root netem delay $delay $jitter" || {
    warn "failed to add latency on $svc"
    return 0
  }
  sleep "$seconds"
  docker exec "$id" sh -c "tc qdisc del dev eth0 root netem" || true
}

cpu_spike() {
  local svc="$1" seconds="${2:-10}"
  local id
  id="$(require_running "$svc")" || return 0
  log "CPU spike on $svc for ${seconds}s"
  docker exec "$id" sh -c "sh -c 'dd if=/dev/zero of=/dev/null' & echo \$!" >/tmp/chaos_pid.txt
  local pid
  pid="$(cat /tmp/chaos_pid.txt)"
  sleep "$seconds"
  docker exec "$id" sh -c "kill $pid >/dev/null 2>&1" || true
  rm -f /tmp/chaos_pid.txt
}

usage() {
  cat <<EOF
Usage: bash scripts/chaos.sh <scenario>

Scenarios:
  gateway-restart      Restart gateway (5s)
  messages-restart     Restart messages (5s)
  nats-pause           Pause NATS (10s)
  gateway-latency      Add 100ms latency to gateway (10s)
  messages-latency     Add 100ms latency to messages (10s)
  gateway-cpu          CPU spike gateway (10s)
  full                Run a short chaos sequence
EOF
}

scenario="${1:-}"
case "$scenario" in
  gateway-restart) restart_service gateway 5 ;;
  messages-restart) restart_service messages 5 ;;
  nats-pause) pause_service nats 10 ;;
  gateway-latency) add_latency gateway 100ms 20ms 10 ;;
  messages-latency) add_latency messages 100ms 20ms 10 ;;
  gateway-cpu) cpu_spike gateway 10 ;;
  full)
    pause_service nats 10
    restart_service gateway 5
    add_latency gateway 100ms 20ms 10
    cpu_spike gateway 10
    restart_service messages 5
    ;;
  *) usage; exit 1 ;;
esac
