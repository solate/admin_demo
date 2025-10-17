package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserPosition holds the schema definition for the UserPosition entity.
type UserPosition struct {
	ent.Schema
}

func (UserPosition) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_user_positions_relation"},
		entsql.WithComments(true),
		schema.Comment("用户岗位关联"),
	}
}

// Fields of the UserPosition.
func (UserPosition) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id").NotEmpty().Comment("用户ID"),
		field.String("position_id").NotEmpty().Comment("岗位ID"),
	}
}

// Edges of the UserPosition.
func (UserPosition) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (UserPosition) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("position_id"),
		// 确保一个用户在同一个岗位下只有一条记录
		index.Fields("user_id", "position_id").Unique(),
	}
}
