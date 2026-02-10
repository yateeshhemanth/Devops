import type { Vm } from '../types';

interface VmTableProps {
  vms: Vm[];
  onPowerAction: (id: string, action: 'start' | 'stop' | 'reboot' | 'pause') => void;
  onOpenConsole: (id: string) => void;
  onSnapshot: (id: string) => void;
  onMigrate: (id: string) => void;
}

export function VmTable({ vms, onPowerAction, onOpenConsole, onSnapshot, onMigrate }: VmTableProps) {
  return (
    <div className="card">
      <h2>Virtual Machines</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>vCPU</th>
            <th>RAM</th>
            <th>Disk</th>
            <th>Status</th>
            <th>Host</th>
            <th>Snapshots</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {vms.map((vm) => (
            <tr key={vm.id}>
              <td>{vm.name}</td>
              <td>{vm.vcpu}</td>
              <td>{vm.memoryGb} GB</td>
              <td>{vm.storageGb} GB</td>
              <td><span className={`pill ${vm.status}`}>{vm.status}</span></td>
              <td>{vm.host}</td>
              <td>{vm.snapshotCount}</td>
              <td className="actions">
                <button type="button" onClick={() => onPowerAction(vm.id, 'start')}>Start</button>
                <button type="button" onClick={() => onPowerAction(vm.id, 'stop')}>Stop</button>
                <button type="button" onClick={() => onPowerAction(vm.id, 'reboot')}>Reboot</button>
                <button type="button" onClick={() => onPowerAction(vm.id, 'pause')}>Pause</button>
                <button type="button" onClick={() => onSnapshot(vm.id)}>Snapshot</button>
                <button type="button" onClick={() => onMigrate(vm.id)}>Migrate</button>
                <button type="button" className="primary" onClick={() => onOpenConsole(vm.id)}>noVNC</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
