import type { Host } from '../types';

export function HostTable({ hosts, onToggleMaintenance }: { hosts: Host[]; onToggleMaintenance: (id: string, enabled: boolean) => void }) {
  return (
    <div className="card">
      <h2>Hosts</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>CPU</th>
            <th>RAM</th>
            <th>State</th>
            <th>Maintenance</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {hosts.map((host) => (
            <tr key={host.id}>
              <td>{host.name}</td>
              <td>{host.cpuCapacity}</td>
              <td>{host.ramCapacityGb} GB</td>
              <td>{host.state}</td>
              <td>{host.maintenance ? 'enabled' : 'disabled'}</td>
              <td>
                <button type="button" onClick={() => onToggleMaintenance(host.id, !host.maintenance)}>
                  {host.maintenance ? 'Disable' : 'Enable'} Maintenance
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
