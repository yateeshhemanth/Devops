import type { StorageVolume } from '../types';

export function StorageTable({ volumes, onExpand }: { volumes: StorageVolume[]; onExpand: (id: string, sizeGb: number) => void }) {
  return (
    <div className="card">
      <h2>Storage Volumes</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Tier</th>
            <th>Size</th>
            <th>Attached To</th>
            <th>IOPS Limit</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {volumes.map((volume) => (
            <tr key={volume.id}>
              <td>{volume.name}</td>
              <td><span className={`pill ${volume.tier}`}>{volume.tier}</span></td>
              <td>{volume.sizeGb} GB</td>
              <td>{volume.attachedTo ?? 'unattached'}</td>
              <td>{volume.iopsLimit ?? 'default'}</td>
              <td><button type="button" onClick={() => onExpand(volume.id, volume.sizeGb + 100)}>Expand +100GB</button></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
