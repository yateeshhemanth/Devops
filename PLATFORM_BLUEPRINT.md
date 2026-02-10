# Enterprise KVM Platform Blueprint

## 1) Product Goal
Build a multi-tenant virtualization platform on KVM that offers day-to-day operations similar to OpenShift Virtualization + VMware vSphere:
- VM lifecycle: create, clone, start/stop/reboot, migrate, snapshot/backup, resize.
- Storage lifecycle: create pools/classes, attach/detach volumes, clone, snapshot, tiering.
- Network lifecycle: virtual networks, VLAN/overlay segments, security groups, floating IPs, LB integration.
- Day-2 operations: monitoring, alerting, patching, upgrades, quota, chargeback, audit, policy.
- Browser console access through **noVNC** (VNC over WebSockets) with RBAC and session recording.

## 2) High-Level Architecture

```text
[Web UI + API Gateway]
        |
[AuthN/AuthZ + Policy + Tenant Service]
        |
[Control Plane Orchestrator]
   |      |       |       |
 VM svc  Net svc Storage  Image svc
   |      |       |       |
 --------------- Event Bus --------------
   |      |       |       |
[Compute Agents] [Network Agents] [Storage Agents]
      \      |       /
         [KVM/libvirt hosts]
             |
      [Ceph/NFS/Local/ZFS]
```

### Core Components
1. **API Gateway**
   - Northbound API (REST + optional GraphQL) and UI backend.
   - Request validation, rate limits, idempotency keys.
2. **Identity & Access**
   - OIDC/SAML integration, RBAC/ABAC, tenant/project hierarchy.
3. **Orchestrator**
   - Reconciler pattern: desired state vs actual state.
   - Handles long-running workflows with retries and compensation.
4. **Compute Service**
   - Manages VM specs, CPU pinning, NUMA, hugepages, affinity, live migration.
5. **Storage Service**
   - CSI-like abstraction for Ceph RBD, NFS, iSCSI, local NVMe, ZFS.
   - Snapshots, clones, volume expansion, encryption-at-rest.
6. **Network Service**
   - Bridges/OVS, VLAN/VXLAN, NAT, DHCP/DNS, ACL/security groups.
7. **Image Service**
   - Golden images/templates, cloud-init snippets, versioning/signing.
8. **Console Service**
   - Integrates libvirt VNC/SPICE endpoints with **websockify/noVNC**.
   - Ephemeral signed console tickets + audit logs.
9. **Telemetry Service**
   - Metrics (Prometheus), logs (Loki/ELK), traces (OTel), alert rules.
10. **Billing/Quota Service**
    - Quotas for vCPU/RAM/storage/IP and optional chargeback.

## 3) Technology Stack Recommendation

### Control Plane
- **Language**: Go (strong concurrency, Kubernetes/libvirt ecosystem).
- **API**: gRPC internally, REST externally.
- **Workflow engine**: Temporal or durable job queue + DB-backed state machine.
- **Database**: PostgreSQL (transactional metadata), Redis (cache/locks).
- **Event bus**: NATS or Kafka.

### Virtualization & Host Layer
- KVM + libvirt + qemu.
- Optional Kubernetes control plane + KubeVirt if you want cloud-native scheduling semantics.
- Host agents in Go/Rust for secure command execution.

### Storage
- Primary: Ceph (RBD for block + CephFS/object if needed).
- Secondary connectors: NFS, iSCSI, local LVM thin.

### Network
- Open vSwitch + OVN (or Calico/Cilium if Kubernetes-based).
- SR-IOV integration for high-performance workloads.

### UI
- React + TypeScript + a design system.
- xterm.js for serial console; noVNC embedded for graphics console.

