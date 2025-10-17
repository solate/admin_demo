package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Position holds the schema definition for the Position entity.
type Position struct {
	ent.Schema
}

func (Position) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_positions"},
		entsql.WithComments(true),
		schema.Comment("岗位"),
	}
}

// Fields of the Position.
func (Position) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("position_id").Unique().Comment("岗位ID"),
		field.String("name").NotEmpty().Comment("岗位名称"),
		field.String("department_id").NotEmpty().Comment("部门ID"),
		field.Text("description").Optional().Comment("岗位描述"),
	}
}

// Edges of the Position.
func (Position) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Position) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("position_id"),
	}
}
