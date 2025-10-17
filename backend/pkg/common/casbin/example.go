package casbin

// import (
// 	"admin_backend/pkg/common/casbin"
// )

// func Example(pm *CasbinRepo) {

// 	tenant := "tenant1"

// 	// 1. 创建角色并分配给用户
// 	pm.AddRoleForUser("user1", "admin", tenant)

// 	// 2. 设置菜单权限
// 	pm.AddMenuPermission("admin", tenant, "/dashboard")
// 	pm.AddMenuPermission("admin", tenant, "/users")

// 	// 3. 设置按钮权限
// 	pm.AddButtonPermission("admin", tenant, "/users", "create")
// 	pm.AddButtonPermission("admin", tenant, "/users", "edit")

// 	// 4. 设置数据权限
// 	dataAttrs := []DataAttribute{
// 		{Field: "department", Operator: "in", Value: "dep1,dep2"},
// 		{Field: "status", Operator: "eq", Value: "active"},
// 	}
// 	pm.AddDataPermission("admin", tenant, "users", dataAttrs)

// 	// 5. 权限检查
// 	canAccessMenu, _ := pm.CheckMenuPermission("user1", tenant, "/dashboard")
// 	canCreateUser, _ := pm.CheckButtonPermission("user1", tenant, "/users", "create")

// 	// 6. 数据权限检查
// 	data := map[string]interface{}{
// 		"department": "dep1",
// 		"status":     "active",
// 	}
// 	canAccessData, _ := pm.CheckDataPermission("user1", tenant, "users", data)

// 	// 7. 获取用户菜单和按钮权限
// 	menus, _ := pm.GetUserMenus("user1", tenant)
// 	buttons, _ := pm.GetUserButtons("user1", tenant, "/users")
// }

// // 初始化默认角色和权限
// func InitializeDefaultPermissions(pm *casbin.CasbinManager) error {
// 	// 创建角色
// 	policies := [][]string{
// 		{"admin", "domain1", "/api/*", "GET"},
// 		{"admin", "domain1", "/api/*", "POST"},
// 		{"user", "domain1", "/api/public/*", "GET"},
// 	}

// 	// 1. 创建角色
// 	role := "admin"

// 	// 2. 设置菜单权限
// 	menus := []string{"/dashboard", "/users", "/roles"}
// 	for _, menu := range menus {
// 		if err := pm.AddMenuPermission(role, tenantID, menu); err != nil {
// 			return err
// 		}
// 	}

// 	// 3. 设置按钮权限
// 	buttons := map[string][]string{
// 		"/users": {"add", "edit", "delete"},
// 		"/roles": {"add", "edit"},
// 	}
// 	for menu, btns := range buttons {
// 		for _, btn := range btns {
// 			if err := pm.AddButtonPermission(role, tenantID, menu, btn); err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil

// 	return pm.BatchAddPermissions(policies)
// }

// // func AuthMiddleware(pm *casbin.CasbinManager) gin.HandlerFunc {
// // 	return func(c *gin.Context) {
// // 		user := GetCurrentUser(c)     // 从上下文获取用户
// // 		domain := GetCurrentDomain(c) // 获取当前域

// // 		hasPermission, err := pm.CheckPermission(
// // 			user.ID,
// // 			domain,
// // 			c.Request.URL.Path,
// // 			c.Request.Method,
// // 		)

// // 		if err != nil || !hasPermission {
// // 			c.AbortWithStatus(403)
// // 			return
// // 		}

// // 		c.Next()
// // 	}
// // }

// // 更新角色权限
// func UpdateRolePermissions(pm *casbin.CasbinManager, role, domain string, permissions []Permission) error {
// 	// 先清除旧权限
// 	oldPermissions, _ := pm.GetRolePermissions(role, domain)
// 	pm.BatchRemovePermissions(oldPermissions)

// 	// 添加新权限
// 	var newPolicies [][]string
// 	for _, p := range permissions {
// 		newPolicies = append(newPolicies, []string{role, domain, p.Resource, p.Action})
// 	}
// 	return pm.BatchAddPermissions(newPolicies)
// }

// /**

// # 角色关系 (g)
// ptype   v0          v1      v2
// g       user123     admin   domain1    -- 用户user123是domain1域的admin角色

