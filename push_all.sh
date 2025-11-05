#!/usr/bin/env bash
set -euo pipefail
SHA=$(git rev-parse --short HEAD)

docker build -f Dockerfile --target app -t ghcr.io/axlle-com/go-blog:latest -t ghcr.io/axlle-com/go-blog:app-$SHA .
docker push ghcr.io/axlle-com/go-blog:latest
docker push ghcr.io/axlle-com/go-blog:app-$SHA

docker build -f Dockerfile --target cli -t ghcr.io/axlle-com/go-blog:cli-latest -t ghcr.io/axlle-com/go-blog:cli-$SHA .
docker push ghcr.io/axlle-com/go-blog:cli-latest
docker push ghcr.io/axlle-com/go-blog:cli-$SHA