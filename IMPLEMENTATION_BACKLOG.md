# Implementation Backlog (Epics + Initial Stories)

## Epic A: Identity, Tenancy, and RBAC
- A1: Integrate OIDC login and token validation.
- A2: Support org/project roles (`viewer`, `operator`, `admin`).
- A3: Enforce API authorization middleware with policy checks.
- A4: Add audit records for all mutating operations.

## Epic B: Host Lifecycle & Compute
- B1: Host registration with signed bootstrap token.
- B2: Host capability discovery (CPU model, NUMA, hugepages, SR-IOV).
- B3: VM create/read/update/delete with reconciliation loop.
- B4: Power actions and status transitions.
- B5: Placement engine v1 (resource-fit + anti-affinity).

## Epic C: Storage
- C1: StorageClass abstraction (ceph-rbd, nfs, local-lvm).
- C2: Volume provisioning, attach/detach workflow.
- C3: Snapshot and clone API.
- C4: Volume expansion with online/offline validation.
- C5: Backup policy schedules and retention lifecycle.

## Epic D: Network
- D1: Network objects (VLAN/VXLAN-backed).
- D2: DHCP and DNS integration.
- D3: Security groups and rules.
- D4: Floating IP allocation and association.
- D5: Load balancer integration hooks.

## Epic E: Console (noVNC)
- E1: Console ticket API with short-lived JWT.
- E2: WebSocket gateway with ticket verification.
- E3: Libvirt VNC/SPICE proxy handling.
- E4: Session audit logs and timeouts.
- E5: UI embedding and reconnect logic.

## Epic F: Observability and Operations
- F1: Metrics endpoint and Prometheus dashboards.
- F2: Central log pipeline and query UI links.
- F3: Alerts for host down, storage pressure, failed reconciles.
- F4: Upgrade orchestration with host drain/evacuate.
- F5: SLO reports for API and provisioning latency.

## Definition of Done (Global)
- API endpoint documented in OpenAPI.
- Unit tests and integration tests for critical path.
- Audit + metrics emitted.
- RBAC and tenancy checks enforced.
- Runbook updated.
