import type { Alert } from '../types';

export function AlertsPanel({ alerts }: { alerts: Alert[] }) {
  return (
    <div className="card">
      <h2>Operations Alerts</h2>
      <ul>
        {alerts.map((alert) => (
          <li key={alert.id}><strong>{alert.severity.toUpperCase()}:</strong> {alert.message}</li>
        ))}
      </ul>
    </div>
  );
}
