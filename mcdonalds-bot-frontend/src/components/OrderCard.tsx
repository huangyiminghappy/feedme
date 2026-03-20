import type { Order } from '../types';

interface OrderCardProps {
  order: Order;
}

export const OrderCard = ({ order }: OrderCardProps) => {
  const getStatusColor = () => {
    switch (order.status) {
      case 'Pending':
        return '#ffc107';
      case 'Processing':
        return '#17a2b8';
      case 'Complete':
        return '#28a745';
      default:
        return '#6c757d';
    }
  };

  return (
    <div
      style={{
        border: '2px solid #ddd',
        borderRadius: '8px',
        padding: '12px',
        backgroundColor: order.type === 'VIP' ? '#fff3e0' : '#f8f9fa',
        borderLeft: `4px solid ${getStatusColor()}`,
      }}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h4 style={{ margin: 0 }}>
          {order.type === 'VIP' ? '⭐' : '📦'} {order.id}
        </h4>
        <span
          style={{
            padding: '2px 8px',
            backgroundColor: getStatusColor(),
            color: 'white',
            borderRadius: '4px',
            fontSize: '12px',
          }}
        >
          {order.status}
        </span>
      </div>
      <div style={{ marginTop: '8px', fontSize: '14px' }}>
        <p style={{ margin: '2px 0' }}>
          <strong>Type:</strong> {order.type}
        </p>
        {order.bot_id && (
          <p style={{ margin: '2px 0' }}>
            <strong>Bot:</strong> {order.bot_id}
          </p>
        )}
      </div>
    </div>
  );
};
