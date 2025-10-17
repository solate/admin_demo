#!/bin/bash

# 启动所有服务的脚本

echo "🚀 启动管理后台所有服务..."

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

# 启动 PostgreSQL
echo "📦 启动 PostgreSQL..."
cd postgres
if ! docker-compose up -d; then
    echo "❌ PostgreSQL 启动失败"
    exit 1
fi
cd ..

# 启动 Redis
echo "📦 启动 Redis..."
cd redis
if ! docker-compose up -d; then
    echo "❌ Redis 启动失败"
    exit 1
fi
cd ..

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 5

# 验证服务
echo "🔍 验证服务状态..."

# 验证 PostgreSQL
if docker exec postgres_db pg_isready -U root -d testdb > /dev/null 2>&1; then
    echo "✅ PostgreSQL 连接正常"
else
    echo "❌ PostgreSQL 连接失败"
fi

# 验证 Redis
if docker exec my_redis redis-cli -a 123456 ping > /dev/null 2>&1; then
    echo "✅ Redis 连接正常"
else
    echo "❌ Redis 连接失败"
fi

echo ""
echo "🎉 所有服务启动完成！"
echo "📊 查看容器状态: docker ps"
echo "🔧 查看日志: docker logs postgres_db 或 docker logs my_redis"
echo "🛑 停止服务: ./stop-all.sh"
echo ""
echo "现在可以启动 Go 服务了："
echo "cd ../app/admin && go run admin.go -f etc/admin.yaml"
