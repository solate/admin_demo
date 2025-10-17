# Schema 实现总结

## 已实现的 Schema

### 1. Factory (工厂管理)
**文件**: `pkg/ent/schema/factory.go`
**表名**: `factories`

**字段**:
- `factory_id` (string, unique) - 工厂ID
- `factory_name` (string, not empty) - 工厂名称
- `address` (string) - 工厂地址
- `contact_phone` (string) - 联系电话
- `status` (int, default: 1) - 状态: 1:启用, 2:禁用
- 基础字段: `created_at`, `updated_at`, `deleted_at`, `tenant_code`

**索引**: `factory_name`, `tenant_code`

### 2. Product (商品管理)
**文件**: `pkg/ent/schema/product.go`
**表名**: `products`

**字段**:
- `product_id` (string, unique) - 商品ID
- `product_name` (string, not empty) - 商品名称
- `unit` (string) - 单位
- `purchase_price` (decimal(18,4)) - 采购价格
- `sale_price` (decimal(18,4)) - 销售价格
- `current_stock` (int, default: 0) - 当前库存
- `min_stock` (int, default: 0) - 最小库存预警
- `status` (int, default: 1) - 状态: 1:启用, 2:禁用
- `factory_id` (string, optional) - 所属工厂ID
- 基础字段: `created_at`, `updated_at`, `deleted_at`, `tenant_code`

**索引**: `product_name`, `factory_id`, `tenant_code`

### 3. Inventory (库存记录)
**文件**: `pkg/ent/schema/inventory.go`
**表名**: `inventories`

**字段**:
- `inventory_id` (string, unique) - 库存记录ID
- `product_id` (string, not empty) - 商品ID
- `operation_type` (string, not empty) - 操作类型: in-入库, out-出库
- `quantity` (int) - 操作数量
- `unit_price` (decimal(18,4)) - 单价
- `total_amount` (decimal(18,4)) - 总金额
- `operator_id` (string) - 操作人ID
- `remark` (string) - 备注
- `operation_time` (int64) - 操作时间
- `before_stock` (int, default: 0) - 操作前库存
- `after_stock` (int, default: 0) - 操作后库存
- 基础字段: `created_at`, `updated_at`, `deleted_at`, `tenant_code`

**索引**: `product_id`, `operation_type`, `operator_id`, `operation_time`

### 4. ProductStatistics (商品统计)
**文件**: `pkg/ent/schema/product_statistics.go`
**表名**: `product_statistics`

**字段**:
- `statistics_id` (string, unique) - 统计ID
- `statistics_date` (string) - 统计日期 YYYY-MM-DD
- `statistics_type` (string) - 统计类型: daily-日统计, monthly-月统计, yearly-年统计
- `total_products` (int, default: 0) - 商品总数
- `active_products` (int, default: 0) - 启用商品数
- `total_stock` (int, default: 0) - 总库存数量
- `total_stock_value` (decimal(18,4)) - 总库存价值
- `low_stock_products` (int, default: 0) - 低库存商品数
- `total_in_quantity` (int, default: 0) - 总入库数量
- `total_in_amount` (decimal(18,4)) - 总入库金额
- `total_out_quantity` (int, default: 0) - 总出库数量
- `total_out_amount` (decimal(18,4)) - 总出库金额
- `total_sales_amount` (decimal(18,4)) - 总销售金额
- `total_sales_quantity` (int, default: 0) - 总销售数量
- 基础字段: `created_at`, `updated_at`, `tenant_code`

**索引**: `statistics_date`, `statistics_type`, `tenant_code`

## 设计特点

1. **统一结构**: 所有schema都包含基础字段（时间戳、软删除、租户隔离）
2. **价格精度**: 使用decimal(18,4)确保财务计算精确性
3. **租户隔离**: 支持多租户架构
4. **软删除**: 支持数据恢复
5. **索引优化**: 为常用查询字段添加索引
6. **字段简化**: 只保留核心业务字段，提高维护性
