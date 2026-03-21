# 秒杀系统部署文档

## 系统架构

本秒杀系统采用前后端分离架构，主要组件包括：

- **前端**：Vue 3 + Vite
- **后端**：Go + Gin + MySQL + Redis + RabbitMQ
- **安全**：JWT认证 + 速率限制 + 验证码

## 环境要求

### 前端环境
- Node.js 16.x 或更高版本
- npm 7.x 或更高版本

### 后端环境
- Go 1.20 或更高版本
- MySQL 8.0 或更高版本
- Redis 6.0 或更高版本
- RabbitMQ 3.8 或更高版本

## 部署步骤

### 1. 克隆项目

```bash
git clone <项目仓库地址>
cd Gin秒杀API
```

### 2. 后端部署

#### 2.1 配置文件修改

编辑 `seckill-go-master/config.yaml` 文件，修改数据库、Redis和RabbitMQ的连接信息：

```yaml
db:
  dsn: "root:password@tcp(127.0.0.1:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

rabbitmq:
  addr: "127.0.0.1:5672"
  username: "guest"
  password: "guest"
  vhost: "/"
```

#### 2.2 安装依赖

```bash
cd seckill-go-master
go mod tidy
```

#### 2.3 启动后端服务

```bash
go run main.go
```

后端服务默认运行在 `http://localhost:8081`。

### 3. 前端部署

#### 3.1 安装依赖

```bash
cd seckill-vue-master
npm install --legacy-peer-deps
```

#### 3.2 启动前端服务

```bash
npm run dev
```

前端服务默认运行在 `http://localhost:5173`。

### 4. 数据库初始化

执行以下SQL语句初始化数据库：

```sql
CREATE DATABASE seckill;

USE seckill;

-- 用户表
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(100) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品表
CREATE TABLE `product` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `stock` int(11) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 活动表
CREATE TABLE `event` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `status` int(11) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 活动商品表
CREATE TABLE `event_product` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `event_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `event_id` (`event_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `event_product_ibfk_1` FOREIGN KEY (`event_id`) REFERENCES `event` (`id`),
  CONSTRAINT `event_product_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 订单表
CREATE TABLE `order` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `product_id` int(11) NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `status` int(11) NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `order_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`),
  CONSTRAINT `order_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 5. 系统配置

#### 5.1 安全配置

- **JWT密钥**：在 `utils/jwt.go` 文件中修改JWT密钥
- **加密密钥**：在 `utils/encryption.go` 文件中修改加密密钥
- **速率限制**：在 `utils/rate_limit.go` 文件中修改速率限制参数

#### 5.2 性能配置

- **Redis连接池**：在 `db/redis.go` 文件中修改Redis连接池参数
- **RabbitMQ连接**：在 `rabbitmq/connection.go` 文件中修改RabbitMQ连接参数

## 运行测试

### 1. 前端测试

```bash
cd seckill-vue-master
npm run test
```

### 2. 后端测试

```bash
cd seckill-go-master
go test ./...
```

## 压力测试

使用wrk工具进行压力测试：

```bash
wrk -t12 -c400 -d30s http://localhost:8081/api/product/list
```

## 安全测试

### 1. 模拟恶意请求

使用ab工具模拟恶意请求：

```bash
ab -n 1000 -c 100 http://localhost:8081/api/product/seckill
```

### 2. SQL注入测试

使用SQLmap工具进行SQL注入测试：

```bash
sqlmap -u "http://localhost:8081/api/product/list" --batch
```

## 故障排除

### 1. 后端服务无法启动

- 检查数据库连接是否正常
- 检查Redis连接是否正常
- 检查RabbitMQ连接是否正常
- 检查端口是否被占用

### 2. 前端无法访问后端API

- 检查CORS配置是否正确
- 检查前端代理配置是否正确
- 检查后端服务是否正常运行

### 3. 秒杀功能无法正常工作

- 检查Redis是否正常运行
- 检查RabbitMQ是否正常运行
- 检查库存是否充足
- 检查活动是否处于活动状态

## 部署建议

1. **生产环境部署**：
   - 使用Nginx作为反向代理
   - 配置HTTPS
   - 使用PM2管理Node.js进程
   - 使用systemd管理Go进程

2. **性能优化**：
   - 使用Redis集群
   - 使用RabbitMQ集群
   - 数据库读写分离
   - 增加服务器资源

3. **安全建议**：
   - 使用HTTPS
   - 定期更新依赖
   - 配置防火墙
   - 定期备份数据库
