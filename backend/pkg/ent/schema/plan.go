package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Plan 计划定义
type Plan struct {
	ent.Schema
}

func (Plan) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_plans"},
		entsql.WithComments(true),
		schema.Comment("计划"),
	}
}

func (Plan) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		field.String("plan_id").Unique().Comment("计划ID"),
		field.String("name").Comment("计划名称"),
		field.Text("description").Optional().Comment("计划描述"),
		field.String("group").Default("default").Comment("任务分组"),
		field.String("cron_spec").Comment("cron表达式"),
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),
		field.String("plan_type").Comment("计划类型: 例行 routine/特殊special"),
		field.Int("priority").Default(0).Comment("任务优先级"),
		field.Int("timeout").Default(3600).Comment("任务超时时间(秒)"),
		field.Int("retry_times").Default(0).Comment("重试次数"),
		field.Int("retry_interval").Default(0).Comment("重试间隔(秒)"),
		field.Int64("start_time").Optional().Comment("生效开始时间"),
		field.Int64("end_time").Optional().Comment("生效结束时间"),
		field.String("command").Comment("要执行的命令或方法"),
		field.String("params").Optional().Comment("执行参数，支持JSON格式"),
	}
}

func (Plan) Edges() []ent.Edge {
	return []ent.Edge{}
}
