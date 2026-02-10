import type { Health } from '../types';

interface StatCardProps {
  label: string;
  value: string;
  health?: Health;
}

const healthClass: Record<Health, string> = {
  healthy: 'health healthy',
  warning: 'health warning',
  critical: 'health critical'
};

export function StatCard({ label, value, health = 'healthy' }: StatCardProps) {
  return (
    <article className="card stat-card">
      <div>
        <h3>{label}</h3>
        <p className="stat-value">{value}</p>
      </div>
      <span className={healthClass[health]}>{health}</span>
    </article>
  );
}
