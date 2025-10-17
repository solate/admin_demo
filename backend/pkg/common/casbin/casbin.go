package casbin

import (
	"admin_backend/pkg/ent"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

// initCasbin 初始化Casbin配置
func NewCasbin(db *ent.Client) (*casbin.Enforcer, error) {
	a, err := NewAdapterWithClient(db)
	if err != nil {
		return nil, err
	}

	// 定义RBAC模型
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, dom, obj, act, type

[policy_definition]
p = sub, dom, obj, act, type

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act && r.type == p.type || r.act == "*" || r.sub == "admin"
	`)
	if err != nil {
		return nil, err
	}

	// 创建enforcer实例
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return enforcer, nil
}
