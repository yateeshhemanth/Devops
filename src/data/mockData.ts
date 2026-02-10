import type { Dashboard } from '../types';

export const fallbackDashboard: Dashboard = {
  vms: [
    {
      id: 'vm-1',
      name: 'payment-api-01',
      vcpu: 8,
      memoryGb: 32,
      storageGb: 500,
      network: 'prod-app-net',
      status: 'running',
      host: 'kvm-host-03',
      snapshotCount: 4,
      lastBackupAt: new Date().toISOString()
    }
  ],
  volumes: [{ id: 'vol-1', name: 'erp-data-rbd', tier: 'gold', sizeGb: 2048, attachedTo: 'erp-db-01', iopsLimit: 35000 }],
  networks: [{ id: 'net-1', name: 'prod-app-net', cidr: '10.40.10.0/24', vlan: 210, protectedBy: 'sg-app' }],
  hosts: [{ id: 'host-1', name: 'kvm-host-01', cpuCapacity: 96, ramCapacityGb: 512, maintenance: false, state: 'ready' }],
  alerts: [{ id: 'a-1', severity: 'warning', message: 'Using fallback data (backend unavailable).' }]
};
