package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_roles"},
		entsql.WithComments(true),
		schema.Comment("角色"),
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("role_id").Unique().Immutable().Comment("角色ID"),
		field.String("name").NotEmpty().Comment("角色名"),
		field.String("code").Unique().NotEmpty().Comment("角色编码"),
		field.Text("description").Optional().Comment("角色描述"),
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique(),
	}
}
