import type { NetworkSegment } from '../types';

export function NetworkTable({ networks }: { networks: NetworkSegment[] }) {
  return (
    <div className="card">
      <h2>Network Segments</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>CIDR</th>
            <th>VLAN</th>
            <th>Security Group</th>
          </tr>
        </thead>
        <tbody>
          {networks.map((network) => (
            <tr key={network.id}>
              <td>{network.name}</td>
              <td>{network.cidr}</td>
              <td>{network.vlan}</td>
              <td>{network.protectedBy}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
