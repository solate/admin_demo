#!/bin/bash

# 停止所有服务的脚本

echo "🛑 停止管理后台所有服务..."

# 停止 PostgreSQL
echo "📦 停止 PostgreSQL..."
cd postgres
if docker-compose ps -q | grep -q .; then
    docker-compose down
    echo "✅ PostgreSQL 已停止"
else
    echo "ℹ️  PostgreSQL 未运行"
fi
cd ..

# 停止 Redis
echo "📦 停止 Redis..."
cd redis
if docker-compose ps -q | grep -q .; then
    docker-compose down
    echo "✅ Redis 已停止"
else
    echo "ℹ️  Redis 未运行"
fi
cd ..

echo ""
echo "🎉 所有服务已停止！"
echo "📊 查看容器状态: docker ps -a"
echo "🗑️  删除容器和数据: ./reset-all.sh"
