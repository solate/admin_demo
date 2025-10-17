

1. 菜单表(menus) ：
```sql
id  | name     | menu_code        | parent_id | path          | component
1   | 系统管理 | sys_manage       | 0         | /system       | Layout
2   | 用户管理 | sys_user        | 1         | /system/user  | system/user/index
3   | 角色管理 | sys_role        | 1         | /system/role  | system/role/index
 ```

2. 权限表(permissions) ：
```sql
id  | name       | resource_code    | action | type   | menu_id
1   | 查看用户   | sys_user        | view   | menu   | 2
2   | 新增用户   | sys_user_add    | action | button | 2
3   | 编辑用户   | sys_user_edit   | action | button | 2
4   | 删除用户   | sys_user_delete | action | button | 2
5   | 查看角色   | sys_role        | view   | menu   | 3
 ```

3. Casbin 规则表(casbin_rule) ：
```sql
ptype | v0    | v1      | v2         | v3     | v4     | v5
p     | admin | default | sys_user   | view   | menu   |
p     | admin | default | sys_role   | view   | menu   |
p     | admin | default | sys_user_add| action | button |
g     | alice | admin   | default       |        |        |
g     | bob   | user    | default       |        |        |
 ```

4. 角色表(roles)：
```sql
id  | name      | code      | status | description
1   | 超级管理员 | admin     | 1      | 系统超级管理员
2   | 普通用户   | user      | 1      | 普通用户
3   | 测试用户   | test      | 1      | 测试用户
```
5. 用户表(users)：
```sql
id  | username | password | name      | status
1   | admin    | ****     | 管理员     | 1
2   | alice    | ****     | 爱丽丝     | 1
3   | bob      | ****     | 鲍勃       | 1
 ```

6. 用户角色关系表(user_roles)：
```sql
id  | user_id | role_id | tenant_code
1   | 2       | 1       | default      # alice 是 admin 角色
2   | 3       | 2       | default      # bob 是 user 角色

```


8. 权限验证示例：


当用户 alice 访问系统时：

1. 验证菜单访问权限：
```go
// 验证是否可以访问用户管理菜单
enforcer.Enforce("admin", "default", "sys_user", "view", "menu")  // true

// 验证是否可以访问用户新增按钮
enforcer.Enforce("admin", "default", "sys_user_add", "action", "button")  // true



// alice(admin角色) 可以访问用户管理和新增用户按钮
enforcer.Enforce("admin", "default", "sys_user", "view", "menu")      // true
enforcer.Enforce("admin", "default", "sys_user_add", "action", "button") // true

// bob(user角色) 只能访问用户管理，不能新增用户
enforcer.Enforce("user", "default", "sys_user", "view", "menu")       // true
enforcer.Enforce("user", "default", "sys_user_add", "action", "button")  // false



// 为角色分配菜单权限
enforcer.AddPolicy("admin", "default", "sys_user", "view", "menu")

// 为角色分配按钮权限
enforcer.AddPolicy("admin", "default", "sys_user_add", "action", "button")

// 为用户分配角色
enforcer.AddGroupingPolicy("alice", "admin", "default")


// 获取用户的所有角色
roles := enforcer.GetRolesForUser("alice", "default")  // ["admin"]

// 获取角色的所有权限
permissions := enforcer.GetPermissionsForUser("admin", "default")
// 返回结果示例：
// [
//   ["admin", "default", "sys_user", "view", "menu"],
//   ["admin", "default", "sys_user_add", "action", "button"]
// ]
 ```

2. 返回的权限数据：
```json
{
    "menus": [
        {
            "id": "1",
            "name": "系统管理",
            "menuCode": "sys_manage",
            "children": [
                {
                    "id": "2",
                    "name": "用户管理",
                    "menuCode": "sys_user",
                    "buttons": ["sys_user_add", "sys_user_edit", "sys_user_delete"]
                }
            ]
        }
    ],
    "permissions": [
        "sys_user_view_menu",
        "sys_user_add_action_button",
        "sys_user_edit_action_button",
        "sys_user_delete_action_button"
    ]
}
 ```