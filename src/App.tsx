import { useEffect, useMemo, useState } from 'react';
import { api } from './api';
import { AlertsPanel } from './components/AlertsPanel';
import { HostTable } from './components/HostTable';
import { NetworkTable } from './components/NetworkTable';
import { StatCard } from './components/StatCard';
import { StorageTable } from './components/StorageTable';
import { Tabs } from './components/Tabs';
import { VmTable } from './components/VmTable';
import { fallbackDashboard } from './data/mockData';
import type { Dashboard } from './types';

const tabs = ['Overview', 'Compute', 'Storage', 'Network', 'Hosts'] as const;
type Tab = (typeof tabs)[number];

export function App() {
  const [selectedTab, setSelectedTab] = useState<Tab>('Overview');
  const [dashboard, setDashboard] = useState<Dashboard>(fallbackDashboard);
  const [consoleMessage, setConsoleMessage] = useState('');
  const [error, setError] = useState('');

  const reload = async () => {
    try {
      const data = await api.dashboard();
      setDashboard(data);
      setError('');
    } catch {
      setError('Backend unavailable. Showing fallback data.');
    }
  };

  useEffect(() => {
    void reload();
  }, []);

  const stats = useMemo(() => {
    const running = dashboard.vms.filter((vm) => vm.status === 'running').length;
    const totalCpu = dashboard.vms.reduce((acc, vm) => acc + vm.vcpu, 0);
    const totalRam = dashboard.vms.reduce((acc, vm) => acc + vm.memoryGb, 0);
    return { running, totalCpu, totalRam };
  }, [dashboard.vms]);

  const act = async (action: () => Promise<unknown>) => {
    await action();
    await reload();
  };

  return (
    <main className="layout">
      <header>
        <h1>KVM Enterprise Platform (Go + React)</h1>
        <p>OpenShift/vSphere-style day-2 operations: VM lifecycle, snapshots, migration, storage growth, host maintenance, and noVNC.</p>
      </header>

      {error && <p className="error">{error}</p>}

      <section className="stat-grid">
        <StatCard label="Running VMs" value={`${stats.running}/${dashboard.vms.length}`} health="healthy" />
        <StatCard label="Allocated vCPU" value={`${stats.totalCpu}`} health="healthy" />
        <StatCard label="Allocated RAM" value={`${stats.totalRam} GB`} health="warning" />
        <StatCard label="Storage Volumes" value={`${dashboard.volumes.length}`} health="healthy" />
      </section>

      <Tabs selected={selectedTab} items={[...tabs]} onSelect={(tab) => setSelectedTab(tab as Tab)} />

      {selectedTab === 'Overview' && (
        <section className="grid-two">
          <VmTable
            vms={dashboard.vms}
            onPowerAction={(id, power) => act(() => api.power(id, power))}
            onSnapshot={(id) => act(() => api.snapshot(id))}
            onMigrate={(id) => act(() => api.migrate(id, 'kvm-host-01'))}
            onOpenConsole={async (id) => {
              const ticket = await api.consoleTicket(id);
              setConsoleMessage(`Ticket: ${ticket.ticket} (expires ${ticket.expiresInS}s) -> ${ticket.wsURL}`);
            }}
          />
          <AlertsPanel alerts={dashboard.alerts} />
        </section>
      )}

      {selectedTab === 'Compute' && (
        <VmTable
          vms={dashboard.vms}
          onPowerAction={(id, power) => act(() => api.power(id, power))}
          onSnapshot={(id) => act(() => api.snapshot(id))}
          onMigrate={(id) => act(() => api.migrate(id, 'kvm-host-03'))}
          onOpenConsole={async (id) => {
            const ticket = await api.consoleTicket(id);
            setConsoleMessage(`Ticket: ${ticket.ticket} (expires ${ticket.expiresInS}s) -> ${ticket.wsURL}`);
          }}
        />
      )}
      {selectedTab === 'Storage' && <StorageTable volumes={dashboard.volumes} onExpand={(id, size) => act(() => api.expandVolume(id, size))} />}
      {selectedTab === 'Network' && <NetworkTable networks={dashboard.networks} />}
      {selectedTab === 'Hosts' && (
        <HostTable
          hosts={dashboard.hosts}
          onToggleMaintenance={(id, enabled) => act(() => api.setMaintenance(id, enabled))}
        />
      )}

      {consoleMessage && (
        <section className="card console-banner">
          <h3>noVNC Session</h3>
          <p>{consoleMessage}</p>
        </section>
      )}
    </main>
  );
}
