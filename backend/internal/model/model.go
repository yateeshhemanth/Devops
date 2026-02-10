package model

type VMStatus string

const (
	VMRunning   VMStatus = "running"
	VMStopped   VMStatus = "stopped"
	VMPaused    VMStatus = "paused"
	VMMigrating VMStatus = "migrating"
)

type VM struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	VCPU         int      `json:"vcpu"`
	MemoryGB     int      `json:"memoryGb"`
	StorageGB    int      `json:"storageGb"`
	Network      string   `json:"network"`
	Status       VMStatus `json:"status"`
	Host         string   `json:"host"`
	SnapshotCnt  int      `json:"snapshotCount"`
	LastBackupAt string   `json:"lastBackupAt"`
}

type StorageTier string

const (
	TierGold   StorageTier = "gold"
	TierSilver StorageTier = "silver"
	TierBronze StorageTier = "bronze"
)

type Volume struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Tier       StorageTier `json:"tier"`
	SizeGB     int         `json:"sizeGb"`
	AttachedTo string      `json:"attachedTo,omitempty"`
	IOPSLimit  int         `json:"iopsLimit,omitempty"`
}

type NetworkSegment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CIDR        string `json:"cidr"`
	VLAN        int    `json:"vlan"`
	ProtectedBy string `json:"protectedBy"`
}

type Host struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CPUCapacity int    `json:"cpuCapacity"`
	RAMCapacity int    `json:"ramCapacityGb"`
	Maintenance bool   `json:"maintenance"`
	State       string `json:"state"`
}

type Alert struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type Dashboard struct {
	VMs      []VM             `json:"vms"`
	Volumes  []Volume         `json:"volumes"`
	Networks []NetworkSegment `json:"networks"`
	Hosts    []Host           `json:"hosts"`
	Alerts   []Alert          `json:"alerts"`
}
