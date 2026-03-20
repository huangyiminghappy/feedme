import { useEffect, useState } from 'react';
import type { Bot, Order, Stats, Event } from './types';
import { api } from './services/api';
import { useWebSocket } from './hooks/useWebSocket';
import { BotCard } from './components/BotCard';
import { OrderCard } from './components/OrderCard';
import { StatsPanel } from './components/StatsPanel';

function App() {
  const [bots, setBots] = useState<Bot[]>([]);
  const [pendingOrders, setPendingOrders] = useState<Order[]>([]);
  const [completedOrders, setCompletedOrders] = useState<Order[]>([]);
  const [stats, setStats] = useState<Stats>({
    total_bots: 0,
    idle_bots: 0,
    processing_bots: 0,
    pending_orders: 0,
    completed_orders: 0,
    total_orders: 0,
  });

  const [newBotId, setNewBotId] = useState('');
  const [newOrderId, setNewOrderId] = useState('');
  const [orderType, setOrderType] = useState<'Normal' | 'VIP'>('Normal');

  const fetchData = async () => {
    try {
      const [botsRes, pendingRes, completedRes, statsRes] = await Promise.all([
        api.getBots(),
        api.getPendingOrders(),
        api.getCompletedOrders(),
        api.getStats(),
      ]);

      setBots(botsRes.bots || []);
      setPendingOrders(pendingRes.orders || []);
      setCompletedOrders(completedRes.orders || []);
      setStats(statsRes);
    } catch (error) {
      console.error('Failed to fetch data:', error);
    }
  };

  const handleEvent = (event: Event) => {
    console.log('Event received:', event);
    fetchData();
  };

  const { connected } = useWebSocket(handleEvent);

  useEffect(() => {
    fetchData();
  }, []);

  const handleAddBot = async () => {
    if (!newBotId.trim()) return;
    try {
      await api.addBot(newBotId);
      setNewBotId('');
    } catch (error) {
      alert('Failed to add bot: ' + (error as Error).message);
    }
  };

  const handleRemoveBot = async (id: string) => {
    try {
      await api.removeBot(id);
    } catch (error) {
      alert('Failed to remove bot: ' + (error as Error).message);
    }
  };

  const handleCreateOrder = async () => {
    if (!newOrderId.trim()) return;
    try {
      await api.createOrder(newOrderId, orderType);
      setNewOrderId('');
    } catch (error) {
      alert('Failed to create order: ' + (error as Error).message);
    }
  };

  const handleReset = async () => {
    if (!confirm('Are you sure you want to reset the system?')) return;
    try {
      await api.reset();
      fetchData();
    } catch (error) {
      alert('Failed to reset: ' + (error as Error).message);
    }
  };

  return (
    <div style={{ padding: '24px', maxWidth: '1400px', margin: '0 auto', fontFamily: 'Arial, sans-serif' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h1 style={{ margin: 0 }}>🍔 McDonald's Bot Management System</h1>
        <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
          <div
            style={{
              width: '12px',
              height: '12px',
              borderRadius: '50%',
              backgroundColor: connected ? '#28a745' : '#dc3545',
            }}
          />
          <span style={{ fontSize: '14px', color: '#6c757d' }}>
            {connected ? 'Connected' : 'Disconnected'}
          </span>
          <button
            onClick={handleReset}
            style={{
              padding: '8px 16px',
              backgroundColor: '#dc3545',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Reset System
          </button>
        </div>
      </div>

      <StatsPanel stats={stats} />

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '24px', marginBottom: '24px' }}>
        <div
          style={{
            padding: '20px',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
            border: '2px solid #dee2e6',
          }}
        >
          <h2 style={{ marginTop: 0 }}>Add Bot</h2>
          <div style={{ display: 'flex', gap: '8px' }}>
            <input
              type="text"
              value={newBotId}
              onChange={(e) => setNewBotId(e.target.value)}
              placeholder="Bot ID (e.g., bot1)"
              style={{
                flex: 1,
                padding: '8px 12px',
                border: '1px solid #ced4da',
                borderRadius: '4px',
              }}
              onKeyPress={(e) => e.key === 'Enter' && handleAddBot()}
            />
            <button
              onClick={handleAddBot}
              style={{
                padding: '8px 16px',
                backgroundColor: '#28a745',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
              }}
            >
              Add Bot
            </button>
          </div>
        </div>

        <div
          style={{
            padding: '20px',
            backgroundColor: '#f8f9fa',
            borderRadius: '8px',
            border: '2px solid #dee2e6',
          }}
        >
          <h2 style={{ marginTop: 0 }}>Create Order</h2>
          <div style={{ display: 'flex', gap: '8px' }}>
            <input
              type="text"
              value={newOrderId}
              onChange={(e) => setNewOrderId(e.target.value)}
              placeholder="Order ID (e.g., order1)"
              style={{
                flex: 1,
                padding: '8px 12px',
                border: '1px solid #ced4da',
                borderRadius: '4px',
              }}
              onKeyPress={(e) => e.key === 'Enter' && handleCreateOrder()}
            />
            <select
              value={orderType}
              onChange={(e) => setOrderType(e.target.value as 'Normal' | 'VIP')}
              style={{
                padding: '8px 12px',
                border: '1px solid #ced4da',
                borderRadius: '4px',
              }}
            >
              <option value="Normal">Normal</option>
              <option value="VIP">VIP</option>
            </select>
            <button
              onClick={handleCreateOrder}
              style={{
                padding: '8px 16px',
                backgroundColor: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
              }}
            >
              Create
            </button>
          </div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '24px' }}>
        <div>
          <h2>🤖 Bots ({bots.length})</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {bots.length === 0 ? (
              <p style={{ color: '#6c757d' }}>No bots available</p>
            ) : (
              bots.map((bot) => <BotCard key={bot.id} bot={bot} onRemove={handleRemoveBot} />)
            )}
          </div>
        </div>

        <div>
          <h2>⏳ Pending Orders ({pendingOrders.length})</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            {pendingOrders.length === 0 ? (
              <p style={{ color: '#6c757d' }}>No pending orders</p>
            ) : (
              pendingOrders.map((order) => <OrderCard key={order.id} order={order} />)
            )}
          </div>
        </div>

        <div>
          <h2>✅ Completed Orders ({completedOrders.length})</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '12px', maxHeight: '600px', overflowY: 'auto' }}>
            {completedOrders.length === 0 ? (
              <p style={{ color: '#6c757d' }}>No completed orders</p>
            ) : (
              completedOrders.map((order) => <OrderCard key={order.id} order={order} />)
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
