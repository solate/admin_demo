package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Department holds the schema definition for the Department entity.
type Department struct {
	ent.Schema
}

func (Department) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_departments"},
		entsql.WithComments(true),
		schema.Comment("部门"),
	}
}

// Fields of the Department.
func (Department) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("department_id").Unique().Comment("部门ID"),
		field.String("name").NotEmpty().Comment("部门名称"),
		field.String("parent_id").Default("").Comment("父部门ID"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the Department.
func (Department) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Department) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("department_id"),
	}
}
