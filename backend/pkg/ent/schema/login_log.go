package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// LoginLog holds the schema definition for the LoginLog entity.
type LoginLog struct {
	ent.Schema
}

func (LoginLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_login_log"},
	}
}

// Fields of the LoginLog.

func (LoginLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		// field.Int("updated_at").Default(0).Comment("修改时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("log_id").Unique().Immutable().Comment("日志ID"),
		field.String("user_id").Comment("用户ID"),
		field.String("user_name").NotEmpty().Comment("用户名"),
		field.String("ip").NotEmpty().Comment("IP地址"),
		// field.Int("status").Default(1).Comment("状态: 1:成功, 2:失败"),
		field.String("message").Optional().Comment("消息"),

		field.String("user_agent").Optional().Comment("用户代理"),
		field.String("browser").Optional().Comment("浏览器"),
		field.String("os").Optional().Comment("操作系统"),
		field.String("device").Optional().Comment("设备"),
		// field.String("location").Optional().Comment("位置, 归属地"),
		field.Int64("login_time").Optional().Comment("登录时间"),
	}
}

// Edges of the LoginLog.
func (LoginLog) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (LoginLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}
