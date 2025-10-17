package casbinrulerepo

import (
	"context"
	"fmt"
	"strings"

	"admin_backend/pkg/ent"
	"admin_backend/pkg/ent/generated"
)

type CasbinRuleRepo struct {
	db *ent.Client
}

// NewCasbinRuleRepo 创建 casbin 规则仓储实例
func NewCasbinRuleRepo(db *ent.Client) *CasbinRuleRepo {
	return &CasbinRuleRepo{db: db}
}

// QueryBySQL 执行自定义 SQL 查询并返回 CasbinRule 结果集
func (r *CasbinRuleRepo) QueryBySQL(ctx context.Context, query string, args ...any) ([]*generated.CasbinRule, error) {
	// 执行 SQL 查询
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*generated.CasbinRule
	// 扫描结果到结构体切片
	for rows.Next() {
		rule := &generated.CasbinRule{}
		err := rows.Scan(
			&rule.ID,
			&rule.Ptype,
			&rule.V0,
			&rule.V1,
			&rule.V2,
			&rule.V3,
			&rule.V4,
			&rule.V5,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

// generatePlaceholders 生成 PostgreSQL 的占位符 ($1, $2, $3...)
func generatePlaceholders(count int) string {
	var placeholders []string
	for i := 1; i <= count; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}
	return strings.Join(placeholders, ",")
}

// // BuildInClauseQuery 构造带有 IN 子句的 SQL 查询和参数
// func BuildInClauseQuery(tableName, column string, values []interface{}) (string, []interface{}) {
// 	if len(values) == 0 {
// 		return "", nil // 如果没有值，返回空查询和参数
// 	}

// 	// 生成占位符 ($1, $2, $3...)
// 	placeholder := generatePlaceholders(len(values))

// 	// 构造 SQL 查询
// 	query := fmt.Sprintf("SELECT * FROM %s WHERE ptype = 'g' AND %s IN (%s);", tableName, column, placeholder)

// 	return query, values
// }

// QueryByUserID 根据用户 ID 列表查询角色规则
func (r *CasbinRuleRepo) QueryByUserID(ctx context.Context, args ...any) ([]*generated.CasbinRule, error) {
	if len(args) == 0 {
		return nil, nil // 如果没有值，返回空查询和参数
	}

	// 生成占位符 ($1, $2, $3...)
	placeholder := generatePlaceholders(len(args))
	// 构造 SQL 查询
	query := fmt.Sprintf("SELECT * FROM casbin_rules WHERE ptype = 'g' AND v0 IN (%s);", placeholder)
	fmt.Println("sql query: " + query)

	return r.QueryBySQL(ctx, query, args...)
}

// QueryByMenuCode 根据菜单编码查询权限规则
func (r *CasbinRuleRepo) QueryByMenuCode(ctx context.Context, args ...any) ([]*generated.CasbinRule, error) {
	if len(args) == 0 {
		return nil, nil // 如果没有值，返回空查询和参数
	}
	// 生成占位符 ($1, $2, $3...)
	placeholder := generatePlaceholders(len(args))

	// 构造 SQL 查询
	query := fmt.Sprintf("SELECT * FROM casbin_rules WHERE ptype = 'p' AND v2 IN (%s);", placeholder)
	fmt.Println("sql query: " + query)
	return r.QueryBySQL(ctx, query, args...)
}
