package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Task 任务实例
type Task struct {
	ent.Schema
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_tasks"},
		entsql.WithComments(true),
		schema.Comment("任务"),
	}
}

func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").Comment("租户code"),

		// taskID
		field.String("task_id").Unique().Comment("任务ID"),
		field.String("name").Comment("任务名称"),
		field.String("plan_id").Comment("计划ID"),
		field.String("plan_type").Comment("计划类型: 例行 routine/特殊special"),
		field.String("group").Default("default").Comment("任务分组"),
		field.Int("priority").Default(0).Comment("任务优先级"),
		field.String("status").Comment("任务状态: pending/running/success/failed/stop/interrupt"),
		field.Int64("planned_time").Comment("计划执行时间"),
		field.Int64("start_time").Optional().Comment("实际开始时间"),
		field.Int64("end_time").Optional().Comment("实际结束时间"),
		field.Int("duration").Optional().Comment("执行时长(ms)"),
		field.Text("result").Optional().Comment("执行结果"),
		field.Text("error").Optional().Comment("错误信息"),
		field.Int("retry_count").Default(0).Comment("已重试次数"),
		field.Int64("next_retry_time").Optional().Comment("下次重试时间"),
	}
}

func (Task) Edges() []ent.Edge {
	return []ent.Edge{}
}
