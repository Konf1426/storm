#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${NAMESPACE:-storm}"

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl is not installed or not in PATH" >&2
  exit 1
fi

kubectl -n "$NAMESPACE" port-forward svc/gateway 8080:8080 &
GW_PID=$!
kubectl -n "$NAMESPACE" port-forward svc/grafana 3000:3000 &
GF_PID=$!
kubectl -n "$NAMESPACE" port-forward svc/prometheus 9090:9090 &
PR_PID=$!

cleanup() {
  kill "$GW_PID" "$GF_PID" "$PR_PID" >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

echo "Port-forward active:"
echo "  gateway    http://localhost:8080"
echo "  grafana    http://localhost:3000"
echo "  prometheus http://localhost:9090"
echo "Press Ctrl+C to stop."

wait
