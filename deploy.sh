#!/usr/bin/env bash
# ============================================================================
#  deploy.sh — выкладка проекта BeautyMed на production-сервер.
#
#  Что делает:
#    1) Синхронизирует исходники (backend, frontend, nginx, docker-compose.yml)
#       на сервер по rsync через SSH, исключая мусор (.git, node_modules,
#       сборочные артефакты, локальные дампы и т.п.).
#    2) Пересобирает и перезапускает контейнеры через docker compose
#       (в режиме --build --remove-orphans -d).
#    3) Проверяет, что все 5 контейнеров поднялись и API отдаёт /health=ok.
#
#  Использование:
#    ./deploy.sh                 — обычная выкладка
#    DRY_RUN=1 ./deploy.sh       — показать, что будет передано, без действий
#    SKIP_BUILD=1 ./deploy.sh    — синхронизировать файлы, но не пересобирать
#    SKIP_VERIFY=1 ./deploy.sh   — без healthcheck
#
#  Зависимости (локально): bash, rsync, sshpass, ssh, curl.
#    macOS:    brew install sshpass hudochenkov/sshpass/sshpass rsync
#    Ubuntu:   sudo apt-get install -y sshpass rsync curl
# ============================================================================

set -euo pipefail

# ───────────── Настройки сервера ─────────────
# SERVER_PASS обязателен и берётся ТОЛЬКО из окружения.
# Удобно хранить в локальном .deploy.env (он в .gitignore):
#   echo 'SERVER_PASS=...' > .deploy.env && source .deploy.env
# либо передавать одной строкой: SERVER_PASS=... ./deploy.sh
SERVER_HOST="${SERVER_HOST:-85.198.80.245}"
SERVER_USER="${SERVER_USER:-root}"
SERVER_PASS="${SERVER_PASS:-}"
REMOTE_DIR="${REMOTE_DIR:-/opt/beautymed}"

if [[ -z "$SERVER_PASS" ]]; then
  echo "ОШИБКА: переменная SERVER_PASS не задана." >&2
  echo "Запустите: SERVER_PASS=... ./deploy.sh   (или экспортируйте перед запуском)" >&2
  exit 1
fi

# ───────────── Цветной вывод ─────────────
if [[ -t 1 ]]; then
  C_RESET=$'\033[0m';  C_RED=$'\033[31m'; C_GRN=$'\033[32m'
  C_YLW=$'\033[33m';   C_BLU=$'\033[34m'; C_BOLD=$'\033[1m'
else
  C_RESET=""; C_RED=""; C_GRN=""; C_YLW=""; C_BLU=""; C_BOLD=""
fi

step()  { echo "${C_BOLD}${C_BLU}==>${C_RESET} ${C_BOLD}$*${C_RESET}"; }
ok()    { echo "  ${C_GRN}✓${C_RESET} $*"; }
warn()  { echo "  ${C_YLW}!${C_RESET} $*"; }
fail()  { echo "  ${C_RED}✗${C_RESET} $*" >&2; }

# ───────────── Проверка зависимостей ─────────────
require() {
  command -v "$1" >/dev/null 2>&1 || {
    fail "Не найдена утилита '$1'. Установите её и повторите попытку."
    exit 1
  }
}

step "Проверяю зависимости"
require rsync
require ssh
require sshpass
require curl
ok "rsync, ssh, sshpass, curl — на месте"

# ───────────── Локальный корень проекта ─────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

if [[ ! -f docker-compose.yml ]]; then
  fail "docker-compose.yml не найден в $SCRIPT_DIR — скрипт должен лежать в корне проекта."
  exit 1
fi

# ───────────── SSH/rsync параметры ─────────────
SSH_OPTS=(-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR)
SSHPASS_CMD=(sshpass -p "$SERVER_PASS")

ssh_run() {
  "${SSHPASS_CMD[@]}" ssh "${SSH_OPTS[@]}" "${SERVER_USER}@${SERVER_HOST}" "$@"
}

# Что НЕ синхронизировать (мусор, секреты, локальные данные сервера).
RSYNC_EXCLUDES=(
  --exclude=.git
  --exclude=.idea
  --exclude=.vscode
  --exclude=.claude
  --exclude=.superpowers
  --exclude=.cursor
  --exclude=node_modules
  --exclude=.nuxt
  --exclude=.output
  --exclude=.nitro
  --exclude=.cache
  --exclude=dist
  --exclude=test-results
  --exclude=playwright-report
  --exclude=coverage
  --exclude='*.log'
  --exclude='*.rdb'
  --exclude=appendonlydir
  --exclude=pgdata
  --exclude=postgres-data
  --exclude=redis-data
  --exclude='.env'
  --exclude='.env.local'
  --exclude='.env.*.local'
  --exclude='.env.production'
  --exclude='.env.development'
  --exclude='*.docx'
  --exclude='*.pdf'
  --exclude='~$*'
  --exclude=EXPLANATION.txt
  --exclude=clinic_1.png
  --exclude=clinic_2.png
  --exclude=clinic_3.png
  --exclude='*.bak'
  --exclude='*.orig'
  --exclude='backend/api'
  --exclude='backend/__debug_bin*'
  --exclude='backend/vendor'
)

