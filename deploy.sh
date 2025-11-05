#!/usr/bin/env bash
set -euo pipefail

# --- defaults ---
STACK="${STACK:-docker-compose.yml}"      # путь к compose-файлу
SERVICES="${SERVICES:-app cli}"        # какие сервисы обновлять (через пробел)
PROFILE="${PROFILE:-}"             # профиль(-и) compose, например: PROFILE="cli"
GHCR_USER="${GHCR_USER:-}"         # логин GHCR (если образы приватные)
GHCR_TOKEN="${GHCR_TOKEN:-}"       # PAT с read:packages (если приватные)
PRUNE="${PRUNE:-1}"                # 1 = docker image prune -f после апдейта
LOGIN="${LOGIN:-auto}"             # auto|yes|no — логин в GHCR на этой машине
DRY_RUN="${DRY_RUN:-0}"            # 1 = только показать команды

# --- detect docker compose ---
if command -v docker &>/dev/null && docker compose version &>/dev/null; then
  DC="docker compose"
elif command -v docker-compose &>/dev/null; then
  DC="docker-compose"
else
  echo "docker compose не найден"; exit 1
fi

usage() {
  cat <<EOF
Usage:
  STACK=./docker-compose.yml SERVICES="app" ./deploy.sh

Env options:
  STACK=<path>           - путь к docker compose файлу (по умолчанию: docker-compose.yml)
  SERVICES="app cli"     - список сервисов для pull/up (по умолчанию: app)
  PROFILE="cli"          - профиль compose (опционально)
  GHCR_USER=<user>       - GHCR юзер (для приватных образов)
  GHCR_TOKEN=<token>     - GHCR токен (read:packages)
  LOGIN=auto|yes|no      - логин в GHCR (auto: если заданы креды)
  PRUNE=1|0              - чистить dangling образы после апдейта (по умолчанию 1)
  DRY_RUN=1              - показать команды, но не выполнять

Примеры:
  ./deploy.sh
  STACK=/opt/go-app/compose.yml ./deploy.sh
  SERVICES="app" GHCR_USER=axlle-com GHCR_TOKEN=*** ./deploy.sh
  PROFILE=cli SERVICES="cli" ./deploy.sh
EOF
}

# help по ключам
if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage; exit 0
fi

run() {
  echo "+ $*"
  if [[ "$DRY_RUN" != "1" ]]; then
    eval "$@"
  fi
}

# --- GHCR login (если нужно) ---
if [[ "$LOGIN" == "yes" || ( "$LOGIN" == "auto" && -n "$GHCR_USER" && -n "$GHCR_TOKEN" ) ]]; then
  run "echo '$GHCR_TOKEN' | docker login ghcr.io -u '$GHCR_USER' --password-stdin"
else
  echo "GHCR login пропущен (LOGIN=$LOGIN)"
fi

# --- build compose cmd ---
PROFILE_ARG=""
if [[ -n "$PROFILE" ]]; then
  PROFILE_ARG="--profile $PROFILE"
fi

# --- pull ---
for S in $SERVICES; do
  run "$DC -f '$STACK' $PROFILE_ARG pull '$S'"
done

# --- up ---
for S in $SERVICES; do
  run "$DC -f '$STACK' $PROFILE_ARG up -d '$S'"
done

# --- prune dangling images ---
if [[ "$PRUNE" == "1" ]]; then
  run "docker image prune -f"
fi

echo "==> DONE"
