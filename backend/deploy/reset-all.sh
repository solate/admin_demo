#!/bin/bash

# 重置所有服务的脚本（删除容器和数据）

echo "🗑️  重置管理后台所有服务..."

# 确认操作
read -p "⚠️  这将删除所有容器和数据，确定继续吗？(y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 操作已取消"
    exit 1
fi

# 停止并删除 PostgreSQL
echo "📦 重置 PostgreSQL..."
cd postgres
if docker-compose ps -q | grep -q .; then
    docker-compose down -v
    echo "✅ PostgreSQL 容器已删除"
else
    echo "ℹ️  PostgreSQL 容器未运行"
fi
cd ..

# 停止并删除 Redis
echo "📦 重置 Redis..."
cd redis
if docker-compose ps -q | grep -q .; then
    docker-compose down -v
    echo "✅ Redis 容器已删除"
else
    echo "ℹ️  Redis 容器未运行"
fi
cd ..

# 删除数据目录
echo "🗑️  删除数据目录..."
if [ -d "postgres/data" ]; then
    rm -rf postgres/data
    echo "✅ PostgreSQL 数据已删除"
fi

if [ -d "redis/redis/data" ]; then
    rm -rf redis/redis/data
    echo "✅ Redis 数据已删除"
fi

echo ""
echo "🎉 重置完成！"
echo "🚀 重新启动: ./start-all.sh"
