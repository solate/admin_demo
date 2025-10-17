# 管理系统后端实现文档

## 项目概述

基于 Go-Zero 框架开发的管理系统后端，实现工厂管理、商品管理、库存管理和数据统计功能。

## 需求分析

### 核心功能模块
1. **登录认证模块** - 用户登录、注册、权限管理
2. **工厂管理模块** - 工厂的增删改查操作
3. **商品管理模块** - 商品信息管理和库存操作
4. **数据统计模块** - 各类业务数据统计和报表

## 技术架构

### 技术栈
- **框架**: Go-Zero (微服务框架)
- **ORM**: Ent (Facebook 开源)
- **数据库**: PostgreSQL/MySQL
- **缓存**: Redis
- **认证**: JWT + Casbin

### 项目结构
```
backend/
├── app/admin/                 # 管理端应用
├── pkg/                      # 公共包
│   ├── ent/                  # Ent ORM 相关
│   ├── common/               # 公共工具
│   └── constants/            # 常量定义
├── internal/                 # 内部包
│   ├── handler/              # 接口处理器
│   ├── logic/                # 业务逻辑
│   ├── repository/           # 数据访问层
│   └── types/                # 类型定义
└── docs/                     # 文档
```

## 数据库设计

### 核心表结构

#### 1. 工厂表 (factories)
- `factory_id` - 工厂ID (主键)
- `factory_name` - 工厂名称
- `address` - 工厂地址
- `contact_phone` - 联系电话
- `status` - 状态 (1:启用, 2:禁用)

#### 2. 商品表 (products)
- `product_id` - 商品ID (主键)
- `product_name` - 商品名称
- `unit` - 单位
- `purchase_price` - 采购价格 (decimal)
- `sale_price` - 销售价格 (decimal)
- `current_stock` - 当前库存
- `min_stock` - 最小库存预警
- `status` - 状态 (1:启用, 2:禁用)
- `factory_id` - 所属工厂ID

#### 3. 库存记录表 (inventories)
- `inventory_id` - 库存记录ID (主键)
- `product_id` - 商品ID
- `operation_type` - 操作类型 (in:入库, out:出库)
- `quantity` - 操作数量
- `unit_price` - 单价 (decimal)
- `total_amount` - 总金额 (decimal)
- `operator_id` - 操作人ID
- `remark` - 备注
- `operation_time` - 操作时间
- `before_stock` - 操作前库存
- `after_stock` - 操作后库存

#### 4. 商品统计表 (product_statistics)
- `statistics_id` - 统计ID (主键)
- `statistics_date` - 统计日期
- `statistics_type` - 统计类型 (daily/monthly/yearly)
- 各类统计数据字段...

## API 接口设计

### 工厂管理 API
```
POST   /api/factory/create     - 创建工厂
GET    /api/factory/list       - 工厂列表
GET    /api/factory/detail     - 工厂详情
PUT    /api/factory/update     - 更新工厂
DELETE /api/factory/delete     - 删除工厂
```

### 商品管理 API
```
POST   /api/product/create     - 创建商品
GET    /api/product/list       - 商品列表
GET    /api/product/detail     - 商品详情
PUT    /api/product/update     - 更新商品
DELETE /api/product/delete     - 删除商品
```

### 库存管理 API
```
POST   /api/inventory/in       - 商品入库
POST   /api/inventory/out      - 商品出库
GET    /api/inventory/list     - 库存记录列表
GET    /api/inventory/stock    - 商品库存信息
GET    /api/inventory/history  - 库存操作历史
```

### 数据统计 API
```
GET    /api/statistics/daily     - 日统计
GET    /api/statistics/monthly   - 月统计
GET    /api/statistics/yearly    - 年统计
POST   /api/statistics/calculate - 手动计算统计
GET    /api/statistics/overview  - 统计概览
```

## 实现计划

### 第一阶段：工厂管理
- [x] Schema 设计
- [ ] 生成 Ent 代码
- [ ] 实现 CRUD 接口
- [ ] 添加搜索和分页功能

### 第二阶段：商品管理
- [x] Schema 设计
- [ ] 生成 Ent 代码
- [ ] 实现商品管理接口
- [ ] 实现价格计算逻辑

### 第三阶段：库存管理
- [x] Schema 设计
- [ ] 生成 Ent 代码
- [ ] 实现库存操作接口
- [ ] 实现原子性操作

### 第四阶段：数据统计
- [x] Schema 设计
- [ ] 生成 Ent 代码
- [ ] 实现统计计算逻辑
- [ ] 实现统计接口

## 开发指南

### 1. 环境准备
```bash
# 安装依赖
go mod tidy

# 生成 Ent 代码
go run -mod=mod entgo.io/ent/cmd/ent generate --target ./pkg/ent/generated --feature sql/upsert,sql/versioned-migration,sql/modifier ./pkg/ent/schema

# 生成 API 代码
goctl api go -api admin.api -dir .. -style=go_zero
```

### 2. 运行项目
```bash
# 启动服务
go run admin.go -f etc/admin.yaml
```

### 3. 测试接口
- 使用 Postman 导入 API 文档
- 测试各个模块的 CRUD 操作
- 验证业务逻辑正确性

## 注意事项

1. **价格精度**: 所有价格相关字段使用 decimal(18,4) 类型
2. **租户隔离**: 所有数据操作都需要考虑租户隔离
3. **软删除**: 删除操作使用软删除，保留数据历史
4. **事务处理**: 库存操作需要使用数据库事务确保一致性
5. **权限控制**: 所有接口都需要进行权限验证
6. **错误处理**: 统一的错误处理和响应格式

## 文档结构

- `schema实现总结.md` - Schema 设计总结
- `功能需求分析.md` - 功能需求详细分析
- `接口实现计划.md` - 接口实现计划和优先级
- `提示词模板.md` - 各模块实现的提示词模板
