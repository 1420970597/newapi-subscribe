# NewAPI Subscribe - 订阅管理系统

一个与 [new-api](https://github.com/Calcium-Ion/new-api) 集成的订阅管理系统，实现按天/周/月/自定义时间的 AI 模型订阅服务，支持每日额度管理和结转策略。

## 功能特性

### 订阅管理
- **灵活的订阅周期**: 支持按天、周、月或自定义天数订阅
- **每日额度控制**: 每天分配固定额度，用完即止
- **额度结转**: 可配置前一天剩余额度是否结转到第二天，支持设置结转上限
- **自动同步**: 每天 0:00 自动同步额度到 new-api

### 用户功能
- **多种登录方式**: 支持本系统注册登录，也支持 new-api 账号快捷登录
- **订阅购买**: 支持支付宝/微信支付（易支付）
- **续费管理**: 支持订阅续费，自动延长有效期
- **使用统计**: 查看使用记录和模型消费分析
- **到期提醒**: 邮件提醒订阅即将到期

### 管理功能
- **套餐管理**: 创建和管理订阅套餐，绑定 new-api 模型分组
- **用户管理**: 查看用户列表、订阅状态、使用分析
- **订单管理**: 查看所有订单记录
- **系统设置**: 站点信息、访问控制、支付配置等

## 技术栈

| 组件 | 技术 |
|-----|------|
| 后端 | Go 1.21 + Gin + GORM |
| 前端 | React 18 + TypeScript + Ant Design 5 |
| 数据库 | SQLite |
| 部署 | Docker + Docker Compose |

## 快速开始

### 环境要求

- Docker 20.10+
- Docker Compose 2.0+

### 1. 克隆项目

```bash
git clone https://github.com/1420970597/newapi-subscribe.git
cd newapi-subscribe
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，配置以下必要参数：

```env
# new-api 配置（必填）
NEWAPI_URL=https://your-newapi-site.com
NEWAPI_ADMIN_USER=admin
NEWAPI_ADMIN_PASS=your-admin-password

# 易支付配置（必填，用于在线支付）
EPAY_URL=https://pay.example.com
EPAY_PID=10001
EPAY_KEY=your-epay-key

# JWT 密钥（建议修改）
JWT_SECRET=your-random-secret-key
```

### 3. 启动服务

```bash
docker-compose up -d
```

### 4. 访问系统

- 前台地址: `http://localhost:8080`
- 默认管理员: `admin` / `admin123`

> 首次登录后请及时修改管理员密码

## 配置说明

### 环境变量完整列表

```env
# ========== 服务配置 ==========
PORT=8080                          # 服务端口
JWT_SECRET=change-me-in-production # JWT 密钥，请使用随机字符串

# ========== 数据库 ==========
DB_PATH=./data/subscribe.db        # SQLite 数据库路径

# ========== new-api 配置 ==========
NEWAPI_URL=https://your-newapi-site.com  # new-api 站点地址
NEWAPI_ADMIN_USER=admin                   # new-api 管理员用户名
NEWAPI_ADMIN_PASS=password                # new-api 管理员密码

# ========== 易支付配置 ==========
EPAY_URL=https://pay.example.com   # 易支付网关地址
EPAY_PID=10001                     # 商户 ID
EPAY_KEY=your-epay-key             # 商户密钥

# ========== SMTP 邮件配置 ==========
SMTP_SERVER=smtp.example.com       # SMTP 服务器
SMTP_PORT=587                      # SMTP 端口
SMTP_USER=noreply@example.com      # SMTP 用户名
SMTP_PASS=password                 # SMTP 密码
SMTP_FROM=noreply@example.com      # 发件人地址

# ========== 定时任务 ==========
CRON_ENABLED=true                  # 是否启用定时任务
CRON_SCHEDULE=0 0 * * *            # Cron 表达式（默认每天 0:00）
```

### new-api 配置要求

本系统需要使用 new-api 的管理员账号来操作用户数据，请确保：

1. `NEWAPI_ADMIN_USER` 账号具有管理员权限（Role >= 10）
2. new-api 站点允许 API 访问

### 易支付配置

支持标准易支付接口，请联系您的易支付服务商获取：
- 商户 ID (PID)
- 商户密钥 (Key)
- 网关地址

## 使用指南

### 创建订阅套餐

1. 使用管理员账号登录
2. 进入「管理后台」→「套餐管理」
3. 点击「新建套餐」
4. 填写套餐信息：
   - **套餐名称**: 如「基础版」「专业版」
   - **周期类型**: 天/周/月/自定义
   - **周期天数**: 订阅持续天数
   - **每日额度**: 每天分配的额度数量
   - **支持结转**: 是否允许未用完的额度结转
   - **最大结转额度**: 结转上限，0 表示无限制
   - **价格类型**: 固定价格或按天计价
   - **价格**: 套餐价格
   - **new-api 分组**: 绑定的模型分组

### 用户购买流程

1. 用户访问首页，查看套餐列表
2. 选择套餐，点击「立即订阅」
3. 选择 new-api 账号处理方式：
   - **创建新账号**: 系统自动在 new-api 创建账号
   - **绑定现有账号**: 使用已有的 new-api 账号
   - **覆盖当前账号**: 清空现有余额，设置为套餐额度
4. 选择支付方式，完成支付
5. 支付成功后订阅立即生效

### 额度同步机制

系统每天 0:00 自动执行额度同步：

1. 查询所有活跃订阅
2. 检查是否过期，过期则清零余额
3. 计算新的每日额度：
   - 如果支持结转: `新额度 = 每日额度 + min(昨日剩余, 最大结转)`
   - 如果不结转: `新额度 = 每日额度`
4. 更新 new-api 用户余额
5. 发送到期提醒邮件

也可以在管理后台手动触发同步。

## API 接口

### 认证接口

| 方法 | 路径 | 说明 |
|-----|------|-----|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/login/newapi | new-api 账号登录 |
| GET | /api/auth/me | 获取当前用户信息 |

### 套餐接口

| 方法 | 路径 | 说明 |
|-----|------|-----|
| GET | /api/plans | 获取套餐列表 |
| GET | /api/plans/:id | 获取套餐详情 |
| GET | /api/plans/:id/models | 获取套餐可用模型 |

### 订阅接口

| 方法 | 路径 | 说明 |
|-----|------|-----|
| GET | /api/subscriptions/current | 获取当前订阅 |
| POST | /api/subscriptions/purchase | 购买订阅 |
| POST | /api/subscriptions/renew | 续费订阅 |
| GET | /api/subscriptions/usage | 获取使用日志 |

### 订单接口

| 方法 | 路径 | 说明 |
|-----|------|-----|
| GET | /api/orders | 获取订单列表 |
| POST | /api/orders/pay | 发起支付 |
| GET | /api/orders/notify | 支付回调 |

### 管理接口

| 方法 | 路径 | 说明 |
|-----|------|-----|
| GET | /api/admin/users | 获取用户列表 |
| GET | /api/admin/subscriptions | 获取所有订阅 |
| GET | /api/admin/orders | 获取所有订单 |
| POST | /api/admin/plans | 创建套餐 |
| PUT | /api/admin/plans/:id | 更新套餐 |
| DELETE | /api/admin/plans/:id | 删除套餐 |
| GET | /api/admin/settings | 获取系统设置 |
| PUT | /api/admin/settings | 更新系统设置 |
| POST | /api/admin/sync/trigger | 手动触发同步 |

## 项目结构

```
newapi-subscribe/
├── backend/                      # Go 后端
│   ├── cmd/server/main.go        # 入口文件
│   └── internal/
│       ├── config/               # 配置管理
│       ├── model/                # 数据模型
│       ├── controller/           # API 控制器
│       ├── service/              # 业务逻辑
│       │   ├── newapi_client.go  # new-api 客户端
│       │   ├── subscription.go   # 订阅服务
│       │   ├── epay.go           # 易支付服务
│       │   └── email.go          # 邮件服务
│       ├── middleware/           # 中间件
│       ├── router/               # 路由
│       ├── cron/                 # 定时任务
│       └── dto/                  # 数据传输对象
├── frontend/                     # React 前端
│   └── src/
│       ├── pages/                # 页面组件
│       │   ├── Home/             # 首页
│       │   ├── Login/            # 登录
│       │   ├── User/             # 用户中心
│       │   └── Admin/            # 管理后台
│       ├── components/           # 通用组件
│       ├── api/                  # API 调用
│       └── store/                # 状态管理
├── Dockerfile                    # Docker 构建文件
├── docker-compose.yml            # Docker Compose 配置
└── .env.example                  # 环境变量示例
```

## 数据库模型

### 用户表 (users)
| 字段 | 类型 | 说明 |
|-----|------|-----|
| id | INTEGER | 主键 |
| username | VARCHAR(64) | 用户名 |
| password | VARCHAR(255) | 密码哈希 |
| email | VARCHAR(128) | 邮箱 |
| role | INTEGER | 角色 (1=用户, 10=管理员) |
| newapi_user_id | INTEGER | new-api 用户 ID |
| newapi_username | VARCHAR(64) | new-api 用户名 |
| newapi_bound | INTEGER | 是否已绑定 |

### 套餐表 (plans)
| 字段 | 类型 | 说明 |
|-----|------|-----|
| id | INTEGER | 主键 |
| name | VARCHAR(128) | 套餐名称 |
| period_type | VARCHAR(16) | 周期类型 |
| period_days | INTEGER | 周期天数 |
| daily_quota | INTEGER | 每日额度 |
| carry_over | INTEGER | 是否结转 |
| max_carry_over | INTEGER | 最大结转额度 |
| price_type | VARCHAR(16) | 价格类型 |
| price | DECIMAL | 价格 |
| newapi_group | VARCHAR(64) | new-api 分组 |

### 订阅表 (subscriptions)
| 字段 | 类型 | 说明 |
|-----|------|-----|
| id | INTEGER | 主键 |
| user_id | INTEGER | 用户 ID |
| plan_id | INTEGER | 套餐 ID |
| status | VARCHAR(16) | 状态 |
| start_date | DATE | 开始日期 |
| end_date | DATE | 结束日期 |
| today_quota | INTEGER | 今日额度 |
| carried_quota | INTEGER | 结转额度 |

### 订单表 (orders)
| 字段 | 类型 | 说明 |
|-----|------|-----|
| id | INTEGER | 主键 |
| order_no | VARCHAR(64) | 订单号 |
| user_id | INTEGER | 用户 ID |
| plan_id | INTEGER | 套餐 ID |
| order_type | VARCHAR(16) | 订单类型 |
| amount | DECIMAL | 金额 |
| status | VARCHAR(16) | 状态 |

## 开发指南

### 本地开发

**后端**:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

**前端**:
```bash
cd frontend
npm install
npm run dev
```

### 构建部署

```bash
# 构建镜像
docker build -t newapi-subscribe .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v ./data:/app/data \
  --env-file .env \
  newapi-subscribe
```

## 常见问题

### Q: 如何修改管理员密码？
A: 目前需要直接操作数据库修改，后续版本会添加密码修改功能。

### Q: 支付回调失败怎么办？
A: 检查易支付配置是否正确，确保回调地址可以被外网访问。

### Q: 额度同步失败怎么办？
A: 检查 new-api 管理员账号密码是否正确，确保有足够权限。可以在管理后台手动触发同步。

### Q: 如何查看系统日志？
A:
```bash
docker-compose logs -f
```

## 许可证

MIT License

## 相关项目

- [new-api](https://github.com/Calcium-Ion/new-api) - OpenAI API 管理与分发系统
