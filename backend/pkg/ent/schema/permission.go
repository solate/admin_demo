package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

func (Permission) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_permissions"},
		entsql.WithComments(true),
		schema.Comment("权限"),
	}
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		// permission_id
		field.String("permission_id").Unique().Immutable().Comment("权限ID"),
		field.String("name").NotEmpty().Comment("权限名称"),
		field.String("code").Unique().NotEmpty().Comment("权限编码"),
		field.String("type").Comment("类型类型: dir/menu/button/api/data"),
		field.String("resource").NotEmpty().Comment("资源"),
		field.String("action").NotEmpty().Comment("操作类型"),
		field.String("parent_id").Optional().Comment("父级ID"),
		field.Text("description").Optional().Comment("描述"),
		field.Int("status").Default(1).Comment("状态 1:启用 2:禁用"),
		field.String("menu_id").Optional().Comment("菜单ID"),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{}
}
