# KVM Environment Deployment Guide

This directory includes everything required to deploy the platform in a KVM-based environment (bare-metal KVM hosts, nested lab, or VMs running on KVM).

## Contents
- `ansible/` automation for host baseline and service deployment.
- `systemd/kvm-platform.service` runtime service unit.
- `nginx/kvm-platform.conf` reverse proxy + API/UI routing.
- `cloud-init/kvm-node-user-data.yaml` bootstrap template for new KVM VMs.

## 1) Host prerequisites
For dashboard node (.160):
```bash
MODE=dashboard bash deploy/scripts/preflight_kvm.sh
```

For each KVM host (.153/.154):
```bash
MODE=host bash deploy/scripts/preflight_kvm.sh
```

Checks include:
- `qemu-system-x86_64`
- `virsh`
- `virt-host-validate`
- `/dev/kvm`
- libvirt config presence

## 2) Build release bundle
```bash
bash deploy/scripts/package_release.sh
```

This generates `release/kvm-platform-<timestamp>.tar.gz` containing:
- `bin/kvm-api`
- `web/*` built React assets
- systemd + nginx configs

## 3) Deploy to KVM VM/host
1. Copy release tarball to target host.
2. Extract under `/opt/kvm-platform`.
3. Install unit file from `systemd/kvm-platform.service`.
4. Install nginx config from `nginx/kvm-platform.conf`.
5. Start services:
   - `systemctl daemon-reload`
   - `systemctl enable --now kvm-platform`
   - `systemctl enable --now nginx`

## 4) Health verification
```bash
curl -fsS http://<host>/healthz
curl -fsS http://<host>/api/v1/dashboard
```

## 5) Optional automation (Ansible)
Use `deploy/kvm/ansible/site.yml` with your inventory to configure hosts end-to-end.

Topology template is prefilled in:
- `deploy/kvm/ansible/inventory.ini`
- `deploy/kvm/hosts.json`
