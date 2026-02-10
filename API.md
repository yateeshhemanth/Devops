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

### Kubernetes deployment
```bash
kubectl apply -f deploy/kubernetes/backend-deployment.yaml
```
