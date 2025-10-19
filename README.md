# admin_demo

## 📋 项目介绍

这是一个完整的后端管理系统演示项目，包含前后端分离架构、数据库、缓存等完整技术栈。

## 🎯 功能需求

### 登录页面
- ✅ 页面名称、账户、密码、查看密码按钮、登录、注册
- ✅ 验证码功能
- ✅ 密码加密存储

### 首页（Dashboard）
- ✅ 导航栏和侧边栏
- ✅ 核心数据统计卡片
- ✅ 快捷入口导航

### 工厂管理页面
- ✅ 新建工厂
- ✅ 搜索工厂
- ✅ 编辑工厂信息
- ✅ 删除工厂
- ✅ 刷新数据

### 商品管理页面
- ✅ 新建商品（基础信息及价格）
- ✅ 商品入库操作
- ✅ 商品出库操作
- ✅ 展示库存数量
- ✅ 查看操作历史记录
- ✅ 低库存预警

### 数据统计页面（Statistics）
- ✅ 商品总数
- ✅ 总库存数量
- ✅ 总库存价值（按采购价）
- ✅ 总销售价值（按销售价）
- ✅ 低库存商品数
- ✅ 总入库数量和金额
- ✅ 总出库数量和金额
- ✅ 商品明细表格

## 📸 功能截图

### 1. 登录页面
![Login Page](https://github.com/solate/admin_demo/raw/main/frontend/docs/images/1.png)

### 2. 首页（数据统计总览）
![Dashboard](https://github.com/solate/admin_demo/raw/main/frontend/docs/images/2.png)

### 3. 工厂管理
![Factory Management](https://github.com/solate/admin_demo/raw/main/frontend/docs/images/3.png)

### 4. 商品管理
![Product Management]([./docs/images/4.png](https://github.com/solate/admin_demo/raw/main/frontend/docs/images/4.png))

## 🛠 技术栈

### 后端
- **框架**: Go-Zero (go-zero framework)
- **数据库**: PostgreSQL
- **缓存**: Redis
- **ORM**: Ent
- **认证**: JWT
- **权限管理**: Casbin

### 前端
- **框架**: Vue.js 3
- **构建工具**: Vite
- **UI 组件库**: Element Plus
- **HTTP客户端**: Axios
- **路由**: Vue Router
- **状态管理**: Pinia (可选)

## 🚀 快速开始

### 后端启动

```bash
cd backend/app/admin
go run admin.go -f etc/admin.yaml
```

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

### 数据库初始化

```bash
cd backend/cmd/init_db
go run init_db.go
```

## 📊 系统架构

- **多租户支持**: 租户隔离，数据独立管理
- **权限控制**: 基于角色的权限管理（RBAC）
- **日志记录**: 操作日志和登录日志
- **缓存机制**: Redis缓存优化性能
- **异常处理**: 统一的错误处理和响应格式

## 👤 默认用户

- **用户名**: admin
- **密码**: admin@123
- **权限**: 超级管理员

## 📝 主要特性

- ✅ 前后端分离
- ✅ RESTful API 设计
- ✅ 多租户架构
- ✅ 权限管理
- ✅ 数据统计和分析
- ✅ 库存管理
- ✅ 低库存预警
- ✅ 操作历史记录
- ✅ 美观的现代化UI
- ✅ 完整的错误处理

## 📄 文档

- [需求.md](./docs/需求.md) - 项目需求
- [接口实现计划.md](./docs/接口实现计划.md) - 接口实现计划
- [schema实现总结.md](./docs/schema实现总结.md) - 数据库Schema总结
- [功能需求分析.md](./docs/功能需求分析.md) - 功能需求分析

## 📞 联系方式

如有问题或建议，请提交Issue或Pull Request。



