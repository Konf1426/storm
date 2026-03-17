#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SONAR_HOST_URL="${SONAR_HOST_URL:-http://localhost:9000}"
SONAR_PROJECT_KEY="${SONAR_PROJECT_KEY:-storm}"

SCANNER_SONAR_HOST_URL="$SONAR_HOST_URL"
if [[ "$SCANNER_SONAR_HOST_URL" == http://localhost:* ]] || [[ "$SCANNER_SONAR_HOST_URL" == https://localhost:* ]]; then
  SCANNER_SONAR_HOST_URL="${SCANNER_SONAR_HOST_URL/localhost/host.docker.internal}"
fi

if [ -z "${SONAR_TOKEN:-}" ]; then
  echo "SONAR_TOKEN is required" >&2
  exit 1
fi

cd "$ROOT_DIR"
chmod +x scripts/sonar-coverage.sh
./scripts/sonar-coverage.sh

docker run --rm \
  --add-host=host.docker.internal:host-gateway \
  -e SONAR_HOST_URL="$SCANNER_SONAR_HOST_URL" \
  -e SONAR_TOKEN="$SONAR_TOKEN" \
  -v "$ROOT_DIR:/usr/src" \
  sonarsource/sonar-scanner-cli:latest \
  -Dsonar.projectKey="$SONAR_PROJECT_KEY"
