# Currency Exchange App

货币汇率管理后端服务：提供用户注册、登录、JWT 身份认证与汇率数据管理接口。基于 Gin + GORM + MySQL 构建，是完整的 Gin 项目分层实践（配置 / 模型 / 控制器 / 路由 / 中间件 / 工具）。

## 功能特性

- 用户注册（bcrypt 密码加密存储，用户名唯一约束）
- 用户登录（密码校验通过后签发 JWT Token）
- JWT 中间件鉴权（保护需要登录的接口）
- 汇率创建（同货币对同日期唯一，组合唯一索引兜底）
- 请求体校验（binding 标签 + validator）
- 防 SQL 注入（GORM 参数占位符）
- 防用户名枚举（登录失败统一返回 `invalid username or password`）
- 敏感配置不入库（数据库密码 / JWT 密钥仅存于本地 config.yml）

## 技术栈

| 技术 | 用途 |
|---|---|
| Go | 后端语言 |
| Gin | HTTP Web 框架 / 路由 / 中间件 |
| GORM | ORM，数据库操作 / 自动建表 |
| MySQL | 数据存储 |
| Viper | 配置文件读取 |
| golang-jwt/jwt/v5 | JWT 签发与验证 |
| bcrypt | 密码哈希 |

## 目录结构

```
CurrencyExchangeApp/
├── main.go                  # 程序入口：初始化配置 → 注册路由 → 启动
├── config/
│   ├── config.go            # 配置模型 + viper 读取
│   ├── config.yml           # 本地配置
│   ├── config_example.yml   # 配置模板复制为 config.yml 使用
│   └── db.go                # 数据库初始化：连接池 + Ping + AutoMigrate
├── models/
│   ├── article_like.go      # 文章点赞模型
│   ├── article.go           # 文章模型
│   ├── user.go              # 用户模型
│   └── exchange_rate.go     # 汇率模型
├── controllers/
│   ├── article_controller.go # 文章处理
│   ├── auth_controller.go   # 注册 / 登录
│   ├── user_controller.go   # 获取用户信息
│   └── rate_controller.go   # 汇率处理
├── middleware/
│   └── auth.go              # JWT 鉴权中间件
├── utils/
│   └── utils.go             # bcrypt 加密 / JWT 生成解析
├── router/
│   └── router.go            # 路由注册 （分组挂载中间件）
├── global/
│   └── global.go            # 全局变量（DB 句柄）
├── go.mod / go.sum
└── README.md
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+ / 8.x（本地已启动）

### 启动步骤

1. **复制配置模板**：

   ```bash
   cp config/config_example.yml config/config.yml
   ```

2. **修改本地配置** `config/config.yml`：

   ```yaml
   database:
     dsn: "root:你的密码@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
   jwt:
     Secret: "换成你自己的随机密钥"
   ```

3. **创建数据库**（MySQL 中执行）：

   ```sql
   CREATE DATABASE IF NOT EXISTS go_db DEFAULT CHARSET utf8mb4;
   ```

4. **启动服务**：

   ```bash
   go run main.go
   ```

   启动时 GORM 会自动建表（users、exchange_rates），日志输出无报错即成功。服务监听 `:3000`。

## 接口文档

> 所有接口：`Content-Type: application/json`

### 1. 用户注册

```
POST /api/auth/register
```

请求体：

```json
{
  "username": "tom",
  "password": "123456"
}
```

校验规则：`username` 3-20 位必填，`password` 至少 6 位。

| 场景 | 状态码 | 响应 |
|---|---|---|
| 注册成功 | 200 | `{"message": "Create Success"}` |
| 请求体不合法 | 400 | `{"error": "invalid request body"}` |
| 用户名已存在 | 409 | `{"error": "username already exists"}` |

### 2. 用户登录

```
POST /api/auth/login
```

请求体：

```json
{
  "username": "tom",
  "password": "123456"
}
```

| 场景 | 状态码 | 响应 |
|---|---|---|
| 登录成功 | 200 | `{"token": "Bearer eyJhbGci...", "message": "Login Success"}` |
| 用户名或密码错误 | 401 | `{"error": "invalid username or password"}` |

> 登录成功返回的 `token` 字段**已包含 `Bearer ` 前缀**，后续请求直接将整个值放入 `Authorization` 请求头。

### 3. 获取当前用户信息（需登录）

```
GET /api/user/info
```

请求头：

```
Authorization: Bearer eyJhbGci...
```

| 场景 | 状态码 | 响应 |
|---|---|---|
| 成功 | 200 | `{"username": "tom", "createdAt": "..."}` |
| 未携带 token | 401 | `{"error": "missing token"}` |
| token 格式错误 | 401 | `{"error": "invalid token format"}` |
| token 无效 / 过期 | 401 | `{"error": "invalid or expired token"}` |

### 4. 创建汇率

```
POST /api/exchangerate/add
```

请求体：

```json
{
  "from_currency": "USD",
  "to_currency": "CNY",
  "rate": 7.25,
  "date": "2026-10-21"
}
```

> `date` 为日期字符串（如 `2026-10-21`）；货币代码大小写不敏感（服务端自动统一为大写）。同一货币对同一天只能创建一条（组合唯一索引兜底）。

| 场景 | 状态码 | 响应 |
|---|---|---|
| 创建成功 | 200 | `{"message": "Create Success"}` |
| 请求体不合法 | 400 | `{"error": "invalid request body"}` |
| 已存在 | 409 | `{"error": "exchange rate already exists"}` |

### 5. 查询汇率

```
GET /api/exchangerate/query?from=USD&to=CNY&date=2026-10-21
```

查询参数：`from`（必填）、`to`（必填）、`date`（必填）。

| 场景 | 状态码 | 响应 |
|---|---|---|
| 查询成功 | 200 | `{"from_currency": "USD", "to_currency": "CNY", "rate": 7.25, "date": "2026-10-21"}` |
| 参数不合法 | 400 | `{"error": "invalid query parameters"}` |
| 查不到 | 404 | `{"error": "exchange rate not found"}` |

### 6. 发布文章（需登录）

```
POST /api/article/create
```

请求头：

```
Authorization: Bearer eyJhbGci...
```

请求体：

```json
{
  "title": "my first article",
  "content": "hello world"
}
```

校验规则：`title` 必填、长度 1-200；`content` 必填。

| 场景 | 状态码 | 响应 |
|---|---|---|
| 发布成功 | 200 | `{"message": "Create Success", "id": 1}` |
| 请求体不合法 | 400 | `{"error": "invalid request body"}` |
| 未携带 token | 401 | `{"error": "missing token"}` |
| token 格式错误 | 401 | `{"error": "invalid token format"}` |
| token 无效 / 过期 | 401 | `{"error": "invalid or expired token"}` |
| 服务器内部异常 | 500 | `{"error": "server error"}` |

### 7. 文章列表（支持搜索、分页）

```
GET /api/article?keyword=汇率&page=1&pagesize=10
```

查询参数（均可选）：

| 参数 | 默认值 | 说明 |
|---|---|---|
| `keyword` | 空 | 标题模糊匹配（包含即命中），空则查全部 |
| `page` | 1 | 页码，从 1 开始 |
| `pagesize` | 10 | 每页条数，最大 100 |

响应示例：

```json
{
  "list": [
    {
      "ID": 1,
      "CreatedAt": "2026-10-21T10:00:00+08:00",
      "title": "Go入门",
      "content": "...",
      "user_id": 1,
      "author": { "ID": 1, "username": "tom" },
      "like_count": 3
    }
  ],
  "total": 1
}
```

> `list` 按发布时间倒序；`author` 只包含 `id` 和 `username`（不含密码）；`total` 为符合筛选条件的文章总数（分页用）。

| 场景 | 状态码 | 响应 |
|---|---|---|
| 成功 | 200 | `{"list": [...], "total": N}` |
| 参数不合法 | 400 | `{"error": "invalid query parameters"}` |
| 服务器内部异常 | 500 | `{"error": "server error"}` |

### 8. 文章详情

```
GET /api/article/:id
```

| 场景 | 状态码 | 响应 |
|---|---|---|
| 查询成功 | 200 | `{"id": 1, "title": "Go入门", "content": "...", "likecount": 3, "author": {"ID": 1, "username": "tom"}, "created_at": "..."}` |
| id 不合法 | 400 | `{"error": "invalid id"}` |
| 文章不存在 | 404 | `{"error": "article not found"}` |
| 服务器内部异常 | 500 | `{"error": "server error"}` |


## 安全设计

- **密码安全**：bcrypt 加盐哈希（cost=12），数据库不存明文；登录比对用 `CompareHashAndPassword`
- **防 SQL 注入**：所有查询使用 GORM 参数占位符（`?`），不拼接 SQL 字符串
- **防用户名枚举**：登录时用户不存在与密码错误统一返回 `401 invalid username or password`
- **JWT**：HS256 签名，密钥来自本地配置，过期时间 24 小时；中间件三层检查（存在 → 格式 → 签名/过期）
- **敏感配置**：`config/config.yml` 已在 `.gitignore` 中，仓库只保留占位模板 `config_example.yml`

