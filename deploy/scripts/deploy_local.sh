#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
LOG_FILE="${ROOT_DIR}/.local-backend.log"
PID_FILE="${ROOT_DIR}/.local-backend.pid"
BIN_FILE="${ROOT_DIR}/backend/bin/kvm-api"

is_valid_backend_pid() {
  local pid="$1"
  kill -0 "$pid" 2>/dev/null || return 1
  ps -p "$pid" -o args= 2>/dev/null | grep -q "${BIN_FILE}"
}

if [[ -f "${PID_FILE}" ]]; then
  old_pid="$(cat "${PID_FILE}")"
  if is_valid_backend_pid "$old_pid"; then
    echo "backend already running with pid ${old_pid}"
    exit 0
  fi
  echo "[deploy] removing stale pid file"
  : >"${PID_FILE}"
fi

echo "[deploy] building backend binary"
mkdir -p "${ROOT_DIR}/backend/bin"
(cd "${ROOT_DIR}/backend" && go build -o "${BIN_FILE}" ./cmd/api)

echo "[deploy] starting backend on :8080"
nohup "${BIN_FILE}" >"${LOG_FILE}" 2>&1 &
new_pid=$!
echo "$new_pid" >"${PID_FILE}"

sleep 1

echo "[deploy] running smoke test"
BASE_URL="http://localhost:8080" bash "${ROOT_DIR}/deploy/scripts/smoke_backend.sh"

echo "[deploy] success"
echo "- pid file: ${PID_FILE}"
echo "- logs: ${LOG_FILE}"
echo "- stop command: kill ${new_pid} && rm -f ${PID_FILE}"
