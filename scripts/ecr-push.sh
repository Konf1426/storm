#!/usr/bin/env bash
set -euo pipefail

REGION="${REGION:-}"
ACCOUNT_ID="${ACCOUNT_ID:-}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

if [ -z "$REGION" ] || [ -z "$ACCOUNT_ID" ]; then
  echo "REGION and ACCOUNT_ID are required" >&2
  exit 1
fi

GATEWAY_REPO="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/storm-gateway"
MESSAGES_REPO="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/storm-messages"

aws ecr get-login-password --region "$REGION" | docker login --username AWS --password-stdin "${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

docker build -t storm-gateway:${IMAGE_TAG} services/gateway
docker build -t storm-messages:${IMAGE_TAG} services/messages

docker tag storm-gateway:${IMAGE_TAG} "${GATEWAY_REPO}:${IMAGE_TAG}"
docker tag storm-messages:${IMAGE_TAG} "${MESSAGES_REPO}:${IMAGE_TAG}"

docker push "${GATEWAY_REPO}:${IMAGE_TAG}"
docker push "${MESSAGES_REPO}:${IMAGE_TAG}"
