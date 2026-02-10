export type Health = 'healthy' | 'warning' | 'critical';

export interface Vm {
  id: string;
  name: string;
  vcpu: number;
  memoryGb: number;
  storageGb: number;
  network: string;
  status: 'running' | 'stopped' | 'paused' | 'migrating';
  host: string;
  snapshotCount: number;
  lastBackupAt: string;
}

export interface StorageVolume {
  id: string;
  name: string;
  tier: 'gold' | 'silver' | 'bronze';
  sizeGb: number;
  attachedTo?: string;
  iopsLimit?: number;
}

export interface NetworkSegment {
  id: string;
  name: string;
  cidr: string;
  vlan: number;
  protectedBy: string;
}

export interface Host {
  id: string;
  name: string;
  address?: string;
  cpuCapacity: number;
  ramCapacityGb: number;
  maintenance: boolean;
  state: string;
}

export interface Alert {
  id: string;
  severity: string;
  message: string;
}

export interface Dashboard {
  vms: Vm[];
  volumes: StorageVolume[];
  networks: NetworkSegment[];
  hosts: Host[];
  alerts: Alert[];
}
