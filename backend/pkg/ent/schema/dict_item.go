package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DictItem holds the schema definition for the DictItem entity.
type DictItem struct {
	ent.Schema
}

func (DictItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_dict_items"},
		entsql.WithComments(true),
		schema.Comment("字典数据"),
	}
}

// Fields of the DictItem.
func (DictItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("item_id").Unique().Comment("字典项ID"),
		field.String("type_code").NotEmpty().Comment("字典类型code"),
		field.String("label").NotEmpty().Comment("字典标签"),
		field.String("value").NotEmpty().Comment("字典键值"),
		field.Text("description").Optional().Comment("字典项描述"),
		field.Int("sort").Default(0).Comment("排序"),
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),
	}
}

// Edges of the DictItem.
func (DictItem) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (DictItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type_code"),
	}
}
