# KVM Platform (Go + React)

An operator-focused virtualization control-plane prototype inspired by OpenShift Virtualization / VMware day-2 workflows.

This repo includes:
- **Go backend API** for day-2 operations (VM power, snapshot, migrate, volume expand, host maintenance, console ticket).
- **React frontend** dashboard for operations UI.
- **Deployment assets** for local run, container hosting, Kubernetes, and KVM host rollout.

---

## Features

### Day-2 operations
- VM power actions: start/stop/reboot/pause
- VM snapshot and migration
- noVNC console ticket endpoint
- Storage volume expansion
- Host maintenance mode toggling
- Dashboard inventory and alerts

### Deployment targets
- Local process deployment (system service style)
- Single-container hosting (UI + API)
- Kubernetes deployment manifest
- KVM deployment bundle (systemd + nginx + ansible + cloud-init)

---

## Repository structure

- `backend/` — Go API service
- `src/` — React UI
- `deploy/scripts/` — local deploy/smoke/preflight/package scripts
- `deploy/kubernetes/` — k8s manifest
- `deploy/kvm/` — KVM production-style deployment assets
- `Dockerfile` — one-image build (API + UI)
- `API.md` — backend API and command reference
- `HOSTING.md` — hosting options

---

## Quick start (backend)

### 1) Run tests
```bash
make test-backend
```

### 2) Deploy locally + smoke test
```bash
make deploy-local
```

This command builds backend binary, starts service on `:8080`, and executes smoke checks.

### 3) Verify
```bash
curl -fsS http://localhost:8080/healthz
curl -fsS http://localhost:8080/api/v1/dashboard
```

---

## Frontend

The frontend is Vite + React + TypeScript.

```bash
npm install
npm run dev
```

> If your environment blocks npm registry access, installs can fail with 403.

---

## API summary

Base URL: `http://localhost:8080`

- `GET /healthz`
- `GET /api/v1/dashboard`
- `POST /api/v1/vms/{id}/power` `{ "action": "start|stop|reboot|pause" }`
- `POST /api/v1/vms/{id}/migrate` `{ "host": "kvm-host-01" }`
- `POST /api/v1/vms/{id}/snapshot`
- `POST /api/v1/vms/{id}/console-ticket`
- `POST /api/v1/storage/volumes/{id}/expand` `{ "sizeGb": 2148 }`
- `POST /api/v1/hosts/{id}/maintenance` `{ "enabled": true }`

See: `API.md`.

---

## Deploy to KVM environment

Recommended topology for your case:
- `192.168.1.160` = dashboard/API server
- `192.168.1.153` = KVM host 1
- `192.168.1.154` = KVM host 2


### 1) Host preflight
```bash
MODE=dashboard make kvm-preflight
```

On each hypervisor host run:
```bash
MODE=host make kvm-preflight
```

Checks:
- qemu/libvirt tooling (`qemu-system-x86_64`, `virsh`, `virt-host-validate`)
- `/dev/kvm`
- libvirt config presence

### 2) Build KVM release bundle
```bash
make package-release
```

Creates:
- `release/kvm-platform-<timestamp>.tar.gz`

### 3) Install on KVM VM/host
Follow `deploy/kvm/README.md`:
- extract under `/opt/kvm-platform`
- install `systemd` service
- install `nginx` reverse proxy config
- start `kvm-platform` + `nginx`

### 4) Optional automation
Use Ansible playbook:
- `deploy/kvm/ansible/site.yml`
- edit `deploy/kvm/hosts.json` to list your real host IPs

---

## Container hosting (single service)

```bash
docker build -t kvm-platform:latest .
docker run --rm -p 8080:8080 kvm-platform:latest
```

Open:
- `http://localhost:8080`

Also included:
- `render.yaml`
- `fly.toml`

---

## Kubernetes deployment

```bash
kubectl apply -f deploy/kubernetes/backend-deployment.yaml
```

---

## Important notes

- Current backend state is **in-memory** (prototype behavior).
- For production: use persistent datastore + auth (OIDC) + stricter CORS + backup/DR.
- This repo provides practical deployment scaffolding for KVM-style environments, but hardening is still required before internet-facing production.
