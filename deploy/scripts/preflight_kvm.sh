#!/usr/bin/env bash
set -euo pipefail

ALLOW_MISSING="${ALLOW_MISSING:-0}"
MODE="${MODE:-host}" # host|dashboard
failures=0

check_cmd() {
  local cmd="$1"
  if command -v "$cmd" >/dev/null 2>&1; then
    echo "[ok] command found: $cmd"
  else
    echo "[err] missing command: $cmd"
    failures=$((failures + 1))
  fi
}

check_file() {
  local path="$1"
  if [[ -e "$path" ]]; then
    echo "[ok] path exists: $path"
  else
    echo "[err] missing path: $path"
    failures=$((failures + 1))
  fi
}

if [[ "$MODE" == "host" ]]; then
  echo "[preflight] mode=host; checking virtualization dependencies"
  check_cmd qemu-system-x86_64
  check_cmd virsh
  check_cmd virt-host-validate
  check_cmd bridge
  check_file /dev/kvm
  check_file /etc/libvirt/libvirtd.conf

  if command -v virt-host-validate >/dev/null 2>&1; then
    echo "[preflight] running virt-host-validate"
    if ! virt-host-validate qemu >/tmp/virt-host-validate.log 2>&1; then
      echo "[warn] virt-host-validate returned non-zero; see /tmp/virt-host-validate.log"
    else
      echo "[ok] virt-host-validate finished"
    fi
  fi
else
  echo "[preflight] mode=dashboard; checking dashboard dependencies"
  check_cmd curl
  check_cmd nginx
  check_cmd systemctl
fi

if [[ "$failures" -gt 0 ]]; then
  echo "[preflight] failed checks: $failures"
  if [[ "$ALLOW_MISSING" == "1" ]]; then
    echo "[preflight] ALLOW_MISSING=1 set; continuing despite failures"
    exit 0
  fi
  exit 1
fi

echo "[preflight] all required checks passed"
