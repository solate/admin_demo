package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SystemLog holds the schema definition for the SystemLog entity.
type SystemLog struct {
	ent.Schema
}

func (SystemLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_system_logs"},
		entsql.WithComments(true),
		schema.Comment("系统日志"),
	}
}

// Fields of the SystemLog.
func (SystemLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		// field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("module").Default("").Comment("所属模块"),
		field.String("action").Default("").Comment("操作类型"),
		field.String("content").Default("").Comment("操作内容"),
		field.String("operator").Default("").Comment("操作人"),
		field.String("user_id").Default("").Comment("用户ID"),
	}
}

// Edges of the SystemLog.
func (SystemLog) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the SystemLog.
func (SystemLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("module"),
	}
}
