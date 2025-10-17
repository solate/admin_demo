package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// DictType holds the schema definition for the DictType entity.
type DictType struct {
	ent.Schema
}

func (DictType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "sys_dict_types"},
		entsql.WithComments(true),
		schema.Comment("字典类型"),
	}
}

// Fields of the DictType.
func (DictType) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("created_at").Immutable().Default(0).Comment("创建时间"),
		field.Int64("updated_at").Default(0).Comment("修改时间"),
		field.Int64("deleted_at").Optional().Nillable().Comment("删除时间"),
		field.String("tenant_code").NotEmpty().Comment("租户编码"),

		field.String("type_id").Unique().Comment("字典类型ID"),
		field.String("name").NotEmpty().Comment("字典类型名称"),
		field.String("code").NotEmpty().Unique().Comment("字典类型编码"),
		field.Text("description").Optional().Comment("字典类型描述"),
		field.Int("status").Default(1).Comment("状态: 1:启用, 2:禁用"),
	}
}

// Edges of the DictType.
func (DictType) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (DictType) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("code").Unique(),
	}
}