// # 权限策略 (p)
// ptype   v0      v1          v2              v3
// p       admin   domain1     /api/users      GET     -- admin角色可以GET访问/api/users

// // 可以使用通配符
// pm.AddPermissionForRole("admin", "domain1", "/api/*", "GET")  // 所有 API 的读权限
// pm.AddPermissionForRole("editor", "domain1", "/api/articles/*", "POST")  // 仅文章相关的写权限

// // 批量操作比单个操作效率更高
// policies := [][]string{
//     {"role1", "domain1", "/api/users", "GET"},
//     {"role1", "domain1", "/api/users", "POST"},
// }
// pm.BatchAddPermissions(policies)
// */

// /*
// // 为用户分配角色
// pm.AddRoleForUser("alice", "admin", "domain1")

// // 为角色设置权限
// pm.AddPermissionForRole("admin", "domain1", "/api/users", "GET")
// pm.AddPermissionForRole("admin", "domain1", "/api/users", "POST")

// // 检查权限
// hasPermission, _ := pm.CheckPermission("alice", "domain1", "/api/users", "GET")
// fmt.Printf("Has permission: %v\n", hasPermission)
// */

// /**

// CREATE TABLE menus (
//     id VARCHAR(36) PRIMARY KEY,
//     name VARCHAR(100) NOT NULL,
//     path VARCHAR(255) NOT NULL,
//     parent_id VARCHAR(36),
//     icon VARCHAR(50),
//     sort_order INT,
//     tenant_id VARCHAR(36),
//     created_at TIMESTAMP,
//     updated_at TIMESTAMP
// );

// CREATE TABLE buttons (
//     id VARCHAR(36) PRIMARY KEY,
//     menu_id VARCHAR(36) NOT NULL,
//     code VARCHAR(50) NOT NULL,
//     name VARCHAR(100) NOT NULL,
//     tenant_id VARCHAR(36) NOT NULL,
//     created_at TIMESTAMP,
//     updated_at TIMESTAMP
// );

// ptype | v0    | v1      | v2        | v3     | v4 | v5
// g     | user1 | admin   | tenant1   |        |    |    // 用户角色关系
// p     | admin | tenant1 | /api/user | GET    |    |    // 权限策略
// p     | admin | tenant1 | btn:add   | access |    |    // 按钮权限

// func InitializePermissions(pm *CasbinManager, tenantID string) error {
//     // 1. 创建角色
//     role := "admin"

//     // 2. 设置菜单权限
//     menus := []string{"/dashboard", "/users", "/roles"}
//     for _, menu := range menus {
//         if err := pm.AddMenuPermission(role, tenantID, menu); err != nil {
//             return err
//         }
//     }

//     // 3. 设置按钮权限
//     buttons := map[string][]string{
//         "/users": {"add", "edit", "delete"},
//         "/roles": {"add", "edit"},
//     }
//     for menu, btns := range buttons {
//         for _, btn := range btns {
//             if err := pm.AddButtonPermission(role, tenantID, menu, btn); err != nil {
//                 return err
//             }
//         }
//     }

//     return nil
// }

// ptype | v0         | v1        | v2                    | v3
// g     | user123    | admin     | tenant1              |
// p     | admin      | tenant1   | menu:/users          | access
// p     | admin      | tenant1   | button:/users:add    | access
// p     | admin      | tenant1   | button:/users:edit   | access
// p     | admin      | tenant1   | button:/users:delete | access
// */

// /**
// -- 权限类型枚举表
// CREATE TABLE permission_types (
//     id BIGINT PRIMARY KEY AUTO_INCREMENT,
//     name VARCHAR(50) NOT NULL,
//     code VARCHAR(50) NOT NULL UNIQUE,
//     description VARCHAR(200),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// -- 菜单表
// CREATE TABLE menus (
//     id BIGINT PRIMARY KEY AUTO_INCREMENT,
//     parent_id BIGINT,
//     name VARCHAR(100) NOT NULL,
//     path VARCHAR(200),
//     icon VARCHAR(100),
//     sort_order INT DEFAULT 0,
//     status TINYINT DEFAULT 1,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//     FOREIGN KEY (parent_id) REFERENCES menus(id)
// );

