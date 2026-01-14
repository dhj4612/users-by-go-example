# Go Web 用户管理系统

一个简单的 Go Web 用户管理系统，适合初学者学习。包含用户注册、登录、JWT 鉴权以及用户增删改查功能。

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT (JSON Web Token)
- **密码加密**: bcrypt

## 项目结构

```
users-by-go-example/
├── cmd/
│   └── main.go              # 程序入口
├── config/
│   ├── config.go            # 配置结构体
│   └── config.yml           # 配置文件
├── global/
│   └── global.go            # 全局资源管理
├── internal/                # 内部代码（不对外暴露）
│   ├── handler/
│   │   └── user_handler.go  # 用户处理器
│   ├── middleware/
│   │   └── auth.go          # JWT 认证中间件
│   ├── model/
│   │   └── user.go          # 用户模型
│   ├── router/
│   │   └── router.go        # 路由配置
│   └── service/
│       └── user_service.go  # 用户服务层
├── utils/
│   └── jwt.go               # JWT 工具类
├── init.sql                 # 数据库初始化脚本
├── go.mod
└── README.md
```

## 环境要求

- Go 1.25+
- MySQL 5.7+
- Redis 6.0+

## 快速开始

### 1. 初始化数据库

首先，确保 MySQL 和 Redis 服务已启动。

执行 `init.sql` 文件创建数据库和表：

```bash
mysql -u root -p111111 < init.sql
```

或者手动执行 SQL：

```sql
CREATE DATABASE IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `users`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（加密后）',
  `nike_name` varchar(50) DEFAULT NULL COMMENT '昵称',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete` tinyint(1) DEFAULT 0 COMMENT '是否删除 0-未删除 1-已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### 2. 配置文件

配置文件位于 `config/config.yml`，默认配置如下：

```yaml
server:
  port: 8080

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: 111111
  dbname: users

redis:
  host: localhost
  port: 6379
  password: 111111
  db: 0

jwt:
  secret: your-secret-key-change-this-in-production
  expire-time: 24  # 过期时间（小时）

white-list:
  - '/api/v1/login'
  - '/api/v1/register'
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 运行项目

```bash
go run cmd/main.go
```

服务器将在 `http://localhost:8080` 启动。

## API 接口

**注意**: 所有接口统一使用 POST 方法，参数通过 JSON body 传递。

### 1. 用户注册

**请求**:

```
POST /api/v1/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456",
  "nikeName": "测试用户"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "nikeName": "测试用户",
    "createTime": "2026-01-14T10:00:00Z",
    "updateTime": "2026-01-14T10:00:00Z"
  }
}
```

### 2. 用户登录

**请求**:

```
POST /api/v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "nikeName": "测试用户",
      "createTime": "2026-01-14T10:00:00Z",
      "updateTime": "2026-01-14T10:00:00Z"
    }
  }
}
```

### 3. 获取用户列表（需要认证）

**请求**:

```
POST /api/v1/users/list
Authorization: Bearer <token>
Content-Type: application/json

{
  "page": 1,
  "pageSize": 10
}
```

**响应**:

```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "testuser",
        "nikeName": "测试用户",
        "createTime": "2026-01-14T10:00:00Z",
        "updateTime": "2026-01-14T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
}
```

### 4. 获取用户详情（需要认证）

**请求**:

```
POST /api/v1/users/get
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": 1
}
```

**响应**:

```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "nikeName": "测试用户",
    "createTime": "2026-01-14T10:00:00Z",
    "updateTime": "2026-01-14T10:00:00Z"
  }
}
```

### 5. 更新用户信息（需要认证）

**请求**:

```
POST /api/v1/users/update
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": 1,
  "nikeName": "新昵称",
  "password": "newpassword"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "nikeName": "新昵称",
    "createTime": "2026-01-14T10:00:00Z",
    "updateTime": "2026-01-14T10:05:00Z"
  }
}
```

### 6. 删除用户（需要认证）

**请求**:

```
POST /api/v1/users/delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": 1
}
```

**响应**:

```json
{
  "code": 200,
  "message": "删除成功"
}
```

## 测试 API

你可以使用 curl、Postman 或其他 HTTP 客户端测试 API。

### 使用 curl 测试

1. **注册用户**:

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456","nikeName":"测试用户"}'
```

2. **登录获取 token**:

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

3. **获取用户列表**（替换 YOUR_TOKEN）:

```bash
curl -X POST http://localhost:8080/api/v1/users/list \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"page":1,"pageSize":10}'
```

4. **获取用户详情**（替换 YOUR_TOKEN）:

```bash
curl -X POST http://localhost:8080/api/v1/users/get \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

5. **更新用户信息**（替换 YOUR_TOKEN）:

```bash
curl -X POST http://localhost:8080/api/v1/users/update \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":1,"nikeName":"新昵称"}'
```

6. **删除用户**（替换 YOUR_TOKEN）:

```bash
curl -X POST http://localhost:8080/api/v1/users/delete \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

## 注意事项

1. **生产环境配置**: 在生产环境中，请修改 `config.yml` 中的 JWT secret 为更安全的密钥。
2. **密码安全**: 密码使用 bcrypt 加密存储，不会以明文形式保存。
3. **软删除**: 删除用户使用软删除方式，数据不会真正从数据库中删除。
4. **JWT 过期时间**: 默认 token 有效期为 24 小时，可在配置文件中修改。
5. **用户名脱敏**: 所有接口返回的用户名都会进行脱敏处理（保留首尾字符，中间用 * 替代）。
6. **分布式锁**: 注册和更新操作使用 Redis 分布式锁，保证接口幂等性。
7. **事务处理**: 所有写操作（注册、更新）都使用数据库事务，保证数据一致性。

## 核心特性

### 1. 分布式锁
- **注册接口**: 针对 username 加锁，防止重复注册
- **更新接口**: 针对用户 ID 加锁，防止并发更新冲突
- 锁超时时间：10 秒
- 支持重试机制：最多重试 3 次，每次间隔 100ms

### 2. 数据库事务
- 所有写操作都在事务中执行
- 自动回滚机制，保证数据一致性
- 支持 panic 恢复

### 3. 数据安全
- 密码使用 bcrypt 加密
- 用户名自动脱敏
- 登录响应不返回用户信息，只返回 token

## 学习要点

这个项目适合初学者学习以下内容：

1. **Go Web 开发**: 使用 Gin 框架构建 Web API
2. **数据库操作**: 使用 GORM 进行数据库 CRUD 操作
3. **JWT 认证**: 实现基于 JWT 的用户认证
4. **项目结构**: 学习标准的 Go 项目结构（internal 目录）
5. **中间件**: 理解和使用 Gin 中间件
6. **配置管理**: 使用 YAML 文件管理配置
7. **密码加密**: 使用 bcrypt 加密用户密码
8. **Redis 分布式锁**: 实现接口幂等性
9. **数据库事务**: 保证数据一致性
10. **优雅关闭**: 实现服务器优雅关闭

## 许可证

MIT License
