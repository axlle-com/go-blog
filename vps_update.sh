#!/usr/bin/env bash
set -euo pipefail

# ===== Подключение =====
HOST="${HOST:-31.128.42.230}"
PORT="${PORT:-22}"
REMOTE_USER="${REMOTE_USER:-root}"
KEY="${KEY:-$HOME/.ssh/id_rsa}"

SSH_OPTS="-i $KEY -o IdentitiesOnly=yes -o PreferredAuthentications=publickey \
  -o PubkeyAuthentication=yes -o StrictHostKeyChecking=accept-new"

# ===== Проект на VPS =====
APP_DIR="${APP_DIR:-/opt/go-app}"     # путь к git-репозиторию на VPS
BRANCH="${BRANCH:-master}"            # ветка
REPO_URL="${REPO_URL:-git@github.com:axlle-com/go-blog.git}"

# ===== Команда после апдейта =====
MAKE_CMD="${MAKE_CMD:-make rebuild-dev}"

echo "==> SSH: $REMOTE_USER@$HOST:$PORT"
echo "==> Repo dir: $APP_DIR | branch: $BRANCH"
[ -n "$REPO_URL" ] && echo "==> origin: $REPO_URL"
echo "==> Command: $MAKE_CMD"

ssh $SSH_OPTS -p "$PORT" "$REMOTE_USER@$HOST" 'bash -s' <<EOF
set -euo pipefail

# передаём значения надёжно
APP_DIR='$APP_DIR'
BRANCH='$BRANCH'
REPO_URL='$REPO_URL'
MAKE_CMD='$MAKE_CMD'

log(){ printf "\033[1;32m==> %s\033[0m\n" "\$*"; }
warn(){ printf "\033[1;33m[!] %s\033[0m\n" "\$*"; }
err(){ printf "\033[1;31m[✗] %s\033[0m\n" "\$*"; }
need(){ command -v "\$1" >/dev/null; }

: "\${APP_DIR:?}"; : "\${BRANCH:?}"; : "\${MAKE_CMD:?}"

# пакетный менеджер
PM=""
if need apt-get; then PM="apt"; elif need dnf; then PM="dnf"; elif need yum; then PM="yum"; elif need apk; then PM="apk"; fi
install(){
  case "\$PM" in
    apt) apt-get update -y && apt-get install -y "\$@";;
    dnf|yum) "\$PM" -y install "\$@";;
    apk) apk add --no-cache "\$@";;
    *) err "Не удалось определить пакетный менеджер"; exit 1;;
  esac
}

need git  || install git
need make || install make
need docker || { err "docker не установлен"; exit 1; }
docker compose version >/dev/null || { err "docker compose v2 не найден"; exit 1; }

# проверка репозитория
[ -d "\$APP_DIR/.git" ] || { err "В \$APP_DIR нет .git"; exit 1; }
cd "\$APP_DIR"

# при необходимости переустановим origin
if [ -n "\$REPO_URL" ]; then
  log "git remote set-url origin \$REPO_URL"
  git remote set-url origin "\$REPO_URL"
fi

log "git fetch --tags --force --prune"
git fetch --tags --force --prune

# жёстко обновляем ветку до origin/\$BRANCH
if git help switch >/dev/null 2>&1; then
  log "git switch -C \$BRANCH origin/\$BRANCH"
  git switch -C "\$BRANCH" "origin/\$BRANCH"
else
  log "git checkout -B \$BRANCH origin/\$BRANCH"
  git checkout -B "\$BRANCH" "origin/\$BRANCH"
fi

log "git reset --hard origin/\$BRANCH"
git reset --hard "origin/\$BRANCH"

log "git clean -xdf -e .env -e acme.json -e data -e data/**"
git clean -xdf -e .env -e acme.json -e data -e data/**

# сабмодули (если есть)
if [ -f .gitmodules ]; then
  log "git submodule sync/update"
  git submodule sync --recursive
  git submodule update --init --recursive --force
fi

git gc --prune=now --aggressive >/dev/null 2>&1 || true

CUR=\$(git rev-parse --short HEAD || true)
log "HEAD: \$CUR"
[ -f Makefile ] && log "Makefile найден" || warn "Makefile не найден"

# запуск указанной команды
if [[ "\$MAKE_CMD" =~ ^[[:space:]]*make[[:space:]]+(.+)\$ ]]; then
  TARGETS="\${BASH_REMATCH[1]}"
  log "run: make -C \$APP_DIR \$TARGETS"
  make -C "\$APP_DIR" \$TARGETS
else
  log "run: \$MAKE_CMD (cwd=\$APP_DIR)"
  bash -lc "\$MAKE_CMD"
fi

log "ГОТОВО"
EOF
