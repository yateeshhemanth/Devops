# Hosting Guide

Yes — this project can be hosted as **one service** (Go API + React UI) using the root `Dockerfile`, and also deployed to KVM hosts/VMs using systemd + nginx.

## What was added
- Root `Dockerfile` builds React and Go, then serves React from Go via `STATIC_DIR=/srv/web`.
- `render.yaml` for Render deployment.
- `fly.toml` for Fly.io deployment.
- `deploy/kvm/` assets for KVM deployment (Ansible, cloud-init, systemd, nginx).

## Quick local container run
```bash
docker build -t kvm-platform:latest .
docker run --rm -p 8080:8080 kvm-platform:latest
```
Open: `http://localhost:8080`

## Deploy on KVM hosts/VMs
```bash
make kvm-preflight
make package-release
```
Then follow `deploy/kvm/README.md` to install service files and nginx.

Use `deploy/kvm/hosts.json` for your host inventory (e.g., .153 and .154) and run dashboard on .160 with `HOSTS_FILE=/opt/kvm-platform/config/hosts.json`.

## Deploy on Render
1. Push this repo to GitHub.
2. In Render, create **New Web Service** from repo.
3. Render auto-detects `render.yaml`.
4. Deploy and open your generated URL.

## Deploy on Fly.io
```bash
fly launch --no-deploy
fly deploy
fly open
```

## Health check
- `GET /healthz`

## Notes
- For production, replace the in-memory store with PostgreSQL/Ceph-backed state.
- Lock down CORS and add auth (OIDC) before internet-facing use.
