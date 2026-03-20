# 🍔 McDonald's 自动化订单管理系统

这是一个完整的前后端分离项目，实现了麦当劳自动化烹饪机器人订单管理系统。

## 📁 项目结构

```
feedme/
├── mcdonalds-bot-backend/     # Go后端项目
│   ├── cmd/server/            # 服务器入口
│   ├── internal/              # 内部包
│   │   ├── models/            # 数据模型
│   │   ├── queue/             # 优先级队列
│   │   ├── scheduler/         # 调度器
│   │   ├── handlers/          # HTTP处理器
│   │   └── websocket/         # WebSocket Hub
│   └── tests/                 # 测试文件
│
├── mcdonalds-bot-frontend/    # React前端项目
│   ├── src/
│   │   ├── components/        # React组件
│   │   ├── hooks/             # 自定义Hooks
│   │   ├── services/          # API服务
│   │   └── types/             # TypeScript类型
│   └── package.json
│
├── DESIGN.md                  # 完整技术设计文档
└── README.md                  # 本文件
```

## ✨ 核心功能

✅ **订单管理**
- 创建普通订单和VIP订单
- VIP订单自动优先处理（FIFO顺序）
- 实时订单状态更新

✅ **机器人管理**
- 动态添加/删除机器人
- 自动分配订单给空闲机器人
- 机器人移除时订单自动回到队列

✅ **实时通信**
- WebSocket实时推送事件
- 前端自动更新界面
- 连接状态显示

✅ **完整测试**
- 7个单元测试全部通过
- 覆盖所有核心功能
- 测试VIP优先级、并发处理等

## 🚀 Quick Start（零基础也能运行）

### 前置要求

确保你的Mac上已安装：
- **Go 1.21+** - 后端语言
- **Node.js 18+** - 前端运行环境
- **npm** - Node包管理器（通常随Node.js一起安装）

#### 检查是否已安装

```bash
# 检查Go版本
go version

# 检查Node.js版本
node --version

# 检查npm版本
npm --version
```

#### 如果没有安装

**安装Go:**
```bash
# 使用Homebrew安装
brew install go
```

**安装Node.js:**
```bash
# 使用Homebrew安装
brew install node
```

### 步骤1：启动后端服务器

打开终端，执行以下命令：

```bash
# 进入后端目录
cd /Users/huangyiming/Documents/workspace/feedme/mcdonalds-bot-backend

# 设置Go代理（加速下载）
export GOPROXY=https://goproxy.cn,direct

# 启动后端服务器
go run cmd/server/main.go
```

**成功标志：**
你会看到类似这样的输出：
```
[GIN-debug] Listening and serving HTTP on :8080
2026/03/20 01:42:43 Server starting on :8080
```

**保持这个终端窗口打开！**

### 步骤2：启动前端应用

打开**新的**终端窗口，执行：

```bash
# 进入前端目录
cd /Users/huangyiming/Documents/workspace/feedme/mcdonalds-bot-frontend

# 启动前端开发服务器
npm run dev
```

**成功标志：**
你会看到：
```
  VITE v8.0.1  ready in 438 ms
  ➜  Local:   http://localhost:5173/
```

### 步骤3：打开浏览器

在浏览器中访问：**http://localhost:5173/**

你会看到一个漂亮的管理界面！

## 🎮 如何使用

### 1. 添加机器人

在"Add Bot"区域：
- 输入机器人ID（例如：`bot1`）
- 点击"Add Bot"按钮
- 机器人会出现在左侧"Bots"列表中

### 2. 创建订单

在"Create Order"区域：
- 输入订单ID（例如：`order1`）
- 选择订单类型（Normal 或 VIP）
- 点击"Create"按钮

### 3. 观察自动处理

- 如果有空闲机器人，订单会自动分配
- 机器人状态变为"Processing"
- 10秒后订单自动完成
- 完成的订单出现在右侧"Completed Orders"列表

### 4. 测试VIP优先级

1. 先创建一个Normal订单（`order1`）
2. 再创建一个VIP订单（`vip1`）
3. 添加一个机器人（`bot1`）
4. 观察：VIP订单会优先被处理！

### 5. 测试机器人移除

1. 添加机器人并让它处理订单
2. 在处理过程中点击"Remove"按钮
3. 观察：订单会自动回到待处理队列

## 🧪 运行测试

```bash
# 进入后端目录
cd /Users/huangyiming/Documents/workspace/feedme/mcdonalds-bot-backend

# 运行所有测试
go test ./tests/... -v

# 运行特定测试
go test ./tests/... -v -run TestVIPPriority
```

**所有测试都应该通过（PASS）！**

## 📊 API文档

### REST API

**基础URL:** `http://localhost:8080/api`

#### 机器人管理

```bash
# 添加机器人
POST /api/bots
Content-Type: application/json
{"id": "bot1"}

# 删除机器人
DELETE /api/bots
Content-Type: application/json
{"id": "bot1"}

# 获取所有机器人
GET /api/bots
```

#### 订单管理

