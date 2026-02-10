import type { Dashboard, Host, StorageVolume, Vm } from './types';

const BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080';

async function req<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: response.statusText }));
    throw new Error(body.error ?? 'request failed');
  }

  return response.json() as Promise<T>;
}

export const api = {
  dashboard: () => req<Dashboard>('/api/v1/dashboard'),
  power: (id: string, action: 'start' | 'stop' | 'reboot' | 'pause') =>
    req<Vm>(`/api/v1/vms/${id}/power`, { method: 'POST', body: JSON.stringify({ action }) }),
  migrate: (id: string, host: string) =>
    req<Vm>(`/api/v1/vms/${id}/migrate`, { method: 'POST', body: JSON.stringify({ host }) }),
  snapshot: (id: string) => req<Vm>(`/api/v1/vms/${id}/snapshot`, { method: 'POST' }),
  consoleTicket: (id: string) =>
    req<{ ticket: string; expiresInS: number; wsURL: string }>(`/api/v1/vms/${id}/console-ticket`, { method: 'POST' }),
  expandVolume: (id: string, sizeGb: number) =>
    req<StorageVolume>(`/api/v1/storage/volumes/${id}/expand`, { method: 'POST', body: JSON.stringify({ sizeGb }) }),
  setMaintenance: (id: string, enabled: boolean) =>
    req<Host>(`/api/v1/hosts/${id}/maintenance`, { method: 'POST', body: JSON.stringify({ enabled }) })
};
