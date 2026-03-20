import type { Bot } from '../types';

interface BotCardProps {
  bot: Bot;
  onRemove: (id: string) => void;
}

export const BotCard = ({ bot, onRemove }: BotCardProps) => {
  return (
    <div
      style={{
        border: '2px solid #ddd',
        borderRadius: '8px',
        padding: '16px',
        backgroundColor: bot.status === 'Processing' ? '#fff3cd' : '#d4edda',
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h3 style={{ margin: 0 }}>🤖 {bot.id}</h3>
        <button
          onClick={() => onRemove(bot.id)}
          style={{
            padding: '4px 12px',
            backgroundColor: '#dc3545',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
          }}
        >
          Remove
        </button>
      </div>
      <div style={{ marginTop: '12px' }}>
        <p style={{ margin: '4px 0' }}>
          <strong>Status:</strong> {bot.status}
        </p>
        <p style={{ margin: '4px 0' }}>
          <strong>Completed:</strong> {bot.completed_orders}
        </p>
        {bot.current_order && (
          <p style={{ margin: '4px 0' }}>
            <strong>Processing:</strong> {bot.current_order.id}
          </p>
        )}
      </div>
    </div>
  );
};
