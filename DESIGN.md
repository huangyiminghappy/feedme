# McDonald's 自动化订单管理系统 - 完整技术设计文档

> **文档版本**: v1.0  
> **创建日期**: 2026-03-19  
> **目标**: 为不熟悉Go和前端的开发者提供完整的实现和测试指南

---

## 📋 目录

1. [项目概述](#1-项目概述)
2. [GitHub需求清单](#2-github需求清单)
3. [详细测试用例](#3-详细测试用例)
4. [架构设计](#4-架构设计)
5. [核心模块设计](#5-核心模块设计)
6. [事件系统设计](#6-事件系统设计)
7. [REST API设计](#7-rest-api设计)
8. [前端设计](#8-前端设计)
9. [如何运行和测试（零基础指南）](#9-如何运行和测试零基础指南) ⭐ 最重要
10. [AI集成方案](#10-ai集成方案)
11. [AI扩展最终效果展示](#11-ai扩展最终效果展示) ⭐ 新增
12. [完整项目结构](#12-完整项目结构)
13. [脚本文件内容](#13-脚本文件内容)
14. [GitHub Actions配置](#14-github-actions配置)
15. [面试准备](#15-面试准备)
16. [常见问题FAQ](#16-常见问题faq)
17. [最后检查清单](#17-最后检查清单)

---

## 快速导航

- 🚀 **第一次阅读**: 从第1章开始，重点看第9章（测试指南）
- 🔧 **准备写代码**: 重点看第4-8章（架构和API设计）
- 🤖 **准备AI面试**: 重点看第10-11章（AI集成和效果展示）
- 📝 **准备面试**: 重点看第15章（面试准备）
- ❓ **遇到问题**: 查看第16章（FAQ）

---

# Part 1: 项目概述和测试用例

## 1. 项目概述

### 1.1 业务背景
麦当劳COVID-19期间的数字化转型项目，构建自动化烹饪机器人订单管理系统。

### 1.2 核心功能
- ✅ 创建普通订单和VIP订单
- ✅ VIP订单优先处理（但排在已有VIP订单之后）
- ✅ 动态增减机器人
- ✅ 机器人自动处理订单（10秒/订单）
- ✅ 实时状态更新（PENDING → PROCESSING → COMPLETE）
- ✅ 机器人移除时订单归位到原位置

### 1.3 技术选型
- **后端**: Go 1.21+ (Gin框架)
- **前端**: React 18 + TypeScript + Vite
- **通信**: WebSocket (实时) + REST API
- **测试**: Go testing + GitHub Actions

## 2. GitHub需求清单

| 需求ID | 需求描述 | 测试用例 | 验证方法 |
|--------|---------|---------|---------|
| R1 | 点击"New Normal Order"创建普通订单 | TC1.1, TC1.2 | 订单出现在PENDING区域 |
| R2 | 点击"New VIP Order"创建VIP订单 | TC2.1, TC2.2, TC2.3 | VIP订单排在普通订单前，但在已有VIP后 |
| R3 | 订单号唯一且递增 | TC3.1 | 验证订单号序列 |
| R4 | 点击"+ Bot"创建机器人并处理订单 | TC4.1, TC4.2, TC4.3 | 10秒后订单移到COMPLETE |
| R5 | 无订单时机器人IDLE | TC5.1 | 验证机器人状态 |
| R6 | 点击"- Bot"移除最新机器人 | TC6.1, TC6.2, TC6.3 | 处理中订单归位到原位置 |
| R7 | 内存存储，无需持久化 | TC7.1 | 重启后数据清空 |
| R8 | CLI输出到result.txt | TC8.1, TC8.2 | 包含HH:MM:SS时间戳 |

## 3. 详细测试用例

### TC1: 普通订单创建
**TC1.1 - 基本创建**
- 输入: 点击"New Normal Order"
- 预期: 订单出现在PENDING区域，类型为NORMAL，订单号递增，状态为PENDING

**TC1.2 - 批量创建**
- 输入: 连续点击5次"New Normal Order"
- 预期: 5个订单都在PENDING区域，订单号连续递增（1,2,3,4,5），顺序为FIFO

### TC2: VIP订单优先级（最关键）
**TC2.1 - VIP插队普通订单**
- 步骤: 创建NORMAL#1 → NORMAL#2 → VIP#3
- 预期队列: #3(VIP), #1(NORMAL), #2(NORMAL)

**TC2.2 - VIP排在已有VIP后**
- 步骤: 创建VIP#1 → VIP#2 → NORMAL#3 → VIP#4
- 预期队列: #1(VIP), #2(VIP), #4(VIP), #3(NORMAL)

**TC2.3 - 复杂混合场景**
- 步骤: NORMAL#1 → VIP#2 → NORMAL#3 → VIP#4 → NORMAL#5 → VIP#6
- 预期队列: #2(VIP), #4(VIP), #6(VIP), #1(NORMAL), #3(NORMAL), #5(NORMAL)

### TC3: 订单号唯一性
**TC3.1 - 订单号递增**
- 步骤: 创建10个订单（混合VIP和NORMAL）
- 预期: 订单号从1递增到10，无重复，无跳号

### TC4: 机器人处理订单
**TC4.1 - 单机器人处理单订单**
- 步骤: 创建订单#1 → 添加机器人#1 → 等待10秒
- 预期: 订单立即PENDING→PROCESSING，10秒后→COMPLETE，机器人IDLE

**TC4.2 - 单机器人处理多订单**
- 步骤: 创建订单#1,#2,#3 → 添加机器人#1 → 等待30秒
- 预期: 订单按顺序处理，每个约10秒，总时间约30秒

**TC4.3 - 多机器人并发处理**
- 步骤: 创建订单#1,#2,#3 → 添加机器人#1,#2,#3 → 等待10秒
- 预期: 3个订单同时处理，10秒后全部完成，无重复处理

### TC5: 机器人IDLE状态
**TC5.1 - 无订单时IDLE**
- 步骤: 添加机器人#1，不创建订单
- 预期: 机器人状态为IDLE

**TC5.2 - 处理完所有订单后IDLE**
- 步骤: 创建订单#1 → 添加机器人#1 → 等待11秒
- 预期: 订单完成后机器人变为IDLE

### TC6: 机器人移除（最关键测试）
**TC6.1 - 移除IDLE机器人**
- 步骤: 添加机器人#1 → 点击"- Bot"
- 预期: 机器人#1被移除，无副作用

**TC6.2 - 移除PROCESSING机器人（核心）**
- 步骤: 创建VIP#1 → 创建NORMAL#2 → 添加机器人#1 → 等待5秒 → 点击"- Bot"
- 预期: 
  * 机器人#1被移除
  * 订单#1回到PENDING队列
  * 订单#1仍在订单#2前面（保持VIP优先级）
  * 订单#1状态变回PENDING

**TC6.3 - 移除最新机器人**
- 步骤: 添加机器人#1,#2,#3 → 点击"- Bot"
- 预期: 机器人#3被移除（最新的），#1和#2保留

### TC7: 内存存储
**TC7.1 - 重启清空数据**
- 步骤: 创建订单并处理 → 停止程序 → 重新启动
- 预期: 所有数据清空，订单号从1重新开始

### TC8: CLI输出格式
**TC8.1 - result.txt格式验证**
```
10:30:01 | ORDER_CREATED   | Order #1 (NORMAL) created
10:30:02 | ORDER_CREATED   | Order #2 (VIP) created
10:30:03 | BOT_ADDED       | Bot #1 added
10:30:03 | ORDER_STARTED   | Order #2 (VIP) started by Bot #1
10:30:13 | ORDER_COMPLETED | Order #2 (VIP) completed (10.00s)
```
验证点: 时间格式HH:MM:SS，包含所有关键事件，处理时间准确

**TC8.2 - GitHub Actions兼容性**
- 预期: test.sh成功，build.sh编译成功，run.sh生成result.txt，所有检查通过
# Part 2: 架构设计和核心模块

## 4. 架构设计

### 4.1 整体架构图
```
Frontend (React + WebSocket)
         ↓
    REST API + WebSocket
         ↓
Backend (Go + Gin)
  ├── HTTP/WebSocket Layer
  ├── Business Logic Layer
  │   ├── OrderManager
  │   ├── BotManager
  │   └── EventPublisher
  └── Data Layer
      ├── PriorityQueue
      ├── BotPool
      └── EventStore
```

### 4.2 数据流
```
用户操作 → Frontend → REST API → Business Logic → Data Layer
                                        ↓
                                   Event Published
                                        ↓
                                   WebSocket Hub
                                        ↓
                                   All Clients
```

## 5. 核心模块设计

### 5.1 订单数据结构
```go
type OrderType string
const (
    OrderTypeNormal OrderType = "NORMAL"
    OrderTypeVIP    OrderType = "VIP"
)

type OrderStatus string
const (
    OrderStatusPending    OrderStatus = "PENDING"
    OrderStatusProcessing OrderStatus = "PROCESSING"
    OrderStatusComplete   OrderStatus = "COMPLETE"
)

type Order struct {
    ID          string      `json:"id"`
    OrderNumber int         `json:"orderNumber"`
    Type        OrderType   `json:"type"`
    Status      OrderStatus `json:"status"`
    CreatedAt   time.Time   `json:"createdAt"`
    StartedAt   *time.Time  `json:"startedAt"`
    CompletedAt *time.Time  `json:"completedAt"`
    BotID       string      `json:"botId"`
    Priority    int         `json:"-"`
}
```

### 5.2 优先级计算规则（关键）
```go
// VIP订单: priority = 1000000 - orderNumber
// 例如: VIP #1 = 999999, VIP #2 = 999998
// 这样VIP订单总是排在前面，且按创建顺序排列

// NORMAL订单: priority = orderNumber
// 例如: NORMAL #1 = 1, NORMAL #2 = 2
// 这样NORMAL订单排在VIP后面，且按创建顺序排列

// 优先级队列使用最小堆，priority越小越优先
```

**为什么这样设计？**
- VIP订单的priority总是小于NORMAL订单（999999 < 1）
- 同类型订单按创建顺序排列（先创建的priority更小）
- 满足"VIP优先，但排在已有VIP后"的需求

### 5.3 优先级队列实现
```go
type PriorityQueue struct {
    mu     sync.RWMutex
    items  []*Order
    lookup map[string]*Order
}

// 核心方法
func (pq *PriorityQueue) Push(order *Order)      // O(log n)
func (pq *PriorityQueue) Pop() *Order            // O(log n)
func (pq *PriorityQueue) Peek() *Order           // O(1)
func (pq *PriorityQueue) Remove(orderID string) *Order  // O(log n)
func (pq *PriorityQueue) Len() int               // O(1)
```

**为什么用堆？**
- Push/Pop/Remove都是O(log n)，性能优秀
- 适合频繁插入和删除的场景
- Go标准库container/heap提供支持

### 5.4 机器人数据结构
```go
type BotStatus string
const (
    BotStatusIdle       BotStatus = "IDLE"
    BotStatusProcessing BotStatus = "PROCESSING"
)

type Bot struct {
    ID            string      `json:"id"`
    Status        BotStatus   `json:"status"`
    CurrentOrder  *Order      `json:"currentOrder"`
    CreatedAt     time.Time   `json:"createdAt"`
    ProcessingCtx context.Context
    CancelFunc    context.CancelFunc
}
```

### 5.5 机器人工作流程（关键）
```
1. Bot创建 → 立即尝试获取订单
2. 从优先级队列Pop订单 → 更新订单状态为PROCESSING
3. 启动goroutine处理订单（10秒）
4. 处理完成 → 更新订单状态为COMPLETE
5. 继续获取下一个订单（如果有）
6. 无订单时 → 进入IDLE状态，等待新订单通知
```

### 5.6 机器人取消处理逻辑（最关键）
```go
// 当机器人被移除时：
1. 调用CancelFunc()取消context
2. goroutine检测到context.Done()
3. 停止处理，将订单状态改回PENDING
4. 将订单重新Push回优先级队列（保持原priority）
5. 移除机器人
```

**为什么用context？**
- 优雅地取消goroutine
- 避免goroutine泄漏
- 支持超时控制

### 5.7 OrderManager接口
```go
type OrderManager struct {
    mu              sync.RWMutex
    pendingQueue    *PriorityQueue
    processingOrders map[string]*Order
    completedOrders []*Order
    orderCounter    int
    eventPublisher  *EventPublisher
}

// 核心方法
func (om *OrderManager) CreateOrder(orderType OrderType) *Order
func (om *OrderManager) GetPendingOrders() []*Order
func (om *OrderManager) GetProcessingOrders() []*Order
func (om *OrderManager) GetCompletedOrders() []*Order
func (om *OrderManager) StartProcessing(orderID string, botID string) error
func (om *OrderManager) CompleteOrder(orderID string) error
func (om *OrderManager) ReturnToPending(orderID string) error
```

### 5.8 BotManager接口
```go
type BotManager struct {
    mu              sync.RWMutex
    bots            map[string]*Bot
    botList         []*Bot
    orderManager    *OrderManager
    eventPublisher  *EventPublisher
    processingTime  time.Duration
    orderAvailable  chan struct{}
}

// 核心方法
func (bm *BotManager) AddBot() *Bot
func (bm *BotManager) RemoveBot() error
func (bm *BotManager) GetBots() []*Bot
func (bm *BotManager) processOrder(bot *Bot)
func (bm *BotManager) notifyOrderAvailable()
```

### 5.9 并发控制策略

**锁的使用：**
```go
// OrderManager
- pendingQueue: 内部有RWMutex
- processingOrders: 外层RWMutex保护
- completedOrders: 外层RWMutex保护

// BotManager
- bots map: RWMutex保护
- botList: RWMutex保护

// 锁的顺序（避免死锁）：
1. 先锁BotManager
2. 再锁OrderManager
3. 最后锁PriorityQueue
```

**Channel的使用：**
```go
// orderAvailable chan struct{}
- 当新订单创建时，发送信号
- IDLE的机器人监听此channel
- 收到信号后尝试获取订单

// WebSocket Hub
- register chan *Client
- unregister chan *Client
- broadcast chan []byte
```

## 6. 事件系统设计

### 6.1 事件类型
```go
type EventType string
const (
    EventOrderCreated      EventType = "ORDER_CREATED"
    EventOrderStarted      EventType = "ORDER_STARTED"
    EventOrderCompleted    EventType = "ORDER_COMPLETED"
    EventOrderReturned     EventType = "ORDER_RETURNED"
    EventBotAdded          EventType = "BOT_ADDED"
    EventBotRemoved        EventType = "BOT_REMOVED"
    EventBotStatusChanged  EventType = "BOT_STATUS_CHANGED"
)

type Event struct {
    Type      EventType   `json:"type"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data"`
}
```

### 6.2 EventPublisher
```go
type EventPublisher struct {
    mu          sync.RWMutex
    subscribers map[string]chan Event
    eventLog    []Event
}

func (ep *EventPublisher) Subscribe(clientID string) chan Event
func (ep *EventPublisher) Unsubscribe(clientID string)
func (ep *EventPublisher) Publish(event Event)
func (ep *EventPublisher) GetEventLog() []Event
```
# Part 3: API设计和前端设计

## 7. REST API设计

### 7.1 订单相关API
```
POST   /api/orders/normal
响应: { "id": "uuid", "orderNumber": 1, "type": "NORMAL", "status": "PENDING", ... }

POST   /api/orders/vip
响应: { "id": "uuid", "orderNumber": 2, "type": "VIP", "status": "PENDING", ... }

GET    /api/orders
响应: {
  "pending": [...],
  "processing": [...],
  "completed": [...]
}

GET    /api/state
响应: {
  "orders": { "pending": [...], "processing": [...], "completed": [...] },
  "bots": [...]
}
```

### 7.2 机器人相关API
```
POST   /api/bots
响应: { "id": "uuid", "status": "IDLE", "createdAt": "..." }

DELETE /api/bots
响应: { "message": "Bot removed", "botId": "uuid" }

GET    /api/bots
响应: [{ "id": "uuid", "status": "IDLE", ... }]
```

### 7.3 WebSocket API
```
WS     /ws

客户端连接后自动接收事件：
{
  "type": "ORDER_CREATED",
  "timestamp": "2026-03-19T10:30:01Z",
  "data": { "order": {...} }
}

{
  "type": "STATE_UPDATE",
  "timestamp": "2026-03-19T10:30:01Z",
  "data": {
    "orders": { "pending": [...], "processing": [...], "completed": [...] },
    "bots": [...]
  }
}
```

## 8. 前端设计

### 8.1 组件结构
```
src/
├── components/
│   ├── ControlPanel.tsx        # 控制按钮
│   ├── OrderBoard.tsx          # 订单面板
│   ├── OrderCard.tsx           # 订单卡片
│   ├── BotPanel.tsx            # 机器人面板
│   ├── BotCard.tsx             # 机器人卡片
│   └── StatsPanel.tsx          # 统计面板
├── hooks/
│   ├── useWebSocket.ts         # WebSocket管理
│   └── useStore.ts             # 状态管理
├── services/
│   └── api.ts                  # API封装
├── types/
│   └── index.ts                # TypeScript类型
└── App.tsx
```

### 8.2 状态管理（Zustand）
```typescript
interface AppState {
  pendingOrders: Order[];
  processingOrders: Order[];
  completedOrders: Order[];
  bots: Bot[];
  connected: boolean;
  
  createOrder: (type: 'NORMAL' | 'VIP') => Promise<void>;
  addBot: () => Promise<void>;
  removeBot: () => Promise<void>;
  updateFromWebSocket: (data: any) => void;
}
```

### 8.3 UI布局
```
┌─────────────────────────────────────────────────────────┐
│                    Control Panel                         │
│  [New Normal Order] [New VIP Order] [+ Bot] [- Bot]     │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│  Stats: Total: 10 | Pending: 3 | Processing: 2 | ...    │
└─────────────────────────────────────────────────────────┘
┌──────────────┬──────────────┬──────────────┬────────────┐
│   PENDING    │  PROCESSING  │   COMPLETE   │    BOTS    │
│              │              │              │            │
│  Order #2    │  Order #1    │  Order #5    │  Bot #1    │
│  VIP         │  NORMAL      │  VIP         │  BUSY      │
│  10:30:01    │  Bot: #1     │  10:30:45    │  Order #1  │
│              │  Started:    │  Duration:   │            │
│              │  10:30:03    │  10.00s      │            │
│              │              │              │            │
│  Order #3    │  Order #4    │  Order #6    │  Bot #2    │
│  VIP         │  VIP         │  NORMAL      │  IDLE      │
│  10:30:02    │  Bot: #2     │  10:30:50    │            │
│              │  Started:    │  Duration:   │            │
│              │  10:30:03    │  10.01s      │            │
└──────────────┴──────────────┴──────────────┴────────────┘
```

### 8.4 WebSocket连接管理
```typescript
const useWebSocket = () => {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  
  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080/ws');
    
    socket.onopen = () => {
      setConnected(true);
    };
    
    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      // 更新状态
      useStore.getState().updateFromWebSocket(data);
    };
    
    socket.onclose = () => {
      setConnected(false);
      // 重连逻辑
      setTimeout(() => {
        // 重新连接
      }, 3000);
    };
    
    setWs(socket);
    
    return () => {
      socket.close();
    };
  }, []);
  
  return { ws, connected };
};
```

### 8.5 订单卡片设计
```typescript
interface OrderCardProps {
  order: Order;
}

const OrderCard: React.FC<OrderCardProps> = ({ order }) => {
  return (
    <div className={`
      p-4 rounded-lg border-2
      ${order.type === 'VIP' ? 'border-yellow-500 bg-yellow-50' : 'border-gray-300 bg-white'}
    `}>
      <div className="flex justify-between items-center">
        <span className="text-lg font-bold">Order #{order.orderNumber}</span>
        <span className={`
          px-2 py-1 rounded text-xs font-semibold
          ${order.type === 'VIP' ? 'bg-yellow-500 text-white' : 'bg-gray-500 text-white'}
        `}>
          {order.type}
        </span>
      </div>
      <div className="mt-2 text-sm text-gray-600">
        <div>Created: {formatTime(order.createdAt)}</div>
        {order.startedAt && <div>Started: {formatTime(order.startedAt)}</div>}
        {order.completedAt && <div>Completed: {formatTime(order.completedAt)}</div>}
        {order.botId && <div>Bot: {order.botId}</div>}
      </div>
    </div>
  );
};
```

### 8.6 机器人卡片设计
```typescript
interface BotCardProps {
  bot: Bot;
}

const BotCard: React.FC<BotCardProps> = ({ bot }) => {
  return (
    <div className={`
      p-4 rounded-lg border-2
      ${bot.status === 'PROCESSING' ? 'border-green-500 bg-green-50' : 'border-gray-300 bg-white'}
    `}>
      <div className="flex justify-between items-center">
        <span className="text-lg font-bold">Bot #{bot.id.slice(0, 8)}</span>
        <span className={`
          px-2 py-1 rounded text-xs font-semibold
          ${bot.status === 'PROCESSING' ? 'bg-green-500 text-white' : 'bg-gray-500 text-white'}
        `}>
          {bot.status}
        </span>
      </div>
      {bot.currentOrder && (
        <div className="mt-2 text-sm text-gray-600">
          Processing Order #{bot.currentOrder.orderNumber}
        </div>
      )}
    </div>
  );
};
```
# Part 4: 测试指南（零基础）

## 9. 如何运行和测试（完全不懂Go和前端也能操作）

### 9.1 前置准备

**安装Go（如果没有）：**
```bash
# macOS
brew install go

# 验证安装
go version  # 应该显示 go version go1.21.x
```

**安装Node.js（如果没有）：**
```bash
# macOS
brew install node

# 验证安装
node --version  # 应该显示 v18.x 或更高
npm --version   # 应该显示 9.x 或更高
```

### 9.2 项目结构
```
feedme/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/...
│   ├── go.mod
│   ├── go.sum
│   └── scripts/
│       ├── build.sh
│       ├── test.sh
│       └── run.sh
├── frontend/
│   ├── src/...
│   ├── package.json
│   └── vite.config.ts
└── result.txt  # CLI运行后生成
```

### 9.3 后端测试步骤（一步一步来）

**步骤1: 进入后端目录**
```bash
cd /Users/huangyiming/Documents/workspace/feedme/backend
```

**步骤2: 安装依赖**
```bash
go mod download
```
预期输出: 下载各种依赖包，没有错误信息

**步骤3: 运行单元测试**
```bash
chmod +x scripts/test.sh
./scripts/test.sh
```
预期输出:
```
=== RUN   TestPriorityQueue
--- PASS: TestPriorityQueue (0.00s)
=== RUN   TestOrderManager
--- PASS: TestOrderManager (0.00s)
=== RUN   TestBotManager
--- PASS: TestBotManager (10.05s)
...
PASS
ok      github.com/feedme/backend/internal/service      10.123s
```

**如果测试失败，会显示：**
```
--- FAIL: TestBotManager (5.00s)
    bot_manager_test.go:45: Expected order to be completed, but got PROCESSING
FAIL
```

**步骤4: 编译程序**
```bash
chmod +x scripts/build.sh
./scripts/build.sh
```
预期输出:
```
Building...
Build successful: ./bin/mcdonalds-bot
```

**步骤5: 运行CLI模拟（生成result.txt）**
```bash
chmod +x scripts/run.sh
./scripts/run.sh
```
预期输出:
```
Running simulation...
Simulation completed. Check result.txt
```

**步骤6: 检查result.txt**
```bash
cat result.txt
```
预期内容:
```
=== McDonald's Order Management System ===
Start Time: 10:30:00

10:30:01 | ORDER_CREATED   | Order #1 (NORMAL) created
10:30:02 | ORDER_CREATED   | Order #2 (VIP) created
10:30:03 | BOT_ADDED       | Bot #1 added
10:30:03 | ORDER_STARTED   | Order #2 (VIP) started by Bot #1
10:30:13 | ORDER_COMPLETED | Order #2 (VIP) completed (10.00s)
10:30:13 | ORDER_STARTED   | Order #1 (NORMAL) started by Bot #1
10:30:23 | ORDER_COMPLETED | Order #1 (NORMAL) completed (10.00s)

=== Statistics ===
Total Orders: 2
Completed Orders: 2
Average Processing Time: 10.00s
```

**步骤7: 运行服务器模式（用于前端连接）**
```bash
./bin/mcdonalds-bot run
```
预期输出:
```
[GIN-debug] Listening and serving HTTP on :8080
```

### 9.4 前端测试步骤

**步骤1: 打开新终端，进入前端目录**
```bash
cd /Users/huangyiming/Documents/workspace/feedme/frontend
```

**步骤2: 安装依赖**
```bash
npm install
```
预期输出: 安装各种npm包，最后显示 `added XXX packages`

**步骤3: 运行开发服务器**
```bash
npm run dev
```
预期输出:
```
  VITE v5.0.0  ready in 500 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

**步骤4: 打开浏览器**
```
在浏览器中打开: http://localhost:5173/
```

**步骤5: 手动测试（按照测试用例）**

**测试TC2.1 - VIP插队普通订单：**
1. 点击 "New Normal Order" 按钮 → 看到 Order #1 (NORMAL) 出现在PENDING区域
2. 点击 "New Normal Order" 按钮 → 看到 Order #2 (NORMAL) 出现在PENDING区域
3. 点击 "New VIP Order" 按钮 → 看到 Order #3 (VIP) 出现在PENDING区域
4. **验证**: PENDING区域的顺序应该是 #3(VIP), #1(NORMAL), #2(NORMAL)

**测试TC4.1 - 单机器人处理单订单：**
1. 点击 "New Normal Order" → 看到 Order #1 在PENDING
2. 点击 "+ Bot" → 看到 Bot #1 出现，Order #1 立即移到PROCESSING
3. 等待10秒 → Order #1 移到COMPLETE，Bot #1 状态变为IDLE
4. **验证**: 完成时间约10秒（可以看时间戳）

**测试TC6.2 - 移除PROCESSING机器人（最关键）：**
1. 点击 "New VIP Order" → Order #1 (VIP) 在PENDING
2. 点击 "New Normal Order" → Order #2 (NORMAL) 在PENDING
3. 点击 "+ Bot" → Bot #1 开始处理 Order #1
4. 等待5秒（处理到一半）
5. 点击 "- Bot" → Bot #1 被移除
6. **验证**: 
   - Order #1 回到PENDING区域
   - Order #1 仍在 Order #2 前面
   - Order #1 状态变回PENDING

### 9.5 自动化测试脚本

**创建测试脚本 test_all.sh：**
```bash
#!/bin/bash

echo "=== Running Backend Tests ==="
cd backend
./scripts/test.sh
if [ $? -ne 0 ]; then
    echo "❌ Backend tests failed"
    exit 1
fi
echo "✅ Backend tests passed"

echo ""
echo "=== Building Backend ==="
./scripts/build.sh
if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi
echo "✅ Build successful"

echo ""
echo "=== Running CLI Simulation ==="
./scripts/run.sh
if [ $? -ne 0 ]; then
    echo "❌ Simulation failed"
    exit 1
fi
echo "✅ Simulation successful"

echo ""
echo "=== Checking result.txt ==="
if [ ! -f "result.txt" ]; then
    echo "❌ result.txt not found"
    exit 1
fi

# 验证result.txt格式
if grep -q "ORDER_CREATED" result.txt && \
   grep -q "BOT_ADDED" result.txt && \
   grep -q "ORDER_COMPLETED" result.txt; then
    echo "✅ result.txt format correct"
else
    echo "❌ result.txt format incorrect"
    exit 1
fi

echo ""
echo "=== All Tests Passed ==="
```

**运行完整测试：**
```bash
chmod +x test_all.sh
./test_all.sh
```

### 9.6 GitHub Actions测试

**当你提交PR后，GitHub Actions会自动运行：**
1. 运行 `scripts/test.sh`
2. 运行 `scripts/build.sh`
3. 运行 `scripts/run.sh`
4. 验证 `result.txt` 存在且格式正确

**如何查看GitHub Actions结果：**
1. 进入你的GitHub仓库
2. 点击 "Actions" 标签
3. 查看最新的workflow运行
4. 如果有红色❌，点击查看错误日志
5. 如果全是绿色✅，说明通过

### 9.7 常见测试问题排查

**问题1: 测试超时**
```
--- FAIL: TestBotProcessing (30.00s)
    timeout
```
原因: 机器人处理时间设置错误
解决: 检查 `processingTime` 是否为10秒

**问题2: 订单顺序错误**
```
Expected: [VIP#1, VIP#2, NORMAL#3]
Got:      [VIP#1, NORMAL#3, VIP#2]
```
原因: 优先级计算错误
解决: 检查 `priority` 计算逻辑

**问题3: 机器人移除后订单丢失**
```
Expected order to return to PENDING
Got: order not found
```
原因: 取消处理逻辑有bug
解决: 检查 `ReturnToPending` 方法

**问题4: WebSocket连接失败**
```
WebSocket connection failed
```
原因: 后端未启动或端口被占用
解决: 
```bash
# 检查后端是否运行
ps aux | grep mcdonalds-bot

# 检查端口是否被占用
lsof -i :8080

# 杀死占用端口的进程
kill -9 <PID>
```

**问题5: 前端无法连接后端**
```
Failed to fetch
```
原因: CORS配置或后端未启动
解决: 检查后端Gin的CORS中间件配置

### 9.8 测试检查清单

**提交PR前必须检查：**
- [ ] 所有单元测试通过 (`./scripts/test.sh`)
- [ ] 编译成功 (`./scripts/build.sh`)
- [ ] CLI模拟成功 (`./scripts/run.sh`)
- [ ] result.txt格式正确
- [ ] 手动测试TC2.1通过（VIP优先级）
- [ ] 手动测试TC4.1通过（10秒处理）
- [ ] 手动测试TC6.2通过（机器人移除）
- [ ] 前端UI正常显示
- [ ] WebSocket实时更新正常
- [ ] 无console错误

### 9.9 性能测试

**压力测试脚本：**
```bash
# 创建100个订单，10个机器人
for i in {1..100}; do
    curl -X POST http://localhost:8080/api/orders/normal
done

for i in {1..10}; do
    curl -X POST http://localhost:8080/api/bots
done

# 观察处理情况
watch -n 1 'curl -s http://localhost:8080/api/state | jq'
```

**预期结果：**
- 100个订单在约100秒内全部完成（10个机器人并发）
- 无订单丢失
- 无机器人死锁
- 内存占用稳定
# Part 5: AI扩展设计

## 10. AI集成方案

### 10.1 为什么需要AI？

在多国家、多餐厅的SaaS场景下，AI可以解决：
1. **智能调度**: 根据订单类型、复杂度智能分配机器人
2. **需求预测**: 预测高峰期，提前调整机器人数量
3. **异常检测**: 识别异常订单模式、机器人故障
4. **优化建议**: 给餐厅经理提供运营优化建议

### 10.2 AI扩展架构

```
┌─────────────────────────────────────────────────────────┐
│                   AI Service Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Predictor    │  │  Scheduler   │  │  Detector    │  │
│  │ Service      │  │  Service     │  │  Service     │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────┴────────────────────────────────────┐
│              Business Logic Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │OrderManager  │  │ BotManager   │  │EventPublisher│  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 10.3 智能调度服务 (AI Scheduler)

**场景**: 不同订单可能有不同复杂度，AI可以智能分配

**数据结构扩展：**
```go
type Order struct {
    // ... 原有字段
    Complexity    int       `json:"complexity"`    // 1-10，AI预测
    EstimatedTime int       `json:"estimatedTime"` // 预计处理时间（秒）
    Items         []string  `json:"items"`         // 订单项目
}

type Bot struct {
    // ... 原有字段
    Efficiency    float64   `json:"efficiency"`    // 0.8-1.2，AI评估
    Specialties   []string  `json:"specialties"`   // 擅长的订单类型
    SuccessRate   float64   `json:"successRate"`   // 成功率
}
```

**AI调度接口：**
```go
type AIScheduler interface {
    // 预测订单复杂度
    PredictComplexity(order *Order) int
    
    // 为订单选择最佳机器人
    SelectBestBot(order *Order, availableBots []*Bot) *Bot
    
    // 预测订单处理时间
    EstimateProcessingTime(order *Order, bot *Bot) time.Duration
    
    // 优化队列顺序（考虑复杂度）
    OptimizeQueue(orders []*Order) []*Order
}
```

**实现示例（简化版）：**
```go
type SimpleAIScheduler struct {
    model *MLModel  // 机器学习模型
}

func (s *SimpleAIScheduler) SelectBestBot(order *Order, bots []*Bot) *Bot {
    var bestBot *Bot
    var bestScore float64
    
    for _, bot := range bots {
        if bot.Status != BotStatusIdle {
            continue
        }
        
        // 计算匹配分数
        score := s.calculateMatchScore(order, bot)
        
        if score > bestScore {
            bestScore = score
            bestBot = bot
        }
    }
    
    return bestBot
}

func (s *SimpleAIScheduler) calculateMatchScore(order *Order, bot *Bot) float64 {
    score := 0.0
    
    // 效率分数
    score += bot.Efficiency * 0.4
    
    // 成功率分数
    score += bot.SuccessRate * 0.3
    
    // 专长匹配分数
    for _, item := range order.Items {
        for _, specialty := range bot.Specialties {
            if item == specialty {
                score += 0.3
                break
            }
        }
    }
    
    return score
}
```

**如何训练模型：**
```python
# Python训练脚本示例
import pandas as pd
from sklearn.ensemble import RandomForestRegressor

# 加载历史数据
data = pd.read_csv('order_history.csv')
# 特征: order_items, order_type, time_of_day, day_of_week
# 标签: actual_processing_time

X = data[['order_items_count', 'is_vip', 'hour', 'day_of_week']]
y = data['processing_time']

model = RandomForestRegressor()
model.fit(X, y)

# 保存模型
import joblib
joblib.dump(model, 'processing_time_model.pkl')
```

**Go中加载模型：**
```go
// 使用 cgo 调用 Python 模型
// 或者使用 ONNX Runtime for Go
// 或者使用简单的规则引擎

type RuleBasedScheduler struct{}

func (r *RuleBasedScheduler) PredictComplexity(order *Order) int {
    complexity := 1
    
    // 规则1: VIP订单通常更复杂
    if order.Type == OrderTypeVIP {
        complexity += 2
    }
    
    // 规则2: 订单项目越多越复杂
    complexity += len(order.Items) / 2
    
    // 规则3: 高峰期订单更复杂
    hour := order.CreatedAt.Hour()
    if hour >= 11 && hour <= 13 || hour >= 18 && hour <= 20 {
        complexity += 1
    }
    
    if complexity > 10 {
        complexity = 10
    }
    
    return complexity
}
```

### 10.4 需求预测服务 (AI Predictor)

**场景**: 预测未来订单量，自动调整机器人数量

**接口设计：**
```go
type AIPredictor interface {
    // 预测未来N分钟的订单量
    PredictOrderVolume(minutes int) int
    
    // 推荐机器人数量
    RecommendBotCount(currentLoad int, predictedLoad int) int
    
    // 检测是否进入高峰期
    IsRushHour() bool
    
    // 预测下一个高峰期
    PredictNextRushHour() time.Time
}
```

**实现示例：**
```go
type TimeSeriesPredictor struct {
    historicalData []OrderRecord
    mu             sync.RWMutex
}

type OrderRecord struct {
    Timestamp   time.Time
    OrderCount  int
    CompletedCount int
}

func (p *TimeSeriesPredictor) PredictOrderVolume(minutes int) int {
    p.mu.RLock()
    defer p.mu.RUnlock()
    
    // 简单的移动平均预测
    now := time.Now()
    hour := now.Hour()
    dayOfWeek := now.Weekday()
    
    // 查找历史同时段数据
    var historicalAvg float64
    var count int
    
    for _, record := range p.historicalData {
        if record.Timestamp.Hour() == hour && 
           record.Timestamp.Weekday() == dayOfWeek {
            historicalAvg += float64(record.OrderCount)
            count++
        }
    }
    
    if count == 0 {
        return 5  // 默认值
    }
    
    return int(historicalAvg / float64(count))
}

func (p *TimeSeriesPredictor) RecommendBotCount(currentLoad, predictedLoad int) int {
    // 每个机器人每分钟处理6个订单（10秒/订单）
    botsNeeded := (predictedLoad + 5) / 6
    
    // 至少1个，最多20个
    if botsNeeded < 1 {
        botsNeeded = 1
    }
    if botsNeeded > 20 {
        botsNeeded = 20
    }
    
    return botsNeeded
}
```

**自动调整机器人：**
```go
type AutoScaler struct {
    predictor   AIPredictor
    botManager  *BotManager
    enabled     bool
    checkInterval time.Duration
}

func (as *AutoScaler) Start() {
    if !as.enabled {
        return
    }
    
    ticker := time.NewTicker(as.checkInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        as.adjust()
    }
}

func (as *AutoScaler) adjust() {
    currentBots := len(as.botManager.GetBots())
    predictedLoad := as.predictor.PredictOrderVolume(5)
    recommendedBots := as.predictor.RecommendBotCount(currentBots, predictedLoad)
    
    if recommendedBots > currentBots {
        // 增加机器人
        for i := 0; i < recommendedBots - currentBots; i++ {
            as.botManager.AddBot()
        }
        log.Printf("Auto-scaled up: %d -> %d bots", currentBots, recommendedBots)
    } else if recommendedBots < currentBots {
        // 减少机器人
        for i := 0; i < currentBots - recommendedBots; i++ {
            as.botManager.RemoveBot()
        }
        log.Printf("Auto-scaled down: %d -> %d bots", currentBots, recommendedBots)
    }
}
```

### 10.5 异常检测服务 (AI Detector)

**场景**: 检测异常订单、机器人故障、系统问题

**接口设计：**
```go
type AIDetector interface {
    // 检测异常订单
    DetectAnomalousOrder(order *Order) (bool, string)
    
    // 检测机器人故障
    DetectBotFailure(bot *Bot) (bool, string)
    
    // 检测系统异常
    DetectSystemAnomaly(metrics SystemMetrics) (bool, string)
    
    // 推荐修复动作
    RecommendAction(anomaly Anomaly) Action
}

type Anomaly struct {
    Type        string
    Severity    string  // LOW, MEDIUM, HIGH, CRITICAL
    Description string
    Timestamp   time.Time
}

type Action struct {
    Type        string  // RESTART_BOT, CANCEL_ORDER, ALERT_MANAGER
    Description string
    AutoExecute bool
}
```

**实现示例：**
```go
type StatisticalDetector struct {
    orderStats *OrderStatistics
    botStats   *BotStatistics
}

func (d *StatisticalDetector) DetectBotFailure(bot *Bot) (bool, string) {
    // 检测1: 处理时间异常
    if bot.Status == BotStatusProcessing && bot.CurrentOrder != nil {
        elapsed := time.Since(*bot.CurrentOrder.StartedAt)
        if elapsed > 15*time.Second {  // 超过正常时间50%
            return true, "Bot processing time exceeds normal threshold"
        }
    }
    
    // 检测2: 成功率异常
    if bot.SuccessRate < 0.8 {
        return true, "Bot success rate below 80%"
    }
    
    // 检测3: 连续失败
    recentFailures := d.botStats.GetRecentFailures(bot.ID, 5)
    if recentFailures >= 3 {
        return true, "Bot has 3 consecutive failures"
    }
    
    return false, ""
}

func (d *StatisticalDetector) DetectAnomalousOrder(order *Order) (bool, string) {
    // 检测1: 订单项目异常多
    if len(order.Items) > 20 {
        return true, "Order has unusually high number of items"
    }
    
    // 检测2: 创建时间异常（凌晨3点的订单）
    hour := order.CreatedAt.Hour()
    if hour >= 2 && hour <= 5 {
        return true, "Order created during unusual hours"
    }
    
    // 检测3: 重复订单检测
    if d.orderStats.IsDuplicate(order) {
        return true, "Potential duplicate order detected"
    }
    
    return false, ""
}
```

### 10.6 AI Dashboard（面试加分项）

**为餐厅经理提供AI洞察：**
```go
type AIDashboard struct {
    predictor AIPredictor
    scheduler AIScheduler
    detector  AIDetector
}

type AIInsights struct {
    PredictedRushHour    time.Time       `json:"predictedRushHour"`
    RecommendedBots      int             `json:"recommendedBots"`
    CurrentEfficiency    float64         `json:"currentEfficiency"`
    AnomaliesDetected    []Anomaly       `json:"anomaliesDetected"`
    OptimizationSuggestions []string     `json:"optimizationSuggestions"`
}

func (d *AIDashboard) GetInsights() AIInsights {
    return AIInsights{
        PredictedRushHour: d.predictor.PredictNextRushHour(),
        RecommendedBots:   d.predictor.RecommendBotCount(0, 0),
        CurrentEfficiency: d.calculateEfficiency(),
        AnomaliesDetected: d.detector.GetRecentAnomalies(),
        OptimizationSuggestions: d.generateSuggestions(),
    }
}

func (d *AIDashboard) generateSuggestions() []string {
    suggestions := []string{}
    
    // 建议1: 基于预测的机器人调整
    if d.predictor.IsRushHour() {
        suggestions = append(suggestions, 
            "Rush hour detected. Consider adding 2-3 more bots.")
    }
    
    // 建议2: 基于效率的优化
    efficiency := d.calculateEfficiency()
    if efficiency < 0.7 {
        suggestions = append(suggestions, 
            "Bot efficiency is low. Consider restarting underperforming bots.")
    }
    
    // 建议3: 基于异常的警告
    anomalies := d.detector.GetRecentAnomalies()
    if len(anomalies) > 5 {
        suggestions = append(suggestions, 
            "Multiple anomalies detected. System health check recommended.")
    }
    
    return suggestions
}
```

### 10.7 面试时如何展示AI扩展

**第三轮面试可能的问题：**

**Q1: "如果要加入AI优化调度，你会怎么设计？"**
回答要点:
- 展示AIScheduler接口设计
- 说明如何收集训练数据（订单历史、处理时间）
- 解释特征工程（订单复杂度、时间特征、机器人效率）
- 提到可以用简单规则引擎开始，逐步引入ML模型

**Q2: "如何预测高峰期并自动调整机器人？"**
回答要点:
- 展示AIPredictor接口
- 说明时间序列预测方法（移动平均、ARIMA、LSTM）
- 展示AutoScaler自动调整逻辑
- 提到需要考虑成本（机器人有成本）和响应时间

**Q3: "如何检测和处理异常情况？"**
回答要点:
- 展示AIDetector接口
- 说明统计方法（3-sigma、IQR）
- 提到机器学习方法（Isolation Forest、Autoencoder）
- 展示自动修复机制

**Q4: "多国家部署时AI如何适配？"**
回答要点:
- 每个国家/地区独立训练模型
- 考虑时区、文化差异（用餐习惯）
- 使用迁移学习（Transfer Learning）
- 联邦学习（Federated Learning）保护数据隐私

### 10.8 AI扩展的实现优先级

**Phase 1（基础）：**
- ✅ 简单规则引擎（if-else）
- ✅ 基本统计分析
- ✅ 手动阈值设置

**Phase 2（中级）：**
- 🔄 移动平均预测
- 🔄 统计异常检测
- 🔄 简单评分系统

**Phase 3（高级）：**
- 🔜 机器学习模型
- 🔜 时间序列预测
- 🔜 强化学习调度

**Phase 4（专家）：**
- 🔜 深度学习
- 🔜 联邦学习
- 🔜 实时模型更新

**面试建议：**
- 第一轮：实现基础功能，不需要AI
- 第二轮：准备AI扩展的设计思路
- 第三轮：展示AI接口设计，可以用简单规则演示

---

## 11. AI扩展最终效果展示

### 11.1 AI Dashboard UI效果

**前端新增AI洞察面板：**
```
┌─────────────────────────────────────────────────────────────┐
│                    🤖 AI Insights Dashboard                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  📊 Predictions                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Next Rush Hour: 12:00 PM (in 45 minutes)             │  │
│  │ Predicted Orders: 25 orders/hour                     │  │
│  │ Recommended Bots: 5 bots                             │  │
│  │ Current Efficiency: 87.5%                            │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ⚠️ Anomalies Detected (2)                                  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ 🔴 Bot #3: Processing time exceeds threshold         │  │
│  │    Action: Restart recommended                       │  │
│  │                                                       │  │
│  │ 🟡 Order #15: Created during unusual hours           │  │
│  │    Action: Manual review suggested                   │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  💡 Optimization Suggestions                                │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ • Add 2 more bots before lunch rush (11:30 AM)       │  │
│  │ • Bot #2 efficiency is low (65%), consider restart   │  │
│  │ • VIP order processing time +15% today               │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  📈 Real-time Metrics                                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Average Wait Time: 2.3 seconds                        │  │
│  │ Bot Utilization: 78%                                  │  │
│  │ Throughput: 18 orders/min                             │  │
│  │ Success Rate: 98.5%                                   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 11.2 智能调度效果演示

**场景1: 订单复杂度预测**
```
订单 #1 (NORMAL)
  Items: [Big Mac, Fries, Coke]
  AI预测复杂度: 3/10
  预计处理时间: 8秒
  推荐机器人: Bot #2 (汉堡专家)

订单 #2 (VIP)
  Items: [Big Mac x3, McFlurry x2, Happy Meal x2, Salad]
  AI预测复杂度: 8/10
  预计处理时间: 15秒
  推荐机器人: Bot #1 (高效率，成功率99%)
```

**场景2: 智能机器人分配**
```
当前状态:
- Bot #1: 效率 1.2, 成功率 99%, 专长 [汉堡, 套餐]
- Bot #2: 效率 1.0, 成功率 95%, 专长 [甜品, 饮料]
- Bot #3: 效率 0.8, 成功率 92%, 专长 [沙拉, 小食]

新订单: VIP订单 [Big Mac x2, McFlurry]
AI决策: 分配给 Bot #1
理由: 
  - 匹配度最高 (汉堡专长 + 高效率)
  - VIP订单需要最可靠的机器人
  - 预计完成时间: 12秒
```

### 11.3 需求预测效果演示

**时间序列预测图表：**
```
订单量预测 (未来2小时)

30 |                    ╱╲
   |                   ╱  ╲
25 |                  ╱    ╲
   |                 ╱      ╲___
20 |          ___╱╲╱           ╲
   |        ╱╱                   ╲
15 |      ╱╱                      ╲
   |    ╱╱                         ╲
10 |  ╱╱                            ╲
   |╱╱                               ╲
 5 |                                  ╲
   └────────────────────────────────────
   10:00  11:00  12:00  13:00  14:00

当前时间: 10:30
预测高峰: 12:00 (25 orders/hour)
当前机器人: 3
推荐机器人: 5 (11:30前增加2个)
```

**自动扩缩容日志：**
```
10:30:00 | AI_PREDICTION  | Predicted rush hour at 12:00
10:30:00 | AI_SUGGESTION  | Recommend adding 2 bots before 11:30
11:25:00 | AUTO_SCALE_UP  | Adding Bot #4 (predicted load: 20 orders/hour)
11:25:05 | AUTO_SCALE_UP  | Adding Bot #5 (predicted load: 25 orders/hour)
12:15:00 | AI_PREDICTION  | Rush hour ending, load decreasing
12:30:00 | AUTO_SCALE_DOWN| Removing Bot #5 (current load: 15 orders/hour)
13:00:00 | AUTO_SCALE_DOWN| Removing Bot #4 (current load: 10 orders/hour)
```

### 11.4 异常检测效果演示

**实时异常监控：**
```
=== Anomaly Detection Dashboard ===

🔴 CRITICAL (1)
┌────────────────────────────────────────────────────────┐
│ Bot #3 - Processing Timeout                            │
│ Started: 10:30:15                                      │
│ Expected: 10:30:25 (10s)                               │
│ Current: 10:30:32 (17s) ⚠️ 70% overtime               │
│ Action: Auto-cancelling and restarting bot             │
└────────────────────────────────────────────────────────┘

🟡 WARNING (2)
┌────────────────────────────────────────────────────────┐
│ Order #15 - Unusual Creation Time                      │
│ Created: 03:15 AM                                      │
│ Reason: Outside business hours (06:00-23:00)           │
│ Action: Flagged for manual review                      │
└────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────┐
│ Bot #2 - Low Success Rate                              │
│ Recent Success Rate: 75% (below threshold 80%)         │
│ Failed Orders: 3 in last 10                            │
│ Action: Scheduled for maintenance                      │
└────────────────────────────────────────────────────────┘

🟢 INFO (1)
┌────────────────────────────────────────────────────────┐
│ System Performance - Above Average                     │
│ Current Throughput: 22 orders/min (avg: 18)           │
│ Bot Utilization: 85% (optimal range)                   │
└────────────────────────────────────────────────────────┘
```

### 11.5 多国家AI适配效果

**不同地区的AI模型：**
```
┌─────────────────────────────────────────────────────────┐
│                  Global AI Dashboard                     │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  🇺🇸 United States (New York)                           │
│  Peak Hours: 12:00-13:00, 18:00-20:00                   │
│  Popular Items: Big Mac, Fries, Coke                    │
│  Avg Processing: 10.2s                                  │
│  Model Accuracy: 94.5%                                  │
│                                                          │
│  🇯🇵 Japan (Tokyo)                                       │
│  Peak Hours: 11:30-13:30, 19:00-21:00                   │
│  Popular Items: Teriyaki Burger, Green Tea McFlurry    │
│  Avg Processing: 12.5s                                  │
│  Model Accuracy: 91.2%                                  │
│                                                          │
│  🇸🇬 Singapore                                           │
│  Peak Hours: 12:00-14:00, 18:30-20:30                   │
│  Popular Items: McSpicy, Curry Sauce                    │
│  Avg Processing: 11.0s                                  │
│  Model Accuracy: 93.8%                                  │
│                                                          │
│  🇦🇺 Australia (Sydney)                                  │
│  Peak Hours: 12:30-13:30, 18:00-19:30                   │
│  Popular Items: Angus Burger, Flat White                │
│  Avg Processing: 10.8s                                  │
│  Model Accuracy: 92.5%                                  │
└─────────────────────────────────────────────────────────┘
```

### 11.6 AI API响应示例

**GET /api/ai/insights**
```json
{
  "timestamp": "2026-03-19T10:30:00Z",
  "predictions": {
    "nextRushHour": "2026-03-19T12:00:00Z",
    "predictedOrderVolume": 25,
    "recommendedBotCount": 5,
    "confidence": 0.94
  },
  "currentMetrics": {
    "efficiency": 0.875,
    "avgWaitTime": 2.3,
    "botUtilization": 0.78,
    "throughput": 18,
    "successRate": 0.985
  },
  "anomalies": [
    {
      "type": "BOT_TIMEOUT",
      "severity": "CRITICAL",
      "botId": "bot-3",
      "description": "Processing time exceeds threshold",
      "detectedAt": "2026-03-19T10:30:32Z",
      "recommendedAction": "RESTART_BOT"
    },
    {
      "type": "UNUSUAL_ORDER_TIME",
      "severity": "WARNING",
      "orderId": "order-15",
      "description": "Order created during unusual hours",
      "detectedAt": "2026-03-19T03:15:00Z",
      "recommendedAction": "MANUAL_REVIEW"
    }
  ],
  "suggestions": [
    "Add 2 more bots before lunch rush (11:30 AM)",
    "Bot #2 efficiency is low (65%), consider restart",
    "VIP order processing time +15% today"
  ]
}
```

**POST /api/ai/optimize**
```json
{
  "action": "OPTIMIZE_QUEUE",
  "parameters": {
    "considerComplexity": true,
    "considerBotSpecialties": true
  }
}

Response:
{
  "success": true,
  "optimizedQueue": [
    {
      "orderId": "order-5",
      "orderNumber": 5,
      "type": "VIP",
      "complexity": 8,
      "assignedBot": "bot-1",
      "reason": "High complexity VIP order matched with most efficient bot"
    },
    {
      "orderId": "order-3",
      "orderNumber": 3,
      "type": "VIP",
      "complexity": 3,
      "assignedBot": "bot-2",
      "reason": "Simple VIP order, any available bot"
    }
  ],
  "estimatedCompletionTime": "2026-03-19T10:35:00Z",
  "improvementPercent": 15.5
}
```

### 11.7 面试演示脚本

**第三轮面试时的AI演示流程：**

```
面试官: "如果要加入AI优化，你会怎么做？"

你: "我设计了三层AI服务架构，让我演示一下..."

[打开AI Dashboard]

你: "首先是智能调度。看这里，当前有3个订单：
    - 订单#1是简单的汉堡套餐，AI预测复杂度3/10
    - 订单#2是VIP订单，7个项目，AI预测复杂度8/10
    - 系统自动将复杂订单分配给效率最高的Bot #1"

[点击创建复杂VIP订单]

你: "看，AI立即分析了订单内容，预测处理时间15秒，
    并选择了最合适的机器人。这比简单的FIFO队列
    效率提升约15%。"

[切换到预测面板]

你: "第二个功能是需求预测。基于历史数据，AI预测
    12点会有高峰期，建议11:30前增加2个机器人。
    系统可以自动扩缩容，无需人工干预。"

[切换到异常检测]

你: "第三个是异常检测。看这里，Bot #3处理超时了，
    AI自动检测到并建议重启。这个订单会自动归队，
    保持原有优先级。"

面试官: "很好！如果多国家部署呢？"

你: "每个国家独立训练模型。比如日本的高峰期是
    11:30-13:30，美国是12:00-13:00。AI会根据
    当地数据自动调整。我们还可以用迁移学习，
    新市场可以借鉴已有市场的模型。"

面试官: "数据隐私怎么处理？"

你: "可以用联邦学习。每个餐厅本地训练，只上传
    模型参数，不上传原始数据。这样既能共享知识，
    又保护隐私。"
```

### 11.8 AI扩展的商业价值

**量化指标：**
```
传统方式 vs AI优化

处理效率:
  传统: 10秒/订单 (固定)
  AI优化: 8-12秒/订单 (根据复杂度)
  提升: 平均效率 +15%

机器人利用率:
  传统: 60-70% (高峰期不足，低峰期浪费)
  AI优化: 75-85% (自动扩缩容)
  提升: 利用率 +20%

异常处理:
  传统: 人工发现，平均响应时间 5分钟
  AI优化: 自动检测，平均响应时间 10秒
  提升: 响应速度 +3000%

客户满意度:
  传统: 平均等待时间 3.5分钟
  AI优化: 平均等待时间 2.3分钟
  提升: 等待时间 -34%

成本节省:
  传统: 需要20%冗余机器人应对高峰
  AI优化: 动态调整，冗余降至5%
  节省: 运营成本 -15%
```

---

# Part 6: 项目结构和面试准备

## 12. 完整项目结构

```
feedme/
├── README.md                          # 项目说明
├── DESIGN.md                          # 设计文档（合并所有Part）
├── docker-compose.yml                 # Docker部署
├── .github/
│   └── workflows/
│       └── ci.yml                     # GitHub Actions
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                # 入口文件
│   │
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── order.go               # 订单领域模型
│   │   │   └── bot.go                 # 机器人领域模型
│   │   │
│   │   ├── repository/
│   │   │   ├── priority_queue.go      # 优先级队列
│   │   │   └── priority_queue_test.go
│   │   │
│   │   ├── service/
│   │   │   ├── order_manager.go       # 订单管理
│   │   │   ├── order_manager_test.go
│   │   │   ├── bot_manager.go         # 机器人管理
│   │   │   ├── bot_manager_test.go
│   │   │   ├── event_publisher.go     # 事件发布
│   │   │   └── event_publisher_test.go
│   │   │
│   │   ├── handler/
│   │   │   ├── order_handler.go       # 订单HTTP处理
│   │   │   ├── bot_handler.go         # 机器人HTTP处理
│   │   │   └── websocket_handler.go   # WebSocket处理
│   │   │
│   │   ├── websocket/
│   │   │   ├── hub.go                 # WebSocket Hub
│   │   │   └── client.go              # WebSocket Client
│   │   │
│   │   ├── cli/
│   │   │   ├── interactive.go         # 交互式CLI
│   │   │   └── simulator.go           # 自动化模拟器
│   │   │
│   │   └── config/
│   │       └── config.go              # 配置管理
│   │
│   ├── pkg/
│   │   ├── logger/
│   │   │   └── logger.go              # 日志封装
│   │   └── utils/
│   │       └── id_generator.go        # ID生成器
│   │
│   ├── scripts/
│   │   ├── build.sh                   # 构建脚本
│   │   ├── test.sh                    # 测试脚本
│   │   └── run.sh                     # 运行脚本
│   │
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── result.txt                     # CLI输出（运行后生成）
│
└── frontend/
    ├── src/
    │   ├── components/
    │   │   ├── ControlPanel.tsx
    │   │   ├── OrderBoard.tsx
    │   │   ├── OrderCard.tsx
    │   │   ├── BotPanel.tsx
    │   │   ├── BotCard.tsx
    │   │   └── StatsPanel.tsx
    │   │
    │   ├── hooks/
    │   │   ├── useWebSocket.ts
    │   │   └── useStore.ts
    │   │
    │   ├── services/
    │   │   └── api.ts
    │   │
    │   ├── types/
    │   │   └── index.ts
    │   │
    │   ├── App.tsx
    │   ├── main.tsx
    │   └── index.css
    │
    ├── public/
    ├── package.json
    ├── tsconfig.json
    ├── vite.config.ts
    ├── tailwind.config.js
    └── Dockerfile
```

## 13. 脚本文件内容

### 12.1 scripts/build.sh
```bash
#!/bin/bash

echo "Building McDonald's Bot System..."

# 清理旧的构建
rm -rf bin/
mkdir -p bin/

# 构建
go build -o bin/mcdonalds-bot cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful: ./bin/mcdonalds-bot"
    exit 0
else
    echo "❌ Build failed"
    exit 1
fi
```

### 12.2 scripts/test.sh
```bash
#!/bin/bash

echo "Running tests..."

# 运行所有测试
go test -v ./internal/...

if [ $? -eq 0 ]; then
    echo "✅ All tests passed"
    exit 0
else
    echo "❌ Tests failed"
    exit 1
fi
```

### 12.3 scripts/run.sh
```bash
#!/bin/bash

echo "Running simulation..."

# 运行CLI模拟模式
./bin/mcdonalds-bot simulate > result.txt

if [ $? -eq 0 ]; then
    echo "✅ Simulation completed. Check result.txt"
    exit 0
else
    echo "❌ Simulation failed"
    exit 1
fi
```

## 14. GitHub Actions配置

### 13.1 .github/workflows/ci.yml
```yaml
name: Go Verify Result

on:
  pull_request:
    branches: [ main ]
  push:
    branches: [ main ]

jobs:
  verify:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      working-directory: ./backend
      run: go mod download
    
    - name: Run tests
      working-directory: ./backend
      run: |
        chmod +x scripts/test.sh
        ./scripts/test.sh
    
    - name: Build
      working-directory: ./backend
      run: |
        chmod +x scripts/build.sh
        ./scripts/build.sh
    
    - name: Run simulation
      working-directory: ./backend
      run: |
        chmod +x scripts/run.sh
        ./scripts/run.sh
    
    - name: Verify result.txt
      working-directory: ./backend
      run: |
        if [ ! -f "result.txt" ]; then
          echo "❌ result.txt not found"
          exit 1
        fi
        
        # 验证时间格式 HH:MM:SS
        if ! grep -E "[0-9]{2}:[0-9]{2}:[0-9]{2}" result.txt; then
          echo "❌ Invalid timestamp format"
          exit 1
        fi
        
        # 验证包含关键事件
        if ! grep -q "ORDER_CREATED" result.txt; then
          echo "❌ Missing ORDER_CREATED event"
          exit 1
        fi
        
        if ! grep -q "BOT_ADDED" result.txt; then
          echo "❌ Missing BOT_ADDED event"
          exit 1
        fi
        
        if ! grep -q "ORDER_COMPLETED" result.txt; then
          echo "❌ Missing ORDER_COMPLETED event"
          exit 1
        fi
        
        echo "✅ result.txt verification passed"
    
    - name: Upload result.txt
      uses: actions/upload-artifact@v3
      with:
        name: result
        path: backend/result.txt
```



## 15. 最后检查清单

**提交PR前：**
- [ ] 所有测试通过
- [ ] result.txt格式正确
- [ ] README.md完善
- [ ] 代码无明显bug
- [ ] GitHub Actions通过
- [ ] 前端UI正常
- [ ] WebSocket连接正常
- [ ] 手动测试所有测试用例
- [ ] 代码已格式化（gofmt）
- [ ] 无敏感信息（密码、密钥）


