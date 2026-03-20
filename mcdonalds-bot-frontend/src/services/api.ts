const API_BASE_URL = 'http://localhost:8080/api';

export const api = {
  async createOrder(id: string, type: 'Normal' | 'VIP') {
    const response = await fetch(`${API_BASE_URL}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, type }),
    });
    if (!response.ok) throw new Error('Failed to create order');
    return response.json();
  },

  async addBot(id: string) {
    const response = await fetch(`${API_BASE_URL}/bots`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id }),
    });
    if (!response.ok) throw new Error('Failed to add bot');
    return response.json();
  },

  async removeBot(id: string) {
    const response = await fetch(`${API_BASE_URL}/bots`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id }),
    });
    if (!response.ok) throw new Error('Failed to remove bot');
    return response.json();
  },

  async getBots() {
    const response = await fetch(`${API_BASE_URL}/bots`);
    if (!response.ok) throw new Error('Failed to get bots');
    return response.json();
  },

  async getOrders() {
    const response = await fetch(`${API_BASE_URL}/orders`);
    if (!response.ok) throw new Error('Failed to get orders');
    return response.json();
  },

  async getPendingOrders() {
    const response = await fetch(`${API_BASE_URL}/orders/pending`);
    if (!response.ok) throw new Error('Failed to get pending orders');
    return response.json();
  },

  async getCompletedOrders() {
    const response = await fetch(`${API_BASE_URL}/orders/completed`);
    if (!response.ok) throw new Error('Failed to get completed orders');
    return response.json();
  },

  async getStats() {
    const response = await fetch(`${API_BASE_URL}/stats`);
    if (!response.ok) throw new Error('Failed to get stats');
    return response.json();
  },

  async reset() {
    const response = await fetch(`${API_BASE_URL}/reset`, {
      method: 'POST',
    });
    if (!response.ok) throw new Error('Failed to reset');
    return response.json();
  },
};
