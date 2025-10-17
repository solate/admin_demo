package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tenant holds the schema definition for the Tenant entity.
type Tenant struct {
	ent.Schema
}

func (Tenant) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_tenants"},
		entsql.WithComments(true),
		schema.Comment("租户"),
	}
}

// Fields of the Tenant.
func (Tenant) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),

		field.String("tenant_id").Unique().Comment("租户ID"),
		field.String("name").NotEmpty().Default("").Comment("租户名称"),
		field.String("code").Unique().NotEmpty().Comment("租户编码"),
		field.Text("description").Default("").Comment("租户描述"),
		field.Int("status").Default(1).Comment("租户状态：1: 启用, 2: 禁用"),

		// // 配置管理
		// field.JSON("config", map[string]interface{}{}).Optional().Comment("租户配置"),
		// field.String("domain").Optional().Comment("租户域名"),
		// field.String("logo").Optional().Comment("租户logo"),
		// field.String("theme").Optional().Default("default").Comment("租户主题"),

		// // 资源限制
		// field.Int("user_limit").Optional().Default(100).Comment("用户数量限制"),
		// field.Int("storage_limit").Optional().Default(1024).Comment("存储空间限制(MB)"),
		// field.Int("api_rate_limit").Optional().Default(1000).Comment("API调用频率限制(次/分钟)"),

		// // 计费相关
		// field.Time("expire_time").Optional().Comment("过期时间"),
		// field.String("package_type").Default("free").Comment("套餐类型：free-免费版，pro-专业版，enterprise-企业版"),
	}
}

// Edges of the Tenant.
func (Tenant) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Tenant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
	}
}
