# Deployment Runbook (.160 dashboard + .153/.154 KVM hosts)

This is the exact deployment flow for your topology.

## Topology
- `192.168.1.160` → Dashboard/API server
- `192.168.1.153` → KVM host 1
- `192.168.1.154` → KVM host 2

## 1) Prepare this repo machine

```bash
make test-backend
```

Run preflight checks (informational in non-KVM environments):
```bash
MODE=dashboard make kvm-preflight
MODE=host make kvm-preflight
```

Build the release package:
```bash
make package-release
```

## 2) Configure host inventory
Edit `deploy/kvm/hosts.json` with your final host metadata.

Current default is already:
- `192.168.1.153`
- `192.168.1.154`

## 3) Deploy with Ansible

Update SSH usernames/keys as needed in `deploy/kvm/ansible/inventory.ini`.

Copy release package to Ansible control machine location `/tmp/kvm-platform.tar.gz`:
```bash
cp release/kvm-platform-*.tar.gz /tmp/kvm-platform.tar.gz
```

Run playbook:
```bash
ansible-playbook -i deploy/kvm/ansible/inventory.ini deploy/kvm/ansible/site.yml
```

What this does:
- Installs KVM/libvirt packages on `.153`/`.154`
- Installs dashboard runtime on `.160`
- Copies host inventory to `/opt/kvm-platform/config/hosts.json`
- Enables `kvm-platform` service and `nginx`

## 4) Configure production auth (recommended)
On `.160` edit `/opt/kvm-platform/.env`:

```env
API_KEYS=viewer-token:viewer,ops-token:operator,admin-token:admin
HOSTS_FILE=/opt/kvm-platform/config/hosts.json
PORT=8080
STATIC_DIR=/opt/kvm-platform/web
```

Then restart service:
```bash
sudo systemctl restart kvm-platform
```

## 5) Verify deployment
From `.160`:
```bash
curl -fsS http://127.0.0.1/healthz
curl -fsS http://127.0.0.1/readyz
curl -fsS -H 'X-API-Key: viewer-token' http://127.0.0.1/api/v1/dashboard
```

From your laptop/browser:
- Open `http://192.168.1.160`

## 6) Troubleshooting
Check service logs:
```bash
sudo journalctl -u kvm-platform -f
sudo journalctl -u nginx -f
```

If dashboard doesn’t show host list, verify:
- `/opt/kvm-platform/config/hosts.json` exists
- `HOSTS_FILE` in `/opt/kvm-platform/.env`
- `sudo systemctl restart kvm-platform`