```bash
# 创建订单
POST /api/orders
Content-Type: application/json
{"id": "order1", "type": "Normal"}  # 或 "VIP"

# 获取所有订单
GET /api/orders

# 获取待处理订单
GET /api/orders/pending

# 获取已完成订单
GET /api/orders/completed
```

#### 系统管理

```bash
# 获取统计信息
GET /api/stats

# 重置系统
POST /api/reset
```

### WebSocket

**连接URL:** `ws://localhost:8080/ws`

**事件类型:**
- `order_created` - 订单创建
- `order_processing` - 订单开始处理
- `order_complete` - 订单完成
- `bot_added` - 机器人添加
- `bot_removed` - 机器人移除
- `queue_updated` - 队列更新

## 🐛 常见问题

### Q1: 后端启动失败，提示端口被占用

**解决方法:**
```bash
# 查找占用8080端口的进程
lsof -i :8080

# 杀死该进程
kill -9 <PID>
```

### Q2: 前端启动失败，提示端口被占用

**解决方法:**
```bash
# 查找占用5173端口的进程
lsof -i :5173

# 杀死该进程
kill -9 <PID>
```

### Q3: 前端无法连接后端

**检查清单:**
1. 后端是否正在运行？（检查终端输出）
2. 后端是否在8080端口？
3. 浏览器控制台是否有错误？

### Q4: WebSocket连接失败

**解决方法:**
1. 确保后端正在运行
2. 刷新浏览器页面
3. 检查浏览器控制台的错误信息

### Q5: Go依赖下载很慢

**解决方法:**
```bash
# 设置国内代理
export GOPROXY=https://goproxy.cn,direct

# 或者使用阿里云代理
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

## 📝 测试场景

### 场景1：基本功能测试

```bash
# 1. 添加机器人
curl -X POST http://localhost:8080/api/bots \
  -H "Content-Type: application/json" \
  -d '{"id":"bot1"}'

# 2. 创建订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"id":"order1","type":"Normal"}'

# 3. 查看统计
curl http://localhost:8080/api/stats
```

### 场景2：VIP优先级测试

```bash
# 1. 创建Normal订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"id":"normal1","type":"Normal"}'

# 2. 创建VIP订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"id":"vip1","type":"VIP"}'

# 3. 查看待处理队列（VIP应该在前面）
curl http://localhost:8080/api/orders/pending

# 4. 添加机器人（会优先处理VIP订单）
curl -X POST http://localhost:8080/api/bots \
  -H "Content-Type: application/json" \
  -d '{"id":"bot1"}'
```

### 场景3：并发处理测试

```bash
# 1. 添加多个机器人
curl -X POST http://localhost:8080/api/bots \
  -H "Content-Type: application/json" \
  -d '{"id":"bot1"}'

curl -X POST http://localhost:8080/api/bots \
  -H "Content-Type: application/json" \
  -d '{"id":"bot2"}'

# 2. 创建多个订单
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"id":"order1","type":"Normal"}'

curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"id":"order2","type":"VIP"}'

curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"id":"order3","type":"Normal"}'

# 3. 查看统计（应该有2个机器人在处理）
curl http://localhost:8080/api/stats
```

## 🎯 技术亮点

1. **优先级队列实现** - 使用Go的container/heap实现高效的优先级队列
2. **并发安全** - 使用sync.RWMutex保证线程安全
3. **实时通信** - WebSocket实现前后端实时数据同步
4. **事件驱动架构** - 解耦的事件监听器模式
5. **完整测试覆盖** - 7个单元测试覆盖所有核心功能
6. **TypeScript类型安全** - 前端使用TypeScript确保类型安全
7. **响应式UI** - React Hooks实现现代化的响应式界面

## 📚 更多文档

- **完整技术设计文档**: 查看 `DESIGN.md`
- **API详细文档**: 查看 `DESIGN.md` 第7章
- **测试指南**: 查看 `DESIGN.md` 第9章
- **面试准备**: 查看 `DESIGN.md` 第15章

## 🎉 完成检查清单

在提交之前，请确认：

- [ ] 后端服务器可以正常启动
- [ ] 前端应用可以正常启动
- [ ] 所有测试都通过（`go test ./tests/... -v`）
- [ ] 可以通过浏览器访问前端界面
- [ ] 可以添加机器人
- [ ] 可以创建订单
- [ ] VIP订单优先处理
- [ ] WebSocket连接正常（界面右上角显示"Connected"）
- [ ] 机器人移除后订单回到队列

## 💡 提示

1. **第一次运行**: 后端第一次启动可能需要下载依赖，请耐心等待
2. **网络问题**: 如果Go依赖下载失败，使用国内代理：`export GOPROXY=https://goproxy.cn,direct`
3. **端口冲突**: 如果端口被占用，使用`lsof -i :端口号`查找并杀死进程
4. **浏览器缓存**: 如果前端显示异常，尝试硬刷新（Cmd+Shift+R）

## 📞 支持

如果遇到问题：
1. 检查终端输出的错误信息
2. 查看浏览器控制台的错误
3. 参考`DESIGN.md`中的FAQ章节
4. 确保Go和Node.js版本符合要求

---

**祝你使用愉快！Good luck with your interview! 🚀**
