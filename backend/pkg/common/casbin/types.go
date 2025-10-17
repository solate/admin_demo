package casbin

// 权限类型常量
const (
	PermTypeDir    = "dir"    // 目录权限
	PermTypeMenu   = "menu"   // 菜单权限
	PermTypeButton = "button" // 按钮权限
	PermTypeAPI    = "api"    // 接口权限
	PermTypeData   = "data"   // 数据权限
)

// 操作类型常量
const (
	ActionView   = "view"   // 查看
	ActionAdd    = "add"    // 增加
	ActionEdit   = "edit"   // 编辑
	ActionDelete = "delete" // 删除
	ActionExport = "export" // 导出
	ActionImport = "import" // 导入
)

// DataRule 数据权限规则
const (
	DataRuleSelf          = "self"           // 仅本人数据
	DataRuleDepartment    = "department"     // 本部门数据
	DataRuleDepartmentSub = "department_sub" // 本部门及下级部门
	DataRuleAll           = "all"            // 所有数据
)

// ResourceType 资源类型
type ResourceType struct {
	Type      string
	Actions   []string
	DataRules []string
}

// 预定义资源类型
var ResourceTypes = map[string]ResourceType{
	PermTypeDir: {
		Type:    PermTypeDir,
		Actions: []string{ActionView},
	},
	PermTypeMenu: {
		Type:    PermTypeMenu,
		Actions: []string{ActionView},
	},
	PermTypeButton: {
		Type:    PermTypeButton,
		Actions: []string{ActionView, ActionAdd, ActionEdit, ActionDelete, ActionExport, ActionImport},
	},
	PermTypeAPI: {
		Type:    PermTypeAPI,
		Actions: []string{"GET", "POST", "PUT", "DELETE"},
	},
	PermTypeData: {
		Type:    PermTypeData,
		Actions: []string{ActionView, ActionAdd, ActionEdit, ActionDelete},
		DataRules: []string{
			DataRuleSelf,          // 仅本人数据
			DataRuleDepartment,    // 本部门数据
			DataRuleDepartmentSub, // 本部门及下级部门
			DataRuleAll,           // 所有数据
		},
	},
}

// Permission 权限定义
type Permission struct {
	Resource string
	Action   string
	Type     string // 添加权限类型字段
}

// DataPermissionRule 数据权限规则
type DataPermissionRule struct {
	Resource   string   // 资源名称
	Rule       string   // 权限规则（如：dept_id in (1,2,3)）
	Conditions []string // 条件
	Operations []string // 操作
}