## 4) Core Day-2 Features (Must-Have)
1. **HA control plane** (3+ nodes, DB replication, leader election).
2. **Live migration** with pre-checks and rollback.
3. **Snapshots & backup** with policy schedules.
4. **Role-based approvals** for risky operations (delete, migrate, resize).
5. **Audit trail** (who changed what, when, from where).
6. **Policy engine** (OPA): image allow-list, CPU model constraints, regions.
7. **Patch manager** for rolling host upgrades with evacuation.
8. **Disaster recovery** runbooks: control plane restore + image/volume restore.

## 5) noVNC Console Design

### Session Flow
1. User clicks "Open Console" in UI.
2. UI calls `POST /v1/vms/{id}/console-ticket`.
3. Console service validates RBAC and VM power state.
4. Service issues short-lived JWT (30–120 sec) bound to user+VM+tenant.
5. Browser connects to console gateway `/console/ws?ticket=...`.
6. Gateway validates ticket, opens backend VNC/SPICE tunnel via websockify.
7. Session is logged (open/close, source IP, duration).

### Security Requirements
- One-time or very short-lived tickets.
- TLS end-to-end and optional mTLS agent links.
- Rate limiting and brute-force protections.
- Optional session watermark and recording metadata.

## 6) API Surface (Example)

### VM APIs
- `POST /v1/vms`
- `GET /v1/vms/{id}`
- `POST /v1/vms/{id}/power` (`start|stop|reboot|pause`)
- `POST /v1/vms/{id}/migrate`
- `POST /v1/vms/{id}/snapshot`
- `POST /v1/vms/{id}/clone`

### Storage APIs
- `POST /v1/storageclasses`
- `POST /v1/volumes`
- `POST /v1/volumes/{id}/attach`
- `POST /v1/volumes/{id}/snapshot`
- `POST /v1/volumes/{id}/expand`

### Network APIs
- `POST /v1/networks`
- `POST /v1/security-groups`
- `POST /v1/floating-ips`
- `POST /v1/loadbalancers`

## 7) Multi-Tenancy & Governance
- Resource hierarchy: `org -> project -> environment`.
- Hard quotas + soft quotas.
- Namespace isolation on network and storage.
- Tagging and policy inheritance.
- Audit export to SIEM.

## 8) Reliability/SRE Targets
- API availability target: 99.95%.
- RPO: < 15 min (with continuous backup metadata replication).
- RTO: < 1 hour for control plane.
- SLOs for VM provisioning latency (P50/P95).

## 9) Phased Delivery Plan

### Phase 0 (4–6 weeks): Foundation
- Identity, tenancy, project model.
- Host enrollment and health model.
- Basic VM create/start/stop/delete.
- Initial noVNC integration.

### Phase 1 (6–8 weeks): Production MVP
- Storage classes + dynamic volumes.
- Network segments + security groups.
- Snapshots + templates + cloning.
- Metrics/logging/auditing.

### Phase 2 (8–12 weeks): Enterprise
- Live migration + HA + host maintenance mode.
- Backup policies + DR workflows.
- Quotas, approvals, policy engine.
- Chargeback/showback reports.

### Phase 3: Advanced Platform
- GPU/SR-IOV profiles.
- Autoscaling and placement optimization.
- Marketplace/catalog of prebuilt stacks.
- Self-service IaC provider (Terraform-compatible).

## 10) Suggested Repository Layout

```text
/platform
  /cmd
    /api-server
    /controller
    /host-agent
    /console-gateway
  /internal
    /auth
    /compute
    /storage
    /network
    /console
    /scheduler
    /policy
  /pkg
  /deploy
    /kubernetes
    /ansible
  /ui
  /docs
```

## 11) First 90-Day Execution Checklist
- Finalize architecture decision records (ADRs).
- Stand up 3-node management cluster + 4 KVM hosts + Ceph.
- Deliver API + UI for VM lifecycle and console.
- Add canary + integration tests for critical workflows.
- Establish production readiness review gates.

---

If you want, next step can be a **detailed implementation backlog** with epics, user stories, and acceptance criteria for each operation (VM, storage, network, console, IAM, monitoring).
