#!/usr/bin/env bash
# ──────────────────────────────────────────────────────
# azure-load-test-k8s.sh – Run distributed k6 load tests
# ──────────────────────────────────────────────────────
# Usage:
#   ./scripts/azure-load-test-k8s.sh                # default 1000 VUs per pod (20k total)
#   ./scripts/azure-load-test-k8s.sh 2000 10         # 2000 VUs per pod (20k total)
# ──────────────────────────────────────────────────────

set -euo pipefail

TARGET_VUS=${1:-1000}
PODS=${2:-20}
NAMESPACE="storm"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
JOB_FILE="$ROOT_DIR/infra/k8s-azure/k6-runner/k6-job.yaml"
K6_SCRIPT="$ROOT_DIR/scripts/k6/storm-azure.js"

echo -e "\n\033[0;36m═══ STORM Load Test ($PODS pods × $TARGET_VUS VUs = $(($PODS * $TARGET_VUS)) total) ═══\033[0m\n"

# ─── 1. Check dependencies ──────────────────────────
if ! command -v kubectl >/dev/null 2>&1; then
  echo "❌ Error: kubectl is not installed." >&2
  exit 1
fi

# ─── 2. Update k6 script ConfigMap ───────────────────
echo -e "\033[0;33m📦  Updating k6 script ConfigMap...\033[0m"
kubectl create configmap k6-script \
    --from-file=storm-azure.js="$K6_SCRIPT" \
    -n "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

# ─── 3. Clean previous job if exists ─────────────────
echo -e "\033[0;33m🧹  Cleaning previous job...\033[0m"
kubectl delete job k6-load-test -n "$NAMESPACE" --ignore-not-found

# ─── 4. Launch job ──────────────────────────────────
echo -e "\033[0;33m🚀  Launching k6 runners...\033[0m"

# Create a temp copy with the right values using sed
TEMP_JOB=$(mktemp)
sed -E "s/parallelism: [0-9]+/parallelism: $PODS/" "$JOB_FILE" | \
sed -E "s/completions: [0-9]+/completions: $PODS/" | \
sed -E "/name: TARGET_VUS/{n;s/value: \"[0-9]+\"/value: \"$TARGET_VUS\"/}" > "$TEMP_JOB"

kubectl apply -f "$TEMP_JOB"
rm "$TEMP_JOB"

# ─── 5. Follow logs ─────────────────────────────────
echo -e "\n\033[0;33m📊  Streaming logs (Ctrl+C to stop watching)...\033[0m"
echo -e "    Monitor HPA:  kubectl get hpa -n $NAMESPACE -w"
echo -e "    Monitor pods: kubectl get pods -n $NAMESPACE -w"
echo ""

sleep 5
# Follow logs for all pods with the app=k6-runner label
# We don't use -f because it might hang if pods are still starting, but let's try
kubectl logs -f -l app=k6-runner -n "$NAMESPACE" --max-log-requests="$PODS"
