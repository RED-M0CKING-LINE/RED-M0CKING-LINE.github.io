#!/usr/bin/env bash
set -euo pipefail

# Meant to be run with crontab, add the following:
#PROJECT_ROOT=/home/local-admin/ethanashley-net-website-main-go
#*/10 * * * * systemd-cat -t personal-website-main-go-deploy ${WEB_PROJECT_ROOT}/deploy/cron-deploy.sh

# Check logs with: journalctl -t personal-website-main-go-deploy -f

: "${WEB_PROJECT_ROOT:?WEB_PROJECT_ROOT must be set}"
COMPOSE_FILE="${PROJECT_ROOT}/deploy/compose/prod.compose.yaml"
BRANCH="prod"

log() { echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"; }

cd "${WEB_PROJECT_ROOT}"

git fetch origin "${BRANCH}"

LOCAL=$(git rev-parse HEAD)
REMOTE=$(git rev-parse "origin/${BRANCH}")

if [ "${LOCAL}" = "${REMOTE}" ]; then
    log "No changes (${LOCAL:0:8})..."
    exit 0
fi

log "New commits detected: ${LOCAL:0:8} -> ${REMOTE:0:8}"
git pull origin "${BRANCH}"

log "Building new stack..."
podman compose -f "${COMPOSE_FILE}" build

log "Bringing stack down..."
podman compose -f "${COMPOSE_FILE}" down

log "Starting stack..."
podman compose -f "${COMPOSE_FILE}" up -d

log "Deploy complete!"
