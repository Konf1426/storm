#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
K8S_DIR="$ROOT_DIR/infra/k8s"
NAMESPACE="${NAMESPACE:-storm}"
BUILD_IMAGES="${BUILD_IMAGES:-1}"
GATEWAY_IMAGE="${GATEWAY_IMAGE:-storm-gateway:latest}"
MESSAGES_IMAGE="${MESSAGES_IMAGE:-storm-messages:latest}"

log() { printf "\n==> %s\n" "$1"; }

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl is not installed or not in PATH" >&2
  exit 1
fi

if [ "$BUILD_IMAGES" = "1" ]; then
  if ! command -v docker >/dev/null 2>&1; then
    echo "docker is not installed or not in PATH" >&2
    exit 1
  fi
  log "Build local images"
  docker build -t "$GATEWAY_IMAGE" "$ROOT_DIR/services/gateway"
  docker build -t "$MESSAGES_IMAGE" "$ROOT_DIR/services/messages"
fi

log "Apply manifests"
kubectl apply -k "$K8S_DIR"

log "Wait for core deployments"
kubectl -n "$NAMESPACE" rollout status deploy/nats --timeout=180s
kubectl -n "$NAMESPACE" rollout status deploy/postgres --timeout=180s
kubectl -n "$NAMESPACE" rollout status deploy/redis --timeout=180s
kubectl -n "$NAMESPACE" rollout status deploy/gateway --timeout=180s
kubectl -n "$NAMESPACE" rollout status deploy/messages --timeout=180s

log "Pods"
kubectl -n "$NAMESPACE" get pods -o wide

log "Services"
kubectl -n "$NAMESPACE" get svc

cat <<EOF

Kubernetes deployment done.

Useful commands:
  kubectl -n $NAMESPACE port-forward svc/gateway 8080:8080
  kubectl -n $NAMESPACE port-forward svc/grafana 3000:3000
  kubectl -n $NAMESPACE logs deploy/gateway -f

EOF
