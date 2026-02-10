# Backend API (Go)

Base URL: `http://localhost:8080`

## Health
- `GET /healthz`

## Dashboard
- `GET /api/v1/dashboard`

## VM day-2 operations
- `POST /api/v1/vms/{id}/power` body: `{ "action": "start|stop|reboot|pause" }`
- `POST /api/v1/vms/{id}/migrate` body: `{ "host": "kvm-host-01" }`
- `POST /api/v1/vms/{id}/snapshot`
- `POST /api/v1/vms/{id}/console-ticket`

## Storage operations
- `POST /api/v1/storage/volumes/{id}/expand` body: `{ "sizeGb": 2148 }`

## Host operations
- `POST /api/v1/hosts/{id}/maintenance` body: `{ "enabled": true }`

## Test + Deploy

### Backend tests
```bash
make test-backend
```

### Local deployment (backend)
```bash
make deploy-local
```
This command builds the Go backend, starts it in the background, and runs smoke tests.

### Backend smoke-only
```bash
make smoke-backend
```

### KVM readiness checks
```bash
make kvm-preflight
```

### Build KVM release bundle
```bash
make package-release
```

### Kubernetes deployment
```bash
kubectl apply -f deploy/kubernetes/backend-deployment.yaml
```

### One-container hosting (UI + API)
```bash
docker build -t kvm-platform:latest .
docker run --rm -p 8080:8080 kvm-platform:latest
```
Open `http://localhost:8080`.

For multi-host KVM setup set:
- `HOSTS_FILE=/opt/kvm-platform/config/hosts.json` (dashboard .160 loads .153/.154 inventory).
