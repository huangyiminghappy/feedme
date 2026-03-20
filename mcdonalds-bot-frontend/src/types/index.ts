export type OrderType = 'Normal' | 'VIP';
export type OrderStatus = 'Pending' | 'Processing' | 'Complete';
export type BotStatus = 'Idle' | 'Processing';

export interface Order {
  id: string;
  type: OrderType;
  status: OrderStatus;
  bot_id?: string;
  created_at: string;
  started_at?: string;
  completed_at?: string;
}

export interface Bot {
  id: string;
  status: BotStatus;
  current_order?: Order;
  processing_time: number;
  completed_orders: number;
}

export interface Stats {
  total_bots: number;
  idle_bots: number;
  processing_bots: number;
  pending_orders: number;
  completed_orders: number;
  total_orders: number;
}

export type EventType = 
  | 'order_created'
  | 'order_processing'
  | 'order_complete'
  | 'bot_added'
  | 'bot_removed'
  | 'queue_updated';

export interface Event {
  type: EventType;
  order?: Order;
  bot?: Bot;
  timestamp: string;
  message?: string;
}
