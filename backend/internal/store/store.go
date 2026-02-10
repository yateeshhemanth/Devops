package store

import (
	"fmt"
	"sync"
	"time"

	"devops/backend/internal/model"
)

type Store struct {
	mu       sync.RWMutex
	vms      map[string]model.VM
	volumes  map[string]model.Volume
	networks map[string]model.NetworkSegment
	hosts    map[string]model.Host
	alerts   []model.Alert
}

func New() *Store {
	return NewWithHosts(nil)
}

func NewWithHosts(hostList []model.Host) *Store {
	hosts := map[string]model.Host{
		"host-1": {ID: "host-1", Name: "kvm-host-01", Address: "10.0.10.11", CPUCapacity: 96, RAMCapacity: 512, State: "ready"},
		"host-2": {ID: "host-2", Name: "kvm-host-03", Address: "10.0.10.12", CPUCapacity: 64, RAMCapacity: 256, State: "ready"},
		"host-3": {ID: "host-3", Name: "kvm-host-04", Address: "10.0.10.13", CPUCapacity: 64, RAMCapacity: 256, State: "ready"},
	}

	if len(hostList) > 0 {
		hosts = make(map[string]model.Host, len(hostList))
		for _, host := range hostList {
			h := host
			if h.ID == "" {
				h.ID = h.Name
			}
			if h.State == "" {
				h.State = "ready"
			}
			hosts[h.ID] = h
		}
	}

	return &Store{
		vms: map[string]model.VM{
			"vm-1": {ID: "vm-1", Name: "payment-api-01", VCPU: 8, MemoryGB: 32, StorageGB: 500, Network: "prod-app-net", Status: model.VMRunning, Host: "kvm-host-03", SnapshotCnt: 4, LastBackupAt: time.Now().Add(-6 * time.Hour).Format(time.RFC3339)},
			"vm-2": {ID: "vm-2", Name: "erp-db-01", VCPU: 16, MemoryGB: 128, StorageGB: 2048, Network: "prod-db-net", Status: model.VMRunning, Host: "kvm-host-01", SnapshotCnt: 8, LastBackupAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339)},
			"vm-3": {ID: "vm-3", Name: "analytics-worker-07", VCPU: 12, MemoryGB: 64, StorageGB: 1000, Network: "batch-net", Status: model.VMPaused, Host: "kvm-host-04", SnapshotCnt: 3, LastBackupAt: time.Now().Add(-20 * time.Hour).Format(time.RFC3339)},
		},
		volumes: map[string]model.Volume{
			"vol-1": {ID: "vol-1", Name: "erp-data-rbd", Tier: model.TierGold, SizeGB: 2048, AttachedTo: "erp-db-01", IOPSLimit: 35000},
			"vol-2": {ID: "vol-2", Name: "payment-logs", Tier: model.TierSilver, SizeGB: 512, AttachedTo: "payment-api-01", IOPSLimit: 12000},
			"vol-3": {ID: "vol-3", Name: "backup-snapshots", Tier: model.TierBronze, SizeGB: 4096},
		},
		networks: map[string]model.NetworkSegment{
			"net-1": {ID: "net-1", Name: "prod-app-net", CIDR: "10.40.10.0/24", VLAN: 210, ProtectedBy: "sg-app"},
			"net-2": {ID: "net-2", Name: "prod-db-net", CIDR: "10.40.20.0/24", VLAN: 220, ProtectedBy: "sg-db"},
			"net-3": {ID: "net-3", Name: "batch-net", CIDR: "10.40.30.0/24", VLAN: 230, ProtectedBy: "sg-batch"},
		},
		hosts:  hosts,
		alerts: []model.Alert{{ID: "a-1", Severity: "warning", Message: "Host kvm-host-04 reports high memory pressure (82%)."}},
	}
}

func (s *Store) Dashboard() model.Dashboard {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return model.Dashboard{VMs: mapValues(s.vms), Volumes: mapValues(s.volumes), Networks: mapValues(s.networks), Hosts: mapValues(s.hosts), Alerts: append([]model.Alert(nil), s.alerts...)}
}

func (s *Store) SetPower(vmID string, action string) (model.VM, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	vm, ok := s.vms[vmID]
	if !ok {
		return model.VM{}, fmt.Errorf("vm %s not found", vmID)
	}
	switch action {
	case "start", "reboot":
		vm.Status = model.VMRunning
	case "stop":
		vm.Status = model.VMStopped
	case "pause":
		vm.Status = model.VMPaused
	default:
		return model.VM{}, fmt.Errorf("unsupported power action %q", action)
	}
	s.vms[vmID] = vm
	return vm, nil
}

func (s *Store) MigrateVM(vmID, host string) (model.VM, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	vm, ok := s.vms[vmID]
	if !ok {
		return model.VM{}, fmt.Errorf("vm %s not found", vmID)
	}
	if !s.hostExists(host) {
		return model.VM{}, fmt.Errorf("target host %s not found", host)
	}
	vm.Status = model.VMMigrating
	vm.Host = host
	s.vms[vmID] = vm
	vm.Status = model.VMRunning
	s.vms[vmID] = vm
	return vm, nil
}

func (s *Store) SnapshotVM(vmID string) (model.VM, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	vm, ok := s.vms[vmID]
	if !ok {
		return model.VM{}, fmt.Errorf("vm %s not found", vmID)
	}
	vm.SnapshotCnt++
	s.vms[vmID] = vm
	return vm, nil
}

func (s *Store) SetHostMaintenance(hostID string, enabled bool) (model.Host, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h, ok := s.hosts[hostID]
	if !ok {
		return model.Host{}, fmt.Errorf("host %s not found", hostID)
	}
	h.Maintenance = enabled
	if enabled {
		h.State = "maintenance"
	} else {
		h.State = "ready"
	}
	s.hosts[hostID] = h
	return h, nil
}

func (s *Store) ExpandVolume(id string, sizeGB int) (model.Volume, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.volumes[id]
	if !ok {
		return model.Volume{}, fmt.Errorf("volume %s not found", id)
	}
	if sizeGB <= v.SizeGB {
		return model.Volume{}, fmt.Errorf("new size must be greater than current size")
	}
	v.SizeGB = sizeGB
	s.volumes[id] = v
	return v, nil
}

func mapValues[T any](m map[string]T) []T {
	res := make([]T, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func (s *Store) hostExists(host string) bool {
	for _, h := range s.hosts {
		if h.ID == host || h.Name == host || h.Address == host {
			return true
		}
	}
	return false
}
