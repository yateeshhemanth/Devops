#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DIST_DIR="${ROOT_DIR}/release"
VERSION="${VERSION:-$(date +%Y%m%d%H%M%S)}"
PKG_DIR="${DIST_DIR}/kvm-platform-${VERSION}"

rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR/bin" "$PKG_DIR/web" "$PKG_DIR/systemd" "$PKG_DIR/nginx" "$PKG_DIR/config"

echo "[package] building backend binary"
(cd "$ROOT_DIR/backend" && go build -o "$PKG_DIR/bin/kvm-api" ./cmd/api)

echo "[package] building frontend"
(cd "$ROOT_DIR" && npm install && npm run build)
cp -r "$ROOT_DIR/dist"/* "$PKG_DIR/web/"

cp "$ROOT_DIR/deploy/kvm/systemd/kvm-platform.service" "$PKG_DIR/systemd/"
cp "$ROOT_DIR/deploy/kvm/nginx/kvm-platform.conf" "$PKG_DIR/nginx/"
cp "$ROOT_DIR/deploy/kvm/hosts.json" "$PKG_DIR/config/hosts.json"

cat > "$PKG_DIR/.env.example" <<ENV
PORT=8080
STATIC_DIR=/opt/kvm-platform/web
HOSTS_FILE=/opt/kvm-platform/config/hosts.json
API_KEYS=viewer-token:viewer,ops-token:operator,admin-token:admin
ENV

tar -C "$DIST_DIR" -czf "$DIST_DIR/kvm-platform-${VERSION}.tar.gz" "kvm-platform-${VERSION}"

echo "[package] artifact: $DIST_DIR/kvm-platform-${VERSION}.tar.gz"
