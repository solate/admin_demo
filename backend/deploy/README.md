# 部署指南

本目录包含了管理后台项目的 Docker 部署配置。

## 目录结构

```
deploy/
├── postgres/           # PostgreSQL 数据库配置
│   ├── docker-compose.yml
│   ├── data/          # 数据持久化目录
│   └── init-scripts/  # 初始化脚本目录
├── redis/             # Redis 缓存配置
│   ├── docker-compose.yml
│   └── redis/
│       ├── conf/      # Redis 配置文件
│       └── data/      # 数据持久化目录
└── README.md          # 本文档
```

## 服务配置

### PostgreSQL 配置
- **镜像**: postgres:15-alpine
- **容器名**: postgres_db
- **端口**: 5432
- **数据库**: testdb
- **用户名**: root
- **密码**: root
- **数据目录**: `./data/postgres`
- **初始化脚本**: `./init-scripts`

### Redis 配置
- **镜像**: redis:latest
- **容器名**: my_redis
- **端口**: 6379
- **密码**: 123456
- **数据目录**: `./redis/data`
- **配置文件**: `./redis/conf/redis.conf`

## 快速启动

### 1. 启动 PostgreSQL

```bash
cd postgres
docker-compose up -d
```

### 2. 启动 Redis

```bash
cd redis
docker-compose up -d
```

### 3. 验证服务

```bash
# 检查 PostgreSQL
docker exec -it postgres_db psql -U root -d testdb -c "SELECT version();"

# 检查 Redis
docker exec -it my_redis redis-cli -a 123456 ping
```

## 完整部署流程

### 1. 启动所有服务

```bash
# 启动 PostgreSQL
cd postgres && docker-compose up -d && cd ..

# 启动 Redis
cd redis && docker-compose up -d && cd ..
```

### 2. 更新应用配置

确保 `app/admin/etc/admin.yaml` 中的数据库配置正确：

```yaml
# DB 配置
DataSource: postgres://root:root@localhost:5432/testdb?sslmode=disable
ShowSQL: true

# 缓存
Redis:
  Host: localhost:6379
  Type: node
  Pass: "123456"
```

### 3. 初始化数据库

```bash
cd ../app/admin
go run cmd/init_db/init_db.go
```

### 4. 启动应用

```bash
go run admin.go -f etc/admin.yaml
```

## 管理命令

### PostgreSQL 管理

```bash
# 进入 PostgreSQL 容器
docker exec -it postgres_db psql -U root -d testdb

# 查看数据库列表
docker exec -it postgres_db psql -U root -d testdb -c "\l"

# 备份数据库
docker exec postgres_db pg_dump -U root testdb > backup.sql

# 恢复数据库
docker exec -i postgres_db psql -U root -d testdb < backup.sql
```

### Redis 管理

```bash
# 进入 Redis 容器
docker exec -it my_redis redis-cli -a 123456

# 查看 Redis 信息
docker exec -it my_redis redis-cli -a 123456 info

# 清空 Redis 数据
docker exec -it my_redis redis-cli -a 123456 flushall
```

### 容器管理

```bash
# 查看运行状态
docker ps

# 查看日志
docker logs postgres_db
docker logs my_redis

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 删除容器和数据
docker-compose down -v
```

## 数据持久化

### PostgreSQL 数据
- 数据存储在 `postgres/data/` 目录
- 即使删除容器，数据也会保留
- 要完全删除数据：`rm -rf postgres/data/`

### Redis 数据
- 数据存储在 `redis/redis/data/` 目录
- 配置文件在 `redis/redis/conf/redis.conf`
- 要完全删除数据：`rm -rf redis/redis/data/`

## 初始化脚本

PostgreSQL 支持在容器启动时自动执行初始化脚本：

1. 将 SQL 脚本放在 `postgres/init-scripts/` 目录
2. 脚本会按文件名顺序执行
3. 只有 `.sql` 和 `.sh` 文件会被执行

示例初始化脚本 `postgres/init-scripts/01-init.sql`：

```sql
-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建用户
CREATE USER app_user WITH PASSWORD 'app_password';

-- 授权
GRANT ALL PRIVILEGES ON DATABASE testdb TO app_user;
```

## 故障排除

### 常见问题

1. **端口冲突**
   - PostgreSQL 端口 5432 被占用
   - Redis 端口 6379 被占用
   - 解决方案：修改 `docker-compose.yml` 中的端口映射

2. **权限问题**
   - 数据目录权限不足
   - 解决方案：`chmod -R 755 postgres/data redis/redis/data`

3. **连接失败**
   - 检查容器是否正常运行：`docker ps`
   - 检查日志：`docker logs <container_name>`
   - 检查网络连接：`docker network ls`

4. **数据丢失**
   - 检查数据目录是否存在
   - 检查卷挂载是否正确
   - 检查容器重启后数据是否恢复

### 重置环境

```bash
# 停止所有服务
cd postgres && docker-compose down -v && cd ..
cd redis && docker-compose down -v && cd ..

# 删除数据目录
rm -rf postgres/data redis/redis/data

# 重新启动
cd postgres && docker-compose up -d && cd ..
cd redis && docker-compose up -d && cd ..
```

## 生产环境建议

1. **安全配置**
   - 修改默认密码
   - 启用 SSL 连接
   - 配置防火墙规则

2. **性能优化**
   - 调整 PostgreSQL 配置参数
   - 配置 Redis 内存策略
   - 启用连接池

3. **监控告警**
   - 配置日志收集
   - 设置性能监控
   - 配置健康检查

4. **备份策略**
   - 定期备份数据库
   - 配置自动备份脚本
   - 测试恢复流程
