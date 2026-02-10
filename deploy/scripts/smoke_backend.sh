#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "[smoke] health check"
curl -fsS "${BASE_URL}/healthz" >/dev/null

echo "[smoke] dashboard"
curl -fsS "${BASE_URL}/api/v1/dashboard" >/dev/null

echo "[smoke] vm power"
curl -fsS -X POST "${BASE_URL}/api/v1/vms/vm-1/power" \
  -H 'Content-Type: application/json' \
  -d '{"action":"stop"}' >/dev/null

echo "[smoke] vm snapshot"
curl -fsS -X POST "${BASE_URL}/api/v1/vms/vm-1/snapshot" >/dev/null

echo "[smoke] vm migrate"
curl -fsS -X POST "${BASE_URL}/api/v1/vms/vm-1/migrate" \
  -H 'Content-Type: application/json' \
  -d '{"host":"kvm-host-03"}' >/dev/null

echo "[smoke] console ticket"
curl -fsS -X POST "${BASE_URL}/api/v1/vms/vm-1/console-ticket" >/dev/null

echo "[smoke] volume expand"
curl -fsS -X POST "${BASE_URL}/api/v1/storage/volumes/vol-2/expand" \
  -H 'Content-Type: application/json' \
  -d '{"sizeGb":900}' >/dev/null

echo "[smoke] host maintenance"
curl -fsS -X POST "${BASE_URL}/api/v1/hosts/host-1/maintenance" \
  -H 'Content-Type: application/json' \
  -d '{"enabled":true}' >/dev/null

echo "[smoke] all checks passed"
