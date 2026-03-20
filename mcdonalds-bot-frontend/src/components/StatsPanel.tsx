import type { Stats } from '../types';

interface StatsPanelProps {
  stats: Stats;
}

export const StatsPanel = ({ stats }: StatsPanelProps) => {
  return (
    <div
      style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
        gap: '16px',
        marginBottom: '24px',
      }}
    >
      <StatCard label="Total Bots" value={stats.total_bots} color="#007bff" />
      <StatCard label="Idle Bots" value={stats.idle_bots} color="#28a745" />
      <StatCard label="Processing" value={stats.processing_bots} color="#ffc107" />
      <StatCard label="Pending Orders" value={stats.pending_orders} color="#17a2b8" />
      <StatCard label="Completed" value={stats.completed_orders} color="#6c757d" />
      <StatCard label="Total Orders" value={stats.total_orders} color="#343a40" />
    </div>
  );
};

interface StatCardProps {
  label: string;
  value: number;
  color: string;
}

const StatCard = ({ label, value, color }: StatCardProps) => {
  return (
    <div
      style={{
        padding: '16px',
        backgroundColor: 'white',
        border: `2px solid ${color}`,
        borderRadius: '8px',
        textAlign: 'center',
      }}
    >
      <div style={{ fontSize: '32px', fontWeight: 'bold', color }}>{value}</div>
      <div style={{ fontSize: '14px', color: '#6c757d', marginTop: '4px' }}>{label}</div>
    </div>
  );
};