// -- 修改权限表，增加更多字段
// CREATE TABLE permissions (
//     id BIGINT PRIMARY KEY AUTO_INCREMENT,
//     name VARCHAR(100) NOT NULL,
//     code VARCHAR(100) NOT NULL UNIQUE,
//     type_id BIGINT NOT NULL,  -- 关联权限类型
//     resource_type VARCHAR(50) NOT NULL, -- MENU/PAGE/BUTTON/API/DATA
//     resource_id VARCHAR(100),  -- 关联的资源ID
//     action VARCHAR(50),        -- 操作类型：VIEW/ADD/EDIT/DELETE等
//     description VARCHAR(200),
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//     FOREIGN KEY (type_id) REFERENCES permission_types(id)
// );

// -- 数据权限规则表
// CREATE TABLE data_permission_rules (
//     id BIGINT PRIMARY KEY AUTO_INCREMENT,
//     permission_id BIGINT NOT NULL,
//     rule_type VARCHAR(50) NOT NULL,  -- 规则类型：DEPARTMENT/USER/CUSTOM
//     rule_value TEXT NOT NULL,        -- 规则值，可以是JSON格式
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//     FOREIGN KEY (permission_id) REFERENCES permissions(id)
// );

// package service

// type MenuService struct {
//     // 依赖注入
// }

// // GetUserMenus 获取用户有权限的菜单
// func (s *MenuService) GetUserMenus(tenantID, userID int64) ([]*Menu, error) {
//     // 1. 获取用户角色
//     // 2. 获取角色关联的菜单权限
//     // 3. 构建菜单树
//     return menus, nil
// }

// package service

// type DataPermissionService struct {
//     // 依赖注入
// }

// // BuildDataPermissionFilter 构建数据权限过滤条件
// func (s *DataPermissionService) BuildDataPermissionFilter(ctx context.Context, userID int64, resourceType string) (*DataFilter, error) {
//     // 1. 获取用户数据权限规则
//     // 2. 解析规则生成过滤条件
//     // 3. 返回过滤器
//     return filter, nil
// }

// package middleware

// type PermissionMiddleware struct {
//     // 依赖注入
// }

// // CheckPermission 综合权限检查中间件
// func (m *PermissionMiddleware) CheckPermission() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         // 1. 获取当前用户信息
//         tenantID := c.GetInt64("tenant_id")
//         userID := c.GetInt64("user_id")

//         // 2. 检查API权限
//         if !m.checkAPIPermission(c) {
//             c.JSON(403, gin.H{"error": "没有接口访问权限"})
//             c.Abort()
//             return
//         }

//         // 3. 检查数据权限
//         filter, err := m.buildDataFilter(c)
//         if err != nil {
//             c.JSON(403, gin.H{"error": "数据权限验证失败"})
//             c.Abort()
//             return
//         }

//         // 4. 将数据权限过滤器注入上下文
//         c.Set("data_filter", filter)

//         c.Next()
//     }
// }

// package middleware

// import (
//     "github.com/gin-gonic/gin"
// )

// type PermissionMiddleware struct {
//     casbinService *service.CasbinService
// }

// // CheckAPIPermission API权限中间件
// func (m *PermissionMiddleware) CheckAPIPermission() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         userID := c.GetString("user_id")
//         tenantID := c.GetString("tenant_id")
//         path := c.Request.URL.Path
//         method := c.Request.Method

//         // 检查API权限
//         hasPermission, err := m.casbinService.CheckPermission(userID, tenantID, path, method, "api")
//         if err != nil || !hasPermission {
//             c.JSON(403, gin.H{"error": "没有接口访问权限"})
//             c.Abort()
//             return
//         }
//         c.Next()
//     }
// }

// // CheckDataPermission 数据权限中间件
// func (m *PermissionMiddleware) CheckDataPermission() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         userID := c.GetString("user_id")
//         tenantID := c.GetString("tenant_id")
//         resourceCode := c.GetString("resource_code")

//         // 获取数据权限规则
//         rules, err := m.casbinService.GetDataPermissionRules(userID, tenantID, resourceCode)
//         if err != nil {
//             c.JSON(403, gin.H{"error": "数据权限获取失败"})
//             c.Abort()
//             return
//         }

//         // 将数据权限规则注入上下文
//         c.Set("data_rules", rules)
//         c.Next()
//     }
// }
// */