RSYNC_OPTS=(-az --delete --stats)
# На macOS системный rsync (BSD-форк, 2.6.9) не понимает --info=stats2 и
# --human-readable — используем универсальный --stats. Для подробного вывода
# в DRY_RUN добавляем -v и --itemize-changes.
[[ "${DRY_RUN:-0}" == "1" ]] && RSYNC_OPTS+=(--dry-run -v --itemize-changes)

# Что передаём (только нужное для прод-сборки).
SYNC_TARGETS=(
  backend
  frontend
  nginx
  docker-compose.yml
)

# ───────────── 1. Синхронизация ─────────────
step "Синхронизирую файлы на ${SERVER_USER}@${SERVER_HOST}:${REMOTE_DIR}"
ssh_run "mkdir -p '${REMOTE_DIR}'"

"${SSHPASS_CMD[@]}" rsync "${RSYNC_OPTS[@]}" "${RSYNC_EXCLUDES[@]}" \
  -e "ssh ${SSH_OPTS[*]}" \
  "${SYNC_TARGETS[@]}" \
  "${SERVER_USER}@${SERVER_HOST}:${REMOTE_DIR}/"

ok "Файлы синхронизированы"

if [[ "${DRY_RUN:-0}" == "1" ]]; then
  warn "DRY_RUN=1 — пересборка и проверка пропущены."
  exit 0
fi

# ───────────── 2. Пересборка контейнеров ─────────────
if [[ "${SKIP_BUILD:-0}" == "1" ]]; then
  warn "SKIP_BUILD=1 — docker compose не запускается."
else
  step "Пересобираю и поднимаю контейнеры (docker compose up -d --build)"
  ssh_run "cd '${REMOTE_DIR}' && docker compose up -d --build --remove-orphans"
  ok "docker compose отработал"
fi

# ───────────── 3. Проверка результата ─────────────
if [[ "${SKIP_VERIFY:-0}" == "1" ]]; then
  warn "SKIP_VERIFY=1 — healthcheck пропущен."
  exit 0
fi

step "Жду, пока контейнеры стабилизируются"
sleep 8

step "Проверяю состояние контейнеров"
PS_OUT="$(ssh_run "cd '${REMOTE_DIR}' && docker compose ps --format '{{.Service}}|{{.State}}|{{.Status}}'")"
echo "${PS_OUT}" | sed 's/^/  /'

EXPECTED=(postgres redis api frontend nginx)
ALL_UP=1
for svc in "${EXPECTED[@]}"; do
  if echo "${PS_OUT}" | awk -F'|' -v s="$svc" '$1==s && $2=="running"' | grep -q .; then
    ok "контейнер '$svc' — running"
  else
    fail "контейнер '$svc' НЕ running"
    ALL_UP=0
  fi
done

step "Пингую HTTP API через nginx"
HEALTH_CODE="$(curl -sk -o /tmp/bm_health.$$ -w '%{http_code}' "https://${SERVER_HOST}/health" || true)"
HEALTH_BODY="$(cat /tmp/bm_health.$$ 2>/dev/null || true)"
rm -f /tmp/bm_health.$$

if [[ "$HEALTH_CODE" != "200" ]]; then
  # Фоллбек на http (на случай если у nginx нет TLS)
  HEALTH_CODE="$(curl -s -o /tmp/bm_health.$$ -w '%{http_code}' "http://${SERVER_HOST}/health" || true)"
  HEALTH_BODY="$(cat /tmp/bm_health.$$ 2>/dev/null || true)"
  rm -f /tmp/bm_health.$$
fi

if [[ "$HEALTH_CODE" == "200" && "$HEALTH_BODY" == *'"ok"'* ]]; then
  ok "/health → 200, body: ${HEALTH_BODY}"
else
  fail "/health отвечает ${HEALTH_CODE}, body: ${HEALTH_BODY:-<пусто>}"
  ALL_UP=0
fi

step "Проверяю один публичный endpoint /api/specialties"
SPEC_CODE="$(curl -sk -o /dev/null -w '%{http_code}' "https://${SERVER_HOST}/api/specialties" || true)"
[[ "$SPEC_CODE" == "200" ]] || SPEC_CODE="$(curl -s -o /dev/null -w '%{http_code}' "http://${SERVER_HOST}/api/specialties" || true)"
if [[ "$SPEC_CODE" == "200" ]]; then
  ok "/api/specialties → 200"
else
  fail "/api/specialties → ${SPEC_CODE}"
  ALL_UP=0
fi

echo
if [[ "$ALL_UP" == "1" ]]; then
  echo "${C_BOLD}${C_GRN}✓ Деплой успешен. Все контейнеры работают, API отвечает.${C_RESET}"
  exit 0
else
  echo "${C_BOLD}${C_RED}✗ Деплой завершён с ошибками. Смотрите вывод выше.${C_RESET}"
  echo "  Подсказка: ${C_BOLD}ssh ${SERVER_USER}@${SERVER_HOST} 'cd ${REMOTE_DIR} && docker compose logs --tail=80'${C_RESET}"
  exit 1
fi
