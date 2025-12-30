#!/usr/bin/env bash
set -euo pipefail

IMAGE="ghcr.io/axlle-com/go-blog"
SHA="$(git rev-parse --short HEAD)"

PLATFORM="linux/amd64"

# Логин (ожидает переменную GHCR_TOKEN в окружении)
# export GHCR_TOKEN=...
echo "${GHCR_TOKEN:?set GHCR_TOKEN}" | docker login ghcr.io -u axlle-com --password-stdin

# APP
docker build --platform "$PLATFORM" -f Dockerfile --target app \
  -t "$IMAGE:latest" \
  -t "$IMAGE:app-$SHA" \
  .

docker push "$IMAGE:latest"
docker push "$IMAGE:app-$SHA"

# CLI
docker build --platform "$PLATFORM" -f Dockerfile --target cli \
  -t "$IMAGE:cli-latest" \
  -t "$IMAGE:cli-$SHA" \
  .

docker push "$IMAGE:cli-latest"
docker push "$IMAGE:cli-$SHA"
